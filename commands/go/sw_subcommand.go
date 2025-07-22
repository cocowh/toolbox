package _go

import (
	"context"
	"fmt"
	"github.com/cocowh/toolbox/pkg/logger"
	"github.com/urfave/cli/v3"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func newSwitchGoSubcommand() *cli.Command {
	return &cli.Command{
		Name:    "sw",
		Usage:   "switch the golang version(requires root permissions)",
		Aliases: []string{"chgo", "switch", "sg", "cg"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "version, v",
				Usage: "Go version to switch to (e.g., 1.22.0)",
			},
			&cli.StringFlag{
				Name:  "install-dir, d",
				Usage: fmt.Sprintf("Base directory where Go versions are installed (default: %s)", defaultInstallDir),
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			if os.Geteuid() != 0 {
				return cli.Exit("Must run as root, use 'sudo toolbox go switch --version=1.22.0'", 1)
			}
			version := c.String("version")
			if version == "" {
				return cli.Exit("Please specify --version (e.g. --version 1.22.0)", 1)
			}

			currentVersion, err := getCurrentGoVersion()
			if err != nil {
				return cli.Exit("Failed to detect current Go version: "+err.Error(), 1)
			}

			if currentVersion == version {
				logger.Info("Current Go version is already %s" + version)
				return nil
			}
			installBase := c.String("install-dir")
			if installBase == "" {
				installBase = defaultInstallDir
			}
			targetDir := filepath.Join(installBase, "go"+version)

			// Check if the target Go directory exists
			if _, err := os.Stat(targetDir); os.IsNotExist(err) {
				return cli.Exit("Go version "+version+" not found in "+targetDir, 1)
			}

			// Replace /usr/local/go with symlink
			logger.Info("Switching Go version to %s from %s", version, targetDir)
			fi, err := os.Lstat(defaultGoRoot)
			if err == nil {
				if fi.Mode()&os.ModeSymlink == 0 {
					return cli.Exit("/usr/local/go is detected and not a soft link. To prevent error deletion, please delete it manually or try again after backup.", 1)
				}
				if err := os.Remove(defaultGoRoot); err != nil {
					return cli.Exit("Failed to remove old soft link `/usr/local/go` : "+err.Error(), 1)
				}
			}

			if err := os.Symlink(targetDir, defaultGoRoot); err != nil {
				return cli.Exit("Failed to create symlink: "+err.Error(), 1)
			}
			logger.Info("Go version switched to %s successfully", version)
			return nil
		},
	}
}

func getCurrentGoVersion() (string, error) {
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Example output: go version go1.22.0 darwin/amd64
	parts := strings.Fields(string(output))
	if len(parts) < 3 {
		return "", fmt.Errorf("unexpected output: %s", string(output))
	}
	version := strings.TrimPrefix(parts[2], "go")
	return version, nil
}
