package cpi_test

import (
	"libvirt-bosh-cpi/cpi"
	"testing"
)

func TestPersistentDiskName(t *testing.T) {
	const uuid = "fake-uuid"
	const expected = "pdisk-" + uuid

	c := cpi.CPI{}
	actual := c.PersistentDiskName(uuid)
	if expected != actual {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func TestIsPersistentDisk(t *testing.T) {
	const goodName = "pdisk-fake-uuid"
	const badName = "edisk-fake-uuid"

	c := cpi.CPI{}
	if c.IsPersistentDisk(goodName) == false {
		t.Errorf("Expected '%s' to be a persistent disk name", goodName)
	}
	if c.IsPersistentDisk(badName) == true {
		t.Errorf("Expected '%s' to not be a persistent disk name", badName)
	}
}
