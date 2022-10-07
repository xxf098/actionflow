package main

import (
	"errors"
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
)

var ErrNotTask = errors.New("not a task")

// cue.mod in the root of go project
// import mod inside pkg
func loadCue() {
	ctx := cuecontext.New()
	entrypoints := []string{"./main.cue"}

	bis := load.Instances(entrypoints, nil)

	for _, bi := range bis {

		if bi.Err != nil {
			fmt.Println("Error during load:", bi.Err)
			continue
		}

		value := ctx.BuildInstance(bi)
		if value.Err() != nil {
			fmt.Println("Error during build:", value.Err())
			continue
		}

		p := cue.MakePath(cue.Str("AddHello"), cue.Str("write"))
		w := value.LookupPath(p)
		if w.Exists() {
			if t, err := lookupType(&w); err == nil {
				fmt.Println("type: ", t)
			}
		}

		fmt.Println("kind: ", value.Kind())
		fmt.Printf("main value: %v", value)
	}
}

func lookupType(v *cue.Value) (string, error) {

	typePath := cue.MakePath(
		cue.Str("$dagger"),
		cue.Str("task"),
		cue.Hid("_name", "github.com/xxf098/dagflow"))
	corePath := cue.MakePath(
		cue.Str("$dagger"),
		cue.Str("task"),
		cue.Hid("_name", "github.com/xxf098/dagflow/core"))

	paths := []cue.Path{corePath, typePath}
	for _, path := range paths {
		typ := v.LookupPath(path)
		if typ.Exists() {
			return typ.String()
		}
	}
	return "", ErrNotTask
}
