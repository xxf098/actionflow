package task

import (
	"context"
	"errors"
	"fmt"
	"os/exec"

	"cuelang.org/go/cue"
	"github.com/xxf098/dagflow/compiler"
)

func init() {
	Register("Exec", func() Task { return &execTask{} })
}

type execTask struct {
}

// redirect output to current shell
func (t *execTask) Run(ctx context.Context, v *cue.Value) (*cue.Value, error) {
	common, err := parseCommon(v)
	if err != nil {
		return nil, err
	}
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
	name := common.args[0]
	args := common.args[1:]
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = common.workdir
	cmd.Env = append(cmd.Env, envs...)
	err = cmd.Run()
	if err != nil {
		return nil, err
	}
	Then(ctx, v)
	// FIXME: pipe output
	value := compiler.NewValue()
	output := value.FillPath(cue.ParsePath("output"), "")
	return &output, nil
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
