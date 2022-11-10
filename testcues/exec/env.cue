package testcues

import (
    "github.com/xxf098/actionflow"
    "github.com/xxf098/actionflow/core"
)

actionflow.#Plan & {
	actions: {

		verify: core.#Exec & {
			env: TEST: "hello world"
			cmd: [
				"sh", "-c",
				#"""
					test "$TEST" = "hello world"
					"""#,
			]
		}        
	
	}
}
