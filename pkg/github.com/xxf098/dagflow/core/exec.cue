package core

#Exec: {
	$dagger: task: _name: "Exec"

	input: _

    // Command to execute
	// Example: ["echo", "hello, world!"]
	args: [...string]

	// Environment variables
	env: [key=string]: string

    // Working directory
	workdir: string | *"."

    // User ID or name
	user: string | *"root:root"
}