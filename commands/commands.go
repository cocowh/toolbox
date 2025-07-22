package commands

import (
	gocmd "github.com/cocowh/toolbox/commands/go"
	"github.com/cocowh/toolbox/commands/web"
	"github.com/urfave/cli/v2"
)

var (
	commands []*cli.Command
)

func init() {
	RegistryCommand(gocmd.NewGoCommand())
	RegistryCommand(web.NewWebCommand())
}

func GetAllCommands() []*cli.Command {
	return commands
}

func RegistryCommand(cmd *cli.Command) {
	commands = append(commands, cmd)
}
