package web

import (
	"github.com/urfave/cli"
)

func NewWebCommand() cli.Command {
	cmd := cli.Command{
		Name:  "web",
		Usage: "web tool",
	}
	cmd.Subcommands = []cli.Command{
		newFileServerSubcommand(),
	}
	return cmd
}
