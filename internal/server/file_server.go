package server

import (
	"fmt"
	"github.com/cocowh/toolbox/pkg/logger"
	"github.com/cocowh/toolbox/pkg/netutils"
	"github.com/urfave/cli/v3"
	"net/http"
	"os"
	"path/filepath"
)

func StartFileServer(c *cli.Command) error {
	port := c.Int("port")
	dir := c.String("dir")
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return cli.Exit("unable to get current working directory: "+err.Error(), 1)
		}
	}
	absDir, _ := filepath.Abs(dir)
	fs := http.FileServer(http.Dir(absDir))
	logger.Info("Listening on port %d, serving files from %s", port, absDir)
	ips, err := netutils.GetLocalIp()
	if err != nil || len(ips) == 0 {
		logger.Warn("Unable to obtain the ip address of the machine")
		logger.Warn("Default listening address: http://localhost:%d/", port)
	} else {
		logger.Info("You can access it from the following address：")
		for _, ip := range ips {
			logger.Info("   http://%s:%d/", ip, port)
		}
	}
	return http.ListenAndServe(fmt.Sprintf(":%d", port), fs)
}
