package dagflow

import (
	"context"
	"fmt"
	"sync"
	"time"

	"cuelang.org/go/cue"
	cueflow "cuelang.org/go/tools/flow"
	"github.com/rs/zerolog/log"
	"github.com/xxf098/dagflow/compiler"
	"github.com/xxf098/dagflow/plan/task"
)

// initTasks() {
// addTask(t *cueflow.Task) {

type Runner struct {
	target cue.Path
	tasks  sync.Map
	mirror cue.Value
	l      sync.Mutex
}

func NewRunner(target cue.Path) *Runner {
	return &Runner{
		target: target,
		mirror: *compiler.NewValue(),
	}
}

// context
func (r *Runner) Run(ctx context.Context, src cue.Value) error {
	if !src.LookupPath(r.target).Exists() {
		return fmt.Errorf("%s not found", r.target.String())
	}
	if err := r.update(cue.MakePath(), &src); err != nil {
		return err
	}
	// r.computed
	flow := cueflow.New(
		&cueflow.Config{
			FindHiddenTasks: true,
		},
		src,
		r.taskFunc,
	)

	if err := flow.Run(ctx); err != nil {
		return err
	}
	return nil
}

func (r *Runner) update(p cue.Path, v *cue.Value) error {
	r.l.Lock()
	defer r.l.Unlock()

	value := r.mirror.FillPath(p, v)
	if value.Value().Err() != nil {
		return value.Value().Err()
	}
	r.initTasks(v)
	return nil
}

func (r *Runner) initTasks(v *cue.Value) {
	flow := cueflow.New(
		&cueflow.Config{
			FindHiddenTasks: true,
		},
		v,
		noOpRunner,
	)

	// Allow tasks under the target
	for _, t := range flow.Tasks() {
		if cuePathHasPrefix(t.Path(), r.target) {
			r.addTask(t)
		}
	}

}

func (r *Runner) addTask(t *cueflow.Task) {
	if _, ok := r.tasks.Load(t.Path().String()); ok {
		return
	}

	r.tasks.Store(t.Path().String(), struct{}{})
	for _, dep := range t.Dependencies() {
		r.addTask(dep)
	}
}

func (r *Runner) shouldRun(p cue.Path) bool {
	_, ok := r.tasks.Load(p.String())
	return ok
}

func (r *Runner) taskFunc(v cue.Value) (cueflow.Runner, error) {

	handler, err := task.Lookup(&v)
	if err != nil {
		if err == task.ErrNotTask {
			return nil, nil
		}
		return nil, err
	}
	if !r.shouldRun(v.Path()) {
		return nil, nil
	}
	return cueflow.RunnerFunc(func(t *cueflow.Task) error {
		ctx := t.Context()
		taskPath := t.Path().String()
		lg := log.Ctx(ctx).With().Str("task", taskPath).Logger()
		ctx = lg.WithContext(ctx)

		for _, dep := range t.Dependencies() {
			lg.Debug().Str("dependency", dep.Path().String()).Msg("dependency detected")
		}

		// fixme
		tval := t.Value()
		start := time.Now()
		result, err := handler.Run(ctx, &tval)
		if err != nil {
			return fmt.Errorf("%s: %s", t.Path().String(), err.Error())
		}
		lg.Info().Dur("duration", time.Since(start)).Str("state", task.StateCompleted.String()).Msg(task.StateCompleted.String())

		if !result.IsConcrete() {
			return nil
		}

		// set output
		if err := t.Fill(result); err != nil {
			lg.Error().Err(err).Msg("failed to fill task")
			return err
		}
		return nil
	}), nil
}

func cuePathHasPrefix(p cue.Path, prefix cue.Path) bool {
	pathSelectors := p.Selectors()
	prefixSelectors := prefix.Selectors()

	if len(pathSelectors) < len(prefixSelectors) {
		return false
	}

	for i, sel := range prefixSelectors {
		if pathSelectors[i] != sel {
			return false
		}
	}

	return true
}

// empty runner just do nothing
func noOpRunner(v cue.Value) (cueflow.Runner, error) {
	_, err := task.Lookup(&v)
	if err != nil {
		// Not a task
		if err == task.ErrNotTask {
			return nil, nil
		}
		return nil, err
	}

	// Return a no op runner
	return cueflow.RunnerFunc(func(t *cueflow.Task) error {
		return nil
	}), nil
}
