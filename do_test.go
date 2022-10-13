package main

import "testing"

func TestWriteFile(t *testing.T) {
	Do("./main.cue", "hello")
}
