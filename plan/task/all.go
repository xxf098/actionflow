package task

import (
	"context"
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
	for _, task := range tasks {
		wg.Add(1)
		go func(ctx context.Context, v cue.Value) {
			defer wg.Done()
			task, err := Lookup(&v)
			if err != nil {
				taskErr = err
				return
			}
			_, taskErr = task.Run(ctx, &v)
		}(ctx, task)
	}
	wg.Wait()
	value := compiler.NewValue()
	output := value.FillPath(cue.ParsePath("output"), "")
	return &output, taskErr
}
