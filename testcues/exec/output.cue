package testcues

import (
    "github.com/xxf098/dagflow"
    "github.com/xxf098/dagflow/core"
)

dagflow.#Plan & {
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