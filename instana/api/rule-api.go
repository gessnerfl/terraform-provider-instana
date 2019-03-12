package api

import (
	"encoding/json"
	"errors"
	"fmt"
)

//NewRuleAPI constructs a new instance of RuleApi
func NewRuleAPI(client RestClient) *RuleAPI {
	return &RuleAPI{
		client:       client,
		resourcePath: "/rules",
	}
}

//RuleAPI is the GO representation of the Rule API of Instana
type RuleAPI struct {
	client       RestClient
	resourcePath string
}

//Rule is the representation of a custom rule in Instana
type Rule struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	EntityType        string  `json:"entityType"`
	MetricName        string  `json:"metricName"`
	Rollup            int     `json:"rollup"`
	Window            int     `json:"window"`
	Aggregation       string  `json:"aggregation"`
	ConditionOperator string  `json:"conditionOperator"`
	ConditionValue    float32 `json:"conditionValue"`
}

//GetID implemention of the interface InstanaDataObject
func (rule Rule) GetID() string {
	return rule.ID
}

//Validate implementation of the interface InstanaDataObject to verify if data object is correct
func (rule Rule) Validate() error {
	if len(rule.ID) == 0 {
		return errors.New("ID is missing")
	}
	if len(rule.Name) == 0 {
		return errors.New("Name is missing")
	}
	if len(rule.EntityType) == 0 {
		return errors.New("EntityType is missing")
	}
	if len(rule.MetricName) == 0 {
		return errors.New("MetricName is missing")
	}
	if len(rule.Aggregation) == 0 {
		return errors.New("Aggregation is missing")
	}
	if len(rule.ConditionOperator) == 0 {
		return errors.New("ConditionOperator is missing")
	}
	return nil
}

//Delete deletes a custom rule
func (resource *RuleAPI) Delete(rule Rule) error {
	return resource.DeleteByID(rule.ID)
}

//DeleteByID deletes a custom rule by its ID
func (resource *RuleAPI) DeleteByID(ruleID string) error {
	return resource.client.Delete(ruleID, resource.resourcePath)
}

//Upsert creates or updates a custom rule
func (resource *RuleAPI) Upsert(rule Rule) error {
	return resource.client.Put(rule, resource.resourcePath)
}

//GetAll returns all configured custom rules from Instana API
func (resource *RuleAPI) GetAll() ([]Rule, error) {
	rules := make([]Rule, 0)

	data, err := resource.client.GetAll(resource.resourcePath)
	if err != nil {
		return rules, err
	}

	if err := json.Unmarshal(data, &rules); err != nil {
		return rules, fmt.Errorf("failed to parse json; %s", err)
	}

	if err := validateAllRules(rules); err != nil {
		return make([]Rule, 0), fmt.Errorf("Received invalid JSON message, %s", err)
	}

	return rules, nil
}

func validateAllRules(rules []Rule) error {
	for _, r := range rules {
		err := r.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}
