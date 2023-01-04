package task

import (
	"context"

	"cuelang.org/go/cue"
	"github.com/xxf098/actionflow/compiler"
)

func init() {
	Register("API", func() Task { return &apiCallTask{} })
}

type apiCallTask struct {
}

func (t *apiCallTask) Run(ctx context.Context, v *cue.Value) (*cue.Value, error) {
	var httpFetch struct {
		Source      string
		Checksum    string
		Dest        string
		Permissions int
		UID         int
		GID         int
	}

	if err := v.Decode(&httpFetch); err != nil {
		return nil, err
	}
	Then(ctx, v)
	value := compiler.NewValue()
	output := value.FillPath(cue.ParsePath("output"), httpFetch.Source)
	return &output, nil
}

func (t *apiCallTask) Name() string {
	return "API"
}
