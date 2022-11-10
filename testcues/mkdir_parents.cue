package testcues

import (
    "github.com/xxf098/actionflow"
    "github.com/xxf098/actionflow/core"
)

actionflow.#Plan & {
	actions: {
		mkdir: core.#Mkdir & {
			path:  "./test/baz"
		}

		writeChecker: core.#WriteFile & {
			input:       mkdir.output
			path:        "./test/baz/foo"
			contents:    "bar"
			permissions: 700
		}

		readChecker: core.#ReadFile & {
			input: writeChecker.output
			path:  "./test/baz/foo"
		} & {
			// assert result
			output: "bar"
		}       
	
	}
}
