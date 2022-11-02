package dagflow

import (
	"context"
	"testing"

	"cuelang.org/go/cue"
	_ "github.com/xxf098/dagflow/plan/task"
)

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
