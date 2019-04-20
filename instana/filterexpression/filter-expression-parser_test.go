package filterexpression_test

import (
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"

	. "github.com/gessnerfl/terraform-provider-instana/instana/filterexpression"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/google/go-cmp/cmp"
)

const (
	keyEntityName = "entity.name"
	keyEntityKind = "entity.kind"
	keyEntityType = "entity.type"

	valueMyValue = "my value"

	messageExpectedNormalizedExpression = "Expected normalized rendered result of comparision expression but got:  %s"

	entityNameEqualsValueExpression = "entity.name EQUALS 'my value'"
)

func TestShouldSuccessfullyParseComplexExpression(t *testing.T) {
	expression := "entity.name CONTAINS 'foo bar' OR entity.kind EQUALS '2.34' AND entity.type EQUALS 'true' AND span.name NOT_EMPTY OR span.id NOT_EQUAL  '1234'"

	logicalAnd := Operator(restapi.LogicalAnd)
	logicalOr := Operator(restapi.LogicalOr)
	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					Comparision: &ComparisionExpression{
						Key:      keyEntityName,
						Operator: Operator(restapi.ContainsOperator),
						Value:    "foo bar",
					},
				},
			},
			Operator: &logicalOr,
			Right: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						Comparision: &ComparisionExpression{
							Key:      keyEntityKind,
							Operator: Operator(restapi.EqualsOperator),
							Value:    "2.34",
						},
					},
					Operator: &logicalAnd,
					Right: &LogicalAndExpression{
						Left: &PrimaryExpression{
							Comparision: &ComparisionExpression{
								Key:      keyEntityType,
								Operator: Operator(restapi.EqualsOperator),
								Value:    "true",
							},
						},
						Operator: &logicalAnd,
						Right: &LogicalAndExpression{
							Left: &PrimaryExpression{
								UnaryOperation: &UnaryOperationExpression{
									Key:      "span.name",
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
							Comparision: &ComparisionExpression{
								Key:      "span.id",
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
					Comparision: &ComparisionExpression{
						Key:      keyEntityName,
						Operator: Operator(restapi.ContainsOperator),
						Value:    "foo",
					},
				},
				Operator: &logicalAnd,
				Right: &LogicalAndExpression{
					Left: &PrimaryExpression{
						Comparision: &ComparisionExpression{
							Key:      keyEntityType,
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

func TestShouldParseComparisionOperationsCaseInsensitive(t *testing.T) {
	expression := "entity.name EQUALS 'foo'"

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					Comparision: &ComparisionExpression{
						Key:      keyEntityName,
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
	expression := "entity.name NOT_EMPTY"

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					UnaryOperation: &UnaryOperationExpression{
						Key:      keyEntityName,
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

	if err != nil {
		t.Fatalf("Did not expected error but got %s", err)
	}

	if !cmp.Equal(expectedResult, result) {
		t.Fatalf("Expected parse expression %v but got %v; diff %s", expectedResult, result, cmp.Diff(expectedResult, result))
	}
}

func TestShouldFailToParseInvalidExpression(t *testing.T) {
	expression := "Foo invalidToken 'bar'"

	sut := NewParser()
	_, err := sut.Parse(expression)

	if err == nil {
		t.Fatal("Expected parsing error")
	}
}

func TestShouldRenderComplexExpressionNormalizedForm(t *testing.T) {
	expression := "entity.name CONTAINS 'foo' OR entity.kind EQUALS '2.34'    and  entity.type EQUALS 'true'  AND span.name  NOT_EMPTY   OR span.id  NOT_EQUAL  '1234'"
	normalizedExpression := "entity.name CONTAINS 'foo' OR entity.kind EQUALS '2.34' AND entity.type EQUALS 'true' AND span.name NOT_EMPTY OR span.id NOT_EQUAL '1234'"

	sut := NewParser()
	result, err := sut.Parse(expression)

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}

	rendered := result.Render()
	if rendered != normalizedExpression {
		t.Fatalf("Expected to get normalized expression rendered but got: %s", rendered)
	}
}

func TestShouldRenderLogicalOrExpressionWhenOrIsSet(t *testing.T) {
	expectedResult := "foo EQUALS 'bar' OR foo CONTAINS 'bar'"

	logicalOr := Operator(restapi.LogicalOr)
	sut := LogicalOrExpression{
		Left: &LogicalAndExpression{
			Left: &PrimaryExpression{
				Comparision: &ComparisionExpression{
					Key:      "foo",
					Operator: Operator(restapi.EqualsOperator),
					Value:    "bar",
				},
			},
		},
		Operator: &logicalOr,
		Right: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					Comparision: &ComparisionExpression{
						Key:      "foo",
						Operator: Operator(restapi.ContainsOperator),
						Value:    "bar",
					},
				},
			},
		},
	}

	rendered := sut.Render()

	if rendered != expectedResult {
		t.Fatalf("Expected normalized rendered result of logical OR but got:  %s", rendered)
	}
}

func TestShouldRenderPrimaryExpressionOnLogicalOrExpressionWhenNeitherOrNorAndIsSet(t *testing.T) {
	sut := LogicalOrExpression{
		Left: &LogicalAndExpression{
			Left: &PrimaryExpression{
				Comparision: &ComparisionExpression{
					Key:      keyEntityName,
					Operator: Operator(restapi.EqualsOperator),
					Value:    valueMyValue,
				},
			},
		},
	}

	rendered := sut.Render()

	if rendered != entityNameEqualsValueExpression {
		t.Fatalf(messageExpectedNormalizedExpression, rendered)
	}
}

func TestShouldRenderLogicalAndExpressionWhenAndIsSet(t *testing.T) {
	expectedResult := "foo EQUALS 'bar' AND foo CONTAINS 'bar'"

	logicalAnd := Operator(restapi.LogicalAnd)
	sut := LogicalAndExpression{
		Left: &PrimaryExpression{
			Comparision: &ComparisionExpression{
				Key:      "foo",
				Operator: Operator(restapi.EqualsOperator),
				Value:    "bar",
			},
		},
		Operator: &logicalAnd,
		Right: &LogicalAndExpression{
			Left: &PrimaryExpression{
				Comparision: &ComparisionExpression{
					Key:      "foo",
					Operator: Operator(restapi.ContainsOperator),
					Value:    "bar",
				},
			},
		},
	}

	rendered := sut.Render()

	if rendered != expectedResult {
		t.Fatalf("Expected normalized rendered result of logical AND but got:  %s", rendered)
	}
}

func TestShouldRenderPrimaryExpressionOnLogicalAndExpressionWhenAndIsNotSet(t *testing.T) {
	sut := LogicalAndExpression{
		Left: &PrimaryExpression{
			Comparision: &ComparisionExpression{
				Key:      keyEntityName,
				Operator: Operator(restapi.EqualsOperator),
				Value:    valueMyValue,
			},
		},
	}

	rendered := sut.Render()

	if rendered != entityNameEqualsValueExpression {
		t.Fatalf(messageExpectedNormalizedExpression, rendered)
	}
}

func TestShouldRenderComparisionOnPrimaryExpressionWhenComparsionIsSet(t *testing.T) {
	sut := PrimaryExpression{
		Comparision: &ComparisionExpression{
			Key:      keyEntityName,
			Operator: Operator(restapi.EqualsOperator),
			Value:    valueMyValue,
		},
	}

	rendered := sut.Render()

	if rendered != entityNameEqualsValueExpression {
		t.Fatalf(messageExpectedNormalizedExpression, rendered)
	}
}

func TestShouldRenderUnaryOperationExpressionOnPrimaryExpressionWhenUnaryOperationIsSet(t *testing.T) {
	expectedResult := "foo IS_EMPTY"

	sut := PrimaryExpression{
		UnaryOperation: &UnaryOperationExpression{
			Key:      "foo",
			Operator: Operator(restapi.IsEmptyOperator),
		},
	}

	rendered := sut.Render()

	if rendered != expectedResult {
		t.Fatalf(messageExpectedNormalizedExpression, rendered)
	}
}
