package restapi

import (
	"encoding/json"
	"fmt"
)

//NewAlertingChannelUnmarshaller creates a new Unmarshaller instance for AlertingChannels
func NewAlertingChannelUnmarshaller() Unmarshaller {
	return &alertingChannelUnmarshaller{}
}

type alertingChannelUnmarshaller struct{}

//Unmarshal Unmarshaller interface implementation
func (u *alertingChannelUnmarshaller) Unmarshal(data []byte) (InstanaDataObject, error) {
	alertingChannel := AlertingChannel{}
	if err := json.Unmarshal(data, &alertingChannel); err != nil {
		return alertingChannel, fmt.Errorf("failed to parse json; %s", err)
	}
	return alertingChannel, nil
}
