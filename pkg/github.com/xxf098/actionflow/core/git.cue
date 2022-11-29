package core

#Git: #Task & {
	$dagger: task: _name: "Git"
	args: 	[...string] | *[]
	repo: string | *""
	ref: string | *"HEAD"
	dir: string | *""
	token: string | *""
}