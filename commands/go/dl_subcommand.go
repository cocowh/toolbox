package _go

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/urfave/cli"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	downloadUrlTemplate = "https://golang.google.cn/dl/%s"
	defaultInstallDir   = "/usr/local/"
)

func newInstallGoSubcommand() cli.Command {
	return cli.Command{
		Name:      "install",
		Usage:     "Download and install a specific GO version(requires root permissions)",
		ShortName: "install",
		Aliases:   []string{"download", "dl", "install", "ins"},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "version, v",
				Usage: "Go version to download(e.g., 1.22.0)",
			},
			cli.StringFlag{
				Name:  "install-dir, d",
				Usage: fmt.Sprintf("Directory to install Go (default: %s)", defaultInstallDir),
			},
		},
		Action: func(c *cli.Context) error {
			if os.Geteuid() != 0 {
				return cli.NewExitError("Must run as root, use 'sudo toolbox go install --version=1.22.0'", 1)
			}
			osType := runtime.GOOS
			if osType == "windows" {
				return cli.NewExitError("Windows is not supported", 1)
			}
			version := c.String("version")
			if version == "" {
				return cli.NewExitError("Please specify --version (e.g. --version 1.22.0)", 1)
			}
			installBase := c.String("install-dir")
			if installBase == "" {
				installBase = defaultInstallDir
			}
			destDir := filepath.Join(installBase, "go"+version)
			if _, err := os.Stat(destDir); err == nil {
				fmt.Println("Go version", version, "already installed at", destDir)
				return nil
			}
			arch := runtime.GOARCH
			tarball := fmt.Sprintf("go%s.%s-%s.tar.gz", version, osType, arch)
			url := fmt.Sprintf(downloadUrlTemplate, tarball)

			fmt.Println("Downloading:", url)
			resp, err := http.Get(url)
			if err != nil {
				return fmt.Errorf("failed to download Go tarball: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("failed to download Go tarball: %s", resp.Status)
			}

			tmpFile := filepath.Join(os.TempDir(), tarball)
			out, err := os.Create(tmpFile)
			if err != nil {
				return fmt.Errorf("failed to create temp file: %w", err)
			}
			_, err = io.Copy(out, resp.Body)
			out.Close()
			if err != nil {
				return fmt.Errorf("failed to save tarball: %w", err)
			}

			fmt.Println("Extracting to:", destDir)
			err = extractTarGz(tmpFile, installBase, version)
			if err != nil {
				return fmt.Errorf("extraction failed: %w", err)
			}

			fmt.Println("Go", version, "downloaded and extracted to", destDir)
			return nil
		},
	}
}

func extractTarGz(tarPath, installBase, version string) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Extracting tar.gz failed:", r)
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
