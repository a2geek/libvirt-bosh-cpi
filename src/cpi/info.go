package cpi

import "github.com/cppforlife/bosh-cpi-go/apiv1"

func (c CPI) Info() (apiv1.Info, error) {
	return apiv1.Info{StemcellFormats: c.config.Stemcell.Formats}, nil
}
