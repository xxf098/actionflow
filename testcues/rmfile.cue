package testcues

import (
    "github.com/xxf098/dagflow"
    "github.com/xxf098/dagflow/core"
)


dagflow.#Plan & {
      actions: {
        test: {

            rmFile: {
                // Write a new file
                write: core.#WriteFile & {
                    path:     "./test.txt"
                    contents: "1,2,3"
                }

                // Remove file
                rm: core.#Rm & {
                    path:  write.path
                }

                verify: core.#Exec & {
                        input: rm.output
                        args: ["/bin/sh", "-e", "-c", """
                            test ! -e ./test.txt
                            """]
                    }
            }
            

		    rmWildcard: {
				// Write directory
				write: core.#Exec & {
					args: ["/bin/sh", "-e", "-c", """
						touch ./foo.txt
						touch ./bar.txt
						touch ./data.json
						"""]
				}

				// Remove all .txt file
				rm: core.#Rm & {
					input: write.output
					path:  "./*.txt"
				}

				verify: core.#Exec & {
					input: rm.output
					args: ["/bin/sh", "-e", "-c", """
							test ! -e ./foo.txt
							test ! -e ./bar.txt
							test -e ./data.json
						"""]
				    }
			    }            
            }
      }
}
