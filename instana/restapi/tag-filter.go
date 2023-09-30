package restapi

import (
	"fmt"
)

// TagFilterExpressionElementType type for TagFilterExpressionElement discriminator type
type TagFilterExpressionElementType string

const (
	//TagFilterExpressionType discriminator type for expression TagFilterExpressionElement
	TagFilterExpressionType TagFilterExpressionElementType = "EXPRESSION"
	//TagFilterType discriminator type for leaf tag_filter TagFilterExpressionElementType
	TagFilterType TagFilterExpressionElementType = "TAG_FILTER"
)

// NewLogicalOrTagFilter creates a new logical OR expression
func NewLogicalOrTagFilter(elements []*TagFilter) *TagFilter {
	operator := LogicalOr
	return &TagFilter{
		Type:            TagFilterExpressionType,
		LogicalOperator: &operator,
		Elements:        elements,
	}
}

// NewLogicalAndTagFilter creates a new logical AND expression
func NewLogicalAndTagFilter(elements []*TagFilter) *TagFilter {
	operator := LogicalAnd
	return &TagFilter{
		Type:            TagFilterExpressionType,
		LogicalOperator: &operator,
		Elements:        elements,
	}
}

// TagFilterEntity type representing the matcher expression entity of a Matcher Expression (either source or destination or not applicable)
type TagFilterEntity string

// TagFilterEntities custom type representing a slice of TagFilterEntity
type TagFilterEntities []TagFilterEntity

// ToStringSlice Returns the string representations fo the aggregations
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

// SupportedTagFilterEntities slice of supported matcher expression entity types
var SupportedTagFilterEntities = TagFilterEntities{TagFilterEntitySource, TagFilterEntityDestination, TagFilterEntityNotApplicable}

// NewStringTagFilter creates a new TagFilter for comparing string values
func NewStringTagFilter(entity TagFilterEntity, name string, operator ExpressionOperator, value string) *TagFilter {
	return &TagFilter{
		Entity:      &entity,
		Name:        &name,
		Operator:    &operator,
		StringValue: &value,
		Value:       value,
		Type:        TagFilterType,
	}
}

// NewNumberTagFilter creates a new TagFilter for comparing number values
func NewNumberTagFilter(entity TagFilterEntity, name string, operator ExpressionOperator, value int64) *TagFilter {
	return &TagFilter{
		Entity:      &entity,
		Name:        &name,
		Operator:    &operator,
		NumberValue: &value,
		Value:       value,
		Type:        TagFilterType,
	}
}

// NewTagTagFilter creates a new TagFilter for comparing tags
func NewTagTagFilter(entity TagFilterEntity, name string, operator ExpressionOperator, key string, value string) *TagFilter {
	fullString := fmt.Sprintf("%s=%s", key, value)
	return &TagFilter{
		Entity:      &entity,
		Name:        &name,
		Operator:    &operator,
		Key:         &key,
		Value:       value,
		StringValue: &fullString,
		Type:        TagFilterType,
	}
}

// NewBooleanTagFilter creates a new TagFilter for comparing tags
func NewBooleanTagFilter(entity TagFilterEntity, name string, operator ExpressionOperator, value bool) *TagFilter {
	return &TagFilter{
		Entity:       &entity,
		Name:         &name,
		Operator:     &operator,
		BooleanValue: &value,
		Value:        value,
		Type:         TagFilterType,
	}
}

// NewUnaryTagFilter creates a new TagFilter for unary expressions
func NewUnaryTagFilter(entity TagFilterEntity, name string, operator ExpressionOperator) *TagFilter {
	return &TagFilter{
		Entity:   &entity,
		Name:     &name,
		Operator: &operator,
		Type:     TagFilterType,
	}
}

// NewUnaryTagFilterWithTagKey creates a new TagFilter for unary expressions supporting tagKeys
func NewUnaryTagFilterWithTagKey(entity TagFilterEntity, name string, tagKey *string, operator ExpressionOperator) *TagFilter {
	return &TagFilter{
		Entity:   &entity,
		Name:     &name,
		Key:      tagKey,
		Operator: &operator,
		Type:     TagFilterType,
	}
}

// TagFilter data structure of a Tag Filter from the Instana API
type TagFilter struct {
	Elements        []*TagFilter                   `json:"elements"`
	LogicalOperator *LogicalOperatorType           `json:"logicalOperator"`
	Entity          *TagFilterEntity               `json:"entity"`
	Name            *string                        `json:"name"`
	Operator        *ExpressionOperator            `json:"operator"`
	BooleanValue    *bool                          `json:"booleanValue"`
	NumberValue     *int64                         `json:"numberValue"`
	StringValue     *string                        `json:"stringValue"`
	Key             *string                        `json:"key"`
	Value           interface{}                    `json:"value"`
	Type            TagFilterExpressionElementType `json:"type"`
}

// GetType Implementation of the TagFilterExpressionElement type
func (f *TagFilter) GetType() TagFilterExpressionElementType {
	return f.Type
}

// PrependElement adds a TagFilterExpressionElement to the end of the list of elements
func (e *TagFilter) PrependElement(element *TagFilter) {
	e.Elements = append(e.Elements, element)
}
