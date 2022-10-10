package task

import (
	"context"
	"errors"
	"os/exec"

	"cuelang.org/go/cue"
)

func init() {
	Register("Exec", func() Task { return &execTask{} })
}

type execTask struct {
}

func (t *execTask) Run(ctx context.Context, v *cue.Value) (*cue.Value, error) {
	common, err := parseCommon(v)
	if err != nil {
		return nil, err
	}
	// env

	name := common.args[0]
	args := common.args[1:]
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = common.workdir
	// cmd.Env
	err = cmd.Run()
	if err != nil {
		return nil, err
	}
	return NewValue()
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
