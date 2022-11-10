package testcues

import (
    "github.com/xxf098/actionflow"
    "github.com/xxf098/actionflow/core"
)


#AddHello: {
	dir: string
    name: string

    write: core.#WriteFile & {
        input: dir
        path: "hello-\(name).txt"
        contents: "hello, \(name)!"
    }
}


actionflow.#Plan & {
      actions: hello: #AddHello & {
        dir: "."
        name: "fileName"
    }
}
