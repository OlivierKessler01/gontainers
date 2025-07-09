package process  

import (
	"os"
	"encoding/json"
	"strconv"
	"fmt"
)

var TrackedProcesses []TrackedProcess 
var LOCK_FILE = "db.lock"

type TrackedProcess struct {
	PID int
	cgroups []string
	namespaces []string
}

func Add(pid int, cgroups, namespaces []string) (TrackedProcess, error) {
	t := TrackedProcess{
        PID:        pid,
        cgroups:    cgroups,
        namespaces: namespaces,
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
        fmt.Println("Cannot acquire lock, someone already has it:", LOCK_FILE)
		return err
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
        return  err
    }
	
    if err := json.Unmarshal(data, &TrackedProcesses); err != nil {
		return err
    }

    files, err := os.ReadDir("/proc")
    if err != nil {
        return err
    }

    var pids []TrackedProcess
	//runtime.Breakpoint()
    for _, f := range files {
        if f.IsDir() {
            pid, err := strconv.Atoi(f.Name())
            if err == nil {
				pids = append(pids, TrackedProcess{PID:pid})
            }
        }
    }
	
    return nil
}

func Save() error {
	binary, err := json.Marshal(TrackedProcesses)
	if err != nil {
		return err
	}

	err = os.WriteFile("db.json", binary, 0644)
	if err != nil {
		return err
	}
	defer releaseLock()
	return nil
}


