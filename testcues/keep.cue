package testcues

import (
    "github.com/xxf098/actionflow"
    "github.com/xxf098/actionflow/core"
)

#Write: core.#WriteFile & {
	fileName: string,
	path:        "./testkeep/\(fileName).txt"
	contents:    "\(fileName)"
	permissions: 700	
}

actionflow.#Plan & {
	actions: {

		mkdir: core.#Mkdir & {
			path:  "./testkeep"
		}

		writeFiles: core.#All & {
			@$()
			max: 2
			tasks: [ 
				#Write &{ fileName: "foo" },
				#Write &{ fileName: "foo1" },
				#Write &{ fileName: "bar" },
				#Write &{ fileName: "bar1" },
			]
		}

		keepFile: core.#Keep & {
            @$()
			path:  "./testkeep/foo*.txt"
		}     
	
	}
}
