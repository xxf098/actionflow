package main

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
)

// run a action in cue
// find action dependency
func Do(filePath string, actionName string) {

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

		p := cue.MakePath(cue.Str("actions"), cue.Str(actionName))
		a := value.LookupPath(p)
		if !a.Exists() {
			continue
		}
		taskType := lookupAction(&a)
		if len(taskType) < 1 {
			continue
		}
		fmt.Println(taskType)

	}

}

func lookupAction(v *cue.Value) string {
	for iter, _ := v.Fields(cue.Optional(true)); iter.Next(); {
		vn := iter.Value()
		ik := vn.IncompleteKind()
		if ik.IsAnyOf(cue.StructKind) && v.IsConcrete() {
			t, err := lookupType(&vn)
			if err != nil {
				continue
			}
			return t
		}
	}
	return ""
}
