package agentmgr

import (
	"fmt"

	"libvirt-bosh-cpi/config"

	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

// AgentManager is an abstraction to update the AgentEnv into a VM
type AgentManager interface {
	Update(apiv1.AgentEnv) error
	ToBytes() ([]byte, error)
}

// NewAgentManager will initalize a new config drive for AgentEnv settings
func NewAgentManager(config config.Config) (AgentManager, error) {
	var a AgentManager
	var err error
	switch config.Stemcell.Type {
	case "ConfigDrive":
		a, err = NewConfigDriveManager(config)
	}
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, fmt.Errorf("Unknown stemcell configuration type '%s'", config.Stemcell.Type)
	}
	return a, nil
}

// NewAgentManagerFromData will allow AgentEnv settings updates on an existing config drive
func NewAgentManagerFromData(config config.Config, data []byte) (AgentManager, error) {
	var a AgentManager
	var err error
	switch config.Stemcell.Type {
	case "ConfigDrive":
		a, err = NewConfigDriveManagerFromData(config, data)
	}
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, fmt.Errorf("Unknown stemcell configuration type '%s'", config.Stemcell.Type)
	}
	return a, nil
}
