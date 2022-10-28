package main

import "testing"

func TestWriteFile(t *testing.T) {
	Do("./testcues/writefile.cue", "hello")
}

func TestExec(t *testing.T) {
	Do("./testcues/exec.cue", "touch")
}

func TestRmFile(t *testing.T) {
	Do("./testcues/rmfile.cue", "hello")
}
