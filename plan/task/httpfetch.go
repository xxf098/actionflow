package task

import (
	"context"

	"cuelang.org/go/cue"
	"github.com/xxf098/dagflow/compiler"
)

func init() {
	Register("HTTPFetch", func() Task { return &httpFetchTask{} })
}

type httpFetchTask struct {
}

func (t *httpFetchTask) Run(ctx context.Context, v *cue.Value) (*cue.Value, error) {
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

	value := compiler.NewValue()
	output := value.FillPath(cue.ParsePath("output"), httpFetch.Source)
	return &output, nil
}
