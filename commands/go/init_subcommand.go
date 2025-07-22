package _go

import (
	"bufio"
	"fmt"
	"github.com/cocowh/toolbox/pkg/logger"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultGoProxy        = "https://goproxy.cn,direct"
	defaultGoPath         = "$HOME/go"
	defaultGoRoot         = "/usr/local/go"
	goENVInitTextTemplate = `
# added by toolbox ` + "`toolbox go init-go-env`" + `
export GOPROXY=%s
export GOPATH=%s
export GOROOT=%s
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
`
)

func newInitGoEnvSubcommand() *cli.Command {
	return &cli.Command{
		Name:    "init-go-env",
		Usage:   "initialize go environment",
		Aliases: []string{"ige", "i"},
		Action: func(c *cli.Context) error {
			shell := os.Getenv("SHELL")
			profileFile := ""
			switch {
			case strings.Contains(shell, "zsh"):
				profileFile = ".zshrc"
			case strings.Contains(shell, "bash"):
				profileFile = ".bashrc"
			default:
				return cli.Exit("unsupported shell: "+shell, 1)
			}
			home, err := os.UserHomeDir()
			if err != nil {
				return cli.Exit("failed to get current user: "+err.Error(), 1)
			}
			fullPath := filepath.Join(home, profileFile)
			if alreadyContainsInit(fullPath) {
				logger.Info("go environment already initialized in %s", fullPath)
				return nil
			}
			f, err := os.OpenFile(fullPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				return cli.Exit("failed to open file: "+err.Error(), 1)
			}
			defer f.Close()
			if _, err := f.WriteString(fmt.Sprintf(goENVInitTextTemplate, defaultGoProxy, defaultGoPath, defaultGoRoot)); err != nil {
				return cli.Exit("failed to write to profile: "+err.Error(), 1)
			}
			logger.Info("Go environment initialized in %s", fullPath)
			logger.Info("Please run `source ~/%s` or restart your terminal to apply changes.", profileFile)
			return nil
		},
	}
}

func alreadyContainsInit(fullPath string) bool {
	f, err := os.Open(fullPath)
	if err != nil {
		return false
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "added by toolbox") {
			return true
		}
	}
	return false
}
