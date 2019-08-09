package connection

import (
	"errors"
	"net"
	"time"
)

type socketFactory struct {
	Socket string
}

func (f socketFactory) Validate() error {
	if f.Socket == "" {
		return errors.New("connection type socket requires a socket file")
	}
	return nil
}

func (f socketFactory) Connect() (net.Conn, error) {
	return net.DialTimeout("unix", f.Socket, 2*time.Second)
}
