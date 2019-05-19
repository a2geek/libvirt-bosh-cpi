package cpi

import (
	"encoding/json"

	"libvirt-bosh-cpi/connection"

	"github.com/cloudfoundry/bosh-utils/errors"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

type Config struct {
	Connection     connection.Config
	ConnectFactory connection.Factory
	Settings       LibvirtSettings
	Agent          apiv1.AgentOptions
}
type LibvirtSettings struct {
	StoragePoolName string
	StorageVolXml   string
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
	if c.Settings.StoragePoolName == "" {
		return errors.Error("storage pool name is required")
	}

	return nil
}
