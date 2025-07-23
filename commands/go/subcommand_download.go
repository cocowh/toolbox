package _go

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/cocowh/toolbox/pkg/logger"
	"github.com/urfave/cli/v3"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
)

const (
	downloadUrlTemplate = "https://golang.google.cn/dl/%s"
	defaultInstallDir   = "/usr/local/"
)

func newInstallGoSubcommand() *cli.Command {
	return &cli.Command{
		Name:    "install",
		Usage:   "Download and install a specific GO version(requires root permissions)",
		Aliases: []string{"download", "dl", "ins"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "version",
				Usage:    "Go version to download(e.g., 1.22.0)",
				Required: true,
				Aliases:  []string{"v"},
			},
			&cli.StringFlag{
				Name:    "dir",
				Usage:   fmt.Sprintf("Directory to install Go (default: %s)", defaultInstallDir),
				Aliases: []string{"d"},
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			if os.Geteuid() != 0 {
				return cli.Exit("Must run as root, use 'sudo toolbox go install --version=1.22.0'", 1)
			}
			osType := runtime.GOOS
			if osType == "windows" {
				return cli.Exit("Windows is not supported", 1)
			}
			version := c.String("version")
			if version == "" {
				return cli.Exit("Please specify --version (e.g. --version 1.22.0)", 1)
			}
			installBase := c.String("install-dir")
			if installBase == "" {
				installBase = defaultInstallDir
			}
			destDir := filepath.Join(installBase, "go"+version)
			if _, err := os.Stat(destDir); err == nil {
				logger.Info("Go version %s already installed at %s", version, destDir)
				return nil
			}
			arch := runtime.GOARCH
			tarball := fmt.Sprintf("go%s.%s-%s.tar.gz", version, osType, arch)
			url := fmt.Sprintf(downloadUrlTemplate, tarball)

			logger.Info("Downloading: %s", url)
			resp, err := http.Get(url)
			if err != nil {
				return cli.Exit("failed to download Go tarball: "+err.Error(), 1)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return cli.Exit("failed to download Go tarball: "+resp.Status, 1)
			}

			tmpFile := filepath.Join(os.TempDir(), tarball)
			out, err := os.Create(tmpFile)
			if err != nil {
				return cli.Exit("failed to create temp file: "+err.Error(), 1)
			}
			_, err = io.Copy(out, resp.Body)
			out.Close()
			if err != nil {
				return cli.Exit("failed to save tarball: "+err.Error(), 1)
			}

			logger.Info("Extracting to: %s", destDir)
			err = extractTarGz(tmpFile, installBase, version)
			if err != nil {
				return cli.Exit("extraction failed: "+err.Error(), 1)
			}
			logger.Info("Go %s downloaded and extracted to %s", version, destDir)
			return nil
		},
	}
}

func extractTarGz(tarPath, installBase, version string) error {
	defer func() {
		if r := recover(); r != nil {
			logger.Fatal("Extracting tar.gz failed: %s", string(debug.Stack()))
		}
	}()
	f, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		relPath := strings.TrimPrefix(hdr.Name, "go/")
		target := filepath.Join(installBase, "go"+version, relPath)

		if hdr.FileInfo().IsDir() {
			if err := os.MkdirAll(target, hdr.FileInfo().Mode()); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, hdr.FileInfo().Mode())
		if err != nil {
			return err
		}
		if _, err := io.Copy(outFile, tr); err != nil {
			outFile.Close()
			return err
		}
		outFile.Close()
	}
	return nil
}
