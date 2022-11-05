package testcues

import (
    "github.com/xxf098/dagflow"
    "github.com/xxf098/dagflow/core"
)

dagflow.#Plan & {
	actions: {

		mkdir: core.#Mkdir & {
			path:  "./test"
            then: core.#WriteFile & {
                path:        "./test/foo"
                contents:    "bar"
                permissions: 700
            }
		}

        gitrm: core.#Git & {
            args: ["clone", "--depth=1", "https://github.com/SimaRyan/simaryan.github.io.git", "simaryan"]
            then: core.#Rm & {
                path: ["./simaryan/_config.yml", "README.md"]
            }
        }

        clone: core.#Git & {
             args: ["clone", "--depth=1", "https://github.com/changfengoss/pub.git"]
             then: core.#Exec & {
				args: ["sh", "-c", "ls ./pub/data | sort | tail -n 3 | xargs -I{} mv ./pub/data/{} ./pub/"]
				then:  core.#Rm & {					
					path: "./pub/data"
				}
			}
        }
	
	}
}