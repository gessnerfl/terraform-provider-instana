package filterexpression

import "github.com/gessnerfl/terraform-provider-instana/instana/restapi"

//ToAPIModel Implementation of the mapping form filter expression model to the Instana API model
func (m *matchExpressionMapperImpl) ToAPIModel(input *FilterExpression) restapi.MatchExpression {
	return m.mapLogicalOrToAPIModel(input.Expression)
}

func (m *matchExpressionMapperImpl) mapLogicalOrToAPIModel(input *LogicalOrExpression) restapi.MatchExpression {
	left := m.mapLogicalAndToAPIModel(input.Left)
	if input.Operator != nil {
		right := m.mapLogicalOrToAPIModel(input.Right)
		return restapi.NewBinaryOperator(left, restapi.LogicalOr, right)
	}
	return left
}

func (m *matchExpressionMapperImpl) mapLogicalAndToAPIModel(input *LogicalAndExpression) restapi.MatchExpression {
	left := m.mapPrimaryExpressionToAPIModel(input.Left)
	if input.Operator != nil {
		right := m.mapLogicalAndToAPIModel(input.Right)
		return restapi.NewBinaryOperator(left, restapi.LogicalAnd, right)
	}
	return left
}

func (m *matchExpressionMapperImpl) mapPrimaryExpressionToAPIModel(input *PrimaryExpression) restapi.MatchExpression {
	if input.UnaryOperation != nil {
		return m.mapUnaryOperatorExpressionToAPIModel(input.UnaryOperation)
	}
	return m.mapComparisionExpressionToAPIModel(input.Comparision)
}

func (m *matchExpressionMapperImpl) mapUnaryOperatorExpressionToAPIModel(input *UnaryOperationExpression) restapi.MatchExpression {
	return restapi.NewUnaryOperationExpression(input.Entity.Key, input.Entity.Origin.MatcherExpressionEntity(), restapi.MatcherOperator(input.Operator))
}

func (m *matchExpressionMapperImpl) mapComparisionExpressionToAPIModel(input *ComparisionExpression) restapi.MatchExpression {
	return restapi.NewComparisionExpression(input.Entity.Key, input.Entity.Origin.MatcherExpressionEntity(), restapi.MatcherOperator(input.Operator), input.Value)
}
