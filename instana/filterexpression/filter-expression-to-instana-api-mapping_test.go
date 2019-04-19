package filterexpression_test

import (
	"fmt"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/filterexpression"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/google/go-cmp/cmp"
)

func TestShouldMapSimpleComparisionToRepresentationOfInstanaAPI(t *testing.T) {
	for _, v := range restapi.SupportedComparisionOperators {
		t.Run(fmt.Sprintf("test mapping of %s", v), createTestShouldMapComparisionToRepresentationOfInstanaAPI(v))
	}
}

func createTestShouldMapComparisionToRepresentationOfInstanaAPI(operator restapi.MatcherOperator) func(*testing.T) {
	return func(t *testing.T) {
		expr := &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						Comparision: &ComparisionExpression{
							Key:      "key",
							Operator: Operator(operator),
							Value:    "value",
						},
					},
				},
			},
		}

		expectedResult := restapi.NewComparisionExpression("key", operator, "value")
		runTestCaseForMappingToAPI(expr, expectedResult, t)
	}
}

func TestShouldMapUnaryOperatorToRepresentationOfInstanaAPI(t *testing.T) {
	for _, v := range restapi.SupportedUnaryOperatorExpressionOperators {
		t.Run(fmt.Sprintf("test mapping of %s", v), createTestShouldMapUnaryOperatorToRepresentationOfInstanaAPI(v))
	}
}

func createTestShouldMapUnaryOperatorToRepresentationOfInstanaAPI(operatorName restapi.MatcherOperator) func(*testing.T) {
	return func(t *testing.T) {
		expr := &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Key:      "key",
							Operator: Operator(operatorName),
						},
					},
				},
			},
		}

		expectedResult := restapi.NewUnaryOperationExpression("key", operatorName)
		runTestCaseForMappingToAPI(expr, expectedResult, t)
	}
}

func TestShouldMapLogicalAndExpression(t *testing.T) {
	logicalAnd := Operator(restapi.LogicalAnd)
	primaryExpression := PrimaryExpression{
		UnaryOperation: &UnaryOperationExpression{
			Key:      "key",
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

	expectedPrimaryExpression := restapi.NewUnaryOperationExpression("key", restapi.IsEmptyOperator)
	expectedResult := restapi.NewBinaryOperator(expectedPrimaryExpression, restapi.LogicalAnd, expectedPrimaryExpression)
	runTestCaseForMappingToAPI(expr, expectedResult, t)
}

func TestShouldMapLogicalOrExpression(t *testing.T) {
	logicalOr := Operator(restapi.LogicalOr)
	primaryExpression := PrimaryExpression{
		UnaryOperation: &UnaryOperationExpression{
			Key:      "key",
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

	expectedPrimaryExpression := restapi.NewUnaryOperationExpression("key", restapi.IsEmptyOperator)
	expectedResult := restapi.NewBinaryOperator(expectedPrimaryExpression, restapi.LogicalOr, expectedPrimaryExpression)
	runTestCaseForMappingToAPI(expr, expectedResult, t)
}

func runTestCaseForMappingToAPI(input *FilterExpression, expectedResult restapi.MatchExpression, t *testing.T) {
	mapper := NewMapper()
	result := mapper.ToAPIModel(input)

	if !cmp.Equal(result, expectedResult) {
		t.Fatalf("Parse result does not match; diff %s", cmp.Diff(expectedResult, result))
	}
}
