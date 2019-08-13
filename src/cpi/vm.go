package cpi

import (
	"fmt"
	"libvirt-bosh-cpi/agentmgr"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

func (c CPI) CreateVM(
	agentID apiv1.AgentID, stemcellCID apiv1.StemcellCID,
	cloudProps apiv1.VMCloudProps, networks apiv1.Networks,
	associatedDiskCIDs []apiv1.DiskCID, env apiv1.VMEnv) (apiv1.VMCID, error) {

	vmCID, _, err := c.CreateVMV2(agentID, stemcellCID, cloudProps, networks, associatedDiskCIDs, env)
	return vmCID, err
}

func (c CPI) CreateVMV2(
	agentID apiv1.AgentID, stemcellCID apiv1.StemcellCID,
	cloudProps apiv1.VMCloudProps, networks apiv1.Networks,
	associatedDiskCIDs []apiv1.DiskCID, env apiv1.VMEnv) (apiv1.VMCID, apiv1.Networks, error) {

	// UUID used for both VM, boot, and ephemeral disk
	uuid, err := c.uuidGen.Generate()
	if err != nil {
		return apiv1.VMCID{}, apiv1.Networks{}, bosherr.WrapError(err, "generating uuid")
	}

	vmName := c.vmName(uuid)

	// Clone stemcell for boot disk
	bootName := c.bootDiskName(uuid)
	stemcellName := c.stemcellName(stemcellCID.AsString())
	_, err = c.manager.CloneStorageVolumeFromStemcell(bootName, stemcellName)
	if err != nil {
		return apiv1.VMCID{}, apiv1.Networks{}, bosherr.WrapErrorf(err, "cloning stemcell '%s'", bootName)
	}

	// Create ephemeral disk
	// BUG? Always creating the ephemeral disk.
	var props LibvirtVMCloudProps
	err = cloudProps.As(&props)
	ephemeralName := c.ephemeralDiskName(uuid)
	ephemeralDiskInBytes := props.EphemeralDisk * bytesPerMegabyte
	_, err = c.manager.CreateStorageVolume(ephemeralName, ephemeralDiskInBytes)
	if err != nil {
		return apiv1.VMCID{}, apiv1.Networks{}, bosherr.WrapErrorf(err, "creating ephemeral disk '%s'", ephemeralName)
	}

	// AgentEnv
	vmCID := apiv1.NewVMCID(vmName)
	agentEnvFactory := apiv1.NewAgentEnvFactory()
	agentEnv := agentEnvFactory.ForVM(agentID, vmCID, networks, env, c.config.Agent)
	agentEnv.AttachSystemDisk(apiv1.NewDiskHintFromString("/dev/vda"))
	agentEnv.AttachEphemeralDisk(apiv1.NewDiskHintFromString("/dev/vdb"))

	// Create config disk
	agentMgr, err := agentmgr.NewAgentManager(c.config.VMPublicKey)
	if err != nil {
		return apiv1.VMCID{}, apiv1.Networks{}, bosherr.WrapError(err, "creating agent manager")
	}
	err = agentMgr.Update(agentEnv)
	if err != nil {
		return apiv1.VMCID{}, apiv1.Networks{}, bosherr.WrapError(err, "updating new config")
	}
	configName := c.configDiskName(uuid)
	configDiskInBytes, err := agentMgr.ToBytes()
	if err != nil {
		return apiv1.VMCID{}, apiv1.Networks{}, bosherr.WrapErrorf(err, "prepping config disk '%s'", configName)
	}
	_, err = c.manager.CreateStorageVolumeFromBytes(configName, configDiskInBytes)
	if err != nil {
		return apiv1.VMCID{}, apiv1.Networks{}, bosherr.WrapErrorf(err, "creating ephemeral disk '%s'", configName)
	}

	// Create VM
	dom, err := c.manager.CreateDomain(vmName, uuid, props.Memory, props.CPU)
	if err != nil {
		return apiv1.VMCID{}, apiv1.Networks{}, bosherr.WrapErrorf(err, "creating domain for vm '%s'", vmName)
	}

	// Attach disks
	if err := c.attachBootDevice(vmName, bootName, "vda"); err != nil {
		return apiv1.VMCID{}, apiv1.Networks{}, bosherr.WrapErrorf(err, "attaching boot disk for vm '%s'", vmName)
	}
	if err := c.attachDiskDevice(vmName, ephemeralName, "vdb"); err != nil {
		return apiv1.VMCID{}, apiv1.Networks{}, bosherr.WrapErrorf(err, "attaching ephemeral disk for vm '%s'", vmName)
	}
	if err := c.attachDiskDevice(vmName, configName, "vdc"); err != nil {
		return apiv1.VMCID{}, apiv1.Networks{}, bosherr.WrapErrorf(err, "attaching config disk for vm '%s'", vmName)
	}

	// Create network interface XML
	for _, network := range networks {
		switch network.Type() {
		case "manual":
			err = c.manager.DomainAttachManualNetworkInterface(dom, network.IP())
			if err != nil {
				return apiv1.VMCID{}, apiv1.Networks{}, bosherr.WrapErrorf(err, "creating network for '%s'", network.IP())
			}
		default:
			return apiv1.VMCID{}, apiv1.Networks{}, bosherr.WrapErrorf(err, "network type '%s' not supported at this time", network.Type())
		}
	}

	if err := c.manager.DomainStart(dom); err != nil {
		return apiv1.VMCID{}, apiv1.Networks{}, bosherr.WrapErrorf(err, "starting domain '%s'", vmName)
	}

	return vmCID, networks, nil
}

func (c CPI) DeleteVM(cid apiv1.VMCID) error {
	dom, err := c.manager.DomainLookupByName(cid.AsString())
	if err != nil {
		return bosherr.WrapError(err, "unable to locate VM")
	}

	diskcids, err := c.discoverDisks(dom)
	if err != nil {
		return bosherr.WrapError(err, "discovering attached disks")
	}

	for _, diskCID := range diskcids {
		err = c.DetachDisk(cid, diskCID)
		if err != nil {
			return bosherr.WrapErrorf(err, "unable to detach disk '%s' from vm '%s'", diskCID.AsString(), cid.AsString())
		}

		if c.isPersistentDisk(cid.AsString()) {
			// Ensure persistent disks are detached but not deleted!
		} else {
			err = c.DeleteDisk(diskCID)
			if err != nil {
				return bosherr.WrapErrorf(err, "unable to delete disk '%s' (was attached to vm '%s')", diskCID.AsString(), cid.AsString())
			}

		}
	}

	err = c.manager.DomainDestroy(cid.AsString())
	if err != nil {
		return bosherr.WrapErrorf(err, "unable to delete vm '%s'", cid.AsString())
	}

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

func (c CPI) bootDiskName(cid string) string {
	return fmt.Sprintf("bdisk-%s", cid)
}

func (c CPI) ephemeralDiskName(cid string) string {
	return fmt.Sprintf("edisk-%s", cid)
}

func (c CPI) configDiskName(cid string) string {
	return fmt.Sprintf("cdisk-%s", cid)
}

func (c CPI) vmName(cid string) string {
	return fmt.Sprintf("vm-%s", cid)
}