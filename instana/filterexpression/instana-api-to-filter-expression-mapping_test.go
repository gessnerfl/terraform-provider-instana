package filterexpression_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	. "github.com/gessnerfl/terraform-provider-instana/instana/filterexpression"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

func TestShouldMapValidOperatorsOfTagExpression(t *testing.T) {
	for _, v := range restapi.SupportedComparisionOperators {
		t.Run(fmt.Sprintf("test mapping of %s", v), testMappingOfOperatorsOfTagExpression(v))
	}
}

func testMappingOfOperatorsOfTagExpression(operatorName string) func(t *testing.T) {
	return func(t *testing.T) {
		key := "key"
		value := "value"
		input := restapi.TagMatcherExpression{
			Dtype:    restapi.LeafExpressionType,
			Key:      key,
			Operator: operatorName,
			Value:    &value,
		}

		expectedResult := &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						Comparision: &ComparisionExpression{
							Key:      key,
							Operator: Operator(operatorName),
							Value:    value,
						},
					},
				},
			},
		}

		runParsingTest(input, expectedResult, t)
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

	if err == nil || !strings.HasPrefix(err.Error(), "invalid operator") || !strings.Contains(err.Error(), "comparision") {
		t.Fatal("Expected to get invalid operation error")
	}
}

func TestShouldMapValidUnaryOperationsOfTagExpression(t *testing.T) {
	for _, v := range restapi.SupportedUnaryOperatorExpressionOperators {
		t.Run(fmt.Sprintf("test mapping of %s ", v), testMappingOfUnaryOperationOfTagExpression(v))
	}
}

func testMappingOfUnaryOperationOfTagExpression(operatorName string) func(t *testing.T) {
	return func(t *testing.T) {
		key := "key"
		input := restapi.TagMatcherExpression{
			Dtype:    restapi.LeafExpressionType,
			Key:      key,
			Operator: operatorName,
		}

		expectedResult := &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Key:      key,
							Operator: Operator(operatorName),
						},
					},
				},
			},
		}

		runParsingTest(input, expectedResult, t)
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

	if err == nil || !strings.HasPrefix(err.Error(), "invalid operator") || !strings.Contains(err.Error(), "unary operator") {
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
	operator := Operator("IS_EMPTY")
	and := Operator("AND")
	primaryExpression := restapi.NewUnaryOperationExpression(key, "IS_EMPTY")
	input := restapi.NewBinaryOperator(primaryExpression, "AND", primaryExpression)

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

	runParsingTest(input, expectedResult, t)
}

func TestShouldMapLogicalAndWhenLeftIsAPrimaryExpressionAndRightIsAnotherAndExpression(t *testing.T) {
	key := "key"
	operator := Operator("IS_EMPTY")
	and := Operator("AND")
	primaryExpression := restapi.NewUnaryOperationExpression(key, "IS_EMPTY")
	nestedAnd := restapi.NewBinaryOperator(primaryExpression, "AND", primaryExpression)
	input := restapi.NewBinaryOperator(primaryExpression, "AND", nestedAnd)

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

	runParsingTest(input, expectedResult, t)
}

func TestShouldFailToMapLogicalAndWhenLeftIsOrExpression(t *testing.T) {
	key := "key"
	primaryExpression := restapi.NewUnaryOperationExpression(key, "IS_EMPTY")
	nestedOr := restapi.NewBinaryOperator(primaryExpression, "OR", primaryExpression)
	input := restapi.NewBinaryOperator(nestedOr, "AND", primaryExpression)

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	if err == nil || !strings.Contains(err.Error(), "logical or is not allowed for left side") {
		t.Fatal("Expected to get invalid logical AND error")
	}
}

func TestShouldFailToMapLogicalAndWhenRightIsOrExpression(t *testing.T) {
	key := "key"
	primaryExpression := restapi.NewUnaryOperationExpression(key, "IS_EMPTY")
	nestedOr := restapi.NewBinaryOperator(primaryExpression, "OR", primaryExpression)
	input := restapi.NewBinaryOperator(primaryExpression, "AND", nestedOr)

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	if err == nil || !strings.Contains(err.Error(), "logical or is not allowed for right side") {
		t.Fatal("Expected to get invalid logical AND error")
	}
}

func TestShouldFailToMapLogicalAndWhenLeftIsAndExpression(t *testing.T) {
	key := "key"
	primaryExpression := restapi.NewUnaryOperationExpression(key, "IS_EMPTY")
	nestedAnd := restapi.NewBinaryOperator(primaryExpression, "AND", primaryExpression)
	input := restapi.NewBinaryOperator(nestedAnd, "AND", primaryExpression)

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	if err == nil || !strings.Contains(err.Error(), "logical and is not allowed for left side") {
		t.Fatal("Expected to get invalid logical AND error")
	}
}

func TestShouldMapLogiclOrWhenLeftAndRightSideIsPrimaryExpression(t *testing.T) {
	key := "key"
	operator := Operator("IS_EMPTY")
	or := Operator("OR")
	primaryExpression := restapi.NewUnaryOperationExpression(key, "IS_EMPTY")
	input := restapi.NewBinaryOperator(primaryExpression, "OR", primaryExpression)

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

	runParsingTest(input, expectedResult, t)
}

func TestShouldMapLogiclOrWhenLeftSideIsALogicalAndAndRightSideIsPrimaryExpression(t *testing.T) {
	key := "key"
	operator := Operator("IS_EMPTY")
	or := Operator("OR")
	and := Operator("AND")
	primaryExpression := restapi.NewUnaryOperationExpression(key, "IS_EMPTY")
	nestedAnd := restapi.NewBinaryOperator(primaryExpression, "AND", primaryExpression)
	input := restapi.NewBinaryOperator(nestedAnd, "OR", primaryExpression)

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

	runParsingTest(input, expectedResult, t)
}

func TestShouldMapLogiclOrWhenLeftSideIsAPrimaryExpressionAndRightSideIsALogicalOr(t *testing.T) {
	key := "key"
	operator := Operator("IS_EMPTY")
	or := Operator("OR")
	primaryExpression := restapi.NewUnaryOperationExpression(key, "IS_EMPTY")
	nestedOr := restapi.NewBinaryOperator(primaryExpression, "OR", primaryExpression)
	input := restapi.NewBinaryOperator(primaryExpression, "OR", nestedOr)

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

	runParsingTest(input, expectedResult, t)
}

func TestShouldMapLogiclOrWhenLeftSideIsAPrimaryExpressionAndRightSideIsALogicalAnd(t *testing.T) {
	key := "key"
	operator := Operator("IS_EMPTY")
	or := Operator("OR")
	and := Operator("AND")
	primaryExpression := restapi.NewUnaryOperationExpression(key, "IS_EMPTY")
	nestedAnd := restapi.NewBinaryOperator(primaryExpression, "AND", primaryExpression)
	input := restapi.NewBinaryOperator(primaryExpression, "OR", nestedAnd)

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

	runParsingTest(input, expectedResult, t)
}

func TestShouldFailToMapLogicalOrWhenLeftIsOrExpression(t *testing.T) {
	key := "key"
	primaryExpression := restapi.NewUnaryOperationExpression(key, "IS_EMPTY")
	nestedOr := restapi.NewBinaryOperator(primaryExpression, "OR", primaryExpression)
	input := restapi.NewBinaryOperator(nestedOr, "OR", primaryExpression)

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	if err == nil || !strings.Contains(err.Error(), "logical or is not allowed for left side") {
		t.Fatal("Expected to get invalid logical AND error")
	}
}

func TestShouldFailToMapBinaryExpressionWhenConjunctionTypeIsNotValid(t *testing.T) {
	key := "key"
	primaryExpression := restapi.NewUnaryOperationExpression(key, "IS_EMPTY")
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
	primaryExpressionRight := restapi.NewUnaryOperationExpression(key, "IS_EMPTY")
	input := restapi.NewBinaryOperator(primaryExpressionLeft, "OR", primaryExpressionRight)

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	if err == nil || !strings.HasPrefix(err.Error(), "invalid operator") || !strings.Contains(err.Error(), "unary operator") {
		t.Fatal("Expected to get invalid logical AND error")
	}
}

func TestShouldReturnMappingErrorIfRightSideOfConjunctionIsNotValid(t *testing.T) {
	key := "key"
	primaryExpressionLeft := restapi.NewUnaryOperationExpression(key, "IS_EMPTY")
	primaryExpressionRight := restapi.NewUnaryOperationExpression(key, "INVALID")
	input := restapi.NewBinaryOperator(primaryExpressionLeft, "OR", primaryExpressionRight)

	mapper := NewMapper()
	_, err := mapper.FromAPIModel(input)

	if err == nil || !strings.HasPrefix(err.Error(), "invalid operator") || !strings.Contains(err.Error(), "unary operator") {
		t.Fatal("Expected to get invalid logical AND error")
	}
}

func runParsingTest(input restapi.MatchExpression, expectedResult *FilterExpression, t *testing.T) {
	mapper := NewMapper()
	result, err := mapper.FromAPIModel(input)

	if err != nil {
		t.Fatalf("Expected no error but got %s", err)
	}
	if !cmp.Equal(result, expectedResult) {
		t.Fatalf("Parse result does not match; diff %s", cmp.Diff(expectedResult, result))
	}
}
