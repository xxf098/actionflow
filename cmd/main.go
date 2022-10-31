package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/xxf098/dagflow/cmd/project"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("dagflow version: v0.1.0")
		fmt.Println(`useage:
	dagflow init [project path]		
`)
		return
	}
	args := os.Args[1:]
	if args[0] == "init" {
		dir := "."
		if len(args) > 1 {
			dir = args[1]
		}
		projectInit(dir)
		return
	}
}

func projectInit(dir string) {

	err := project.Init(context.Background(), dir, "")
	if err != nil {
		log.Fatal(err)
	}
	if dir == "." {
		fmt.Println("Project initialized! To install dagger packages, run `dagflow update`")
	} else {
		fmt.Printf("Project initialized in \"%s\"! To install dagger packages, go to subfolder \"%s\" and run \"dagger project update\"", dir, dir)
	}
}
