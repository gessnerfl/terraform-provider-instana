package restapi

import "errors"

//CustomEventSpecificationResource represents the REST resource of custom event specification at Instana
type CustomEventSpecificationResource interface {
	GetOne(id string) (CustomEventSpecification, error)
	Upsert(spec CustomEventSpecification) (CustomEventSpecification, error)
	Delete(spec CustomEventSpecification) error
	DeleteByID(specID string) error
}

//RuleType custom type representing the type of the custom event specification rule
type RuleType string

const (
	//SystemRuleType const for RuleType of System Rules
	SystemRuleType = "system"
)

//NewSystemRuleSpecification creates a new instance of a System Rule
func NewSystemRuleSpecification(systemRuleID string, severity int) RuleSpecification {
	return RuleSpecification{
		DType:        SystemRuleType,
		SystemRuleID: systemRuleID,
		Severity:     severity,
	}
}

//RuleSpecification representation of a rule specification for a CustomEventSpecification
type RuleSpecification struct {
	DType        RuleType `json:"ruleType"`
	SystemRuleID string   `json:"systemRuleId"`
	Severity     int      `json:"severity"`
}

//Validate Rule interface implementation for SystemRule
func (r *RuleSpecification) Validate() error {
	if len(r.DType) == 0 {
		return errors.New("type of system rule is missing")
	}
	if r.DType == SystemRuleType {
		return r.validateSystemRule()
	}
	return nil
}

func (r *RuleSpecification) validateSystemRule() error {
	if len(r.SystemRuleID) == 0 {
		return errors.New("id of system rule is missing")
	}
	return nil
}

//EventSpecificationDownstream definition of downstream reporting for the event specification
type EventSpecificationDownstream struct {
	IntegrationIds                []string `json:"integrationIds"`
	BroadcastToAllAlertingConfigs bool     `json:"broadcastToAllAlertingConfigs"`
}

//Validate validates the consitency of an EventSpecificationDownstream
func (d EventSpecificationDownstream) Validate() error {
	if len(d.IntegrationIds) == 0 {
		return errors.New("At least one integration id must be defined for a downstream specification")
	}
	return nil
}

//CustomEventSpecification is the representation of a custom event specification in Instana
type CustomEventSpecification struct {
	ID             string                        `json:"id"`
	Name           string                        `json:"name"`
	EntityType     string                        `json:"entityType"`
	Query          *string                       `json:"query"`
	Triggering     bool                          `json:"triggering"`
	Description    *string                       `json:"string"`
	ExpirationTime *int                          `json:"expirationTime"`
	Enabled        bool                          `json:"enabled"`
	Rules          []RuleSpecification           `json:"rules"`
	Downstream     *EventSpecificationDownstream `json:"downstream"`
}

//GetID implemention of the interface InstanaDataObject
func (spec CustomEventSpecification) GetID() string {
	return spec.ID
}

//Validate implementation of the interface InstanaDataObject to verify if data object is correct
func (spec CustomEventSpecification) Validate() error {
	if len(spec.ID) == 0 {
		return errors.New("ID is missing")
	}
	if len(spec.Name) == 0 {
		return errors.New("name is missing")
	}
	if len(spec.EntityType) == 0 {
		return errors.New("entity type is missing")
	}
	if len(spec.Rules) != 1 {
		return errors.New("exactly one rule must be defined")
	}
	for _, r := range spec.Rules {
		if err := r.Validate(); err != nil {
			return err
		}
	}

	if spec.Downstream != nil {
		if err := spec.Downstream.Validate(); err != nil {
			return err
		}
	}
	return nil
}
