package main

import (
    "github.com/xxf098/actionflow"
    "github.com/xxf098/actionflow/core"
)


#GoTest: {
	fileName: string
	actionName: string
	test: string
 
    write: core.#Exec & {
        cmd: ["sh", "-c", """
		cp ../testcues/\(fileName).cue ./
		./flow do \(actionName)
		rm -f \(fileName).cue
		\(test)
"""]
    }
}

actionflow.#Plan & {
	actions: {
		setup: core.#Step & {
			uses: "actions/setup-go@v3"
			with: { "go-version": "1.19" }
		}

		build: core.#Exec & {
			cmd: ["sh", "-c", """
			go mod tidy
			make flow
"""]
		}

		testAll: core.#All & {
			tasks: [
				#GoTest & {
					fileName: "writefile"
					actionName: "hello"
					test: "test -f hello-fileName.txt"
				}
			]
		}
	}
}