package main

import (
	"context"
	"testing"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	_ "github.com/xxf098/dagflow/task"
)

func loadFile(filePath string) cue.Value {
	ctx := cuecontext.New()
	entrypoints := []string{filePath}

	bis := load.Instances(entrypoints, nil)
	return ctx.BuildInstance(bis[0])
}

func TestRun(t *testing.T) {
	v := loadFile("./testcues/writefile1.cue")
	target := cue.ParsePath("actions.hello")
	runner := NewRunner(target)
	err := runner.Run(context.Background(), v)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRun1(t *testing.T) {
	v := loadFile("./testcues/writefile1.cue")
	target := cue.ParsePath("actions.hello")
	runner := NewRunner(target)
	err := runner.Run(context.Background(), v)
	if err != nil {
		t.Fatal(err)
	}
}
