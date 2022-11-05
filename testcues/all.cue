package testcues

import (
    "github.com/xxf098/dagflow"
    "github.com/xxf098/dagflow/core"
)

dagflow.#Plan & {
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
