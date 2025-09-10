package process

import (
	"fmt"
	"os"
	"os/exec"
//	"runtime"
	"strings"
	"syscall"
	"text/tabwriter"

	"github.com/google/uuid"
)

type Process struct {
	PID int
	Cmd string
	Args []string
}

type ContainerStatus int

const (
    Created ContainerStatus = iota
   	Running 
   	Exited 
)

type Container struct {
	Name string
	Id string
	Process *Process	
	Cgroups []string
	Namespaces []string
	Status ContainerStatus 
}

func createContainer(name string, cmd []string) (string, error) {
	var id uuid.UUID
	var process Process
	var container *Container

	id = uuid.New()
	if name == "" {
		name = id.String()
	}
	
	AcquireLock()
	defer ReleaseLock()
	err := Load()
	if err != nil {
		return "", err
	}

	_, ok := TrackedContainers[name]
	if ok == true{
		return "", fmt.Errorf("Cannot create container %s, this name is already used.", name)
	}

	process = Process {
		Cmd: cmd[0],
		Args: cmd[1:],
		PID: 0,
	}

	container = &Container {
		Process: &process,
		Status: Created,
		Id: id.String(),
		Name: name,
	}
	TrackedContainers[name] = container

	err = Save()
	if err != nil {
		panic("We create a container but can't write it into the DB, wroooong")
	}

	return container.Id, nil
}

func runContainer(name string, cmd []string) (string, error) {
	var containerId string
	var err error 
	
	containerId, err = createContainer(name, cmd)
	if err != nil {
		return containerId, err
	}
	//runtime.Breakpoint()

	err = startContainer(containerId)
	if err != nil {
		return containerId, err
	}

	return containerId, nil
}


func startContainer(containerId string) error {
	var cgroups []string
	var namespaces []string
	var container *Container
	var exist bool

	AcquireLock()
	defer ReleaseLock()
	err := Load()
	if err != nil {
		return err
	}

	container, exist = TrackedContainers[containerId]
	if exist == false {
		return fmt.Errorf("Cannot start container %s, container does not exist.", containerId)
	}
	
	proc := exec.Command(container.Process.Cmd, container.Process.Args...)

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
		return err
	}

	container.Cgroups = cgroups
	container.Namespaces = namespaces
	
	err = Save()
	if err != nil {
		panic("We started a container but can't write it into the DB, wroooong")
	}

	return nil 
}

func removeContainer(containerId string) error {
	AcquireLock()
	defer ReleaseLock()

	err := Load()

	if err != nil {
		return err
	}

	if !IsTracked(containerId) {
		return fmt.Errorf("Container with PID %s doesn't exist.", containerId)
	}

	pid := TrackedContainers[containerId].Process.PID
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("Cannot find container %s main process with PID %d", containerId, pid)
	}	

	// The main process and all its children, by sending signal to negative PID
	err = syscall.Kill(-process.Pid, syscall.SIGKILL)
	if err != nil {
        return fmt.Errorf("Failed to kill process: %s", err)
    } 
    
	fmt.Println("Container killed.")
	delete(TrackedContainers, containerId)

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
    fmt.Fprintln(w, "ContainerName\tContainerId\tcgroups\tnamespaces")
    for _, t := range TrackedContainers {
		row := []string {
			fmt.Sprintf("%s", t.Name),
			fmt.Sprintf("%s", t.Id),
			strings.Join(t.Cgroups, ","),
			strings.Join(t.Namespaces, ","),
		}
    	fmt.Fprintln(w, strings.Join(row, "\t"))
    }

    w.Flush()
	
    return nil
}
