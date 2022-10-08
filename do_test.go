package main

import "testing"

func TestDo(t *testing.T) {
	Do("./main.cue", "hello")
}
