package web

import (
	"github.com/cocowh/toolbox/core/commands"
	"github.com/urfave/cli"
)

var (
	webCmd cli.Command
)

func init() {
	webCmd = newWebCommand()
	webCmd.Subcommands = []cli.Command{
		newFileServerSubcommand(),
	}
	commands.RegistryCommand(webCmd)
}
func newWebCommand() cli.Command {
	return cli.Command{
		Name:  "web",
		Usage: "web tool",
	}
}
