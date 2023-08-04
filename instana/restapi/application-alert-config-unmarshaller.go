package restapi

import (
	"encoding/json"
)

// NewApplicationAlertConfigUnmarshaller creates a new Unmarshaller instance for ApplicationAlertConfigs
func NewApplicationAlertConfigUnmarshaller() JSONUnmarshaller[*ApplicationAlertConfig] {
	return &applicationAlertConfigUnmarshaller{
		tagFilterUnmarshaller: NewTagFilterUnmarshaller(),
	}
}

type applicationAlertConfigUnmarshaller struct {
	tagFilterUnmarshaller TagFilterUnmarshaller
}

// UnmarshalArray Unmarshaller interface implementation
func (u *applicationAlertConfigUnmarshaller) UnmarshalArray(data []byte) (*[]*ApplicationAlertConfig, error) {
	return unmarshalArray[*ApplicationAlertConfig](data, u.Unmarshal)
}

func (u *applicationAlertConfigUnmarshaller) Unmarshal(data []byte) (*ApplicationAlertConfig, error) {
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
	for i, v := range temp.CustomerPayloadFields {
		temp.CustomerPayloadFields[i] = u.mapCustomPayloadField(v)
	}
	return temp, nil
}

func (u *applicationAlertConfigUnmarshaller) mapCustomPayloadField(field CustomPayloadField[any]) CustomPayloadField[any] {
	if field.Type == DynamicCustomPayloadType {
		data := field.Value.(map[string]interface{})
		var keyPtr *string
		if val, ok := data["key"]; ok {
			key := val.(string)
			keyPtr = &key
		}
		field.Value = DynamicCustomPayloadFieldValue{
			TagName: data["tagName"].(string),
			Key:     keyPtr,
		}
	} else {
		field.Value = StaticStringCustomPayloadFieldValue(field.Value.(string))
	}
	return field
}
