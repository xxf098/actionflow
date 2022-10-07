package main

import (
    "github.com/xxf098/dagflow"
    "github.com/xxf098/dagflow/core"
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


dagflow.#Plan & {
      actions: hello: #AddHello & {
        dir: "."
        name: "fileName"
    }
}
