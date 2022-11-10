package task

import (
	"context"
	"io/fs"
	"os"
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
	path, err := v.Lookup("path").String()
	if err != nil {
		return nil, err
	}

	// Permissions (int)
	permissions, err := v.Lookup("permissions").Int64()
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

	if opts.Parents {
		err = os.MkdirAll(path, fs.FileMode(permissions))
	} else {
		dir, err := os.Stat(path)
		if err == nil && dir.IsDir() {
			return nil, nil
		}
		err = os.Mkdir(path, fs.FileMode(permissions))
	}
	if err != nil {
		return nil, err
	}
	lg.Info().Dur("duration", time.Since(start)).Str("task", v.Path().String()).Msg(t.Name())
	Then(ctx, v)
	value := compiler.NewValue()
	output := value.FillPath(cue.ParsePath("output"), path)
	return &output, nil
}

func (t *mkdirTask) Name() string {
	return "Mkdir"
}
