package core


// Read the contents of a UTF-8 encoded file into a CUE string. Any non-UTF-8
// encoded content may have UTF replacement characters instead of the expected data.
#ReadFile: {
	$dagger: task: _name: "ReadFile"

	// Filesystem tree holding the file
	input: _
	// Path of the file to read
	path: string
	// Contents of the file
	output: string @dagger(generated)
}

// Write a file to a filesystem tree, creating it if needed
#WriteFile: {
	$dagger: task: _name: "WriteFile"

	// Input filesystem tree
	input: _
	// Path of the file to write
	path: string
	// Contents to write
	contents: string
	// Permissions of the file
	permissions: *0o644 | int
	// Output file path
	output: string @dagger(generated)		
}


// Remove file or directory from a filesystem tree
#Rm: {
	$dagger: task: _name: "Rm"

	input: _
	// Path to delete (handle wildcard)
	// (e.g. /file.txt or /*.txt)
	path: string | [...string]

	// Allow wildcard selection
	// Default to: true
	allowWildcard: *true | bool

	// Output filesystem tree
	output: string @dagger(generated)
}

// Create one or multiple directory in a container
#Mkdir: {
	$dagger: task: _name: "Mkdir"

	// Path of the directory to create
	// It can be nested (e.g : "/foo" or "/foo/bar")
	path: string

	// Permissions of the directory
	permissions: *0o755 | int

	// If set, it creates parents' directory if they do not exist
	parents: *true | false

	// dir full path
	output: string @dagger(generated)
}