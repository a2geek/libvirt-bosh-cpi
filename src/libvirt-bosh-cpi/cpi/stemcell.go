package cpi

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

func (c CPI) CreateStemcell(imagePath string, _ apiv1.StemcellCloudProps) (apiv1.StemcellCID, error) {
	uuid, err := c.uuidGen.Generate()
	if err != nil {
		return apiv1.StemcellCID{}, bosherr.WrapError(err, "generating uuid")
	}

	name := c.stemcellName(uuid)

	_, err = c.manager.CreateStorageVolumeFromImage(name, imagePath)
	if err != nil {
		return apiv1.StemcellCID{}, bosherr.WrapError(err, "unable to create storage volume")
	}

	return apiv1.NewStemcellCID(uuid), nil
}

func (c CPI) DeleteStemcell(cid apiv1.StemcellCID) error {
	name := c.stemcellName(cid.AsString())

	if err := c.manager.StorageVolDeleteByName(name); err != nil {
		return bosherr.WrapErrorf(err, "unable to delete '%s' storage volume", name)
	}

	return nil
}

func (c CPI) stemcellName(cid string) string {
	return fmt.Sprintf("sc-%s", cid)
}
