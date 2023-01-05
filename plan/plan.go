package plan

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"cuelang.org/go/cue"
	"github.com/xxf098/actionflow/compiler"
	"github.com/xxf098/actionflow/pkg"
)

var (
	ActionSelector = cue.Str("actions")
)

type Plan struct {
	config Config
	source cue.Value
}

type Config struct {
	Args   []string
	With   []string
	Target string
	DryRun bool
}

func Load(ctx context.Context, cfg Config) (*Plan, error) {

	planFileInfo, _ := os.Stat(cfg.Args[0])

	src := ""
	args := cfg.Args[0]

	var cueModExists bool

	if planFileInfo.IsDir() && filepath.IsAbs(args) {
		src, cueModExists = pkg.GetCueModParent(cfg.Args...)
		args = "."
	} else {
		_, cueModExists = pkg.GetCueModParent()
	}

	if !cueModExists {
		return nil, fmt.Errorf("project not found. Run `flow init`")
	}

	v, err := compiler.Build(ctx, src, args)
	if err != nil {
		return nil, err
	}
	if v.Err() != nil {
		return nil, v.Err()
	}

	p := &Plan{
		config: cfg,
		source: *v,
	}
	return p, nil
}

func (p *Plan) Do(ctx context.Context, path cue.Path) error {
	r := NewRunner(path)
	err := r.Run(ctx, p.source)
	return err
}
