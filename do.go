package actionflow

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/xxf098/actionflow/plan"
	"github.com/xxf098/actionflow/plan/task"
)

func Do(ctx context.Context, dir string, action string) error {
	targetPath := getTargetPath([]string{action})
	daggerPlan, err := loadPlan(ctx, dir)
	if err != nil {
		return err
	}
	err = daggerPlan.Do(ctx, targetPath)
	if err != nil {
		return err
	}
	return nil
}

func flowTest(cueFile string, action string) error {

	v := plan.LoadFile(cueFile)
	iter, _ := v.Fields()
	for iter.Next() {
		fmt.Println(iter.Label())
	}
	target := cue.ParsePath(fmt.Sprintf(`actions.%s`, action))
	runner := plan.NewRunner(target)
	err := runner.Run(context.Background(), v)
	return err
}

func getTargetPath(args []string) cue.Path {
	selectors := []cue.Selector{plan.ActionSelector}
	for _, arg := range args {
		selectors = append(selectors, cue.Str(arg))
	}
	return cue.MakePath(selectors...)
}

func loadPlan(ctx context.Context, planPath string) (*plan.Plan, error) {
	absPlanPath, err := filepath.Abs(planPath)
	if err != nil {
		return nil, err
	}

	_, err = os.Stat(absPlanPath)
	if err != nil {
		return nil, err
	}
	os.Chdir(absPlanPath)
	return plan.Load(ctx, plan.Config{
		Args: []string{absPlanPath},
	})
}

// run a action in cue
// TODO: run action with dependency
func doTest(filePath string, actionName string) (*cue.Value, error) {

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
