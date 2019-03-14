package resources

import (
	"encoding/json"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

//NewRuleBindingResource constructs a new instance of RuleBindingResource
func NewRuleBindingResource(client restapi.RestClient) restapi.RuleBindingResource {
	return &RuleBindingResourceImpl{
		client:       client,
		resourcePath: "/ruleBindings",
	}
}

//RuleBindingResourceImpl is the GO representation of the rule binding API of the Instana
type RuleBindingResourceImpl struct {
	client       restapi.RestClient
	resourcePath string
}

//GetOne retrieves a single rule binding from Instana API by its ID
func (resource *RuleBindingResourceImpl) GetOne(id string) (restapi.RuleBinding, error) {
	data, err := resource.client.GetOne(id, resource.resourcePath)
	if err != nil {
		return restapi.RuleBinding{}, err
	}
	return resource.validateResponseAndConvertToStruct(data)
}

func (resource *RuleBindingResourceImpl) validateResponseAndConvertToStruct(data []byte) (restapi.RuleBinding, error) {
	ruleBinding := restapi.RuleBinding{}
	if err := json.Unmarshal(data, &ruleBinding); err != nil {
		return ruleBinding, fmt.Errorf("failed to parse json; %s", err)
	}

	if err := ruleBinding.Validate(); err != nil {
		return ruleBinding, err
	}
	return ruleBinding, nil
}

//GetAll returns all configured rule binding from Instana API
func (resource *RuleBindingResourceImpl) GetAll() ([]restapi.RuleBinding, error) {
	bindings := make([]restapi.RuleBinding, 0)

	data, err := resource.client.GetAll(resource.resourcePath)
	if err != nil {
		return bindings, err
	}

	if err := json.Unmarshal(data, &bindings); err != nil {
		return bindings, fmt.Errorf("failed to parse json; %s", err)
	}

	if err := resource.validateAllRuleBindings(bindings); err != nil {
		return make([]restapi.RuleBinding, 0), fmt.Errorf("Received invalid JSON message, %s", err)
	}

	return bindings, nil
}

func (resource *RuleBindingResourceImpl) validateAllRuleBindings(bindings []restapi.RuleBinding) error {
	for _, b := range bindings {
		err := b.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

//Upsert creates or updates a rule binding
func (resource *RuleBindingResourceImpl) Upsert(binding restapi.RuleBinding) (restapi.RuleBinding, error) {
	if err := binding.Validate(); err != nil {
		return binding, err
	}
	data, err := resource.client.Put(binding, resource.resourcePath)
	if err != nil {
		return binding, err
	}
	return resource.validateResponseAndConvertToStruct(data)
}

//Delete deletes a rule binding
func (resource *RuleBindingResourceImpl) Delete(binding restapi.RuleBinding) error {
	return resource.DeleteByID(binding.ID)
}

//DeleteByID deletes a rule binding by its ID
func (resource *RuleBindingResourceImpl) DeleteByID(ruleBindingID string) error {
	return resource.client.Delete(ruleBindingID, resource.resourcePath)
}
