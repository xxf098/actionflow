package task

import (
	"context"
	"io/fs"
	"os"

	"cuelang.org/go/cue"
	"github.com/xxf098/dagflow/compiler"
)

func init() {
	Register("Mkdir", func() Task { return &mkdirTask{} })
}

type mkdirTask struct {
}

func (t *mkdirTask) Run(ctx context.Context, v *cue.Value) (*cue.Value, error) {
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
		err = os.Mkdir(path, fs.FileMode(permissions))
	}
	if err != nil {
		return nil, err
	}
	Then(ctx, v)
	value := compiler.NewValue()
	output := value.FillPath(cue.ParsePath("output"), path)
	return &output, nil
}
