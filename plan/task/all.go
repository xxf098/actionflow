package task

import (
	"context"
	"sync"
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
	if max == 0 || max > int64(len(tasks)) {
		max = int64(len(tasks))
	}
	ch := make(chan struct{}, max)
	var wg sync.WaitGroup

	var taskErr error
	// ignore error anyway
	lg := log.Ctx(ctx)
	start := time.Now()
	for i, task := range tasks {
		wg.Add(1)
		ch <- struct{}{}
		go func(ctx context.Context, index int, v cue.Value) {
			defer wg.Done()
			defer func() {
				<-ch
			}()
			t, err := Lookup(&v)
			if err != nil {
				taskErr = err
				lg.Error().Err(err).Msgf("index: %d", index)
				return
			}
			_, err = t.Run(ctx, &v)
			if err != nil {
				taskErr = err
				lg.Error().Err(err).Msgf("index: %d name: %s", index, t.Name())
			}
		}(ctx, i, task)
	}
	wg.Wait()
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
