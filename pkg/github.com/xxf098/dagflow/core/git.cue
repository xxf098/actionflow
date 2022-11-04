package core

#Git: {
	$dagger: task: _name: "Git"
	input: 			_
	remote:     	string
    depth?:     	uint
	directory?: 	string
	auth?:      {
		username: string
		password: string // can be password or personal access token
	} | {
		authToken: string
	} | {
		authHeader: string
	}
	output: string @dagger(generated)
}