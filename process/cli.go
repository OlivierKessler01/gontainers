package process

import ( 
	"os"
	"os/exec"
	"strconv"
	"fmt"
	//"runtime"
)


//Run a container
func Run(args []string) ([]TrackedProcess, error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	//runtime.Breakpoint()

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	var processes []TrackedProcess
	processes = append(processes, TrackedProcess{})
	return processes, err 
}

//List containers
func List(args []string) ([]TrackedProcess, error) {

    files, err := os.ReadDir("/proc")
    if err != nil {
        return nil, err
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
    return pids, nil
}



