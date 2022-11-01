package compiler

import (
	"context"
	"errors"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	cueload "cuelang.org/go/cue/load"
)

// Build a cue configuration tree from the files
func Build(ctx context.Context, src string, args ...string) (*cue.Value, error) {

	buildConfig := &cueload.Config{
		Dir:     src,
		Overlay: map[string]cueload.Source{},
	}

	instances := cueload.Instances(args, buildConfig)
	if len(instances) != 1 {
		return nil, errors.New("only one package is supported at a time")
	}

	instance := instances[0]
	if err := instance.Err; err != nil {
		return nil, err
	}

	c := cuecontext.New()
	v := c.BuildInstance(instance)
	if err := v.Err(); err != nil {
		return nil, err
	}
	if err := v.Validate(); err != nil {
		return nil, err
	}
	return &v, nil
}
