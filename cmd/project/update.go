package project

import (
	"context"
	"log"

	"github.com/xxf098/dagflow/pkg"
)

func Update(dir string) {

	ctx := context.Background()

	cueModPath, cueModExists := pkg.GetCueModParent(dir)
	if !cueModExists {
		log.Fatal("dagger project not found. Run `dagger project init`")
	}

	err := pkg.Vendor(ctx, cueModPath)
	if err != nil {
		log.Fatal(err)
	}
}
