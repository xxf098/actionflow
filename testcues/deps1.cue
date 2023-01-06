package testcues

import (
    "github.com/xxf098/actionflow"
    "github.com/xxf098/actionflow/core"
)

actionflow.#Plan & {
	actions: {
        rm: core.#Rm & {
            path: "./test"
        },

        mkdir: core.#Mkdir & {
            @$rm()
            path: "./test"
        }

        write: core.#WriteFile & {
            @$mkdir()
            path: "./test/hello.txt"
            contents: "hello world!"
        }

        write1: core.#WriteFile & {
            @$mkdir()
            path: "./test/hello1.txt"
            contents: "hello1 world!"
        }

        all: core.#All & {
            tasks: [
                core.#ReadFile & {
                    @debug()
                    path: write.path
                },
                core.#ReadFile & {
                    @debug()
                    path: write1.path
                },                
            ]
        }

	}
}
