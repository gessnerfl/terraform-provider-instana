package resources

import (
	"encoding/json"
	"errors"
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
	applicationConfig, err := resource.unmarshalApplicationConfig(data)
	if err != nil {
		return applicationConfig, fmt.Errorf("failed to parse json; %s", err)
	}

	if err := applicationConfig.Validate(); err != nil {
		return applicationConfig, err
	}
	return applicationConfig, nil
}

func (resource *ApplicationConfigResourceImpl) unmarshalApplicationConfig(data []byte) (restapi.ApplicationConfig, error) {
	var matchExpression json.RawMessage
	temp := restapi.ApplicationConfig{
		MatchSpecification: &matchExpression,
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return restapi.ApplicationConfig{}, err
	}
	matchSpecification, err := resource.unmarshalMatchSpecification(matchExpression)
	if err != nil {
		return restapi.ApplicationConfig{}, err
	}
	return restapi.ApplicationConfig{
		ID:                 temp.ID,
		Label:              temp.Label,
		MatchSpecification: matchSpecification,
		Scope:              temp.Scope,
	}, nil
}

func (resource *ApplicationConfigResourceImpl) unmarshalMatchSpecification(raw json.RawMessage) (restapi.MatchExpression, error) {
	temp := struct {
		Dtype restapi.MatchExpressionType `json:"type"`
	}{}

	if err := json.Unmarshal(raw, &temp); err != nil {
		return nil, err
	}

	if temp.Dtype == restapi.BinaryOperatorExpressionType {
		return resource.unmarshalBinaryOperator(raw)
	} else if temp.Dtype == restapi.LeafExpressionType {
		return resource.unmarshalTagMatcherExpression(raw)
	} else {
		return nil, errors.New("invalid expression type")
	}
}

func (resource *ApplicationConfigResourceImpl) unmarshalBinaryOperator(raw json.RawMessage) (restapi.BinaryOperator, error) {
	var leftRaw json.RawMessage
	var rightRaw json.RawMessage
	temp := restapi.BinaryOperator{
		Left:  &leftRaw,
		Right: &rightRaw,
	}

	if err := json.Unmarshal(raw, &temp); err != nil {
		return restapi.BinaryOperator{}, err
	}

	left, err := resource.unmarshalMatchSpecification(leftRaw)
	if err != nil {
		return restapi.BinaryOperator{}, err
	}

	right, err := resource.unmarshalMatchSpecification(rightRaw)
	if err != nil {
		return restapi.BinaryOperator{}, err
	}
	return restapi.BinaryOperator{
		Dtype:       temp.Dtype,
		Left:        left,
		Right:       right,
		Conjunction: temp.Conjunction,
	}, nil
}

func (resource *ApplicationConfigResourceImpl) unmarshalTagMatcherExpression(raw json.RawMessage) (restapi.TagMatcherExpression, error) {
	data := restapi.TagMatcherExpression{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return restapi.TagMatcherExpression{}, err
	}
	return data, nil
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
