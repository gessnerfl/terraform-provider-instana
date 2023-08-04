package restapi

import (
	"encoding/json"
)

// NewWebsiteAlertConfigUnmarshaller creates a new Unmarshaller instance for WebsiteAlertConfigs
func NewWebsiteAlertConfigUnmarshaller() JSONUnmarshaller[*WebsiteAlertConfig] {
	return &websiteAlertConfigUnmarshaller{
		tagFilterUnmarshaller: NewTagFilterUnmarshaller(),
	}
}

type websiteAlertConfigUnmarshaller struct {
	tagFilterUnmarshaller TagFilterUnmarshaller
}

// UnmarshalArray Unmarshaller interface implementation
func (u *websiteAlertConfigUnmarshaller) UnmarshalArray(data []byte) (*[]*WebsiteAlertConfig, error) {
	return unmarshalArray[*WebsiteAlertConfig](data, u.Unmarshal)
}

// Unmarshal Unmarshaller interface implementation
func (u *websiteAlertConfigUnmarshaller) Unmarshal(data []byte) (*WebsiteAlertConfig, error) {
	var rawTagFilterExpression json.RawMessage
	temp := &WebsiteAlertConfig{
		TagFilterExpression: &rawTagFilterExpression,
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return &WebsiteAlertConfig{}, err
	}

	tagFilter, err := u.tagFilterUnmarshaller.Unmarshal(rawTagFilterExpression)
	if err != nil {
		return &WebsiteAlertConfig{}, err
	}
	temp.TagFilterExpression = tagFilter
	return temp, nil
}
