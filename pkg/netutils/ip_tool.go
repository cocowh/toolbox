package netutils

import "net"

func GetLocalIp() ([]string, error) {
	var ips []string
	ifaces, err := net.Interfaces()
	if err != nil {
		return ips, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() || ip.To4() == nil {
				continue
			}
			ips = append(ips, ip.String())
		}
	}
	return ips, nil
}
