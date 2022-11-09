package testcues

import (
    "github.com/xxf098/dagflow"
    "github.com/xxf098/dagflow/core"
)

#ReadHello: {
	p: string
 
    read: core.#ReadFile & {
        path: p
    }
    result: read.output
}

#AddHello: {
	content: string
    name: string

    write: core.#WriteFile & {
        path: "./test/hello-\(name).txt"
        contents: content
    }
}


dagflow.#Plan & {
      actions: {
          hellofile: #ReadHello & {
                p: "./go.mod"
         },
          hello: #AddHello & {
            content: hellofile.result
            name: "fileName"
        }
      }
}
