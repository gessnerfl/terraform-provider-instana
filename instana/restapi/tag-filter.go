package restapi

//TagFilterExpressionElementType type for TagFilterExpressionElement discriminator type
type TagFilterExpressionElementType string

const (
	//TagFilterExpressionType discriminator type for expression TagFilterExpressionElement
	TagFilterExpressionType TagFilterExpressionElementType = "EXPRESSION"
	//TagFilterType discriminator type for leaf tag_filter TagFilterExpressionElementType
	TagFilterType TagFilterExpressionElementType = "TAG_FILTER"
)

//LogicalOperatorType custom type for logical operators
type LogicalOperatorType string

const (
	//LogicalAnd constant for logical AND conjunction
	LogicalAnd = LogicalOperatorType("AND")
	//LogicalOr constant for logical OR conjunction
	LogicalOr = LogicalOperatorType("OR")
)

//TagFilterExpressionElement interface for the Instana API type TagFilterExpressionElement
type TagFilterExpressionElement interface {
	GetType() TagFilterExpressionElementType
	Validate() error
}

func NewLogicalOrTagFilter(elements []TagFilterExpressionElement) *TagFilterExpression {
	return &TagFilterExpression{
		Type:            TagFilterExpressionType,
		LogicalOperator: LogicalOr,
		Elements:        elements,
	}
}

func NewLogicalAndTagFilter(elements []TagFilterExpressionElement) *TagFilterExpression {
	return &TagFilterExpression{
		Type:            TagFilterExpressionType,
		LogicalOperator: LogicalAnd,
		Elements:        elements,
	}
}

//TagFilterExpression data structure of an Instana tag filter expression
type TagFilterExpression struct {
	Elements        []TagFilterExpressionElement
	LogicalOperator LogicalOperatorType
	Type            TagFilterExpressionElementType
}

//GetType Implementation of the TagFilterExpressionElement type
func (e *TagFilterExpression) GetType() TagFilterExpressionElementType {
	return e.Type
}

//Validate Implementation of the TagFilterExpressionElement type
func (e *TagFilterExpression) Validate() error {
	//TODO add implementation
	return nil
}

func (e *TagFilterExpression) PrependElement(element TagFilterExpressionElement) {
	e.Elements = append([]TagFilterExpressionElement{element}, e.Elements...)
}

//TagFilterEntity type representing the matcher expression entity of a Matcher Expression (either source or destination or not applicable)
type TagFilterEntity string

//TagFilterEntities custom type representing a slice of TagFilterEntity
type TagFilterEntities []TagFilterEntity

//ToStringSlice Returns the string representations fo the aggregations
func (entities TagFilterEntities) ToStringSlice() []string {
	result := make([]string, len(entities))
	for i, v := range entities {
		result[i] = string(v)
	}
	return result
}

const (
	//TagFilterEntitySource const for a SOURCE matcher expression entity
	TagFilterEntitySource = TagFilterEntity("SOURCE")
	//TagFilterEntityDestination const for a DESTINATION matcher expression entity
	TagFilterEntityDestination = TagFilterEntity("DESTINATION")
	//TagFilterEntityNotApplicable const for a NOT_APPLICABLE matcher expression entity
	TagFilterEntityNotApplicable = TagFilterEntity("NOT_APPLICABLE")
)

//SupportedTagFilterEntities slice of supported matcher expression entity types
var SupportedTagFilterEntities = TagFilterEntities{TagFilterEntitySource, TagFilterEntityDestination, TagFilterEntityNotApplicable}

//IsSupported check if the provided matcher expression entity is supported
func (entities TagFilterEntities) IsSupported(entity TagFilterEntity) bool {
	for _, v := range entities {
		if v == entity {
			return true
		}
	}
	return false
}

//TagFilterOperator custom type for tag matcher operators
type TagFilterOperator string

const (
	//EqualsOperator constant for the EQUALS operator
	EqualsOperator = TagFilterOperator("EQUALS")
	//NotEqualOperator constant for the NOT_EQUAL operator
	NotEqualOperator = TagFilterOperator("NOT_EQUAL")
	//ContainsOperator constant for the CONTAINS operator
	ContainsOperator = TagFilterOperator("CONTAINS")
	//NotContainOperator constant for the NOT_CONTAIN operator
	NotContainOperator = TagFilterOperator("NOT_CONTAIN")
	//IsEmptyOperator constant for the IS_EMPTY operator
	IsEmptyOperator = TagFilterOperator("IS_EMPTY")
	//NotEmptyOperator constant for the NOT_EMPTY operator
	NotEmptyOperator = TagFilterOperator("NOT_EMPTY")
	//IsBlankOperator constant for the IS_BLANK operator
	IsBlankOperator = TagFilterOperator("IS_BLANK")
	//NotBlankOperator constant for the NOT_BLANK operator
	NotBlankOperator = TagFilterOperator("NOT_BLANK")

	//StartsWithOperator constant for the STARTS_WITH operator
	StartsWithOperator = TagFilterOperator("STARTS_WITH")
	//EndsWithOperator constant for the ENDS_WITH operator
	EndsWithOperator = TagFilterOperator("ENDS_WITH")
	//NotStartsWithOperator constant for the NOT_STARTS_WITH operator
	NotStartsWithOperator = TagFilterOperator("NOT_STARTS_WITH")
	//NotEndsWithOperator constant for the NOT_ENDS_WITH operator
	NotEndsWithOperator = TagFilterOperator("NOT_ENDS_WITH")
	//GreaterOrEqualThanOperator constant for the GREATER_OR_EQUAL_THAN operator
	GreaterOrEqualThanOperator = TagFilterOperator("GREATER_OR_EQUAL_THAN")
	//LessOrEqualThanOperator constant for the LESS_OR_EQUAL_THAN operator
	LessOrEqualThanOperator = TagFilterOperator("LESS_OR_EQUAL_THAN")
	//GreaterThanOperator constant for the GREATER_THAN operator
	GreaterThanOperator = TagFilterOperator("GREATER_THAN")
	//LessThanOperator constant for the LESS_THAN operator
	LessThanOperator = TagFilterOperator("LESS_THAN")
)

//SupportedComparisonOperators list of supported comparison operators of Instana API
var SupportedComparisonOperators = []TagFilterOperator{
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
var SupportedUnaryExpressionOperators = []TagFilterOperator{
	IsEmptyOperator,
	NotEmptyOperator,
	IsBlankOperator,
	NotBlankOperator,
}

func NewStringTagFilter(entity TagFilterEntity, name string, operator TagFilterOperator, value *string) *TagFilter {
	return &TagFilter{
		Entity:      entity,
		Name:        name,
		Operator:    operator,
		StringValue: value,
		Type:        TagFilterType,
	}
}

func NewNumberTagFilter(entity TagFilterEntity, name string, operator TagFilterOperator, value *int64) *TagFilter {
	return &TagFilter{
		Entity:      entity,
		Name:        name,
		Operator:    operator,
		NumberValue: value,
		Type:        TagFilterType,
	}
}

func NewTagTagFilter(entity TagFilterEntity, name string, operator TagFilterOperator, key *string, value *string) *TagFilter {
	return &TagFilter{
		Entity:   entity,
		Name:     name,
		Operator: operator,
		TagKey:   key,
		TagValue: value,
		Type:     TagFilterType,
	}
}

func NewBooleanTagFilter(entity TagFilterEntity, name string, operator TagFilterOperator, value *bool) *TagFilter {
	return &TagFilter{
		Entity:       entity,
		Name:         name,
		Operator:     operator,
		BooleanValue: value,
		Type:         TagFilterType,
	}
}

func NewUnaryTagFilter(entity TagFilterEntity, name string, operator TagFilterOperator) *TagFilter {
	return &TagFilter{
		Entity:   entity,
		Name:     name,
		Operator: operator,
		Type:     TagFilterType,
	}
}

type TagFilter struct {
	Entity       TagFilterEntity
	Name         string
	Operator     TagFilterOperator
	BooleanValue *bool
	NumberValue  *int64
	StringValue  *string
	TagKey       *string
	TagValue     *string
	Type         TagFilterExpressionElementType
}

//GetType Implementation of the TagFilterExpressionElement type
func (f *TagFilter) GetType() TagFilterExpressionElementType {
	return f.Type
}

//Validate Implementation of the TagFilterExpressionElement type
func (f *TagFilter) Validate() error {
	//TODO add implementation
	return nil
}
