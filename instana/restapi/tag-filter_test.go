package restapi_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/stretchr/testify/require"
)

const (
	tagFilterEntityName           = "entity.name"
	errorTagFilterOperatorInvalid = "tag filter operator INVALID"
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
	require.Contains(t, err.Error(), errorTagFilterOperatorInvalid)
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

func TestShouldCreateValidStringTagFilter(t *testing.T) {
	value := "test"

	sut := NewStringTagFilter(TagFilterEntityDestination, tagFilterEntityName, EqualsOperator, value)

	require.Equal(t, TagFilterType, sut.Type)
	require.Equal(t, TagFilterType, sut.GetType())
	require.Equal(t, tagFilterEntityName, sut.Name)
	require.Equal(t, TagFilterEntityDestination, sut.Entity)
	require.Equal(t, EqualsOperator, sut.Operator)
	require.Equal(t, &value, sut.StringValue)
	require.Nil(t, sut.NumberValue)
	require.Nil(t, sut.BooleanValue)
	require.Nil(t, sut.Key)
	require.Equal(t, value, sut.Value)
	require.NoError(t, sut.Validate())
}

func TestShouldCreateValidNumberTagFilter(t *testing.T) {
	value := int64(1234)

	sut := NewNumberTagFilter(TagFilterEntityDestination, tagFilterEntityName, EqualsOperator, value)

	require.Equal(t, TagFilterType, sut.Type)
	require.Equal(t, TagFilterType, sut.GetType())
	require.Equal(t, tagFilterEntityName, sut.Name)
	require.Equal(t, TagFilterEntityDestination, sut.Entity)
	require.Equal(t, EqualsOperator, sut.Operator)
	require.Equal(t, &value, sut.NumberValue)
	require.Nil(t, sut.StringValue)
	require.Nil(t, sut.BooleanValue)
	require.Nil(t, sut.Key)
	require.Equal(t, value, sut.Value)
	require.NoError(t, sut.Validate())
}

func TestShouldCreateValidBooleanTagFilter(t *testing.T) {
	value := true

	sut := NewBooleanTagFilter(TagFilterEntityDestination, tagFilterEntityName, EqualsOperator, value)

	require.Equal(t, TagFilterType, sut.Type)
	require.Equal(t, TagFilterType, sut.GetType())
	require.Equal(t, tagFilterEntityName, sut.Name)
	require.Equal(t, TagFilterEntityDestination, sut.Entity)
	require.Equal(t, EqualsOperator, sut.Operator)
	require.Equal(t, &value, sut.BooleanValue)
	require.Nil(t, sut.StringValue)
	require.Nil(t, sut.NumberValue)
	require.Nil(t, sut.Key)
	require.Equal(t, value, sut.Value)
	require.NoError(t, sut.Validate())
}

func TestShouldCreateValidTagTagFilter(t *testing.T) {
	key := "key"
	value := "value"
	fullString := "key=value"

	sut := NewTagTagFilter(TagFilterEntityDestination, tagFilterEntityName, EqualsOperator, key, value)

	require.Equal(t, TagFilterType, sut.Type)
	require.Equal(t, TagFilterType, sut.GetType())
	require.Equal(t, tagFilterEntityName, sut.Name)
	require.Equal(t, TagFilterEntityDestination, sut.Entity)
	require.Equal(t, EqualsOperator, sut.Operator)
	require.Equal(t, &key, sut.Key)
	require.Equal(t, value, sut.Value)
	require.Equal(t, &fullString, sut.StringValue)
	require.Nil(t, sut.NumberValue)
	require.Nil(t, sut.BooleanValue)
	require.NoError(t, sut.Validate())
}

func TestShouldCreateValidUnaryTagFilter(t *testing.T) {
	sut := NewUnaryTagFilter(TagFilterEntityDestination, tagFilterEntityName, IsEmptyOperator)

	require.Equal(t, TagFilterType, sut.Type)
	require.Equal(t, TagFilterType, sut.GetType())
	require.Equal(t, tagFilterEntityName, sut.Name)
	require.Equal(t, TagFilterEntityDestination, sut.Entity)
	require.Equal(t, IsEmptyOperator, sut.Operator)
	require.Nil(t, sut.StringValue)
	require.Nil(t, sut.NumberValue)
	require.Nil(t, sut.BooleanValue)
	require.Nil(t, sut.Key)
	require.Nil(t, sut.Value)
	require.NoError(t, sut.Validate())
}

func TestShouldReturnNoErrorWhenValidatingTagFilterWithASupportedEntityType(t *testing.T) {
	for _, entity := range SupportedTagFilterEntities {
		sut := NewUnaryTagFilter(entity, tagFilterEntityName, IsEmptyOperator)
		require.NoError(t, sut.Validate())
	}
}

func TestShouldReturnErrorWhenValidatingTagFilterWithAnUnsupportedEntityType(t *testing.T) {
	sut := NewUnaryTagFilter("INVALID", tagFilterEntityName, IsEmptyOperator)

	err := sut.Validate()

	require.Error(t, err)
	require.Contains(t, err.Error(), "tag filter entity type INVALID")
}

func TestShouldReturnErrorWhenValidatingTagFilterWithoutName(t *testing.T) {
	sut := NewUnaryTagFilter(TagFilterEntityDestination, "", IsEmptyOperator)

	err := sut.Validate()

	require.Error(t, err)
	require.Contains(t, err.Error(), "tag filter name")
}

func TestShouldReturnNoErrorWhenValidatingUnaryTagFilterWithASupportedOperator(t *testing.T) {
	for _, op := range SupportedUnaryExpressionOperators {
		sut := NewUnaryTagFilter(TagFilterEntityDestination, tagFilterEntityName, op)
		require.NoError(t, sut.Validate())
	}
}

func TestShouldReturnErrorWhenValidatingUnaryTagFilterWithAnInvalidOperator(t *testing.T) {
	sut := NewUnaryTagFilter(TagFilterEntityDestination, tagFilterEntityName, "INVALID")

	err := sut.Validate()

	require.Error(t, err)
	require.Contains(t, err.Error(), errorTagFilterOperatorInvalid)
}

func TestShouldReturnErrorWhenValidatingUnaryTagFilterWithAStringValueAssigned(t *testing.T) {
	value := "value"
	testUnaryOperatorHasNoValueAssigned(t, func(sut *TagFilter) { sut.StringValue = &value })
}

func TestShouldReturnErrorWhenValidatingUnaryTagFilterWithANumberValueAssigned(t *testing.T) {
	value := int64(1234)
	testUnaryOperatorHasNoValueAssigned(t, func(sut *TagFilter) { sut.NumberValue = &value })
}

func TestShouldReturnErrorWhenValidatingUnaryTagFilterWithABooleanValueAssigned(t *testing.T) {
	value := true
	testUnaryOperatorHasNoValueAssigned(t, func(sut *TagFilter) { sut.BooleanValue = &value })
}

func TestShouldReturnErrorWhenValidatingUnaryTagFilterWithATagKeyAssigned(t *testing.T) {
	key := "key"
	testUnaryOperatorHasNoValueAssigned(t, func(sut *TagFilter) { sut.Key = &key })
}

func TestShouldReturnErrorWhenValidatingUnaryTagFilterWithATagValueAssigned(t *testing.T) {
	value := "value"
	testUnaryOperatorHasNoValueAssigned(t, func(sut *TagFilter) { sut.Value = &value })
}

func testUnaryOperatorHasNoValueAssigned(t *testing.T, valueSetter func(sut *TagFilter)) {
	sut := NewUnaryTagFilter(TagFilterEntityDestination, tagFilterEntityName, IsEmptyOperator)
	valueSetter(sut)

	err := sut.Validate()

	require.Error(t, err)
	require.Contains(t, err.Error(), "no value must be assigned")
}

func TestShouldReturnNoErrorWhenValidatingComparisonTagFilterWithASupportedOperator(t *testing.T) {
	value := "value"
	for _, op := range SupportedComparisonOperators {
		sut := NewStringTagFilter(TagFilterEntityDestination, tagFilterEntityName, op, value)
		require.NoError(t, sut.Validate())
	}
}

func TestShouldReturnErrorWhenValidatingComparisonTagFilterWithAnInvalidOperator(t *testing.T) {
	value := "value"

	sut := NewStringTagFilter(TagFilterEntityDestination, tagFilterEntityName, "INVALID", value)

	err := sut.Validate()

	require.Error(t, err)
	require.Contains(t, err.Error(), errorTagFilterOperatorInvalid)
}

func TestShouldReturnErrorWhenValidatingComparisonTagFilterWithoutValue(t *testing.T) {
	sut := &TagFilter{
		Entity:      TagFilterEntityDestination,
		Name:        tagFilterEntityName,
		Operator:    EqualsOperator,
		Value:       nil,
		StringValue: nil,
	}

	err := sut.Validate()

	require.Error(t, err)
	require.Contains(t, err.Error(), "value missing for comparison")
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
