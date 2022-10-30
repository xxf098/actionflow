package main

import (
	"fmt"
	"testing"

	"cuelang.org/go/cue"
)

func TestWriteFile(t *testing.T) {
	Do("./testcues/writefile.cue", "hello")
}

func TestWriteFile1(t *testing.T) {
	Do("./testcues/writefile1.cue", "hello")
}

func TestExec(t *testing.T) {
	Do("./testcues/exec.cue", "touch")
}

func TestRmFile(t *testing.T) {
	Do("./testcues/rmfile.cue", "hello")
}

func TestMkdir(t *testing.T) {
	output := Do("./testcues/mkdir.cue", "hello")
	v := output.LookupPath(cue.ParsePath("output"))
	s, _ := v.String()
	fmt.Println(s)
}
