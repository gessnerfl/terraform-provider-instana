package restapi_test

import (
	"go.uber.org/mock/gomock"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
)

func TestShouldPrependElementToListOfElementsOfTagFilterExpression(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	element1 := NewStringTagFilter(TagFilterEntityDestination, "entity.type", EqualsOperator, "foo")
	element2 := NewStringTagFilter(TagFilterEntityDestination, "entity.type", EqualsOperator, "bar")
	element3 := NewStringTagFilter(TagFilterEntityDestination, "entity.type", EqualsOperator, "baz")
	elements := []*TagFilter{element1, element2}

	sut := NewLogicalAndTagFilter(elements)
	sut.PrependElement(element3)

	require.Len(t, sut.Elements, 3)
	require.Equal(t, element3, sut.Elements[2])
}

func TestShouldConvertTagFilterEntitiesToStringSlice(t *testing.T) {
	expectedResult := []string{"SOURCE", "DESTINATION", "NOT_APPLICABLE"}
	require.Equal(t, expectedResult, SupportedTagFilterEntities.ToStringSlice())
}
