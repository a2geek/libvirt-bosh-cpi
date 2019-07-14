package cpi

import (
	"libvirt-bosh-cpi/config"
	"libvirt-bosh-cpi/mgr"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	libvirt "github.com/digitalocean/go-libvirt"
)

// Factory implementation.
type Factory struct {
	config config.Config
	logger boshlog.Logger
}

func NewFactory(config config.Config, logger boshlog.Logger) Factory {
	return Factory{config, logger}
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

	m, err := mgr.NewLibvirtManager(l, f.config.Settings, f.logger)
	if err != nil {
		return nil, bosherr.WrapError(err, "failed to create Libvirt manager")
	}

	cpi := CPI{
		manager: m,
		uuidGen: boshuuid.NewGenerator(),
		config:  f.config,
		logger:  f.logger,
	}
	return cpi, nil
}
