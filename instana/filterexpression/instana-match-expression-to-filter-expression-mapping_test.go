package filterexpression_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana/filterexpression"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

const (
	invalidOperator     = "invalid operator"
	invalidEntityOrigin = "Invalid entity origin"
	unaryOperator       = "unary operator"
	comparision         = "comparision"

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
		input := restapi.NewComparisionExpression(key, restapi.MatcherExpressionEntityDestination, operator, value)

		expectedResult := &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						Comparision: &ComparisionExpression{
							Entity:   &EntitySpec{Key: key, Origin: EntityOriginDestination},
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

func TestShouldFailToMapComparisionWhenOperatorOfTagExpressionIsNotValid(t *testing.T) {
	key := "key"
	value := "value"
	input := restapi.NewComparisionExpression(key, restapi.MatcherExpressionEntityDestination, "FOO", value)

	mapper := NewMatchExpressionMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), invalidOperator)
	require.Contains(t, err.Error(), comparision)
}

func TestShouldMapValidUnaryOperationsOfTagExpression(t *testing.T) {
	for _, v := range restapi.SupportedUnaryExpressionOperators {
		t.Run(fmt.Sprintf("test mapping of %s ", v), testMappingOfUnaryOperationOfTagExpression(v))
	}
}

func testMappingOfUnaryOperationOfTagExpression(operator restapi.MatcherOperator) func(t *testing.T) {
	return func(t *testing.T) {
		key := "key"
		input := restapi.NewUnaryOperationExpression(key, restapi.MatcherExpressionEntityDestination, operator)

		expectedResult := &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Entity:   &EntitySpec{Key: key, Origin: EntityOriginDestination},
							Operator: Operator(operator),
						},
					},
				},
			},
		}

		runTestCaseForMappingFromAPI(input, expectedResult, t)
	}
}

func TestShouldFailToMapUnaryOperationWhenOperatorOfTagExpressionIsNotValid(t *testing.T) {
	key := "key"
	input := restapi.NewUnaryOperationExpression(key, restapi.MatcherExpressionEntityDestination, "FOO")

	mapper := NewMatchExpressionMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), invalidOperator)
	require.Contains(t, err.Error(), unaryOperator)
}

func TestShouldFailMapToMapExpressionWhenTypeIsMissing(t *testing.T) {
	key := "key"
	input := restapi.TagMatcherExpression{
		Key:      key,
		Operator: "FOO",
	}

	mapper := NewMatchExpressionMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "unsupported match expression")
}

func TestShouldMapLogicalAndWhenLeftAndRightIsAPrimaryExpression(t *testing.T) {
	key := "key"
	operator := Operator(restapi.IsEmptyOperator)
	and := Operator(restapi.LogicalAnd)
	primaryExpression := restapi.NewUnaryOperationExpression(key, restapi.MatcherExpressionEntityDestination, restapi.IsEmptyOperator)
	input := restapi.NewBinaryOperator(primaryExpression, restapi.LogicalAnd, primaryExpression)

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					UnaryOperation: &UnaryOperationExpression{
						Entity:   &EntitySpec{Key: key, Origin: EntityOriginDestination},
						Operator: operator,
					},
				},
				Operator: &and,
				Right: &LogicalAndExpression{
					Left: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Entity:   &EntitySpec{Key: key, Origin: EntityOriginDestination},
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
	primaryExpression := restapi.NewUnaryOperationExpression(key, restapi.MatcherExpressionEntityDestination, restapi.IsEmptyOperator)
	nestedAnd := restapi.NewBinaryOperator(primaryExpression, restapi.LogicalAnd, primaryExpression)
	input := restapi.NewBinaryOperator(primaryExpression, restapi.LogicalAnd, nestedAnd)

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					UnaryOperation: &UnaryOperationExpression{
						Entity:   &EntitySpec{Key: key, Origin: EntityOriginDestination},
						Operator: operator,
					},
				},
				Operator: &and,
				Right: &LogicalAndExpression{
					Left: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Entity:   &EntitySpec{Key: key, Origin: EntityOriginDestination},
							Operator: operator,
						},
					},
					Operator: &and,
					Right: &LogicalAndExpression{
						Left: &PrimaryExpression{
							UnaryOperation: &UnaryOperationExpression{
								Entity:   &EntitySpec{Key: key, Origin: EntityOriginDestination},
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
	primaryExpression := restapi.NewUnaryOperationExpression(key, restapi.MatcherExpressionEntityDestination, restapi.IsEmptyOperator)
	nestedOr := restapi.NewBinaryOperator(primaryExpression, restapi.LogicalOr, primaryExpression)
	input := restapi.NewBinaryOperator(nestedOr, restapi.LogicalAnd, primaryExpression)

	mapper := NewMatchExpressionMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "logical or is not allowed for left side")
}

func TestShouldFailToMapLogicalAndWhenRightIsOrExpression(t *testing.T) {
	key := "key"
	primaryExpression := restapi.NewUnaryOperationExpression(key, restapi.MatcherExpressionEntityDestination, restapi.IsEmptyOperator)
	nestedOr := restapi.NewBinaryOperator(primaryExpression, restapi.LogicalOr, primaryExpression)
	input := restapi.NewBinaryOperator(primaryExpression, restapi.LogicalAnd, nestedOr)

	mapper := NewMatchExpressionMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "logical or is not allowed for right side")
}

func TestShouldFailToMapLogicalAndWhenLeftIsAndExpression(t *testing.T) {
	key := "key"
	primaryExpression := restapi.NewUnaryOperationExpression(key, restapi.MatcherExpressionEntityDestination, restapi.IsEmptyOperator)
	nestedAnd := restapi.NewBinaryOperator(primaryExpression, restapi.LogicalAnd, primaryExpression)
	input := restapi.NewBinaryOperator(nestedAnd, restapi.LogicalAnd, primaryExpression)

	mapper := NewMatchExpressionMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "logical and is not allowed for left side")
}

func TestShouldMapLogiclOrWhenLeftAndRightSideIsPrimaryExpression(t *testing.T) {
	key := "key"
	operator := Operator(restapi.IsEmptyOperator)
	or := Operator(restapi.LogicalOr)
	primaryExpression := restapi.NewUnaryOperationExpression(key, restapi.MatcherExpressionEntityDestination, restapi.IsEmptyOperator)
	input := restapi.NewBinaryOperator(primaryExpression, restapi.LogicalOr, primaryExpression)

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					UnaryOperation: &UnaryOperationExpression{
						Entity:   &EntitySpec{Key: key, Origin: EntityOriginDestination},
						Operator: operator,
					},
				},
			},
			Operator: &or,
			Right: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Entity:   &EntitySpec{Key: key, Origin: EntityOriginDestination},
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
	primaryExpression := restapi.NewUnaryOperationExpression(key, restapi.MatcherExpressionEntityDestination, restapi.IsEmptyOperator)
	nestedAnd := restapi.NewBinaryOperator(primaryExpression, restapi.LogicalAnd, primaryExpression)
	input := restapi.NewBinaryOperator(nestedAnd, restapi.LogicalOr, primaryExpression)

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					UnaryOperation: &UnaryOperationExpression{
						Entity:   &EntitySpec{Key: key, Origin: EntityOriginDestination},
						Operator: operator,
					},
				},
				Operator: &and,
				Right: &LogicalAndExpression{
					Left: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Entity:   &EntitySpec{Key: key, Origin: EntityOriginDestination},
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
							Entity:   &EntitySpec{Key: key, Origin: EntityOriginDestination},
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
	primaryExpression := restapi.NewUnaryOperationExpression(key, restapi.MatcherExpressionEntityDestination, restapi.IsEmptyOperator)
	nestedOr := restapi.NewBinaryOperator(primaryExpression, restapi.LogicalOr, primaryExpression)
	input := restapi.NewBinaryOperator(primaryExpression, restapi.LogicalOr, nestedOr)

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					UnaryOperation: &UnaryOperationExpression{
						Entity:   &EntitySpec{Key: key, Origin: EntityOriginDestination},
						Operator: operator,
					},
				},
			},
			Operator: &or,
			Right: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Entity:   &EntitySpec{Key: key, Origin: EntityOriginDestination},
							Operator: operator,
						},
					},
				},
				Operator: &or,
				Right: &LogicalOrExpression{
					Left: &LogicalAndExpression{
						Left: &PrimaryExpression{
							UnaryOperation: &UnaryOperationExpression{
								Entity:   &EntitySpec{Key: key, Origin: EntityOriginDestination},
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
	primaryExpression := restapi.NewUnaryOperationExpression(key, restapi.MatcherExpressionEntityDestination, restapi.IsEmptyOperator)
	nestedAnd := restapi.NewBinaryOperator(primaryExpression, restapi.LogicalAnd, primaryExpression)
	input := restapi.NewBinaryOperator(primaryExpression, restapi.LogicalOr, nestedAnd)

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					UnaryOperation: &UnaryOperationExpression{
						Entity:   &EntitySpec{Key: key, Origin: EntityOriginDestination},
						Operator: operator,
					},
				},
			},
			Operator: &or,
			Right: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Entity:   &EntitySpec{Key: key, Origin: EntityOriginDestination},
							Operator: operator,
						},
					},
					Operator: &and,
					Right: &LogicalAndExpression{
						Left: &PrimaryExpression{
							UnaryOperation: &UnaryOperationExpression{
								Entity:   &EntitySpec{Key: key, Origin: EntityOriginDestination},
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
	primaryExpression := restapi.NewUnaryOperationExpression(key, restapi.MatcherExpressionEntityDestination, restapi.IsEmptyOperator)
	nestedOr := restapi.NewBinaryOperator(primaryExpression, restapi.LogicalOr, primaryExpression)
	input := restapi.NewBinaryOperator(nestedOr, restapi.LogicalOr, primaryExpression)

	mapper := NewMatchExpressionMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "logical or is not allowed for left side")
}

func TestShouldFailToMapBinaryExpressionWhenConjunctionTypeIsNotValid(t *testing.T) {
	key := "key"
	primaryExpression := restapi.NewUnaryOperationExpression(key, restapi.MatcherExpressionEntityDestination, restapi.IsEmptyOperator)
	input := restapi.NewBinaryOperator(primaryExpression, "FOO", primaryExpression)

	mapper := NewMatchExpressionMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "invalid conjunction operator")
}

func TestShouldReturnMappingErrorIfLeftSideOfConjunctionIsNotValid(t *testing.T) {
	key := "key"
	primaryExpressionLeft := restapi.NewUnaryOperationExpression(key, restapi.MatcherExpressionEntityDestination, "INVALID")
	primaryExpressionRight := restapi.NewUnaryOperationExpression(key, restapi.MatcherExpressionEntityDestination, restapi.IsEmptyOperator)
	input := restapi.NewBinaryOperator(primaryExpressionLeft, restapi.LogicalOr, primaryExpressionRight)

	mapper := NewMatchExpressionMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), invalidOperator)
	require.Contains(t, err.Error(), unaryOperator)
}

func TestShouldReturnMappingErrorIfRightSideOfConjunctionIsNotValid(t *testing.T) {
	key := "key"
	primaryExpressionLeft := restapi.NewUnaryOperationExpression(key, restapi.MatcherExpressionEntityDestination, restapi.IsEmptyOperator)
	primaryExpressionRight := restapi.NewUnaryOperationExpression(key, restapi.MatcherExpressionEntityDestination, "INVALID")
	input := restapi.NewBinaryOperator(primaryExpressionLeft, restapi.LogicalOr, primaryExpressionRight)

	mapper := NewMatchExpressionMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), invalidOperator)
	require.Contains(t, err.Error(), unaryOperator)
}

func runTestCaseForMappingFromAPI(input restapi.MatchExpression, expectedResult *FilterExpression, t *testing.T) {
	mapper := NewMatchExpressionMapper()
	result, err := mapper.FromAPIModel(input)

	require.Nil(t, err)
	require.Equal(t, expectedResult, result)
}
