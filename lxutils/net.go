package lxutils
import (
	"net"
	"github.com/layer-x/layerx-commons/lxerrors"
)

func GetLocalIp() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	// handle err
	for _, i := range ifaces {
		if i.Name == "eth1" {
			addrs, err := i.Addrs()
			if err != nil {
				return nil, err
			}
			// handle err
			for _, addr := range addrs {
				switch v := addr.(type) {
				case *net.IPNet:
					return v.IP, nil
				case *net.IPAddr:
					return v.IP, nil
				}
				// process IP address
			}
		}
	}
	return nil, lxerrors.New("Could not find IP in network interfaces", nil)
}