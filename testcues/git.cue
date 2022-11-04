package testcues

import (
    "github.com/xxf098/dagflow"
    "github.com/xxf098/dagflow/core"
)


#HelloClone: {
	url: string
 
    write: core.#Git & {
        remote: url
    }
}


dagflow.#Plan & {
      actions: hello: #HelloClone & {
        url: "https://github.com/SimaRyan/simaryan.github.io.git"
    }
}