package _go

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"strings"
	"time"
)

func newGoGtagCommand() cli.Command {
	return cli.Command{
		Name:      "tag",
		Usage:     "generate go pseudo-version from current git commit",
		ShortName: "gt",
		Aliases:   []string{"gt", "tag", "gtag"},
		Action: func(c *cli.Context) error {
			commitHash, err := getGitCommit()
			if err != nil {
				return err
			}
			commitTime, err := getGitCommitTime()
			if err != nil {
				return err
			}
			modPath, majorVersion, err := getModuleInfo()
			if err != nil {
				return err
			}
			pseudoVersion := fmt.Sprintf("%s.0.0-%s-%s", majorVersion, commitTime, commitHash[:12])
			fmt.Printf("Module: %s\n", modPath)
			fmt.Println("Pseudo-version:", pseudoVersion)
			return nil
		},
	}
}

func getGitCommit() (string, error) {
	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git commit: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func getGitCommitTime() (string, error) {
	out, err := exec.Command("git", "show", "-s", "--format=%cI", "HEAD").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git commit time: %w", err)
	}

	t, err := time.Parse(time.RFC3339, strings.TrimSpace(string(out)))
	if err != nil {
		return "", fmt.Errorf("failed to parse time: %w", err)
	}

	return t.UTC().Format("20060102150405"), nil
}

func getModuleInfo() (modPath string, majorVersion string, err error) {
	file, err := os.Open("go.mod")
	if err != nil {
		return "", "", fmt.Errorf("failed to open go.mod: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			modPath = strings.TrimSpace(strings.TrimPrefix(line, "module"))
			break
		}
	}

	if modPath == "" {
		return "", "", fmt.Errorf("module path not found in go.mod")
	}

	// 检查是否为 v2+，模块路径末尾可能是 /v2
	parts := strings.Split(modPath, "/")
	last := parts[len(parts)-1]
	if strings.HasPrefix(last, "v") && len(last) > 1 && last[1] >= '2' && last[1] <= '9' {
		majorVersion = last
	} else {
		majorVersion = "v0"
	}

	return modPath, majorVersion, nil
}
