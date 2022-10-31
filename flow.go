package dagflow

import (
	"context"
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/tools/flow"
)

var CTX cue.Context

var input = `
tasks: {
	a: {
		foo: 1
		hello: string
	}
	b: {
		foo: 2
	}
	c: {
		foo: a.foo * 3
		goo: b.foo * 3
	}
}
`

func flowTask() {

	var err error
	fmt.Println("Custom Flow Task")

	// create context
	ctx := cuecontext.New()

	// Setup the flow Config
	cfg := &flow.Config{Root: cue.ParsePath("tasks")}

	// compile our input
	value := ctx.CompileString(input, cue.Filename("input.cue"))
	if value.Err() != nil {
		fmt.Println("Error:", value.Err())
		return
	}

	// create the workflow whiich will build the task graph
	workflow := flow.New(cfg, value, TaskFactory)

	fmt.Println("===RUN===")
	// run our custom workflow
	err = workflow.Run(context.Background())
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

}

func TaskFactory(val cue.Value) (flow.Runner, error) {
	// You can see the recursive values with this

	// Check that we have something that looks like a task
	foo := val.Lookup("foo")
	if !foo.Exists() {
		fmt.Println("TF 1: ", val)
		return nil, nil
	}
	fmt.Println("TF 2: ", val)

	num, err := foo.Int64()
	if err != nil {
		return nil, err
	}

	// Create and return a flow.Runner
	ct := &CustomTask{
		Val: int(num),
	}
	return ct, nil
}

type CustomTask struct {
	Val int
}

func (C *CustomTask) Run(t *flow.Task, pErr error) error {
	// not sure this is OK, but the value which was used for this task
	val := t.Value()
	fmt.Println("CustomTask:", C.Val, val)

	// Do some work
	next := map[string]interface{}{
		"bar": C.Val + 1,
	}
	hello := val.LookupPath(cue.ParsePath("foo"))
	if hello.Exists() {
		next["hello"] = "world"
	}

	// Use fill to "return" a result to the workflow engine
	t.Fill(next)

	return nil
}
