package main

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
)

func main() {
	// flowTask()
	// print()
	loadCue("./main.cue")
}

func print() {

	const val = `
i: int
s: "hello"
b: c: int
`

	var (
		c *cue.Context
		v cue.Value
	)

	// create a context
	c = cuecontext.New()

	// compile some CUE into a Value
	v = c.CompileString(val)

	// print the value
	fmt.Println(v)

}
