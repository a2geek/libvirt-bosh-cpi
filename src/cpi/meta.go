package cpi

import (
	"encoding/json"
	"time"

	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

type ActualVMMeta struct {
	Director   string `json:"director"`
	Deployment string `json:"deployment"`
	Name       string `json:"name"`
	Job        string `json:"job"`
	Id         string `json:"id"`
	Index      string `json:"index"` // technically this is a string
}

type ActualDiskMeta struct {
	Director      string    `json:"director"`
	Deployment    string    `json:"deployment"`
	InstanceID    string    `json:"instance_id"`
	Job           string    `json:"job"`
	InstanceIndex string    `json:"instance_index"` // technically this is a string
	InstanceName  string    `json:"instance_name"`
	AttachedAt    time.Time `json:"attached_at"`
}

func NewActualVMMeta(metadata apiv1.VMMeta) (ActualVMMeta, error) {
	bytes, err := metadata.MarshalJSON()
	if err != nil {
		return ActualVMMeta{}, err
	}

	actual := ActualVMMeta{}
	if err = json.Unmarshal(bytes, &actual); err != nil {
		return ActualVMMeta{}, err
	}

	return actual, nil
}

func NewActualDiskMeta(metadata apiv1.DiskMeta) (ActualDiskMeta, error) {
	bytes, err := metadata.MarshalJSON()
	if err != nil {
		return ActualDiskMeta{}, err
	}

	actual := ActualDiskMeta{}
	if err = json.Unmarshal(bytes, &actual); err != nil {
		return ActualDiskMeta{}, err
	}

	return actual, nil
}
