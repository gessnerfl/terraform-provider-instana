package filterexpression_test

import (
	"fmt"
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana/filterexpression"
)

const (
	keyEntityName = "entity.name"
	keyEntityKind = "entity.kind"
	keyEntityType = "entity.type"

	valueMyValue = "my value"

	entityNameEqualsValueExpression = "entity.name@dest EQUALS 'my value'"
)

func TestShouldSuccessfullyParseComplexExpression(t *testing.T) {
	expression := "entity.name CONTAINS 'foo bar' OR entity.kind EQUALS '2.34' AND entity.type EQUALS 'true' AND span.name NOT_EMPTY OR span.id NOT_EQUAL  '1234'"

	logicalAnd := Operator(restapi.LogicalAnd)
	logicalOr := Operator(restapi.LogicalOr)
	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					Comparison: &ComparisonExpression{
						Entity:   &EntitySpec{Identifier: keyEntityName, Origin: EntityOriginDestination},
						Operator: Operator(restapi.ContainsOperator),
						Value:    "foo bar",
					},
				},
			},
			Operator: &logicalOr,
			Right: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:   &EntitySpec{Identifier: keyEntityKind, Origin: EntityOriginDestination},
							Operator: Operator(restapi.EqualsOperator),
							Value:    "2.34",
						},
					},
					Operator: &logicalAnd,
					Right: &LogicalAndExpression{
						Left: &PrimaryExpression{
							Comparison: &ComparisonExpression{
								Entity:   &EntitySpec{Identifier: keyEntityType, Origin: EntityOriginDestination},
								Operator: Operator(restapi.EqualsOperator),
								Value:    "true",
							},
						},
						Operator: &logicalAnd,
						Right: &LogicalAndExpression{
							Left: &PrimaryExpression{
								UnaryOperation: &UnaryOperationExpression{
									Entity:   &EntitySpec{Identifier: "span.name", Origin: EntityOriginDestination},
									Operator: Operator(restapi.NotEmptyOperator),
								},
							},
						},
					},
				},
				Operator: &logicalOr,
				Right: &LogicalOrExpression{
					Left: &LogicalAndExpression{
						Left: &PrimaryExpression{
							Comparison: &ComparisonExpression{
								Entity:   &EntitySpec{Identifier: "span.id", Origin: EntityOriginDestination},
								Operator: Operator(restapi.NotEqualOperator),
								Value:    "1234",
							},
						},
					},
				},
			},
		},
	}

	shouldSuccessfullyParseExpression(expression, expectedResult, t)
}

func TestShouldParseKeywordsCaseInsensitive(t *testing.T) {
	expression := "entity.name CONTAINS 'foo' and entity.type EQUALS 'bar'"

	logicalAnd := Operator(restapi.LogicalAnd)
	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					Comparison: &ComparisonExpression{
						Entity:   &EntitySpec{Identifier: keyEntityName, Origin: EntityOriginDestination},
						Operator: Operator(restapi.ContainsOperator),
						Value:    "foo",
					},
				},
				Operator: &logicalAnd,
				Right: &LogicalAndExpression{
					Left: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:   &EntitySpec{Identifier: keyEntityType, Origin: EntityOriginDestination},
							Operator: Operator(restapi.EqualsOperator),
							Value:    "bar",
						},
					},
				},
			},
		},
	}

	shouldSuccessfullyParseExpression(expression, expectedResult, t)
}

func TestShouldParseAllSupportedComparisionOperators(t *testing.T) {
	for _, o := range restapi.SupportedComparisonOperators {
		t.Run(fmt.Sprintf("TestShouldParseComparisionOperator%s", string(o)), createTestCaseForParsingSupportedComparisionOperators(o))
	}
}

func createTestCaseForParsingSupportedComparisionOperators(operator restapi.TagFilterOperator) func(*testing.T) {
	return func(t *testing.T) {
		expression := fmt.Sprintf("entity.name %s 'foo'", string(operator))

		expectedResult := &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:   &EntitySpec{Identifier: keyEntityName, Origin: EntityOriginDestination},
							Operator: Operator(operator),
							Value:    "foo",
						},
					},
				},
			},
		}

		shouldSuccessfullyParseExpression(expression, expectedResult, t)
	}
}

func TestShouldParseAllSupportedUnaryOperators(t *testing.T) {
	for _, o := range restapi.SupportedUnaryExpressionOperators {
		t.Run(fmt.Sprintf("TestShouldParseUnaryOperator%s", string(o)), createTestCaseForParsingSupportedUnaryOperators(o))
	}
}

func createTestCaseForParsingSupportedUnaryOperators(operator restapi.TagFilterOperator) func(*testing.T) {
	return func(t *testing.T) {
		expression := fmt.Sprintf("entity.name %s", string(operator))

		expectedResult := &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Entity:   &EntitySpec{Identifier: keyEntityName, Origin: EntityOriginDestination},
							Operator: Operator(operator),
						},
					},
				},
			},
		}

		shouldSuccessfullyParseExpression(expression, expectedResult, t)
	}
}

func TestShouldParseComparisionOperationsCaseInsensitive(t *testing.T) {
	expression := "entity.name Equals 'foo'"

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					Comparison: &ComparisonExpression{
						Entity:   &EntitySpec{Identifier: keyEntityName, Origin: EntityOriginDestination},
						Operator: Operator(restapi.EqualsOperator),
						Value:    "foo",
					},
				},
			},
		},
	}

	shouldSuccessfullyParseExpression(expression, expectedResult, t)
}

func TestShouldParseUnaryOperationsCaseInsensitive(t *testing.T) {
	expression := "entity.name not_Empty"

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					UnaryOperation: &UnaryOperationExpression{
						Entity:   &EntitySpec{Identifier: keyEntityName, Origin: EntityOriginDestination},
						Operator: Operator(restapi.NotEmptyOperator),
					},
				},
			},
		},
	}

	shouldSuccessfullyParseExpression(expression, expectedResult, t)
}

func TestShouldParseIdentifiersWithDashes(t *testing.T) {
	expression := "call.http.header.x-example-foo EQUALS 'test'"

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					Comparison: &ComparisonExpression{
						Entity:   &EntitySpec{Identifier: "call.http.header.x-example-foo", Origin: EntityOriginDestination},
						Operator: Operator(restapi.EqualsOperator),
						Value:    "test",
					},
				},
			},
		},
	}

	shouldSuccessfullyParseExpression(expression, expectedResult, t)
}

func TestShouldParseIdentifierWithSlashes(t *testing.T) {
	expression := "kubernetes.pod.label.foo/bar EQUALS 'test'"

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					Comparison: &ComparisonExpression{
						Entity:   &EntitySpec{Identifier: "kubernetes.pod.label.foo/bar", Origin: EntityOriginDestination},
						Operator: Operator(restapi.EqualsOperator),
						Value:    "test",
					},
				},
			},
		},
	}

	shouldSuccessfullyParseExpression(expression, expectedResult, t)
}

func TestShouldParseEntityOriginFromComparisionExpression(t *testing.T) {
	expression := "entity.name@src EQUALS 'test'"

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					Comparison: &ComparisonExpression{
						Entity:   &EntitySpec{Identifier: keyEntityName, Origin: EntityOriginSource, OriginDefined: true},
						Operator: Operator(restapi.EqualsOperator),
						Value:    "test",
					},
				},
			},
		},
	}

	shouldSuccessfullyParseExpression(expression, expectedResult, t)
}

func TestShouldParseEntityOriginFromUnaryExpression(t *testing.T) {
	expression := "entity.name@src NOT_EMPTY"

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					UnaryOperation: &UnaryOperationExpression{
						Entity:   &EntitySpec{Identifier: keyEntityName, Origin: EntityOriginSource, OriginDefined: true},
						Operator: Operator(restapi.NotEmptyOperator),
					},
				},
			},
		},
	}

	shouldSuccessfullyParseExpression(expression, expectedResult, t)
}

func shouldSuccessfullyParseExpression(input string, expectedResult *FilterExpression, t *testing.T) {
	sut := NewParser()
	result, err := sut.Parse(input)

	require.Nil(t, err)
	require.Equal(t, expectedResult, result)
}

func TestShouldFailToParseInvalidExpression(t *testing.T) {
	expression := "Foo invalidToken 'bar'"

	sut := NewParser()
	_, err := sut.Parse(expression)

	require.NotNil(t, err)
}

func TestShouldRenderComplexExpressionNormalizedForm(t *testing.T) {
	expression := "entity.name CONTAINS 'foo' OR entity.kind EQUALS '2.34'    and  entity.type EQUALS 'true'  AND span.name  NOT_EMPTY   OR span.id  NOT_EQUAL  '1234'"
	normalizedExpression := "entity.name@dest CONTAINS 'foo' OR entity.kind@dest EQUALS '2.34' AND entity.type@dest EQUALS 'true' AND span.name@dest NOT_EMPTY OR span.id@dest NOT_EQUAL '1234'"

	sut := NewParser()
	result, err := sut.Parse(expression)
	require.Nil(t, err)

	rendered := result.Render()
	require.Equal(t, normalizedExpression, rendered)
}

func TestShouldRenderLogicalOrExpressionWhenOrIsSet(t *testing.T) {
	expectedResult := "foo@dest EQUALS 'bar' OR foo@dest CONTAINS 'bar'"

	logicalOr := Operator(restapi.LogicalOr)
	sut := LogicalOrExpression{
		Left: &LogicalAndExpression{
			Left: &PrimaryExpression{
				Comparison: &ComparisonExpression{
					Entity:   &EntitySpec{Identifier: "foo", Origin: EntityOriginDestination},
					Operator: Operator(restapi.EqualsOperator),
					Value:    "bar",
				},
			},
		},
		Operator: &logicalOr,
		Right: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					Comparison: &ComparisonExpression{
						Entity:   &EntitySpec{Identifier: "foo", Origin: EntityOriginDestination},
						Operator: Operator(restapi.ContainsOperator),
						Value:    "bar",
					},
				},
			},
		},
	}

	rendered := sut.Render()

	require.Equal(t, expectedResult, rendered)
}

func TestShouldRenderPrimaryExpressionOnLogicalOrExpressionWhenNeitherOrNorAndIsSet(t *testing.T) {
	sut := LogicalOrExpression{
		Left: &LogicalAndExpression{
			Left: &PrimaryExpression{
				Comparison: &ComparisonExpression{
					Entity:   &EntitySpec{Identifier: keyEntityName, Origin: EntityOriginDestination},
					Operator: Operator(restapi.EqualsOperator),
					Value:    valueMyValue,
				},
			},
		},
	}

	rendered := sut.Render()

	require.Equal(t, entityNameEqualsValueExpression, rendered)
}

func TestShouldRenderLogicalAndExpressionWhenAndIsSet(t *testing.T) {
	expectedResult := "foo@dest EQUALS 'bar' AND foo@dest CONTAINS 'bar'"

	logicalAnd := Operator(restapi.LogicalAnd)
	sut := LogicalAndExpression{
		Left: &PrimaryExpression{
			Comparison: &ComparisonExpression{
				Entity:   &EntitySpec{Identifier: "foo", Origin: EntityOriginDestination},
				Operator: Operator(restapi.EqualsOperator),
				Value:    "bar",
			},
		},
		Operator: &logicalAnd,
		Right: &LogicalAndExpression{
			Left: &PrimaryExpression{
				Comparison: &ComparisonExpression{
					Entity:   &EntitySpec{Identifier: "foo", Origin: EntityOriginDestination},
					Operator: Operator(restapi.ContainsOperator),
					Value:    "bar",
				},
			},
		},
	}

	rendered := sut.Render()
	require.Equal(t, expectedResult, rendered)
}

func TestShouldRenderPrimaryExpressionOnLogicalAndExpressionWhenAndIsNotSet(t *testing.T) {
	sut := LogicalAndExpression{
		Left: &PrimaryExpression{
			Comparison: &ComparisonExpression{
				Entity:   &EntitySpec{Identifier: keyEntityName, Origin: EntityOriginDestination},
				Operator: Operator(restapi.EqualsOperator),
				Value:    valueMyValue,
			},
		},
	}

	rendered := sut.Render()

	require.Equal(t, entityNameEqualsValueExpression, rendered)
}

func TestShouldRenderComparisionOnPrimaryExpressionWhenComparsionIsSet(t *testing.T) {
	sut := PrimaryExpression{
		Comparison: &ComparisonExpression{
			Entity:   &EntitySpec{Identifier: keyEntityName, Origin: EntityOriginDestination},
			Operator: Operator(restapi.EqualsOperator),
			Value:    valueMyValue,
		},
	}

	rendered := sut.Render()

	require.Equal(t, entityNameEqualsValueExpression, rendered)
}

func TestShouldRenderUnaryOperationExpressionOnPrimaryExpressionWhenUnaryOperationIsSet(t *testing.T) {
	expectedResult := "foo@dest IS_EMPTY"

	sut := PrimaryExpression{
		UnaryOperation: &UnaryOperationExpression{
			Entity:   &EntitySpec{Identifier: "foo", Origin: EntityOriginDestination},
			Operator: Operator(restapi.IsEmptyOperator),
		},
	}

	rendered := sut.Render()

	require.Equal(t, expectedResult, rendered)
}

func TestShouldGetEntityOriginByKey(t *testing.T) {
	for _, o := range SupportedEntityOrigins {
		t.Run(fmt.Sprintf("TestShouldGetEntityOriginForKey%s", o.Key()), func(t *testing.T) {
			require.Equal(t, o, SupportedEntityOrigins.ForKey(o.Key()))
		})
	}
}

func TestShouldReturnEntityOriginDestinationAsFallbackValueWhenKeyIsNotValid(t *testing.T) {
	require.Equal(t, EntityOriginDestination, SupportedEntityOrigins.ForKey("invalid"))
}

func TestShouldGetEntityOriginByInstanaAPIEntity(t *testing.T) {
	for _, e := range restapi.SupportedMatcherExpressionEntities {
		t.Run(fmt.Sprintf("TestShouldGetEntityOriginForInstanaAPIEntity%s", e), func(t *testing.T) {
			require.Equal(t, e, SupportedEntityOrigins.ForInstanaAPIEntity(e).MatcherExpressionEntity())
		})
	}
}

func TestShouldReturnEntityOriginDestinationAsFallbackValueWhenMatcherExpressionEntityIsNotValid(t *testing.T) {
	require.Equal(t, EntityOriginDestination, SupportedEntityOrigins.ForInstanaAPIEntity(restapi.MatcherExpressionEntity("invalid")))
}

func TestShouldNormalizeExpression(t *testing.T) {
	input := "entity.name    NOT_EMPTY"
	expectedResult := "entity.name@dest NOT_EMPTY"

	result, err := Normalize(input)
	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestShouldFailToNormalizeExpressionWehnExpressionIsNotValied(t *testing.T) {
	input := "entity.name    bla bla bla"

	result, err := Normalize(input)
	require.Error(t, err)
	require.Equal(t, input, result)
}
