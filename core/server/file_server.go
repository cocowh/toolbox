package server

import (
	"fmt"
	"github.com/cocowh/toolbox/utils"
	"github.com/urfave/cli"
	"net/http"
	"os"
	"path/filepath"
)

func StartFileServer(c *cli.Context) error {
	port := c.Int("port")
	dir := c.String("dir")
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("unable to get current working directory: %w", err)
		}
	}
	absDir, _ := filepath.Abs(dir)
	fs := http.FileServer(http.Dir(absDir))
	fmt.Printf("Listening on port %d, serving files from %s\n", port, absDir)
	ips, err := utils.GetLocalIp()
	if err != nil || len(ips) == 0 {
		fmt.Println("Unable to obtain the ip address of the machine")
		fmt.Printf("Default listening address: http://localhost:%d/\n", port)
	} else {
		fmt.Println("You can access it from the following address：")
		for _, ip := range ips {
			fmt.Printf("   http://%s:%d/\n", ip, port)
		}
	}
	return http.ListenAndServe(fmt.Sprintf(":%d", port), fs)
}
