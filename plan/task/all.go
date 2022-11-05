package task

import (
	"context"
	"log"
	"sync"

	"cuelang.org/go/cue"
	"github.com/xxf098/dagflow/compiler"
)

func init() {
	Register("All", func() Task { return &allTasks{} })
}

type allTasks struct {
}

// TODO: max concurrency, exit on error
func (t *allTasks) Run(ctx context.Context, v *cue.Value) (*cue.Value, error) {
	ignoreError, err := v.Lookup("ignore_error").Bool()
	if err != nil {
		return nil, err
	}
	iter, err := v.Lookup("tasks").List()
	if err != nil {
		return nil, err
	}
	tasks := []cue.Value{}
	for iter.Next() {
		tasks = append(tasks, iter.Value())
	}
	var wg sync.WaitGroup
	var taskErr error
	// ignore error anyway
	for i, task := range tasks {
		wg.Add(1)
		go func(ctx context.Context, index int, v cue.Value) {
			defer wg.Done()
			task, err := Lookup(&v)
			if err != nil {
				taskErr = err
				log.Println("index:", index, err)
				return
			}
			_, err = task.Run(ctx, &v)
			if err != nil {
				taskErr = err
				log.Println(task.Name(), "index:", index, err)
			}
		}(ctx, i, task)
	}
	wg.Wait()
	value := compiler.NewValue()
	output := value.FillPath(cue.ParsePath("output"), "")
	if ignoreError {
		taskErr = nil
	}
	return &output, taskErr
}

func (t *allTasks) Name() string {
	return "All"
}
