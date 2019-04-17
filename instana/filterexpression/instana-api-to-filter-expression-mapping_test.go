package filterexpression_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	. "github.com/gessnerfl/terraform-provider-instana/instana/filterexpression"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

func TestShouldMapValidOperatorsOfTagExpression(t *testing.T) {
	for k, v := range OperatorMappingInstanaAPIToFilterExpression {
		t.Run(fmt.Sprintf("test mapping of %s t %s", k, v), testMappingOfOperatorsOfTagExpression(k, v))
	}
}

func testMappingOfOperatorsOfTagExpression(apiOperatorName string, filterExprOperatorName Operator) func(t *testing.T) {
	return func(t *testing.T) {
		key := "key"
		value := "value"
		input := restapi.TagMatcherExpression{
			Dtype:    restapi.LeafExpressionType,
			Key:      key,
			Operator: apiOperatorName,
			Value:    &value,
		}

		expectedResult := &FilterExpression{
			Expression: &LogicalOrExpression{
				Left: &LogicalAndExpression{
					Left: &PrimaryExpression{
						Comparision: &ComparisionExpression{
							Key:      key,
							Operator: filterExprOperatorName,
							Value:    value,
						},
					},
				},
			},
		}

		mapper := NewMapper()
		result, err := mapper.FromAPIModel(input)

		if err != nil {
			t.Fatalf("Expected no error but got %s", err)
		}
		if !cmp.Equal(result, expectedResult) {
			t.Fatalf("Expected parse expression %v but got %v; diff %s", expectedResult, result, cmp.Diff(expectedResult, result))
		}
	}
}
