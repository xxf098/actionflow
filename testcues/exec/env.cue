package testcues

import (
    "github.com/xxf098/dagflow"
    "github.com/xxf098/dagflow/core"
)

dagflow.#Plan & {
	actions: {

		verify: core.#Exec & {
			env: TEST: "hello world"
			args: [
				"sh", "-c",
				#"""
					test "$TEST" = "hello world"
					"""#,
			]
		}        
	
	}
}
