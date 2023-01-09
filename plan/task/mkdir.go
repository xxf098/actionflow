package task

import (
	"context"
	"io/fs"
	"os"
	"strings"
	"time"

	"cuelang.org/go/cue"
	"github.com/rs/zerolog/log"
	"github.com/xxf098/actionflow/compiler"
)

func init() {
	Register("Mkdir", func() Task { return &mkdirTask{} })
}

type mkdirTask struct {
}

func (t *mkdirTask) Run(ctx context.Context, v *cue.Value) (*cue.Value, error) {
	lg := log.Ctx(ctx)
	start := time.Now()
	paths := []string{}
	pValue := v.LookupPath(cue.ParsePath("path"))

	if pValue.IncompleteKind().IsAnyOf(cue.ListKind) {
		iter, err := pValue.List()
		if err != nil {
			return nil, err
		}
		for iter.Next() {
			path, err := iter.Value().String()
			if err != nil {
				return nil, err
			}
			paths = append(paths, path)
		}
	} else {
		path, err := pValue.String()
		if err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}

	// Permissions (int)
	permissions, err := v.LookupPath(cue.ParsePath("permissions")).Int64()
	if err != nil {
		return nil, err
	}

	// Retrieve options
	var opts struct {
		Parents bool
	}

	if err := v.Decode(&opts); err != nil {
		return nil, err
	}

	for _, path := range paths {
		if opts.Parents {
			err = os.MkdirAll(path, fs.FileMode(permissions))
		} else {
			dir, err := os.Stat(path)
			if err == nil && dir.IsDir() {
				continue
			}
			if err = os.Mkdir(path, fs.FileMode(permissions)); err != nil {
				return nil, err
			}
		}
	}

	if err != nil {
		return nil, err
	}
	lg.Info().Dur("duration", time.Since(start)).Str("task", v.Path().String()).Msg(t.Name())
	Then(ctx, v)
	value := compiler.NewValue()
	output := value.FillPath(cue.ParsePath("output"), strings.Join(paths, "\n"))
	return &output, nil
}

func (t *mkdirTask) Name() string {
	return "Mkdir"
}
