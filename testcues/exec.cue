package testcues

import (
    "github.com/xxf098/dagflow"
    "github.com/xxf098/dagflow/core"
)

#TouchFile: {
    name: string

    exec: core.#Exec & {
        cmd: ["touch", name]
    }
}

dagflow.#Plan & {
      actions: touch: #TouchFile & {
        name: "hello.txt"
    }
}