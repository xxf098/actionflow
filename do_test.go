package main

import "testing"

func TestWriteFile(t *testing.T) {
	Do("./testcues/writefile.cue", "hello")
}
