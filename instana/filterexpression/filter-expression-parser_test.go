package filterexpression_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/filterexpression"
	"github.com/google/go-cmp/cmp"
)

func TestShouldSuccessfullyParseComplexExpression(t *testing.T) {
	expression := "entity.name CO 'foo bar' OR entity.kind EQ '2.34' AND entity.type EQ 'true' AND span.name NOT EMPTY OR ( span.id NE  '1234' OR span.id NE '6789' )"

	logicalAnd := Operator("AND")
	logicalOr := Operator("OR")
	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					Comparision: &ComparisionExpression{
						Key:      "entity.name",
						Operator: "CO",
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
							Operator: "EQ",
							Value:    "2.34",
						},
					},
					Operator: &logicalAnd,
					Right: &LogicalAndExpression{
						Left: &PrimaryExpression{
							Comparision: &ComparisionExpression{
								Key:      "entity.type",
								Operator: "EQ",
								Value:    "true",
							},
						},
						Operator: &logicalAnd,
						Right: &LogicalAndExpression{
							Left: &PrimaryExpression{
								UnaryOperation: &UnaryOperationExpression{
									Key:      "span.name",
									Operator: "NOT EMPTY",
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
											Operator: "NE",
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
												Operator: "NE",
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
	expression := "entity.name co 'foo' and entity.type EQ 'bar'"

	logicalAnd := Operator("AND")
	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					Comparision: &ComparisionExpression{
						Key:      "entity.name",
						Operator: "CO",
						Value:    "foo",
					},
				},
				Operator: &logicalAnd,
				Right: &LogicalAndExpression{
					Left: &PrimaryExpression{
						Comparision: &ComparisionExpression{
							Key:      "entity.type",
							Operator: "EQ",
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
	expression := "entity.name eQ 'foo'"

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &PrimaryExpression{
					Comparision: &ComparisionExpression{
						Key:      "entity.name",
						Operator: "EQ",
						Value:    "foo",
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
					UnaryOperation: &UnaryOperationExpression{
						Key:      "entity.name",
						Operator: "NOT EMPTY",
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
	expression := "entity.name co 'foo' OR entity.kind EQ '2.34'    and  entity.type EQ 'true'  AND span.name  NOT empTy   OR span.id  NE  '1234'"
	normalizedExpression := "entity.name CO 'foo' OR entity.kind EQ '2.34' AND entity.type EQ 'true' AND span.name NOT EMPTY OR span.id NE '1234'"

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
	expectedResult := "foo EQ 'bar' OR foo CO 'bar'"

	logicalOr := Operator("OR")
	sut := LogicalOrExpression{
		Left: &LogicalAndExpression{
			Left: &PrimaryExpression{
				Comparision: &ComparisionExpression{
					Key:      "foo",
					Operator: "EQ",
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
						Operator: "CO",
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
	expectedResult := "foo EQ 'bar'"

	sut := LogicalOrExpression{
		Left: &LogicalAndExpression{
			Left: &PrimaryExpression{
				Comparision: &ComparisionExpression{
					Key:      "foo",
					Operator: "EQ",
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
	expectedResult := "foo EQ 'bar' AND foo CO 'bar'"

	logicalAnd := Operator("AND")
	sut := LogicalAndExpression{
		Left: &PrimaryExpression{
			Comparision: &ComparisionExpression{
				Key:      "foo",
				Operator: "EQ",
				Value:    "bar",
			},
		},
		Operator: &logicalAnd,
		Right: &LogicalAndExpression{
			Left: &PrimaryExpression{
				Comparision: &ComparisionExpression{
					Key:      "foo",
					Operator: "CO",
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
	expectedResult := "foo EQ 'bar'"

	sut := LogicalAndExpression{
		Left: &PrimaryExpression{
			Comparision: &ComparisionExpression{
				Key:      "foo",
				Operator: "EQ",
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
	expectedResult := "foo EQ 'bar'"

	sut := PrimaryExpression{
		Comparision: &ComparisionExpression{
			Key:      "foo",
			Operator: "EQ",
			Value:    "bar",
		},
	}

	rendered := sut.Render()

	if rendered != expectedResult {
		t.Fatalf("Expected normalized rendered result of comparision expression but got:  %s", rendered)
	}
}

func TestShouldRenderUnaryOperationExpressionOnPrimaryExpressionWhenUnaryOperationIsSet(t *testing.T) {
	expectedResult := "foo IS EMPTY"

	sut := PrimaryExpression{
		UnaryOperation: &UnaryOperationExpression{
			Key:      "foo",
			Operator: "IS EMPTY",
		},
	}

	rendered := sut.Render()

	if rendered != expectedResult {
		t.Fatalf("Expected normalized rendered result of comparision expression but got:  %s", rendered)
	}
}

func TestShouldRenderSubExpression(t *testing.T) {
	expectedResult := "foo IS EMPTY AND ( a EQ 'b' OR a EQ 'c' )"

	logicalOr := Operator("OR")
	logicalAnd := Operator("AND")
	sut := &LogicalAndExpression{
		Left: &PrimaryExpression{
			UnaryOperation: &UnaryOperationExpression{
				Key:      "foo",
				Operator: "IS EMPTY",
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
								Operator: "EQ",
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
									Operator: "EQ",
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
