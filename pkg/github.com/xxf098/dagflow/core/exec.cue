package core

#Exec: #Task & {
	$dagger: task: _name: "Exec"

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