package connection

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
)

type tlsFactory struct {
	Host              string
	Port              string
	ClientCertificate string
	ClientPrivateKey  string
	CACertificate     string
}

func (f tlsFactory) Validate() error {
	if f.Host == "" || f.Port == "" || f.ClientCertificate == "" || f.ClientPrivateKey == "" || f.CACertificate == "" {
		return errors.New("connection type tls requires a host, port, client certificate/private key, and a CA certificate specified")
	}
	return nil
}

func (f tlsFactory) Connect() (net.Conn, error) {
	cert, err := tls.X509KeyPair([]byte(f.ClientCertificate), []byte(f.ClientPrivateKey))
	if err != nil {
		return nil, err
	}

	roots := x509.NewCertPool()
	roots.AppendCertsFromPEM([]byte(f.CACertificate))

	cfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      roots,
	}

	connectString := fmt.Sprintf("%s:%s", f.Host, f.Port)
	return tls.Dial("tcp", connectString, cfg)
}
