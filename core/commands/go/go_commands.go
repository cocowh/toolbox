package _go

import (
	"github.com/cocowh/toolbox/core/commands"
	"github.com/urfave/cli"
)

var (
	goCmd cli.Command
)

func init() {
	goCmd = newGoCommand()
	goCmd.Subcommands = []cli.Command{
		newInitGoEnvSubcommand(),
		newGoGtagCommand(),
		newInstallGoSubcommand(),
		newChgoSubcommand(),
	}
	commands.RegistryCommand(goCmd)
}

func newGoCommand() cli.Command {
	return cli.Command{
		Name:  "go",
		Usage: "go tool",
	}
}
