package process

import (
	"encoding/json"
	"fmt"
	"olivierkessler01/gontainers/config"
	"os"
	"path/filepath"
	"strconv"
)

var TrackedProcesses []TrackedProcess 
const LOCK_FILE = "db.lock"
const DB_FILE = "db.json"

func getDBFilePath() string {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
    return filepath.Join(cfg.DBPath, DB_FILE)
}

func getLockFilePath() string {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
    return filepath.Join(cfg.DBPath, LOCK_FILE)
}

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

func AcquireLock() error {
	if _, err := os.Stat(getLockFilePath()); os.IsNotExist(err) {
        file, err := os.Create(getLockFilePath())
        if err != nil {
            fmt.Println("Error acquiring lock:", err)
            return err
        }
        defer file.Close()
        fmt.Println("Lock acquired:", getLockFilePath())
    } else {
		return fmt.Errorf("Cannot acquire lock, someone already has it: %s", getLockFilePath())
    }

	return nil
}

func ReleaseLock() error {
	if _, err := os.Stat(getLockFilePath()); os.IsNotExist(err) {
        fmt.Println("Lock already released:", getLockFilePath())
    } else {
		err := os.Remove(getLockFilePath())
		if err != nil {
        	fmt.Println("Failure releasing lock:", getLockFilePath())
			return err
		}
		fmt.Println("Lock successfully released:", getLockFilePath())
    }

	return nil
}

func IsTracked(pid int) bool {
    for _, proc := range TrackedProcesses {
		if proc.PID == pid {
			return true
		}
    }
	
	return false
}

func Load() error {
	err := AcquireLock()

	if err != nil {
		return err
	}

 	data, err := os.ReadFile(getDBFilePath())
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
	binary, err := json.Marshal(TrackedProcesses)
	if err != nil {
		return err
	}

	err = os.WriteFile(getDBFilePath(), binary, 0644)
	if err != nil {
		return err
	}
	return nil
}


