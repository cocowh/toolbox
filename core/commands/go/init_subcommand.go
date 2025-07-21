package _go

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

const goENVInitText = `
# added by toolbox ` + "`toolbox go init-go-env`" + `
export GOPROXY=https://goproxy.cn,direct
export GOPATH=$HOME/go
export GOROOT=/usr/local/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
`

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
			usr, err := user.Current()
			if err != nil {
				return cli.NewExitError("failed to get current user: "+err.Error(), 1)
			}
			fullPath := filepath.Join(usr.HomeDir, profileFile)
			if alreadyContainsInit(fullPath) {
				fmt.Println("go environment already initialized in", fullPath)
				return nil
			}
			f, err := os.OpenFile(fullPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				return cli.NewExitError("failed to open file: "+err.Error(), 1)
			}
			defer f.Close()
			if _, err := f.WriteString(goENVInitText); err != nil {
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
