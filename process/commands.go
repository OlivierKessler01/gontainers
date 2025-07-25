package process

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"text/tabwriter"

	"github.com/urfave/cli/v3"
	"google.golang.org/grpc"
)

const GRPC_SOCKET = "/var/run/gontainers.sock"

// Setup and start your gRPC server
func ServeGRPC(ctx context.Context, cmd *cli.Command) error {
	listener, err := net.Listen("unix", GRPC_SOCKET)
	if err != nil {
		return err
	}

	defer func() {
		listener.Close()
		os.Remove(GRPC_SOCKET)
	}()

	grpcServer := grpc.NewServer()
	RegisterMyRuntime(grpcServer)

	// Run server in background
	errCh := make(chan error, 1)
	go func() {
		slog.Info("Starting gRPC server...")
		errCh <- grpcServer.Serve(listener)
	}()

	// Watch for context cancellation
	select {
	case <-ctx.Done():
		slog.Info("Context canceled. Stopping server...")
		grpcServer.GracefulStop()
		return nil
	case err := <-errCh:
		return fmt.Errorf("gRPC server error: %w", err)
	}
}

//Run a container
func Run(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args()

	AcquireLock()
	defer ReleaseLock()
	err := Load()

	if err != nil {
		return err
	}
	var cgroups []string
	var namespaces []string
	
	var procArgs []string
	procArgs = strings.Fields(args.First())
	proc := exec.Command(procArgs[0], procArgs[1:]...)

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
	
	Add(proc.Process.Pid, cgroups, namespaces)
	err = Save()
	if err != nil {
		panic("We launched a container but can't write it into the DB, wroooong")
	}

	return err 
}

//List containers
func List(ctx context.Context, cmd *cli.Command) error {
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

//Init the database
func Init(ctx context.Context, cmd *cli.Command) error {
	if _, err := os.Stat(getDBFilePath()); os.IsNotExist(err) {
		source, err := os.Open(DB_DEFAULT_FILE)
		if err != nil {
			return err
		}
		defer source.Close()

		destination, err := os.Create(getDBFilePath())
		if err != nil {
			return err
		}
		defer destination.Close()
		_, err = io.Copy(destination, source)

		return err
	} else {
		return fmt.Errorf("Database already initialized.")
	}
}

//Remove containers
func Remove(ctx context.Context, cmd *cli.Command) error {
	AcquireLock()
	defer ReleaseLock()

	err := Load()

	if err != nil {
		return err
	}

	var pid int
	pid, err = strconv.Atoi(cmd.Args().First())
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


