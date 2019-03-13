package resources

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

//NewRuleBindingResource constructs a new instance of RuleBindingResource
func NewRuleBindingResource(client restapi.RestClient) *RuleBindingResource {
	return &RuleBindingResource{
		client:       client,
		resourcePath: "/ruleBindings",
	}
}

//RuleBindingResource is the GO representation of the rule binding API of the Instana
type RuleBindingResource struct {
	client       restapi.RestClient
	resourcePath string
}

//RuleBinding is the representation of a rule binding in Instana
type RuleBinding struct {
	ID             string   `json:"id"`
	Enabled        bool     `json:"enabled"`
	Triggering     bool     `json:"triggering"`
	Severity       int      `json:"severity"`
	Text           string   `json:"text"`
	Description    string   `json:"description"`
	ExpirationTime int      `json:"expirationTime"`
	Query          string   `json:"query"`
	RuleIds        []string `json:"ruleIds"`
}

//Validate implementation of the interface InstanaDataObject to verify if data object is correct
func (binding RuleBinding) Validate() error {
	if len(binding.ID) == 0 {
		return errors.New("ID is missing")
	}
	if len(binding.Text) == 0 {
		return errors.New("Text is missing")
	}
	if len(binding.Description) == 0 {
		return errors.New("Description is missing")
	}
	if len(binding.RuleIds) == 0 {
		return errors.New("RuleIds are missing")
	}
	return nil
}

//GetID implemention of the interface InstanaDataObject
func (binding RuleBinding) GetID() string {
	return binding.ID
}

//Delete deletes a rule binding
func (resource *RuleBindingResource) Delete(binding RuleBinding) error {
	return resource.DeleteByID(binding.ID)
}

//DeleteByID deletes a rule binding by its ID
func (resource *RuleBindingResource) DeleteByID(ruleBindingID string) error {
	return resource.client.Delete(ruleBindingID, resource.resourcePath)
}

//Upsert creates or updates a rule binding
func (resource *RuleBindingResource) Upsert(binding RuleBinding) (RuleBinding, error) {
	if err := binding.Validate(); err != nil {
		return binding, err
	}
	data, err := resource.client.Put(binding, resource.resourcePath)
	if err != nil {
		return binding, err
	}
	return resource.validateResponseAndConvertToStruct(data)
}

//GetOne retrieves a single rule binding from Instana API by its ID
func (resource *RuleBindingResource) GetOne(id string) (RuleBinding, error) {
	data, err := resource.client.GetOne(id, resource.resourcePath)
	if err != nil {
		return RuleBinding{}, err
	}
	return resource.validateResponseAndConvertToStruct(data)
}

func (resource *RuleBindingResource) validateResponseAndConvertToStruct(data []byte) (RuleBinding, error) {
	ruleBinding := RuleBinding{}
	if err := json.Unmarshal(data, &ruleBinding); err != nil {
		return ruleBinding, fmt.Errorf("failed to parse json; %s", err)
	}

	if err := ruleBinding.Validate(); err != nil {
		return ruleBinding, err
	}
	return ruleBinding, nil
}

//GetAll returns all configured rule binding from Instana API
func (resource *RuleBindingResource) GetAll() ([]RuleBinding, error) {
	bindings := make([]RuleBinding, 0)

	data, err := resource.client.GetAll(resource.resourcePath)
	if err != nil {
		return bindings, err
	}

	if err := json.Unmarshal(data, &bindings); err != nil {
		return bindings, fmt.Errorf("failed to parse json; %s", err)
	}

	if err := validateAllRuleBindings(bindings); err != nil {
		return make([]RuleBinding, 0), fmt.Errorf("Received invalid JSON message, %s", err)
	}

	return bindings, nil
}

func validateAllRuleBindings(bindings []RuleBinding) error {
	for _, b := range bindings {
		err := b.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}
