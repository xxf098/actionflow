package core

#Git: #Task & {
	$dagger: task: _name: "Git"
	args: 	[...string] | *[]
	repo: string | *""
	ref: string | *""
	dir: string | *""
	token: string | *""
}