package _go

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli"
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

func newInitGoEnvSubcommand() cli.Command {
	return cli.Command{
		Name:      "init-go-env",
		Usage:     "initialize go environment",
		ShortName: "i",
		Aliases:   []string{"init-go-env", "i"},
		Action: func(c *cli.Context) error {
			shell := os.Getenv("SHELL")
			profileFile := ""
			switch {
			case strings.Contains(shell, "zsh"):
				profileFile = ".zshrc"
			case strings.Contains(shell, "bash"):
				profileFile = ".bashrc"
			default:
				return cli.NewExitError("unsupported shell: "+shell, 1)
			}
			home, err := os.UserHomeDir()
			if err != nil {
				return cli.NewExitError("failed to get current user: "+err.Error(), 1)
			}
			fullPath := filepath.Join(home, profileFile)
			if alreadyContainsInit(fullPath) {
				fmt.Println("go environment already initialized in", fullPath)
				return nil
			}
			f, err := os.OpenFile(fullPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				return cli.NewExitError("failed to open file: "+err.Error(), 1)
			}
			defer f.Close()
			if _, err := f.WriteString(fmt.Sprintf(goENVInitTextTemplate, defaultGoProxy, defaultGoPath, defaultGoRoot)); err != nil {
				return cli.NewExitError("failed to write to profile: "+err.Error(), 1)
			}
			fmt.Println("✅ Go environment initialized in", fullPath)
			fmt.Println("👉 Please run `source ~/" + profileFile + "` or restart your terminal to apply changes.")
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
