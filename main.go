package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/olivierkessler01/gontainers/process"
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc"
)

const GRPC_SOCKET = "/var/run/gontainers.sock"

func serveGRPC(ctx context.Context, cmd *cli.Command) error {
	// Setup and start your gRPC server here
	listener, err := net.Listen("unix", GRPC_SOCKET)
	if err != nil {
		return err
	}

	defer func() {
		listener.Close()
		os.Remove(GRPC_SOCKET)
	}()

	grpcServer := grpc.NewServer()
	process.RegisterMyRuntime(grpcServer)

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

func run(args []string) {
	var logLevel slog.Level = slog.LevelError
	var verboseOutArgs []string

	cancelleableContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	verboseOutArgs = make([]string, 0)
	
	for _, arg := range os.Args {
		if arg == "--verbose=true" || arg == "-v" {
			logLevel = slog.LevelInfo
		} else {
			verboseOutArgs = append(verboseOutArgs, arg)
		}
	}

	args = verboseOutArgs

	slog.SetLogLoggerLevel(logLevel)
	process.CURRENT_GOROUTINE_ID = uuid.New()

	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:   "run",
				Usage:  "Run a container, get a PID.",
				Action: process.Run,
			},
			{
				Name:   "list",
				Usage:  "List containers.",
				Action: process.List,
			},
			{
				Name:   "remove",
				Usage:  "Remove a container.",
				Action: process.Remove,
			},
			{
				Name:   "server",
				Usage:  "Server the CR-API gRPC server.",
				Action: serveGRPC,
			},
			{
				Name:   "init",
				Usage:  "Init the container database.",
				Action: process.Init,
			},
		},
	}

	if err := cmd.Run(cancelleableContext, args); err != nil {
		slog.Error(fmt.Sprintf("%s", err))
		return
	}
}

func main() {
	run(os.Args)
}
