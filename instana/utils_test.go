package instana_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/google/go-cmp/cmp"
)

func TestRandomID(t *testing.T) {
	id := RandomID()

	if len(id) == 0 {
		t.Fatal("Expected to get a new id generated")
	}
}

func TestReadStringArrayParameterFromResourceWhenParameterIsProvided(t *testing.T) {
	ruleIds := []string{"test1", "test2"}
	data := make(map[string]interface{})
	data[RuleBindingFieldRuleIds] = ruleIds
	resourceData := NewTestHelper(t).CreateRuleBindingResourceData(data)
	result := ReadStringArrayParameterFromResource(resourceData, RuleBindingFieldRuleIds)

	if result == nil {
		t.Fatal("Expected result to available")
	}
	if !cmp.Equal(result, ruleIds) {
		t.Fatal("Expected to get rule ids in string array")
	}
}

func TestReadStringArrayParameterFromResourceWhenParameterIsMissing(t *testing.T) {
	resourceData := NewTestHelper(t).CreateEmptyRuleBindingResourceData()
	result := ReadStringArrayParameterFromResource(resourceData, RuleBindingFieldRuleIds)

	if result != nil {
		t.Fatal("Expected result to be nil as no data is provided")
	}
}

func TestShouldReturnStringRepresentationOfSeverityWarning(t *testing.T) {
	testShouldReturnStringRepresentationOfSeverity(SeverityWarning, t)
}

func TestShouldReturnStringRepresentationOfSeverityCritical(t *testing.T) {
	testShouldReturnStringRepresentationOfSeverity(SeverityCritical, t)
}

func testShouldReturnStringRepresentationOfSeverity(severity Severity, t *testing.T) {
	result, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(severity.GetAPIRepresentation())

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}

	if result != severity.GetTerraformRepresentation() {
		t.Fatal("Expected to get proper terraform representation")
	}
}

func TestShouldFailToConvertStringRepresentationForSeverityWhenIntValueIsNotValid(t *testing.T) {
	result, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(1)

	if err == nil {
		t.Fatal("Expected error")
	}

	if result != "INVALID" {
		t.Fatal("Expected to get INVALID string")
	}
}

func TestShouldReturnIntRepresentationOfSeverityWarning(t *testing.T) {
	testShouldReturnIntRepresentationOfSeverity(SeverityWarning, t)
}

func TestShouldReturnIntRepresentationOfSeverityCritical(t *testing.T) {
	testShouldReturnIntRepresentationOfSeverity(SeverityCritical, t)
}

func testShouldReturnIntRepresentationOfSeverity(severity Severity, t *testing.T) {
	result, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(severity.GetTerraformRepresentation())

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}

	if result != severity.GetAPIRepresentation() {
		t.Fatal("Expected to get proper terraform representation")
	}
}

func TestShouldFailToConvertIntRepresentationForSeverityWhenStringValueIsNotValid(t *testing.T) {
	result, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation("foo")

	if err == nil {
		t.Fatal("Expected error")
	}

	if result != -1 {
		t.Fatal("Expected to get -1")
	}
}
