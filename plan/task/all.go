package task

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"cuelang.org/go/cue"
	"github.com/xxf098/actionflow/compiler"
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
	// max
	max, err := v.Lookup("max").Int64()
	if err != nil {
		return nil, err
	}
	lenTasks := len(tasks)
	if max == 0 || max > int64(lenTasks) {
		max = int64(lenTasks)
	}
	ch := make(chan struct{}, max)
	// var wg sync.WaitGroup

	// ignore error anyway
	lg := log.Ctx(ctx)
	start := time.Now()
	// timeout
	errChan := make(chan error, lenTasks)
	for i, task := range tasks {
		// wg.Add(1)
		ch <- struct{}{}
		go func(ctx context.Context, index int, v cue.Value) {
			// defer wg.Done()
			var err error
			defer func() {
				<-ch
				errChan <- err
			}()
			t, err := Lookup(&v)
			if err != nil {
				lg.Error().Err(fmt.Errorf("Lookup error: %s", v.Path().String())).Msgf("name: %s", t.Name())
				return
			}
			_, err = t.Run(ctx, &v)
			if err != nil {
				lg.Error().Err(fmt.Errorf("Run error: %s", v.Path().String())).Msgf("name: %s", t.Name())
			}
		}(ctx, i, task)
	}
	var taskErr error
	for i := 0; i < lenTasks; i++ {
		taskErr = <-errChan
		if !ignoreError && taskErr != nil {
			break
		}
	}
	// wg.Wait()
	lg.Info().Dur("duration", time.Since(start)).Str("task", v.Path().String()).Msg(t.Name())
	Then(ctx, v)
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
