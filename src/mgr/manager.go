package mgr

import (
	"io"

	"github.com/digitalocean/go-libvirt"
)

type Manager interface {
	CloneStorageVolumeFromStemcell(name, stemcell string) (libvirt.StorageVol, error)
	CreateDomain(name, uuid string, memory, cpu uint) (libvirt.Domain, error)
	CreateStorageVolume(name string, sizeInBytes uint64) (libvirt.StorageVol, error)
	CreateStorageVolumeFromBytes(name string, data []byte) (libvirt.StorageVol, error)
	CreateStorageVolumeFromImage(name string, image io.Reader, sizeInBytes uint64) (libvirt.StorageVol, error)
	DomainAttachDisk(vmName string, storageVol StorageVolXml) error
	DomainAttachManualNetworkInterface(dom libvirt.Domain, ip string) error
	DomainDestroy(name string) error
	DomainDetachDisk(vmName string, storageVol StorageVolXml) error
	DomainGetXMLDescByName(name string) (string, error)
	DomainListDevices(dom libvirt.Domain) (DevicesXml, error)
	DomainLookupByName(name string) (libvirt.Domain, error)
	DomainReboot(name string) error
	DomainSetDescription(dom libvirt.Domain, description string) error
	DomainSetTitle(dom libvirt.Domain, title string) error
	DomainStart(dom libvirt.Domain) error
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
	Bus string `xml:"bus,attr"`
	Dev string `xml:"dev,attr"`
}

type StorageVolXml struct {
	Type         string `xml:"type,attr"`
	Device       string `xml:"-"`
	Name         string `xml:"name"`
	TargetPath   string `xml:"target>path"`
	TargetBus    string `xml:"-"`
	TargetDevice string `xml:"-"`
}
