package restapi_test

import (
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
}

func TestShouldFailToValidateApplicationConfigWhenIDIsMissing(t *testing.T) {
	config :=
		ApplicationConfig{
			Label:              "label",
			MatchSpecification: NewComparisionExpression("key", "EQUALS", "value"),
			Scope:              "scope",
		}

	if err := config.Validate(); err == nil || !strings.Contains(err.Error(), "ID") {
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

	if err := config.Validate(); err == nil || !strings.Contains(err.Error(), "Label") {
		t.Fatal("Expected invalid application config because of missing Label")
	}
}

func TestShouldFailToValidateApplicationConfigWhenMatchSpecificationIsMissing(t *testing.T) {
	config :=
		ApplicationConfig{
			ID:    "id",
			Label: "label",
			Scope: "scope",
		}

	if err := config.Validate(); err == nil || !strings.Contains(err.Error(), "MatchSpecification") {
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

	if err := config.Validate(); err == nil || !strings.Contains(err.Error(), "Key") {
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

	if err := config.Validate(); err == nil || !strings.Contains(err.Error(), "Scope") {
		t.Fatal("Expected invalid application config because of missing Scope")
	}
}

func TestShouldSuccessfullyValidateConsistentBinaryExpression(t *testing.T) {
	left := NewComparisionExpression("keyLeft", "EQUALS", "valueLeft")
	right := NewComparisionExpression("keyRight", "EQUALS", "valueRight")

	exp := NewBinaryOperator(left, "AND", right)

	if err := exp.Validate(); err != nil {
		t.Fatalf("Expected no error but got, %s", err)
	}
}

func TestShouldFailToValidateBinaryExpressionWhenLeftOperatorIsMissing(t *testing.T) {
	right := NewComparisionExpression("keyRight", "EQUALS", "valueRight")

	exp := NewBinaryOperator(nil, "AND", right)

	if err := exp.Validate(); err == nil || !strings.Contains(err.Error(), "Left") {
		t.Fatal("Expected invalid application config because of missing Left operator")
	}
}

func TestShouldFailToValidateBinaryExpressionWhenLeftOperatorIsNotValid(t *testing.T) {
	left := NewComparisionExpression("", "EQUALS", "valueLeft")
	right := NewComparisionExpression("keyRight", "EQUALS", "valueRight")

	exp := NewBinaryOperator(left, "AND", right)

	if err := exp.Validate(); err == nil || !strings.Contains(err.Error(), "Key") {
		t.Fatal("Expected invalid application config because of invalid Left operator")
	}
}

func TestShouldFailToValidateBinaryExpressionWhenRightOperatorIsMissing(t *testing.T) {
	left := NewComparisionExpression("keyLeft", "EQUALS", "valueLeft")

	exp := NewBinaryOperator(left, "AND", nil)

	if err := exp.Validate(); err == nil || !strings.Contains(err.Error(), "Right") {
		t.Fatal("Expected invalid application config because of missing right operator")
	}
}

func TestShouldFailToValidateBinaryExpressionWhenRightOperatorIsNotValid(t *testing.T) {
	left := NewComparisionExpression("keyLeft", "EQUALS", "valueLeft")
	right := NewComparisionExpression("", "EQUALS", "valueRight")

	exp := NewBinaryOperator(left, "AND", right)

	if err := exp.Validate(); err == nil || !strings.Contains(err.Error(), "Key") {
		t.Fatal("Expected invalid application config because of invalid right operator")
	}
}

func TestShouldFailToValidateBinaryExpressionWhenConjunctionIsMissing(t *testing.T) {
	left := NewComparisionExpression("keyLeft", "EQUALS", "valueLeft")
	right := NewComparisionExpression("keyRight", "EQUALS", "valueRight")

	exp := NewBinaryOperator(left, "", right)

	if err := exp.Validate(); err == nil || !strings.Contains(err.Error(), "Conjunction") {
		t.Fatal("Expected invalid application config because of missing conjunction")
	}
}
