package resources

import (
	"encoding/json"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

//NewCustomEventSpecificationResource constructs a new instance of CustomEventSpecificationResource
func NewCustomEventSpecificationResource(client restapi.RestClient) restapi.CustomEventSpecificationResource {
	return &customEventSpecificationResourceImpl{
		client:       client,
		resourcePath: restapi.CustomEventSpecificationResourcePath,
	}
}

//customEventSpecificationResourceImpl is the GO representation of the application config API of the Instana
type customEventSpecificationResourceImpl struct {
	client       restapi.RestClient
	resourcePath string
}

//GetOne retrieves a single application config from Instana API by its ID
func (resource *customEventSpecificationResourceImpl) GetOne(id string) (restapi.CustomEventSpecification, error) {
	data, err := resource.client.GetOne(id, resource.resourcePath)
	if err != nil {
		return restapi.CustomEventSpecification{}, err
	}
	return resource.validateResponseAndConvertToStruct(data)
}

func (resource *customEventSpecificationResourceImpl) validateResponseAndConvertToStruct(data []byte) (restapi.CustomEventSpecification, error) {
	spec := restapi.CustomEventSpecification{}
	if err := json.Unmarshal(data, &spec); err != nil {
		return spec, fmt.Errorf("failed to parse json; %s", err)
	}

	if err := spec.Validate(); err != nil {
		return spec, err
	}
	return spec, nil
}

//Upsert creates or updates a application config
func (resource *customEventSpecificationResourceImpl) Upsert(specification restapi.CustomEventSpecification) (restapi.CustomEventSpecification, error) {
	if err := specification.Validate(); err != nil {
		return specification, err
	}
	data, err := resource.client.Put(specification, resource.resourcePath)
	if err != nil {
		return specification, err
	}
	return resource.validateResponseAndConvertToStruct(data)
}

//Delete deletes a application config
func (resource *customEventSpecificationResourceImpl) Delete(specification restapi.CustomEventSpecification) error {
	return resource.DeleteByID(specification.ID)
}

//DeleteByID deletes a application config by its ID
func (resource *customEventSpecificationResourceImpl) DeleteByID(specificationID string) error {
	return resource.client.Delete(specificationID, resource.resourcePath)
}
