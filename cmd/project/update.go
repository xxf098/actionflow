package project

import (
	"context"
	"fmt"

	"github.com/xxf098/actionflow/pkg"
)

func Update(ctx context.Context, dir string) error {

	cueModPath, cueModExists := pkg.GetCueModParent(dir)
	if !cueModExists {
		return fmt.Errorf("project not found. Run `flow init`")
	}

	err := pkg.Vendor(ctx, cueModPath)
	return err
}
