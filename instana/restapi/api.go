package restapi

import "errors"

//InstanaDataObject is a marker interface for any data object provided by any resource of the Instana REST API
type InstanaDataObject interface {
	GetID() string
	Validate() error
}

//RestClient interface to access REST resources of the Instana API
type RestClient interface {
	GetOne(id string, resourcePath string) ([]byte, error)
	GetAll(resourcePath string) ([]byte, error)
	Put(data InstanaDataObject, resourcePath string) ([]byte, error)
	Delete(resourceID string, resourceBasePath string) error
}

//InstanaAPI is the interface to all resources of the Instana Rest API
type InstanaAPI interface {
	Rules() RuleResource
	RuleBindings() RuleBindingResource
}

//ErrEntityNotFound error message which is returned when the entity cannot be found at the server
var ErrEntityNotFound = errors.New("Failed to get resource from Instana API. 404 - Resource not found")
