package project

import (
	"context"

	"github.com/xxf098/dagflow/pkg"
)

func Init(ctx context.Context, parentDir, module string) error {
	return pkg.CueModInit(ctx, parentDir, module)
}
