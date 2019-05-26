package cpi

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

func (c CPI) CreateVM(
	agentID apiv1.AgentID, stemcellCID apiv1.StemcellCID,
	cloudProps apiv1.VMCloudProps, networks apiv1.Networks,
	associatedDiskCIDs []apiv1.DiskCID, env apiv1.VMEnv) (apiv1.VMCID, error) {

	return apiv1.NewVMCID("vm-cid"), nil
}

func (c CPI) CreateVMV2(
	agentID apiv1.AgentID, stemcellCID apiv1.StemcellCID,
	cloudProps apiv1.VMCloudProps, networks apiv1.Networks,
	associatedDiskCIDs []apiv1.DiskCID, env apiv1.VMEnv) (apiv1.VMCID, apiv1.Networks, error) {

	return apiv1.NewVMCID("vm-cid"), networks, nil
}

func (c CPI) DeleteVM(cid apiv1.VMCID) error {
	return nil
}

func (c CPI) CalculateVMCloudProperties(res apiv1.VMResources) (apiv1.VMCloudProps, error) {
	props := make(map[string]interface{})
	props["cpu"] = res.CPU
	props["memory"] = res.RAM
	props["ephemeral_disk"] = res.EphemeralDiskSize
	return apiv1.NewVMCloudPropsFromMap(props), nil
}

func (c CPI) SetVMMetadata(cid apiv1.VMCID, metadata apiv1.VMMeta) error {
	return nil
}

func (c CPI) HasVM(cid apiv1.VMCID) (bool, error) {
	vm, err := c.manager.DomainLookupByName(cid.AsString())
	if err != nil {
		return false, bosherr.WrapErrorf(err, "unable to find '%s' VM", cid.AsString())
	}

	return cid.AsString() != vm.Name, nil
}

func (c CPI) RebootVM(cid apiv1.VMCID) error {
	return c.manager.DomainReboot(cid.AsString())
}

func (c CPI) ephemeralDiskName(cid string) string {
	return fmt.Sprintf("edisk-%s", cid)
}

func (c CPI) vmName(cid string) string {
	return fmt.Sprintf("vm-%s", cid)
}
