package core

// HTTP operations

// Raw buildkit API
//
// package llb // import "github.com/moby/buildkit/client/llb"
//
// func HTTP(url string, opts ...HTTPOption) State
//
// type HTTPOption interface {
//         SetHTTPOption(*HTTPInfo)
// }
// func Checksum(dgst digest.Digest) HTTPOption
// func Chmod(perm os.FileMode) HTTPOption
// func Chown(uid, gid int) HTTPOption
// func Filename(name string) HTTPOption

Method: *"GET" | "POST" | "PUT" | "DELETE" | "OPTIONS" | "HEAD" | "CONNECT" | "TRACE" | "PATCH"

// Fetch a file over HTTP
#API: #Task & {
	$dagger: task: _name: "API"

	method: Method
	host:   string
	path:   string | *""
	auth?:  string
	headers?: [string]: string
	query?: [string]:   string
	form?: [string]:    string
	data?:    string | {...}
	timeout?: string
	// curl?: string
	retry?: {
		count: int | *3
		timer: string | *"6s"
		codes: [...int]
	}

	// filled by task
	resp: {
		status:     string
		statusCode: int

		body: *{} | bytes | string
		header: [string]:  string | [...string]
		trailer: [string]: string | [...string]
	}	
}
