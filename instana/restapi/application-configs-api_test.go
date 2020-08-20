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
	scopeFieldValue      = "scope"
	keyFieldValue        = "key"
	valueFieldValue      = "value"
	keyLeftFieldValue    = "keyLeft"
	keyRightFieldValue   = "keyRight"
	valueLeftFieldValue  = "valueLeft"
	valueRightFieldValue = "valueRight"

	operatorFieldName = "operator"
)

func TestShouldSuccussullyValididateConsistentApplicationConfig(t *testing.T) {
	config := ApplicationConfig{
		ID:                 idFieldValue,
		Label:              labelFieldValue,
		MatchSpecification: NewComparisionExpression(keyFieldValue, EqualsOperator, valueFieldValue),
		Scope:              scopeFieldValue,
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
			MatchSpecification: NewComparisionExpression(keyFieldValue, EqualsOperator, valueFieldValue),
			Scope:              scopeFieldValue,
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
			MatchSpecification: NewComparisionExpression(keyFieldValue, EqualsOperator, valueFieldValue),
			Scope:              scopeFieldValue,
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
			MatchSpecification: NewComparisionExpression(keyFieldValue, EqualsOperator, valueFieldValue),
			Scope:              scopeFieldValue,
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
			MatchSpecification: NewComparisionExpression(keyFieldValue, EqualsOperator, valueFieldValue),
			Scope:              scopeFieldValue,
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
			Scope:         scopeFieldValue,
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
		MatchSpecification: NewComparisionExpression("", EqualsOperator, valueFieldValue),
		Scope:              scopeFieldValue,
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
			MatchSpecification: NewComparisionExpression(keyFieldValue, EqualsOperator, valueFieldValue),
			BoundaryScope:      BoundaryScopeInbound,
		}

	err := config.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "scope")
}

func TestShouldFailToValidateApplicationConfigWhenScopeIsBlank(t *testing.T) {
	config :=
		ApplicationConfig{
			ID:                 idFieldValue,
			Label:              labelFieldValue,
			MatchSpecification: NewComparisionExpression(keyFieldValue, EqualsOperator, valueFieldValue),
			Scope:              " ",
			BoundaryScope:      BoundaryScopeInbound,
		}

	err := config.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "scope")
}

func TestShouldFailToValidateApplicationConfigWhenBoundaryScopeIsMissing(t *testing.T) {
	config :=
		ApplicationConfig{
			ID:                 idFieldValue,
			Label:              labelFieldValue,
			MatchSpecification: NewComparisionExpression(keyFieldValue, EqualsOperator, valueFieldValue),
			Scope:              scopeFieldValue,
		}

	err := config.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "boundary scope")
}

func TestShouldFailToValidateApplicationConfigWhenBoundaryScopeIsBlank(t *testing.T) {
	config :=
		ApplicationConfig{
			ID:                 idFieldValue,
			Label:              labelFieldValue,
			MatchSpecification: NewComparisionExpression(keyFieldValue, EqualsOperator, valueFieldValue),
			Scope:              scopeFieldValue,
			BoundaryScope:      " ",
		}

	err := config.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "boundary scope")
}

func TestShouldSuccessfullyValidateConsistentBinaryExpression(t *testing.T) {
	for _, operator := range SupportedConjunctionTypes {
		t.Run(fmt.Sprintf("TestShouldSuccessfullyValidateConsistentBinaryExpressionOfType%s", operator), createTestShouldSuccessfullyValidateConsistentBinaryExpression(operator))
	}
}

func createTestShouldSuccessfullyValidateConsistentBinaryExpression(operator ConjunctionType) func(t *testing.T) {
	return func(t *testing.T) {
		left := NewComparisionExpression(keyLeftFieldValue, EqualsOperator, valueLeftFieldValue)
		right := NewComparisionExpression(keyRightFieldValue, EqualsOperator, valueRightFieldValue)

		exp := NewBinaryOperator(left, operator, right)

		err := exp.Validate()

		assert.Nil(t, err)
		assert.Equal(t, BinaryOperatorExpressionType, exp.GetType())
	}
}

func TestShouldFailToValidateBinaryExpressionWhenLeftOperatorIsMissing(t *testing.T) {
	right := NewComparisionExpression(keyRightFieldValue, EqualsOperator, valueRightFieldValue)

	exp := NewBinaryOperator(nil, LogicalAnd, right)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "left")
}

func TestShouldFailToValidateBinaryExpressionWhenLeftOperatorIsNotValid(t *testing.T) {
	left := NewComparisionExpression("", EqualsOperator, valueLeftFieldValue)
	right := NewComparisionExpression(keyRightFieldValue, EqualsOperator, valueRightFieldValue)

	exp := NewBinaryOperator(left, LogicalAnd, right)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "key")
}

func TestShouldFailToValidateBinaryExpressionWhenConjunctionIsNotValid(t *testing.T) {
	left := NewComparisionExpression("leftKey", EqualsOperator, valueLeftFieldValue)
	right := NewComparisionExpression(keyRightFieldValue, EqualsOperator, valueRightFieldValue)

	exp := NewBinaryOperator(left, "FOO", right)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "conjunction")
}

func TestShouldFailToValidateBinaryExpressionWhenRightOperatorIsMissing(t *testing.T) {
	left := NewComparisionExpression(keyLeftFieldValue, EqualsOperator, valueLeftFieldValue)

	exp := NewBinaryOperator(left, LogicalAnd, nil)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "right")
}

func TestShouldFailToValidateBinaryExpressionWhenRightOperatorIsNotValid(t *testing.T) {
	left := NewComparisionExpression(keyLeftFieldValue, EqualsOperator, valueLeftFieldValue)
	right := NewComparisionExpression("", EqualsOperator, valueRightFieldValue)

	exp := NewBinaryOperator(left, LogicalAnd, right)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "key")
}

func TestShouldFailToValidateBinaryExpressionWhenConjunctionIsMissing(t *testing.T) {
	left := NewComparisionExpression(keyLeftFieldValue, EqualsOperator, valueLeftFieldValue)
	right := NewComparisionExpression(keyRightFieldValue, EqualsOperator, valueRightFieldValue)

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
		exp := NewComparisionExpression(keyFieldValue, operator, valueFieldValue)

		err := exp.Validate()

		assert.Nil(t, err)
		assert.Equal(t, LeafExpressionType, exp.GetType())
	}
}

func TestShouldFailToValidateComparisionExpressionWhenKeyIsMissing(t *testing.T) {
	exp := NewComparisionExpression("", EqualsOperator, valueFieldValue)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "key")
}

func TestShouldFailToValidateComparisionExpressionWhenOperatorIsMissing(t *testing.T) {
	exp := NewComparisionExpression(keyFieldValue, "", valueFieldValue)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), operatorFieldName)
}

func TestShouldFailToValidateComparisionExpressionWhenOperatorIsNotValid(t *testing.T) {
	exp := NewComparisionExpression(keyFieldValue, "FOO", valueFieldValue)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), operatorFieldName)
}

func TestShouldFailToValidateComparisionExpressionWhenValueIsMissing(t *testing.T) {
	exp := NewComparisionExpression(keyFieldValue, EqualsOperator, "")

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
		exp := NewUnaryOperationExpression(keyFieldValue, operator)

		err := exp.Validate()

		assert.Nil(t, err)
		assert.Equal(t, LeafExpressionType, exp.GetType())
	}
}

func TestShouldFailToValidateUnaryOperatorExpressionWhenKeyIsMissing(t *testing.T) {
	exp := NewUnaryOperationExpression("", IsEmptyOperator)

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "key")
}

func TestShouldFailToValidateUnaryOperatorExpressionWhenOperatorIsMissing(t *testing.T) {
	exp := NewUnaryOperationExpression(keyFieldValue, "")

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), operatorFieldName)
}

func TestShouldFailToValidateUnaryOperatorExpressionWhenOperatorIsNotValid(t *testing.T) {
	exp := NewUnaryOperationExpression(keyFieldValue, "FOO")

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), operatorFieldName)
}

func TestShouldFailToValidateUnaryOperatorExpressionWhenValueIsSet(t *testing.T) {
	exp := NewComparisionExpression(keyFieldValue, IsEmptyOperator, "")

	err := exp.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "value")
}
