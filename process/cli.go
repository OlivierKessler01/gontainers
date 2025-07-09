package process

import ( 
	"os"
	"os/exec"
	"fmt"
)



//Run a container
func Run(args []string) error {
	Load()
	var cgroups []string
	var namespaces []string

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	//runtime.Breakpoint()

	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	Add(cmd.Process.Pid, cgroups, namespaces)
	Save()

	return err 
}

//List containers
func List(args []string) error {
	Load()
    for _, t := range TrackedProcesses {
		fmt.Println(t.PID)
    }
	
    return nil
}



