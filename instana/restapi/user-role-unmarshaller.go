package restapi

import (
	"encoding/json"
	"fmt"
)

//NewUserRoleUnmarshaller creates a new Unmarshaller instance for user roles
func NewUserRoleUnmarshaller() Unmarshaller {
	return &userRoleUnmarshaller{}
}

type userRoleUnmarshaller struct{}

//Unmarshal Unmarshaller interface implementation
func (u *userRoleUnmarshaller) Unmarshal(data []byte) (InstanaDataObject, error) {
	userRole := UserRole{}
	if err := json.Unmarshal(data, &userRole); err != nil {
		return userRole, fmt.Errorf("failed to parse json; %s", err)
	}
	return userRole, nil
}
