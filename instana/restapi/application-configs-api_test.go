package restapi_test

import (
	"fmt"
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

func TestShouldSuccussullyValididateConsistentApplicationConfig(t *testing.T) {
	config := ApplicationConfig{
		ID:                 idFieldValue,
		Label:              labelFieldValue,
		MatchSpecification: NewComparisionExpression(keyFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueFieldValue),
		Scope:              ApplicationConfigScopeIncludeAllDownstream,
		BoundaryScope:      BoundaryScopeInbound,
	}

	err := config.Validate()

	assert.Nil(t, err)
	assert.Equal(t, idFieldValue, config.GetID())
}

func TestShouldFailToValidateApplicationConfigWhenIDIsMissing(t *testing.T) {
	config :=
		ApplicationConfig{
			Label:              labelFieldValue,
			MatchSpecification: NewComparisionExpression(keyFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueFieldValue),
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
			MatchSpecification: NewComparisionExpression(keyFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueFieldValue),
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
			MatchSpecification: NewComparisionExpression(keyFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueFieldValue),
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
			MatchSpecification: NewComparisionExpression(keyFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueFieldValue),
			Scope:              ApplicationConfigScopeIncludeAllDownstream,
			BoundaryScope:      BoundaryScopeInbound,
		}

	err := config.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "label")
}

func TestShouldFailToValidateApplicationConfigWhenMatchSpecificationIsMissing(t *testing.T) {
	config :=
		ApplicationConfig{
			ID:            idFieldValue,
			Label:         labelFieldValue,
			Scope:         ApplicationConfigScopeIncludeAllDownstream,
			BoundaryScope: BoundaryScopeInbound,
		}

	err := config.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "match specification")
}

func TestShouldFailToValidateApplicationConfigWhenMatchSpecificationIsNotValid(t *testing.T) {
	config := ApplicationConfig{
		ID:                 idFieldValue,
		Label:              labelFieldValue,
		MatchSpecification: NewComparisionExpression("", MatcherExpressionEntityDestination, EqualsOperator, valueFieldValue),
		Scope:              ApplicationConfigScopeIncludeAllDownstream,
		BoundaryScope:      BoundaryScopeInbound,
	}

	err := config.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "key")
}

func TestShouldFailToValidateApplicationConfigWhenScopeIsMissing(t *testing.T) {
	config :=
		ApplicationConfig{
			ID:                 idFieldValue,
			Label:              labelFieldValue,
			MatchSpecification: NewComparisionExpression(keyFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueFieldValue),
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
			MatchSpecification: NewComparisionExpression(keyFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueFieldValue),
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
			MatchSpecification: NewComparisionExpression(keyFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueFieldValue),
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
			MatchSpecification: NewComparisionExpression(keyFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueFieldValue),
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
			MatchSpecification: NewComparisionExpression(keyFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueFieldValue),
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
			MatchSpecification: NewComparisionExpression(keyFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueFieldValue),
			Scope:              ApplicationConfigScopeIncludeAllDownstream,
			BoundaryScope:      "invalid",
		}

	err := config.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), messageLabelBoundaryScope)
}

func TestShouldSuccessfullyValidateConsistentBinaryExpression(t *testing.T) {
	for _, operator := range SupportedConjunctionTypes {
		t.Run(fmt.Sprintf("TestShouldSuccessfullyValidateConsistentBinaryExpressionOfType%s", operator), createTestShouldSuccessfullyValidateConsistentBinaryExpression(operator))
	}
}

func createTestShouldSuccessfullyValidateConsistentBinaryExpression(operator ConjunctionType) func(t *testing.T) {
	return func(t *testing.T) {
		left := NewComparisionExpression(keyLeftFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueLeftFieldValue)
		right := NewComparisionExpression(keyRightFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueRightFieldValue)

		exp := NewBinaryOperator(left, operator, right)

		err := exp.Validate()

		assert.Nil(t, err)
		assert.Equal(t, BinaryOperatorExpressionType, exp.GetType())
	}
}

func TestShouldFailToValidateBinaryExpressionWhenLeftOperatorIsMissing(t *testing.T) {
	right := NewComparisionExpression(keyRightFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueRightFieldValue)

	exp := NewBinaryOperator(nil, LogicalAnd, right)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "left")
}

func TestShouldFailToValidateBinaryExpressionWhenLeftOperatorIsNotValid(t *testing.T) {
	left := NewComparisionExpression("", MatcherExpressionEntityDestination, EqualsOperator, valueLeftFieldValue)
	right := NewComparisionExpression(keyRightFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueRightFieldValue)

	exp := NewBinaryOperator(left, LogicalAnd, right)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "key")
}

func TestShouldFailToValidateBinaryExpressionWhenConjunctionIsNotValid(t *testing.T) {
	left := NewComparisionExpression("leftKey", MatcherExpressionEntityDestination, EqualsOperator, valueLeftFieldValue)
	right := NewComparisionExpression(keyRightFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueRightFieldValue)

	exp := NewBinaryOperator(left, "FOO", right)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "conjunction")
}

func TestShouldFailToValidateBinaryExpressionWhenRightOperatorIsMissing(t *testing.T) {
	left := NewComparisionExpression(keyLeftFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueLeftFieldValue)

	exp := NewBinaryOperator(left, LogicalAnd, nil)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "right")
}

func TestShouldFailToValidateBinaryExpressionWhenRightOperatorIsNotValid(t *testing.T) {
	left := NewComparisionExpression(keyLeftFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueLeftFieldValue)
	right := NewComparisionExpression("", MatcherExpressionEntityDestination, EqualsOperator, valueRightFieldValue)

	exp := NewBinaryOperator(left, LogicalAnd, right)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "key")
}

func TestShouldFailToValidateBinaryExpressionWhenConjunctionIsMissing(t *testing.T) {
	left := NewComparisionExpression(keyLeftFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueLeftFieldValue)
	right := NewComparisionExpression(keyRightFieldValue, MatcherExpressionEntityDestination, EqualsOperator, valueRightFieldValue)

	exp := NewBinaryOperator(left, "", right)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "conjunction")
}

func TestShouldCreateValidComparisionExpression(t *testing.T) {
	for _, operator := range SupportedComparisionOperators {
		t.Run(fmt.Sprintf("TestShouldCreateValidComparisionExpressionOfType%s", operator), createTestShouldCreateValidComparisionExpression(operator))
	}
}

func createTestShouldCreateValidComparisionExpression(operator MatcherOperator) func(*testing.T) {
	return func(t *testing.T) {
		exp := NewComparisionExpression(keyFieldValue, MatcherExpressionEntityDestination, operator, valueFieldValue)

		err := exp.Validate()

		assert.Nil(t, err)
		assert.Equal(t, LeafExpressionType, exp.GetType())
	}
}

func TestShouldFailToValidateComparisionExpressionWhenKeyIsMissing(t *testing.T) {
	exp := NewComparisionExpression("", MatcherExpressionEntityDestination, EqualsOperator, valueFieldValue)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "key")
}

func TestShouldFailToValidateComparisionExpressionWhenOperatorIsMissing(t *testing.T) {
	exp := NewComparisionExpression(keyFieldValue, MatcherExpressionEntityDestination, "", valueFieldValue)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), operatorFieldName)
}

func TestShouldFailToValidateComparisionExpressionWhenOperatorIsNotValid(t *testing.T) {
	exp := NewComparisionExpression(keyFieldValue, MatcherExpressionEntityDestination, "FOO", valueFieldValue)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), operatorFieldName)
}

func TestShouldFailToValidateComparisionExpressionWhenValueIsMissing(t *testing.T) {
	exp := NewComparisionExpression(keyFieldValue, MatcherExpressionEntityDestination, EqualsOperator, "")

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "value")
}

func TestShouldCreateValidUnaryOperatorExpression(t *testing.T) {
	for _, operator := range SupportedUnaryExpressionOperators {
		t.Run(fmt.Sprintf("TestShouldCreateValidUnaryOperatorExpressionOfType%s", operator), createTestShouldCreateValidUnaryOperatorExpression(operator))
	}
}

func createTestShouldCreateValidUnaryOperatorExpression(operator MatcherOperator) func(*testing.T) {
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
	exp := NewComparisionExpression(keyFieldValue, MatcherExpressionEntityDestination, IsEmptyOperator, "")

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
