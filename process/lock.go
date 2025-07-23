package process

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"github.com/olivierkessler01/gontainers/config"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

const LOCK_FILE = "db.lock"
var CURRENT_GOROUTINE_ID uuid.UUID

func getLockFilePath() string {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
    return filepath.Join(cfg.DBPath, LOCK_FILE)
}

func AcquireLock() error {
	if _, err := os.Stat(getLockFilePath()); os.IsNotExist(err) {
        file, err := os.Create(getLockFilePath())
        if err != nil {
            slog.Error(fmt.Sprintf("Error acquiring lock: %s", err))
            return err
        }
        defer file.Close()

		_, err = file.WriteString(CURRENT_GOROUTINE_ID.String())
		if err != nil {
            slog.Error(fmt.Sprintf("Error acquiring lock: %s", err))
			return err
		}

        slog.Info(fmt.Sprintf("Lock acquired: %s", getLockFilePath()))
    } else {
		return fmt.Errorf("Cannot acquire lock, someone already has it: %s", getLockFilePath())
    }

	return nil
}

func ReleaseLock() error {
	var isLockHeld bool

	if _, err := os.Stat(getLockFilePath()); os.IsNotExist(err) {
        fmt.Println("Lock already released:", getLockFilePath())
		return nil
    } 

	isLockHeld, err := IsLockHeld()
	if err != nil {
		return err
	}

	if isLockHeld {
		err := os.Remove(getLockFilePath())
		if err != nil {
			fmt.Println("Failure releasing lock:", getLockFilePath())
			return err
		}
        slog.Info(fmt.Sprintf("Lock successfully released: %s", getLockFilePath()))
		return nil
	} 

	return errors.New("Cannot release the lock as it's held by another goroutine.")
}

func IsLockHeld() (bool, error) {
	if _, err := os.Stat(getLockFilePath()); os.IsNotExist(err) {
		return false, nil
    } 

	file, err := os.Open(getLockFilePath())
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		var lockHolder string
		lockHolder =  scanner.Text()
		if lockHolder == CURRENT_GOROUTINE_ID.String() {
			return true, nil
		} else {
			return false, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, errors.New("Unexpected error while checking lock.")

}
