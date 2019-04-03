package resources

import (
	"encoding/json"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

//NewRuleResource constructs a new instance of RuleResource
func NewRuleResource(client restapi.RestClient) restapi.RuleResource {
	return &RuleResourceImpl{
		client:       client,
		resourcePath: restapi.RulesResourcePath,
	}
}

//RuleResourceImpl is the GO representation of the Rule Resource of Instana
type RuleResourceImpl struct {
	client       restapi.RestClient
	resourcePath string
}

//GetOne retrieves a single custom rule from Instana API by its ID
func (resource *RuleResourceImpl) GetOne(id string) (restapi.Rule, error) {
	data, err := resource.client.GetOne(id, resource.resourcePath)
	if err != nil {
		return restapi.Rule{}, err
	}
	return resource.validateResponseAndConvertToStruct(data)
}

func (resource *RuleResourceImpl) validateResponseAndConvertToStruct(data []byte) (restapi.Rule, error) {
	rule := restapi.Rule{}
	if err := json.Unmarshal(data, &rule); err != nil {
		return rule, fmt.Errorf("failed to parse json; %s", err)
	}

	if err := rule.Validate(); err != nil {
		return rule, err
	}
	return rule, nil
}

func (resource *RuleResourceImpl) validateAllRules(rules []restapi.Rule) error {
	for _, r := range rules {
		err := r.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

//Upsert creates or updates a custom rule
func (resource *RuleResourceImpl) Upsert(rule restapi.Rule) (restapi.Rule, error) {
	if err := rule.Validate(); err != nil {
		return rule, err
	}
	data, err := resource.client.Put(rule, resource.resourcePath)
	if err != nil {
		return rule, err
	}
	return resource.validateResponseAndConvertToStruct(data)
}

//Delete deletes a custom rule
func (resource *RuleResourceImpl) Delete(rule restapi.Rule) error {
	return resource.DeleteByID(rule.ID)
}

//DeleteByID deletes a custom rule by its ID
func (resource *RuleResourceImpl) DeleteByID(ruleID string) error {
	return resource.client.Delete(ruleID, resource.resourcePath)
}
