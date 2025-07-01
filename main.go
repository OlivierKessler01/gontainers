package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {

	args = os.Args[:1]
	fmt.Println("Hello world")
	cmd := exec.Command("ls", "-la")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
