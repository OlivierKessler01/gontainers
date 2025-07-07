package main

import (
	"fmt"
	"olivierkessler01/gontainers/process"
	"os"
	//"runtime"
)

func main() {
	funcMap := map[string]func(args []string) error {
		"run":  process.Run,
		"list": process.List,
	}

	var args []string
	args = os.Args[1:]
	//runtime.Breakpoint()
	fmt.Println("Program starting")
	funcMap[args[0]](args[1:])
	fmt.Println("Program finished")
}
