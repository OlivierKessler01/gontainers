package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/olivierkessler01/gontainers/process"
	"github.com/urfave/cli/v3"
)

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
				Usage:  "Server the CR-API gRPC server.\n" + 	
						"You can then use `crictl --runtime-endpoint unix:///var/run/gontainers.sock version` to test it.",
				Action: process.ServeGRPC,
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
