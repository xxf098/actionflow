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
            @$write()
            path: "./test/hello1.txt"
            contents: "hello1 world!"
        }        

        read: core.#ReadFile & {
            @debug()
            path: write1.path
        }

	}
}
