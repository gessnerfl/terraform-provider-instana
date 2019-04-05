package resources

import (
	"encoding/json"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

//NewApplicationConfigResource constructs a new instance of ApplicationConfigResource
func NewApplicationConfigResource(client restapi.RestClient) restapi.ApplicationConfigResource {
	return &ApplicationConfigResourceImpl{
		client:       client,
		resourcePath: restapi.ApplicationConfigsResourcePath,
	}
}

//ApplicationConfigResourceImpl is the GO representation of the application config API of the Instana
type ApplicationConfigResourceImpl struct {
	client       restapi.RestClient
	resourcePath string
}

//GetOne retrieves a single application config from Instana API by its ID
func (resource *ApplicationConfigResourceImpl) GetOne(id string) (restapi.ApplicationConfig, error) {
	data, err := resource.client.GetOne(id, resource.resourcePath)
	if err != nil {
		return restapi.ApplicationConfig{}, err
	}
	return resource.validateResponseAndConvertToStruct(data)
}

func (resource *ApplicationConfigResourceImpl) validateResponseAndConvertToStruct(data []byte) (restapi.ApplicationConfig, error) {
	applicationConfig := restapi.ApplicationConfig{}
	if err := json.Unmarshal(data, &applicationConfig); err != nil {
		return applicationConfig, fmt.Errorf("failed to parse json; %s", err)
	}

	if err := applicationConfig.Validate(); err != nil {
		return applicationConfig, err
	}
	return applicationConfig, nil
}

func (resource *ApplicationConfigResourceImpl) validateAllApplicationConfigs(bindings []restapi.ApplicationConfig) error {
	for _, b := range bindings {
		err := b.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

//Upsert creates or updates a application config
func (resource *ApplicationConfigResourceImpl) Upsert(binding restapi.ApplicationConfig) (restapi.ApplicationConfig, error) {
	if err := binding.Validate(); err != nil {
		return binding, err
	}
	data, err := resource.client.Put(binding, resource.resourcePath)
	if err != nil {
		return binding, err
	}
	return resource.validateResponseAndConvertToStruct(data)
}

//Delete deletes a application config
func (resource *ApplicationConfigResourceImpl) Delete(binding restapi.ApplicationConfig) error {
	return resource.DeleteByID(binding.ID)
}

//DeleteByID deletes a application config by its ID
func (resource *ApplicationConfigResourceImpl) DeleteByID(applicationConfigID string) error {
	return resource.client.Delete(applicationConfigID, resource.resourcePath)
}
