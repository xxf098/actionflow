package task

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
)

var (
	ErrNotTask = errors.New("not a task")
	tasks      sync.Map
	typePath   = cue.MakePath(
		cue.Str("$dagger"),
		cue.Str("task"),
		cue.Hid("_name", "github.com/xxf098/dagflow"))
	corePath = cue.MakePath(
		cue.Str("$dagger"),
		cue.Str("task"),
		cue.Hid("_name", "github.com/xxf098/dagflow/core"))
	paths = []cue.Path{corePath, typePath}
)

// return result
type Task interface {
	Run(ctx context.Context, v *cue.Value) (*cue.Value, error)
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

	typeString, err := lookupType(v)
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
			// fmt.Println(v.Cue())
			return typ.String()
		}
	}
	return "", ErrNotTask
}

func NewValue() (*cue.Value, error) {
	c := cuecontext.New()
	v := c.CompileString("", cue.Filename("_"))
	if v.Err() != nil {
		return nil, v.Err()
	}
	return &v, nil
}
