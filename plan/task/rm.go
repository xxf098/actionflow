package task

import (
	"context"
	"os"

	"cuelang.org/go/cue"
	"github.com/xxf098/dagflow/compiler"
)

func init() {
	Register("Rm", func() Task { return &rmTask{} })
}

type rmTask struct {
}

func (t *rmTask) Run(ctx context.Context, v *cue.Value) (*cue.Value, error) {
	path, err := v.Lookup("path").String()
	if err != nil {
		return nil, err
	}
	err = os.RemoveAll(path)
	if err != nil {
		return nil, err
	}
	return compiler.NewValue(), nil
}
