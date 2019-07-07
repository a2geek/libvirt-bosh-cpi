package mgr

import (
	"bufio"
	"bytes"
	"libvirt-bosh-cpi/config"
	"os"
	"text/template"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/digitalocean/go-libvirt"
)

func NewLibvirtManager(client *libvirt.Libvirt, settings config.LibvirtSettings) (Manager, error) {
	pool, err := client.StoragePoolLookupByName(settings.StoragePoolName)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "unable to locate '%s' storage pool", settings.StoragePoolName)
	}

	m := libvirtManager{
		client:      client,
		settings:    settings,
		defaultPool: pool,
	}

	return m, nil
}

type libvirtManager struct {
	client      *libvirt.Libvirt
	settings    config.LibvirtSettings
	defaultPool libvirt.StoragePool
}

func (m libvirtManager) CreateStorageVolume(name string, sizeInBytes uint64) (libvirt.StorageVol, error) {
	xml, err := m.generateStorageVolumeXML(name, sizeInBytes)
	if err != nil {
		return libvirt.StorageVol{}, bosherr.WrapError(err, "unable to generate storage volume XML")
	}

	vol, err := m.client.StorageVolCreateXML(m.defaultPool, xml, 0)
	if err != nil {
		return libvirt.StorageVol{}, bosherr.WrapError(err, "unable to create storage volume of a specific size")
	}

	return vol, nil
}

func (m libvirtManager) generateStorageVolumeXML(name string, sizeInBytes uint64) (string, error) {
	tmpl, err := template.New("storage-volume").Parse(m.settings.StorageVolXml)
	if err != nil {
		return "", bosherr.WrapError(err, "unable to parse storage volume XML")
	}

	var xml bytes.Buffer
	tvars := map[string]interface{}{
		"Name": name,
		"Size": sizeInBytes,
		"Unit": "bytes",
	}
	err = tmpl.Execute(&xml, tvars)
	if err != nil {
		return "", bosherr.WrapError(err, "unable to generate storage volume XML template")
	}

	return xml.String(), nil
}

func (m libvirtManager) CreateDomain(name, uuid string, memory, cpu uint) (libvirt.Domain, error) {
	tmpl, err := template.New("domain").Parse(m.settings.VmDomainXml)
	if err != nil {
		return libvirt.Domain{}, bosherr.WrapError(err, "unable to parse VM domain XML")
	}

	var xml bytes.Buffer
	tvars := map[string]interface{}{
		"Name":   name,
		"UUID":   uuid,
		"Memory": memory,
		"CPU":    cpu,
	}
	err = tmpl.Execute(&xml, tvars)
	if err != nil {
		return libvirt.Domain{}, bosherr.WrapError(err, "unable to generate VM domain XML template")
	}

	return m.client.DomainCreateXML(xml.String(), 0)
}

func (m libvirtManager) CreateStorageVolumeFromBytes(name string, data []byte) (libvirt.StorageVol, error) {
	sizeInBytes := uint64(len(data))

	vol, err := m.CreateStorageVolume(name, sizeInBytes)
	if err != nil {
		return libvirt.StorageVol{}, bosherr.WrapError(err, "unable to create storage volume from bytes")
	}

	r := bytes.NewReader(data)

	err = m.client.StorageVolUpload(vol, r, 0, sizeInBytes, 0)
	if err != nil {
		return libvirt.StorageVol{}, bosherr.WrapError(err, "unable to upload stemcell")
	}

	return vol, nil
}

func (m libvirtManager) ReadStorageVolumeBytes(name string) ([]byte, error) {
	vol, err := m.StorageVolLookupByName(name)
	if err != nil {
		return nil, bosherr.WrapError(err, "unable to locate volume")
	}

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	err = m.client.StorageVolDownload(vol, w, 0, 0, 0)
	if err != nil {
		return nil, bosherr.WrapError(err, "reading data from volume")
	}

	err = w.Flush()
	if err != nil {
		return nil, bosherr.WrapError(err, "unable to flush data from volume")
	}

	return b.Bytes(), nil
}

func (m libvirtManager) CreateStorageVolumeFromImage(name, imagePath string, diskSizeInBytes uint64) (libvirt.StorageVol, error) {
	vol, err := m.CreateStorageVolume(name, diskSizeInBytes)
	if err != nil {
		return libvirt.StorageVol{}, bosherr.WrapError(err, "unable to create storage volume from image")
	}

	r, err := os.Open(imagePath)
	if err != nil {
		return libvirt.StorageVol{}, bosherr.WrapError(err, "unable to open stemcell file")
	}
	defer r.Close()

	finfo, err := os.Stat(imagePath)
	if err != nil {
		return libvirt.StorageVol{}, bosherr.WrapError(err, "determining stemcell size")
	}

	imageSize := finfo.Size()

	err = m.client.StorageVolUpload(vol, r, 0, uint64(imageSize), 0)
	if err != nil {
		return libvirt.StorageVol{}, bosherr.WrapError(err, "unable to upload stemcell")
	}

	return vol, nil
}

func (m libvirtManager) CloneStorageVolumeFromStemcell(name, stemcell string) (libvirt.StorageVol, error) {
	stemcellVol, err := m.StorageVolLookupByName(stemcell)
	if err != nil {
		return libvirt.StorageVol{}, bosherr.WrapError(err, "unable to locate stemcell")
	}

	xml, err := m.generateStorageVolumeXML(name, 0)
	if err != nil {
		return libvirt.StorageVol{}, bosherr.WrapError(err, "unable to generate storage volume XML")
	}

	vol, err := m.client.StorageVolCreateXMLFrom(m.defaultPool, xml, stemcellVol, 0)
	if err != nil {
		return libvirt.StorageVol{}, bosherr.WrapError(err, "unable to clone stemcell")
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

	err = m.client.StorageVolWipe(vol, 0)
	if err != nil {
		return bosherr.WrapErrorf(err, "unable to wipe '%s' storage volume", name)
	}

	err = m.client.StorageVolDelete(vol, 0)
	if err != nil {
		return bosherr.WrapErrorf(err, "unable to delete '%s' storage volume", name)
	}

	return nil
}

func (m libvirtManager) StorageVolResize(name string, capacityInBytes uint64) error {
	vol, err := m.client.StorageVolLookupByName(m.defaultPool, name)
	if err != nil {
		return bosherr.WrapErrorf(err, "unable to locate '%s' storage volume", name)
	}

	return m.client.StorageVolResize(vol, capacityInBytes, 0)
}

func (m libvirtManager) DomainGetXMLDescByName(name string) (string, error) {
	vm, err := m.client.DomainLookupByName(name)
	if err != nil {
		return "", bosherr.WrapError(err, "unable to find VM")
	}

	return m.client.DomainGetXMLDesc(vm, 0)
}

func (m libvirtManager) DomainAttachBootDisk(vmName string, storageVol StorageVolXml) error {
	vm, err := m.client.DomainLookupByName(vmName)
	if err != nil {
		return bosherr.WrapErrorf(err, "unable to find '%s' VM", vmName)
	}

	tmpl, err := template.New("attach-boot-device").Parse(m.settings.RootDeviceXml)
	if err != nil {
		return bosherr.WrapError(err, "unable to parse root device XML")
	}

	var xml bytes.Buffer
	err = tmpl.Execute(&xml, storageVol)
	if err != nil {
		return bosherr.WrapError(err, "unable to generate root device XML template")
	}

	return m.client.DomainAttachDevice(vm, xml.String())
}

func (m libvirtManager) DomainAttachDisk(vmName string, storageVol StorageVolXml) error {
	vm, err := m.client.DomainLookupByName(vmName)
	if err != nil {
		return bosherr.WrapErrorf(err, "unable to find '%s' VM", vmName)
	}

	tmpl, err := template.New("attach-disk-device").Parse(m.settings.DiskDeviceXml)
	if err != nil {
		return bosherr.WrapError(err, "unable to parse disk device XML")
	}

	var xml bytes.Buffer
	err = tmpl.Execute(&xml, storageVol)
	if err != nil {
		return bosherr.WrapError(err, "unable to generate disk device XML template")
	}

	return m.client.DomainAttachDevice(vm, xml.String())
}

func (m libvirtManager) DomainDetachDisk(vmName string, storageVol StorageVolXml) error {
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

func (m libvirtManager) DomainLookupByName(name string) (libvirt.Domain, error) {
	return m.client.DomainLookupByName(name)
}

func (m libvirtManager) DomainReboot(name string) error {
	vm, err := m.client.DomainLookupByName(name)
	if err != nil {
		return bosherr.WrapErrorf(err, "unable to find '%s' VM", name)
	}
	return m.client.DomainReboot(vm, 0)
}

func (m libvirtManager) DomainDestroy(name string) error {
	vm, err := m.client.DomainLookupByName(name)
	if err != nil {
		return bosherr.WrapErrorf(err, "unable to find vm '%s'", name)
	}

	return m.client.DomainDestroy(vm)
}

func (m libvirtManager) DomainAttachManualNetworkInterface(dom libvirt.Domain, ip string) error {
	networkDeviceXML, err := m.createNetworkXML(m.settings.ManualNetworkInterfaceXml, dom, ip)
	if err != nil {
		return bosherr.WrapError(err, "unable to create manual network xml")
	}

	if err := m.client.DomainAttachDevice(dom, networkDeviceXML); err != nil {
		return bosherr.WrapErrorf(err, "unable to attach network device to domain '%s'", dom.Name)
	}

	if m.settings.NetworkDhcpIpXml != "" {
		networkDhcpXML, err := m.createNetworkXML(m.settings.NetworkDhcpIpXml, dom, ip)
		if err != nil {
			return bosherr.WrapError(err, "unable to create network DHCP entry")
		}

		net, err := m.client.NetworkLookupByName(m.settings.NetworkName)
		if err != nil {
			return bosherr.WrapErrorf(err, "unable to locate network named '%s'", m.settings.NetworkName)
		}

		cmd := uint32(libvirt.NetworkUpdateCommandAddLast)
		section := uint32(libvirt.NetworkSectionIPDhcpHost)
		flags := libvirt.NetworkUpdateAffectLive
		m.client.NetworkUpdate(net, cmd, section, -1, networkDhcpXML, flags)
	}

	return nil
}

func (m libvirtManager) createNetworkXML(xmlString string, dom libvirt.Domain, ip string) (string, error) {
	tmpl, err := template.New("network-xml").Parse(xmlString)
	if err != nil {
		return "", bosherr.WrapError(err, "unable to parse network XML")
	}

	var xml bytes.Buffer
	tvars := map[string]interface{}{
		"NetworkName": m.settings.NetworkName,
		"IpAddress":   ip,
		"VmName":      dom.Name,
	}
	err = tmpl.Execute(&xml, tvars)
	if err != nil {
		return "", bosherr.WrapError(err, "unable to generate network XML template")
	}

	return xml.String(), nil
}
