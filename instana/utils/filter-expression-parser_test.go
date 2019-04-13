package utils_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/utils"
	"github.com/google/go-cmp/cmp"
)

func TestShouldSuccessfullyParseComplexExpression(t *testing.T) {
	expression := "entity.name CO 'foo'OR entity.kind EQ 2.34 AND entity.type EQ true AND span.name NOT EMPTY OR span.id NE 1234"

	foo := "foo"
	decNum := 2.34
	intNum := float64(1234)
	isTrue := Boolean(true)
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
			Operator: "OR",
			Right: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						Comparision: &ComparisionExpression{
							Key:      "entity.kind",
							Operator: "EQ",
							Value:    &Value{Number: &decNum},
						},
					},
					Operator: "AND",
					Right: &LogicalAndExpression{
						Left: &PrimaryExpression{
							Comparision: &ComparisionExpression{
								Key:      "entity.type",
								Operator: "EQ",
								Value:    &Value{Boolean: &isTrue},
							},
						},
						Operator: "AND",
						Right: &LogicalAndExpression{
							Left: &PrimaryExpression{
								UnaryExpression: &UnaryExpression{
									Key:      "span.name",
									Function: "NOTEMPTY",
								},
							},
						},
					},
				},
				Operator: "OR",
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
	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					Comparision: &ComparisionExpression{
						Key:      "entity.name",
						Operator: "co",
						Value:    &Value{String: &foo},
					},
				},
				Operator: "and",
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
						Operator: "eQ",
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
						Function: "NotEmptY",
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
	expr := "Foo bla 'bar'"

	sut := NewFilterExpressionParser()
	_, err := sut.Parse(expr)

	if err == nil {
		t.Fatal("Expected parsing error")
	}
}
