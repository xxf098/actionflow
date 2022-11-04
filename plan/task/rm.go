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
	paths := []string{}
	pValue := v.Lookup("path")
	if pValue.IncompleteKind().IsAnyOf(cue.ListKind) {
		iter, err := pValue.List()
		if err != nil {
			return nil, err
		}
		for iter.Next() {
			path, err := iter.Value().String()
			if err != nil {
				return nil, err
			}
			paths = append(paths, path)
		}
	} else {
		path, err := pValue.String()
		if err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}

	var err error
	for _, path := range paths {
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
	}

	value := compiler.NewValue()
	output := value.FillPath(cue.ParsePath("output"), paths[0])
	return &output, nil
}
