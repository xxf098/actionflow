package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/xxf098/actionflow/task"
)

// Go api
func main() {
	if len(os.Args) < 2 {
		fmt.Println("actionflow version: v0.1.0")
		fmt.Print(`useage:
	flow init [project path]
	flow update [project path]
	flow do [action name]	
`)
		return
	}
	args := os.Args[1:]
	if args[0] == "init" {
		projectInit(args)
		return
	}

	if args[0] == "update" {
		projectUpdate(args)
		return
	}
	if args[0] == "do" {
		flowAction(args)
		return
	}
}

func projectInit(args []string) {
	dir := "."
	if len(args) > 1 {
		dir = args[1]
	}
	err := task.Init(context.Background(), dir, "")
	if err != nil {
		log.Fatal(err)
	}
	if dir == "." {
		log.Println("Project initialized! To install actionflow packages, run `flow update`")
	} else {
		log.Printf("Project initialized in \"%s\"! To install actionflow packages, go to subfolder \"%s\" and run \"flow update\"\n", dir, dir)
	}
}

func projectUpdate(args []string) {
	dir, _ := os.Getwd()
	if len(args) > 1 {
		dir = args[1]
	}
	err := task.Update(context.Background(), dir)
	if err != nil {
		log.Fatalln(err)
	}
}

func flowAction(args []string) {
	action := args[1]
	dir, _ := os.Getwd()
	if len(args) > 2 {
		dir = args[2]
	}
	// Flow(dir, action)
	task.Do(dir, action)
}
