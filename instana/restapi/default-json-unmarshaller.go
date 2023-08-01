package restapi

import (
	"encoding/json"
	"fmt"
)

// NewDefaultJSONUnmarshaller creates a new instance of a generic JSONUnmarshaller without specific nested marshalling
func NewDefaultJSONUnmarshaller[T InstanaDataObject](objectType T) JSONUnmarshaller[T] {
	return &defaultJSONUnmarshaller[T]{
		objectType: objectType,
	}
}

// NewDefaultJSONArrayUnmarshaller creates a new instance of a generic JSONUnmarshaller without specific nested marshalling for an array of the given type
func NewDefaultJSONArrayUnmarshaller[T InstanaDataObject, A *[]T](objectType A) JSONUnmarshaller[A] {
	return &defaultJSONUnmarshaller[A]{
		objectType: objectType,
	}
}

type defaultJSONUnmarshaller[T any] struct {
	objectType T
}

// Unmarshal JSONUnmarshaller interface implementation
func (u *defaultJSONUnmarshaller[T]) Unmarshal(data []byte) (T, error) {
	target := u.objectType
	if err := json.Unmarshal(data, &target); err != nil {
		return target, fmt.Errorf("failed to parse json; %s", err)
	}
	return target, nil
}
