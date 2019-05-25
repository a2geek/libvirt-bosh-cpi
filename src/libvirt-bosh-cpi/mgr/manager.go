package mgr

import (
	"github.com/digitalocean/go-libvirt"
)

type Manager interface {
	CreateStorageVolume(name string, size uint64) (libvirt.StorageVol, error)
	CreateStorageVolumeFromImage(name, imagePath string) (libvirt.StorageVol, error)
	DomainAttachDevice(vmName string, storageVol StorageVolXml) error
	DomainDetachDevice(vmName string, storageVol StorageVolXml) error
	DomainGetXMLDescByName(name string, flags libvirt.DomainXMLFlags) (string, error)
	StorageVolDeleteByName(name string) error
	StorageVolGetXMLByName(name string) (string, error)
	StorageVolLookupByName(name string) (libvirt.StorageVol, error)
	StorageVolLookupByPath(path string) (libvirt.StorageVol, error)
	StorageVolResize(name string, capacity uint64) error
}

type DiskDeviceXml struct {
	Type       string `xml:"domain>devices>disk>type,attr"`
	Device     string `xml:"domain>devices>disk>device,attr"`
	SourceFile string `xml:"domain>devices>disk>source>file,attr"`
	TargetDev  string `xml:"domain>devices>disk>target>dev,attr"`
}

type StorageVolXml struct {
	Type         string `xml:"volume>type,attr"`
	Name         string `xml:"volume>name"`
	TargetPath   string `xml:"volume>target>path"`
	TargetDevice string `xml:"-"`
}
