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
	rm workflow.cue
	cp ../testcues/\(fileName).cue ./
	./flow do \(actionName)
	\(test)
	rm *.cue && cp ../workflow.cue ./
	sed -i '1 s/^.*$/package testcues/' workflow.cue
"""]
}

#TestFunc: core.#Exec & {
	name: string
	cmd: ["sh", "-c", """
go test -run \(name)
"""]
}

funcs: ["TestWriteFile", "TestWriteFile1", "TestWriteFile2", "TestExec", "TestRun", "TestRm1",
	"TestMkDir1", "TestMkDirs", "TestMkDirParent", "TestMkdir", "TestGitPull", "TestWorkDir", "TestArgs", "TestOutput", "TestAll", "TestStep", "TestStdout", "TestReadfile", "TestKeep", "TestKeep1", "TestKeep2"]

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

		testKeep: #GoTest & {
			fileName: "keep", 
			actionName: "keepFile", 
			test: "test -f testkeep/foo.txt && test -f testkeep/foo1.txt"
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
				#GoTest & { fileName: "then", actionName: "mkdirWrite", test: "test -f ./test/foo4 && test -f ./test/foo5" },
				#GoTest & { fileName: "mkdir_parents", actionName: "writeChecker", test: "test -f ./test/baz/foo" },
				#GoTest & { fileName: "all", actionName: "writeAll", test: "test -f ./testt/foo1 && test -f ./testt/foo2" },
				#GoTest & { fileName: "git", actionName: "actionflow", test: "test -f ./actionflow/go.mod && test -f ./actionflow/workflow.cue" },
				#GoTest & { fileName: "stdout", actionName: "print", test: "" },
				#GoTest & { fileName: "api", actionName: "print", test: "" },
				#GoTest & { fileName: "deps", actionName: "read", test: "" },
				#GoTest & { fileName: "deps1", actionName: "all", test: "" },
				#GoTest & { fileName: "mkdirs", actionName: "zero", test: "test -d zero/world && test -d yellow" },
			]
		}

		testFuncs: core.#All & {
			max: 2
			tasks: [ for f in funcs { #TestFunc & { name: f } } ]
		}
	}
}