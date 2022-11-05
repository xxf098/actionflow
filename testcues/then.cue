package testcues

import (
    "github.com/xxf098/dagflow"
    "github.com/xxf098/dagflow/core"
)

dagflow.#Plan & {
	actions: {

		mkdir: core.#Mkdir & {
			path:  "./test"
            then: core.#WriteFile & {
                path:        "./test/foo"
                contents:    "bar"
                permissions: 700
            }
		}
	
	}
}