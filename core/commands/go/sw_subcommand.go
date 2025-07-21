package _go

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func newChgoSubcommand() cli.Command {
	return cli.Command{
		Name:      "sw",
		Usage:     "switch the golang version",
		ShortName: "sw",
		Aliases:   []string{"sw", "chgo", "switch"},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "version, v",
				Usage: "Go version to switch to (e.g., 1.22.0)",
			},
			cli.StringFlag{
				Name:  "go-root-dir, d",
				Usage: "Base directory where Go versions are installed (default: $GOPATH)",
			},
		},
		Action: func(c *cli.Context) error {
			version := c.String("version")
			if version == "" {
				return cli.NewExitError("Please specify --version (e.g. --version 1.22.0)", 1)
			}

			currentVersion, err := getCurrentGoVersion()
			if err != nil {
				return cli.NewExitError("Failed to detect current Go version: "+err.Error(), 1)
			}

			if currentVersion == version {
				fmt.Println("Current Go version is already", version)
				return nil
			}

			baseDir := c.String("go-root-dir")
			if baseDir == "" {
				baseDir = os.Getenv("GOPATH")
				if baseDir == "" {
					baseDir = filepath.Join(os.Getenv("HOME"), "go")
				}
			}

			targetDir := filepath.Join(baseDir, "go"+version)

			// Check if target Go directory exists
			if _, err := os.Stat(targetDir); os.IsNotExist(err) {
				return cli.NewExitError("Go version "+version+" not found in "+targetDir, 1)
			}

			// Replace /usr/local/go with symlink
			fmt.Println("Switching Go version to", version, "from", targetDir)

			if err := os.RemoveAll("/usr/local/go"); err != nil {
				return cli.NewExitError("Failed to remove existing /usr/local/go: "+err.Error(), 1)
			}

			if err := os.Symlink(targetDir, "/usr/local/go"); err != nil {
				return cli.NewExitError("Failed to create symlink: "+err.Error(), 1)
			}

			fmt.Println("Go version switched to", version, "successfully.")
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
