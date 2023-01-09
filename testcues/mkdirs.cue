package testcues

import (
    "github.com/xxf098/actionflow"
    "github.com/xxf098/actionflow/core"
)


actionflow.#Plan & {
      actions: zero: core.#Mkdir & {
        path: ["zero/", "zero/world", "yellow/"]
    }
}
