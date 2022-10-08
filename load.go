package main

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/xxf098/dagflow/task"
)

const (
	ScalarKind cue.Kind = cue.StringKind | cue.NumberKind | cue.BoolKind
)

// cue.mod in the root of go project
// import mod inside pkg
func loadCue(filePath string) {
	ctx := cuecontext.New()
	entrypoints := []string{filePath}

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

		// p := cue.MakePath(cue.Str("AddHello"), cue.Str("write"))
		// w := value.LookupPath(p)
		// if w.Exists() {
		// 	if t, err := lookupType(&w); err == nil {
		// 		// use type to find task action
		// 		fmt.Println("type: ", t)
		// 	}
		// }

		p := cue.MakePath(cue.Str("actions"), cue.Str("hello"))
		a := value.LookupPath(p)
		if a.Exists() {
			inputs := lookupInput(&a)
			fmt.Println("hello: ", a)
			fmt.Println(inputs)
		}

		fmt.Println("kind: ", value.Kind())
		fmt.Printf("main value: %v", value)
	}
}

type Input struct {
	Name          string
	Type          string
	Documentation string
}

func lookupInput(v *cue.Value) []Input {
	inputs := []Input{}
	for iter, _ := v.Fields(cue.Optional(true)); iter.Next(); {
		vn := iter.Value()
		ik := vn.IncompleteKind()
		if ik.IsAnyOf(ScalarKind) && v.IsConcrete() {
			inputs = append(inputs, Input{
				Name: iter.Label(),
				Type: ik.String(),
			})
		}
	}
	return inputs
}

//  find action type
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
	return "", task.ErrNotTask
}
