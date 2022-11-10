package testcues

import (
    "github.com/xxf098/actionflow"
    "github.com/xxf098/actionflow/core"
)

actionflow.#Plan & {
	actions: {

		verify: core.#Exec & {
			workdir: "/tmp"
			cmd: [
				"sh", "-c",
				#"""
					test "$(pwd)" = "/tmp"
					"""#,
			]
		}          
	
	}
}
