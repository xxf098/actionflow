package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/xxf098/dagflow"
	"github.com/xxf098/dagflow/plan"
)

// https://cuelang.org/docs/concepts/packages/#import-path
func Do(dir string, action string) {
	ctx := context.Background()
	targetPath := getTargetPath([]string{action})
	daggerPlan, err := loadPlan(ctx, dir)
	if err != nil {
		log.Fatal(err)
	}
	err = daggerPlan.Do(ctx, targetPath)
	if err != nil {
		log.Fatal(err)
	}
}

func Flow(dir string, action string) {
	mainCue := path.Join(dir, "main.cue")
	fmt.Println(mainCue)
	v := loadFile(mainCue)
	iter, _ := v.Fields()
	for iter.Next() {
		fmt.Println(iter.Label())
	}
	target := cue.ParsePath(fmt.Sprintf(`actions.%s`, action))
	runner := dagflow.NewRunner(target)
	err := runner.Run(context.Background(), v)
	if err != nil {
		log.Fatal(err)
	}
}

func loadFile(filePath string) cue.Value {
	ctx := cuecontext.New()
	entrypoints := []string{filePath}

	bis := load.Instances(entrypoints, nil)
	return ctx.BuildInstance(bis[0])
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
	return plan.Load(ctx, plan.Config{
		Args: []string{planPath},
	})
}

func getTargetPath(args []string) cue.Path {
	selectors := []cue.Selector{plan.ActionSelector}
	for _, arg := range args {
		selectors = append(selectors, cue.Str(arg))
	}
	return cue.MakePath(selectors...)
}
