package testcues

import (
    "github.com/xxf098/actionflow"
    "github.com/xxf098/actionflow/core"
)

actionflow.#Plan & {
	actions: {

		exec: core.#Run & {
			run: "echo  $(date '+%Y-%m-%d %H:%M:%S')"
		}

        
		save: core.#WriteFile & {
			contents: exec.output
			path:  "./runoutput.txt"
		}      
	
	}
}