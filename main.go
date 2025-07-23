package main

import (
	"fmt"
	"log/slog"
	"net"
	"olivierkessler01/gontainers/process"
	"os"
	//"runtime"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func serveGRPC(args []string) (int, error) {
	// Setup and start your gRPC server here
	listener, err := net.Listen("unix", "/var/run/gontainers.sock")
	if err != nil {
		return 0, fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	process.RegisterMyRuntime(grpcServer)

	slog.Info("Starting gRPC server...")
	return 1, grpcServer.Serve(listener)
}

func main() {
	funcMap := map[string]func(args []string) (int, error) {
		"run":    process.Run,
		"list":   process.List,
		"remove": process.Remove,
		"serve":  serveGRPC,
		"init":   process.Init,
	}

	var logLevel slog.Level = slog.LevelError
	for _, arg := range os.Args {
		if arg == "--verbose" || arg == "-v" {
			logLevel = slog.LevelInfo
			break
		}
	}

	slog.SetLogLoggerLevel(logLevel)
	//runtime.Breakpoint()

	process.CURRENT_GOROUTINE_ID = uuid.New()

	var args []string
	args = os.Args[1:]
	_, err := funcMap[args[0]](args[1:])
	if err != nil {
		slog.Error(fmt.Sprintf("Error: %s", err))
	}
}
