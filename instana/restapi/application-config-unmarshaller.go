package restapi

import (
	"encoding/json"
)

// NewApplicationConfigUnmarshaller creates a new Unmarshaller instance for application configs
func NewApplicationConfigUnmarshaller() JSONUnmarshaller[*ApplicationConfig] {
	return &applicationConfigUnmarshaller{
		tagFilterUnmarshaller: NewTagFilterUnmarshaller(),
	}
}

type applicationConfigUnmarshaller struct {
	tagFilterUnmarshaller TagFilterUnmarshaller
}

// UnmarshalArray Unmarshaller interface implementation
func (u *applicationConfigUnmarshaller) UnmarshalArray(data []byte) (*[]*ApplicationConfig, error) {
	return unmarshalArray[*ApplicationConfig](data, u.Unmarshal)
}

// Unmarshal Unmarshaller interface implementation
func (u *applicationConfigUnmarshaller) Unmarshal(data []byte) (*ApplicationConfig, error) {
	var rawTagFilterExpression json.RawMessage
	temp := &ApplicationConfig{
		TagFilterExpression: &rawTagFilterExpression,
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return &ApplicationConfig{}, err
	}
	tagFilter, err := u.tagFilterUnmarshaller.Unmarshal(rawTagFilterExpression)
	if err != nil {
		return &ApplicationConfig{}, err
	}
	return &ApplicationConfig{
		ID:                  temp.ID,
		Label:               temp.Label,
		TagFilterExpression: tagFilter,
		Scope:               temp.Scope,
		BoundaryScope:       temp.BoundaryScope,
	}, nil
}
