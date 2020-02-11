package restapi

import (
	"encoding/json"
	"fmt"
)

//NewCustomEventSpecificationUnmarshaller creates a new instance of Unmarshaller for custom event specifications
func NewCustomEventSpecificationUnmarshaller() Unmarshaller {
	return &customEventSpecificationUnmarshaller{}
}

type customEventSpecificationUnmarshaller struct{}

//Unmarshal Unmarshaller interface implementation
func (u *customEventSpecificationUnmarshaller) Unmarshal(data []byte) (InstanaDataObject, error) {
	customEventSpecification := CustomEventSpecification{}
	if err := json.Unmarshal(data, &customEventSpecification); err != nil {
		return customEventSpecification, fmt.Errorf("failed to parse json; %s", err)
	}
	return customEventSpecification, nil
}
