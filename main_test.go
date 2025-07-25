package main

import (
	"os"
	"testing"
)

func Test(t *testing.T) {
	args := os.Args[0:1] // Name of the program.
	args = append(args, []string{"run", "tail -f /dev/null"}...)
	run(args)

	args = os.Args[0:1] // Name of the program.
	args = append(args, "list")
	run(args)

	args = os.Args[0:1] // Name of the program.
	args = append(args, "remove")
	run(args)

	args = os.Args[0:1] // Name of the program.
	args = append(args, "list")
	run(args)
}
