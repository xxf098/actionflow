package compiler

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"

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

	fileSystem := os.DirFS(src)
	err := fs.WalkDir(fileSystem, ".", func(p string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !entry.Type().IsRegular() {
			return nil
		}

		if filepath.Ext(entry.Name()) != ".cue" {
			return nil
		}

		contents, err := fs.ReadFile(fileSystem, p)
		if err != nil {
			return fmt.Errorf("%s: %w", p, err)
		}

		overlayPath := path.Join(buildConfig.Dir, p)
		buildConfig.Overlay[overlayPath] = cueload.FromBytes(contents)
		return nil
	})
	if err != nil {
		return nil, err
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
