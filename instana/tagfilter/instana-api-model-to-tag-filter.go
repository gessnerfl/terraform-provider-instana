package tagfilter

import (
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/utils"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

//FromAPIModel Implementation of the mapping from the Instana API model to the filter expression model
func (m *tagFilterMapper) FromAPIModel(input restapi.TagFilterExpressionElement) (*FilterExpression, error) {
	expr, err := m.mapExpressionElement(input)
	if err != nil {
		return nil, err
	}
	if expr.or != nil {
		return &FilterExpression{Expression: expr.or}, nil
	}
	if expr.and != nil {
		return &FilterExpression{Expression: &LogicalOrExpression{Left: expr.and}}, nil
	}
	return &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: expr.primary,
			},
		},
	}, nil
}

func (m *tagFilterMapper) mapExpressionElement(input restapi.TagFilterExpressionElement) (*expressionHandle, error) {
	if input.GetType() == restapi.TagFilterExpressionType {
		expression := input.(*restapi.TagFilterExpression)
		return m.mapExpression(expression)
	} else if input.GetType() == restapi.TagFilterType {
		tagFilter := input.(*restapi.TagFilter)
		primaryExpression, err := m.mapTagFilter(tagFilter)
		if err != nil {
			return nil, err
		}
		return &expressionHandle{
			primary: primaryExpression,
		}, nil
	}
	return nil, fmt.Errorf("unsupported tag filter expression of type %s", input.GetType())
}

func (m *tagFilterMapper) mapExpression(operator *restapi.TagFilterExpression) (*expressionHandle, error) {
	elements := make([]*expressionHandle, len(operator.Elements))
	var err error
	for i := 0; i < len(operator.Elements); i++ {
		elements[i], err = m.mapExpressionElement(operator.Elements[i])
		if err != nil {
			return nil, err
		}
	}

	if operator.LogicalOperator == restapi.LogicalAnd {
		return m.mapLogicalAndFromAPIModel(elements)
	}
	if operator.LogicalOperator == restapi.LogicalOr {
		return m.mapLogicalOr(elements)
	}
	return nil, fmt.Errorf("invalid logical operator %s", operator.LogicalOperator)

}

func (m *tagFilterMapper) mapLogicalOr(elements []*expressionHandle) (*expressionHandle, error) {
	total := len(elements)
	if total < 2 {
		return nil, fmt.Errorf("at least two elements are expected for logical or")
	}
	if elements[0].or != nil {
		return nil, fmt.Errorf("invalid logical or expression: logical or is not allowed for first element")
	}

	operator := Operator(restapi.LogicalOr)
	var expression *LogicalOrExpression

	for i := total - 2; i >= 0; i-- {
		if expression == nil {
			expression = &LogicalOrExpression{
				Left:     m.mapLeftOfLogicalOr(elements[i]),
				Operator: &operator,
				Right:    m.mapRightOfLogicalOr(elements[i+1]),
			}
		} else {
			expression = &LogicalOrExpression{
				Left:     m.mapLeftOfLogicalOr(elements[i]),
				Operator: &operator,
				Right:    expression,
			}
		}
	}

	return &expressionHandle{or: expression}, nil
}

func (m *tagFilterMapper) mapLeftOfLogicalOr(left *expressionHandle) *LogicalAndExpression {
	if left.and != nil {
		return left.and
	}
	return &LogicalAndExpression{
		Left: left.primary,
	}
}

func (m *tagFilterMapper) mapRightOfLogicalOr(right *expressionHandle) *LogicalOrExpression {
	if right.or != nil {
		return right.or
	} else if right.and != nil {
		return &LogicalOrExpression{Left: right.and}
	} else {
		return &LogicalOrExpression{Left: &LogicalAndExpression{Left: right.primary}}
	}
}

func (m *tagFilterMapper) mapLogicalAndFromAPIModel(elements []*expressionHandle) (*expressionHandle, error) {
	total := len(elements)
	if total < 2 {
		return nil, fmt.Errorf("at least two elements are expected for logical and")
	}
	for _, e := range elements {
		if e.or != nil {
			return nil, fmt.Errorf("invalid logical and expression: logical or is not allowed")
		}
	}
	if elements[0].and != nil {
		return nil, fmt.Errorf("invalid logical and expression: logical and is not allowed for first element")
	}

	operator := Operator(restapi.LogicalAnd)
	var expression *LogicalAndExpression

	for i := total - 2; i >= 0; i-- {
		if expression == nil {
			expression = &LogicalAndExpression{
				Left:     elements[i].primary,
				Operator: &operator,
				Right:    m.mapRightOfLogicalAnd(elements[i+1]),
			}
		} else {
			expression = &LogicalAndExpression{
				Left:     elements[i].primary,
				Operator: &operator,
				Right:    expression,
			}
		}
	}
	return &expressionHandle{and: expression}, nil
}

func (m *tagFilterMapper) mapRightOfLogicalAnd(right *expressionHandle) *LogicalAndExpression {
	if right.and != nil {
		return right.and
	}
	return &LogicalAndExpression{Left: right.primary}
}

func (m *tagFilterMapper) mapTagFilter(tagFilter *restapi.TagFilter) (*PrimaryExpression, error) {
	origin := SupportedEntityOrigins.ForInstanaAPIEntity(tagFilter.Entity)
	if restapi.SupportedUnaryExpressionOperators.IsSupported(tagFilter.Operator) {
		return &PrimaryExpression{
			UnaryOperation: &UnaryOperationExpression{
				Entity:   &EntitySpec{Identifier: tagFilter.Name, TagKey: tagFilter.Key, Origin: utils.StringPtr(origin.Key())},
				Operator: Operator(tagFilter.Operator),
			},
		}, nil
	}
	if !restapi.SupportedComparisonOperators.IsSupported(tagFilter.Operator) {
		return nil, fmt.Errorf("invalid operator: %s is not a supported tag filter operator", tagFilter.Operator)
	}
	return &PrimaryExpression{
		Comparison: &ComparisonExpression{
			Entity:       &EntitySpec{Identifier: tagFilter.Name, TagKey: tagFilter.Key, Origin: utils.StringPtr(origin.Key())},
			Operator:     Operator(tagFilter.Operator),
			StringValue:  m.mapStringOrTagValue(tagFilter),
			BooleanValue: tagFilter.BooleanValue,
			NumberValue:  tagFilter.NumberValue,
		},
	}, nil
}

func (m *tagFilterMapper) mapStringOrTagValue(tagFilter *restapi.TagFilter) *string {
	if tagFilter.Key != nil {
		tagValue := tagFilter.Value.(string)
		return &tagValue
	}
	return tagFilter.StringValue
}

type expressionHandle struct {
	or      *LogicalOrExpression
	and     *LogicalAndExpression
	primary *PrimaryExpression
}
