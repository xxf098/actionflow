package main

import (
	"fmt"

	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
)

// cue.mod in the root of go project
// import mod inside pkg
func loadCue() {
	ctx := cuecontext.New()
	entrypoints := []string{"./main.cue"}

	bis := load.Instances(entrypoints, nil)

	for _, bi := range bis {

		if bi.Err != nil {
			fmt.Println("Error during load:", bi.Err)
			continue
		}

		value := ctx.BuildInstance(bi)
		if value.Err() != nil {
			fmt.Println("Error during build:", value.Err())
			continue
		}

		fmt.Printf("main value: %v", value)
	}
}
