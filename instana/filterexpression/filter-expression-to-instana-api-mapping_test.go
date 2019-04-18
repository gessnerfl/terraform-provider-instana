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

func createTestShouldMapComparisionToRepresentationOfInstanaAPI(operatorName string) func(*testing.T) {
	return func(t *testing.T) {
		expr := &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						Comparision: &ComparisionExpression{
							Key:      "key",
							Operator: Operator(operatorName),
							Value:    "value",
						},
					},
				},
			},
		}

		expectedResult := restapi.NewComparisionExpression("key", operatorName, "value")
		runTestCaseForMappingToAPI(expr, expectedResult, t)
	}
}

func TestShouldMapUnaryOperatorToRepresentationOfInstanaAPI(t *testing.T) {
	for _, v := range restapi.SupportedUnaryOperatorExpressionOperators {
		t.Run(fmt.Sprintf("test mapping of %s", v), createTestShouldMapUnaryOperatorToRepresentationOfInstanaAPI(v))
	}
}

func createTestShouldMapUnaryOperatorToRepresentationOfInstanaAPI(operatorName string) func(*testing.T) {
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
	logicalAnd := Operator("AND")
	primaryExpression := PrimaryExpression{
		UnaryOperation: &UnaryOperationExpression{
			Key:      "key",
			Operator: Operator("IS_EMPTY"),
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

	expectedPrimaryExpression := restapi.NewUnaryOperationExpression("key", "IS_EMPTY")
	expectedResult := restapi.NewBinaryOperator(expectedPrimaryExpression, "AND", expectedPrimaryExpression)
	runTestCaseForMappingToAPI(expr, expectedResult, t)
}

func TestShouldMapLogicalOrExpression(t *testing.T) {
	logicalOr := Operator("OR")
	primaryExpression := PrimaryExpression{
		UnaryOperation: &UnaryOperationExpression{
			Key:      "key",
			Operator: Operator("IS_EMPTY"),
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

	expectedPrimaryExpression := restapi.NewUnaryOperationExpression("key", "IS_EMPTY")
	expectedResult := restapi.NewBinaryOperator(expectedPrimaryExpression, "OR", expectedPrimaryExpression)
	runTestCaseForMappingToAPI(expr, expectedResult, t)
}

func TestShouldMapSubExpression(t *testing.T) {
	logicalOr := Operator("OR")
	logicalAnd := Operator("AND")
	expr := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					UnaryOperation: &UnaryOperationExpression{
						Key:      "foo",
						Operator: "IS_EMPTY",
					},
				},
				Operator: &logicalAnd,
				Right: &LogicalAndExpression{
					Left: &PrimaryExpression{
						SubExpression: &LogicalOrExpression{
							Left: &LogicalAndExpression{
								Left: &PrimaryExpression{
									Comparision: &ComparisionExpression{
										Key:      "a",
										Operator: "EQUALS",
										Value:    "b",
									},
								},
							},
							Operator: &logicalOr,
							Right: &LogicalOrExpression{
								Left: &LogicalAndExpression{
									Left: &PrimaryExpression{
										Comparision: &ComparisionExpression{
											Key:      "a",
											Operator: "EQUALS",
											Value:    "c",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	primaryExpression1 := restapi.NewUnaryOperationExpression("foo", "IS_EMPTY")
	primaryExpression2 := restapi.NewComparisionExpression("a", "EQUALS", "b")
	primaryExpression3 := restapi.NewComparisionExpression("a", "EQUALS", "c")
	nestedOr := restapi.NewBinaryOperator(primaryExpression2, "OR", primaryExpression3)
	expectedResult := restapi.NewBinaryOperator(primaryExpression1, "AND", nestedOr)
	runTestCaseForMappingToAPI(expr, expectedResult, t)
}

func runTestCaseForMappingToAPI(input *FilterExpression, expectedResult restapi.MatchExpression, t *testing.T) {
	mapper := NewMapper()
	result := mapper.ToAPIModel(input)

	if !cmp.Equal(result, expectedResult) {
		t.Fatalf("Parse result does not match; diff %s", cmp.Diff(expectedResult, result))
	}
}
