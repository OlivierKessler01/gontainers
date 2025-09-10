package main

import (
	"os"
	"testing"
)

func Test(t *testing.T) {
	testCases := [][]string{
		{"run", "--name=test", "--command=tail -f /dev/null"},
		{"list"},
		{"remove"},
		{"list"},
	}

	for _, args := range testCases {
		progArgs := append(os.Args[0:1], args...)
		if err := run(progArgs); err != nil {
			t.Fatalf("Command %v failed: %v", args, err)
		}
	}
}
