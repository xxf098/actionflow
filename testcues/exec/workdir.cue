package testcues

import (
    "github.com/xxf098/dagflow"
    "github.com/xxf098/dagflow/core"
)

dagflow.#Plan & {
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
