package cpi

import (
	"fmt"
	"os"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

func (c CPI) CreateStemcell(imagePath string, _ apiv1.StemcellCloudProps) (apiv1.StemcellCID, error) {
	pool, err := c.client.StoragePoolLookupByName(c.config.Settings.StoragePoolName)
	if err != nil {
		return apiv1.StemcellCID{}, bosherr.WrapErrorf(err, "unable to locate '%s' storage pool", c.config.Settings.StoragePoolName)
	}

	uuid, err := c.uuidGen.Generate()
	if err != nil {
		return apiv1.StemcellCID{}, bosherr.WrapError(err, "generating uuid")
	}

	finfo, err := os.Stat(imagePath)
	if err != nil {
		return apiv1.StemcellCID{}, bosherr.WrapError(err, "determining stemcell size")
	}

	name := c.stemcellName(uuid)
	size := finfo.Size()

	vol, err := c.createStorageVolume(pool, name, uint64(size))
	if err != nil {
		return apiv1.StemcellCID{}, bosherr.WrapError(err, "unable to create storage volume")
	}

	r, err := os.Open(imagePath)
	if err != nil {
		return apiv1.StemcellCID{}, bosherr.WrapError(err, "unable to open stemcell file")
	}

	err = c.client.StorageVolUpload(vol, r, 0, uint64(size), 0)
	if err != nil {
		return apiv1.StemcellCID{}, bosherr.WrapError(err, "unable to upload stemcell")
	}

	return apiv1.NewStemcellCID(uuid), nil
}

func (c CPI) DeleteStemcell(cid apiv1.StemcellCID) error {
	pool, err := c.client.StoragePoolLookupByName(c.config.Settings.StoragePoolName)
	if err != nil {
		return bosherr.WrapErrorf(err, "unable to locate '%s' storage pool", c.config.Settings.StoragePoolName)
	}

	name := c.stemcellName(cid.AsString())

	vol, err := c.client.StorageVolLookupByName(pool, name)
	if err != nil {
		return bosherr.WrapErrorf(err, "unable to locate '%s' storage volume", name)
	}

	err = c.client.StorageVolDelete(vol, 0)
	if err != nil {
		return bosherr.WrapErrorf(err, "unable to delete '%s' storage volume", name)
	}

	return nil
}

func (c CPI) stemcellName(cid string) string {
	return fmt.Sprintf("sc-%s", cid)
}
