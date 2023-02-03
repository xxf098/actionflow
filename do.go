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

func doPlan[T PlanLoader](ctx context.Context, p string, action string) error {
	targetPath := getTargetPath([]string{action})
	var t T
	daggerPlan, err := t.load(ctx, p)
	if err != nil {
		return err
	}
	err = daggerPlan.Do(ctx, targetPath)
	if err != nil {
		return err
	}
	return nil
}

// load files in dir path
func Do(ctx context.Context, p string, action string) error {
	// targetPath := getTargetPath([]string{action})
	// daggerPlan, err := loadPlan(ctx, p)
	// if err != nil {
	// 	return err
	// }
	// err = daggerPlan.Do(ctx, targetPath)
	// if err != nil {
	// 	return err
	// }
	// return nil
	return doPlan[DirPlanLoader](ctx, p, action)
}

// load one file
func doFlowTest(cueFile string, action string) error {
	return doPlan[FilePlanLoader](context.Background(), cueFile, action)
}

type PlanLoader interface {
	load(ctx context.Context, p string) (*plan.Plan, error)
}

type DirPlanLoader struct {
}

func (l DirPlanLoader) load(ctx context.Context, planPath string) (*plan.Plan, error) {
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

type FilePlanLoader struct {
}

func (l FilePlanLoader) load(ctx context.Context, cueFile string) (*plan.Plan, error) {
	v := plan.LoadFile(cueFile)
	if v.Err() != nil {
		return nil, v.Err()
	}
	iter, _ := v.Fields()
	for iter.Next() {
		fmt.Println(iter.Label())
	}
	return plan.NewPlan(v), nil
}

func getTargetPath(args []string) cue.Path {
	selectors := []cue.Selector{plan.ActionSelector}
	for _, arg := range args {
		selectors = append(selectors, cue.Str(arg))
	}
	return cue.MakePath(selectors...)
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
