package plan

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	cueflow "cuelang.org/go/tools/flow"
	"github.com/rs/zerolog/log"
	"github.com/xxf098/actionflow/plan/task"
)

// initTasks() {
// addTask(t *cueflow.Task) {

type Runner struct {
	target cue.Path
	tasks  sync.Map
	deps   sync.Map // dependency by attributes
	// mirror cue.Value
	l sync.Mutex
}

func NewRunner(target cue.Path) *Runner {
	return &Runner{
		target: target,
		// mirror: *compiler.NewValue(),
	}
}

// runSequence

// context
func (r *Runner) Run(ctx context.Context, src cue.Value) error {
	if !src.LookupPath(r.target).Exists() {
		return fmt.Errorf("%s not found", r.target.String())
	}

	if err := r.initDeps(ctx, &src); err != nil {
		return err
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

	// add deps
	r.updateDeps(flow)

	if err := flow.Run(ctx); err != nil {
		return err
	}
	return nil
}

func (r *Runner) update(p cue.Path, v *cue.Value) error {
	r.l.Lock()
	defer r.l.Unlock()

	// value := r.mirror.FillPath(p, v)
	// if value.Value().Err() != nil {
	// 	return value.Value().Err()
	// }
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

	r.updateDeps(flow)

	// Allow tasks under the target
	for _, t := range flow.Tasks() {
		if cuePathHasPrefix(t.Path(), r.target) {
			r.addTask(t)
		}
	}
	// add dep's deps

}

func (r *Runner) initDeps(ctx context.Context, v *cue.Value) error {
	flow := cueflow.New(
		&cueflow.Config{
			FindHiddenTasks: true,
		},
		v,
		r.depsRunner,
	)
	if err := flow.Run(ctx); err != nil {
		return err
	}
	return nil
}

func (r *Runner) updateDeps(flow *cueflow.Controller) {
	tasks := flow.Tasks()
	for i, t := range tasks {
		// if cuePathHasPrefix(t.Path(), r.target) {
		// add deps from attributes(@$)
		path := t.Path().String()
		if val, ok := r.deps.Load(path); ok {
			deps := val.([]string)
			for j, t1 := range tasks {
				if i == j {
					continue
				}
				for _, dep := range deps {
					if t1.Path().String() == dep {
						// add deps
						t.AddDep(path, t1)
						break
					}
				}
			}
		}
		// }
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
			lg.Trace().Str("dependency", dep.Path().String()).Msg("dependency detected")
		}

		// fixme
		tval := t.Value()
		start := time.Now()
		result, err := handler.Run(ctx, &tval)
		if err != nil {
			lg.Error().Err(err).Dur("duration", time.Since(start)).Str("state", task.StateCompleted.String()).Msg(task.StateCompleted.String())
			return fmt.Errorf("%s: %s", t.Path().String(), err.Error())
		}
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

// find attrs deps
func (r *Runner) depsRunner(v cue.Value) (cueflow.Runner, error) {
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
		tval := t.Value()
		attrs := tval.Attributes(cue.ValueAttr)
		for _, attr := range attrs {
			name := attr.Name()
			if !strings.HasPrefix(name, "$") {
				continue
			}
			depName := fmt.Sprintf("actions.%s", strings.TrimPrefix(name, "$"))
			taskPath := t.Path().String()
			// self dependency
			if taskPath == depName {
				return fmt.Errorf("self dependency found: %s", depName)
			}
			if val, ok := r.deps.Load(taskPath); ok {
				deps := val.([]string)
				// check already add
				found := false
				for _, dep := range deps {
					if dep == depName {
						found = true
						break
					}
				}
				if !found {
					deps = append(deps, depName)
					r.deps.Store(taskPath, deps)
				}
			} else {
				r.deps.Store(taskPath, []string{depName})
			}
		}
		return nil
	}), nil
}

// match against target
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

func LoadFile(filePath string) cue.Value {
	ctx := cuecontext.New()
	entrypoints := []string{filePath}

	bis := load.Instances(entrypoints, nil)
	return ctx.BuildInstance(bis[0])
}
