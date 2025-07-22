package _go

import (
	"github.com/urfave/cli"
)

func NewGoCommand() cli.Command {
	cmd := cli.Command{
		Name:  "go",
		Usage: "go tool",
	}
	cmd.Subcommands = []cli.Command{
		newInitGoEnvSubcommand(),
		newGoGtagCommand(),
		newInstallGoSubcommand(),
		newSwitchGoSubcommand(),
	}
	return cmd
}
