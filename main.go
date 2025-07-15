package main

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
	"olivierkessler01/gontainers/process"
	"os"
)

func serveGRPC(args []string) error {
	// Setup and start your gRPC server here
	listener, err := net.Listen("unix", "/var/run/gontainers.sock")
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	process.RegisterMyRuntime(grpcServer)

	fmt.Println("Starting gRPC server...")
	return grpcServer.Serve(listener)
}

func main() {
	funcMap := map[string]func(args []string) error{
		"run":   process.Run,
		"list":  process.List,
		"serve": serveGRPC,
	}
	defer process.ReleaseLock()

	var args []string
	args = os.Args[1:]
	err := funcMap[args[0]](args[1:])
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
