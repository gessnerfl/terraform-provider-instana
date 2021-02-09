package restapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

//NewDefaultJSONUnmarshaller creates a new instance of a generic JSONUnmarshaller without specific nested marshalling
func NewDefaultJSONUnmarshaller(objectType interface{}) JSONUnmarshaller {
	if reflect.TypeOf(objectType).Kind() != reflect.Ptr {
		err := errors.New("objectType of defaultJSONUnmarshaller must be a pointer")
		panic(err)
	}
	return &defaultJSONUnmarshaller{
		objectType: objectType,
	}
}

type defaultJSONUnmarshaller struct {
	objectType interface{}
}

//Unmarshal JSONUnmarshaller interface implementation
func (u *defaultJSONUnmarshaller) Unmarshal(data []byte) (interface{}, error) {
	target := u.objectType
	if err := json.Unmarshal(data, &target); err != nil {
		return target, fmt.Errorf("failed to parse json; %s", err)
	}
	return target, nil
}
