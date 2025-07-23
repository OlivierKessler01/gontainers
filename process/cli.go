package process

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"text/tabwriter"
	"io"
)

//Run a container
func Run(args []string) (int, error) {
	AcquireLock()
	defer ReleaseLock()
	err := Load()

	if err != nil {
		return 0, err
	}
	var cgroups []string
	var namespaces []string

	cmd := exec.Command(args[0], args[1:]...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWUTS |   // new UTS namespace (hostname)
            syscall.CLONE_NEWPID |            // new PID namespace
            syscall.CLONE_NEWNS |             // new mount namespace
            syscall.CLONE_NEWNET,             // new network namespace
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	namespaces, err = GetNamespaces(cmd.Process.Pid)
	if err != nil {
		return 0, err
	}
	
	Add(cmd.Process.Pid, cgroups, namespaces)
	err = Save()
	if err != nil {
		panic("We launched a container but can't write it into the DB, wroooong")
	}

	return cmd.Process.Pid, err 
}

//List containers
func List(args []string) (int, error) {
	AcquireLock()
	defer ReleaseLock()

	err := Load()

	if err != nil {
		return 0, err
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
	
    return len(TrackedProcesses), nil
}

//Init the database
func Init(args []string) (int, error) {
	if _, err := os.Stat(getDBFilePath()); os.IsNotExist(err) {
		source, err := os.Open(DB_DEFAULT_FILE)
		if err != nil {
			return 0, err
		}
		defer source.Close()

		destination, err := os.Create(getDBFilePath())
		if err != nil {
			return 0, err
		}
		defer destination.Close()
		_, err = io.Copy(destination, source)

		return 1, err
	} else {
		return 0, fmt.Errorf("Database already initialized.")
	}
}

//Remove containers
func Remove(args []string) (int, error) {
	AcquireLock()
	defer ReleaseLock()

	err := Load()

	if err != nil {
		return 0, err
	}

	var pid int
	pid, err = strconv.Atoi(args[0])
	if err != nil {
		return 0, err
	}

	if !IsTracked(pid) {
		return 0, fmt.Errorf("Container with PID %d doesn't exist.", pid)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return 0, err
	}
	err = process.Kill()

	if err != nil {
        return 0, fmt.Errorf("Failed to kill process: %s", err)
    } 
    
	fmt.Println("Container killed.")

	delete(TrackedProcesses, pid)

	err = Save()
	if err != nil {
		return 0, err
	}

    return 1, nil
}


