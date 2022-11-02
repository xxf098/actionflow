package task

import (
	"context"
	"fmt"
	"os/exec"

	"cuelang.org/go/cue"
	"github.com/xxf098/dagflow/compiler"
)

func init() {
	Register("GitPull", func() Task { return &gitPullTask{} })
}

type gitPullTask struct {
}

// FIXME: auth
func (t *gitPullTask) Run(ctx context.Context, v *cue.Value) (*cue.Value, error) {
	var gitPull struct {
		Remote string
		Depth  int
		Auth   struct {
			Username string
		}
	}

	if err := v.Decode(&gitPull); err != nil {
		return nil, err
	}

	args := []string{"clone"}
	if gitPull.Depth > 0 {
		args = append(args, "--depth", fmt.Sprintf("%d", gitPull.Depth))
	}
	args = append(args, gitPull.Remote)

	cmd := exec.Command("git", args...)
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	value := compiler.NewValue()
	output := value.FillPath(cue.ParsePath("output"), gitPull.Remote)
	return &output, nil
}
