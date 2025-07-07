package process  

import (
	"os"
	"encoding/json"
)

var TrackedProcesses []TrackedProcess 

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

func Load() error {
 	data, err := os.ReadFile("db.json")
    if err != nil {
        return  err
    }
	
    if err := json.Unmarshal(data, &TrackedProcesses); err != nil {
		return err
    }

    return nil
}

func Save() error {
	binary, err := json.Marshal(TrackedProcesses)
	if err != nil {
		return err
	}
	os.WriteFile("db.json", binary, 0644)
	return nil
}


