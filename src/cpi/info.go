package cpi

import "github.com/cppforlife/bosh-cpi-go/apiv1"

func (c CPI) Info() (apiv1.Info, error) {
	return apiv1.Info{StemcellFormats: []string{"openstack-qcow2", "openstack-raw"}}, nil
}
