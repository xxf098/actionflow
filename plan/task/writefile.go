package task

import (
	"context"
	"errors"
	"os"
	"time"

	"cuelang.org/go/cue"
	"github.com/rs/zerolog/log"
	"github.com/xxf098/dagflow/compiler"
)

func init() {
	Register("WriteFile", func() Task { return &writeFileTask{} })
}

type writeFileTask struct {
}

func (t writeFileTask) Run(ctx context.Context, v *cue.Value) (*cue.Value, error) {
	lg := log.Ctx(ctx)
	start := time.Now()
	p, err := v.Lookup("path").String()
	if err != nil {
		return nil, errors.New("fail to parse path")
	}
	contents, err := v.Lookup("contents").String()
	if err != nil {
		return nil, errors.New("fail to parse contents")
	}
	err = os.WriteFile(p, []byte(contents), 0644)
	if err != nil {
		return nil, err
	}
	lg.Info().Dur("duration", time.Since(start)).Str("task", v.Path().String()).Msg(t.Name())

	value := compiler.NewValue()
	output := value.FillPath(cue.ParsePath("output"), p)
	return &output, nil
}

func (t *writeFileTask) Name() string {
	return "WriteFile"
}
