package process

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"strings"
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
	var command []string
	var name string
	command = strings.Fields(cmd.String("command"))
	name = cmd.String("name")
	id, err := runContainer(name, command)
	if err == nil {
		fmt.Println("Container launched %d", id)
	}
	return err 
}

//List containers
func List(ctx context.Context, cmd *cli.Command) error {
	return listContainers()
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
	containerId := cmd.Args().First()
	return removeContainer(containerId)
}


