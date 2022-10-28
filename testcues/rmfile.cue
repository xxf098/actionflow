package testcues

import (
    "github.com/xxf098/dagflow"
    "github.com/xxf098/dagflow/core"
)


#RmHello: {
	dir: string
 
    write: core.#Rm & {
        path: "hello.txt"
    }
}


dagflow.#Plan & {
      actions: hello: #RmHello & {
        dir: "."
    }
}
