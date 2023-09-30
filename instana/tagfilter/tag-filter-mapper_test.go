package tagfilter_test

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	. "github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestShouldReturnNilWhenMappingEmptyTagFilterExpressionToNormalizedString(t *testing.T) {
	operatorType := restapi.LogicalOr
	input := &restapi.TagFilter{
		Type:            restapi.TagFilterExpressionType,
		LogicalOperator: &operatorType,
	}

	result, err := MapTagFilterToNormalizedString(input)

	require.NoError(t, err)
	require.Nil(t, result)
}

func TestShouldReturnErrorWhenMappingAnInvalidTagFilterExpressionToNormalizedString(t *testing.T) {
	input := &restapi.TagFilter{
		Type: restapi.TagFilterExpressionElementType("invalid"),
	}

	result, err := MapTagFilterToNormalizedString(input)

	require.Error(t, err)
	require.Nil(t, result)
}

func TestShouldReturnStringWhenMappingAValidTagFilterExpressionToNormalizedString(t *testing.T) {
	value := int64(1234)
	input := restapi.NewNumberTagFilter(restapi.TagFilterEntityDestination, tagFilterName, restapi.EqualsOperator, value)

	result, err := MapTagFilterToNormalizedString(input)

	require.NoError(t, err)
	require.Equal(t, "name@dest EQUALS 1234", *result)
}
