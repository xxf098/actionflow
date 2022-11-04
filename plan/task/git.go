package task

import (
	"context"
	"fmt"
	"os/exec"

	"cuelang.org/go/cue"
	"github.com/xxf098/dagflow/compiler"
)

func init() {
	Register("Git", func() Task { return &gitTask{} })
}

type gitTask struct {
}

// FIXME: auth
func (t *gitTask) Run(ctx context.Context, v *cue.Value) (*cue.Value, error) {
	var gitArgs struct {
		Args []string
	}

	if err := v.Decode(&gitArgs); err != nil {
		return nil, err
	}

	if len(gitArgs.Args) < 1 {
		return nil, fmt.Errorf("not enough args")
	}

	cmd := exec.Command("git", gitArgs.Args...)
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	value := compiler.NewValue()
	output := value.FillPath(cue.ParsePath("output"), "")
	return &output, nil
}
