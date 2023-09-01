package restapi

import (
	"encoding/json"
)

// NewSliConfigUnmarshaller creates a new Unmarshaller instance for SliConfigs
func NewSliConfigUnmarshaller() JSONUnmarshaller[*SliConfig] {
	return &sliConfigUnmarshaller{
		tagFilterUnmarshaller: NewTagFilterUnmarshaller(),
	}
}

type sliConfigUnmarshaller struct {
	tagFilterUnmarshaller TagFilterUnmarshaller
}

// UnmarshalArray Unmarshaller interface implementation
func (u *sliConfigUnmarshaller) UnmarshalArray(data []byte) (*[]*SliConfig, error) {
	return unmarshalArray[*SliConfig](data, u.Unmarshal)
}

// Unmarshal Unmarshaller interface implementation
func (u *sliConfigUnmarshaller) Unmarshal(data []byte) (*SliConfig, error) {
	var rawTagFilterExpression json.RawMessage
	var rawGoodEventTagFilterExpression json.RawMessage
	var rawBadEventTagFilterExpression json.RawMessage
	temp := &SliConfig{
		SliEntity: SliEntity{
			FilterExpression:          &rawTagFilterExpression,
			GoodEventFilterExpression: &rawGoodEventTagFilterExpression,
			BadEventFilterExpression:  &rawBadEventTagFilterExpression,
		},
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return &SliConfig{}, err
	}

	filterExpression, err := u.tagFilterUnmarshaller.Unmarshal(rawTagFilterExpression)
	if err != nil {
		return &SliConfig{}, err
	}
	goodEventFilterExpression, err := u.tagFilterUnmarshaller.Unmarshal(rawGoodEventTagFilterExpression)
	if err != nil {
		return &SliConfig{}, err
	}
	badEventFilterExpression, err := u.tagFilterUnmarshaller.Unmarshal(rawBadEventTagFilterExpression)
	if err != nil {
		return &SliConfig{}, err
	}
	temp.SliEntity.FilterExpression = filterExpression
	temp.SliEntity.GoodEventFilterExpression = goodEventFilterExpression
	temp.SliEntity.BadEventFilterExpression = badEventFilterExpression
	return temp, nil
}
