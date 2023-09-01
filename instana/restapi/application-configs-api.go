package restapi

const (
	//ApplicationMonitoringBasePath path to application monitoring resource of Instana RESTful API
	ApplicationMonitoringBasePath = InstanaAPIBasePath + "/application-monitoring"
	//ApplicationMonitoringSettingsBasePath path to application monitoring settings resource of Instana RESTful API
	ApplicationMonitoringSettingsBasePath = ApplicationMonitoringBasePath + settingsPathElement
	//ApplicationConfigsResourcePath path to application config resource of Instana RESTful API
	ApplicationConfigsResourcePath = ApplicationMonitoringSettingsBasePath + "/application"
)

// MatchExpressionType type for MatchExpression discriminator type
type MatchExpressionType string

const (
	//BinaryOperatorExpressionType discriminator type for binary operations
	BinaryOperatorExpressionType MatchExpressionType = "BINARY_OP"
	//LeafExpressionType discriminator type for leaf operations
	LeafExpressionType MatchExpressionType = "LEAF"
)

// ApplicationConfigScope type definition of the application config scope of the Instana Web REST API
type ApplicationConfigScope string

// ApplicationConfigScopes type definition of slice of ApplicationConfigScope
type ApplicationConfigScopes []ApplicationConfigScope

// ToStringSlice returns a slice containing the string representations of the given scopes
func (scopes ApplicationConfigScopes) ToStringSlice() []string {
	result := make([]string, len(scopes))
	for i, s := range scopes {
		result[i] = string(s)
	}
	return result
}

const (
	//ApplicationConfigScopeIncludeNoDownstream constant for the scope INCLUDE_NO_DOWNSTREAM
	ApplicationConfigScopeIncludeNoDownstream = ApplicationConfigScope("INCLUDE_NO_DOWNSTREAM")
	//ApplicationConfigScopeIncludeImmediateDownstreamDatabaseAndMessaging constant for the scope INCLUDE_IMMEDIATE_DOWNSTREAM_DATABASE_AND_MESSAGING
	ApplicationConfigScopeIncludeImmediateDownstreamDatabaseAndMessaging = ApplicationConfigScope("INCLUDE_IMMEDIATE_DOWNSTREAM_DATABASE_AND_MESSAGING")
	//ApplicationConfigScopeIncludeAllDownstream constant for the scope INCLUDE_ALL_DOWNSTREAM
	ApplicationConfigScopeIncludeAllDownstream = ApplicationConfigScope("INCLUDE_ALL_DOWNSTREAM")
)

// SupportedApplicationConfigScopes supported ApplicationConfigScopes of the Instana Web REST API
var SupportedApplicationConfigScopes = ApplicationConfigScopes{
	ApplicationConfigScopeIncludeNoDownstream,
	ApplicationConfigScopeIncludeImmediateDownstreamDatabaseAndMessaging,
	ApplicationConfigScopeIncludeAllDownstream,
}

// ApplicationConfigResource represents the REST resource of application perspective configuration at Instana
type ApplicationConfigResource interface {
	GetOne(id string) (ApplicationConfig, error)
	Upsert(rule ApplicationConfig) (ApplicationConfig, error)
	Delete(rule ApplicationConfig) error
	DeleteByID(applicationID string) error
}

// MatchExpression is the interface definition of a match expression in Instana
type MatchExpression interface {
	GetType() MatchExpressionType
}

// NewBinaryOperator creates and new binary operator MatchExpression
func NewBinaryOperator(left MatchExpression, conjunction LogicalOperatorType, right MatchExpression) MatchExpression {
	return &BinaryOperator{
		Dtype:       BinaryOperatorExpressionType,
		Left:        left,
		Right:       right,
		Conjunction: conjunction,
	}
}

// BinaryOperator is the representation of a binary operator expression in Instana
type BinaryOperator struct {
	Dtype       MatchExpressionType `json:"type"`
	Left        interface{}         `json:"left"`
	Right       interface{}         `json:"right"`
	Conjunction LogicalOperatorType `json:"conjunction"`
}

// NewComparisonExpression creates and new tag matcher expression for a comparision
func NewComparisonExpression(key string, entity MatcherExpressionEntity, operator ExpressionOperator, value string) MatchExpression {
	return &TagMatcherExpression{
		Dtype:    LeafExpressionType,
		Key:      key,
		Entity:   entity,
		Operator: operator,
		Value:    &value,
	}
}

// NewUnaryOperationExpression creates and new tag matcher expression for a unary operation
func NewUnaryOperationExpression(key string, entity MatcherExpressionEntity, operator ExpressionOperator) MatchExpression {
	return &TagMatcherExpression{
		Dtype:    LeafExpressionType,
		Key:      key,
		Entity:   entity,
		Operator: operator,
	}
}

// MatcherExpressionEntity type representing the matcher expression entity of a Matcher Expression (either source or destination or not applicable)
type MatcherExpressionEntity string

// MatcherExpressionEntities custom type representing a slice of MatcherExpressionEntity
type MatcherExpressionEntities []MatcherExpressionEntity

// ToStringSlice Returns the string representations fo the aggregations
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

// SupportedMatcherExpressionEntities slice of supported matcher expression entity types
var SupportedMatcherExpressionEntities = MatcherExpressionEntities{MatcherExpressionEntitySource, MatcherExpressionEntityDestination, MatcherExpressionEntityNotApplicable}

// TagMatcherExpression is the representation of a tag matcher expression in Instana
type TagMatcherExpression struct {
	Dtype    MatchExpressionType     `json:"type"`
	Key      string                  `json:"key"`
	Entity   MatcherExpressionEntity `json:"entity"`
	Operator ExpressionOperator      `json:"operator"`
	Value    *string                 `json:"value"`
}

// ApplicationConfig is the representation of a application perspective configuration in Instana
type ApplicationConfig struct {
	ID                  string                 `json:"id"`
	Label               string                 `json:"label"`
	MatchSpecification  interface{}            `json:"matchSpecification"`
	TagFilterExpression interface{}            `json:"tagFilterExpression"`
	Scope               ApplicationConfigScope `json:"scope"`
	BoundaryScope       BoundaryScope          `json:"boundaryScope"`
}

// GetIDForResourcePath implementation of the interface InstanaDataObject
func (a *ApplicationConfig) GetIDForResourcePath() string {
	return a.ID
}

// GetType implementation of the interface MatchExpression for BinaryOperator
func (b *BinaryOperator) GetType() MatchExpressionType {
	return b.Dtype
}

// GetType implementation of the interface MatchExpression for TagMatcherExpression
func (t *TagMatcherExpression) GetType() MatchExpressionType {
	return t.Dtype
}
