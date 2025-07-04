package main

import (
	"fmt"
	"olivierkessler01/gontainers/process"
	"os"
	"runtime"
)


func main() {
	funcMap := map[string]func(args []string) bool{
		"List":  process.List,
	}

	var args []string
	runtime.Breakpoint()
	args = os.Args[1:]
	fmt.Println("Program starting")
	funcMap[args[0]](args[1:])
	fmt.Println("Program finished")
}
