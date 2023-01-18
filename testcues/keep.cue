package testcues

import (
    "github.com/xxf098/actionflow"
    "github.com/xxf098/actionflow/core"
)

actionflow.#Plan & {
	actions: {

		mkdir: core.#Mkdir & {
			path:  "./test1"
		}

		writeChecker: core.#WriteFile & {
			@$()
			path:        "./test1/foo.txt"
			contents:    "bar"
			permissions: 700
		}

        writeChecker1: core.#WriteFile & {
			@$()
			path:        "./test1/foo1.txt"
			contents:    "bar"
			permissions: 700
		}

        writeChecker2: core.#WriteFile & {
			@$()
			path:        "./test1/bar.txt"
			contents:    "bar1"
			permissions: 700
		}

        writeChecker3: core.#WriteFile & {
			@$()
			path:        "./test1/bar1.txt"
			contents:    "bar2"
			permissions: 700
		}

		keepFile: core.#Keep & {
            @$()
			path:  "./test1/foo*.txt"
		}     
	
	}
}
