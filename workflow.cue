package main

import (
    "github.com/xxf098/actionflow"
    "github.com/xxf098/actionflow/core"
)

actionflow.#Plan & {
	actions: {
		setup: core.#Step & {
			uses: "actions/setup-go@v3"
			with: { "go-version": "1.19" }
		}

		build: core.#Exec & {
			cmd: ["sh", "-c", """
			go mod tidy
			make flow
"""]
		}
	}
}