package testcues

import (
    "github.com/xxf098/actionflow"
    "github.com/xxf098/actionflow/core"
)

actionflow.#Plan & {
	actions: {
		setup: core.#Step & {
			uses: "actions/setup-node@v3"
			with: { "node-version": 16 }
		}
	}
}