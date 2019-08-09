package connection

type Config struct {
	Type              string
	Socket            string
	Host              string
	Port              string
	ClientCertificate string
	ClientPrivateKey  string
	CACertificate     string
}
