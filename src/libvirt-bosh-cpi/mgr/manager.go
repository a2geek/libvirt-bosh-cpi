package mgr

import (
	"github.com/digitalocean/go-libvirt"
)

type Manager interface {
	CreateStorageVolume(name string, size uint64) (libvirt.StorageVol, error)
	CreateStorageVolumeFromImage(name, imagePath string) (libvirt.StorageVol, error)
	DomainGetXMLDescByName(name string, flags libvirt.DomainXMLFlags) (string, error)
	StorageVolDeleteByName(name string) error
	StorageVolLookupByPath(name string) (libvirt.StorageVol, error)
}
