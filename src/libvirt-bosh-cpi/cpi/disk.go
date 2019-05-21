package cpi

import (
	"encoding/xml"
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

func (c CPI) GetDisks(cid apiv1.VMCID) ([]apiv1.DiskCID, error) {
	name := c.vmName(cid.AsString())
	xmlstring, err := c.manager.DomainGetXMLDescByName(name, 0)
	if err != nil {
		return []apiv1.DiskCID{}, bosherr.WrapError(err, "unable to retrieve VM description")
	}

	type DiskDevice struct {
		Type       string `xml:"domain>devices>disk>type,attr"`
		Device     string `xml:"domain>devices>disk>device,attr"`
		SourceFile string `xml:"domain>devices>disk>source>file,attr"`
		TargetDev  string `xml:"domain>devices>disk>target>dev,attr"`
	}

	var disks []DiskDevice
	err = xml.Unmarshal([]byte(xmlstring), &disks)
	if err != nil {
		return []apiv1.DiskCID{}, bosherr.WrapError(err, "unable to unmarshal disk XML")
	}

	var diskcids []apiv1.DiskCID
	for _, disk := range disks {
		if disk.Type == "file" && disk.Device == "disk" {
			vol, err := c.manager.StorageVolLookupByPath(disk.SourceFile)
			if err != nil {
				return []apiv1.DiskCID{}, bosherr.WrapError(err, "unable to locate storage volume")
			}

			diskcids = append(diskcids, apiv1.NewDiskCID(vol.Name))
		}
	}

	return diskcids, nil
}

func (c CPI) CreateDisk(size int,
	cloudProps apiv1.DiskCloudProps, associatedVMCID *apiv1.VMCID) (apiv1.DiskCID, error) {

	uuid, err := c.uuidGen.Generate()
	if err != nil {
		return apiv1.DiskCID{}, bosherr.WrapError(err, "generating uuid")
	}

	name := c.persistantDiskName(uuid)

	_, err = c.manager.CreateStorageVolume(name, uint64(size))
	if err != nil {
		return apiv1.DiskCID{}, bosherr.WrapError(err, "creating disk")
	}

	return apiv1.NewDiskCID(name), nil
}

func (c CPI) DeleteDisk(cid apiv1.DiskCID) error {
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
	return apiv1.NewDiskHintFromString(""), nil
}

func (c CPI) DetachDisk(vmCID apiv1.VMCID, diskCID apiv1.DiskCID) error {
	return nil
}

func (c CPI) HasDisk(cid apiv1.DiskCID) (bool, error) {
	vol, err := c.manager.StorageVolLookupByName(cid.AsString())
	if err != nil {
		return false, bosherr.WrapErrorf(err, "has disk on '%s'", cid.AsString())
	}

	return cid.AsString() == vol.Name, nil
}

func (c CPI) SetDiskMetadata(cid apiv1.DiskCID, metadata apiv1.DiskMeta) error {
	return nil
}

func (c CPI) ResizeDisk(cid apiv1.DiskCID, size int) error {
	return nil
}

func (c CPI) SnapshotDisk(cid apiv1.DiskCID, meta apiv1.DiskMeta) (apiv1.SnapshotCID, error) {
	return apiv1.NewSnapshotCID("snap-cid"), nil
}

func (c CPI) DeleteSnapshot(cid apiv1.SnapshotCID) error {
	return nil
}

func (c CPI) persistantDiskName(cid string) string {
	return fmt.Sprintf("pdisk-%s", cid)
}
