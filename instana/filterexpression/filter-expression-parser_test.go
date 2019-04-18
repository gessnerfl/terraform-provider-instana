package filterexpression_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/filterexpression"
	"github.com/google/go-cmp/cmp"
)

func TestShouldSuccessfullyParseComplexExpression(t *testing.T) {
	expression := "entity.name CONTAINS 'foo bar' OR entity.kind EQUALS '2.34' AND entity.type EQUALS 'true' AND span.name NOT_EMPTY OR ( span.id NOT_EQUAL  '1234' OR span.id NOT_EQUAL '6789' )"

	logicalAnd := Operator("AND")
	logicalOr := Operator("OR")
	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					Comparision: &ComparisionExpression{
						Key:      "entity.name",
						Operator: "CONTAINS",
						Value:    "foo bar",
					},
				},
			},
			Operator: &logicalOr,
			Right: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						Comparision: &ComparisionExpression{
							Key:      "entity.kind",
							Operator: "EQUALS",
							Value:    "2.34",
						},
					},
					Operator: &logicalAnd,
					Right: &LogicalAndExpression{
						Left: &PrimaryExpression{
							Comparision: &ComparisionExpression{
								Key:      "entity.type",
								Operator: "EQUALS",
								Value:    "true",
							},
						},
						Operator: &logicalAnd,
						Right: &LogicalAndExpression{
							Left: &PrimaryExpression{
								UnaryOperation: &UnaryOperationExpression{
									Key:      "span.name",
									Operator: "NOT_EMPTY",
								},
							},
						},
					},
				},
				Operator: &logicalOr,
				Right: &LogicalOrExpression{
					Left: &LogicalAndExpression{
						Left: &PrimaryExpression{
							SubExpression: &LogicalOrExpression{
								Left: &LogicalAndExpression{
									Left: &PrimaryExpression{
										Comparision: &ComparisionExpression{
											Key:      "span.id",
											Operator: "NOT_EQUAL",
											Value:    "1234",
										},
									},
								},
								Operator: &logicalOr,
								Right: &LogicalOrExpression{
									Left: &LogicalAndExpression{
										Left: &PrimaryExpression{
											Comparision: &ComparisionExpression{
												Key:      "span.id",
												Operator: "NOT_EQUAL",
												Value:    "6789",
											},
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

	shouldSuccessfullyParseExpression(expression, expectedResult, t)
}

func TestShouldParseKeywordsCaseInsensitive(t *testing.T) {
	expression := "entity.name CONTAINS 'foo' and entity.type EQUALS 'bar'"

	logicalAnd := Operator("AND")
	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					Comparision: &ComparisionExpression{
						Key:      "entity.name",
						Operator: "CONTAINS",
						Value:    "foo",
					},
				},
				Operator: &logicalAnd,
				Right: &LogicalAndExpression{
					Left: &PrimaryExpression{
						Comparision: &ComparisionExpression{
							Key:      "entity.type",
							Operator: "EQUALS",
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
						Key:      "entity.name",
						Operator: "EQUALS",
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
						Key:      "entity.name",
						Operator: "NOT_EMPTY",
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
		t.Fatalf("Expected no error but got '%s'", err)
	}

	rendered := result.Render()
	if rendered != normalizedExpression {
		t.Fatalf("Expected to get normalized expression rendered but got: %s", rendered)
	}
}

func TestShouldRenderLogicalOrExpressionWhenOrIsSet(t *testing.T) {
	expectedResult := "foo EQUALS 'bar' OR foo CONTAINS 'bar'"

	logicalOr := Operator("OR")
	sut := LogicalOrExpression{
		Left: &LogicalAndExpression{
			Left: &PrimaryExpression{
				Comparision: &ComparisionExpression{
					Key:      "foo",
					Operator: "EQUALS",
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
						Operator: "CONTAINS",
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
	expectedResult := "foo EQUALS 'bar'"

	sut := LogicalOrExpression{
		Left: &LogicalAndExpression{
			Left: &PrimaryExpression{
				Comparision: &ComparisionExpression{
					Key:      "foo",
					Operator: "EQUALS",
					Value:    "bar",
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
	expectedResult := "foo EQUALS 'bar' AND foo CONTAINS 'bar'"

	logicalAnd := Operator("AND")
	sut := LogicalAndExpression{
		Left: &PrimaryExpression{
			Comparision: &ComparisionExpression{
				Key:      "foo",
				Operator: "EQUALS",
				Value:    "bar",
			},
		},
		Operator: &logicalAnd,
		Right: &LogicalAndExpression{
			Left: &PrimaryExpression{
				Comparision: &ComparisionExpression{
					Key:      "foo",
					Operator: "CONTAINS",
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
	expectedResult := "foo EQUALS 'bar'"

	sut := LogicalAndExpression{
		Left: &PrimaryExpression{
			Comparision: &ComparisionExpression{
				Key:      "foo",
				Operator: "EQUALS",
				Value:    "bar",
			},
		},
	}

	rendered := sut.Render()

	if rendered != expectedResult {
		t.Fatalf("Expected normalized rendered result of comparision expression but got:  %s", rendered)
	}
}

func TestShouldRenderComparisionOnPrimaryExpressionWhenComparsionIsSet(t *testing.T) {
	expectedResult := "foo EQUALS 'bar'"

	sut := PrimaryExpression{
		Comparision: &ComparisionExpression{
			Key:      "foo",
			Operator: "EQUALS",
			Value:    "bar",
		},
	}

	rendered := sut.Render()

	if rendered != expectedResult {
		t.Fatalf("Expected normalized rendered result of comparision expression but got:  %s", rendered)
	}
}

func TestShouldRenderUnaryOperationExpressionOnPrimaryExpressionWhenUnaryOperationIsSet(t *testing.T) {
	expectedResult := "foo IS_EMPTY"

	sut := PrimaryExpression{
		UnaryOperation: &UnaryOperationExpression{
			Key:      "foo",
			Operator: "IS_EMPTY",
		},
	}

	rendered := sut.Render()

	if rendered != expectedResult {
		t.Fatalf("Expected normalized rendered result of comparision expression but got:  %s", rendered)
	}
}

func TestShouldRenderSubExpression(t *testing.T) {
	expectedResult := "foo IS_EMPTY AND ( a EQUALS 'b' OR a EQUALS 'c' )"

	logicalOr := Operator("OR")
	logicalAnd := Operator("AND")
	sut := &LogicalAndExpression{
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
	}

	rendered := sut.Render()

	if rendered != expectedResult {
		t.Fatalf("Expected normalized rendered result of comparision expression but got:  %s", rendered)
	}
}
