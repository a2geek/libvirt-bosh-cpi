package cpi

import "github.com/cppforlife/bosh-cpi-go/apiv1"

func (c CPI) Info() (apiv1.Info, error) {
	defer c.manager.Disconnect()
	return apiv1.Info{StemcellFormats: c.config.Stemcell.Formats}, nil
}
