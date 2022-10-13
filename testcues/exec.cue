package testcues

import (
    "github.com/xxf098/dagflow"
    "github.com/xxf098/dagflow/core"
)

#TouchFile: {
    name: string

    exec: core.#Exec & {
        args: ["touch", name]
    }
}

dagflow.#Plan & {
      actions: touch: #TouchFile & {
        name: "hello.txt"
    }
}