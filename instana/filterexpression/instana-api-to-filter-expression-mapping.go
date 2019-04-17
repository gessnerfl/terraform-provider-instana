package filterexpression

import (
	"errors"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

//OperatorMappingInstanaAPIToFilterExpression map defining the mapping of the operators from instana API to the representation of the Filter Expression
var OperatorMappingInstanaAPIToFilterExpression = map[string]Operator{
	"EQUALS":      Operator("EQ"),
	"NOT_EQUAL":   Operator("NE"),
	"CONTAINS":    Operator("CO"),
	"NOT_CONTAIN": Operator("NC"),
}

//UnaryOperatorMappingInstanaAPIToFilterExpression map defining the mapping of the unary operator names from instana API to the representation of the Filter Expression
var UnaryOperatorMappingInstanaAPIToFilterExpression = map[string]UnaryOperator{
	"IS_EMPTY":  UnaryOperator("IS EMPTY"),
	"NOT_EMPTY": UnaryOperator("NOT EMPTY"),
	"IS_BLANK":  UnaryOperator("IS BLANK"),
	"NOT_BLANK": UnaryOperator("NOT BLANK"),
}

//FromAPIModel Implementation of the mapping from the Instana API model to the filter expression model
func (m *mapperImpl) FromAPIModel(input restapi.MatchExpression) (*FilterExpression, error) {
	expr, err := m.mapExpression(input)
	if err != nil {
		return nil, err
	}
	if expr.or != nil {
		return &FilterExpression{Expression: expr.or}, nil
	} else if expr.and != nil {
		return &FilterExpression{Expression: &LogicalOrExpression{Left: expr.and}}, nil
	} else if expr.primary != nil {
		return &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: expr.primary,
				},
			},
		}, nil
	}
	return nil, errors.New("expected exactly one expression to be returned")
}

func (m *mapperImpl) mapExpression(input restapi.MatchExpression) (*expressionHandle, error) {
	if input.GetType() == restapi.BinaryOperatorExpressionType {
		binaryOp := input.(restapi.BinaryOperator)
		return m.mapBinaryOperator(&binaryOp)
	} else if input.GetType() == restapi.LeafExpressionType {
		tagMatcher := input.(restapi.TagMatcherExpression)
		primaryExpression, err := m.mapPrimaryExpression(&tagMatcher)
		if err != nil {
			return nil, err
		}
		return &expressionHandle{
			primary: primaryExpression,
		}, nil
	}
	return nil, fmt.Errorf("unsupported match expression of type %s", input.GetType())
}

func (m *mapperImpl) mapBinaryOperator(operator *restapi.BinaryOperator) (*expressionHandle, error) {
	left, err := m.mapExpression(operator.Left.(restapi.MatchExpression))
	if err != nil {
		return nil, err
	}
	right, err := m.mapExpression(operator.Right.(restapi.MatchExpression))
	if err != nil {
		return nil, err
	}

	if operator.Conjunction == "AND" {
		return m.mapLogicalAnd(left, right)
	}
	if operator.Conjunction == "OR" {
		return m.mapLogicalOr(left, right)
	}
	return nil, fmt.Errorf("Invalid conjunction operator %s", operator.Conjunction)

}

func (m *mapperImpl) mapLogicalOr(left *expressionHandle, right *expressionHandle) (*expressionHandle, error) {
	if left.or != nil {
		return nil, fmt.Errorf("invalid logical or expression: logical or is not allowed for left side")
	}

	operator := Operator("OR")
	return &expressionHandle{
		or: &LogicalOrExpression{
			Left:     m.mapLeftOfLogicalOr(left),
			Operator: &operator,
			Right:    m.mapRightOfLogicalOr(right),
		},
	}, nil
}

func (m *mapperImpl) mapLeftOfLogicalOr(left *expressionHandle) *LogicalAndExpression {
	if left.and != nil {
		return left.and
	}
	return &LogicalAndExpression{
		Left: left.primary,
	}
}

func (m *mapperImpl) mapRightOfLogicalOr(right *expressionHandle) *LogicalOrExpression {
	if right.or != nil {
		return right.or
	} else if right.and != nil {
		return &LogicalOrExpression{Left: right.and}
	} else {
		return &LogicalOrExpression{Left: &LogicalAndExpression{Left: right.primary}}
	}
}

func (m *mapperImpl) mapLogicalAnd(left *expressionHandle, right *expressionHandle) (*expressionHandle, error) {
	if left.or != nil {
		return nil, fmt.Errorf("invalid logical and expression: logical or is not allowed for left side")
	}

	if right.or != nil {
		return nil, fmt.Errorf("invalid logical and expression: logical or is not allowed for right side")
	}

	if left.and != nil {
		return nil, fmt.Errorf("invalid logical and expression: logical and is not allowed for left side")
	}

	operator := Operator("AND")
	if right.and != nil {
		return &expressionHandle{
			and: &LogicalAndExpression{
				Left:     left.primary,
				Operator: &operator,
				Right:    right.and,
			},
		}, nil
	}

	return &expressionHandle{
		and: &LogicalAndExpression{
			Left:     left.primary,
			Operator: &operator,
			Right:    &LogicalAndExpression{Left: right.primary},
		},
	}, nil
}

func (m *mapperImpl) mapPrimaryExpression(matcher *restapi.TagMatcherExpression) (*PrimaryExpression, error) {
	if matcher.Value != nil {
		operator, err := m.mapOperator(matcher.Operator)
		if err != nil {
			return nil, err
		}
		return &PrimaryExpression{
			Comparision: &ComparisionExpression{
				Key:      matcher.Key,
				Operator: operator,
				Value:    *matcher.Value,
			},
		}, nil
	}

	unaryOperator, err := m.mapUnaryOperator(matcher.Operator)
	if err != nil {
		return nil, err
	}
	return &PrimaryExpression{
		UnaryOperation: &UnaryOperationExpression{
			Key:      matcher.Key,
			Operator: unaryOperator,
		},
	}, nil
}

func (m *mapperImpl) mapOperator(apiName string) (Operator, error) {
	value := OperatorMappingInstanaAPIToFilterExpression[apiName]
	if value == "" {
		return value, fmt.Errorf("invalid operation: operation '%s' not supported", apiName)
	}
	return value, nil
}

func (m *mapperImpl) mapUnaryOperator(apiName string) (UnaryOperator, error) {
	value := UnaryOperatorMappingInstanaAPIToFilterExpression[apiName]
	if value == "" {
		return value, fmt.Errorf("invalid unary operation: unary operator '%s' not supported", apiName)
	}
	return value, nil
}

type expressionHandle struct {
	or      *LogicalOrExpression
	and     *LogicalAndExpression
	primary *PrimaryExpression
}
