package main

import (
	"context"
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/xxf098/dagflow/task"
)

// run a action in cue
// TODO: run action with dependency
func Do(filePath string, actionName string) *cue.Value {

	ctx := cuecontext.New()
	entrypoints := []string{filePath}

	bis := load.Instances(entrypoints, nil)

	var output *cue.Value

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
		taskType, actionValue := task.LookupAction(&a)
		if len(taskType) < 1 {
			continue
		}
		fmt.Println(taskType)
		t := task.New(taskType)
		if t == nil {
			continue
		}
		var err error
		output, err = t.Run(context.Background(), actionValue)
		if err != nil {
			fmt.Println(err)
		}
	}
	return output
}
