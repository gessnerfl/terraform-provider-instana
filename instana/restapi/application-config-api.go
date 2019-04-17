package restapi

import "errors"

//MatchExpressionType type for MatchExpression discriminator type
type MatchExpressionType string

const (
	//BinaryOperatorExpressionType discriminator type for binary operations
	BinaryOperatorExpressionType MatchExpressionType = "BINARY_OP"
	//LeafExpressionType discriminator type for leaf operations
	LeafExpressionType MatchExpressionType = "LEAF"
)

//ApplicationConfigResource represents the REST resource of application perspective configuration at Instana
type ApplicationConfigResource interface {
	GetOne(id string) (ApplicationConfig, error)
	Upsert(rule ApplicationConfig) (ApplicationConfig, error)
	Delete(rule ApplicationConfig) error
	DeleteByID(applicationID string) error
}

//MatchExpression is the interface definition of a match expression in Instana
type MatchExpression interface {
	GetType() MatchExpressionType
	Validate() error
}

//NewBinaryOperator creates and new binary operator MatchExpression
func NewBinaryOperator(left MatchExpression, conjunction string, right MatchExpression) MatchExpression {
	return BinaryOperator{
		Dtype:       BinaryOperatorExpressionType,
		Left:        left,
		Right:       right,
		Conjunction: conjunction,
	}
}

//BinaryOperator is the representation of a binary operator expression in Instana
type BinaryOperator struct {
	Dtype       MatchExpressionType `json:"type"`
	Left        interface{}         `json:"left"`
	Right       interface{}         `json:"right"`
	Conjunction string              `json:"conjunction"`
}

//NewTagMatcherExpression creates and new tag matcher MatchExpression
func NewTagMatcherExpression(key string, operator string, value string) MatchExpression {
	return TagMatcherExpression{
		Dtype:    LeafExpressionType,
		Key:      key,
		Operator: operator,
		Value:    &value,
	}
}

//TagMatcherExpression is the representation of a tag matcher expression in Instana
type TagMatcherExpression struct {
	Dtype    MatchExpressionType `json:"type"`
	Key      string              `json:"key"`
	Operator string              `json:"operator"`
	Value    *string             `json:"value"`
}

//ApplicationConfig is the representation of a application perspective configuration in Instana
type ApplicationConfig struct {
	ID                 string      `json:"id"`
	Label              string      `json:"label"`
	MatchSpecification interface{} `json:"matchSpecification"`
	Scope              string      `json:"scope"`
}

//GetID implemention of the interface InstanaDataObject
func (a ApplicationConfig) GetID() string {
	return a.ID
}

//Validate implemention of the interface InstanaDataObject for ApplicationConfig
func (a ApplicationConfig) Validate() error {
	if len(a.ID) == 0 {
		return errors.New("ID is missing")
	}
	if len(a.Label) == 0 {
		return errors.New("Label is missing")
	}
	if a.MatchSpecification == nil {
		return errors.New("MatchSpecification is missing")
	}

	if err := a.MatchSpecification.(MatchExpression).Validate(); err != nil {
		return err
	}

	if len(a.Scope) == 0 {
		return errors.New("Scope is missing")
	}
	return nil
}

//GetType implemention of the interface MatchExpression for BinaryOperator
func (b BinaryOperator) GetType() MatchExpressionType {
	return b.Dtype
}

//Validate implemention of the interface MatchExpression for BinaryOperator
func (b BinaryOperator) Validate() error {
	if b.Left == nil {
		return errors.New("Left expression is missing")
	}
	if err := b.Left.(MatchExpression).Validate(); err != nil {
		return err
	}

	if len(b.Conjunction) == 0 {
		return errors.New("Conjunction of expressions is missing")
	}

	if b.Right == nil {
		return errors.New("Right expression is missing")
	}
	if err := b.Right.(MatchExpression).Validate(); err != nil {
		return err
	}
	return nil
}

//GetType implemention of the interface MatchExpression for TagMatcherExpression
func (t TagMatcherExpression) GetType() MatchExpressionType {
	return t.Dtype
}

//Validate implemention of the interface MatchExpression for TagMatcherExpression
func (t TagMatcherExpression) Validate() error {
	if len(t.Key) == 0 {
		return errors.New("Key of tag expression is missing")
	}
	if len(t.Operator) == 0 {
		return errors.New("Operator of tag expression is missing")
	}
	return nil
}
