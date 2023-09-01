package restapi_test

import (
	"go.uber.org/mock/gomock"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/stretchr/testify/require"
)

func TestShouldPrependElementToListOfElementsOfTagFilterExpression(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	element1 := mocks.NewMockTagFilterExpressionElement(ctrl)
	element2 := mocks.NewMockTagFilterExpressionElement(ctrl)
	element3 := mocks.NewMockTagFilterExpressionElement(ctrl)
	elements := []TagFilterExpressionElement{element1, element2}

	sut := NewLogicalAndTagFilter(elements)
	sut.PrependElement(element3)

	require.Len(t, sut.Elements, 3)
	require.Equal(t, element3, sut.Elements[2])
}

func TestShouldConvertTagFilterEntitiesToStringSlice(t *testing.T) {
	expectedResult := []string{"SOURCE", "DESTINATION", "NOT_APPLICABLE"}
	require.Equal(t, expectedResult, SupportedTagFilterEntities.ToStringSlice())
}
