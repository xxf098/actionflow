package testcues

import (
    "github.com/xxf098/actionflow"
    "github.com/xxf098/actionflow/core"
)

actionflow.#Plan & {
	actions: {

        dir: core.#Mkdir & {
            path: "./testy"
        }

		writeEach: core.#Each & {
            input: dir.output
			tasks:  [
                core.#WriteFile & {
                    path:        "./testy/foo3"
                    contents:    "bar1"
                    permissions: 700
                },

                core.#WriteFile & {
                    path:        "./testy/foo4"
                    contents:    "bar2"
                    permissions: 700
                },                
            ]
		}
	
	}
}
