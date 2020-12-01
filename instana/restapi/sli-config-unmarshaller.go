package restapi

import (
	"encoding/json"
	"fmt"
)

func NewSliConfigUnmarshaller() Unmarshaller {
	return &sliConfigUnmarshaller{}
}

type sliConfigUnmarshaller struct{}

func (u *sliConfigUnmarshaller) Unmarshal(data []byte) (InstanaDataObject, error) {
	config := SliConfig{}
	if err := json.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("failed to parse json; %s", err)
	}
	return config, nil
}
