package task

import (
	"context"
	"errors"
	"fmt"

	"cuelang.org/go/cue"
	"github.com/xxf098/actionflow/plan/github"
	"github.com/xxf098/actionflow/plan/github/model"
)

func init() {
	Register("Step", func() Task { return &stepTask{} })
}

// run github step
type stepTask struct {
}

func (t *stepTask) Run(ctx context.Context, v *cue.Value) (*cue.Value, error) {
	usesValue := v.Lookup("uses")
	uses, err := usesValue.String()
	if err != nil {
		return nil, err
	}
	withValue := v.Lookup("with")
	withs := map[string]string{}
	if withValue.Exists() {
		ik := withValue.IncompleteKind()
		if !(ik.IsAnyOf(cue.StructKind) && v.IsConcrete()) {
			return nil, errors.New("")
		}
		iter, _ := withValue.Fields()
		for iter.Next() {
			label := iter.Label()
			if v, err := iter.Value().String(); err == nil {
				withs[label] = v
			} else if v, err := iter.Value().Bool(); err == nil {
				withs[label] = fmt.Sprintf("%t", v)
			} else if v, err := iter.Value().Int64(); err == nil {
				withs[label] = fmt.Sprintf("%d", v)
			} else if v, err := iter.Value().Float64(); err == nil {
				withs[label] = fmt.Sprintf("%f", v)
			} else {
				return nil, fmt.Errorf("wrong field %s", label)
			}
		}
	}
	step := model.NewStep(uses, withs)
	if err := github.Executor(ctx, &step); err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *stepTask) Name() string {
	return "Step"
}