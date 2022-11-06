package testcues

import (
    "github.com/xxf098/dagflow"
    "github.com/xxf098/dagflow/core"
)

dagflow.#Plan & {
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
