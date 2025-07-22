package main

import (
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"net"
	"olivierkessler01/gontainers/process"
	"os"
)

func serveGRPC(args []string) (int, error) {
	// Setup and start your gRPC server here
	listener, err := net.Listen("unix", "/var/run/gontainers.sock")
	if err != nil {
		return 0, fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	process.RegisterMyRuntime(grpcServer)

	fmt.Println("Starting gRPC server...")
	return 1, grpcServer.Serve(listener)
}

func main() {
	funcMap := map[string]func(args []string) (int, error){
		"run":   process.Run,
		"list":  process.List,
		"serve": serveGRPC,
	}

	process.CURRENT_GOROUTINE_ID = uuid.New()

	var args []string
	args = os.Args[1:]
	_, err := funcMap[args[0]](args[1:])
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
