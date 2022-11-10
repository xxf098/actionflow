package core

#Git: #Task & {
	$dagger: task: _name: "Git"
	args: 		[...string]

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