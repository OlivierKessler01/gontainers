package process

import ( 
	"os"
	"os/exec"
	"strconv"
	"fmt"
	//"runtime"
)


//Run a container
func Run(args []string) error {
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

    files, err := os.ReadDir("/proc")
    if err != nil {
        return err
    }
    var pids []TrackedProcess
	//runtime.Breakpoint()
    for _, f := range files {
        if f.IsDir() {
            pid, err := strconv.Atoi(f.Name())
            if err == nil {
				pids = append(pids, TrackedProcess{PID:pid})
            }
        }
    }
	
	fmt.Println(pids)
    return nil
}



