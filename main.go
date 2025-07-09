package main

import (
	"olivierkessler01/gontainers/process"
	"os"
)

func main() {
	funcMap := map[string]func(args []string) error{
		"run":  process.Run,
		"list": process.List,
	}

	var args []string
	args = os.Args[1:]
	funcMap[args[0]](args[1:])
}
