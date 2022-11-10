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
}


actionflow.#Plan & {
      actions: hello: #MkHello & {
        dir: "dir"
    }
}
