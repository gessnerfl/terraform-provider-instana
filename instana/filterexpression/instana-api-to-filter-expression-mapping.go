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

//FunctionMappingInstanaAPIToFilterExpression map defining the mapping of the function names of unary operations from instana API to the representation of the Filter Expression
var FunctionMappingInstanaAPIToFilterExpression = map[string]Function{
	"IS_EMPTY":  Function("IS EMPTY"),
	"NOT_EMPTY": Function("NOT EMPTY"),
	"IS_BLANK":  Function("IS BLANK"),
	"NOT_BLANK": Function("NOT BLANK"),
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
		return nil, fmt.Errorf("invalid logical or expression: logical or is not allowed for the left side")
	}

	if left.and == nil && left.primary == nil {
		return nil, errors.New("invalid logical or expression: left side of logical or is not defined")
	}

	if right.or == nil && right.and == nil && right.primary == nil {
		return nil, errors.New("invalid logical or expression: right side of logical or is not defined")
	}

	operator := Operator("OR")
	orExpression := LogicalOrExpression{Operator: &operator}

	if left.and != nil {
		orExpression.Left = left.and
	} else {
		orExpression.Left = &LogicalAndExpression{
			Left: left.primary,
		}
	}

	if right.or != nil {
		orExpression.Right = right.or
	} else if right.and != nil {
		orExpression.Right = &LogicalOrExpression{Left: left.and}
	} else {
		orExpression.Right = &LogicalOrExpression{Left: &LogicalAndExpression{Left: right.primary}}
	}

	return &expressionHandle{or: &orExpression}, nil
}

func (m *mapperImpl) mapLogicalAnd(left *expressionHandle, right *expressionHandle) (*expressionHandle, error) {
	if left.or != nil {
		return nil, fmt.Errorf("invalid logical and expression: logical or is not allowed for the left side")
	}

	if right.or != nil {
		return nil, fmt.Errorf("invalid logical and expression: logical or is not allowed for the right side")
	}

	if left.and != nil {
		return nil, fmt.Errorf("invalid logical and expression: logical and is not allowed for left side")
	}

	if left.primary == nil {
		return nil, errors.New("invalid logical and expression: left side of logical and is not defined")
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

	if right.primary != nil {
		return &expressionHandle{
			and: &LogicalAndExpression{
				Left:     left.primary,
				Operator: &operator,
				Right:    &LogicalAndExpression{Left: right.primary},
			},
		}, nil
	}

	return nil, errors.New("invalid logical and expression: right side of logical and is not defined")
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

	function, err := m.mapFunction(matcher.Operator)
	if err != nil {
		return nil, err
	}
	return &PrimaryExpression{
		UnaryExpression: &UnaryExpression{
			Key:      matcher.Key,
			Function: function,
		},
	}, nil
}

func (m *mapperImpl) mapOperator(apiName string) (Operator, error) {
	value := OperatorMappingInstanaAPIToFilterExpression[apiName]
	if value == "" {
		return value, fmt.Errorf("invalid operation: operation %s not supported", apiName)
	}
	return value, nil
}

func (m *mapperImpl) mapFunction(apiName string) (Function, error) {
	value := FunctionMappingInstanaAPIToFilterExpression[apiName]
	if value == "" {
		return value, fmt.Errorf("invalid operation: operation %s not supported", apiName)
	}
	return value, nil
}

type expressionHandle struct {
	or      *LogicalOrExpression
	and     *LogicalAndExpression
	primary *PrimaryExpression
}
