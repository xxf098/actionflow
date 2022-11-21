package core

#env: [string]: bool | number | string

#Step: #Task & {
	$dagger: task: _name: "Step"
    // A name for your step to display on GitHub.
    uses: string
    with?: #env & {
            args?: string, entrypoint?: string, ...
    }
}