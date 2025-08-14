package process

import (
	"encoding/json"
	"errors"
	"log/slog"
	"github.com/olivierkessler01/gontainers/config"
	"os"
	"path/filepath"
	"strconv"
)

var TrackedContainers map[string]*Container
const DB_FILE = "db.json"
const DB_DEFAULT_FILE = "default.db.json"

func getDBFilePath() string {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
    return filepath.Join(cfg.DBPath, DB_FILE)
}

func IsTracked(containerId string) bool {
	_, exist := TrackedContainers[containerId]
	return exist
}

func Load() error {
	held, err := IsLockHeld()
	if err != nil {
		return err
	}

	if held == false {
		return errors.New("Cannot Load process " +
			"DB because lock is held by another goroutine.")	
	}

 	data, err := os.ReadFile(getDBFilePath())
    if err != nil {
		slog.Error("Error reaching for the database, did you " +
		"run `./gontainer init ?")
        return err
    }
	
    if err := json.Unmarshal(data, &TrackedContainers); err != nil {
		return err
    }
	
    files, err := os.ReadDir("/proc")
    if err != nil {
        return err
    }

    var systemPids map[int]bool 
	systemPids = make(map[int]bool)

    for _, f := range files {
        if f.IsDir() {
            pid, err := strconv.Atoi(f.Name())
            if err == nil {
				systemPids[pid] = true
            }
        }
    }

	for _,proc := range TrackedContainers {
		if _, present := systemPids[proc.PID]; !present {
			delete(TrackedContainers, proc.PID)
		}
	}
	
    return nil
}

func Save() error {
	held, err := IsLockHeld()
	if err != nil {
		return err
	}

	if held == false {
		return errors.New("Cannot Save process DB because lock is held by another goroutine.")	
	}

	binary, err := json.Marshal(TrackedContainers)
	if err != nil {
		return err
	}

	err = os.WriteFile(getDBFilePath(), binary, 0644)
	if err != nil {
		return err
	}
	return nil
}


