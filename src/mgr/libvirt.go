package mgr

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"io"
	"libvirt-bosh-cpi/config"
	"libvirt-bosh-cpi/util"
	"text/template"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/digitalocean/go-libvirt"
)

func NewLibvirtManager(client *libvirt.Libvirt, settings config.LibvirtSettings, logger boshlog.Logger) (Manager, error) {
	pool, err := client.StoragePoolLookupByName(settings.StoragePoolName)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "unable to locate '%s' storage pool", settings.StoragePoolName)
	}

	m := libvirtManager{
		client:      client,
		settings:    settings,
		logger:      logger,
		defaultPool: pool,
	}

	return m, nil
}

type libvirtManager struct {
	client      *libvirt.Libvirt
	settings    config.LibvirtSettings
	logger      boshlog.Logger
	defaultPool libvirt.StoragePool
}

func (m libvirtManager) CreateStorageVolume(name string, sizeInBytes uint64) (libvirt.StorageVol, error) {
	xml, err := m.generateStorageVolumeXML(name, sizeInBytes)
	if err != nil {
		return libvirt.StorageVol{}, bosherr.WrapError(err, "unable to generate storage volume XML")
	}
	m.logger.Debug("libvirt", "CreateStorageVolume XML=%s", xml)

	vol, err := m.client.StorageVolCreateXML(m.defaultPool, xml, 0)
	if err != nil {
		return libvirt.StorageVol{}, bosherr.WrapError(err, "unable to create storage volume of a specific size")
	}
	m.logger.Debug("libvirt", "CreateStorageVolume Volume=%v", vol)

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

	dom, err := m.client.DomainDefineXML(xml.String())
	if err != nil {
		return libvirt.Domain{}, bosherr.WrapError(err, "unable to create VM from domain XML")
	}

	err = m.client.DomainSetAutostart(dom, 1)
	if err != nil {
		return libvirt.Domain{}, bosherr.WrapError(err, "unable to set domain to autostart")
	}

	return dom, nil
}

func (m libvirtManager) DomainStart(dom libvirt.Domain) error {
	active, err := m.client.DomainIsActive(dom)
	if err != nil {
		return bosherr.WrapErrorf(err, "unable to determine if '%s' is running", dom.Name)
	}

	if active == 0 {
		err = m.client.DomainCreate(dom)
		if err != nil {
			return bosherr.WrapErrorf(err, "unable to start '%s'", dom.Name)
		}
	}

	return nil
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

func (m libvirtManager) CreateStorageVolumeFromImage(name string, imageReader io.Reader, diskSizeInBytes uint64) (libvirt.StorageVol, error) {
	vol, err := m.CreateStorageVolume(name, diskSizeInBytes)
	if err != nil {
		return libvirt.StorageVol{}, bosherr.WrapError(err, "unable to create storage volume from image")
	}

	err = m.client.StorageVolUpload(vol, imageReader, 0, 0, 0)
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

	active, err := m.client.DomainIsActive(vm)
	if err != nil {
		return bosherr.WrapErrorf(err, "unable to determine if '%s' is running", vmName)
	}

	flags := libvirt.DomainDeviceModifyConfig
	if active != 0 {
		flags |= libvirt.DomainDeviceModifyLive
	}
	return m.client.DomainAttachDeviceFlags(vm, xml.String(), uint32(flags))
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

	flags := libvirt.DomainDeviceModifyConfig | libvirt.DomainDeviceModifyLive
	return m.client.DomainDetachDeviceFlags(vm, xml.String(), uint32(flags))
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

	err = m.client.DomainDestroy(vm)
	if err != nil {
		return err
	}

	return m.client.DomainUndefine(vm)
}

func (m libvirtManager) DomainAttachManualNetworkInterface(dom libvirt.Domain, ip string) error {
	macaddr, err := util.GenerateRandomHardwareAddr()
	if err != nil {
		return bosherr.WrapError(err, "unable to generate mac address")
	}

	networkDeviceXML, err := m.createNetworkXML(m.settings.ManualNetworkInterfaceXml, dom, ip, macaddr.String())
	if err != nil {
		return bosherr.WrapError(err, "unable to create manual network xml")
	}

	if err := m.client.DomainAttachDeviceFlags(dom, networkDeviceXML, 0); err != nil {
		return bosherr.WrapErrorf(err, "unable to attach network device to domain '%s'", dom.Name)
	}

	if m.settings.NetworkDhcpIpXml != "" {
		networkDhcpXML, err := m.createNetworkXML(m.settings.NetworkDhcpIpXml, dom, ip, macaddr.String())
		if err != nil {
			return bosherr.WrapError(err, "unable to create network DHCP entry")
		}

		net, err := m.client.NetworkLookupByName(m.settings.NetworkName)
		if err != nil {
			return bosherr.WrapErrorf(err, "unable to locate network named '%s'", m.settings.NetworkName)
		}

		cmd := uint32(libvirt.NetworkUpdateCommandAddLast)
		section := uint32(libvirt.NetworkSectionIPDhcpHost)
		// NOTE: Seems to be a bug, argument #2 picts the section to modify. Guessing arg #3 is cmd!
		err = m.client.NetworkUpdate(net, section, cmd, -1, networkDhcpXML, 0)
		if err != nil {
			bosherr.WrapErrorf(err, "unable to attach to network")
		}
	}

	return nil
}

func (m libvirtManager) createNetworkXML(xmlString string, dom libvirt.Domain, ip, mac string) (string, error) {
	tmpl, err := template.New("network-xml").Parse(xmlString)
	if err != nil {
		return "", bosherr.WrapError(err, "unable to parse network XML")
	}

	var xml bytes.Buffer
	tvars := map[string]interface{}{
		"NetworkName": m.settings.NetworkName,
		"IpAddress":   ip,
		"MacAddress":  mac,
		"VmName":      dom.Name,
	}
	err = tmpl.Execute(&xml, tvars)
	if err != nil {
		return "", bosherr.WrapError(err, "unable to generate network XML template")
	}

	return xml.String(), nil
}

func (m libvirtManager) DomainListDevices(dom libvirt.Domain) (DevicesXml, error) {
	xmlstring, err := m.DomainGetXMLDescByName(dom.Name)
	if err != nil {
		return DevicesXml{}, bosherr.WrapError(err, "unable to retrieve VM description")
	}

	var devices DevicesXml
	err = xml.Unmarshal([]byte(xmlstring), &devices)
	if err != nil {
		return DevicesXml{}, bosherr.WrapErrorf(err, "unable to unmarshal devices XML: '%s'", xmlstring)
	}

	return devices, nil
}

func (m libvirtManager) DomainSetDescription(dom libvirt.Domain, description string) error {
	return m.client.DomainSetMetadata(dom, int32(libvirt.DomainMetadataDescription), []string{description}, nil, nil, libvirt.DomainAffectLive|libvirt.DomainAffectConfig)
}

func (m libvirtManager) DomainSetTitle(dom libvirt.Domain, title string) error {
	return m.client.DomainSetMetadata(dom, int32(libvirt.DomainMetadataTitle), []string{title}, nil, nil, libvirt.DomainAffectLive|libvirt.DomainAffectConfig)
}
