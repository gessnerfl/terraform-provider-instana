package restapi

import "errors"

//RuleResource represents the REST resource of custom rules at Instana
type RuleResource interface {
	GetOne(id string) (Rule, error)
	GetAll() ([]Rule, error)
	Upsert(rule Rule) (Rule, error)
	Delete(rule Rule) error
	DeleteByID(ruleID string) error
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
	ConditionValue    float64 `json:"conditionValue"`
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
