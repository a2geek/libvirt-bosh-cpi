package config

import (
	"encoding/json"

	"libvirt-bosh-cpi/connection"

	"github.com/cloudfoundry/bosh-utils/errors"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

type Config struct {
	Agent          apiv1.AgentOptions
	ConnectFactory connection.Factory
	Connection     connection.Config
	Settings       LibvirtSettings
}
type LibvirtSettings struct {
	DiskDeviceXml             string
	ManualNetworkInterfaceXml string
	NetworkName               string
	NetworkDhcpIpXml          string
	RootDeviceXml             string
	StoragePoolName           string
	StorageVolXml             string
	VmDomainXml               string
}

func NewConfigFromPath(path string, fs boshsys.FileSystem) (Config, error) {
	// This includes any default values
	config := Config{}

	bytes, err := fs.ReadFile(path)
	if err != nil {
		return config, bosherr.WrapErrorf(err, "reading config '%s'", path)
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return config, bosherr.WrapError(err, "unmarshalling config")
	}

	cf, err := connection.NewFactory(config.Connection)
	if err != nil {
		return config, bosherr.WrapError(err, "identifying connection factory")
	}

	err = cf.Validate()
	if err != nil {
		return config, bosherr.WrapError(err, "validating connection config")
	}

	err = config.Validate()
	if err != nil {
		return config, bosherr.WrapError(err, "validating config")
	}

	config.ConnectFactory = cf
	return config, nil
}

func (c Config) Validate() error {
	var e []error

	if c.Settings.DiskDeviceXml == "" {
		e = append(e, errors.Error("disk device xml is required"))
	}

	if c.Settings.ManualNetworkInterfaceXml == "" {
		e = append(e, errors.Error("manual network interface xml is required"))
	}

	if c.Settings.NetworkName == "" {
		e = append(e, errors.Error("network name is required"))
	}

	if c.Settings.RootDeviceXml == "" {
		e = append(e, errors.Error("root device xml is required"))
	}

	if c.Settings.StoragePoolName == "" {
		e = append(e, errors.Error("storage pool name is required"))
	}

	if c.Settings.StorageVolXml == "" {
		e = append(e, errors.Error("storage volume xml is required"))
	}

	if c.Settings.VmDomainXml == "" {
		e = append(e, errors.Error("domain xml is required"))
	}

	if len(e) > 0 {
		return errors.NewMultiError(e...)
	}

	return nil
}
