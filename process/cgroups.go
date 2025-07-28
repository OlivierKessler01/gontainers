package process

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// CreateCgroup creates a new cgroup directory under the unified cgroup v2 hierarchy.
func CreateCgroup(name string) (string, error) {
    cgroupRoot := "/sys/fs/cgroup" // typical mountpoint for cgroup v2
    cgroupPath := filepath.Join(cgroupRoot, name)

    if err := os.Mkdir(cgroupPath, 0755); err != nil {
        if os.IsExist(err) {
            // cgroup already exists, not an error for us
            return cgroupPath, nil
        }
        return "", fmt.Errorf("failed to create cgroup directory: %w", err)
    }
    return cgroupPath, nil
}

// AddPidToCgroup writes the PID of the process to the cgroup.procs file.
func AddPidToCgroup(cgroupPath string, pid int) error {
    procsFile := filepath.Join(cgroupPath, "cgroup.procs")
    pidStr := strconv.Itoa(pid)
    return ioutil.WriteFile(procsFile, []byte(pidStr), 0644)
}
