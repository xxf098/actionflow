package compiler

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
)

func NewValue() *cue.Value {
	c := cuecontext.New()
	v := c.CompileString("", cue.Filename("_"))
	if v.Err() != nil {
		panic(v.Err())
	}
	return &v
}
