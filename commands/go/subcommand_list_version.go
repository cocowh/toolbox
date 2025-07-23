package _go

import (
	"context"
	"github.com/cocowh/toolbox/pkg/logger"
	"github.com/urfave/cli/v3"
	"os"
	"path"
	"regexp"
	"sort"
)

func newListVersionCommand() *cli.Command {
	return &cli.Command{
		Name:    "list-version",
		Usage:   "List all available golang versions",
		Aliases: []string{"lv", "list-versions"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "dir",
				Usage: "The directory to search for golang versions",
				Value: defaultInstallDir,
			},
		},
		Action: func(ctx context.Context, command *cli.Command) error {
			dir := command.String("dir")
			if dir == "" {
				dir = defaultInstallDir
			}
			re := regexp.MustCompile(`^go(\d+\.\d+(\.\d+)?)$`)
			entries, err := os.ReadDir(dir)
			if err != nil {
				return cli.Exit("failed to read dir: "+dir+" err:"+err.Error(), 1)
			}
			var versions []string
			for _, entry := range entries {
				if !entry.IsDir() {
					continue
				}
				fullPaht := path.Join(dir, entry.Name())
				info, err := os.Lstat(fullPaht)
				if err != nil {
					continue
				}
				if info.Mode()&os.ModeSymlink != 0 {
					continue
				}
				if re.MatchString(entry.Name()) {
					versions = append(versions, entry.Name()[2:])
				}
			}
			logger.Info("Search go version in %s", dir)
			sort.Strings(versions)
			for _, version := range versions {
				logger.Info(version)
			}
			return nil
		},
	}
}
