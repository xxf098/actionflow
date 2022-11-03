package testcues

import (
    "github.com/xxf098/dagflow"
    "github.com/xxf098/dagflow/core"
)

dagflow.#Plan & {
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
