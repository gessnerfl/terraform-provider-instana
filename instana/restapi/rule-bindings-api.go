package restapi

import "errors"

//RuleBindingResource represents the REST resource of  rule bindings at Instana
type RuleBindingResource interface {
	GetOne(id string) (RuleBinding, error)
	Upsert(rule RuleBinding) (RuleBinding, error)
	Delete(rule RuleBinding) error
	DeleteByID(ruleID string) error
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

//GetID implemention of the interface InstanaDataObject
func (binding RuleBinding) GetID() string {
	return binding.ID
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
