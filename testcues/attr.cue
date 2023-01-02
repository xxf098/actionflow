package testcues

import (
    "github.com/xxf098/actionflow"
    "github.com/xxf098/actionflow/core"
)

#ReadHello: {
	p: string
 
    read: core.#ReadFile & {
        @debug()
        path: p
    }
    result: read.output
}

#Print: {
    name: string

    exec: core.#Stdout & {
        text: name
    }
}


actionflow.#Plan & {
      actions: {
        hellofile: #ReadHello & {
            p: "./Makefile"
        },
        print: #Print & {
            name: hellofile.result
        }
      }
}
