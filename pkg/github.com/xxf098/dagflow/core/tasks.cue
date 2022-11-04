package core

#Task: {
	$dagger: task: _name: string
	input: 		_
	output: string @dagger(generated)
    ...
}


#Tasks: #Task & {
	$dagger: task: _name: "Tasks"
	tasks: 		[...]
}
