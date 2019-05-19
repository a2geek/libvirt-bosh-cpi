package connection

import (
	"net"

	"github.com/pkg/errors"
)

// Factory is an abstraction to support multiple connection types.
type Factory interface {
	Validate() error
	Connect() (net.Conn, error)
}

// NewFactory will setup the appropriate Factory.
func NewFactory(config Config) (Factory, error) {
	switch config.Type {
	case "socket":
		factory := socketFactory{
			Socket: config.Socket,
		}
		return factory, nil
	case "tcp":
		factory := tcpFactory{
			Host: config.Host,
			Port: config.Port,
		}
		return factory, nil
	case "tls":
		factory := tlsFactory{
			Host:              config.Host,
			Port:              config.Port,
			ClientCertificate: config.ClientCertificate,
			ClientPrivateKey:  config.ClientPrivateKey,
			CACertificate:     config.CACertificate,
		}
		return factory, nil
	}
	return nil, errors.Errorf("unknown connection type of '%s'.", config.Type)
}
