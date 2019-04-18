package restapi

import (
	"errors"
	"fmt"
)

//MatchExpressionType type for MatchExpression discriminator type
type MatchExpressionType string

const (
	//BinaryOperatorExpressionType discriminator type for binary operations
	BinaryOperatorExpressionType MatchExpressionType = "BINARY_OP"
	//LeafExpressionType discriminator type for leaf operations
	LeafExpressionType MatchExpressionType = "LEAF"
)

//SupportedComparisionOperators list of supported comparision operators of Instana API
var SupportedComparisionOperators = []string{"EQUALS", "NOT_EQUAL", "CONTAINS", "NOT_CONTAIN"}

//SupportedUnaryOperatorExpressionOperators list of supported unary expression operators of Instana API
var SupportedUnaryOperatorExpressionOperators = []string{"IS_EMPTY", "NOT_EMPTY", "IS_BLANK", "NOT_BLANK"}

//SupportedConjunctions list of supported binary expression operators of Instana API
var SupportedConjunctions = []string{"AND", "OR"}

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

//NewComparisionExpression creates and new tag matcher expression for a comparision
func NewComparisionExpression(key string, operator string, value string) MatchExpression {
	return TagMatcherExpression{
		Dtype:    LeafExpressionType,
		Key:      key,
		Operator: operator,
		Value:    &value,
	}
}

//NewUnaryOperationExpression creates and new tag matcher expression for a unary operation
func NewUnaryOperationExpression(key string, operator string) MatchExpression {
	return TagMatcherExpression{
		Dtype:    LeafExpressionType,
		Key:      key,
		Operator: operator,
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
		return errors.New("id is missing")
	}
	if len(a.Label) == 0 {
		return errors.New("label is missing")
	}
	if a.MatchSpecification == nil {
		return errors.New("match specification is missing")
	}

	if err := a.MatchSpecification.(MatchExpression).Validate(); err != nil {
		return err
	}

	if len(a.Scope) == 0 {
		return errors.New("scope is missing")
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
		return errors.New("left expression is missing")
	}
	if err := b.Left.(MatchExpression).Validate(); err != nil {
		return err
	}

	if len(b.Conjunction) == 0 {
		return errors.New("conjunction of expressions is missing")
	}

	if !IsSupportedConjunction(b.Conjunction) {
		return fmt.Errorf("conjunction of type '%s' is not supported", b.Conjunction)
	}

	if b.Right == nil {
		return errors.New("right expression is missing")
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
		return errors.New("key of tag expression is missing")
	}
	if len(t.Operator) == 0 {
		return errors.New("operator of tag expression is missing")
	}

	if IsSupportedComparision(t.Operator) {
		if t.Value == nil || len(*t.Value) == 0 {
			return errors.New("value missing for comparision expression")
		}
	} else if IsSupportedUnaryOperatorExpression(t.Operator) {
		if t.Value != nil {
			return errors.New("value not allowed for unary operator expression")
		}
	} else {
		return fmt.Errorf("operator of tag expression is not supported")
	}

	return nil
}

//IsSupportedComparision returns true if the provided operator is a valid comparision type
func IsSupportedComparision(operator string) bool {
	return isInSlice(SupportedComparisionOperators, operator)
}

//IsSupportedUnaryOperatorExpression returns true if the provided operator is a valid unary operator type
func IsSupportedUnaryOperatorExpression(operator string) bool {
	return isInSlice(SupportedUnaryOperatorExpressionOperators, operator)
}

//IsSupportedConjunction returns true if the provided operator is a valid conjunction type
func IsSupportedConjunction(operator string) bool {
	return isInSlice(SupportedConjunctions, operator)
}

func isInSlice(allOperators []string, operator string) bool {
	for _, v := range allOperators {
		if v == operator {
			return true
		}
	}
	return false
}
