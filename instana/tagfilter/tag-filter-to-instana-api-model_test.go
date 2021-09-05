package filterexpression_test

import (
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	. "github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/stretchr/testify/require"
)

const (
	entitySpecKey = "key"
)

func TestShouldMapComparisonToRepresentationOfInstanaAPI(t *testing.T) {
	for _, v := range restapi.SupportedComparisonOperators {
		t.Run(fmt.Sprintf("test comparison of string value using operatore %s", v), createTestShouldMapStringComparisonToRepresentationOfInstanaAPI(v))
		t.Run(fmt.Sprintf("test comparison of number value using operatore of %s", v), createTestShouldMapNumberComparisonToRepresentationOfInstanaAPI(v))
		t.Run(fmt.Sprintf("test comparison of boolean value using operatore of %s", v), createTestShouldMapBooleanComparisonToRepresentationOfInstanaAPI(v))
		t.Run(fmt.Sprintf("test comparison of tag using operatore of %s", v), createTestShouldMapTagComparisonToRepresentationOfInstanaAPI(v))
	}
}

func createTestShouldMapStringComparisonToRepresentationOfInstanaAPI(operator restapi.TagFilterOperator) func(*testing.T) {
	return func(t *testing.T) {
		expr := &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:      &EntitySpec{Identifier: entitySpecKey, Origin: EntityOriginDestination},
							Operator:    Operator(operator),
							StringValue: utils.StringPtr("value"),
						},
					},
				},
			},
		}

		expectedResult := restapi.NewStringTagFilter(restapi.TagFilterEntityDestination, entitySpecKey, operator, utils.StringPtr("value"))
		runTestCaseForMappingToAPI(expr, expectedResult, t)
	}
}

func createTestShouldMapNumberComparisonToRepresentationOfInstanaAPI(operator restapi.TagFilterOperator) func(*testing.T) {
	numberValue := int64(1234)
	return func(t *testing.T) {
		expr := &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:      &EntitySpec{Identifier: entitySpecKey, Origin: EntityOriginDestination},
							Operator:    Operator(operator),
							NumberValue: &numberValue,
						},
					},
				},
			},
		}

		expectedResult := restapi.NewNumberTagFilter(restapi.TagFilterEntityDestination, entitySpecKey, operator, &numberValue)
		runTestCaseForMappingToAPI(expr, expectedResult, t)
	}
}

func createTestShouldMapBooleanComparisonToRepresentationOfInstanaAPI(operator restapi.TagFilterOperator) func(*testing.T) {
	boolValue := true
	return func(t *testing.T) {
		expr := &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:       &EntitySpec{Identifier: entitySpecKey, Origin: EntityOriginDestination},
							Operator:     Operator(operator),
							BooleanValue: &boolValue,
						},
					},
				},
			},
		}

		expectedResult := restapi.NewBooleanTagFilter(restapi.TagFilterEntityDestination, entitySpecKey, operator, &boolValue)
		runTestCaseForMappingToAPI(expr, expectedResult, t)
	}
}

func createTestShouldMapTagComparisonToRepresentationOfInstanaAPI(operator restapi.TagFilterOperator) func(*testing.T) {
	key := "key"
	value := "value"
	return func(t *testing.T) {
		expr := &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:   &EntitySpec{Identifier: entitySpecKey, Origin: EntityOriginDestination},
							Operator: Operator(operator),
							TagValue: &TagValue{
								Key:   key,
								Value: value,
							},
						},
					},
				},
			},
		}

		expectedResult := restapi.NewTagTagFilter(restapi.TagFilterEntityDestination, entitySpecKey, operator, &key, &value)
		runTestCaseForMappingToAPI(expr, expectedResult, t)
	}
}

func TestShouldMapUnaryOperatorToRepresentationOfInstanaAPI(t *testing.T) {
	for _, v := range restapi.SupportedUnaryExpressionOperators {
		t.Run(fmt.Sprintf("test mapping of %s", v), createTestShouldMapUnaryOperatorToRepresentationOfInstanaAPI(v))
	}
}

func createTestShouldMapUnaryOperatorToRepresentationOfInstanaAPI(operatorName restapi.TagFilterOperator) func(*testing.T) {
	return func(t *testing.T) {
		expr := &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Entity:   &EntitySpec{Identifier: entitySpecKey, Origin: EntityOriginDestination},
							Operator: Operator(operatorName),
						},
					},
				},
			},
		}

		expectedResult := restapi.NewUnaryTagFilter(restapi.TagFilterEntityDestination, entitySpecKey, operatorName)
		runTestCaseForMappingToAPI(expr, expectedResult, t)
	}
}

func TestShouldMapLogicalAndExpression(t *testing.T) {
	logicalAnd := Operator(restapi.LogicalAnd)
	primaryExpression := PrimaryExpression{
		UnaryOperation: &UnaryOperationExpression{
			Entity:   &EntitySpec{Identifier: entitySpecKey, Origin: EntityOriginDestination},
			Operator: Operator(restapi.IsEmptyOperator),
		},
	}
	expr := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left:     &primaryExpression,
				Operator: &logicalAnd,
				Right: &LogicalAndExpression{
					Left: &primaryExpression,
				},
			},
		},
	}

	expectedPrimaryExpression := restapi.NewUnaryTagFilter(restapi.TagFilterEntityDestination, entitySpecKey, restapi.IsEmptyOperator)
	expectedResult := restapi.NewLogicalAndTagFilter([]restapi.TagFilterExpressionElement{expectedPrimaryExpression, expectedPrimaryExpression})
	runTestCaseForMappingToAPI(expr, expectedResult, t)
}

func TestShouldMapLogicalOrExpression(t *testing.T) {
	logicalOr := Operator(restapi.LogicalOr)
	primaryExpression := PrimaryExpression{
		UnaryOperation: &UnaryOperationExpression{
			Entity:   &EntitySpec{Identifier: entitySpecKey, Origin: EntityOriginDestination},
			Operator: Operator(restapi.IsEmptyOperator),
		},
	}
	expr := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &primaryExpression,
			},
			Operator: &logicalOr,
			Right: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &primaryExpression,
				},
			},
		},
	}

	expectedPrimaryExpression := restapi.NewUnaryTagFilter(restapi.TagFilterEntityDestination, entitySpecKey, restapi.IsEmptyOperator)
	expectedResult := restapi.NewLogicalOrTagFilter([]restapi.TagFilterExpressionElement{expectedPrimaryExpression, expectedPrimaryExpression})
	runTestCaseForMappingToAPI(expr, expectedResult, t)
}

func runTestCaseForMappingToAPI(input *FilterExpression, expectedResult restapi.TagFilterExpressionElement, t *testing.T) {
	mapper := NewTagFilterMapper()
	result := mapper.ToAPIModel(input)

	require.Equal(t, expectedResult, result)
}
