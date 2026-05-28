package collector

import (
	netConnections "github.com/shirou/gopsutil/v3/net"
)

func GetOpenPortsCount() (int, error) {

	connections, err := netConnections.Connections("inet")
	if err != nil {
		return 0, err
	}

	portMap := make(map[uint32]bool)

	for _, conn := range connections {

		if conn.Status == "LISTEN" {
			portMap[conn.Laddr.Port] = true
		}
	}

	return len(portMap), nil
}