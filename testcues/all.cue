package testcues

import (
    "github.com/xxf098/actionflow"
    "github.com/xxf098/actionflow/core"
)

actionflow.#Plan & {
	actions: {

        dir: core.#Mkdir & {
            path: "./testt"
        }

		writeAll: core.#All & {
            input: dir.output
			tasks:  [
                core.#WriteFile & {
                    path:        "./testt/foo1"
                    contents:    "bar1"
                    permissions: 700
                },

                core.#WriteFile & {
                    path:        "./testt/foo2"
                    contents:    "bar2"
                    permissions: 700
                },                
            ]
		}
	
	}
}
