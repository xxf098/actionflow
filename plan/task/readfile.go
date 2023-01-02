package task

import (
	"context"
	"fmt"
	"os"
	"time"

	"cuelang.org/go/cue"
	"github.com/rs/zerolog/log"
	"github.com/xxf098/actionflow/compiler"
)

func init() {
	Register("ReadFile", func() Task { return &readFileTask{} })
}

type readFileTask struct {
}

func (t readFileTask) Run(ctx context.Context, v *cue.Value) (*cue.Value, error) {
	lg := log.Ctx(ctx)
	start := time.Now()
	path, err := v.LookupPath(cue.ParsePath("path")).String()
	if err != nil {
		return nil, err
	}

	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ReadFile %s: %w", path, err)
	}
	text := string(contents)
	debug(v, text)
	lg.Info().Dur("duration", time.Since(start)).Str("task", v.Path().String()).Msg(t.Name())
	Then(ctx, v)
	value := compiler.NewValue()
	output := value.FillPath(cue.ParsePath("output"), text)
	return &output, nil
}

func (t *readFileTask) Name() string {
	return "ReadFile"
}
