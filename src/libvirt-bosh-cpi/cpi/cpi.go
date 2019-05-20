package cpi

import (
	"libvirt-bosh-cpi/config"
	"libvirt-bosh-cpi/mgr"

	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"
)

type CPI struct {
	manager mgr.Manager
	uuidGen boshuuid.Generator
	config  config.Config
}
