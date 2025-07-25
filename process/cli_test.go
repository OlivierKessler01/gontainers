package process

import (
	"strconv"
	"testing"
)

func TestList(t *testing.T) {
	args := []string {""}
	List(args)
}

func TestCreateListRemove(t *testing.T) {
	args := []string {"tail", "-f", "/dev/null"}
	 err := Run(args)
	if err != nil {
        t.Fatalf("Failed running the container. %s", err)
    }

	var nbRemoved, nbContainers int
	
	err = List([]string{})
	if err != nil || nbContainers != 1 {
		t.Fatalf("Wrong number of containers")
	}

	err = Remove([]string{strconv.Itoa(pid)})

	if err != nil || nbRemoved != 1 {
        t.Fatalf("Failed deleting container. %s", err)
    }
}

