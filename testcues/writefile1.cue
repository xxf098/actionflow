package testcues

import (
    "github.com/xxf098/actionflow"
    "github.com/xxf098/actionflow/core"
)

#MkHello: {
	dir: string
 
    write: core.#Mkdir & {
        path: "\(dir)/hello"
    }
    result: write.output
}

#AddHello: {
	dir: string
    name: string

    write: core.#WriteFile & {
        path: "\(dir)/hello-\(name).txt"
        contents: "hello, \(name)!"
    }
}


actionflow.#Plan & {
      actions: {
          hellodir: #MkHello & {
                dir: "test"
         },
          hello: #AddHello & {
            dir: hellodir.result
            name: "fileName"
        }
      }
}
