package testcues

import (
    "github.com/xxf098/actionflow"
    "github.com/xxf098/actionflow/core"
)

actionflow.#Plan & {
	actions: {

		mkdir: core.#Mkdir & {
			path:  "./test"
		}

		writeChecker: core.#WriteFile & {
			input:       mkdir.output
			path:        "./test/foo"
			contents:    "bar"
			permissions: 700
		}

		readChecker: core.#ReadFile & {
			input: writeChecker.output
			path:  "./test/foo"
		} & {
			// assert result
			output: "bar"
		}        
	
	}
}
