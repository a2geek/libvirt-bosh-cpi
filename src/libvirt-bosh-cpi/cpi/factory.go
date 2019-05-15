package cpi

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	libvirt "github.com/digitalocean/go-libvirt"
)

// Factory implementation.
type Factory struct {
	config Config
}

func NewFactory(config Config) Factory {
	return Factory{config}
}

func (f Factory) New(_ apiv1.CallContext) (apiv1.CPI, error) {
	c, err := f.config.ConnectFactory.Connect()
	if err != nil {
		return nil, bosherr.Errorf("failed to dial libvirt: %v", err)
	}

	l := libvirt.New(c)
	if err := l.Connect(); err != nil {
		return nil, bosherr.Errorf("failed to connect: %v", err)
	}

	cpi := CPI{
		client:  l,
		uuidGen: boshuuid.NewGenerator(),
		config:  f.config,
	}
	return cpi, nil
}
