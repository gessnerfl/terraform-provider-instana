package restapi_test

import (
	"errors"
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/golang/mock/gomock"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

const (
	idFieldValue         = "id"
	labelFieldValue      = "label"
	keyFieldValue        = "key"
	valueFieldValue      = "value"
	keyLeftFieldValue    = "keyLeft"
	keyRightFieldValue   = "keyRight"
	valueLeftFieldValue  = "valueLeft"
	valueRightFieldValue = "valueRight"

	operatorFieldName = "operator"

	messageLabelScope         = "scope"
	messageLabelBoundaryScope = "boundary scope"
)

func TestShouldSuccessfullyValidateConsistentApplicationConfigWithMatchSpecification(t *testing.T) {
	config := ApplicationConfig{
		ID:                 idFieldValue,
		Label:              labelFieldValue,
		MatchSpecification: NewComparisonExpression(keyFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueFieldValue),
		Scope:              ApplicationConfigScopeIncludeAllDownstream,
		BoundaryScope:      BoundaryScopeInbound,
	}

	err := config.Validate()

	assert.Nil(t, err)
	assert.Equal(t, idFieldValue, config.GetIDForResourcePath())
}

func TestShouldSuccessfullyValidateConsistentApplicationConfigWithTagFilter(t *testing.T) {
	value := valueFieldValue
	config := ApplicationConfig{
		ID:                  idFieldValue,
		Label:               labelFieldValue,
		TagFilterExpression: NewStringTagFilter(TagFilterEntityDestination, keyFieldValue, EqualsOperator, &value),
		Scope:               ApplicationConfigScopeIncludeAllDownstream,
		BoundaryScope:       BoundaryScopeInbound,
	}

	err := config.Validate()

	assert.Nil(t, err)
	assert.Equal(t, idFieldValue, config.GetIDForResourcePath())
}

func TestShouldFailToValidateApplicationConfigWhenIDIsMissing(t *testing.T) {
	config :=
		ApplicationConfig{
			Label:              labelFieldValue,
			MatchSpecification: NewComparisonExpression(keyFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueFieldValue),
			Scope:              ApplicationConfigScopeIncludeAllDownstream,
			BoundaryScope:      BoundaryScopeInbound,
		}

	err := config.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "id")
}

func TestShouldFailToValidateApplicationConfigWhenIDIsBlank(t *testing.T) {
	config :=
		ApplicationConfig{
			ID:                 " ",
			Label:              labelFieldValue,
			MatchSpecification: NewComparisonExpression(keyFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueFieldValue),
			Scope:              ApplicationConfigScopeIncludeAllDownstream,
			BoundaryScope:      BoundaryScopeInbound,
		}

	err := config.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "id")
}

func TestShouldFailToValidateApplicationConfigWhenLabelIsMissing(t *testing.T) {
	config :=
		ApplicationConfig{
			ID:                 idFieldValue,
			MatchSpecification: NewComparisonExpression(keyFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueFieldValue),
			Scope:              ApplicationConfigScopeIncludeAllDownstream,
			BoundaryScope:      BoundaryScopeInbound,
		}

	err := config.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "label")
}

func TestShouldFailToValidateApplicationConfigWhenLabelIsBlank(t *testing.T) {
	config :=
		ApplicationConfig{
			ID:                 idFieldValue,
			Label:              " ",
			MatchSpecification: NewComparisonExpression(keyFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueFieldValue),
			Scope:              ApplicationConfigScopeIncludeAllDownstream,
			BoundaryScope:      BoundaryScopeInbound,
		}

	err := config.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "label")
}

func TestShouldFailToValidateApplicationConfigWhenMatchSpecificationAndTagFilterIsMissing(t *testing.T) {
	config :=
		ApplicationConfig{
			ID:            idFieldValue,
			Label:         labelFieldValue,
			Scope:         ApplicationConfigScopeIncludeAllDownstream,
			BoundaryScope: BoundaryScopeInbound,
		}

	err := config.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "either match specification or tag filter")
}

func TestShouldFailToValidateApplicationConfigWhenMatchSpecificationAndTagFilterAreProvided(t *testing.T) {
	value := valueFieldValue
	config :=
		ApplicationConfig{
			ID:                  idFieldValue,
			Label:               labelFieldValue,
			Scope:               ApplicationConfigScopeIncludeAllDownstream,
			BoundaryScope:       BoundaryScopeInbound,
			MatchSpecification:  NewComparisonExpression(keyFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueFieldValue),
			TagFilterExpression: NewStringTagFilter(TagFilterEntityDestination, keyFieldValue, EqualsOperator, &value),
		}

	err := config.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "either match specification or tag filter")
}

func TestShouldFailToValidateApplicationConfigWhenMatchSpecificationIsNotValid(t *testing.T) {
	config := ApplicationConfig{
		ID:                 idFieldValue,
		Label:              labelFieldValue,
		MatchSpecification: NewComparisonExpression("", MatcherExpressionEntityDestination, EqualsOperator, valueFieldValue),
		Scope:              ApplicationConfigScopeIncludeAllDownstream,
		BoundaryScope:      BoundaryScopeInbound,
	}

	err := config.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "key")
}

func TestShouldFailToValidateApplicationConfigWhenTagFilterExpressionIsNotValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedError := errors.New("test")
	tagFilter := mocks.NewMockTagFilterExpressionElement(ctrl)
	tagFilter.EXPECT().Validate().Times(1).Return(expectedError)

	config := ApplicationConfig{
		ID:                  idFieldValue,
		Label:               labelFieldValue,
		TagFilterExpression: tagFilter,
		Scope:               ApplicationConfigScopeIncludeAllDownstream,
		BoundaryScope:       BoundaryScopeInbound,
	}

	err := config.Validate()

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}

func TestShouldFailToValidateApplicationConfigWhenScopeIsMissing(t *testing.T) {
	config :=
		ApplicationConfig{
			ID:                 idFieldValue,
			Label:              labelFieldValue,
			MatchSpecification: NewComparisonExpression(keyFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueFieldValue),
			BoundaryScope:      BoundaryScopeInbound,
		}

	err := config.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), messageLabelScope)
}

func TestShouldFailToValidateApplicationConfigWhenScopeIsBlank(t *testing.T) {
	config :=
		ApplicationConfig{
			ID:                 idFieldValue,
			Label:              labelFieldValue,
			MatchSpecification: NewComparisonExpression(keyFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueFieldValue),
			Scope:              " ",
			BoundaryScope:      BoundaryScopeInbound,
		}

	err := config.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), messageLabelScope)
}

func TestShouldFailToValidateApplicationConfigWhenScopeIsNotSupported(t *testing.T) {
	config :=
		ApplicationConfig{
			ID:                 idFieldValue,
			Label:              labelFieldValue,
			MatchSpecification: NewComparisonExpression(keyFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueFieldValue),
			Scope:              "invalid",
			BoundaryScope:      BoundaryScopeInbound,
		}

	err := config.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), messageLabelScope)
}

func TestShouldFailToValidateApplicationConfigWhenBoundaryScopeIsMissing(t *testing.T) {
	config :=
		ApplicationConfig{
			ID:                 idFieldValue,
			Label:              labelFieldValue,
			MatchSpecification: NewComparisonExpression(keyFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueFieldValue),
			Scope:              ApplicationConfigScopeIncludeAllDownstream,
		}

	err := config.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), messageLabelBoundaryScope)
}

func TestShouldFailToValidateApplicationConfigWhenBoundaryScopeIsBlank(t *testing.T) {
	config :=
		ApplicationConfig{
			ID:                 idFieldValue,
			Label:              labelFieldValue,
			MatchSpecification: NewComparisonExpression(keyFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueFieldValue),
			Scope:              ApplicationConfigScopeIncludeAllDownstream,
			BoundaryScope:      " ",
		}

	err := config.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), messageLabelBoundaryScope)
}

func TestShouldFailToValidateApplicationConfigWhenBoundaryScopeIsNotSupported(t *testing.T) {
	config :=
		ApplicationConfig{
			ID:                 idFieldValue,
			Label:              labelFieldValue,
			MatchSpecification: NewComparisonExpression(keyFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueFieldValue),
			Scope:              ApplicationConfigScopeIncludeAllDownstream,
			BoundaryScope:      "invalid",
		}

	err := config.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), messageLabelBoundaryScope)
}

func TestShouldSuccessfullyValidateConsistentBinaryExpression(t *testing.T) {
	for _, operator := range SupportedLogicalOperatorTypes {
		t.Run(fmt.Sprintf("TestShouldSuccessfullyValidateConsistentBinaryExpressionOfType%s", operator), createTestShouldSuccessfullyValidateConsistentBinaryExpression(operator))
	}
}

func createTestShouldSuccessfullyValidateConsistentBinaryExpression(operator LogicalOperatorType) func(t *testing.T) {
	return func(t *testing.T) {
		left := NewComparisonExpression(keyLeftFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueLeftFieldValue)
		right := NewComparisonExpression(keyRightFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueRightFieldValue)

		exp := NewBinaryOperator(left, operator, right)

		err := exp.Validate()

		assert.Nil(t, err)
		assert.Equal(t, BinaryOperatorExpressionType, exp.GetType())
	}
}

func TestShouldFailToValidateBinaryExpressionWhenLeftOperatorIsMissing(t *testing.T) {
	right := NewComparisonExpression(keyRightFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueRightFieldValue)

	exp := NewBinaryOperator(nil, LogicalAnd, right)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "left")
}

func TestShouldFailToValidateBinaryExpressionWhenLeftOperatorIsNotValid(t *testing.T) {
	left := NewComparisonExpression("", MatcherExpressionEntityDestination, EqualsOperator, valueLeftFieldValue)
	right := NewComparisonExpression(keyRightFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueRightFieldValue)

	exp := NewBinaryOperator(left, LogicalAnd, right)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "key")
}

func TestShouldFailToValidateBinaryExpressionWhenConjunctionIsNotValid(t *testing.T) {
	left := NewComparisonExpression("leftKey", MatcherExpressionEntityDestination, EqualsOperator, valueLeftFieldValue)
	right := NewComparisonExpression(keyRightFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueRightFieldValue)

	exp := NewBinaryOperator(left, "FOO", right)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "conjunction")
}

func TestShouldFailToValidateBinaryExpressionWhenRightOperatorIsMissing(t *testing.T) {
	left := NewComparisonExpression(keyLeftFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueLeftFieldValue)

	exp := NewBinaryOperator(left, LogicalAnd, nil)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "right")
}

func TestShouldFailToValidateBinaryExpressionWhenRightOperatorIsNotValid(t *testing.T) {
	left := NewComparisonExpression(keyLeftFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueLeftFieldValue)
	right := NewComparisonExpression("", MatcherExpressionEntityDestination, EqualsOperator, valueRightFieldValue)

	exp := NewBinaryOperator(left, LogicalAnd, right)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "key")
}

func TestShouldFailToValidateBinaryExpressionWhenConjunctionIsMissing(t *testing.T) {
	left := NewComparisonExpression(keyLeftFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueLeftFieldValue)
	right := NewComparisonExpression(keyRightFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueRightFieldValue)

	exp := NewBinaryOperator(left, "", right)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "conjunction")
}

func TestShouldCreateValidComparisionExpression(t *testing.T) {
	for _, operator := range SupportedComparisonOperators {
		t.Run(fmt.Sprintf("TestShouldCreateValidComparisionExpressionOfType%s", operator), createTestShouldCreateValidComparisionExpression(operator))
	}
}

func createTestShouldCreateValidComparisionExpression(operator TagFilterOperator) func(*testing.T) {
	return func(t *testing.T) {
		exp := NewComparisonExpression(keyFieldValue, MatcherExpressionEntityDestination, operator, valueFieldValue)

		err := exp.Validate()

		assert.Nil(t, err)
		assert.Equal(t, LeafExpressionType, exp.GetType())
	}
}

func TestShouldFailToValidateComparisionExpressionWhenKeyIsMissing(t *testing.T) {
	exp := NewComparisonExpression("", MatcherExpressionEntityDestination, EqualsOperator, valueFieldValue)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "key")
}

func TestShouldFailToValidateComparisionExpressionWhenEntityIsNotValid(t *testing.T) {
	exp := NewComparisonExpression(keyFieldValue, MatcherExpressionEntity("invalid"), EqualsOperator, valueFieldValue)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "entity")
}

func TestShouldFailToValidateComparisionExpressionWhenOperatorIsMissing(t *testing.T) {
	exp := NewComparisonExpression(keyFieldValue, MatcherExpressionEntityDestination, "", valueFieldValue)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), operatorFieldName)
}

func TestShouldFailToValidateComparisionExpressionWhenOperatorIsNotValid(t *testing.T) {
	exp := NewComparisonExpression(keyFieldValue, MatcherExpressionEntityDestination, "FOO", valueFieldValue)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), operatorFieldName)
}

func TestShouldFailToValidateComparisionExpressionWhenValueIsMissing(t *testing.T) {
	exp := NewComparisonExpression(keyFieldValue, MatcherExpressionEntityDestination, EqualsOperator, "")

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "value")
}

func TestShouldCreateValidUnaryOperatorExpression(t *testing.T) {
	for _, operator := range SupportedUnaryExpressionOperators {
		t.Run(fmt.Sprintf("TestShouldCreateValidUnaryOperatorExpressionOfType%s", operator), createTestShouldCreateValidUnaryOperatorExpression(operator))
	}
}

func createTestShouldCreateValidUnaryOperatorExpression(operator TagFilterOperator) func(*testing.T) {
	return func(t *testing.T) {
		exp := NewUnaryOperationExpression(keyFieldValue, MatcherExpressionEntityDestination, operator)

		err := exp.Validate()

		assert.Nil(t, err)
		assert.Equal(t, LeafExpressionType, exp.GetType())
	}
}

func TestShouldFailToValidateUnaryOperatorExpressionWhenKeyIsMissing(t *testing.T) {
	exp := NewUnaryOperationExpression("", MatcherExpressionEntityDestination, IsEmptyOperator)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "key")
}

func TestShouldFailToValidateUnaryOperatorExpressionWhenEntityIsNotValid(t *testing.T) {
	exp := NewUnaryOperationExpression(keyFieldValue, MatcherExpressionEntity("invalid"), IsEmptyOperator)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "entity")
}

func TestShouldFailToValidateUnaryOperatorExpressionWhenOperatorIsMissing(t *testing.T) {
	exp := NewUnaryOperationExpression(keyFieldValue, MatcherExpressionEntityDestination, "")

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), operatorFieldName)
}

func TestShouldFailToValidateUnaryOperatorExpressionWhenOperatorIsNotValid(t *testing.T) {
	exp := NewUnaryOperationExpression(keyFieldValue, MatcherExpressionEntityDestination, "FOO")

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), operatorFieldName)
}

func TestShouldFailToValidateUnaryOperatorExpressionWhenValueIsSet(t *testing.T) {
	exp := NewComparisonExpression(keyFieldValue, MatcherExpressionEntityDestination, IsEmptyOperator, "")

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "value")
}

func TestShouldReturnTrueForAllSupportedApplicationConfigScopes(t *testing.T) {
	for _, scope := range SupportedApplicationConfigScopes {
		t.Run(fmt.Sprintf("TestShouldReturnTrueForSupportedApplicationConfigScope%s", string(scope)), createTestCaseToVerifySupportedApplicationConfigScope(scope))
	}
}

func createTestCaseToVerifySupportedApplicationConfigScope(scope ApplicationConfigScope) func(t *testing.T) {
	return func(t *testing.T) {
		assert.True(t, SupportedApplicationConfigScopes.IsSupported(scope))
	}
}

func TestShouldReturnfalseWhenApplicationConfigScopeIsNotSupported(t *testing.T) {
	assert.False(t, SupportedApplicationConfigScopes.IsSupported(ApplicationConfigScope(valueInvalid)))
}

func TestShouldReturnStringRepresentationOfSupporedApplicationConfigScopes(t *testing.T) {
	assert.Equal(t, []string{"INCLUDE_NO_DOWNSTREAM", "INCLUDE_IMMEDIATE_DOWNSTREAM_DATABASE_AND_MESSAGING", "INCLUDE_ALL_DOWNSTREAM"}, SupportedApplicationConfigScopes.ToStringSlice())
}

func TestShouldReturnTrueForAllSupportedApplicationConfigBoundaryScopes(t *testing.T) {
	for _, scope := range SupportedBoundaryScopes {
		t.Run(fmt.Sprintf("TestShouldReturnTrueForSupportedBoundaryScope%s", string(scope)), createTestCaseToVerifySupportedApplicationConfigBoundaryScope(scope))
	}
}

func createTestCaseToVerifySupportedApplicationConfigBoundaryScope(scope BoundaryScope) func(t *testing.T) {
	return func(t *testing.T) {
		assert.True(t, SupportedBoundaryScopes.IsSupported(scope))
	}
}

func TestShouldReturnfalseWhenApplicationConfigBoundaryScopeIsNotSupported(t *testing.T) {
	assert.False(t, SupportedBoundaryScopes.IsSupported(BoundaryScope(valueInvalid)))
}

func TestShouldReturnStringRepresentationOfSupporedApplicationConfigBoundaryScopes(t *testing.T) {
	assert.Equal(t, []string{"ALL", "INBOUND", "DEFAULT"}, SupportedBoundaryScopes.ToStringSlice())
}

func TestShouldReturnTrueForAllSupportedMatcherExpressionEntities(t *testing.T) {
	for _, entity := range SupportedMatcherExpressionEntities {
		t.Run(fmt.Sprintf("TestShouldReturnTrueForSupportedMatcherExpressionEntity%s", string(entity)), createTestCaseToVerifySupportedMatcherExpressionEntity(entity))
	}
}

func createTestCaseToVerifySupportedMatcherExpressionEntity(entity MatcherExpressionEntity) func(t *testing.T) {
	return func(t *testing.T) {
		assert.True(t, SupportedMatcherExpressionEntities.IsSupported(entity))
	}
}

func TestShouldReturnfalseWhenMatcherExpressionEntityIsNotSupported(t *testing.T) {
	assert.False(t, SupportedMatcherExpressionEntities.IsSupported(MatcherExpressionEntity(valueInvalid)))
}

func TestShouldReturnStringRepresentationOfSupporedMatcherExpressionEntities(t *testing.T) {
	assert.Equal(t, []string{"SOURCE", "DESTINATION", "NOT_APPLICABLE"}, SupportedMatcherExpressionEntities.ToStringSlice())
}
