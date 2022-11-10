package testcues

import (
    "github.com/xxf098/actionflow"
    "github.com/xxf098/actionflow/core"
)

#TouchFile: {
    name: string

    exec: core.#Exec & {
        cmd: ["touch", name]
    }
}

actionflow.#Plan & {
      actions: touch: #TouchFile & {
        name: "hello.txt"
    }
}