package restapi

import (
	"encoding/json"
	"fmt"
)

//NewAlertingConfigurationUnmarshaller creates a new Unmarshaller instance for AlertingConfiguration
func NewAlertingConfigurationUnmarshaller() Unmarshaller {
	return &alertingConfigurationUnmarshaller{}
}

type alertingConfigurationUnmarshaller struct{}

//Unmarshal Unmarshaller interface implementation
func (u *alertingConfigurationUnmarshaller) Unmarshal(data []byte) (InstanaDataObject, error) {
	config := AlertingConfiguration{}
	if err := json.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("failed to parse json; %s", err)
	}
	return config, nil
}
