package web

import (
	"github.com/urfave/cli/v3"
)

func NewWebCommand() *cli.Command {
	cmd := &cli.Command{
		Name:  "web",
		Usage: "web tool",
	}
	cmd.Commands = []*cli.Command{
		newFileServerSubcommand(),
	}
	return cmd
}
