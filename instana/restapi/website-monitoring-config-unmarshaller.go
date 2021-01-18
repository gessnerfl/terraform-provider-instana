package restapi

import (
	"encoding/json"
	"fmt"
)

//NewWebsiteMonitoringConfigUnmarshaller creates a new Unmarshaller instance for sli configs
func NewWebsiteMonitoringConfigUnmarshaller() Unmarshaller {
	return &websiteMonitoringConfigUnmarshaller{}
}

type websiteMonitoringConfigUnmarshaller struct{}

//Unmarshal Unmarshaller interface implementation
func (u *websiteMonitoringConfigUnmarshaller) Unmarshal(data []byte) (InstanaDataObject, error) {
	config := WebsiteMonitoringConfig{}
	if err := json.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("failed to parse json; %s", err)
	}
	return config, nil
}
