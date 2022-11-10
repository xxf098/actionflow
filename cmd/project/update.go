package project

import (
	"context"
	"log"

	"github.com/xxf098/actionflow/pkg"
)

func Update(ctx context.Context, dir string) {

	cueModPath, cueModExists := pkg.GetCueModParent(dir)
	if !cueModExists {
		log.Fatal("dagger project not found. Run `dagger project init`")
	}

	err := pkg.Vendor(ctx, cueModPath)
	if err != nil {
		log.Fatal(err)
	}
}
