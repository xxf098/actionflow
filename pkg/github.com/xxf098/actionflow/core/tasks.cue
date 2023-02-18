package core

// link list
#Task: {
	$dagger: task: _name: string
	input: 		_
	output: string @dagger(generated)
	then: null | #Task
	// run task in background
	then_: 	null | #Task
    ...
}


#All: #Task & {
	$dagger: task: _name: "All"
	ignore_error: bool | *false
	// api limit
	max: int | *0
	tasks: 	[...#Task]
}


#Each: #All & {
	max: 1
}