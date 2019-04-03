package resources

import (
	"encoding/json"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

//NewUserRoleResource constructs a new instance of UserRoleResource
func NewUserRoleResource(client restapi.RestClient) restapi.UserRoleResource {
	return &UserRoleResourceImpl{
		client:       client,
		resourcePath: restapi.UserRolesResourcePath,
	}
}

//UserRoleResourceImpl is the GO representation of the User Role Resource of Instana
type UserRoleResourceImpl struct {
	client       restapi.RestClient
	resourcePath string
}

//GetOne retrieves a single custom User Role from Instana API by its ID
func (resource *UserRoleResourceImpl) GetOne(id string) (restapi.UserRole, error) {
	data, err := resource.client.GetOne(id, resource.resourcePath)
	if err != nil {
		return restapi.UserRole{}, err
	}
	return resource.validateResponseAndConvertToStruct(data)
}

func (resource *UserRoleResourceImpl) validateResponseAndConvertToStruct(data []byte) (restapi.UserRole, error) {
	userRole := restapi.UserRole{}
	if err := json.Unmarshal(data, &userRole); err != nil {
		return userRole, fmt.Errorf("failed to parse json; %s", err)
	}

	if err := userRole.Validate(); err != nil {
		return userRole, err
	}
	return userRole, nil
}

func (resource *UserRoleResourceImpl) validateAllUserRoles(userRoles []restapi.UserRole) error {
	for _, r := range userRoles {
		err := r.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

//Upsert creates or updates a user role
func (resource *UserRoleResourceImpl) Upsert(userRole restapi.UserRole) (restapi.UserRole, error) {
	if err := userRole.Validate(); err != nil {
		return userRole, err
	}
	data, err := resource.client.Put(userRole, resource.resourcePath)
	if err != nil {
		return userRole, err
	}
	return resource.validateResponseAndConvertToStruct(data)
}

//Delete deletes a user role
func (resource *UserRoleResourceImpl) Delete(userRole restapi.UserRole) error {
	return resource.DeleteByID(userRole.ID)
}

//DeleteByID deletes a user role by its ID
func (resource *UserRoleResourceImpl) DeleteByID(userRoleID string) error {
	return resource.client.Delete(userRoleID, resource.resourcePath)
}
