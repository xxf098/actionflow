package flow

import (
	"context"

	"github.com/xxf098/actionflow/pkg"
)

func Init(ctx context.Context, parentDir, module string) error {
	return pkg.CueModInit(ctx, parentDir, module)
}
