package restapi

import (
	"encoding/json"
	"fmt"
)

//NewSliConfigUnmarshaller creates a new Unmarshaller instance for sli configs
func NewSliConfigUnmarshaller() Unmarshaller {
	return &sliConfigUnmarshaller{}
}

type sliConfigUnmarshaller struct{}

//Unmarshal Unmarshaller interface implementation
func (u *sliConfigUnmarshaller) Unmarshal(data []byte) (InstanaDataObject, error) {
	config := SliConfig{}
	if err := json.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("failed to parse json; %s", err)
	}
	return config, nil
}
