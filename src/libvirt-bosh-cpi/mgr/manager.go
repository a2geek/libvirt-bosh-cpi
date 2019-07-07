package mgr

import (
	"github.com/digitalocean/go-libvirt"
)

type Manager interface {
	CloneStorageVolumeFromStemcell(name, stemcell string) (libvirt.StorageVol, error)
	CreateDomain(name, uuid string, memory, cpu uint) (libvirt.Domain, error)
	CreateStorageVolume(name string, sizeInBytes uint64) (libvirt.StorageVol, error)
	CreateStorageVolumeFromBytes(name string, data []byte) (libvirt.StorageVol, error)
	CreateStorageVolumeFromImage(name, imagePath string, sizeInBytes uint64) (libvirt.StorageVol, error)
	DomainAttachBootDisk(vmName string, storageVol StorageVolXml) error
	DomainAttachDisk(vmName string, storageVol StorageVolXml) error
	DomainAttachManualNetworkInterface(dom libvirt.Domain, ip string) error
	DomainDestroy(name string) error
	DomainDetachDisk(vmName string, storageVol StorageVolXml) error
	DomainGetXMLDescByName(name string) (string, error)
	DomainLookupByName(name string) (libvirt.Domain, error)
	DomainReboot(name string) error
	ReadStorageVolumeBytes(name string) ([]byte, error)
	StorageVolDeleteByName(name string) error
	StorageVolGetXMLByName(name string) (string, error)
	StorageVolLookupByName(name string) (libvirt.StorageVol, error)
	StorageVolLookupByPath(path string) (libvirt.StorageVol, error)
	StorageVolResize(name string, capacity uint64) error
}

type DevicesXml struct {
	Disks []DiskDeviceXml `xml:"devices>disk"`
}
type DiskDeviceXml struct {
	Type   string              `xml:"type,attr"`
	Device string              `xml:"device,attr"`
	Source SourceDiskDeviceXml `xml:"source"`
	Target TargetDiskDeviceXml `xml:"target"`
}
type SourceDiskDeviceXml struct {
	File string `xml:"file,attr"`
}
type TargetDiskDeviceXml struct {
	Dev string `xml:"dev,attr"`
}

type StorageVolXml struct {
	Type         string `xml:"type,attr"`
	Name         string `xml:"name"`
	TargetPath   string `xml:"target>path"`
	TargetDevice string `xml:"-"`
}
