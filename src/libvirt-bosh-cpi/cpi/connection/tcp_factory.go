package connection

import (
	"errors"
	"fmt"
	"net"
	"time"
)

type tcpFactory struct {
	Host string
	Port string
}

func (f tcpFactory) Validate() error {
	if f.Host == "" || f.Port == "" {
		return errors.New("connection type tcp requires a host and port specified")
	}
	return nil
}

func (f tcpFactory) Connect() (net.Conn, error) {
	connectString := fmt.Sprintf("%s:%s", f.Host, f.Port)
	return net.DialTimeout("tcp", connectString, 2*time.Second)
}
