package core

#GitPull: {
	$dagger: task: _name: "GitPull"
	input: 		_
	remote:     string
    depth?:     uint
	directory: 	string
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