package core

// link list
#Task: {
	$dagger: task: _name: string
	input: 		_
	output: string @dagger(generated)
	then: 	null | #Task
    ...
}


#All: #Task & {
	$dagger: task: _name: "All"
	ignore_error: bool | *false
	tasks: 		[...#Task]
}
