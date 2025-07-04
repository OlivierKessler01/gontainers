package process

import ( 
	"os"
	"os/exec"
	//"runtime"
)

func run(args []string) bool {
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

func List(args []string) bool {
	run(args)
	return true
}

