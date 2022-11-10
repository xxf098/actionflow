package testcues

import (
    "github.com/xxf098/actionflow"
    "github.com/xxf098/actionflow/core"
)

actionflow.#Plan & {
	actions: {

		exec: core.#Exec & {
			cmd: ["sh", "-c", "echo  $(date)"]
		}

        
		save: core.#WriteFile & {
			contents: exec.output
			path:  "./output.txt"
		}      
	
	}
}