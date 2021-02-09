package instana_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestRandomID(t *testing.T) {
	id := RandomID()

	assert.NotEqual(t, 0, len(id))
}

func TestReadStringArrayParameterFromResourceWhenParameterIsProvided(t *testing.T) {
	ruleIds := []interface{}{"test1", "test2"}
	data := make(map[string]interface{})
	data[AlertingChannelOpsGenieFieldTags] = ruleIds
	resourceData := NewTestHelper(t).CreateResourceDataForResourceHandle(NewAlertingChannelOpsGenieResourceHandle(), data)
	result := ReadStringArrayParameterFromResource(resourceData, AlertingChannelOpsGenieFieldTags)

	assert.NotNil(t, result)
	assert.Equal(t, []string{"test1", "test2"}, result)
}

func TestReadStringArrayParameterFromResourceWhenParameterIsMissing(t *testing.T) {
	resourceData := NewTestHelper(t).CreateEmptyResourceDataForResourceHandle(NewAlertingChannelOpsGenieResourceHandle())
	result := ReadStringArrayParameterFromResource(resourceData, AlertingChannelOpsGenieFieldTags)

	assert.Nil(t, result)
}

func TestReadStringSetParameterFromResourceWhenParameterIsProvided(t *testing.T) {
	emails := []interface{}{"test1", "test2"}
	data := make(map[string]interface{})
	data[AlertingChannelEmailFieldEmails] = emails
	resourceData := NewTestHelper(t).CreateResourceDataForResourceHandle(NewAlertingChannelEmailResourceHandle(), data)
	result := ReadStringSetParameterFromResource(resourceData, AlertingChannelEmailFieldEmails)

	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Contains(t, result, "test1")
	assert.Contains(t, result, "test2")
}

func TestReadStringSetParameterFromResourceWhenParameterIsMissing(t *testing.T) {
	resourceData := NewTestHelper(t).CreateEmptyResourceDataForResourceHandle(NewAlertingChannelEmailResourceHandle())
	result := ReadStringSetParameterFromResource(resourceData, AlertingChannelEmailFieldEmails)

	assert.Nil(t, result)
}

func TestShouldReturnStringRepresentationOfSeverityWarning(t *testing.T) {
	testShouldReturnStringRepresentationOfSeverity(restapi.SeverityWarning, t)
}

func TestShouldReturnStringRepresentationOfSeverityCritical(t *testing.T) {
	testShouldReturnStringRepresentationOfSeverity(restapi.SeverityCritical, t)
}

func testShouldReturnStringRepresentationOfSeverity(severity restapi.Severity, t *testing.T) {
	result, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(severity.GetAPIRepresentation())

	assert.Nil(t, err)
	assert.Equal(t, severity.GetTerraformRepresentation(), result)
}

func TestShouldFailToConvertStringRepresentationForSeverityWhenIntValueIsNotValid(t *testing.T) {
	result, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(1)

	assert.NotNil(t, err)
	assert.Equal(t, "INVALID", result)
}

func TestShouldReturnIntRepresentationOfSeverityWarning(t *testing.T) {
	testShouldReturnIntRepresentationOfSeverity(restapi.SeverityWarning, t)
}

func TestShouldReturnIntRepresentationOfSeverityCritical(t *testing.T) {
	testShouldReturnIntRepresentationOfSeverity(restapi.SeverityCritical, t)
}

func testShouldReturnIntRepresentationOfSeverity(severity restapi.Severity, t *testing.T) {
	result, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(severity.GetTerraformRepresentation())

	assert.Nil(t, err)
	assert.Equal(t, severity.GetAPIRepresentation(), result)
}

func TestShouldFailToConvertIntRepresentationForSeverityWhenStringValueIsNotValid(t *testing.T) {
	result, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation("foo")

	assert.NotNil(t, err)
	assert.Equal(t, -1, result)
}

func TestShoulSuccessfullyMergeTheTwoPrividedMapsIntoASingleMap(t *testing.T) {
	mapA := map[string]*schema.Schema{"a": {}, "b": {}}
	mapB := map[string]*schema.Schema{"c": {}, "d": {}}
	mapMerged := map[string]*schema.Schema{"a": {}, "b": {}, "c": {}, "d": {}}

	result := MergeSchemaMap(mapA, mapB)

	assert.Equal(t, mapMerged, result)
}
