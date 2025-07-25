package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/olivierkessler01/gontainers/process"

	//"runtime"

	"github.com/google/uuid"
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc"
)

func serveGRPC(ctx context.Context, cmd *cli.Command) error {
	// Setup and start your gRPC server here
	listener, err := net.Listen("unix", "/var/run/gontainers.sock")
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	process.RegisterMyRuntime(grpcServer)

	slog.Info("Starting gRPC server...")
	return grpcServer.Serve(listener)
}

func main() {
	var logLevel slog.Level = slog.LevelError
	for _, arg := range os.Args {
		if arg == "--verbose" || arg == "-v" {
			logLevel = slog.LevelInfo
			break
		}
	}

	slog.SetLogLoggerLevel(logLevel)
	process.CURRENT_GOROUTINE_ID = uuid.New()

    cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Run a container, get a PID.",
				Action: process.Run,
			},
			{
				Name:  "list",
				Usage: "List containers.",
				Action: process.List,
			},
			{
				Name:  "remove",
				Usage: "Remove a container.",
				Action: process.Remove,
			},
			{
				Name:  "server",
				Usage: "Server the CR-API gRPC server.",
				Action: serveGRPC,
			},
			{
				Name:  "init",
				Usage: "Init the container database.",
				Action: process.Init,
			},
		},
	}

    if err := cmd.Run(context.Background(), os.Args); err != nil {
		slog.Error(fmt.Sprintf("Error: %s", err))
		return
    }
}
