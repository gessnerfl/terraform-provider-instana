package instana_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

func TestRandomID(t *testing.T) {
	id := RandomID()

	if len(id) == 0 {
		t.Fatal("Expected to get a new id generated")
	}
}

func TestReadStringArrayParameterFromResourceWhenParameterIsProvided(t *testing.T) {
	ruleIds := make([]interface{}, 2)
	ruleIds[0] = "test1"
	ruleIds[1] = "test2"
	data := make(map[string]interface{})
	data[CustomEventSpecificationDownstreamIntegrationIds] = ruleIds
	resourceData := NewTestHelper(t).CreateResourceDataForResourceHandle(NewCustomEventSpecificationWithSystemRuleResourceHandle(), data)
	result := ReadStringArrayParameterFromResource(resourceData, CustomEventSpecificationDownstreamIntegrationIds)

	if result == nil {
		t.Fatal("Expected result to available")
	}

	if len(result) != len(ruleIds) {
		t.Fatal("Expected that result has the same number of elements as ruleIds")
	}

	for i, v := range ruleIds {
		if result[i] != v {
			t.Fatalf("Expected to get rule id %s in result on position %d, value is %s", v, i, result[i])
		}
	}
}

func TestReadStringArrayParameterFromResourceWhenParameterIsMissing(t *testing.T) {
	resourceData := NewTestHelper(t).CreateEmptyResourceDataForResourceHandle(NewCustomEventSpecificationWithSystemRuleResourceHandle())
	result := ReadStringArrayParameterFromResource(resourceData, CustomEventSpecificationDownstreamIntegrationIds)

	if result != nil {
		t.Fatal("Expected result to be nil as no data is provided")
	}
}

func TestShouldReturnStringRepresentationOfSeverityWarning(t *testing.T) {
	testShouldReturnStringRepresentationOfSeverity(restapi.SeverityWarning, t)
}

func TestShouldReturnStringRepresentationOfSeverityCritical(t *testing.T) {
	testShouldReturnStringRepresentationOfSeverity(restapi.SeverityCritical, t)
}

func testShouldReturnStringRepresentationOfSeverity(severity restapi.Severity, t *testing.T) {
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
	testShouldReturnIntRepresentationOfSeverity(restapi.SeverityWarning, t)
}

func TestShouldReturnIntRepresentationOfSeverityCritical(t *testing.T) {
	testShouldReturnIntRepresentationOfSeverity(restapi.SeverityCritical, t)
}

func testShouldReturnIntRepresentationOfSeverity(severity restapi.Severity, t *testing.T) {
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
