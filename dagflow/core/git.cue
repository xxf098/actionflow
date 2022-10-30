package core

#GitPull: {
	$dagger: task: _name: "GitPull"
	remote:     string
    depth?:      uint
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