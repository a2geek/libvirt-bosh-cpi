package util

import (
	"crypto/rand"
	"net"
)

// GenerateRandomHardwareAddr per https://stackoverflow.com/questions/21018729/generate-mac-address-in-go
func GenerateRandomHardwareAddr() (net.HardwareAddr, error) {
	buf := make([]byte, 6)
	for {
		_, err := rand.Read(buf)
		if err != nil {
			return nil, err
		}
	
		// Set local bit, ensure unicast address
		buf[0] = (buf[0] | 2) & 0xfe

		// Libvirt apparently uses 0xfe for tap devices; only keeping MAC addresses without 0xfe
		if buf[0] != 0xfe {
			break;
		}
	}

	return net.HardwareAddr(buf), nil
}
