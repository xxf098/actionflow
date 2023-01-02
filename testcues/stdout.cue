package testcues

import (
    "github.com/xxf098/actionflow"
    "github.com/xxf098/actionflow/core"
)

#Print: {
    name: string

    exec: core.#Stdout & {
        text: name
    }
}

actionflow.#Plan & {
      actions: {
        print: #Print & {
                name: "hello world"
            }
    }
}