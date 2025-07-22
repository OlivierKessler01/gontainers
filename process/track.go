package process

import (
	"encoding/json"
	"errors"
	"olivierkessler01/gontainers/config"
	"os"
	"path/filepath"
	"strconv"
)

var TrackedProcesses []TrackedProcess 
const DB_FILE = "db.json"

func getDBFilePath() string {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
    return filepath.Join(cfg.DBPath, DB_FILE)
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


func IsTracked(pid int) bool {
    for _, proc := range TrackedProcesses {
		if proc.PID == pid {
			return true
		}
    }
	
	return false
}


func Load() error {
	held, err := IsLockHeld()
	if err != nil {
		return err
	}

	if held == false {
		return errors.New("Cannot Load process DB because lock is held by another goroutine.")	
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


