package process

import ( 
	"os"
	"os/exec"
	"strconv"
)

type TrackedProcess struct {
	PID int
	cgroups []string
	namespaces []string
}

func Run(args []string) bool {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	//runtime.Breakpoint()

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	return true
}

func List(args []string) ([]TrackedProcess, error) {
    files, err := os.ReadDir("/proc")
    if err != nil {
        return nil, err
    }

    var pids []TrackedProcess
    for _, f := range files {
        if f.IsDir() {
            pid, err := strconv.Atoi(f.Name())
            if err == nil {
				pids = append(pids, TrackedProcess{PID:pid})
            }
        }
    }

    return pids, nil
}



