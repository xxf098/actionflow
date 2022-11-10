SHELL := bash# we want bash behaviour in all shell invocations

flow: # Build a dev flow binary
	CGO_ENABLED=0 go build -o ./flow -ldflags '-s -w' ./cmd/