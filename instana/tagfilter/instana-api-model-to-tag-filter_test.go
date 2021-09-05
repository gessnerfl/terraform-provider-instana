package filterexpression_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	. "github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
)

const (
	invalidOperator   = "invalid operator"
	tagFilterOperator = "tag filter operator"
)

func TestShouldMapValidOperatorsOfTagFilter(t *testing.T) {
	for _, v := range restapi.SupportedComparisonOperators {
		t.Run(fmt.Sprintf("test mapping of %s", v), testMappingOfOperatorsOfTagFilter(v))
	}
}

func testMappingOfOperatorsOfTagFilter(operator restapi.TagFilterOperator) func(t *testing.T) {
	return func(t *testing.T) {
		name := "name"
		value := "value"
		input := restapi.NewStringTagFilter(restapi.TagFilterEntityDestination, name, operator, &value)

		expectedResult := &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:      &EntitySpec{Identifier: name, Origin: EntityOriginDestination},
							Operator:    Operator(operator),
							StringValue: &value,
						},
					},
				},
			},
		}

		runTestCaseForMappingFromAPI(input, expectedResult, t)
	}
}

func TestShouldFailToMapComparisonWhenOperatorOfTagFilterIsNotValid(t *testing.T) {
	name := "name"
	value := "value"
	input := restapi.NewStringTagFilter(restapi.TagFilterEntityDestination, name, "FOO", &value)

	mapper := NewTagFilterMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), invalidOperator)
	require.Contains(t, err.Error(), tagFilterOperator)
}

func TestShouldMapValidUnaryOperationsOfTagFilter(t *testing.T) {
	for _, v := range restapi.SupportedUnaryExpressionOperators {
		t.Run(fmt.Sprintf("test mapping of %s ", v), testMappingOfUnaryOperationOfTagFilter(v))
	}
}

func testMappingOfUnaryOperationOfTagFilter(operator restapi.TagFilterOperator) func(t *testing.T) {
	return func(t *testing.T) {
		name := "name"
		input := restapi.NewUnaryTagFilter(restapi.TagFilterEntityDestination, name, operator)

		expectedResult := &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Entity:   &EntitySpec{Identifier: name, Origin: EntityOriginDestination},
							Operator: Operator(operator),
						},
					},
				},
			},
		}

		runTestCaseForMappingFromAPI(input, expectedResult, t)
	}
}

func TestShouldFailToMapUnaryOperationWhenOperatorOfTagFilterIsNotValid(t *testing.T) {
	name := "name"
	input := restapi.NewUnaryTagFilter(restapi.TagFilterEntityDestination, name, "FOO")

	mapper := NewTagFilterMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), invalidOperator)
	require.Contains(t, err.Error(), tagFilterOperator)
}

func TestShouldFailMapExpressionWhenTypeIsMissing(t *testing.T) {
	name := "name"
	input := &restapi.TagFilter{
		Name:     name,
		Operator: "FOO",
	}

	mapper := NewTagFilterMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "unsupported tag filter expression")
}

func TestShouldMapLogicalAndWhenLeftAndRightIsAPrimaryExpression(t *testing.T) {
	name := "name"
	operator := Operator(restapi.IsEmptyOperator)
	and := Operator(restapi.LogicalAnd)
	primaryExpression := restapi.NewUnaryTagFilter(restapi.TagFilterEntityDestination, name, restapi.IsEmptyOperator)
	input := restapi.NewLogicalAndTagFilter([]restapi.TagFilterExpressionElement{primaryExpression, primaryExpression})

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					UnaryOperation: &UnaryOperationExpression{
						Entity:   &EntitySpec{Identifier: name, Origin: EntityOriginDestination},
						Operator: operator,
					},
				},
				Operator: &and,
				Right: &LogicalAndExpression{
					Left: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Entity:   &EntitySpec{Identifier: name, Origin: EntityOriginDestination},
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
	name := "name"
	operator := Operator(restapi.IsEmptyOperator)
	and := Operator(restapi.LogicalAnd)
	primaryExpression := restapi.NewUnaryTagFilter(restapi.TagFilterEntityDestination, name, restapi.IsEmptyOperator)
	nestedAnd := restapi.NewLogicalAndTagFilter([]restapi.TagFilterExpressionElement{primaryExpression, primaryExpression})
	input := restapi.NewLogicalAndTagFilter([]restapi.TagFilterExpressionElement{primaryExpression, nestedAnd})

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					UnaryOperation: &UnaryOperationExpression{
						Entity:   &EntitySpec{Identifier: name, Origin: EntityOriginDestination},
						Operator: operator,
					},
				},
				Operator: &and,
				Right: &LogicalAndExpression{
					Left: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Entity:   &EntitySpec{Identifier: name, Origin: EntityOriginDestination},
							Operator: operator,
						},
					},
					Operator: &and,
					Right: &LogicalAndExpression{
						Left: &PrimaryExpression{
							UnaryOperation: &UnaryOperationExpression{
								Entity:   &EntitySpec{Identifier: name, Origin: EntityOriginDestination},
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
	name := "name"
	primaryExpression := restapi.NewUnaryTagFilter(restapi.TagFilterEntityDestination, name, restapi.IsEmptyOperator)
	nestedOr := restapi.NewLogicalOrTagFilter([]restapi.TagFilterExpressionElement{primaryExpression, primaryExpression})
	input := restapi.NewLogicalAndTagFilter([]restapi.TagFilterExpressionElement{nestedOr, primaryExpression})

	mapper := NewTagFilterMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "logical or is not allowed")
}

func TestShouldFailToMapLogicalAndWhenRightIsOrExpression(t *testing.T) {
	name := "name"
	primaryExpression := restapi.NewUnaryTagFilter(restapi.TagFilterEntityDestination, name, restapi.IsEmptyOperator)
	nestedOr := restapi.NewLogicalOrTagFilter([]restapi.TagFilterExpressionElement{primaryExpression, primaryExpression})
	input := restapi.NewLogicalAndTagFilter([]restapi.TagFilterExpressionElement{primaryExpression, nestedOr})

	mapper := NewTagFilterMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "logical or is not allowed")
}

func TestShouldFailToMapLogicalAndWhenLeftIsAndExpression(t *testing.T) {
	name := "name"
	primaryExpression := restapi.NewUnaryTagFilter(restapi.TagFilterEntityDestination, name, restapi.IsEmptyOperator)
	nestedAnd := restapi.NewLogicalAndTagFilter([]restapi.TagFilterExpressionElement{primaryExpression, primaryExpression})
	input := restapi.NewLogicalAndTagFilter([]restapi.TagFilterExpressionElement{nestedAnd, primaryExpression})

	mapper := NewTagFilterMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "logical and is not allowed for first element")
}

func TestShouldMapLogiclOrWhenLeftAndRightSideIsPrimaryExpression(t *testing.T) {
	name := "name"
	operator := Operator(restapi.IsEmptyOperator)
	or := Operator(restapi.LogicalOr)
	primaryExpression := restapi.NewUnaryTagFilter(restapi.TagFilterEntityDestination, name, restapi.IsEmptyOperator)
	input := restapi.NewLogicalOrTagFilter([]restapi.TagFilterExpressionElement{primaryExpression, primaryExpression})

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					UnaryOperation: &UnaryOperationExpression{
						Entity:   &EntitySpec{Identifier: name, Origin: EntityOriginDestination},
						Operator: operator,
					},
				},
			},
			Operator: &or,
			Right: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Entity:   &EntitySpec{Identifier: name, Origin: EntityOriginDestination},
							Operator: operator,
						},
					},
				},
			},
		},
	}

	runTestCaseForMappingFromAPI(input, expectedResult, t)
}

func TestShouldMapLogicalOrWhenLeftSideIsALogicalAndAndRightSideIsPrimaryExpression(t *testing.T) {
	name := "name"
	operator := Operator(restapi.IsEmptyOperator)
	or := Operator(restapi.LogicalOr)
	and := Operator(restapi.LogicalAnd)
	primaryExpression := restapi.NewUnaryTagFilter(restapi.TagFilterEntityDestination, name, restapi.IsEmptyOperator)
	nestedAnd := restapi.NewLogicalAndTagFilter([]restapi.TagFilterExpressionElement{primaryExpression, primaryExpression})
	input := restapi.NewLogicalOrTagFilter([]restapi.TagFilterExpressionElement{nestedAnd, primaryExpression})

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					UnaryOperation: &UnaryOperationExpression{
						Entity:   &EntitySpec{Identifier: name, Origin: EntityOriginDestination},
						Operator: operator,
					},
				},
				Operator: &and,
				Right: &LogicalAndExpression{
					Left: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Entity:   &EntitySpec{Identifier: name, Origin: EntityOriginDestination},
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
							Entity:   &EntitySpec{Identifier: name, Origin: EntityOriginDestination},
							Operator: operator,
						},
					},
				},
			},
		},
	}

	runTestCaseForMappingFromAPI(input, expectedResult, t)
}

func TestShouldMapLogicalOrWhenLeftSideIsAPrimaryExpressionAndRightSideIsALogicalOr(t *testing.T) {
	name := "name"
	operator := Operator(restapi.IsEmptyOperator)
	or := Operator(restapi.LogicalOr)
	primaryExpression := restapi.NewUnaryTagFilter(restapi.TagFilterEntityDestination, name, restapi.IsEmptyOperator)
	nestedOr := restapi.NewLogicalOrTagFilter([]restapi.TagFilterExpressionElement{primaryExpression, primaryExpression})
	input := restapi.NewLogicalOrTagFilter([]restapi.TagFilterExpressionElement{primaryExpression, nestedOr})

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					UnaryOperation: &UnaryOperationExpression{
						Entity:   &EntitySpec{Identifier: name, Origin: EntityOriginDestination},
						Operator: operator,
					},
				},
			},
			Operator: &or,
			Right: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Entity:   &EntitySpec{Identifier: name, Origin: EntityOriginDestination},
							Operator: operator,
						},
					},
				},
				Operator: &or,
				Right: &LogicalOrExpression{
					Left: &LogicalAndExpression{
						Left: &PrimaryExpression{
							UnaryOperation: &UnaryOperationExpression{
								Entity:   &EntitySpec{Identifier: name, Origin: EntityOriginDestination},
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

func TestShouldMapLogicalOrWhenLeftSideIsAPrimaryExpressionAndRightSideIsALogicalAnd(t *testing.T) {
	name := "name"
	operator := Operator(restapi.IsEmptyOperator)
	or := Operator(restapi.LogicalOr)
	and := Operator(restapi.LogicalAnd)
	primaryExpression := restapi.NewUnaryTagFilter(restapi.TagFilterEntityDestination, name, restapi.IsEmptyOperator)
	nestedAnd := restapi.NewLogicalAndTagFilter([]restapi.TagFilterExpressionElement{primaryExpression, primaryExpression})
	input := restapi.NewLogicalOrTagFilter([]restapi.TagFilterExpressionElement{primaryExpression, nestedAnd})

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					UnaryOperation: &UnaryOperationExpression{
						Entity:   &EntitySpec{Identifier: name, Origin: EntityOriginDestination},
						Operator: operator,
					},
				},
			},
			Operator: &or,
			Right: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Entity:   &EntitySpec{Identifier: name, Origin: EntityOriginDestination},
							Operator: operator,
						},
					},
					Operator: &and,
					Right: &LogicalAndExpression{
						Left: &PrimaryExpression{
							UnaryOperation: &UnaryOperationExpression{
								Entity:   &EntitySpec{Identifier: name, Origin: EntityOriginDestination},
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
	name := "name"
	primaryExpression := restapi.NewUnaryTagFilter(restapi.TagFilterEntityDestination, name, restapi.IsEmptyOperator)
	nestedOr := restapi.NewLogicalOrTagFilter([]restapi.TagFilterExpressionElement{primaryExpression, primaryExpression})
	input := restapi.NewLogicalOrTagFilter([]restapi.TagFilterExpressionElement{nestedOr, primaryExpression})

	mapper := NewTagFilterMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "logical or is not allowed for first element")
}

func TestShouldFailToMapBinaryExpressionWhenConjunctionTypeIsNotValid(t *testing.T) {
	name := "name"
	primaryExpression := restapi.NewUnaryTagFilter(restapi.TagFilterEntityDestination, name, restapi.IsEmptyOperator)
	input := &restapi.TagFilterExpression{
		Type:            restapi.TagFilterExpressionType,
		LogicalOperator: restapi.LogicalOperatorType("FOO"),
		Elements:        []restapi.TagFilterExpressionElement{primaryExpression, primaryExpression},
	}

	mapper := NewTagFilterMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "invalid conjunction operator")
}

func TestShouldReturnMappingErrorIfLeftSideOfConjunctionIsNotValid(t *testing.T) {
	name := "name"
	primaryExpressionLeft := restapi.NewUnaryTagFilter(restapi.TagFilterEntityDestination, name, "INVALID")
	primaryExpressionRight := restapi.NewUnaryTagFilter(restapi.TagFilterEntityDestination, name, restapi.IsEmptyOperator)
	input := restapi.NewLogicalOrTagFilter([]restapi.TagFilterExpressionElement{primaryExpressionLeft, primaryExpressionRight})

	mapper := NewTagFilterMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), invalidOperator)
	require.Contains(t, err.Error(), tagFilterOperator)
}

func TestShouldReturnMappingErrorIfRightSideOfConjunctionIsNotValid(t *testing.T) {
	name := "name"
	primaryExpressionLeft := restapi.NewUnaryTagFilter(restapi.TagFilterEntityDestination, name, restapi.IsEmptyOperator)
	primaryExpressionRight := restapi.NewUnaryTagFilter(restapi.TagFilterEntityDestination, name, "INVALID")
	input := restapi.NewLogicalOrTagFilter([]restapi.TagFilterExpressionElement{primaryExpressionLeft, primaryExpressionRight})

	mapper := NewTagFilterMapper()
	_, err := mapper.FromAPIModel(input)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), invalidOperator)
	require.Contains(t, err.Error(), tagFilterOperator)
}

func runTestCaseForMappingFromAPI(input restapi.TagFilterExpressionElement, expectedResult *FilterExpression, t *testing.T) {
	mapper := NewTagFilterMapper()
	result, err := mapper.FromAPIModel(input)

	require.Nil(t, err)
	require.Equal(t, expectedResult, result)
}
