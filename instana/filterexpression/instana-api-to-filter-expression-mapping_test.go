package filterexpression_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	. "github.com/gessnerfl/terraform-provider-instana/instana/filterexpression"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

const (
	invalidOperator = "invalid operator"
	unaryOperator   = "unary operator"
	comparision     = "comparision"

	messageExpectedToGetInvalidLogicalAndError  = "Expected to get invalid logical AND error"
	messageExpectedToGetInvalidLogicalOrError   = "Expected to get invalid logical OR error"
	messageExpectedToGetInvalidConjunctionError = "Expected to get invalid conjunction error"
)

func TestShouldMapValidOperatorsOfTagExpression(t *testing.T) {
	for _, v := range restapi.SupportedComparisionOperators {
		t.Run(fmt.Sprintf("test mapping of %s", v), testMappingOfOperatorsOfTagExpression(v))
	}
}

func testMappingOfOperatorsOfTagExpression(operator restapi.MatcherOperator) func(t *testing.T) {
	return func(t *testing.T) {
		key := "key"
		value := "value"
		input := restapi.TagMatcherExpression{
			Dtype:    restapi.LeafExpressionType,
			Key:      key,
			Operator: operator,
			Value:    &value,
		}

		expectedResult := &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						Comparision: &ComparisionExpression{
							Key:      key,
							Operator: Operator(operator),
							Value:    value,
						},
					},
				},
			},
		}

		runTestCaseForMappingFromAPI(input, expectedResult, t)
	}
}

func TestShouldFailMapToMapComparisionWhenOperatorOfTagExpressionIsNotValid(t *testing.T) {
	key := "key"
	value := "value"
	input := restapi.TagMatcherExpression{
		Dtype:    restapi.LeafExpressionType,
		Key:      key,
		Operator: "FOO",
		Value:    &value,
	}

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	if err == nil || !strings.HasPrefix(err.Error(), invalidOperator) || !strings.Contains(err.Error(), comparision) {
		t.Fatal("Expected to get invalid operation error")
	}
}

func TestShouldMapValidUnaryOperationsOfTagExpression(t *testing.T) {
	for _, v := range restapi.SupportedUnaryOperatorExpressionOperators {
		t.Run(fmt.Sprintf("test mapping of %s ", v), testMappingOfUnaryOperationOfTagExpression(v))
	}
}

func testMappingOfUnaryOperationOfTagExpression(operator restapi.MatcherOperator) func(t *testing.T) {
	return func(t *testing.T) {
		key := "key"
		input := restapi.TagMatcherExpression{
			Dtype:    restapi.LeafExpressionType,
			Key:      key,
			Operator: operator,
		}

		expectedResult := &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Key:      key,
							Operator: Operator(operator),
						},
					},
				},
			},
		}

		runTestCaseForMappingFromAPI(input, expectedResult, t)
	}
}

func TestShouldFailMapToMapUnaryOperationWhenOperatorOfTagExpressionIsNotValid(t *testing.T) {
	key := "key"
	input := restapi.TagMatcherExpression{
		Dtype:    restapi.LeafExpressionType,
		Key:      key,
		Operator: "FOO",
	}

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	if err == nil || !strings.HasPrefix(err.Error(), invalidOperator) || !strings.Contains(err.Error(), unaryOperator) {
		t.Fatal("Expected to get invalid unary operation error")
	}
}

func TestShouldFailMapToMapExpressionWhenTypeIsMissing(t *testing.T) {
	key := "key"
	input := restapi.TagMatcherExpression{
		Key:      key,
		Operator: "FOO",
	}

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	if err == nil || !strings.HasPrefix(err.Error(), "unsupported match expression") {
		t.Fatal("Expected to get unsupported match expression error")
	}
}

func TestShouldMapLogicalAndWhenLeftAndRightIsAPrimaryExpression(t *testing.T) {
	key := "key"
	operator := Operator(restapi.IsEmptyOperator)
	and := Operator(restapi.LogicalAnd)
	primaryExpression := restapi.NewUnaryOperationExpression(key, restapi.IsEmptyOperator)
	input := restapi.NewBinaryOperator(primaryExpression, restapi.LogicalAnd, primaryExpression)

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					UnaryOperation: &UnaryOperationExpression{
						Key:      key,
						Operator: operator,
					},
				},
				Operator: &and,
				Right: &LogicalAndExpression{
					Left: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Key:      key,
							Operator: operator,
						},
					},
				},
			},
		},
	}

	runTestCaseForMappingFromAPI(input, expectedResult, t)
}

func TestShouldMapLogicalAndWhenLeftIsAPrimaryExpressionAndRightIsAnotherAndExpression(t *testing.T) {
	key := "key"
	operator := Operator(restapi.IsEmptyOperator)
	and := Operator(restapi.LogicalAnd)
	primaryExpression := restapi.NewUnaryOperationExpression(key, restapi.IsEmptyOperator)
	nestedAnd := restapi.NewBinaryOperator(primaryExpression, restapi.LogicalAnd, primaryExpression)
	input := restapi.NewBinaryOperator(primaryExpression, restapi.LogicalAnd, nestedAnd)

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					UnaryOperation: &UnaryOperationExpression{
						Key:      key,
						Operator: operator,
					},
				},
				Operator: &and,
				Right: &LogicalAndExpression{
					Left: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Key:      key,
							Operator: operator,
						},
					},
					Operator: &and,
					Right: &LogicalAndExpression{
						Left: &PrimaryExpression{
							UnaryOperation: &UnaryOperationExpression{
								Key:      key,
								Operator: operator,
							},
						},
					},
				},
			},
		},
	}

	runTestCaseForMappingFromAPI(input, expectedResult, t)
}

func TestShouldFailToMapLogicalAndWhenLeftIsOrExpression(t *testing.T) {
	key := "key"
	primaryExpression := restapi.NewUnaryOperationExpression(key, restapi.IsEmptyOperator)
	nestedOr := restapi.NewBinaryOperator(primaryExpression, restapi.LogicalOr, primaryExpression)
	input := restapi.NewBinaryOperator(nestedOr, restapi.LogicalAnd, primaryExpression)

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	if err == nil || !strings.Contains(err.Error(), "logical or is not allowed for left side") {
		t.Fatal(messageExpectedToGetInvalidLogicalAndError)
	}
}

func TestShouldFailToMapLogicalAndWhenRightIsOrExpression(t *testing.T) {
	key := "key"
	primaryExpression := restapi.NewUnaryOperationExpression(key, restapi.IsEmptyOperator)
	nestedOr := restapi.NewBinaryOperator(primaryExpression, restapi.LogicalOr, primaryExpression)
	input := restapi.NewBinaryOperator(primaryExpression, restapi.LogicalAnd, nestedOr)

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	if err == nil || !strings.Contains(err.Error(), "logical or is not allowed for right side") {
		t.Fatal(messageExpectedToGetInvalidLogicalAndError)
	}
}

func TestShouldFailToMapLogicalAndWhenLeftIsAndExpression(t *testing.T) {
	key := "key"
	primaryExpression := restapi.NewUnaryOperationExpression(key, restapi.IsEmptyOperator)
	nestedAnd := restapi.NewBinaryOperator(primaryExpression, restapi.LogicalAnd, primaryExpression)
	input := restapi.NewBinaryOperator(nestedAnd, restapi.LogicalAnd, primaryExpression)

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	if err == nil || !strings.Contains(err.Error(), "logical and is not allowed for left side") {
		t.Fatal(messageExpectedToGetInvalidLogicalAndError)
	}
}

func TestShouldMapLogiclOrWhenLeftAndRightSideIsPrimaryExpression(t *testing.T) {
	key := "key"
	operator := Operator(restapi.IsEmptyOperator)
	or := Operator(restapi.LogicalOr)
	primaryExpression := restapi.NewUnaryOperationExpression(key, restapi.IsEmptyOperator)
	input := restapi.NewBinaryOperator(primaryExpression, restapi.LogicalOr, primaryExpression)

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					UnaryOperation: &UnaryOperationExpression{
						Key:      key,
						Operator: operator,
					},
				},
			},
			Operator: &or,
			Right: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Key:      key,
							Operator: operator,
						},
					},
				},
			},
		},
	}

	runTestCaseForMappingFromAPI(input, expectedResult, t)
}

func TestShouldMapLogiclOrWhenLeftSideIsALogicalAndAndRightSideIsPrimaryExpression(t *testing.T) {
	key := "key"
	operator := Operator(restapi.IsEmptyOperator)
	or := Operator(restapi.LogicalOr)
	and := Operator(restapi.LogicalAnd)
	primaryExpression := restapi.NewUnaryOperationExpression(key, restapi.IsEmptyOperator)
	nestedAnd := restapi.NewBinaryOperator(primaryExpression, restapi.LogicalAnd, primaryExpression)
	input := restapi.NewBinaryOperator(nestedAnd, restapi.LogicalOr, primaryExpression)

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					UnaryOperation: &UnaryOperationExpression{
						Key:      key,
						Operator: operator,
					},
				},
				Operator: &and,
				Right: &LogicalAndExpression{
					Left: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Key:      key,
							Operator: operator,
						},
					},
				},
			},
			Operator: &or,
			Right: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Key:      key,
							Operator: operator,
						},
					},
				},
			},
		},
	}

	runTestCaseForMappingFromAPI(input, expectedResult, t)
}

func TestShouldMapLogiclOrWhenLeftSideIsAPrimaryExpressionAndRightSideIsALogicalOr(t *testing.T) {
	key := "key"
	operator := Operator(restapi.IsEmptyOperator)
	or := Operator(restapi.LogicalOr)
	primaryExpression := restapi.NewUnaryOperationExpression(key, restapi.IsEmptyOperator)
	nestedOr := restapi.NewBinaryOperator(primaryExpression, restapi.LogicalOr, primaryExpression)
	input := restapi.NewBinaryOperator(primaryExpression, restapi.LogicalOr, nestedOr)

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					UnaryOperation: &UnaryOperationExpression{
						Key:      key,
						Operator: operator,
					},
				},
			},
			Operator: &or,
			Right: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Key:      key,
							Operator: operator,
						},
					},
				},
				Operator: &or,
				Right: &LogicalOrExpression{
					Left: &LogicalAndExpression{
						Left: &PrimaryExpression{
							UnaryOperation: &UnaryOperationExpression{
								Key:      key,
								Operator: operator,
							},
						},
					},
				},
			},
		},
	}

	runTestCaseForMappingFromAPI(input, expectedResult, t)
}

func TestShouldMapLogiclOrWhenLeftSideIsAPrimaryExpressionAndRightSideIsALogicalAnd(t *testing.T) {
	key := "key"
	operator := Operator(restapi.IsEmptyOperator)
	or := Operator(restapi.LogicalOr)
	and := Operator(restapi.LogicalAnd)
	primaryExpression := restapi.NewUnaryOperationExpression(key, restapi.IsEmptyOperator)
	nestedAnd := restapi.NewBinaryOperator(primaryExpression, restapi.LogicalAnd, primaryExpression)
	input := restapi.NewBinaryOperator(primaryExpression, restapi.LogicalOr, nestedAnd)

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					UnaryOperation: &UnaryOperationExpression{
						Key:      key,
						Operator: operator,
					},
				},
			},
			Operator: &or,
			Right: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Key:      key,
							Operator: operator,
						},
					},
					Operator: &and,
					Right: &LogicalAndExpression{
						Left: &PrimaryExpression{
							UnaryOperation: &UnaryOperationExpression{
								Key:      key,
								Operator: operator,
							},
						},
					},
				},
			},
		},
	}

	runTestCaseForMappingFromAPI(input, expectedResult, t)
}

func TestShouldFailToMapLogicalOrWhenLeftIsOrExpression(t *testing.T) {
	key := "key"
	primaryExpression := restapi.NewUnaryOperationExpression(key, restapi.IsEmptyOperator)
	nestedOr := restapi.NewBinaryOperator(primaryExpression, restapi.LogicalOr, primaryExpression)
	input := restapi.NewBinaryOperator(nestedOr, restapi.LogicalOr, primaryExpression)

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	if err == nil || !strings.Contains(err.Error(), "logical or is not allowed for left side") {
		t.Fatal(messageExpectedToGetInvalidLogicalOrError)
	}
}

func TestShouldFailToMapBinaryExpressionWhenConjunctionTypeIsNotValid(t *testing.T) {
	key := "key"
	primaryExpression := restapi.NewUnaryOperationExpression(key, restapi.IsEmptyOperator)
	input := restapi.NewBinaryOperator(primaryExpression, "FOO", primaryExpression)

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	if err == nil || !strings.HasPrefix(err.Error(), "invalid conjunction operator") {
		t.Fatal("Expected to get invalid match expression error")
	}
}

func TestShouldReturnMappingErrorIfLeftSideOfConjunctionIsNotValid(t *testing.T) {
	key := "key"
	primaryExpressionLeft := restapi.NewUnaryOperationExpression(key, "INVALID")
	primaryExpressionRight := restapi.NewUnaryOperationExpression(key, restapi.IsEmptyOperator)
	input := restapi.NewBinaryOperator(primaryExpressionLeft, restapi.LogicalOr, primaryExpressionRight)

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	if err == nil || !strings.HasPrefix(err.Error(), invalidOperator) || !strings.Contains(err.Error(), unaryOperator) {
		t.Fatal(messageExpectedToGetInvalidConjunctionError)
	}
}

func TestShouldReturnMappingErrorIfRightSideOfConjunctionIsNotValid(t *testing.T) {
	key := "key"
	primaryExpressionLeft := restapi.NewUnaryOperationExpression(key, restapi.IsEmptyOperator)
	primaryExpressionRight := restapi.NewUnaryOperationExpression(key, "INVALID")
	input := restapi.NewBinaryOperator(primaryExpressionLeft, restapi.LogicalOr, primaryExpressionRight)

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	if err == nil || !strings.HasPrefix(err.Error(), invalidOperator) || !strings.Contains(err.Error(), unaryOperator) {
		t.Fatal(messageExpectedToGetInvalidConjunctionError)
	}
}

func runTestCaseForMappingFromAPI(input restapi.MatchExpression, expectedResult *FilterExpression, t *testing.T) {
	mapper := NewMapper()
	result, err := mapper.FromAPIModel(input)

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
	if !cmp.Equal(result, expectedResult) {
		t.Fatalf("Parse result does not match; diff %s", cmp.Diff(expectedResult, result))
	}
}
