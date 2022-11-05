package task

import (
	"context"
	"errors"
	"os"

	"cuelang.org/go/cue"
	"github.com/xxf098/dagflow/compiler"
)

func init() {
	Register("WriteFile", func() Task { return &writeFileTask{} })
}

type writeFileTask struct {
}

func (t writeFileTask) Run(ctx context.Context, v *cue.Value) (*cue.Value, error) {
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
	value := compiler.NewValue()
	output := value.FillPath(cue.ParsePath("output"), p)
	return &output, nil
}

func (t *writeFileTask) Name() string {
	return "WriteFile"
}
