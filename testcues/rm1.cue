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
			path:        "./test/www.github.com.txt"
			contents:    "bar"
			permissions: 700
		}

		readChecker: core.#ReadFile & {
			input: writeChecker.output
			path:  "./test/www.github.com.txt"
            then: core.#Rm & { path: ["./test/www.github.com.txt"] }
		} & {
			// assert result
			output: "bar"
		}     

        
	
	}
}
