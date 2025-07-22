package main

import (
	commands2 "github.com/cocowh/toolbox/commands"
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
		Usage:   "Toolbox is a command line tool that provides a set of tools for developers.",
		Version: Version,
	}
	app.Commands = commands2.GetAllCommands()
	return app
}
