package process

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"text/tabwriter"
)

func createContainer(cmd []string) error {
	//TODO:
	return nil
}

func runContainer(cmd []string) (int, error) {
	var cgroups []string
	var namespaces []string

	AcquireLock()
	defer ReleaseLock()
	err := Load()
	if err != nil {
		return 0, err
	}
	
	proc := exec.Command(cmd[0], cmd[1:]...)

	proc.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWUTS |   // new UTS namespace (hostname)
            syscall.CLONE_NEWPID |            // new PID namespace
            syscall.CLONE_NEWNS |             // new mount namespace
            syscall.CLONE_NEWNET,             // new network namespace
	}
	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr

	err = proc.Start()
	if err != nil {
		panic(err)
	}

	namespaces, err = GetNamespaces(proc.Process.Pid)
	if err != nil {
		return 0, err
	}
	
	Add(proc.Process.Pid, cgroups, namespaces)
	err = Save()
	if err != nil {
		panic("We launched a container but can't write it into the DB, wroooong")
	}

	return proc.Process.Pid, err 
}

func removeContainer(pid int) error {
	AcquireLock()
	defer ReleaseLock()

	err := Load()

	if err != nil {
		return err
	}

	if !IsTracked(pid) {
		return fmt.Errorf("Container with PID %d doesn't exist.", pid)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	err = process.Kill()

	if err != nil {
        return fmt.Errorf("Failed to kill process: %s", err)
    } 
    
	fmt.Println("Container killed.")

	delete(TrackedProcesses, pid)

	err = Save()
	if err != nil {
		return err
	}

    return nil
}

func listContainers() error {
	AcquireLock()
	defer ReleaseLock()

	err := Load()

	if err != nil {
		return err
	}

    w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
    fmt.Fprintln(w, "PID\tcgroups\tnamespaces")
    for _, t := range TrackedProcesses {
		row := []string {
			fmt.Sprintf("%d", t.PID),
			strings.Join(t.Cgroups, ","),
			strings.Join(t.Namespaces, ","),
		}
    	fmt.Fprintln(w, strings.Join(row, "\t"))
    }

    w.Flush()
	
    return nil
}
