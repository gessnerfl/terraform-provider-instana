package filterexpression_test

import (
	"fmt"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/filterexpression"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
)

const (
	entitySpecKey = "key"
)

func TestShouldMapSimpleComparisionToRepresentationOfInstanaAPI(t *testing.T) {
	for _, v := range restapi.SupportedComparisonOperators {
		t.Run(fmt.Sprintf("test mapping of %s", v), createTestShouldMapComparisionToRepresentationOfInstanaAPI(v))
	}
}

func createTestShouldMapComparisionToRepresentationOfInstanaAPI(operator restapi.ExpressionOperator) func(*testing.T) {
	return func(t *testing.T) {
		expr := &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:   &EntitySpec{Identifier: entitySpecKey, Origin: EntityOriginDestination},
							Operator: Operator(operator),
							Value:    "value",
						},
					},
				},
			},
		}

		expectedResult := restapi.NewComparisonExpression(entitySpecKey, restapi.MatcherExpressionEntityDestination, operator, "value")
		runTestCaseForMappingToAPI(expr, expectedResult, t)
	}
}

func TestShouldMapUnaryOperatorToRepresentationOfInstanaAPI(t *testing.T) {
	for _, v := range restapi.SupportedUnaryExpressionOperators {
		t.Run(fmt.Sprintf("test mapping of %s", v), createTestShouldMapUnaryOperatorToRepresentationOfInstanaAPI(v))
	}
}

func createTestShouldMapUnaryOperatorToRepresentationOfInstanaAPI(operatorName restapi.ExpressionOperator) func(*testing.T) {
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

		expectedResult := restapi.NewUnaryOperationExpression(entitySpecKey, restapi.MatcherExpressionEntityDestination, operatorName)
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

	expectedPrimaryExpression := restapi.NewUnaryOperationExpression(entitySpecKey, restapi.MatcherExpressionEntityDestination, restapi.IsEmptyOperator)
	expectedResult := restapi.NewBinaryOperator(expectedPrimaryExpression, restapi.LogicalAnd, expectedPrimaryExpression)
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

	expectedPrimaryExpression := restapi.NewUnaryOperationExpression(entitySpecKey, restapi.MatcherExpressionEntityDestination, restapi.IsEmptyOperator)
	expectedResult := restapi.NewBinaryOperator(expectedPrimaryExpression, restapi.LogicalOr, expectedPrimaryExpression)
	runTestCaseForMappingToAPI(expr, expectedResult, t)
}

func runTestCaseForMappingToAPI(input *FilterExpression, expectedResult restapi.MatchExpression, t *testing.T) {
	mapper := NewMatchExpressionMapper()
	result := mapper.ToAPIModel(input)

	require.Equal(t, expectedResult, result)
}
