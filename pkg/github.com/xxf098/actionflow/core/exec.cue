package core

#Exec: #Task & {
	$dagger: task: _name: "Exec"

	// cmd is the command to run.
	// Example: ["echo", "hello, world!"]
	cmd: string | [string, ...string]

	// Environment variables
	env: [key=string]: string

    // Working directory
	workdir: string | *"."

    // User ID or name
	user: string | *"root:root"

}

#Run: #Exec & {
	run: string
	cmd: ["sh", "-c", run]
}