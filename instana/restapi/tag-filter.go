package restapi

import (
	"errors"
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"strings"
)

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

//LogicalOperatorTypes custom type for slice of logical operators
type LogicalOperatorTypes []LogicalOperatorType

//IsSupported check if the provided logical operator is supported
func (operators LogicalOperatorTypes) IsSupported(o LogicalOperatorType) bool {
	for _, v := range operators {
		if v == o {
			return true
		}
	}
	return false
}

const (
	//LogicalAnd constant for logical AND conjunction
	LogicalAnd = LogicalOperatorType("AND")
	//LogicalOr constant for logical OR conjunction
	LogicalOr = LogicalOperatorType("OR")
)

//SupportedLogicalOperatorTypes list of supported logical operators of Instana API
var SupportedLogicalOperatorTypes = LogicalOperatorTypes{LogicalAnd, LogicalOr}

//TagFilterExpressionElement interface for the Instana API type TagFilterExpressionElement
type TagFilterExpressionElement interface {
	GetType() TagFilterExpressionElementType
	Validate() error
}

//NewLogicalOrTagFilter creates a new logical OR expression
func NewLogicalOrTagFilter(elements []TagFilterExpressionElement) *TagFilterExpression {
	return &TagFilterExpression{
		Type:            TagFilterExpressionType,
		LogicalOperator: LogicalOr,
		Elements:        elements,
	}
}

//NewLogicalAndTagFilter creates a new logical AND expression
func NewLogicalAndTagFilter(elements []TagFilterExpressionElement) *TagFilterExpression {
	return &TagFilterExpression{
		Type:            TagFilterExpressionType,
		LogicalOperator: LogicalAnd,
		Elements:        elements,
	}
}

//TagFilterExpression data structure of an Instana tag filter expression
type TagFilterExpression struct {
	Elements        []TagFilterExpressionElement   `json:"elements"`
	LogicalOperator LogicalOperatorType            `json:"logicalOperator"`
	Type            TagFilterExpressionElementType `json:"type"`
}

//GetType Implementation of the TagFilterExpressionElement type
func (e *TagFilterExpression) GetType() TagFilterExpressionElementType {
	return e.Type
}

//Validate Implementation of the TagFilterExpressionElement type
func (e *TagFilterExpression) Validate() error {
	if len(e.Elements) < 2 {
		return errors.New("at least two elements are expected for a tag filter expression")
	}
	if !SupportedLogicalOperatorTypes.IsSupported(e.LogicalOperator) {
		return fmt.Errorf("tag filter operator %s is not supported", e.LogicalOperator)
	}
	if strings.ToUpper(string(e.Type)) != string(TagFilterExpressionType) {
		return fmt.Errorf("tag filter expression must be of type EXPRESSION but %s is provided", e.Type)
	}
	for _, element := range e.Elements {
		err := element.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

//PrependElement adds a TagFilterExpressionElement to the end of the list of elements
func (e *TagFilterExpression) PrependElement(element TagFilterExpressionElement) {
	e.Elements = append([]TagFilterExpressionElement{element}, e.Elements...)
}

//TagFilterEntity type representing the matcher expression entity of a Matcher Expression (either source or destination or not applicable)
type TagFilterEntity string

//TagFilterEntities custom type representing a slice of TagFilterEntity
type TagFilterEntities []TagFilterEntity

//IsSupported check if the provided tag filter entity is supported
func (entities TagFilterEntities) IsSupported(entity TagFilterEntity) bool {
	for _, v := range entities {
		if v == entity {
			return true
		}
	}
	return false
}

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

//TagFilterOperator custom type for tag matcher operators
type TagFilterOperator string

//TagFilterOperators custom type representing a slice of TagFilterOperator
type TagFilterOperators []TagFilterOperator

//IsSupported check if the provided tag filter operator is supported
func (operators TagFilterOperators) IsSupported(o TagFilterOperator) bool {
	for _, v := range operators {
		if v == o {
			return true
		}
	}
	return false
}

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
var SupportedComparisonOperators = TagFilterOperators{
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
var SupportedUnaryExpressionOperators = TagFilterOperators{
	IsEmptyOperator,
	NotEmptyOperator,
	IsBlankOperator,
	NotBlankOperator,
}

//NewStringTagFilter creates a new TagFilter for comparing string values
func NewStringTagFilter(entity TagFilterEntity, name string, operator TagFilterOperator, value string) *TagFilter {
	return &TagFilter{
		Entity:      entity,
		Name:        name,
		Operator:    operator,
		StringValue: &value,
		Value:       value,
		Type:        TagFilterType,
	}
}

//NewNumberTagFilter creates a new TagFilter for comparing number values
func NewNumberTagFilter(entity TagFilterEntity, name string, operator TagFilterOperator, value int64) *TagFilter {
	return &TagFilter{
		Entity:      entity,
		Name:        name,
		Operator:    operator,
		NumberValue: &value,
		Value:       value,
		Type:        TagFilterType,
	}
}

//NewTagTagFilter creates a new TagFilter for comparing tags
func NewTagTagFilter(entity TagFilterEntity, name string, operator TagFilterOperator, key string, value string) *TagFilter {
	fullString := fmt.Sprintf("%s=%s", key, value)
	return &TagFilter{
		Entity:      entity,
		Name:        name,
		Operator:    operator,
		Key:         &key,
		Value:       value,
		StringValue: &fullString,
		Type:        TagFilterType,
	}
}

//NewBooleanTagFilter creates a new TagFilter for comparing tags
func NewBooleanTagFilter(entity TagFilterEntity, name string, operator TagFilterOperator, value bool) *TagFilter {
	return &TagFilter{
		Entity:       entity,
		Name:         name,
		Operator:     operator,
		BooleanValue: &value,
		Value:        value,
		Type:         TagFilterType,
	}
}

//NewUnaryTagFilter creates a new TagFilter for unary expressions
func NewUnaryTagFilter(entity TagFilterEntity, name string, operator TagFilterOperator) *TagFilter {
	return &TagFilter{
		Entity:   entity,
		Name:     name,
		Operator: operator,
		Type:     TagFilterType,
	}
}

//TagFilter data structure of a Tag Filter from the Instana API
type TagFilter struct {
	Entity       TagFilterEntity                `json:"entity"`
	Name         string                         `json:"name"`
	Operator     TagFilterOperator              `json:"operator"`
	BooleanValue *bool                          `json:"booleanValue"`
	NumberValue  *int64                         `json:"numberValue"`
	StringValue  *string                        `json:"stringValue"`
	Key          *string                        `json:"key"`
	Value        interface{}                    `json:"value"`
	Type         TagFilterExpressionElementType `json:"type"`
}

//GetType Implementation of the TagFilterExpressionElement type
func (f *TagFilter) GetType() TagFilterExpressionElementType {
	return f.Type
}

//Validate Implementation of the TagFilterExpressionElement type
func (f *TagFilter) Validate() error {
	if !SupportedTagFilterEntities.IsSupported(f.Entity) {
		return fmt.Errorf("tag filter entity type %s is not supported", f.Entity)
	}
	if utils.IsBlank(f.Name) {
		return errors.New("tag filter name is missing")
	}
	isSupportedComparisonOperation := SupportedComparisonOperators.IsSupported(f.Operator)
	isSupportedUnaryOperation := SupportedUnaryExpressionOperators.IsSupported(f.Operator)
	if !isSupportedUnaryOperation && !isSupportedComparisonOperation {
		return fmt.Errorf("tag filter operator %s is not supported", f.Operator)
	}
	if isSupportedComparisonOperation && !f.isValueAssigned() {
		return errors.New("value missing for comparison operation")
	}
	if isSupportedUnaryOperation && f.isValueAssigned() {
		return errors.New("no value must be assigned for unary operation")
	}
	return nil
}

func (f *TagFilter) isValueAssigned() bool {
	return f.Value != nil
}
