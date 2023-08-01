package restapi

import (
	"encoding/json"
	"fmt"
)

// NewDefaultJSONUnmarshaller creates a new instance of a generic JSONUnmarshaller without specific nested marshalling
func NewDefaultJSONUnmarshaller[T InstanaDataObject](objectType T) JSONUnmarshaller[T] {
	arrayType := make([]T, 0)
	return &defaultJSONUnmarshaller[T]{
		objectType: objectType,
		arrayType:  &arrayType,
	}
}

type defaultJSONUnmarshaller[T any] struct {
	objectType T
	arrayType  *[]T
}

// Unmarshal JSONUnmarshaller interface implementation
func (u *defaultJSONUnmarshaller[T]) Unmarshal(data []byte) (T, error) {
	target := u.objectType
	if err := json.Unmarshal(data, &target); err != nil {
		return target, fmt.Errorf("failed to parse json; %s", err)
	}
	return target, nil
}

// UnmarshalArray JSONUnmarshaller interface implementation
func (u *defaultJSONUnmarshaller[T]) UnmarshalArray(data []byte) (*[]T, error) {
	target := u.arrayType
	if err := json.Unmarshal(data, &target); err != nil {
		return target, fmt.Errorf("failed to parse json; %s", err)
	}
	return target, nil
}
