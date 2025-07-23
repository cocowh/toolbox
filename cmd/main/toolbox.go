package main

import (
	"context"
	"fmt"
	commands2 "github.com/cocowh/toolbox/commands"
	"github.com/cocowh/toolbox/pkg/logger"
	"github.com/urfave/cli/v3"
	"os"
)

var (
	Version = ""
)

func main() {
	cmd := newToolboxCommand()
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		logger.Fatal(err.Error())
		os.Exit(1)
	}
}

func newToolboxCommand() *cli.Command {
	app := &cli.Command{
		Name:    "toolbox",
		Usage:   "Toolbox is a command line tool that provides a set of tools for developers.",
		Version: Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log-level",
				Value:   logger.InfoLevel.ToStringValue(),
				Usage:   fmt.Sprintf("set log level: %s", logger.GetAllLogLevelsString()),
				Aliases: []string{"ll"},
			},
			&cli.BoolFlag{
				Name:    "log-show-tag",
				Usage:   "show log tag",
				Value:   false,
				Aliases: []string{"stag"},
			},
			&cli.BoolFlag{
				Name:    "log-show-time",
				Usage:   "show log time",
				Value:   false,
				Aliases: []string{"stime"},
			},
		},
		Before: func(ctx context.Context, command *cli.Command) (context.Context, error) {
			initLogger(command)
			return ctx, nil
		},
		Commands: commands2.GetAllCommands(),
		ExitErrHandler: func(ctx context.Context, command *cli.Command, err error) {
			if err != nil {
				logger.Error(err.Error())
			}
			os.Exit(1)
		},
	}
	return app
}

func initLogger(command *cli.Command) {
	logLevel := command.Int("log-level")
	logger.SetLevel(logger.LogLevel(logLevel))
	if command.Bool("log-show-tag") {
		logger.EnableTag()
	}
	if command.Bool("log-show-time") {
		logger.EnableTimestamp()
	}
}
