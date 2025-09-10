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

func run(args []string) (error) {
	var logLevel slog.Level = slog.LevelError

	cancelleableContext, stop := signal.NotifyContext(
		context.Background(), 
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	for _, arg := range os.Args {
		if arg == "--verbose" || arg == "-v" {
			logLevel = slog.LevelInfo
			break
		} 	
	}

	slog.SetLogLoggerLevel(logLevel)
	process.CURRENT_GOROUTINE_ID = uuid.New()

	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Enable verbose logging",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Run a container, get a PID.",
				Action: process.Run,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "name",
						Aliases:  []string{"n"},
						Usage:    "Name of the container",
						Required: true, // ensures user must provide it
					},
					&cli.StringFlag{
						Name:     "command",
						Aliases:  []string{"c"},
						Usage:    "Command to run in the container",
						Required: true,
					},
				},
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
				Name: "server",
				Usage: "Serve the CR-API gRPC server.\n" +
					"You can then use `crictl --runtime-endpoint " +
					"unix:///var/run/gontainers.sock version` to test it.",
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
		return err
	}

	return nil
}

func main() {
	run(os.Args)
}
