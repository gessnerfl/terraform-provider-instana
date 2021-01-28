package filterexpression

import (
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

//FromAPIModel Implementation of the mapping from the Instana API model to the filter expression model
func (m *mapperImpl) FromAPIModel(input restapi.MatchExpression) (*FilterExpression, error) {
	expr, err := m.mapExpressionFromAPIModel(input)
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

func (m *mapperImpl) mapExpressionFromAPIModel(input restapi.MatchExpression) (*expressionHandle, error) {
	if input.GetType() == restapi.BinaryOperatorExpressionType {
		binaryOp := input.(restapi.BinaryOperator)
		return m.mapBinaryOperatorFromAPIModel(&binaryOp)
	} else if input.GetType() == restapi.LeafExpressionType {
		tagMatcher := input.(restapi.TagMatcherExpression)
		primaryExpression, err := m.mapPrimaryExpressionFromAPIModel(&tagMatcher)
		if err != nil {
			return nil, err
		}
		return &expressionHandle{
			primary: primaryExpression,
		}, nil
	}
	return nil, fmt.Errorf("unsupported match expression of type %s", input.GetType())
}

func (m *mapperImpl) mapBinaryOperatorFromAPIModel(operator *restapi.BinaryOperator) (*expressionHandle, error) {
	left, err := m.mapExpressionFromAPIModel(operator.Left.(restapi.MatchExpression))
	if err != nil {
		return nil, err
	}
	right, err := m.mapExpressionFromAPIModel(operator.Right.(restapi.MatchExpression))
	if err != nil {
		return nil, err
	}

	if operator.Conjunction == restapi.LogicalAnd {
		return m.mapLogicalAndFromAPIModel(left, right)
	}
	if operator.Conjunction == restapi.LogicalOr {
		return m.mapLogicalOrFromAPIModel(left, right)
	}
	return nil, fmt.Errorf("invalid conjunction operator %s", operator.Conjunction)

}

func (m *mapperImpl) mapLogicalOrFromAPIModel(left *expressionHandle, right *expressionHandle) (*expressionHandle, error) {
	if left.or != nil {
		return nil, fmt.Errorf("invalid logical or expression: logical or is not allowed for left side")
	}

	operator := Operator(restapi.LogicalOr)
	return &expressionHandle{
		or: &LogicalOrExpression{
			Left:     m.mapLeftOfLogicalOrFromAPIModel(left),
			Operator: &operator,
			Right:    m.mapRightOfLogicalOrFromAPIModel(right),
		},
	}, nil
}

func (m *mapperImpl) mapLeftOfLogicalOrFromAPIModel(left *expressionHandle) *LogicalAndExpression {
	if left.and != nil {
		return left.and
	}
	return &LogicalAndExpression{
		Left: left.primary,
	}
}

func (m *mapperImpl) mapRightOfLogicalOrFromAPIModel(right *expressionHandle) *LogicalOrExpression {
	if right.or != nil {
		return right.or
	} else if right.and != nil {
		return &LogicalOrExpression{Left: right.and}
	} else {
		return &LogicalOrExpression{Left: &LogicalAndExpression{Left: right.primary}}
	}
}

func (m *mapperImpl) mapLogicalAndFromAPIModel(left *expressionHandle, right *expressionHandle) (*expressionHandle, error) {
	if left.or != nil {
		return nil, fmt.Errorf("invalid logical and expression: logical or is not allowed for left side")
	}

	if right.or != nil {
		return nil, fmt.Errorf("invalid logical and expression: logical or is not allowed for right side")
	}

	if left.and != nil {
		return nil, fmt.Errorf("invalid logical and expression: logical and is not allowed for left side")
	}

	operator := Operator(restapi.LogicalAnd)
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

func (m *mapperImpl) mapPrimaryExpressionFromAPIModel(matcher *restapi.TagMatcherExpression) (*PrimaryExpression, error) {
	origin, err := m.mapOrigin(matcher.Entity)
	if err != nil {
		return nil, err
	}
	if matcher.Value != nil {
		if !restapi.IsSupportedComparision(matcher.Operator) {
			return nil, fmt.Errorf("invalid operator: %s is not a supported comparision operator", matcher.Operator)
		}
		return &PrimaryExpression{
			Comparision: &ComparisionExpression{
				Entity:   &EntitySpec{Key: matcher.Key, Origin: origin},
				Operator: Operator(matcher.Operator),
				Value:    *matcher.Value,
			},
		}, nil
	}
	if !restapi.IsSupportedUnaryOperatorExpression(matcher.Operator) {
		return nil, fmt.Errorf("invalid operator: %s is not a supported unary operator", matcher.Operator)
	}
	return &PrimaryExpression{
		UnaryOperation: &UnaryOperationExpression{
			Entity:   &EntitySpec{Key: matcher.Key, Origin: origin},
			Operator: Operator(matcher.Operator),
		},
	}, nil
}

func (m *mapperImpl) mapOrigin(entity restapi.MatcherExpressionEntity) (EntityOrigin, error) {
	if entity == restapi.MatcherExpressionEntityDestination {
		return EntityOriginDestination, nil
	} else if entity == restapi.MatcherExpressionEntityDestination {
		return EntityOriginSource, nil
	}
	return EntityOriginDestination, fmt.Errorf("Invalid entity origin: %s is not a supported entity of an entity match specification leaf element", entity)
}

type expressionHandle struct {
	or      *LogicalOrExpression
	and     *LogicalAndExpression
	primary *PrimaryExpression
}
