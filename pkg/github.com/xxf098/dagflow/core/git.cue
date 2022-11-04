package core

#Git: {
	$dagger: task: _name: "Git"
	input: 		_
	args: 		[...string]
	output: string @dagger(generated)
}

#GitClone: #Git & {
	args: ["clone", ...string]
}

#GitPull: #Git & {
	args: 	["pull", ...string]
}

#GitPush: #Git & {
	args: 	["push", ...string]
}