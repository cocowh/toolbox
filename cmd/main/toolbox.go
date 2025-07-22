package main

import (
	"fmt"
	commands2 "github.com/cocowh/toolbox/commands"
	"github.com/cocowh/toolbox/pkg/logger"
	"github.com/urfave/cli/v2"
	"os"
)

var (
	Version = ""
)

func main() {
	app := newToolboxApp()
	if err := app.Run(os.Args); err != nil {
		logger.Error("Error: %v", err)
		os.Exit(1)
	}
}

func newToolboxApp() *cli.App {
	app := &cli.App{
		Name:    "toolbox",
		Usage:   "Toolbox is a command line tool that provides a set of tools for developers.",
		Version: Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "log-level, ll",
				Value: logger.DebugLevel.ToStringValue(),
				Usage: fmt.Sprintf("set log level: %s", logger.GetAllLogLevelsString()),
			},
			&cli.BoolFlag{
				Name:  "log-hide-tag, lht",
				Usage: "hide log tag",
				Value: true,
			},
		},
		Before: func(c *cli.Context) error {
			logLevel := c.Int("log-level")
			logger.SetLevel(logger.LogLevel(logLevel))
			if c.Bool("log-hide-tag") {
				logger.EnableHideTag()
			}
			return nil
		},
		Commands: commands2.GetAllCommands(),
		ExitErrHandler: func(c *cli.Context, err error) {
			if err != nil {
				logger.Error(err.Error())
			}
			os.Exit(1)
		},
	}
	return app
}
