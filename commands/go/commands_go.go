package _go

import (
	"github.com/urfave/cli/v3"
)

func NewGoCommand() *cli.Command {
	cmd := &cli.Command{
		Name:  "go",
		Usage: "go tool",
	}
	cmd.Commands = []*cli.Command{
		newInitGoEnvSubcommand(),
		newGoGtagCommand(),
		newInstallGoSubcommand(),
		newSwitchGoSubcommand(),
		newListVersionCommand(),
	}
	return cmd
}
