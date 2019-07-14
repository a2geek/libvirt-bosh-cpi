package cpi

import (
	"libvirt-bosh-cpi/config"
	"libvirt-bosh-cpi/mgr"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"
)

type CPI struct {
	manager mgr.Manager
	uuidGen boshuuid.Generator
	config  config.Config
	logger  boshlog.Logger
}
