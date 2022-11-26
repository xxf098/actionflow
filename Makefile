SHELL := bash# we want bash behaviour in all shell invocations

.PHONY: flow

flow: # Build a dev flow binary
	CGO_ENABLED=0 go build -o ./cmd/flow -ldflags '-s -w' ./cmd/