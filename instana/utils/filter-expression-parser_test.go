package utils_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/utils"
	"github.com/google/go-cmp/cmp"
)

func TestShouldSuccessfullyParseComplexExpression(t *testing.T) {
	expression := "entity.name CO 'foo'OR entity.kind EQ 2.34 AND entity.type EQ true AND span.name NOT EMPTY OR span.id NE  1234"

	foo := "foo"
	decNum := 2.34
	intNum := float64(1234)
	isTrue := Boolean(true)
	logicalAnd := Operator("AND")
	logicalOr := Operator("OR")
	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					Comparision: &ComparisionExpression{
						Key:      "entity.name",
						Operator: "CO",
						Value:    &Value{String: &foo},
					},
				},
			},
			Operator: &logicalOr,
			Right: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						Comparision: &ComparisionExpression{
							Key:      "entity.kind",
							Operator: "EQ",
							Value:    &Value{Number: &decNum},
						},
					},
					Operator: &logicalAnd,
					Right: &LogicalAndExpression{
						Left: &PrimaryExpression{
							Comparision: &ComparisionExpression{
								Key:      "entity.type",
								Operator: "EQ",
								Value:    &Value{Boolean: &isTrue},
							},
						},
						Operator: &logicalAnd,
						Right: &LogicalAndExpression{
							Left: &PrimaryExpression{
								UnaryExpression: &UnaryExpression{
									Key:      "span.name",
									Function: "NOT EMPTY",
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
								Operator: "NE",
								Value:    &Value{Number: &intNum},
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
	expression := "entity.name co 'foo' and entity.type EQ 'bar'"

	foo := "foo"
	bar := "bar"
	logicalAnd := Operator("AND")
	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					Comparision: &ComparisionExpression{
						Key:      "entity.name",
						Operator: "CO",
						Value:    &Value{String: &foo},
					},
				},
				Operator: &logicalAnd,
				Right: &LogicalAndExpression{
					Left: &PrimaryExpression{
						Comparision: &ComparisionExpression{
							Key:      "entity.type",
							Operator: "EQ",
							Value:    &Value{String: &bar},
						},
					},
				},
			},
		},
	}

	shouldSuccessfullyParseExpression(expression, expectedResult, t)
}

func TestShouldParseComparisionOperationsCaseInsensitive(t *testing.T) {
	expression := "entity.name eQ 'foo'"

	foo := "foo"
	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					Comparision: &ComparisionExpression{
						Key:      "entity.name",
						Operator: "EQ",
						Value:    &Value{String: &foo},
					},
				},
			},
		},
	}

	shouldSuccessfullyParseExpression(expression, expectedResult, t)
}

func TestShouldParseUnaryOperationsCaseInsensitive(t *testing.T) {
	expression := "entity.name Not EmptY"

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					UnaryExpression: &UnaryExpression{
						Key:      "entity.name",
						Function: "NOT EMPTY",
					},
				},
			},
		},
	}

	shouldSuccessfullyParseExpression(expression, expectedResult, t)
}

func shouldSuccessfullyParseExpression(input string, expectedResult *FilterExpression, t *testing.T) {
	sut := NewFilterExpressionParser()
	result, err := sut.Parse(input)

	if err != nil {
		t.Fatalf("Did not expected error but got %s", err)
	}

	if !cmp.Equal(expectedResult, result) {
		t.Fatalf("Expected parse expression %v but got %v; diff %s", expectedResult, result, cmp.Diff(expectedResult, result))
	}
}

func TestShouldFailToParseInvalidExpression(t *testing.T) {
	expression := "Foo bla 'bar'"

	sut := NewFilterExpressionParser()
	_, err := sut.Parse(expression)

	if err == nil {
		t.Fatal("Expected parsing error")
	}
}

func TestShouldRenderComplexExpressionNormalizedForm(t *testing.T) {
	expression := "entity.name co 'foo' OR entity.kind EQ 2.34    and  entity.type EQ TRUE  AND span.name  NOT empTy   OR span.id  NE  1234"
	normalizedExpression := "entity.name CO 'foo' OR entity.kind EQ 2.340 AND entity.type EQ true AND span.name NOT EMPTY OR span.id NE 1234.000"

	sut := NewFilterExpressionParser()
	result, err := sut.Parse(expression)

	if err != nil {
		t.Fatalf("Expected no error but got '%s'", err)
	}

	rendered := result.Render()
	if rendered != normalizedExpression {
		t.Fatalf("Expected to get normalized expression rendered but got: %s", rendered)
	}
}

func TestShouldRenderLogicalOrExpressionWhenOrIsSet(t *testing.T) {
	expectedResult := "foo EQ 'bar' OR foo CO 'bar'"
	bar := "bar"

	logicalOr := Operator("OR")
	sut := LogicalOrExpression{
		Left: &LogicalAndExpression{
			Left: &PrimaryExpression{
				Comparision: &ComparisionExpression{
					Key:      "foo",
					Operator: "EQ",
					Value: &Value{
						String: &bar,
					},
				},
			},
		},
		Operator: &logicalOr,
		Right: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					Comparision: &ComparisionExpression{
						Key:      "foo",
						Operator: "CO",
						Value: &Value{
							String: &bar,
						},
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
	expectedResult := "foo EQ 'bar'"
	bar := "bar"

	sut := LogicalOrExpression{
		Left: &LogicalAndExpression{
			Left: &PrimaryExpression{
				Comparision: &ComparisionExpression{
					Key:      "foo",
					Operator: "EQ",
					Value: &Value{
						String: &bar,
					},
				},
			},
		},
	}

	rendered := sut.Render()

	if rendered != expectedResult {
		t.Fatalf("Expected normalized rendered result of comparision expression but got:  %s", rendered)
	}
}

func TestShouldRenderLogicalAndExpressionWhenAndIsSet(t *testing.T) {
	expectedResult := "foo EQ 'bar' AND foo CO 'bar'"
	bar := "bar"

	logicalAnd := Operator("AND")
	sut := LogicalAndExpression{
		Left: &PrimaryExpression{
			Comparision: &ComparisionExpression{
				Key:      "foo",
				Operator: "EQ",
				Value: &Value{
					String: &bar,
				},
			},
		},
		Operator: &logicalAnd,
		Right: &LogicalAndExpression{
			Left: &PrimaryExpression{
				Comparision: &ComparisionExpression{
					Key:      "foo",
					Operator: "CO",
					Value: &Value{
						String: &bar,
					},
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
	expectedResult := "foo EQ 'bar'"
	bar := "bar"

	sut := LogicalAndExpression{
		Left: &PrimaryExpression{
			Comparision: &ComparisionExpression{
				Key:      "foo",
				Operator: "EQ",
				Value: &Value{
					String: &bar,
				},
			},
		},
	}

	rendered := sut.Render()

	if rendered != expectedResult {
		t.Fatalf("Expected normalized rendered result of comparision expression but got:  %s", rendered)
	}
}

func TestShouldRenderComparisionOnPrimaryExpressionWhenComparsionIsSet(t *testing.T) {
	expectedResult := "foo EQ 'bar'"
	bar := "bar"

	sut := PrimaryExpression{
		Comparision: &ComparisionExpression{
			Key:      "foo",
			Operator: "EQ",
			Value: &Value{
				String: &bar,
			},
		},
	}

	rendered := sut.Render()

	if rendered != expectedResult {
		t.Fatalf("Expected normalized rendered result of comparision expression but got:  %s", rendered)
	}
}

func TestShouldRenderUnaryExpressionOnPrimaryExpressionWhenUnaryExpressionIsSet(t *testing.T) {
	expectedResult := "foo IS EMPTY"

	sut := PrimaryExpression{
		UnaryExpression: &UnaryExpression{
			Key:      "foo",
			Function: "IS EMPTY",
		},
	}

	rendered := sut.Render()

	if rendered != expectedResult {
		t.Fatalf("Expected normalized rendered result of comparision expression but got:  %s", rendered)
	}
}

func TestShouldRenderValueWhenStringIsSet(t *testing.T) {
	expectedResult := "'foo'"

	value := "foo"
	sut := Value{
		String: &value,
	}

	rendered := sut.Render()

	if rendered != expectedResult {
		t.Fatalf("Expected normalized rendered result of comparision expression but got:  %s", rendered)
	}
}

func TestShouldRenderValueWhenDecimalNumberIsSet(t *testing.T) {
	expectedResult := "1.230"

	value := 1.23
	sut := Value{
		Number: &value,
	}

	rendered := sut.Render()

	if rendered != expectedResult {
		t.Fatalf("Expected normalized rendered result of comparision expression but got:  %s", rendered)
	}
}

func TestShouldRenderValueWhenIntegerNumberIsSet(t *testing.T) {
	expectedResult := "123.000"

	value := float64(123)
	sut := Value{
		Number: &value,
	}

	rendered := sut.Render()

	if rendered != expectedResult {
		t.Fatalf("Expected normalized rendered result of comparision expression but got:  %s", rendered)
	}
}

func TestShouldRenderValueWhenBooleanIsSet(t *testing.T) {
	expectedResult := "true"

	boolTrue := Boolean(true)
	sut := Value{
		Boolean: &boolTrue,
	}

	rendered := sut.Render()

	if rendered != expectedResult {
		t.Fatalf("Expected normalized rendered result of comparision expression but got:  %s", rendered)
	}
}
