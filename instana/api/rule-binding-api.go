package api

import (
	"encoding/json"
	"errors"
	"fmt"
)

//NewRuleBindingAPI constructs a new instance of RuleBindingAPI
func NewRuleBindingAPI(client RestClient) *RuleBindingAPI {
	return &RuleBindingAPI{
		client:       client,
		resourcePath: "/ruleBindings",
	}
}

//RuleBindingAPI is the GO representation of the rule binding API of the Instana
type RuleBindingAPI struct {
	client       RestClient
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
	return nil
}

//GetID implemention of the interface InstanaDataObject
func (binding RuleBinding) GetID() string {
	return binding.ID
}

//Delete deletes a rule binding
func (resource *RuleBindingAPI) Delete(binding RuleBinding) error {
	return resource.DeleteByID(binding.ID)
}

//DeleteByID deletes a rule binding by its ID
func (resource *RuleBindingAPI) DeleteByID(ruleBindingID string) error {
	return resource.client.Delete(ruleBindingID, resource.resourcePath)
}

//Upsert creates or updates a rule binding
func (resource *RuleBindingAPI) Upsert(binding RuleBinding) error {
	return resource.client.Put(binding, resource.resourcePath)
}

//GetAll returns all configured rule binding from Instana API
func (resource *RuleBindingAPI) GetAll() ([]RuleBinding, error) {
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
