package restapi

import (
	"errors"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/utils"
)

const (
	//ApplicationMonitoringBasePath path to application monitoring resource of Instana RESTful API
	ApplicationMonitoringBasePath = InstanaAPIBasePath + "/application-monitoring"
	//ApplicationMonitoringSettingsBasePath path to application monitoring settings resource of Instana RESTful API
	ApplicationMonitoringSettingsBasePath = ApplicationMonitoringBasePath + settingsPathElement
	//ApplicationConfigsResourcePath path to application config resource of Instana RESTful API
	ApplicationConfigsResourcePath = ApplicationMonitoringSettingsBasePath + "/application"
)

//MatchExpressionType type for MatchExpression discriminator type
type MatchExpressionType string

const (
	//BinaryOperatorExpressionType discriminator type for binary operations
	BinaryOperatorExpressionType MatchExpressionType = "BINARY_OP"
	//LeafExpressionType discriminator type for leaf operations
	LeafExpressionType MatchExpressionType = "LEAF"
)

//ConjunctionType custom type for conjunctions
type ConjunctionType string

const (
	//LogicalAnd constant for logical AND conjunction
	LogicalAnd = ConjunctionType("AND")
	//LogicalOr constant for logical OR conjunction
	LogicalOr = ConjunctionType("OR")
)

//MatcherOperator custom type for tag matcher operators
type MatcherOperator string

const (
	//EqualsOperator constant for the EQUALS operator
	EqualsOperator = MatcherOperator("EQUALS")
	//NotEqualOperator constant for the NOT_EQUAL operator
	NotEqualOperator = MatcherOperator("NOT_EQUAL")
	//ContainsOperator constant for the CONTAINS operator
	ContainsOperator = MatcherOperator("CONTAINS")
	//NotContainOperator constant for the NOT_CONTAIN operator
	NotContainOperator = MatcherOperator("NOT_CONTAIN")
	//IsEmptyOperator constant for the IS_EMPTY operator
	IsEmptyOperator = MatcherOperator("IS_EMPTY")
	//NotEmptyOperator constant for the NOT_EMPTY operator
	NotEmptyOperator = MatcherOperator("NOT_EMPTY")
	//IsBlankOperator constant for the IS_BLANK operator
	IsBlankOperator = MatcherOperator("IS_BLANK")
	//NotBlankOperator constant for the NOT_BLANK operator
	NotBlankOperator = MatcherOperator("NOT_BLANK")

	//StartsWithOperator constant for the STARTS_WITH operator
	StartsWithOperator = MatcherOperator("STARTS_WITH")
	//EndsWithOperator constant for the ENDS_WITH operator
	EndsWithOperator = MatcherOperator("ENDS_WITH")
	//NotStartsWithOperator constant for the NOT_STARTS_WITH operator
	NotStartsWithOperator = MatcherOperator("NOT_STARTS_WITH")
	//NotEndsWithOperator constant for the NOT_ENDS_WITH operator
	NotEndsWithOperator = MatcherOperator("NOT_ENDS_WITH")
	//GreaterOrEqualThanOperator constant for the GREATER_OR_EQUAL_THAN operator
	GreaterOrEqualThanOperator = MatcherOperator("GREATER_OR_EQUAL_THAN")
	//LessOrEqualThanOperator constant for the LESS_OR_EQUAL_THAN operator
	LessOrEqualThanOperator = MatcherOperator("LESS_OR_EQUAL_THAN")
	//GreaterThanOperator constant for the GREATER_THAN operator
	GreaterThanOperator = MatcherOperator("GREATER_THAN")
	//LessThanOperator constant for the LESS_THAN operator
	LessThanOperator = MatcherOperator("LESS_THAN")
)

//SupportedComparisionOperators list of supported comparision operators of Instana API
var SupportedComparisionOperators = []MatcherOperator{
	EqualsOperator,
	NotEqualOperator,
	ContainsOperator,
	NotContainOperator,
	StartsWithOperator,
	EndsWithOperator,
	NotStartsWithOperator,
	NotEndsWithOperator,
	GreaterOrEqualThanOperator,
	LessOrEqualThanOperator,
	GreaterThanOperator,
	LessThanOperator,
}

//SupportedUnaryExpressionOperators list of supported unary expression operators of Instana API
var SupportedUnaryExpressionOperators = []MatcherOperator{
	IsEmptyOperator,
	NotEmptyOperator,
	IsBlankOperator,
	NotBlankOperator,
}

//SupportedConjunctionTypes list of supported binary expression operators of Instana API
var SupportedConjunctionTypes = []ConjunctionType{LogicalAnd, LogicalOr}

//BoundaryScope type definition of the application config boundary scope of the Instana Web REST API
type BoundaryScope string

//BoundaryScopes type definition of slice of BoundaryScopes
type BoundaryScopes []BoundaryScope

//ToStringSlice returns a slice containing the string representations of the given boundary scopes
func (scopes BoundaryScopes) ToStringSlice() []string {
	result := make([]string, len(scopes))
	for i, s := range scopes {
		result[i] = string(s)
	}
	return result
}

//IsSupported checks if the given BoundaryScope is defined as a supported BoundaryScope of the underlying slice
func (scopes BoundaryScopes) IsSupported(s BoundaryScope) bool {
	for _, scope := range scopes {
		if s == scope {
			return true
		}
	}
	return false
}

const (
	//BoundaryScopeAll constant value for the boundary scope ALL of an application config of the Instana Web REST API
	BoundaryScopeAll = BoundaryScope("ALL")
	//BoundaryScopeInbound constant value for the boundary scope INBOUND of an application config of the Instana Web REST API
	BoundaryScopeInbound = BoundaryScope("INBOUND")
	//BoundaryScopeDefault constant value for the boundary scope DEFAULT of an application config of the Instana Web REST API
	BoundaryScopeDefault = BoundaryScope("DEFAULT")
)

//SupportedBoundaryScopes supported BoundaryScopes of the Instana Web REST API
var SupportedBoundaryScopes = BoundaryScopes{BoundaryScopeAll, BoundaryScopeInbound, BoundaryScopeDefault}

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
func NewBinaryOperator(left MatchExpression, conjunction ConjunctionType, right MatchExpression) MatchExpression {
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
	Conjunction ConjunctionType     `json:"conjunction"`
}

//NewComparisionExpression creates and new tag matcher expression for a comparision
func NewComparisionExpression(key string, operator MatcherOperator, value string) MatchExpression {
	return TagMatcherExpression{
		Dtype:    LeafExpressionType,
		Key:      key,
		Operator: operator,
		Value:    &value,
	}
}

//NewUnaryOperationExpression creates and new tag matcher expression for a unary operation
func NewUnaryOperationExpression(key string, operator MatcherOperator) MatchExpression {
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
	Operator MatcherOperator     `json:"operator"`
	Value    *string             `json:"value"`
}

//ApplicationConfig is the representation of a application perspective configuration in Instana
type ApplicationConfig struct {
	ID                 string        `json:"id"`
	Label              string        `json:"label"`
	MatchSpecification interface{}   `json:"matchSpecification"`
	Scope              string        `json:"scope"`
	BoundaryScope      BoundaryScope `json:"boundaryScope"`
}

//GetID implemention of the interface InstanaDataObject
func (a ApplicationConfig) GetID() string {
	return a.ID
}

//Validate implemention of the interface InstanaDataObject for ApplicationConfig
func (a ApplicationConfig) Validate() error {
	if utils.IsBlank(a.ID) {
		return errors.New("id is missing")
	}
	if utils.IsBlank(a.Label) {
		return errors.New("label is missing")
	}
	if a.MatchSpecification == nil {
		return errors.New("match specification is missing")
	}

	if err := a.MatchSpecification.(MatchExpression).Validate(); err != nil {
		return err
	}

	if utils.IsBlank(a.Scope) {
		return errors.New("scope is missing")
	}
	if utils.IsBlank(string(a.BoundaryScope)) {
		return errors.New("boundary scope is missing")
	}
	if !SupportedBoundaryScopes.IsSupported(a.BoundaryScope) {
		return errors.New("boundary scope is not supported")
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

	if !IsSupportedConjunctionType(b.Conjunction) {
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
func IsSupportedComparision(operator MatcherOperator) bool {
	return isInMatcherOperatorSlice(SupportedComparisionOperators, operator)
}

//IsSupportedUnaryOperatorExpression returns true if the provided operator is a valid unary operator type
func IsSupportedUnaryOperatorExpression(operator MatcherOperator) bool {
	return isInMatcherOperatorSlice(SupportedUnaryExpressionOperators, operator)
}

func isInMatcherOperatorSlice(allOperators []MatcherOperator, operator MatcherOperator) bool {
	for _, v := range allOperators {
		if v == operator {
			return true
		}
	}
	return false
}

//IsSupportedConjunctionType returns true if the provided operator is a valid conjunction type
func IsSupportedConjunctionType(operator ConjunctionType) bool {
	return isInConjunctionTypeSlice(SupportedConjunctionTypes, operator)
}

func isInConjunctionTypeSlice(allOperators []ConjunctionType, operator ConjunctionType) bool {
	for _, v := range allOperators {
		if v == operator {
			return true
		}
	}
	return false
}
