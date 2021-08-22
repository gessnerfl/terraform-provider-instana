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

//SupportedConjunctionTypes list of supported binary expression operators of Instana API
var SupportedConjunctionTypes = []LogicalOperatorType{LogicalAnd, LogicalOr}

//ApplicationConfigScope type definition of the application config scope of the Instana Web REST API
type ApplicationConfigScope string

//ApplicationConfigScopes type definition of slice of ApplicationConfigScope
type ApplicationConfigScopes []ApplicationConfigScope

//ToStringSlice returns a slice containing the string representations of the given scopes
func (scopes ApplicationConfigScopes) ToStringSlice() []string {
	result := make([]string, len(scopes))
	for i, s := range scopes {
		result[i] = string(s)
	}
	return result
}

//IsSupported checks if the given ApplicationConfigScope is defined as a supported ApplicationConfigScope of the underlying slice
func (scopes ApplicationConfigScopes) IsSupported(s ApplicationConfigScope) bool {
	for _, scope := range scopes {
		if s == scope {
			return true
		}
	}
	return false
}

const (
	//ApplicationConfigScopeIncludeNoDownstream constant for the scope INCLUDE_NO_DOWNSTREAM
	ApplicationConfigScopeIncludeNoDownstream = ApplicationConfigScope("INCLUDE_NO_DOWNSTREAM")
	//ApplicationConfigScopeIncludeImmediateDownstreamDatabaseAndMessaging constant for the scope INCLUDE_IMMEDIATE_DOWNSTREAM_DATABASE_AND_MESSAGING
	ApplicationConfigScopeIncludeImmediateDownstreamDatabaseAndMessaging = ApplicationConfigScope("INCLUDE_IMMEDIATE_DOWNSTREAM_DATABASE_AND_MESSAGING")
	//ApplicationConfigScopeIncludeAllDownstream constant for the scope INCLUDE_ALL_DOWNSTREAM
	ApplicationConfigScopeIncludeAllDownstream = ApplicationConfigScope("INCLUDE_ALL_DOWNSTREAM")
)

//SupportedApplicationConfigScopes supported ApplicationConfigScopes of the Instana Web REST API
var SupportedApplicationConfigScopes = ApplicationConfigScopes{
	ApplicationConfigScopeIncludeNoDownstream,
	ApplicationConfigScopeIncludeImmediateDownstreamDatabaseAndMessaging,
	ApplicationConfigScopeIncludeAllDownstream,
}

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
func NewBinaryOperator(left MatchExpression, conjunction LogicalOperatorType, right MatchExpression) MatchExpression {
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
	Conjunction LogicalOperatorType `json:"conjunction"`
}

//NewComparisionExpression creates and new tag matcher expression for a comparision
func NewComparisionExpression(key string, entity MatcherExpressionEntity, operator TagFilterOperator, value string) MatchExpression {
	return TagMatcherExpression{
		Dtype:    LeafExpressionType,
		Key:      key,
		Entity:   entity,
		Operator: operator,
		Value:    &value,
	}
}

//NewUnaryOperationExpression creates and new tag matcher expression for a unary operation
func NewUnaryOperationExpression(key string, entity MatcherExpressionEntity, operator TagFilterOperator) MatchExpression {
	return TagMatcherExpression{
		Dtype:    LeafExpressionType,
		Key:      key,
		Entity:   entity,
		Operator: operator,
	}
}

//MatcherExpressionEntity type representing the matcher expression entity of a Matcher Expression (either source or destination or not applicable)
type MatcherExpressionEntity string

//MatcherExpressionEntities custom type representing a slice of MatcherExpressionEntity
type MatcherExpressionEntities []MatcherExpressionEntity

//ToStringSlice Returns the string representations fo the aggregations
func (entities MatcherExpressionEntities) ToStringSlice() []string {
	result := make([]string, len(entities))
	for i, v := range entities {
		result[i] = string(v)
	}
	return result
}

const (
	//MatcherExpressionEntitySource const for a SOURCE matcher expression entity
	MatcherExpressionEntitySource = MatcherExpressionEntity("SOURCE")
	//MatcherExpressionEntityDestination const for a DESTINATION matcher expression entity
	MatcherExpressionEntityDestination = MatcherExpressionEntity("DESTINATION")
	//MatcherExpressionEntityNotApplicable const for a NOT_APPLICABLE matcher expression entity
	MatcherExpressionEntityNotApplicable = MatcherExpressionEntity("NOT_APPLICABLE")
)

//SupportedMatcherExpressionEntities slice of supported matcher expression entity types
var SupportedMatcherExpressionEntities = MatcherExpressionEntities{MatcherExpressionEntitySource, MatcherExpressionEntityDestination, MatcherExpressionEntityNotApplicable}

//IsSupported check if the provided matcher expression entity is supported
func (entities MatcherExpressionEntities) IsSupported(entity MatcherExpressionEntity) bool {
	for _, v := range entities {
		if v == entity {
			return true
		}
	}
	return false
}

//TagMatcherExpression is the representation of a tag matcher expression in Instana
type TagMatcherExpression struct {
	Dtype    MatchExpressionType     `json:"type"`
	Key      string                  `json:"key"`
	Entity   MatcherExpressionEntity `json:"entity"`
	Operator TagFilterOperator       `json:"operator"`
	Value    *string                 `json:"value"`
}

//ApplicationConfig is the representation of a application perspective configuration in Instana
type ApplicationConfig struct {
	ID                 string                 `json:"id"`
	Label              string                 `json:"label"`
	MatchSpecification interface{}            `json:"matchSpecification"`
	Scope              ApplicationConfigScope `json:"scope"`
	BoundaryScope      BoundaryScope          `json:"boundaryScope"`
}

//GetIDForResourcePath implemention of the interface InstanaDataObject
func (a *ApplicationConfig) GetIDForResourcePath() string {
	return a.ID
}

//Validate implementation of the interface InstanaDataObject for ApplicationConfig
func (a *ApplicationConfig) Validate() error {
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

	if utils.IsBlank(string(a.Scope)) {
		return errors.New("scope is missing")
	}
	if !SupportedApplicationConfigScopes.IsSupported(a.Scope) {
		return errors.New("scope is not supported")
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
	if !SupportedMatcherExpressionEntities.IsSupported(t.Entity) {
		return fmt.Errorf("entity %s of tag expression is not supported", t.Entity)
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
func IsSupportedComparision(operator TagFilterOperator) bool {
	return isInMatcherOperatorSlice(SupportedComparisonOperators, operator)
}

//IsSupportedUnaryOperatorExpression returns true if the provided operator is a valid unary operator type
func IsSupportedUnaryOperatorExpression(operator TagFilterOperator) bool {
	return isInMatcherOperatorSlice(SupportedUnaryExpressionOperators, operator)
}

func isInMatcherOperatorSlice(allOperators []TagFilterOperator, operator TagFilterOperator) bool {
	for _, v := range allOperators {
		if v == operator {
			return true
		}
	}
	return false
}

//IsSupportedConjunctionType returns true if the provided operator is a valid conjunction type
func IsSupportedConjunctionType(operator LogicalOperatorType) bool {
	return isInConjunctionTypeSlice(SupportedConjunctionTypes, operator)
}

func isInConjunctionTypeSlice(allOperators []LogicalOperatorType, operator LogicalOperatorType) bool {
	for _, v := range allOperators {
		if v == operator {
			return true
		}
	}
	return false
}
