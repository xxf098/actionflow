package main

import (
    "github.com/xxf098/actionflow"
    "github.com/xxf098/actionflow/core"
)


#GoTest: core.#Exec & {
	fileName: string
	actionName: string
	test: string
 
	cmd: ["sh", "-c", """
	rm *.cue && cp ../workflow.cue ./
	sed -i '1 s/^.*$/package testcues/' workflow.cue
	cp ../testcues/\(fileName).cue ./
	./flow do \(actionName)
	\(test)
"""]
}

actionflow.#Plan & {
	actions: {
		setup: core.#Step & {
			uses: "actions/setup-go@v3"
			with: { 
				"go-version": "1.19"
				"cache": true 
			}
		}

		build: core.#Exec & {
			cmd: ["sh", "-c", """
			go mod tidy
			make flow
"""]
		}

		testAll: core.#All & {
			max: 1
			tasks: [
				#GoTest & { fileName: "writefile", actionName: "hello", test: "test -f hello-fileName.txt" },
				#GoTest & { fileName: "exec", actionName: "touch", test: "test -f ./hello.txt" },
				#GoTest & { fileName: "exec/output", actionName: "save", test: "test -f ./output.txt" },
				#GoTest & { fileName: "exec/run", actionName: "save", test: "test -f ./runoutput.txt" },
				#GoTest & { fileName: "exec/env", actionName: "verifyEnv", test: "" },
				#GoTest & { fileName: "exec/workdir", actionName: "verify", test: "" },
				#GoTest & { fileName: "then", actionName: "mkdirWrite", test: "test -f ./test/foo4" },
				#GoTest & { fileName: "mkdir_parents", actionName: "writeChecker", test: "test -f ./test/baz/foo" },
				#GoTest & { fileName: "all", actionName: "writeAll", test: "test -f ./test/foo1 && test -f ./test/foo2" },
				#GoTest & { fileName: "git", actionName: "actionflow", test: "test -f ./actionflow/go.mod && test -f ./actionflow/workflow.cue" },
				#GoTest & { fileName: "stdout", actionName: "print", test: "" },
			]
		}
	}
}