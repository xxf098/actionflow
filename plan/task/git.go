package task

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os/exec"
	"strings"
	"time"

	"cuelang.org/go/cue"
	"github.com/rs/zerolog/log"
	"github.com/xxf098/actionflow/common/git"
	"github.com/xxf098/actionflow/compiler"
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
		Args  []string
		Repo  string
		Ref   string
		Dir   string
		Token string
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

	if len(args) < 1 && len(strings.TrimSpace(gitArgs.Repo)) < 1 {
		return nil, fmt.Errorf("not enough args")
	}

	if len(args) > 0 {
		var errBuf bytes.Buffer
		cmd := exec.Command("git", args...)
		cmd.Stderr = &errBuf
		err := cmd.Run()
		if err != nil {
			return nil, fmt.Errorf("%s: %s", err.Error(), errBuf.String())
		}
	} else {
		dir := gitArgs.Dir
		if len(dir) < 1 {
			u, err := url.Parse(gitArgs.Repo)
			if err != nil {
				return nil, err
			}
			splites := strings.Split(u.Path, "/")
			dir = splites[len(splites)-1]
			dir = strings.TrimSuffix(dir, ".git")
		}

		cfg := git.GitCloneConfig{
			URL:   gitArgs.Repo,
			Dir:   dir,
			Ref:   gitArgs.Ref,
			Token: gitArgs.Token,
		}
		if err := git.Clone(ctx, cfg); err != nil {
			return nil, err
		}
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
