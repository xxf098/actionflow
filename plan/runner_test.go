package plan

import (
	"context"
	"testing"

	"cuelang.org/go/cue"
	_ "github.com/xxf098/actionflow/plan/task"
)

func TestRun(t *testing.T) {
	v := LoadFile("../testcues/writefile1.cue")
	target := cue.ParsePath("actions.hello")
	runner := NewRunner(target)
	err := runner.Run(context.Background(), v)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRun1(t *testing.T) {
	v := LoadFile("../testcues/writefile1.cue")
	target := cue.ParsePath("actions.hello")
	runner := NewRunner(target)
	err := runner.Run(context.Background(), v)
	if err != nil {
		t.Fatal(err)
	}
}
