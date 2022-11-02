package dagflow

import (
	"context"
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/xxf098/dagflow/plan/task"
)

// run a action in cue
// TODO: run action with dependency
func Do(filePath string, actionName string) (*cue.Value, error) {

	ctx := cuecontext.New()
	entrypoints := []string{filePath}

	bis := load.Instances(entrypoints, nil)

	var output *cue.Value

	var err error
	selectors := []cue.Selector{cue.Str("actions")}
	for _, v := range strings.Split(actionName, ".") {
		selectors = append(selectors, cue.Str(v))
	}

	for _, bi := range bis {

		if bi.Err != nil {
			fmt.Println("Error during load:", bi.Err)
			err = bi.Err
			continue
		}

		value := ctx.BuildInstance(bi)
		if value.Err() != nil {
			fmt.Println("Error during build:", value.Err())
			err = value.Err()
			continue
		}

		p := cue.MakePath(selectors...)
		a := value.LookupPath(p)
		if !a.Exists() {
			err = fmt.Errorf("path not found")
			continue
		}
		taskType, actionValue := task.LookupAction(&a)
		if len(taskType) < 1 {
			err = fmt.Errorf("task not found")
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
	return output, err
}
