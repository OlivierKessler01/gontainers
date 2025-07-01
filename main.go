package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	var args []string 
	args = os.Args[:1]
	fmt.Println(args)
	fmt.Println("Hello world")
	//cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
