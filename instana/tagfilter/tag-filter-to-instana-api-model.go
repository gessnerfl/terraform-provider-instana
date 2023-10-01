package tagfilter

import (
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

// ToAPIModel Implementation of the mapping form filter expression model to the Instana API model
func (m *tagFilterMapper) ToAPIModel(input *FilterExpression) *restapi.TagFilter {
	return m.mapLogicalOrToAPIModel(input.Expression)
}

func (m *tagFilterMapper) mapLogicalOrToAPIModel(input *LogicalOrExpression) *restapi.TagFilter {
	left := m.mapLogicalAndToAPIModel(input.Left)
	if input.Operator != nil {
		right := m.mapLogicalOrToAPIModel(input.Right)
		leftElements := m.unwrapExpressionElements(left, restapi.LogicalOr)
		rightElements := m.unwrapExpressionElements(right, restapi.LogicalOr)
		return restapi.NewLogicalOrTagFilter(append(leftElements, rightElements...))
	}
	return left
}

func (m *tagFilterMapper) mapLogicalAndToAPIModel(input *LogicalAndExpression) *restapi.TagFilter {
	left := m.mapBracketExpressionToAPIModel(input.Left)
	if input.Operator != nil {
		right := m.mapLogicalAndToAPIModel(input.Right)
		leftElements := m.unwrapExpressionElements(left, restapi.LogicalAnd)
		rightElements := m.unwrapExpressionElements(right, restapi.LogicalAnd)
		return restapi.NewLogicalAndTagFilter(append(leftElements, rightElements...))
	}
	return left
}

func (m *tagFilterMapper) unwrapExpressionElements(element *restapi.TagFilter, operator restapi.LogicalOperatorType) []*restapi.TagFilter {
	if element.GetType() == restapi.TagFilterExpressionType && element.LogicalOperator != nil && *element.LogicalOperator == operator {
		return element.Elements
	}
	return []*restapi.TagFilter{element}
}

func (m *tagFilterMapper) mapBracketExpressionToAPIModel(input *BracketExpression) *restapi.TagFilter {
	if input.Bracket != nil {
		return m.mapLogicalOrToAPIModel(input.Bracket)
	}
	return m.mapPrimaryExpressionToAPIModel(input.Primary)
}

func (m *tagFilterMapper) mapPrimaryExpressionToAPIModel(input *PrimaryExpression) *restapi.TagFilter {
	if input.UnaryOperation != nil {
		return m.mapUnaryOperatorExpressionToAPIModel(input.UnaryOperation)
	}
	return m.mapComparisonExpressionToAPIModel(input.Comparison)
}

func (m *tagFilterMapper) mapUnaryOperatorExpressionToAPIModel(input *UnaryOperationExpression) *restapi.TagFilter {
	origin := EntityOriginDestination.TagFilterEntity()
	if input.Entity.Origin != nil {
		origin = SupportedEntityOrigins.ForKey(*input.Entity.Origin).TagFilterEntity()
	}
	return restapi.NewUnaryTagFilterWithTagKey(origin, input.Entity.Identifier, input.Entity.TagKey, restapi.ExpressionOperator(input.Operator))
}

func (m *tagFilterMapper) mapComparisonExpressionToAPIModel(input *ComparisonExpression) *restapi.TagFilter {
	origin := EntityOriginDestination.TagFilterEntity()
	if input.Entity.Origin != nil {
		origin = SupportedEntityOrigins.ForKey(*input.Entity.Origin).TagFilterEntity()
	}
	if input.Entity.TagKey != nil {
		return restapi.NewTagTagFilter(origin, input.Entity.Identifier, restapi.ExpressionOperator(input.Operator), *input.Entity.TagKey, m.mapValueAsString(input))
	} else if input.NumberValue != nil {
		return restapi.NewNumberTagFilter(origin, input.Entity.Identifier, restapi.ExpressionOperator(input.Operator), *input.NumberValue)
	} else if input.BooleanValue != nil {
		return restapi.NewBooleanTagFilter(origin, input.Entity.Identifier, restapi.ExpressionOperator(input.Operator), *input.BooleanValue)
	}
	return restapi.NewStringTagFilter(origin, input.Entity.Identifier, restapi.ExpressionOperator(input.Operator), *input.StringValue)
}

func (m *tagFilterMapper) mapValueAsString(input *ComparisonExpression) string {
	if input.NumberValue != nil {
		return fmt.Sprintf("%d", *input.NumberValue)
	} else if input.BooleanValue != nil {
		return fmt.Sprintf("%t", *input.BooleanValue)
	}
	return *input.StringValue
}
