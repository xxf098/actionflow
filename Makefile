SHELL := bash# we want bash behaviour in all shell invocations

dagger: # Build a dev dagger binary
	CGO_ENABLED=0 go build -o ./cmd/dagflow -ldflags '-s -w' ./cmd/