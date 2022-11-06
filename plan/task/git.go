package task

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"cuelang.org/go/cue"
	"github.com/rs/zerolog/log"
	"github.com/xxf098/dagflow/compiler"
)

func init() {
	Register("Git", func() Task { return &gitTask{} })
}

type gitTask struct {
}

// FIXME: auth
func (t *gitTask) Run(ctx context.Context, v *cue.Value) (*cue.Value, error) {
	lg := log.Ctx(ctx)
	start := time.Now()
	var gitArgs struct {
		Args []string
	}

	if err := v.Decode(&gitArgs); err != nil {
		return nil, err
	}
	args := []string{}
	for _, arg := range gitArgs.Args {
		if len(strings.TrimSpace(arg)) > 0 {
			args = append(args, arg)
		}
	}

	if len(args) < 1 {
		return nil, fmt.Errorf("not enough args")
	}

	cmd := exec.Command("git", args...)
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	lg.Info().Dur("duration", time.Since(start)).Str("task", v.Path().String()).Msg(t.Name())
	Then(ctx, v)
	value := compiler.NewValue()
	output := value.FillPath(cue.ParsePath("output"), "")
	return &output, nil
}

func (t *gitTask) Name() string {
	return "Git"
}
