package process

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

var TrackedProcesses []TrackedProcess 
var LOCK_FILE = "db.lock"

type TrackedProcess struct {
	PID int
	Cgroups []string
	Namespaces []string
}

func Add(pid int, cgroups, namespaces []string) (TrackedProcess, error) {
	t := TrackedProcess{
        PID:        pid,
        Cgroups:    cgroups,
        Namespaces: namespaces,
    }
	TrackedProcesses = append(TrackedProcesses, t)
	return t, nil 
}

func acquireLock() error {
	if _, err := os.Stat(LOCK_FILE); os.IsNotExist(err) {
        file, err := os.Create(LOCK_FILE)
        if err != nil {
            fmt.Println("Error acquiring lock:", err)
            return err
        }
        defer file.Close()
        fmt.Println("Lock acquired:", LOCK_FILE)
    } else {
		return fmt.Errorf("Cannot acquire lock, someone already has it: %s", LOCK_FILE)
    }

	return nil
}

func releaseLock() error {
	if _, err := os.Stat(LOCK_FILE); os.IsNotExist(err) {
        fmt.Println("Lock already released:", LOCK_FILE)
    } else {
		err := os.Remove(LOCK_FILE)
		if err != nil {
        	fmt.Println("Failure releasing lock:", LOCK_FILE)
			return err
		}
		fmt.Println("Lock successfully released:", LOCK_FILE)
    }

	return nil
}

func Load() error {
	err := acquireLock()

	if err != nil {
		return err
	}

 	data, err := os.ReadFile("db.json")
    if err != nil {
        return err
    }
	
    if err := json.Unmarshal(data, &TrackedProcesses); err != nil {
		return err
    }
	
    files, err := os.ReadDir("/proc")
    if err != nil {
        return err
    }

    var pids map[int]bool 
	pids = make(map[int]bool)

    for _, f := range files {
        if f.IsDir() {
            pid, err := strconv.Atoi(f.Name())
            if err == nil {
				pids[pid] = true
            }
        }
    }

	var newTrackedProcesses []TrackedProcess
	for _,proc := range TrackedProcesses {
		if _, present := pids[proc.PID]; present {
			newTrackedProcesses = append(newTrackedProcesses, proc)
		}
	}
	
	TrackedProcesses = newTrackedProcesses
	
    return nil
}

func Save() error {
	defer releaseLock()
	binary, err := json.Marshal(TrackedProcesses)
	if err != nil {
		return err
	}

	err = os.WriteFile("db.json", binary, 0644)
	if err != nil {
		return err
	}
	return nil
}


