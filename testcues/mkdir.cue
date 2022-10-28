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
}


dagflow.#Plan & {
      actions: hello: #MkHello & {
        dir: "dir"
    }
}
