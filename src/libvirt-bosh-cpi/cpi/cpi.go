package cpi

import (
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"
	libvirt "github.com/digitalocean/go-libvirt"
)

type CPI struct {
	client  *libvirt.Libvirt
	uuidGen boshuuid.Generator
	config  Config
}
