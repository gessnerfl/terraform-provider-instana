package restapi

import (
	"encoding/json"
)

//NewApplicationAlertConfigUnmarshaller creates a new Unmarshaller instance for ApplicationAlertConfigs
func NewApplicationAlertConfigUnmarshaller() JSONUnmarshaller {
	return &applicationAlertConfigUnmarshaller{
		tagFilterUnmarshaller: NewTagFilterUnmarshaller(),
	}
}

type applicationAlertConfigUnmarshaller struct {
	tagFilterUnmarshaller TagFilterUnmarshaller
}

//Unmarshal Unmarshaller interface implementation
func (u *applicationAlertConfigUnmarshaller) Unmarshal(data []byte) (interface{}, error) {
	var rawTagFilterExpression json.RawMessage
	temp := &ApplicationAlertConfig{
		TagFilterExpression: &rawTagFilterExpression,
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return &ApplicationAlertConfig{}, err
	}

	tagFilter, err := u.tagFilterUnmarshaller.Unmarshal(rawTagFilterExpression)
	if err != nil {
		return &ApplicationAlertConfig{}, err
	}
	temp.TagFilterExpression = tagFilter
	return temp, nil
}
