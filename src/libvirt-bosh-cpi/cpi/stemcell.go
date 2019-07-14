package cpi

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

func (c CPI) CreateStemcell(imagePath string, cloudProps apiv1.StemcellCloudProps) (apiv1.StemcellCID, error) {
	uuid, err := c.uuidGen.Generate()
	if err != nil {
		return apiv1.StemcellCID{}, bosherr.WrapError(err, "generating uuid")
	}

	err = c.writeStemcellCloudProps(uuid, cloudProps)
	if err != nil {
		return apiv1.StemcellCID{}, bosherr.WrapError(err, "writing stemcell cloud properties")
	}

	props, err := c.readStemcellCloudProps(uuid)
	if err != nil {
		return apiv1.StemcellCID{}, bosherr.WrapError(err, "reading stemcell cloud properties")
	}

	name := c.stemcellName(uuid)

	f, err := os.Open(imagePath)
	defer f.Close()
	if err != nil {
		return apiv1.StemcellCID{}, bosherr.WrapErrorf(err, "open imagePath '%s'", imagePath)
	}

	gz, err := gzip.NewReader(f)
	defer gz.Close()
	if err != nil {
		return apiv1.StemcellCID{}, bosherr.WrapErrorf(err, "gunzip imagePath '%s'", imagePath)
	}

	// Assume first file is the one we want. Flexible or failure?
	tf := tar.NewReader(gz)
	_, err = tf.Next()
	if err != nil {
		return apiv1.StemcellCID{}, bosherr.WrapErrorf(err, "tar.next for imagePath '%s'", imagePath)
	}

	_, err = c.manager.CreateStorageVolumeFromImage(name, tf, props.Disk)
	if err != nil {
		return apiv1.StemcellCID{}, bosherr.WrapError(err, "unable to create stemcell storage volume")
	}

	return apiv1.NewStemcellCID(uuid), nil
}

func (c CPI) writeStemcellCloudProps(cid string, cloudProps apiv1.StemcellCloudProps) error {
	var v map[string]interface{}
	cloudProps.As(&v)

	data, err := json.Marshal(v)
	if err != nil {
		return bosherr.WrapError(err, "unable to write stemcell metadata")
	}

	name := c.stemcellMetadataName(cid)

	_, err = c.manager.CreateStorageVolumeFromBytes(name, data)
	return err
}

func (c CPI) readStemcellCloudProps(cid string) (LibvirtStemcellCloudProps, error) {
	name := c.stemcellMetadataName(cid)

	data, err := c.manager.ReadStorageVolumeBytes(name)
	if err != nil {
		return LibvirtStemcellCloudProps{}, bosherr.WrapError(err, "unable to read stemcell metdata")
	}

	var props LibvirtStemcellCloudProps
	err = json.Unmarshal(data, &props)
	if err != nil {
		return LibvirtStemcellCloudProps{}, bosherr.WrapError(err, "unable to unmarshall stemcell metadata")
	}

	return props, nil
}

func (c CPI) DeleteStemcell(cid apiv1.StemcellCID) error {
	name := c.stemcellName(cid.AsString())
	if err := c.manager.StorageVolDeleteByName(name); err != nil {
		return bosherr.WrapErrorf(err, "unable to delete '%s' storage volume", name)
	}

	metadataName := c.stemcellMetadataName(cid.AsString())
	if err := c.manager.StorageVolDeleteByName(metadataName); err != nil {
		return bosherr.WrapErrorf(err, "unable to delete '%s' storage volume metadata", metadataName)
	}

	return nil
}

func (c CPI) stemcellName(cid string) string {
	name := cid
	if !strings.HasPrefix(name, "sc-") {
		name = fmt.Sprintf("sc-%s", name)
	}
	return name
}

func (c CPI) stemcellMetadataName(cid string) string {
	name := fmt.Sprintf("%s.json", cid)
	if !strings.HasPrefix(name, "sc-") {
		name = fmt.Sprintf("sc-%s", name)
	}
	return name
}
