package process

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"text/tabwriter"
)

//Run a container
func Run(args []string) error {
	err := Load()
	if err != nil {
		return err
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
		return err
	}
	
	runtime.Breakpoint()
	Add(cmd.Process.Pid, cgroups, namespaces)
	Save()

	return err 
}

//List containers
func List(args []string) error {
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



