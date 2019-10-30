package cpi

import (
	"encoding/xml"
	"fmt"
	"libvirt-bosh-cpi/mgr"
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	libvirt "github.com/digitalocean/go-libvirt"
)

const bytesPerKilobyte = 1024
const bytesPerMegabyte = bytesPerKilobyte * 1024

func (c CPI) GetDisks(cid apiv1.VMCID) ([]apiv1.DiskCID, error) {
	defer c.manager.Disconnect()
	name := c.vmName(cid.AsString())

	dom, err := c.manager.DomainLookupByName(name)
	if err != nil {
		return []apiv1.DiskCID{}, bosherr.WrapError(err, "unable to locate VM")
	}

	allDisks, err := c.discoverDisks(dom)
	if err != nil {
		return []apiv1.DiskCID{}, bosherr.WrapError(err, "unable to discover disks")
	}

	// Keep Only persistent
	var diskcids []apiv1.DiskCID
	for _, diskcid := range allDisks {
		if c.isPersistentDisk(diskcid.AsString()) {
			diskcids = append(diskcids, diskcid)
		}
	}

	return diskcids, nil
}

func (c CPI) discoverDisks(dom libvirt.Domain) ([]apiv1.DiskCID, error) {
	devices, err := c.manager.DomainListDevices(dom)
	if err != nil {
		return []apiv1.DiskCID{}, bosherr.WrapErrorf(err, "unable to list devices")
	}

	var diskcids []apiv1.DiskCID
	for _, disk := range devices.Disks {
		if disk.Type == "file" && disk.Device == "disk" {
			vol, err := c.manager.StorageVolLookupByPath(disk.Source.File)
			if err != nil {
				return []apiv1.DiskCID{}, bosherr.WrapError(err, "unable to locate storage volume")
			}

			diskcids = append(diskcids, apiv1.NewDiskCID(vol.Name))
		}
	}

	return diskcids, nil
}

func (c CPI) CreateDisk(sizeInMegabytes int,
	cloudProps apiv1.DiskCloudProps, associatedVMCID *apiv1.VMCID) (apiv1.DiskCID, error) {

	defer c.manager.Disconnect()

	uuid, err := c.uuidGen.Generate()
	if err != nil {
		return apiv1.DiskCID{}, bosherr.WrapError(err, "generating uuid")
	}

	sizeInBytes := sizeInMegabytes * bytesPerMegabyte
	name := c.persistentDiskName(uuid)

	_, err = c.manager.CreateStorageVolume(name, uint64(sizeInBytes))
	if err != nil {
		return apiv1.DiskCID{}, bosherr.WrapError(err, "creating disk")
	}

	return apiv1.NewDiskCID(name), nil
}

func (c CPI) DeleteDisk(cid apiv1.DiskCID) error {
	defer c.manager.Disconnect()
	return c.deleteDisk(cid)
}

func (c CPI) deleteDisk(cid apiv1.DiskCID) error {
	if err := c.manager.StorageVolDeleteByName(cid.AsString()); err != nil {
		return bosherr.WrapErrorf(err, "deleting disk '%s'", cid.AsString())
	}

	return nil
}

func (c CPI) AttachDisk(vmCID apiv1.VMCID, diskCID apiv1.DiskCID) error {
	_, err := c.AttachDiskV2(vmCID, diskCID)
	return err
}

func (c CPI) AttachDiskV2(vmCID apiv1.VMCID, diskCID apiv1.DiskCID) (apiv1.DiskHint, error) {
	defer c.manager.Disconnect()

	err := c.attachDiskDevice(vmCID.AsString(), diskCID.AsString(), "vdd")
	if err != nil {
		return apiv1.NewDiskHintFromString(""), bosherr.WrapErrorf(err, "attaching disk '%s' to vm '%s'", diskCID.AsString(), vmCID.AsString())
	}

	diskHint := apiv1.NewDiskHintFromMap(map[string]interface{}{"path": "/dev/vdd"})
	return diskHint, nil
}

func (c CPI) attachDiskDevice(vmName, diskName, deviceName string) error {
	xmlstring, err := c.manager.StorageVolGetXMLByName(diskName)
	if err != nil {
		return bosherr.WrapErrorf(err, "unable to locate storage volume '%s'", diskName)
	}

	var storageVol mgr.StorageVolXml
	err = xml.Unmarshal([]byte(xmlstring), &storageVol)
	if err != nil {
		return bosherr.WrapErrorf(err, "unable to unmarshal storage volume XML: %s", xmlstring)
	}
	if c.isConfigDisk(diskName) && c.config.Stemcell.Type == "CDROM" {
		storageVol.Device = "cdrom"
		storageVol.TargetDevice = "sr0"
		storageVol.TargetBus = "scsi"
	} else {
		storageVol.Device = "disk"
		storageVol.TargetDevice = deviceName
		storageVol.TargetBus = "virtio"
	}

	return c.manager.DomainAttachDisk(vmName, storageVol)
}

func (c CPI) DetachDisk(vmCID apiv1.VMCID, diskCID apiv1.DiskCID) error {
	defer c.manager.Disconnect()
	return c.detachDisk(vmCID, diskCID)
}

func (c CPI) detachDisk(vmCID apiv1.VMCID, diskCID apiv1.DiskCID) error {
	xmlstring, err := c.manager.StorageVolGetXMLByName(diskCID.AsString())
	if err != nil {
		return bosherr.WrapErrorf(err, "detach/unable to locate storage volume '%s'", diskCID.AsString())
	}

	var storageVol mgr.StorageVolXml
	err = xml.Unmarshal([]byte(xmlstring), &storageVol)
	if err != nil {
		return bosherr.WrapError(err, "detach/unable to unmarshal storage volume XML")
	}

	dom, err := c.manager.DomainLookupByName(vmCID.AsString())
	if err != nil {
		return bosherr.WrapError(err, "detach/unable to locate vm")
	}

	devices, err := c.manager.DomainListDevices(dom)
	if err != nil {
		return bosherr.WrapError(err, "detach/unable to list vm devices")
	}

	for _, disk := range devices.Disks {
		if strings.Contains(disk.Source.File, diskCID.AsString()) {
			storageVol.Device = disk.Device
			storageVol.TargetBus = disk.Target.Bus
			storageVol.TargetDevice = disk.Target.Dev
		}
	}

	return c.manager.DomainDetachDisk(vmCID.AsString(), storageVol)
}

func (c CPI) HasDisk(cid apiv1.DiskCID) (bool, error) {
	defer c.manager.Disconnect()

	vol, err := c.manager.StorageVolLookupByName(cid.AsString())
	if err != nil {
		return false, bosherr.WrapErrorf(err, "has disk on '%s'", cid.AsString())
	}

	return cid.AsString() == vol.Name, nil
}

func (c CPI) SetDiskMetadata(cid apiv1.DiskCID, metadata apiv1.DiskMeta) error {
	defer c.manager.Disconnect()
	// Libvirt does not appear to have metadata capability for a storage volume
	return nil
}

func (c CPI) ResizeDisk(cid apiv1.DiskCID, sizeInMegabytes int) error {
	defer c.manager.Disconnect()
	sizeInBytes := sizeInMegabytes * bytesPerMegabyte
	return c.manager.StorageVolResize(cid.AsString(), uint64(sizeInBytes))
}

func (c CPI) SnapshotDisk(cid apiv1.DiskCID, meta apiv1.DiskMeta) (apiv1.SnapshotCID, error) {
	defer c.manager.Disconnect()
	return apiv1.NewSnapshotCID("snap-cid"), nil
}

func (c CPI) DeleteSnapshot(cid apiv1.SnapshotCID) error {
	defer c.manager.Disconnect()
	return nil
}

func (c CPI) persistentDiskName(cid string) string {
	return fmt.Sprintf("pdisk-%s", cid)
}

func (c CPI) isPersistentDisk(cid string) bool {
	return strings.HasPrefix(cid, "pdisk-")
}
