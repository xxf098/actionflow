package testcues

import (
    "github.com/xxf098/actionflow"
    "github.com/xxf098/actionflow/core"
)

actionflow.#Plan & {
	actions: {

		exec: core.#Exec & {
			cmd: ["sh", "-c", "echo -n hello world > ./output.txt"]
		}

        
		verify: core.#ReadFile & {
			input: exec.output
			path:  "./output.txt"
		} & {
			// assert result
			output: "hello world"
		}        
	
	}
}
