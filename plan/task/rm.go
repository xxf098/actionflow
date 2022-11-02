package task

import (
	"context"
	"os"
	"path/filepath"
	"strings"

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

	if strings.Contains(path, "*") {
		paths, err := filepath.Glob(path)
		if err != nil {
			panic(err)
		}
		for _, p := range paths {
			err = os.RemoveAll(p)
			if err != nil {
				break
			}
		}
	} else {
		err = os.RemoveAll(path)
	}

	if err != nil {
		return nil, err
	}

	value := compiler.NewValue()
	output := value.FillPath(cue.ParsePath("output"), path)
	return &output, nil
}
