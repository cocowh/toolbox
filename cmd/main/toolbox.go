package main

import (
	"github.com/urfave/cli"
	"os"
)

var (
	Version = ""
)

func main() {
	app := newToolboxApp()
	_ = app.Run(os.Args)
}

func newToolboxApp() *cli.App {
	app := &cli.App{
		Name:    "toolbox",
		Usage:   "toolbox is a command line tool that provides a set of tools for developers.",
		Version: Version,
	}
	var commands []*cli.Command
	commands = append(commands, webCmd)
	app.Commands = commands
	return app
}
