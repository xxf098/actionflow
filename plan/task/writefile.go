package task

import (
	"context"
	"errors"
	"os"
	"time"

	"cuelang.org/go/cue"
	"github.com/rs/zerolog/log"
	"github.com/xxf098/actionflow/compiler"
)

func init() {
	Register("WriteFile", func() Task { return &writeFileTask{} })
}

type writeFileTask struct {
}

func (t writeFileTask) Run(ctx context.Context, v *cue.Value) (*cue.Value, error) {
	lg := log.Ctx(ctx)
	start := time.Now()
	p, err := v.LookupPath(cue.ParsePath("path")).String()
	if err != nil {
		return nil, errors.New("fail to parse path")
	}
	contents, err := v.LookupPath(cue.ParsePath("contents")).String()
	if err != nil {
		return nil, errors.New("fail to parse contents")
	}
	append, err := v.LookupPath(cue.ParsePath("append")).Bool()
	if err != nil {
		return nil, errors.New("fail to parse append")
	}
	if append {
		err = appendFile(p, contents)
	} else {
		err = os.WriteFile(p, []byte(contents), 0644)
	}
	if err != nil {
		return nil, err
	}
	lg.Info().Dur("duration", time.Since(start)).Str("task", v.Path().String()).Msg(t.Name())
	Then(ctx, v)
	value := compiler.NewValue()
	output := value.FillPath(cue.ParsePath("output"), p)
	return &output, nil
}

func appendFile(path, contents string) error {
	f, err := os.OpenFile(path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(contents)
	return err
}

func (t *writeFileTask) Name() string {
	return "WriteFile"
}
