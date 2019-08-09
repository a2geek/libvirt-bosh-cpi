package util

import (
	"crypto/rand"
	"net"
)

// GenerateRandomHardwareAddr per https://stackoverflow.com/questions/21018729/generate-mac-address-in-go
func GenerateRandomHardwareAddr() (net.HardwareAddr, error) {
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		return nil, err
	}

	// Set local bit, ensure unicast address
	buf[0] = (buf[0] | 2) & 0xfe

	return net.HardwareAddr(buf), nil
}
