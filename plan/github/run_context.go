package github

import (
	"context"

	"github.com/xxf098/actionflow/plan/github/model"
)

func Executor(ctx context.Context, step *model.Step) error {
	sar := newStepActionRemote(step)
	if err := sar.pre(ctx); err != nil {
		return err
	}
	if err := sar.main(ctx); err != nil {
		return err
	}
	return sar.post(ctx)
}
