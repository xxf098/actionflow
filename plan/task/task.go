package task

import (
	"context"
	"errors"
	"fmt"

	"sync"

	"cuelang.org/go/cue"
	"github.com/rs/zerolog/log"
)

var (
	ErrNotTask = errors.New("not a task")
	tasks      sync.Map
	typePath   = cue.MakePath(
		cue.Str("$dagger"),
		cue.Str("task"),
		cue.Hid("_name", "github.com/xxf098/actionflow"))
	corePath = cue.MakePath(
		cue.Str("$dagger"),
		cue.Str("task"),
		cue.Hid("_name", "github.com/xxf098/actionflow/core"))
	paths = []cue.Path{corePath, typePath}
)

// State is the state of the task.
type State int8

func (s State) String() string {
	return [...]string{"computing", "skipped", "completed", "cancelled", "failed"}[s]
}

func ParseState(s string) (State, error) {
	switch s {
	case "computing":
		return StateComputing, nil
	case "skipped":
		return StateSkipped, nil
	case "cancelled":
		return StateCanceled, nil
	case "failed":
		return StateFailed, nil
	case "completed":
		return StateCompleted, nil
	}

	return -1, fmt.Errorf("invalid state [%s]", s)
}

func (s State) CanTransition(t State) bool {
	return s <= t
}

const (
	// state order is important here since it defines the  order
	// on how states can transition only forwards
	// computing > completed > canceled > failed
	StateComputing State = iota
	StateSkipped
	StateCompleted
	StateCanceled
	StateFailed
)

// return result
type Task interface {
	Run(ctx context.Context, v *cue.Value) (*cue.Value, error)
	Name() string
}

type NewFunc func() Task

// Register a task type and initializer
func Register(typ string, f NewFunc) {
	tasks.Store(typ, f)
}

// New creates a new Task of the given type.
func New(typ string) Task {
	v, ok := tasks.Load(typ)
	if !ok {
		return nil
	}
	fn := v.(NewFunc)
	return fn()
}

// find action type by path
func Lookup(v *cue.Value) (Task, error) {
	if v.Kind() != cue.StructKind {
		return nil, ErrNotTask
	}

	typeString, err := lookupActionType(v)
	if err != nil {
		return nil, err
	}

	t := New(typeString)
	if t == nil {
		return nil, fmt.Errorf("unknown type %q", typeString)
	}

	return t, nil
}

func lookupType(v *cue.Value) (string, error) {
	for _, path := range paths {
		typ := v.LookupPath(path)
		if typ.Exists() {
			return typ.String()
		}
	}
	return "", ErrNotTask
}

// lookup action type in cue
func lookupActionType(v *cue.Value) (string, error) {
	ik := v.IncompleteKind()
	if ik.IsAnyOf(cue.StructKind) && v.IsConcrete() {
		t, err := lookupType(v)
		if err != nil {
			return "", err
		}
		return t, nil
	}
	return "", ErrNotTask
}

// lookup action type in cue
func LookupAction(v *cue.Value) (string, *cue.Value) {
	if t, err := lookupType(v); err == nil && len(t) > 0 {
		return t, v
	}
	for iter, _ := v.Fields(cue.Optional(true)); iter.Next(); {
		vn := iter.Value()
		ik := vn.IncompleteKind()
		if ik.IsAnyOf(cue.StructKind) && v.IsConcrete() {
			t, err := lookupType(&vn)
			if err != nil {
				continue
			}
			return t, &vn
		}
	}
	return "", nil
}

// do then action
func Then(ctx context.Context, v *cue.Value) error {
	tv := v.LookupPath(cue.ParsePath("then"))
	isBackground := false
	if !tv.Exists() {
		tv = v.LookupPath(cue.ParsePath("then_"))
		isBackground = true
	}
	tk := tv.Kind()
	if tk.IsAnyOf(cue.StructKind) && tv.IsConcrete() {
		return runThenTask(ctx, tv, isBackground)
	}
	// else if tv.IsConcrete() && tk.IsAnyOf(cue.ListKind) {
	// 	iter, err := tv.List()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	for iter.Next() {
	// 		v := iter.Value()
	// 		if err := runThenTask(ctx, v, isBackground); err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	return nil
}

func runThenTask(ctx context.Context, tv cue.Value, isBackground bool) error {
	lg := log.Ctx(ctx)
	task, err := Lookup(&tv)
	if err != nil {
		lg.Error().Err(err).Str("then task", tv.Path().String()).Msg(task.Name())
		return err
	} else {
		if isBackground {
			go func() {
				if _, err := task.Run(ctx, &tv); err != nil {
					lg.Error().Err(err).Str("then task", tv.Path().String()).Msg(task.Name())
				}
			}()
			return nil
		}
		_, err := task.Run(ctx, &tv)
		if err != nil {
			return err
		}
	}
	return nil
}
