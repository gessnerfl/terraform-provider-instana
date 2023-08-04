package instana_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestRandomID(t *testing.T) {
	id := RandomID()

	require.NotEqual(t, 0, len(id))
}

func TestReadStringArrayParameterFromResourceWhenParameterIsProvided(t *testing.T) {
	ruleIds := []interface{}{"test1", "test2"}
	data := make(map[string]interface{})
	data[AlertingChannelOpsGenieFieldTags] = ruleIds
	resourceData := NewTestHelper[*restapi.AlertingChannel](t).CreateResourceDataForResourceHandle(NewAlertingChannelOpsGenieResourceHandle(), data)
	result := ReadStringArrayParameterFromResource(resourceData, AlertingChannelOpsGenieFieldTags)

	require.NotNil(t, result)
	require.Equal(t, []string{"test1", "test2"}, result)
}

func TestReadStringArrayParameterFromResourceWhenParameterIsMissing(t *testing.T) {
	resourceData := NewTestHelper[*restapi.AlertingChannel](t).CreateEmptyResourceDataForResourceHandle(NewAlertingChannelOpsGenieResourceHandle())
	result := ReadStringArrayParameterFromResource(resourceData, AlertingChannelOpsGenieFieldTags)

	require.Nil(t, result)
}

func TestReadStringSetParameterFromResourceWhenParameterIsProvided(t *testing.T) {
	emails := []interface{}{"test1", "test2"}
	data := make(map[string]interface{})
	data[AlertingChannelEmailFieldEmails] = emails
	resourceData := NewTestHelper[*restapi.AlertingChannel](t).CreateResourceDataForResourceHandle(NewAlertingChannelEmailResourceHandle(), data)
	result := ReadStringSetParameterFromResource(resourceData, AlertingChannelEmailFieldEmails)

	require.NotNil(t, result)
	require.Len(t, result, 2)
	require.Contains(t, result, "test1")
	require.Contains(t, result, "test2")
}

func TestReadStringSetParameterFromResourceWhenParameterIsMissing(t *testing.T) {
	resourceData := NewTestHelper[*restapi.AlertingChannel](t).CreateEmptyResourceDataForResourceHandle(NewAlertingChannelEmailResourceHandle())
	result := ReadStringSetParameterFromResource(resourceData, AlertingChannelEmailFieldEmails)

	require.Nil(t, result)
}

func TestShouldReturnStringRepresentationOfSeverityWarning(t *testing.T) {
	testShouldReturnStringRepresentationOfSeverity(restapi.SeverityWarning, t)
}

func TestShouldReturnStringRepresentationOfSeverityCritical(t *testing.T) {
	testShouldReturnStringRepresentationOfSeverity(restapi.SeverityCritical, t)
}

func testShouldReturnStringRepresentationOfSeverity(severity restapi.Severity, t *testing.T) {
	result, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(severity.GetAPIRepresentation())

	require.Nil(t, err)
	require.Equal(t, severity.GetTerraformRepresentation(), result)
}

func TestShouldFailToConvertStringRepresentationForSeverityWhenIntValueIsNotValid(t *testing.T) {
	result, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(1)

	require.NotNil(t, err)
	require.Equal(t, "INVALID", result)
}

func TestShouldReturnIntRepresentationOfSeverityWarning(t *testing.T) {
	testShouldReturnIntRepresentationOfSeverity(restapi.SeverityWarning, t)
}

func TestShouldReturnIntRepresentationOfSeverityCritical(t *testing.T) {
	testShouldReturnIntRepresentationOfSeverity(restapi.SeverityCritical, t)
}

func testShouldReturnIntRepresentationOfSeverity(severity restapi.Severity, t *testing.T) {
	result, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(severity.GetTerraformRepresentation())

	require.Nil(t, err)
	require.Equal(t, severity.GetAPIRepresentation(), result)
}

func TestShouldFailToConvertIntRepresentationForSeverityWhenStringValueIsNotValid(t *testing.T) {
	result, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation("foo")

	require.NotNil(t, err)
	require.Equal(t, -1, result)
}

func TestShouldSuccessfullyMergeTheTwoProvidedMapsIntoASingleMap(t *testing.T) {
	mapA := map[string]*schema.Schema{"a": {}, "b": {}}
	mapB := map[string]*schema.Schema{"c": {}, "d": {}}
	mapMerged := map[string]*schema.Schema{"a": {}, "b": {}, "c": {}, "d": {}}

	result := MergeSchemaMap(mapA, mapB)

	require.Equal(t, mapMerged, result)
}

func TestShouldReturnIntPointerFromResource(t *testing.T) {
	resourceSchema := make(map[string]*schema.Schema)
	resourceSchema["test"] = &schema.Schema{
		Type:     schema.TypeInt,
		Required: true,
	}

	value := 12
	data := make(map[string]interface{})
	data["test"] = value

	resourceData := schema.TestResourceDataRaw(t, resourceSchema, data)

	require.Equal(t, &value, GetIntPointerFromResourceData(resourceData, "test"))
}

func TestShouldReturnNilWhenIntPointerIsRequestedButNotSetInResource(t *testing.T) {
	resourceSchema := make(map[string]*schema.Schema)
	resourceSchema["test"] = &schema.Schema{
		Type:     schema.TypeInt,
		Required: true,
	}

	data := make(map[string]interface{})

	resourceData := schema.TestResourceDataRaw(t, resourceSchema, data)

	require.Nil(t, GetIntPointerFromResourceData(resourceData, "test"))
}

func TestShouldReturnInt32PointerFromResource(t *testing.T) {
	resourceSchema := make(map[string]*schema.Schema)
	resourceSchema["test"] = &schema.Schema{
		Type:     schema.TypeInt,
		Required: true,
	}

	value := int32(12)
	data := make(map[string]interface{})
	data["test"] = int(value)

	resourceData := schema.TestResourceDataRaw(t, resourceSchema, data)

	require.Equal(t, &value, GetInt32PointerFromResourceData(resourceData, "test"))
}

func TestShouldReturnNilWhenInt32PointerIsRequestedButNotSetInResource(t *testing.T) {
	resourceSchema := make(map[string]*schema.Schema)
	resourceSchema["test"] = &schema.Schema{
		Type:     schema.TypeInt,
		Required: true,
	}

	data := make(map[string]interface{})

	resourceData := schema.TestResourceDataRaw(t, resourceSchema, data)

	require.Nil(t, GetInt32PointerFromResourceData(resourceData, "test"))
}

func TestShouldReturnFloat64PointerFromResource(t *testing.T) {
	resourceSchema := make(map[string]*schema.Schema)
	resourceSchema["test"] = &schema.Schema{
		Type:     schema.TypeFloat,
		Required: true,
	}

	value := float64(12.1)
	data := make(map[string]interface{})
	data["test"] = value

	resourceData := schema.TestResourceDataRaw(t, resourceSchema, data)

	require.Equal(t, &value, GetFloat64PointerFromResourceData(resourceData, "test"))
}

func TestShouldReturnNilWhenFloat64PointerIsRequestedButNotSetInResource(t *testing.T) {
	resourceSchema := make(map[string]*schema.Schema)
	resourceSchema["test"] = &schema.Schema{
		Type:     schema.TypeFloat,
		Required: true,
	}

	data := make(map[string]interface{})

	resourceData := schema.TestResourceDataRaw(t, resourceSchema, data)

	require.Nil(t, GetFloat64PointerFromResourceData(resourceData, "test"))
}

func TestShouldReturnFloat32PointerFromResource(t *testing.T) {
	resourceSchema := make(map[string]*schema.Schema)
	resourceSchema["test"] = &schema.Schema{
		Type:     schema.TypeFloat,
		Required: true,
	}

	value := float32(12.1)
	data := make(map[string]interface{})
	data["test"] = float64(value)

	resourceData := schema.TestResourceDataRaw(t, resourceSchema, data)

	require.Equal(t, &value, GetFloat32PointerFromResourceData(resourceData, "test"))
}

func TestShouldReturnNilWhenFloat32PointerIsRequestedButNotSetInResource(t *testing.T) {
	resourceSchema := make(map[string]*schema.Schema)
	resourceSchema["test"] = &schema.Schema{
		Type:     schema.TypeFloat,
		Required: true,
	}

	data := make(map[string]interface{})

	resourceData := schema.TestResourceDataRaw(t, resourceSchema, data)

	require.Nil(t, GetFloat32PointerFromResourceData(resourceData, "test"))
}

func TestShouldReturnStringPointerFromResource(t *testing.T) {
	resourceSchema := make(map[string]*schema.Schema)
	resourceSchema["test"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}

	value := "test"
	data := make(map[string]interface{})
	data["test"] = value

	resourceData := schema.TestResourceDataRaw(t, resourceSchema, data)

	require.Equal(t, &value, GetStringPointerFromResourceData(resourceData, "test"))
}

func TestShouldReturnNilWhenStringPointerIsRequestedButNotSetInResource(t *testing.T) {
	resourceSchema := make(map[string]*schema.Schema)
	resourceSchema["test"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}

	data := make(map[string]interface{})

	resourceData := schema.TestResourceDataRaw(t, resourceSchema, data)

	require.Nil(t, GetStringPointerFromResourceData(resourceData, "test"))
}

func TestShouldConvertInterfaceSliceToTargetIntSlice(t *testing.T) {
	input := []interface{}{12, 34, 56}
	expectedResult := []int{12, 34, 56}

	result := ConvertInterfaceSlice[int](input)

	require.Equal(t, expectedResult, result)
}
