package task

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"cuelang.org/go/cue"
	"github.com/rs/zerolog/log"
	"github.com/xxf098/actionflow/compiler"
)

func init() {
	Register("Exec", func() Task { return &execTask{} })
}

type execTask struct {
}

// redirect output to current shell
func (t *execTask) Run(ctx context.Context, v *cue.Value) (*cue.Value, error) {
	// common, err := parseCommon(v)
	// if err != nil {
	// 	return nil, err
	// }
	// // env
	// it, err := v.Lookup("env").Fields()
	// if err != nil {
	// 	return nil, err
	// }
	// envs := []string{}
	// for it.Next() {
	// 	key := it.Label()
	// 	value, err := it.Value().String()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	env := fmt.Sprintf("%s=%v", key, value)
	// 	envs = append(envs, env)
	// }
	// name := common.args[0]
	// args := common.args[1:]
	// cmd := exec.CommandContext(ctx, name, args...)
	// cmd.Dir = common.workdir
	// cmd.Env = append(cmd.Env, envs...)
	lg := log.Ctx(ctx)
	start := time.Now()
	cmd, err := mkCommand(ctx, v)
	if err != nil {
		return nil, err
	}
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", err.Error(), errBuf.String())
	}
	s := outBuf.String()
	if len(s) > 0 {
		lg.Debug().Str("task", v.Path().String()).Msg(s)

	}
	lg.Info().Dur("duration", time.Since(start)).Str("task", v.Path().String()).Msg(t.Name())
	Then(ctx, v)
	value := compiler.NewValue()
	output := value.FillPath(cue.ParsePath("output"), s)
	return &output, nil
}

func (t *execTask) Name() string {
	return "Exec"
}

type execCommon struct {
	args    []string
	workdir string
	user    string
	hosts   map[string]string
}

func parseCommon(v *cue.Value) (*execCommon, error) {
	e := &execCommon{}

	// args
	var cmd struct {
		Args []string
	}

	if err := v.Decode(&cmd); err != nil {
		return nil, err
	}
	if len(cmd.Args) < 1 {
		return nil, errors.New("not enough argument")
	}
	e.args = cmd.Args

	// workdir
	workdir, err := v.Lookup("workdir").String()
	if err != nil {
		return nil, err
	}
	e.workdir = workdir

	// user
	user, err := v.Lookup("user").String()
	if err != nil {
		return nil, err
	}
	e.user = user
	return e, nil
}

func mkCommand(ctx context.Context, v *cue.Value) (c *exec.Cmd, err error) {
	var bin string
	var args []string

	cv := v.Lookup("cmd")
	switch cv.Kind() {
	case cue.StringKind:
		s, err := cv.String()
		if err != nil {
			return nil, err
		}
		list := strings.Fields(s)
		bin = list[0]
		args = append(args, list[1:]...)
	case cue.ListKind:
		list, err := cv.List()
		if err != nil {
			return nil, err
		}
		if !list.Next() {
			return nil, errors.New("empty command list")
		}
		bin, err = list.Value().String()
		if err != nil {
			return nil, err
		}
		for list.Next() {
			str, err := list.Value().String()
			if err != nil {
				return nil, err
			}
			args = append(args, str)
		}
	}

	if bin == "" {
		return nil, errors.New("empty command")
	}

	cmd := exec.CommandContext(ctx, bin, args...)

	workdir, err := v.Lookup("workdir").String()
	if err != nil {
		return nil, err
	}
	cmd.Dir = workdir

	// env
	it, err := v.Lookup("env").Fields()
	if err != nil {
		return nil, err
	}
	envs := []string{}
	for it.Next() {
		key := it.Label()
		value, err := it.Value().String()
		if err != nil {
			return nil, err
		}
		env := fmt.Sprintf("%s=%v", key, value)
		envs = append(envs, env)
	}
	cmd.Env = append(cmd.Env, envs...)

	return cmd, nil
}
