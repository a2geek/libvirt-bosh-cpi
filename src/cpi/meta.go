package cpi

import (
	"encoding/json"

	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

type ActualVMMeta struct {
	Name string `json:"name"`
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
