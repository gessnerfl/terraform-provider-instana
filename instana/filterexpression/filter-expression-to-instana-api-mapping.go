package filterexpression

import "github.com/gessnerfl/terraform-provider-instana/instana/restapi"

//ToAPIModel Implementation of the mapping form filter expression model to the Instana API model
func (m *mapperImpl) ToAPIModel(input *FilterExpression) restapi.MatchExpression {
	return m.mapLogicalOrToAPIModel(input.Expression)
}

func (m *mapperImpl) mapLogicalOrToAPIModel(input *LogicalOrExpression) restapi.MatchExpression {
	left := m.mapLogicalAndToAPIModel(input.Left)
	if input.Operator != nil {
		right := m.mapLogicalOrToAPIModel(input.Right)
		return restapi.NewBinaryOperator(left, "OR", right)
	}
	return left
}

func (m *mapperImpl) mapLogicalAndToAPIModel(input *LogicalAndExpression) restapi.MatchExpression {
	left := m.mapPrimaryExpressionToAPIModel(input.Left)
	if input.Operator != nil {
		right := m.mapLogicalAndToAPIModel(input.Right)
		return restapi.NewBinaryOperator(left, "AND", right)
	}
	return left
}

func (m *mapperImpl) mapPrimaryExpressionToAPIModel(input *PrimaryExpression) restapi.MatchExpression {
	if input.SubExpression != nil {
		return m.mapSubExpressionToAPIModel(input.SubExpression)
	}
	if input.UnaryOperation != nil {
		return m.mapUnaryOperatorExpressionToAPIModel(input.UnaryOperation)
	}
	return m.mapComparisionExpressionToAPIModel(input.Comparision)
}

func (m *mapperImpl) mapUnaryOperatorExpressionToAPIModel(input *UnaryOperationExpression) restapi.MatchExpression {
	return restapi.NewUnaryOperationExpression(input.Key, string(input.Operator))
}

func (m *mapperImpl) mapComparisionExpressionToAPIModel(input *ComparisionExpression) restapi.MatchExpression {
	return restapi.NewComparisionExpression(input.Key, string(input.Operator), input.Value)
}

func (m *mapperImpl) mapSubExpressionToAPIModel(input *LogicalOrExpression) restapi.MatchExpression {
	return m.mapLogicalOrToAPIModel(input)
}
