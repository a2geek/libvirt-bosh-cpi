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

	return apiv1.NewDiskCID("disk-cid"), nil
}

func (c CPI) DeleteDisk(cid apiv1.DiskCID) error {
	return nil
}

func (c CPI) AttachDisk(vmCID apiv1.VMCID, diskCID apiv1.DiskCID) error {
	return nil
}

func (c CPI) AttachDiskV2(vmCID apiv1.VMCID, diskCID apiv1.DiskCID) (apiv1.DiskHint, error) {
	return apiv1.NewDiskHintFromString(""), nil
}

func (c CPI) DetachDisk(vmCID apiv1.VMCID, diskCID apiv1.DiskCID) error {
	return nil
}

func (c CPI) HasDisk(cid apiv1.DiskCID) (bool, error) {
	return false, nil
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
