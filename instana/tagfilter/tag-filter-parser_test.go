package tagfilter_test

import (
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
)

const (
	keyEntityName = "entity.name"
	keyEntityKind = "entity.kind"
	keyEntityType = "entity.type"
	keyAgentTags  = "agent.tag"

	valueMyValue = "my value"

	entityNameEqualsValueExpression = "entity.name@dest EQUALS 'my value'"
)

func TestShouldParseStringComparisonExpression(t *testing.T) {
	expression := "entity.name EQUALS 'foo'"
	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:      &EntitySpec{Identifier: keyEntityName},
							Operator:    Operator(restapi.EqualsOperator),
							StringValue: utils.StringPtr("foo"),
						},
					},
				},
			},
		},
	}

	shouldSuccessfullyParseExpression(expression, expectedResult, t)
}

func TestShouldParseNumberComparisonExpression(t *testing.T) {
	expression := "entity.name EQUALS 123"
	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:      &EntitySpec{Identifier: keyEntityName},
							Operator:    Operator(restapi.EqualsOperator),
							NumberValue: utils.Int64Ptr(int64(123)),
						},
					},
				},
			},
		},
	}

	shouldSuccessfullyParseExpression(expression, expectedResult, t)
}

func TestShouldParseBoolComparisonExpression(t *testing.T) {
	expression := "entity.name EQUALS TRUE"
	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:       &EntitySpec{Identifier: keyEntityName},
							Operator:     Operator(restapi.EqualsOperator),
							BooleanValue: utils.BoolPtr(true),
						},
					},
				},
			},
		},
	}

	shouldSuccessfullyParseExpression(expression, expectedResult, t)
}

func TestShouldParseStringTagComparisonExpression(t *testing.T) {
	expression := "agent.tag:key EQUALS 'value'"
	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:      &EntitySpec{Identifier: keyAgentTags, TagKey: utils.StringPtr("key")},
							Operator:    Operator(restapi.EqualsOperator),
							StringValue: utils.StringPtr("value"),
						},
					},
				},
			},
		},
	}

	shouldSuccessfullyParseExpression(expression, expectedResult, t)
}

func TestShouldParseNumberTagComparisonExpression(t *testing.T) {
	expression := "agent.tag:key EQUALS 1234"
	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:      &EntitySpec{Identifier: keyAgentTags, TagKey: utils.StringPtr("key")},
							Operator:    Operator(restapi.EqualsOperator),
							NumberValue: utils.Int64Ptr(1234),
						},
					},
				},
			},
		},
	}

	shouldSuccessfullyParseExpression(expression, expectedResult, t)
}

func TestShouldParseBooleanTagComparisonExpression(t *testing.T) {
	expression := "agent.tag:key EQUALS true"
	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:       &EntitySpec{Identifier: keyAgentTags, TagKey: utils.StringPtr("key")},
							Operator:     Operator(restapi.EqualsOperator),
							BooleanValue: utils.BoolPtr(true),
						},
					},
				},
			},
		},
	}

	shouldSuccessfullyParseExpression(expression, expectedResult, t)
}

func TestShouldHandlePrecedenceWhenParseExpressions(t *testing.T) {
	logicalAnd := Operator(restapi.LogicalAnd)
	logicalOr := Operator(restapi.LogicalOr)
	expression := "entity.name EQUALS 'foo' OR entity.name EQUALS 'bar' AND agent.tag:key EQUALS 'value'"
	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:      &EntitySpec{Identifier: keyEntityName},
							Operator:    Operator(restapi.EqualsOperator),
							StringValue: utils.StringPtr("foo"),
						},
					},
				},
			},
			Operator: &logicalOr,
			Right: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &BracketExpression{
						Primary: &PrimaryExpression{
							Comparison: &ComparisonExpression{
								Entity:      &EntitySpec{Identifier: keyEntityName},
								Operator:    Operator(restapi.EqualsOperator),
								StringValue: utils.StringPtr("bar"),
							},
						},
					},
					Operator: &logicalAnd,
					Right: &LogicalAndExpression{
						Left: &BracketExpression{
							Primary: &PrimaryExpression{
								Comparison: &ComparisonExpression{
									Entity:      &EntitySpec{Identifier: keyAgentTags, TagKey: utils.StringPtr("key")},
									Operator:    Operator(restapi.EqualsOperator),
									StringValue: utils.StringPtr("value"),
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

func TestShouldParseExpressionStartingWithABracketExpression(t *testing.T) {
	logicalAnd := Operator(restapi.LogicalAnd)
	logicalOr := Operator(restapi.LogicalOr)
	expression := "( entity.name EQUALS 'foo' OR entity.name EQUALS 'bar' ) AND agent.tag:key EQUALS 'value'"
	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Bracket: &LogicalOrExpression{
						Left: &LogicalAndExpression{
							Left: &BracketExpression{
								Primary: &PrimaryExpression{
									Comparison: &ComparisonExpression{
										Entity:      &EntitySpec{Identifier: keyEntityName},
										Operator:    Operator(restapi.EqualsOperator),
										StringValue: utils.StringPtr("foo"),
									},
								},
							},
						},
						Operator: &logicalOr,
						Right: &LogicalOrExpression{
							Left: &LogicalAndExpression{
								Left: &BracketExpression{
									Primary: &PrimaryExpression{
										Comparison: &ComparisonExpression{
											Entity:      &EntitySpec{Identifier: keyEntityName},
											Operator:    Operator(restapi.EqualsOperator),
											StringValue: utils.StringPtr("bar"),
										},
									},
								},
							},
						},
					},
				},
				Operator: &logicalAnd,
				Right: &LogicalAndExpression{
					Left: &BracketExpression{
						Primary: &PrimaryExpression{
							Comparison: &ComparisonExpression{
								Entity:      &EntitySpec{Identifier: keyAgentTags, TagKey: utils.StringPtr("key")},
								Operator:    Operator(restapi.EqualsOperator),
								StringValue: utils.StringPtr("value"),
							},
						},
					},
				},
			},
		},
	}

	shouldSuccessfullyParseExpression(expression, expectedResult, t)
}

func TestShouldParseExpressionEndingWithABracketExpression(t *testing.T) {
	logicalAnd := Operator(restapi.LogicalAnd)
	logicalOr := Operator(restapi.LogicalOr)
	expression := "agent.tag:key EQUALS 'value' AND ( entity.name EQUALS 'foo' OR entity.name EQUALS 'bar' )"
	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:      &EntitySpec{Identifier: keyAgentTags, TagKey: utils.StringPtr("key")},
							Operator:    Operator(restapi.EqualsOperator),
							StringValue: utils.StringPtr("value"),
						},
					},
				},
				Operator: &logicalAnd,
				Right: &LogicalAndExpression{
					Left: &BracketExpression{
						Bracket: &LogicalOrExpression{
							Left: &LogicalAndExpression{
								Left: &BracketExpression{
									Primary: &PrimaryExpression{
										Comparison: &ComparisonExpression{
											Entity:      &EntitySpec{Identifier: keyEntityName},
											Operator:    Operator(restapi.EqualsOperator),
											StringValue: utils.StringPtr("foo"),
										},
									},
								},
							},
							Operator: &logicalOr,
							Right: &LogicalOrExpression{
								Left: &LogicalAndExpression{
									Left: &BracketExpression{
										Primary: &PrimaryExpression{
											Comparison: &ComparisonExpression{
												Entity:      &EntitySpec{Identifier: keyEntityName},
												Operator:    Operator(restapi.EqualsOperator),
												StringValue: utils.StringPtr("bar"),
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

func TestShouldSuccessfullyParseComplexExpression(t *testing.T) {
	expression := "entity.name CONTAINS 'foo bar' OR entity.kind EQUALS 234 AND entity.type EQUALS true AND ( span.name NOT_EMPTY OR span.id@src NOT_EQUAL  '1234' )"

	logicalAnd := Operator(restapi.LogicalAnd)
	logicalOr := Operator(restapi.LogicalOr)
	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:      &EntitySpec{Identifier: keyEntityName},
							Operator:    Operator(restapi.ContainsOperator),
							StringValue: utils.StringPtr("foo bar"),
						},
					},
				},
			},
			Operator: &logicalOr,
			Right: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &BracketExpression{
						Primary: &PrimaryExpression{
							Comparison: &ComparisonExpression{
								Entity:      &EntitySpec{Identifier: keyEntityKind},
								Operator:    Operator(restapi.EqualsOperator),
								NumberValue: utils.Int64Ptr(int64(234)),
							},
						},
					},
					Operator: &logicalAnd,
					Right: &LogicalAndExpression{
						Left: &BracketExpression{
							Primary: &PrimaryExpression{
								Comparison: &ComparisonExpression{
									Entity:       &EntitySpec{Identifier: keyEntityType},
									Operator:     Operator(restapi.EqualsOperator),
									BooleanValue: utils.BoolPtr(true),
								},
							},
						},
						Operator: &logicalAnd,
						Right: &LogicalAndExpression{
							Left: &BracketExpression{
								Bracket: &LogicalOrExpression{
									Left: &LogicalAndExpression{
										Left: &BracketExpression{
											Primary: &PrimaryExpression{
												UnaryOperation: &UnaryOperationExpression{
													Entity:   &EntitySpec{Identifier: "span.name"},
													Operator: Operator(restapi.NotEmptyOperator),
												},
											},
										},
									},
									Operator: &logicalOr,
									Right: &LogicalOrExpression{
										Left: &LogicalAndExpression{
											Left: &BracketExpression{
												Primary: &PrimaryExpression{
													Comparison: &ComparisonExpression{
														Entity:      &EntitySpec{Identifier: "span.id", Origin: utils.StringPtr(EntityOriginSource.Key())},
														Operator:    Operator(restapi.NotEqualOperator),
														StringValue: utils.StringPtr("1234"),
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
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:      &EntitySpec{Identifier: keyEntityName},
							Operator:    Operator(restapi.ContainsOperator),
							StringValue: utils.StringPtr("foo"),
						},
					},
				},
				Operator: &logicalAnd,
				Right: &LogicalAndExpression{
					Left: &BracketExpression{
						Primary: &PrimaryExpression{
							Comparison: &ComparisonExpression{
								Entity:      &EntitySpec{Identifier: keyEntityType},
								Operator:    Operator(restapi.EqualsOperator),
								StringValue: utils.StringPtr("bar"),
							},
						},
					},
				},
			},
		},
	}

	shouldSuccessfullyParseExpression(expression, expectedResult, t)
}

func TestShouldParseAllSupportedComparisonOperators(t *testing.T) {
	for _, o := range restapi.SupportedComparisonOperators {
		t.Run(fmt.Sprintf("TestShouldParseComparisionOperator%s", string(o)), createTestCaseForParsingSupportedComparisonOperators(o))
	}
}

func createTestCaseForParsingSupportedComparisonOperators(operator restapi.ExpressionOperator) func(*testing.T) {
	return func(t *testing.T) {
		expression := fmt.Sprintf("entity.name %s 'foo'", string(operator))

		expectedResult := &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &BracketExpression{
						Primary: &PrimaryExpression{
							Comparison: &ComparisonExpression{
								Entity:      &EntitySpec{Identifier: keyEntityName},
								Operator:    Operator(operator),
								StringValue: utils.StringPtr("foo"),
							},
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

func createTestCaseForParsingSupportedUnaryOperators(operator restapi.ExpressionOperator) func(*testing.T) {
	return func(t *testing.T) {
		expression := fmt.Sprintf("entity.name %s", string(operator))

		expectedResult := &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &BracketExpression{
						Primary: &PrimaryExpression{
							UnaryOperation: &UnaryOperationExpression{
								Entity:   &EntitySpec{Identifier: keyEntityName},
								Operator: Operator(operator),
							},
						},
					},
				},
			},
		}

		shouldSuccessfullyParseExpression(expression, expectedResult, t)
	}
}

func TestShouldParseFullySpecifiedUnaryTagFilter(t *testing.T) {
	//expression := "agent.tag:stage@dest NOT_EMPTY AND service.name@src EQUALS 'foo'"
	expression := "agent.tag:stage@dest NOT_EMPTY"

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Entity:   &EntitySpec{Identifier: keyAgentTags, TagKey: utils.StringPtr("stage"), Origin: utils.StringPtr(EntityOriginDestination.Key())},
							Operator: Operator(restapi.NotEmptyOperator),
						},
					},
				},
			},
		},
	}

	shouldSuccessfullyParseExpression(expression, expectedResult, t)
}

func TestShouldParseComparisonOperationsCaseInsensitive(t *testing.T) {
	expression := "entity.name Equals 'foo'"

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:      &EntitySpec{Identifier: keyEntityName},
							Operator:    Operator(restapi.EqualsOperator),
							StringValue: utils.StringPtr("foo"),
						},
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
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Entity:   &EntitySpec{Identifier: keyEntityName},
							Operator: Operator(restapi.NotEmptyOperator),
						},
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
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:      &EntitySpec{Identifier: "call.http.header.x-example-foo"},
							Operator:    Operator(restapi.EqualsOperator),
							StringValue: utils.StringPtr("test"),
						},
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
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:      &EntitySpec{Identifier: "kubernetes.pod.label.foo/bar"},
							Operator:    Operator(restapi.EqualsOperator),
							StringValue: utils.StringPtr("test"),
						},
					},
				},
			},
		},
	}

	shouldSuccessfullyParseExpression(expression, expectedResult, t)
}

func TestShouldParseEntityOriginFromComparisonExpression(t *testing.T) {
	expression := "entity.name@src EQUALS 'test'"

	expectedResult := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:      &EntitySpec{Identifier: keyEntityName, Origin: utils.StringPtr(EntityOriginSource.Key())},
							Operator:    Operator(restapi.EqualsOperator),
							StringValue: utils.StringPtr("test"),
						},
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
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Entity:   &EntitySpec{Identifier: keyEntityName, Origin: utils.StringPtr(EntityOriginSource.Key())},
							Operator: Operator(restapi.NotEmptyOperator),
						},
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

func TestShouldRenderComplexExpressionInNormalizedForm(t *testing.T) {
	expression := "entity.name CONTAINS 'foo' OR entity.kind EQUALS '2.34'    and  entity.type EQUALS 'true'  AND ( span.name  NOT_EMPTY   OR span.id  NOT_EQUAL  '1234' )"
	normalizedExpression := "entity.name@dest CONTAINS 'foo' OR entity.kind@dest EQUALS '2.34' AND entity.type@dest EQUALS 'true' AND ( span.name@dest NOT_EMPTY OR span.id@dest NOT_EQUAL '1234' )"

	sut := NewParser()
	result, err := sut.Parse(expression)
	require.Nil(t, err)

	rendered := result.Render()
	require.Equal(t, normalizedExpression, rendered)
}

func TestShouldRenderLogicalOrExpression(t *testing.T) {
	expectedResult := "foo@dest EQUALS 'bar' OR foo@dest CONTAINS 'bar'"

	logicalOr := Operator(restapi.LogicalOr)
	sut := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:      &EntitySpec{Identifier: "foo", Origin: utils.StringPtr(EntityOriginDestination.Key())},
							Operator:    Operator(restapi.EqualsOperator),
							StringValue: utils.StringPtr("bar"),
						},
					},
				},
			},
			Operator: &logicalOr,
			Right: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &BracketExpression{
						Primary: &PrimaryExpression{
							Comparison: &ComparisonExpression{
								Entity:      &EntitySpec{Identifier: "foo", Origin: utils.StringPtr(EntityOriginDestination.Key())},
								Operator:    Operator(restapi.ContainsOperator),
								StringValue: utils.StringPtr("bar"),
							},
						},
					},
				},
			},
		},
	}

	rendered := sut.Render()

	require.Equal(t, expectedResult, rendered)
}

func TestShouldRenderLogicalAndExpression(t *testing.T) {
	expectedResult := "foo@dest EQUALS 'bar' AND foo@dest CONTAINS 'bar'"

	logicalAnd := Operator(restapi.LogicalAnd)
	sut := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:      &EntitySpec{Identifier: "foo", Origin: utils.StringPtr(EntityOriginDestination.Key())},
							Operator:    Operator(restapi.EqualsOperator),
							StringValue: utils.StringPtr("bar"),
						},
					},
				},
				Operator: &logicalAnd,
				Right: &LogicalAndExpression{
					Left: &BracketExpression{
						Primary: &PrimaryExpression{
							Comparison: &ComparisonExpression{
								Entity:      &EntitySpec{Identifier: "foo", Origin: utils.StringPtr(EntityOriginDestination.Key())},
								Operator:    Operator(restapi.ContainsOperator),
								StringValue: utils.StringPtr("bar"),
							},
						},
					},
				},
			},
		},
	}

	rendered := sut.Render()
	require.Equal(t, expectedResult, rendered)
}

func TestShouldRenderBracketExpression(t *testing.T) {
	expectedResult := "( foo@dest EQUALS 'bar' OR foo@dest CONTAINS 'bar' ) AND bar@dest EQUALS 'value'"

	logicalOr := Operator(restapi.LogicalAnd)
	logicalAnd := Operator(restapi.LogicalAnd)
	sut := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Bracket: &LogicalOrExpression{
						Left: &LogicalAndExpression{
							Left: &BracketExpression{
								Primary: &PrimaryExpression{
									Comparison: &ComparisonExpression{
										Entity:      &EntitySpec{Identifier: "foo", Origin: utils.StringPtr(EntityOriginDestination.Key())},
										Operator:    Operator(restapi.EqualsOperator),
										StringValue: utils.StringPtr("bar"),
									},
								},
							},
						},
						Operator: &logicalOr,
						Right: &LogicalOrExpression{
							Left: &LogicalAndExpression{
								Left: &BracketExpression{
									Primary: &PrimaryExpression{
										Comparison: &ComparisonExpression{
											Entity:      &EntitySpec{Identifier: "foo", Origin: utils.StringPtr(EntityOriginDestination.Key())},
											Operator:    Operator(restapi.ContainsOperator),
											StringValue: utils.StringPtr("bar"),
										},
									},
								},
							},
						},
					},
				},
				Operator: &logicalAnd,
				Right: &LogicalAndExpression{
					Left: &BracketExpression{
						Primary: &PrimaryExpression{
							Comparison: &ComparisonExpression{
								Entity:      &EntitySpec{Identifier: "bar", Origin: utils.StringPtr(EntityOriginDestination.Key())},
								Operator:    Operator(restapi.EqualsOperator),
								StringValue: utils.StringPtr("value"),
							},
						},
					},
				},
			},
		},
	}

	rendered := sut.Render()
	require.Equal(t, expectedResult, rendered)
}

func TestShouldRenderPrimaryStringComparisonExpression(t *testing.T) {
	sut := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:      &EntitySpec{Identifier: keyEntityName, Origin: utils.StringPtr(EntityOriginDestination.Key())},
							Operator:    Operator(restapi.EqualsOperator),
							StringValue: utils.StringPtr(valueMyValue),
						},
					},
				},
			},
		},
	}

	rendered := sut.Render()

	require.Equal(t, entityNameEqualsValueExpression, rendered)
}

func TestShouldRenderPrimaryNumberComparisonExpression(t *testing.T) {
	sut := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:      &EntitySpec{Identifier: keyEntityName, Origin: utils.StringPtr(EntityOriginDestination.Key())},
							Operator:    Operator(restapi.EqualsOperator),
							NumberValue: utils.Int64Ptr(int64(1234)),
						},
					},
				},
			},
		},
	}

	rendered := sut.Render()

	require.Equal(t, "entity.name@dest EQUALS 1234", rendered)
}

func TestShouldRenderPrimaryBooleanComparisonExpression(t *testing.T) {
	sut := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:       &EntitySpec{Identifier: keyEntityName, Origin: utils.StringPtr(EntityOriginDestination.Key())},
							Operator:     Operator(restapi.EqualsOperator),
							BooleanValue: utils.BoolPtr(true),
						},
					},
				},
			},
		},
	}

	rendered := sut.Render()

	require.Equal(t, "entity.name@dest EQUALS true", rendered)
}

func TestShouldRenderPrimaryTagComparisonExpression(t *testing.T) {
	sut := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						Comparison: &ComparisonExpression{
							Entity:      &EntitySpec{Identifier: keyEntityName, TagKey: utils.StringPtr("key"), Origin: utils.StringPtr(EntityOriginDestination.Key())},
							Operator:    Operator(restapi.EqualsOperator),
							StringValue: utils.StringPtr("value"),
						},
					},
				},
			},
		},
	}

	rendered := sut.Render()

	require.Equal(t, "entity.name:key@dest EQUALS 'value'", rendered)
}

func TestShouldRenderUnaryOperationExpression(t *testing.T) {
	expectedResult := "foo@dest IS_EMPTY"

	sut := &FilterExpression{
		Expression: &LogicalOrExpression{
			Left: &LogicalAndExpression{
				Left: &BracketExpression{
					Primary: &PrimaryExpression{
						UnaryOperation: &UnaryOperationExpression{
							Entity:   &EntitySpec{Identifier: "foo", Origin: utils.StringPtr(EntityOriginDestination.Key())},
							Operator: Operator(restapi.IsEmptyOperator),
						},
					},
				},
			},
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
	for _, e := range restapi.SupportedTagFilterEntities {
		t.Run(fmt.Sprintf("TestShouldGetEntityOriginForInstanaAPIEntity%s", e), func(t *testing.T) {
			require.Equal(t, e, SupportedEntityOrigins.ForInstanaAPIEntity(e).TagFilterEntity())
		})
	}
}

func TestShouldReturnEntityOriginDestinationAsFallbackValueWhenMatcherExpressionEntityIsNotValid(t *testing.T) {
	require.Equal(t, EntityOriginDestination, SupportedEntityOrigins.ForInstanaAPIEntity("invalid"))
}

func TestShouldNormalizeExpressionAndProduceTheSameResultWhenPerformedMultipleTimesOnOutput(t *testing.T) {
	testSets := []normalizationTestSet{
		{
			name:     "BasicComparison",
			input:    "entity.name    NOT_EMPTY",
			expected: "entity.name@dest NOT_EMPTY",
		},
		{
			name:     "AndWithBracketedOr",
			input:    "agent.tag:key EQUALS 'value' AND ( entity.name EQUALS 'foo' OR entity.name EQUALS 'bar' )",
			expected: "( agent.tag:key@dest EQUALS 'value' AND ( entity.name@dest EQUALS 'foo' OR entity.name@dest EQUALS 'bar' ) )"},
		{
			name:     "BracketedOrWithAnd",
			input:    "( entity.name EQUALS 'foo' OR entity.name EQUALS 'bar' ) AND agent.tag:key EQUALS 'value'",
			expected: "( ( entity.name@dest EQUALS 'foo' OR entity.name@dest EQUALS 'bar' ) AND agent.tag:key@dest EQUALS 'value' )",
		},
		{
			name:     "OrWithBracketedAnd",
			input:    "agent.tag:key EQUALS 'value' OR ( entity.name EQUALS 'foo' AND entity.name EQUALS 'bar' )",
			expected: "( agent.tag:key@dest EQUALS 'value' OR ( entity.name@dest EQUALS 'foo' AND entity.name@dest EQUALS 'bar' ) )",
		},
		{
			name:     "BracketedAndWithOr",
			input:    "( entity.name EQUALS 'foo' AND entity.name EQUALS 'bar' ) OR agent.tag:key EQUALS 'value'",
			expected: "( ( entity.name@dest EQUALS 'foo' AND entity.name@dest EQUALS 'bar' ) OR agent.tag:key@dest EQUALS 'value' )",
		},
		{
			name:     "AndWithBracketedAnd",
			input:    "agent.tag:key EQUALS 'value' AND ( entity.name EQUALS 'foo' AND entity.name EQUALS 'bar' )",
			expected: "( agent.tag:key@dest EQUALS 'value' AND entity.name@dest EQUALS 'foo' AND entity.name@dest EQUALS 'bar' )",
		},
		{
			name:     "BracketedAndWithAnd",
			input:    "( entity.name EQUALS 'foo' AND entity.name EQUALS 'bar' ) AND agent.tag:key EQUALS 'value'",
			expected: "( entity.name@dest EQUALS 'foo' AND entity.name@dest EQUALS 'bar' AND agent.tag:key@dest EQUALS 'value' )",
		},
		{
			name:     "OrWithBracketedOr",
			input:    "agent.tag:key EQUALS 'value' OR ( entity.name EQUALS 'foo' OR entity.name EQUALS 'bar' )",
			expected: "( agent.tag:key@dest EQUALS 'value' OR entity.name@dest EQUALS 'foo' OR entity.name@dest EQUALS 'bar' )",
		},
		{
			name:     "BracketedOrWithOr",
			input:    "( entity.name EQUALS 'foo' OR entity.name EQUALS 'bar' ) OR agent.tag:key EQUALS 'value'",
			expected: "( entity.name@dest EQUALS 'foo' OR entity.name@dest EQUALS 'bar' OR agent.tag:key@dest EQUALS 'value' )",
		},
	}

	for _, s := range testSets {
		input := s.input
		for i := 0; i < 3; i++ {
			t.Run(fmt.Sprintf("TestShouldNormalizeExpression%s%d", s.name, i), func(t *testing.T) {

				result, err := Normalize(input)
				require.NoError(t, err)
				require.Equal(t, s.expected, result)

				input = result
			})
		}
	}
}

func TestShouldFailToNormalizeExpressionWhenExpressionIsNotValid(t *testing.T) {
	input := "entity.name    bla bla bla"

	result, err := Normalize(input)
	require.Error(t, err)
	require.Equal(t, input, result)
}

type normalizationTestSet struct {
	name     string
	input    string
	expected string
}
