package testcues

import (
    "github.com/xxf098/dagflow"
    "github.com/xxf098/dagflow/core"
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


dagflow.#Plan & {
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
