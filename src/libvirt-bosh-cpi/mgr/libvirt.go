package mgr

import (
	"bytes"
	"libvirt-bosh-cpi/config"
	"os"
	"text/template"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/digitalocean/go-libvirt"
)

func NewLibvirtManager(client *libvirt.Libvirt, settings config.LibvirtSettings) (Manager, error) {
	m := libvirtManager{
		client:   client,
		settings: settings,
	}

	if err := m.initialize(); err != nil {
		return m, err
	}
	return m, nil
}

type libvirtManager struct {
	client      *libvirt.Libvirt
	settings    config.LibvirtSettings
	defaultPool libvirt.StoragePool
}

func (m libvirtManager) initialize() error {
	pool, err := m.client.StoragePoolLookupByName(m.settings.StoragePoolName)
	if err != nil {
		return bosherr.WrapErrorf(err, "unable to locate '%s' storage pool", m.settings.StoragePoolName)
	}
	m.defaultPool = pool
	return nil
}

func (m libvirtManager) CreateStorageVolume(name string, size uint64) (libvirt.StorageVol, error) {
	tmpl, err := template.New("storage-volume").Parse(m.settings.StorageVolXml)
	if err != nil {
		return libvirt.StorageVol{}, bosherr.WrapError(err, "unable to parse storage volume XML")
	}

	var xml bytes.Buffer
	tvars := map[string]interface{}{
		"Name": name,
		"Size": size,
		"Unit": "bytes",
	}
	err = tmpl.Execute(&xml, tvars)
	if err != nil {
		return libvirt.StorageVol{}, bosherr.WrapError(err, "unable to generate storage volume XML template")
	}

	vol, err := m.client.StorageVolCreateXML(m.defaultPool, xml.String(), 0)
	if err != nil {
		return libvirt.StorageVol{}, bosherr.WrapError(err, "unable to create storage volume")
	}

	return vol, nil
}

func (m libvirtManager) CreateStorageVolumeFromImage(name, imagePath string) (libvirt.StorageVol, error) {
	finfo, err := os.Stat(imagePath)
	if err != nil {
		return libvirt.StorageVol{}, bosherr.WrapError(err, "determining stemcell size")
	}

	size := finfo.Size()

	vol, err := m.CreateStorageVolume(name, uint64(size))
	if err != nil {
		return libvirt.StorageVol{}, bosherr.WrapError(err, "unable to create storage volume")
	}

	r, err := os.Open(imagePath)
	if err != nil {
		return libvirt.StorageVol{}, bosherr.WrapError(err, "unable to open stemcell file")
	}
	defer r.Close()

	err = m.client.StorageVolUpload(vol, r, 0, uint64(size), 0)
	if err != nil {
		return libvirt.StorageVol{}, bosherr.WrapError(err, "unable to upload stemcell")
	}

	return vol, nil
}

func (m libvirtManager) StorageVolLookupByName(name string) (libvirt.StorageVol, error) {
	return m.client.StorageVolLookupByName(m.defaultPool, name)
}

func (m libvirtManager) StorageVolLookupByPath(path string) (libvirt.StorageVol, error) {
	return m.client.StorageVolLookupByPath(path)
}

func (m libvirtManager) StorageVolGetXMLByName(name string) (string, error) {
	vol, err := m.client.StorageVolLookupByName(m.defaultPool, name)
	if err != nil {
		return "", bosherr.WrapErrorf(err, "unable to locate '%s' storage volume", name)
	}

	return m.client.StorageVolGetXMLDesc(vol, 0)
}

func (m libvirtManager) StorageVolDeleteByName(name string) error {
	vol, err := m.client.StorageVolLookupByName(m.defaultPool, name)
	if err != nil {
		return bosherr.WrapErrorf(err, "unable to locate '%s' storage volume", name)
	}

	err = m.client.StorageVolDelete(vol, 0)
	if err != nil {
		return bosherr.WrapErrorf(err, "unable to delete '%s' storage volume", name)
	}

	return nil
}

func (m libvirtManager) DomainGetXMLDescByName(name string, flags libvirt.DomainXMLFlags) (string, error) {
	vm, err := m.client.DomainLookupByName(name)
	if err != nil {
		return "", bosherr.WrapError(err, "unable to find VM")
	}

	return m.client.DomainGetXMLDesc(vm, 0)
}

func (m libvirtManager) DomainAttachDevice(vmName string, storageVol MgrStorageVol) error {
	vm, err := m.client.DomainLookupByName(vmName)
	if err != nil {
		return bosherr.WrapErrorf(err, "unable to find '%s' VM", vmName)
	}

	tmpl, err := template.New("attach-device").Parse(m.settings.DiskDeviceXml)
	if err != nil {
		return bosherr.WrapError(err, "unable to parse storage volume XML")
	}

	var xml bytes.Buffer
	err = tmpl.Execute(&xml, storageVol)
	if err != nil {
		return bosherr.WrapError(err, "unable to generate storage volume XML template")
	}

	return m.client.DomainAttachDevice(vm, xml.String())
}

func (m libvirtManager) DomainDetachDevice(vmName string, storageVol MgrStorageVol) error {
	vm, err := m.client.DomainLookupByName(vmName)
	if err != nil {
		return bosherr.WrapErrorf(err, "unable to find '%s' VM", vmName)
	}

	tmpl, err := template.New("attach-device").Parse(m.settings.DiskDeviceXml)
	if err != nil {
		return bosherr.WrapError(err, "unable to parse storage volume XML")
	}

	var xml bytes.Buffer
	err = tmpl.Execute(&xml, storageVol)
	if err != nil {
		return bosherr.WrapError(err, "unable to generate storage volume XML template")
	}

	return m.client.DomainDetachDevice(vm, xml.String())
}
