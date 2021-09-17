package tagfilter

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

//ToAPIModel Implementation of the mapping form filter expression model to the Instana API model
func (m *tagFilterMapper) ToAPIModel(input *FilterExpression) restapi.TagFilterExpressionElement {
	return m.mapLogicalOrToAPIModel(input.Expression)
}

func (m *tagFilterMapper) mapLogicalOrToAPIModel(input *LogicalOrExpression) restapi.TagFilterExpressionElement {
	left := m.mapLogicalAndToAPIModel(input.Left)
	if input.Operator != nil {
		right := m.mapLogicalOrToAPIModel(input.Right)
		leftElements := m.unwrapExpressionElements(left, restapi.LogicalOr)
		rightElements := m.unwrapExpressionElements(right, restapi.LogicalOr)
		return restapi.NewLogicalOrTagFilter(append(leftElements, rightElements...))
	}
	return left
}

func (m *tagFilterMapper) mapLogicalAndToAPIModel(input *LogicalAndExpression) restapi.TagFilterExpressionElement {
	left := m.mapPrimaryExpressionToAPIModel(input.Left)
	if input.Operator != nil {
		right := m.mapLogicalAndToAPIModel(input.Right)
		leftElements := m.unwrapExpressionElements(left, restapi.LogicalAnd)
		rightElements := m.unwrapExpressionElements(right, restapi.LogicalAnd)
		return restapi.NewLogicalAndTagFilter(append(leftElements, rightElements...))
	}
	return left
}

func (m *tagFilterMapper) unwrapExpressionElements(element restapi.TagFilterExpressionElement, operator restapi.LogicalOperatorType) []restapi.TagFilterExpressionElement {
	if element.GetType() == restapi.TagFilterExpressionType && element.(*restapi.TagFilterExpression).LogicalOperator == operator {
		return element.(*restapi.TagFilterExpression).Elements
	}
	return []restapi.TagFilterExpressionElement{element}
}

func (m *tagFilterMapper) mapPrimaryExpressionToAPIModel(input *PrimaryExpression) restapi.TagFilterExpressionElement {
	if input.UnaryOperation != nil {
		return m.mapUnaryOperatorExpressionToAPIModel(input.UnaryOperation)
	}
	return m.mapComparisonExpressionToAPIModel(input.Comparison)
}

func (m *tagFilterMapper) mapUnaryOperatorExpressionToAPIModel(input *UnaryOperationExpression) restapi.TagFilterExpressionElement {
	return restapi.NewUnaryTagFilter(input.Entity.Origin.TagFilterEntity(), input.Entity.Identifier, restapi.TagFilterOperator(input.Operator))
}

func (m *tagFilterMapper) mapComparisonExpressionToAPIModel(input *ComparisonExpression) restapi.TagFilterExpressionElement {
	if input.TagValue != nil {
		return restapi.NewTagTagFilter(input.Entity.Origin.TagFilterEntity(), input.Entity.Identifier, restapi.TagFilterOperator(input.Operator), &input.TagValue.Key, &input.TagValue.Value)
	} else if input.NumberValue != nil {
		return restapi.NewNumberTagFilter(input.Entity.Origin.TagFilterEntity(), input.Entity.Identifier, restapi.TagFilterOperator(input.Operator), input.NumberValue)
	} else if input.BooleanValue != nil {
		return restapi.NewBooleanTagFilter(input.Entity.Origin.TagFilterEntity(), input.Entity.Identifier, restapi.TagFilterOperator(input.Operator), input.BooleanValue)
	}
	return restapi.NewStringTagFilter(input.Entity.Origin.TagFilterEntity(), input.Entity.Identifier, restapi.TagFilterOperator(input.Operator), input.StringValue)
}
