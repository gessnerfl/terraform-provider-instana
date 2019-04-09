package utils_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/utils"
	"github.com/google/go-cmp/cmp"
)

func TestShouldSuccessfullyParseComplexExpression(t *testing.T) {
	expression := "entity.name CO \"foo\" AND entity.type EQ \"type\" AND ( NOT_EMPTY(span.name) OR span.id NE \"test\" )"

	expectedResult := Expression{
		Conjunction: &Conjunction{
			Left: &Expression{
				Comparision: &ComparisionExpression{
					Key:      "entity.name",
					Operator: "CO",
					Value:    "foo",
				},
			},
			Operator: "AND",
			Right: &Expression{
				Conjunction: &Conjunction{
					Left: &Expression{
						Comparision: &ComparisionExpression{
							Key:      "entity.type",
							Operator: "EQ",
							Value:    "my type",
						},
					},
					Operator: "AND",
					Right: &Expression{
						SubExpression: &Expression{
							Conjunction: &Conjunction{
								Left: &Expression{
									NotEmpty: &NotEmptyExpression{
										Key: "span.name",
									},
								},
								Operator: "OR",
								Right: &Expression{
									Comparision: &ComparisionExpression{
										Key:      "span.id",
										Operator: "NE",
										Value:    "test",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	sut := NewDynamicFocusFilter()
	result, err := sut.Parse(expression)

	if err != nil {
		t.Fatalf("Did not expected error but got %s", err)
	}

	if !cmp.Equal(expectedResult, result) {
		t.Fatalf("Expected parse expression %v but got %v; diff %s", expectedResult, result, cmp.Diff(expectedResult, result))
	}
}
