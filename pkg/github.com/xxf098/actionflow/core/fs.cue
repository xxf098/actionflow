package core


// Read the contents of a UTF-8 encoded file into a CUE string. Any non-UTF-8
// encoded content may have UTF replacement characters instead of the expected data.
#ReadFile: #Task & {
	$dagger: task: _name: "ReadFile"

	// Path of the file to read
	path: !=""
}

// Write a file to a filesystem tree, creating it if needed
#WriteFile: #Task & {
	$dagger: task: _name: "WriteFile"

	// Path of the file to write
	path: !=""
	// Contents to write
	contents: string
	// Permissions of the file
	permissions: *0o644 | int
	// append to file
	append: bool | *false
}


// Remove file or directory from a filesystem tree
#Rm: #Task & {
	$dagger: task: _name: "Rm"

	// Path to delete (handle wildcard)
	// (e.g. /file.txt or /*.txt)
	path: string | [...string]

	// Allow wildcard selection
	// Default to: true
	allowWildcard: *true | bool
}

// Create one or multiple directory in a container
#Mkdir: #Task & {
	$dagger: task: _name: "Mkdir"

	// Path of the directory to create
	// It can be nested (e.g : "/foo" or "/foo/bar")
	path: string

	// Permissions of the directory
	permissions: *0o755 | int

	// If set, it creates parents' directory if they do not exist
	parents: *true | false
}

// Create one or multiple directory in a container
#Stdout: #Task & {
	$dagger: task: _name: "Stdout"

	// text to write
	text: string
}