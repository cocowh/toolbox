package commands

import "github.com/urfave/cli"

var (
	commands []cli.Command
)

func GetAllCommands() []cli.Command {
	return commands
}

func RegistryCommand(cmd cli.Command) {
	commands = append(commands, cmd)
}
