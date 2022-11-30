package testcues

import (
    "github.com/xxf098/actionflow"
    "github.com/xxf098/actionflow/core"
)


#HelloClone: {
	url: string
 
    write: core.#Git & {
        args: ["clone", "--depth=1", url]
    }
}


actionflow.#Plan & {
      actions: {
        hello: #HelloClone & {
            url: "https://github.com/SimaRyan/simaryan.github.io.git"
        }
        checkout: core.#Git & {
           repo: "https://github.com/SimaRyan/simaryan.github.io.git"
        }
        actionflow: core.#Git & {
           repo: "https://github.com/xxf098/actionflow.git"
        }
      }
}