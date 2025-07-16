package main

import (
	"github.com/cocowh/toolbox/core/server"
	"github.com/urfave/cli"
)

var (
	webCmd = cli.Command{
		Name:  "web",
		Usage: "web tool",
		Subcommands: []cli.Command{
			{
				Name:      "file_serve",
				Usage:     "serve files from a directory",
				ShortName: "fs",
				Aliases:   []string{"fs", "file_serve"},
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "port",
						Value: 8080,
						Usage: "listening port",
					},
					&cli.StringFlag{
						Name:  "dir",
						Usage: "directory to serve files from",
					},
				},
				Action: func(c *cli.Context) error {
					port := c.Int("port")
					dir := c.String("dir")
					return server.StartFileServer(port, dir)
				},
			},
		},
	}
)
