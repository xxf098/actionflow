package github

import (
	"context"

	"github.com/xxf098/actionflow/plan/github/model"
)

type step interface {
	pre(ctx context.Context) error
	main(ctx context.Context) error
	post(ctx context.Context) error

	getStepModel() *model.Step
	getEnv() *map[string]string
}
