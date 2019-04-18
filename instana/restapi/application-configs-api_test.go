package restapi_test

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

func TestShouldSuccussullyValididateConsistentApplicationConfig(t *testing.T) {
	config := ApplicationConfig{
		ID:                 "id",
		Label:              "label",
		MatchSpecification: NewComparisionExpression("key", "EQUALS", "value"),
		Scope:              "scope",
	}

	if err := config.Validate(); err != nil {
		t.Fatalf("Expected no error but got, %s", err)
	}

	if config.GetID() != "id" {
		t.Fatal("Expected id to be valid")
	}
}

func TestShouldFailToValidateApplicationConfigWhenIDIsMissing(t *testing.T) {
	config :=
		ApplicationConfig{
			Label:              "label",
			MatchSpecification: NewComparisionExpression("key", "EQUALS", "value"),
			Scope:              "scope",
		}

	if err := config.Validate(); err == nil || !strings.Contains(err.Error(), "id") {
		t.Fatal("Expected invalid application config because of missing ID")
	}
}

func TestShouldFailToValidateApplicationConfigWhenLabelIsMissing(t *testing.T) {
	config :=
		ApplicationConfig{
			ID:                 "id",
			MatchSpecification: NewComparisionExpression("key", "EQUALS", "value"),
			Scope:              "scope",
		}

	if err := config.Validate(); err == nil || !strings.Contains(err.Error(), "label") {
		t.Fatal("Expected invalid application config because of missing label")
	}
}

func TestShouldFailToValidateApplicationConfigWhenMatchSpecificationIsMissing(t *testing.T) {
	config :=
		ApplicationConfig{
			ID:    "id",
			Label: "label",
			Scope: "scope",
		}

	if err := config.Validate(); err == nil || !strings.Contains(err.Error(), "match specification") {
		t.Fatal("Expected invalid application config because of missing MatchSpecification")
	}
}

func TestShouldFailToValidateApplicationConfigWhenMatchSpecificationIsNotValid(t *testing.T) {
	config := ApplicationConfig{
		ID:                 "id",
		Label:              "label",
		MatchSpecification: NewComparisionExpression("", "EQUALS", "value"),
		Scope:              "scope",
	}

	if err := config.Validate(); err == nil || !strings.Contains(err.Error(), "key") {
		t.Fatal("Expected invalid application config because of invalid match specification")
	}
}

func TestShouldFailToValidateApplicationConfigWhenScopeIsMissing(t *testing.T) {
	config :=
		ApplicationConfig{
			ID:                 "id",
			Label:              "label",
			MatchSpecification: NewComparisionExpression("key", "EQUALS", "value"),
		}

	if err := config.Validate(); err == nil || !strings.Contains(err.Error(), "scope") {
		t.Fatal("Expected invalid application config because of missing Scope")
	}
}

func TestShouldSuccessfullyValidateConsistentBinaryExpression(t *testing.T) {
	for _, operator := range SupportedConjunctions {
		t.Run(fmt.Sprintf("TestShouldSuccessfullyValidateConsistentBinaryExpressionOfType%s", operator), createTestShouldSuccessfullyValidateConsistentBinaryExpression(operator))
	}
}

func createTestShouldSuccessfullyValidateConsistentBinaryExpression(operator string) func(t *testing.T) {
	return func(t *testing.T) {
		left := NewComparisionExpression("keyLeft", "EQUALS", "valueLeft")
		right := NewComparisionExpression("keyRight", "EQUALS", "valueRight")

		exp := NewBinaryOperator(left, operator, right)

		if err := exp.Validate(); err != nil {
			t.Fatalf("Expected no error but got, %s", err)
		}

		if exp.GetType() != BinaryOperatorExpressionType {
			t.Fatal("Expected type to be binary operator")
		}
	}
}

func TestShouldFailToValidateBinaryExpressionWhenLeftOperatorIsMissing(t *testing.T) {
	right := NewComparisionExpression("keyRight", "EQUALS", "valueRight")

	exp := NewBinaryOperator(nil, "AND", right)

	if err := exp.Validate(); err == nil || !strings.Contains(err.Error(), "left") {
		t.Fatal("Expected invalid application config because of missing Left operator")
	}
}

func TestShouldFailToValidateBinaryExpressionWhenLeftOperatorIsNotValid(t *testing.T) {
	left := NewComparisionExpression("", "EQUALS", "valueLeft")
	right := NewComparisionExpression("keyRight", "EQUALS", "valueRight")

	exp := NewBinaryOperator(left, "AND", right)

	if err := exp.Validate(); err == nil || !strings.Contains(err.Error(), "key") {
		t.Fatal("Expected invalid application config because of invalid Left operator")
	}
}

func TestShouldFailToValidateBinaryExpressionWhenConjunctionIsNotValid(t *testing.T) {
	left := NewComparisionExpression("leftKey", "EQUALS", "valueLeft")
	right := NewComparisionExpression("keyRight", "EQUALS", "valueRight")

	exp := NewBinaryOperator(left, "FOO", right)

	if err := exp.Validate(); err == nil || !strings.Contains(err.Error(), "conjunction") {
		t.Fatal("Expected invalid application config because of invalid conjunction")
	}
}

func TestShouldFailToValidateBinaryExpressionWhenRightOperatorIsMissing(t *testing.T) {
	left := NewComparisionExpression("keyLeft", "EQUALS", "valueLeft")

	exp := NewBinaryOperator(left, "AND", nil)

	if err := exp.Validate(); err == nil || !strings.Contains(err.Error(), "right") {
		t.Fatal("Expected invalid application config because of missing right operator")
	}
}

func TestShouldFailToValidateBinaryExpressionWhenRightOperatorIsNotValid(t *testing.T) {
	left := NewComparisionExpression("keyLeft", "EQUALS", "valueLeft")
	right := NewComparisionExpression("", "EQUALS", "valueRight")

	exp := NewBinaryOperator(left, "AND", right)

	if err := exp.Validate(); err == nil || !strings.Contains(err.Error(), "key") {
		t.Fatal("Expected invalid application config because of invalid right operator")
	}
}

func TestShouldFailToValidateBinaryExpressionWhenConjunctionIsMissing(t *testing.T) {
	left := NewComparisionExpression("keyLeft", "EQUALS", "valueLeft")
	right := NewComparisionExpression("keyRight", "EQUALS", "valueRight")

	exp := NewBinaryOperator(left, "", right)

	if err := exp.Validate(); err == nil || !strings.Contains(err.Error(), "conjunction") {
		t.Fatal("Expected invalid application config because of missing conjunction")
	}
}

func TestShouldCreateValidComparisionExpression(t *testing.T) {
	for _, operator := range SupportedComparisionOperators {
		t.Run(fmt.Sprintf("TestShouldCreateValidComparisionExpressionOfType%s", operator), createTestShouldCreateValidComparisionExpression(operator))
	}
}

func createTestShouldCreateValidComparisionExpression(operator string) func(*testing.T) {
	return func(t *testing.T) {
		exp := NewComparisionExpression("key", "EQUALS", "value")

		if err := exp.Validate(); err != nil {
			t.Fatalf("Expected no error but got '%s'", err)
		}

		if exp.GetType() != LeafExpressionType {
			t.Fatal("Expected leaf expression type")
		}
	}
}

func TestShouldFailToValidateComparisionExpressionWhenKeyIsMissing(t *testing.T) {
	exp := NewComparisionExpression("", "EQUALS", "value")

	if err := exp.Validate(); err == nil || !strings.Contains(err.Error(), "key") {
		t.Fatal("Expected invalid tag matcher expression because of missing key")
	}
}

func TestShouldFailToValidateComparisionExpressionWhenOperatorIsMissing(t *testing.T) {
	exp := NewComparisionExpression("key", "", "value")

	if err := exp.Validate(); err == nil || !strings.Contains(err.Error(), "operator") {
		t.Fatal("Expected invalid tag matcher expression because of missing operator")
	}
}

func TestShouldFailToValidateComparisionExpressionWhenOperatorIsNotValid(t *testing.T) {
	exp := NewComparisionExpression("key", "FOO", "value")

	if err := exp.Validate(); err == nil || !strings.Contains(err.Error(), "operator") {
		t.Fatal("Expected invalid tag matcher expression because of missing operator")
	}
}

func TestShouldFailToValidateComparisionExpressionWhenValueIsMissing(t *testing.T) {
	exp := NewComparisionExpression("key", "EQUALS", "")

	if err := exp.Validate(); err == nil || !strings.Contains(err.Error(), "value") {
		t.Fatal("Expected invalid tag matcher expression because of missing value")
	}
}

func TestShouldCreateValidUnaryOperatorExpression(t *testing.T) {
	for _, operator := range SupportedUnaryOperatorExpressionOperators {
		t.Run(fmt.Sprintf("TestShouldCreateValidUnaryOperatorExpressionOfType%s", operator), createTestShouldCreateValidUnaryOperatorExpression(operator))
	}
}

func createTestShouldCreateValidUnaryOperatorExpression(operator string) func(*testing.T) {
	return func(t *testing.T) {
		exp := NewUnaryOperationExpression("key", operator)

		if err := exp.Validate(); err != nil {
			t.Fatalf("Expected no error but got '%s'", err)
		}

		if exp.GetType() != LeafExpressionType {
			t.Fatal("Expected leaf expression type")
		}
	}
}

func TestShouldFailToValidateUnaryOperatorExpressionWhenKeyIsMissing(t *testing.T) {
	exp := NewUnaryOperationExpression("", "IS_EMPTY")

	if err := exp.Validate(); err == nil || !strings.Contains(err.Error(), "key") {
		t.Fatal("Expected invalid tag matcher expression because of missing key")
	}
}

func TestShouldFailToValidateUnaryOperatorExpressionWhenOperatorIsMissing(t *testing.T) {
	exp := NewUnaryOperationExpression("key", "")

	if err := exp.Validate(); err == nil || !strings.Contains(err.Error(), "operator") {
		t.Fatal("Expected invalid tag matcher expression because of missing operator")
	}
}

func TestShouldFailToValidateUnaryOperatorExpressionWhenOperatorIsNotValid(t *testing.T) {
	exp := NewUnaryOperationExpression("key", "FOO")

	if err := exp.Validate(); err == nil || !strings.Contains(err.Error(), "operator") {
		t.Fatal("Expected invalid tag matcher expression because of missing operator")
	}
}

func TestShouldFailToValidateUnaryOperatorExpressionWhenValueIsSet(t *testing.T) {
	exp := NewComparisionExpression("key", "IS_EMPTY", "")

	if err := exp.Validate(); err == nil || !strings.Contains(err.Error(), "value") {
		t.Fatal("Expected invalid tag matcher expression because of missing value")
	}
}
