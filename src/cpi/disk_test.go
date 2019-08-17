package cpi

import "testing"

func (c CPI) TestPersistentDiskName(t *testing.T) {
	const uuid = "fake-uuid"
	const expected = "pdisk-" + uuid
	actual := c.persistentDiskName(uuid)
	if expected != actual {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

func (c CPI) TestIsPersistentDisk(t *testing.T) {
	const goodName = "pdisk-fake-uuid"
	const badName = "edisk-fake-uuid"

	if c.isPersistentDisk(goodName) == false {
		t.Errorf("Expected '%s' to be a persistent disk name", goodName)
	}
	if c.isPersistentDisk(badName) == true {
		t.Errorf("Expected '%s' to not be a persistent disk name", badName)
	}
}
