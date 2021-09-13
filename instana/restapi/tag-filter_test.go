package restapi_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/stretchr/testify/require"
)

func TestShouldCreateValidLogicalOrTagFilterExpression(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	element := mocks.NewMockTagFilterExpressionElement(ctrl)
	element.EXPECT().Validate().Times(2).Return(nil)
	elements := []TagFilterExpressionElement{element, element}
	sut := NewLogicalOrTagFilter(elements)

	require.Equal(t, TagFilterExpressionType, sut.Type)
	require.Equal(t, TagFilterExpressionType, sut.GetType())
	require.Equal(t, LogicalOr, sut.LogicalOperator)
	require.Equal(t, elements, sut.Elements)
	require.NoError(t, sut.Validate())
}

func TestShouldCreateValidLogicalAndTagFilterExpression(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	element := mocks.NewMockTagFilterExpressionElement(ctrl)
	element.EXPECT().Validate().Times(2).Return(nil)
	elements := []TagFilterExpressionElement{element, element}
	sut := NewLogicalAndTagFilter(elements)

	require.Equal(t, TagFilterExpressionType, sut.Type)
	require.Equal(t, TagFilterExpressionType, sut.GetType())
	require.Equal(t, LogicalAnd, sut.LogicalOperator)
	require.Equal(t, elements, sut.Elements)
	require.NoError(t, sut.Validate())
}

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

func TestShouldReturnErrorWhenValidatingTagFilterExpressionWithLessThanTwoElements(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	element1 := mocks.NewMockTagFilterExpressionElement(ctrl)
	elements := []TagFilterExpressionElement{element1}

	sut := NewLogicalAndTagFilter(elements)

	err := sut.Validate()

	require.Error(t, err)
	require.Contains(t, err.Error(), "at least two elements")
}

func TestShouldReturnErrorWhenValidatingTagFilterExpressionAndLogicalOperatorIsNotSupported(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	element1 := mocks.NewMockTagFilterExpressionElement(ctrl)
	elements := []TagFilterExpressionElement{element1, element1}

	sut := NewLogicalAndTagFilter(elements)
	sut.LogicalOperator = "INVALID"

	err := sut.Validate()

	require.Error(t, err)
	require.Contains(t, err.Error(), "tag filter operator INVALID")
}

func TestShouldReturnErrorWhenValidatingTagFilterExpressionAndTagFilterTypeIsNotExpression(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	element1 := mocks.NewMockTagFilterExpressionElement(ctrl)
	elements := []TagFilterExpressionElement{element1, element1}

	sut := NewLogicalAndTagFilter(elements)
	sut.Type = "INVALID"

	err := sut.Validate()

	require.Error(t, err)
	require.Contains(t, err.Error(), "tag filter expression must be of type EXPRESSION")
}

func TestShouldReturnErrorWhenValidatingTagFilterExpressionAndAtLeastOnElementReturnsAnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedError := errors.New("test")
	element1 := mocks.NewMockTagFilterExpressionElement(ctrl)
	element1.EXPECT().Validate().Times(1).Return(nil)
	element2 := mocks.NewMockTagFilterExpressionElement(ctrl)
	element2.EXPECT().Validate().Times(1).Return(expectedError)
	elements := []TagFilterExpressionElement{element1, element2}

	sut := NewLogicalAndTagFilter(elements)

	err := sut.Validate()

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

func TestShouldReturnTrueForAllSupportedLogicalOperatorTypes(t *testing.T) {
	for _, v := range SupportedLogicalOperatorTypes {
		require.True(t, SupportedLogicalOperatorTypes.IsSupported(v))
	}
}

func TestShouldReturnFalseForAllNonSupportedLogicalOperatorTypes(t *testing.T) {
	for _, v := range []string{"FOO", "BAR", "INVALID"} {
		require.False(t, SupportedLogicalOperatorTypes.IsSupported(LogicalOperatorType(v)))
	}
}

func TestShouldReturnTrueForAllSupportedTagFilterEntityTypes(t *testing.T) {
	for _, v := range SupportedTagFilterEntities {
		require.True(t, SupportedTagFilterEntities.IsSupported(v))
	}
}

func TestShouldReturnFalseForAllNonSupportedTagFilterEntityTypes(t *testing.T) {
	for _, v := range []string{"FOO", "BAR", "INVALID"} {
		require.False(t, SupportedTagFilterEntities.IsSupported(TagFilterEntity(v)))
	}
}

func TestShouldConvertTagFilterEntitiesToStringSlice(t *testing.T) {
	expectedResult := []string{"SOURCE", "DESTINATION", "NOT_APPLICABLE"}
	require.Equal(t, expectedResult, SupportedTagFilterEntities.ToStringSlice())
}

func TestShouldReturnTrueForAllSupportedComparisonOperators(t *testing.T) {
	for _, v := range SupportedComparisonOperators {
		require.True(t, SupportedComparisonOperators.IsSupported(v))
	}
}

func TestShouldReturnFalseForAllNonSupportedComparisonOperators(t *testing.T) {
	for _, v := range append(SupportedUnaryExpressionOperators, "INVALID_OPERATOR") {
		require.False(t, SupportedComparisonOperators.IsSupported(v))
	}
}

func TestShouldReturnTrueForAllSupportedUnaryExpressionOperators(t *testing.T) {
	for _, v := range SupportedUnaryExpressionOperators {
		require.True(t, SupportedUnaryExpressionOperators.IsSupported(v))
	}
}

func TestShouldReturnFalseForAllNonSupportedUnaryExpressionOperators(t *testing.T) {
	for _, v := range append(SupportedComparisonOperators, "INVALID_OPERATOR") {
		require.False(t, SupportedUnaryExpressionOperators.IsSupported(v))
	}
}
