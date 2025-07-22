package web

import (
	"github.com/cocowh/toolbox/internal/server"
	"github.com/urfave/cli/v2"
)

func newFileServerSubcommand() *cli.Command {
	return &cli.Command{
		Name:    "file_serve",
		Usage:   "Serve files from a directory",
		Aliases: []string{"fs"},
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "port, p",
				Value: 8080,
				Usage: "Listening port",
			},
			&cli.StringFlag{
				Name:  "dir, d",
				Usage: "Directory to serve files from",
			},
		},
		Action: func(c *cli.Context) error {
			return server.StartFileServer(c)
		},
	}
}
