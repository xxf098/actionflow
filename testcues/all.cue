package testcues

import (
    "github.com/xxf098/actionflow"
    "github.com/xxf098/actionflow/core"
)

actionflow.#Plan & {
	actions: {

        dir: core.#Mkdir & {
            path: "./test"
        }

		writeAll: core.#All & {
            input: dir.output
			tasks:  [
                core.#WriteFile & {
                    path:        "./test/foo1"
                    contents:    "bar1"
                    permissions: 700
                },

                core.#WriteFile & {
                    path:        "./test/foo2"
                    contents:    "bar2"
                    permissions: 700
                },                
            ]
		}
	
	}
}
