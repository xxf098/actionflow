package core

import "github.com/xxf098/dagflow"


#ReadFile: {
	$dagger: task: _name: "ReadFile"

	// Filesystem tree holding the file
	input: dagflow.#FS
	// Path of the file to read
	path: string
	// Contents of the file
	contents: string @dagger(generated)
}


#WriteFile: {
	$dagger: task: _name: "WriteFile"

	// Input filesystem tree
	input: dagger.#FS
	// Path of the file to write
	path: string
	// Contents to write
	contents: string
	// Permissions of the file
	permissions: *0o644 | int
	// Output filesystem tree
	output: dagger.#FS @dagger(generated)
}