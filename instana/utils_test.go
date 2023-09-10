package instana_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestUtils(t *testing.T) {
	unitTest := &utilsUnitTest{}
	t.Run("should generate random id", unitTest.shouldGenerateRandomID)
	t.Run("should read string array parameter from map when parameter is provided", unitTest.shouldReadStringArrayParameterFromMapWhenParameterIsProvided)
	t.Run("should read string array parameter from map when parameter is missing", unitTest.shouldReadStringArrayParameterFromMapWhenParameterIsMissing)
	t.Run("should read string set parameter from resource when parameter is provided", unitTest.shouldReadStringSetParameterFromResourceWhenParameterIsProvided)
	t.Run("should read string set parameter from resource when parameter is missing", unitTest.shouldReadStringSetParameterFromResourceWhenParameterIsMissing)
	t.Run("should read string set parameter from map when parameter is provided", unitTest.shouldReadStringSetParameterFromMapWhenParameterIsProvided)
	t.Run("should read string set parameter from map when parameter is missing", unitTest.shouldReadStringSetParameterFromMapWhenParameterIsMissing)
	t.Run("should return string representation of severity warning", unitTest.shouldShouldReturnStringRepresentationOfSeverityWarning)
	t.Run("should return string representation of severity critical", unitTest.shouldShouldReturnStringRepresentationOfSeverityCritical)
	t.Run("should fail to return string representation of severity when int value is not valid", unitTest.shouldShouldFailToConvertStringRepresentationForSeverityWhenIntValueIsNotValid)
	t.Run("should return int representation of severity warning", unitTest.shouldShouldReturnIntRepresentationOfSeverityWarning)
	t.Run("should return int representation of severity critical", unitTest.shouldShouldReturnIntRepresentationOfSeverityCritical)
	t.Run("should fail to return int representation of severity when string value is not valid", unitTest.shouldShouldFailToConvertIntRepresentationForSeverityWhenStringValueIsNotValid)
	t.Run("should successfully merge two provided maps into a single map", unitTest.shouldShouldSuccessfullyMergeTheTwoProvidedMapsIntoASingleMap)
	t.Run("should return int pointer from resource", unitTest.shouldShouldReturnIntPointerFromResource)
	t.Run("should return nil when int pointer is requested but not set in resource", unitTest.shouldShouldReturnNilWhenIntPointerIsRequestedButNotSetInResource)
	t.Run("should return int32 pointer from resource", unitTest.shouldShouldReturnInt32PointerFromResource)
	t.Run("should return nil when int32 pointer is requested but not set in resource", unitTest.shouldShouldReturnNilWhenInt32PointerIsRequestedButNotSetInResource)
	t.Run("should return float64 pointer from resource", unitTest.shouldShouldReturnFloat64PointerFromResource)
	t.Run("should return nil when float64 pointer is requested but not set in resource", unitTest.shouldShouldReturnNilWhenFloat64PointerIsRequestedButNotSetInResource)
	t.Run("should return float32 pointer from resource", unitTest.shouldShouldReturnFloat32PointerFromResource)
	t.Run("should return nil when float32 pointer is requested but not set in resource", unitTest.shouldShouldReturnNilWhenFloat32PointerIsRequestedButNotSetInResource)
	t.Run("should return string pointer from resource", unitTest.shouldShouldReturnStringPointerFromResource)
	t.Run("should return nil when string pointer is requested but not set in resource", unitTest.shouldShouldReturnNilWhenStringPointerIsRequestedButNotSetInResource)
	t.Run("should return pointer value from map", unitTest.shouldShouldReturnPointerValueFromMap)
	t.Run("should return nil when pointer is requested but value is not set in map", unitTest.shouldShouldReturnNilWhenPointerIsRequestedButNotSetInMap)
	t.Run("should convert interface slice to target int slice", unitTest.shouldShouldConvertInterfaceSliceToTargetIntSlice)
	t.Run("should get pointer from map", unitTest.shouldTestGetPointerFromMap)
}

type utilsUnitTest struct{}

func (r *utilsUnitTest) shouldGenerateRandomID(t *testing.T) {
	id := RandomID()

	require.NotEqual(t, 0, len(id))
}

func (r *utilsUnitTest) shouldReadStringArrayParameterFromMapWhenParameterIsProvided(t *testing.T) {
	emails := []interface{}{"test1", "test2"}
	result := ReadArrayParameterFromMap[string](map[string]interface{}{AlertingChannelEmailFieldEmails: emails}, AlertingChannelEmailFieldEmails)

	require.NotNil(t, result)
	require.Len(t, result, 2)
	require.Contains(t, result, "test1")
	require.Contains(t, result, "test2")
}

func (r *utilsUnitTest) shouldReadStringArrayParameterFromMapWhenParameterIsMissing(t *testing.T) {
	result := ReadArrayParameterFromMap[string](map[string]interface{}{"foo": "bar"}, AlertingChannelOpsGenieFieldTags)

	require.Nil(t, result)
}

func (r *utilsUnitTest) shouldReadStringSetParameterFromResourceWhenParameterIsProvided(t *testing.T) {
	integrationIds := []interface{}{"test1", "test2"}
	data := make(map[string]interface{})
	data[AlertingConfigFieldIntegrationIds] = integrationIds
	resourceData := NewTestHelper[*restapi.AlertingConfiguration](t).CreateResourceDataForResourceHandle(NewAlertingConfigResourceHandle(), data)
	result := ReadStringSetParameterFromResource(resourceData, AlertingConfigFieldIntegrationIds)

	require.NotNil(t, result)
	require.Len(t, result, 2)
	require.Contains(t, result, "test1")
	require.Contains(t, result, "test2")
}

func (r *utilsUnitTest) shouldReadStringSetParameterFromResourceWhenParameterIsMissing(t *testing.T) {
	resourceData := NewTestHelper[*restapi.AlertingConfiguration](t).CreateEmptyResourceDataForResourceHandle(NewAlertingConfigResourceHandle())
	result := ReadStringSetParameterFromResource(resourceData, AlertingConfigFieldIntegrationIds)

	require.Nil(t, result)
}

func (r *utilsUnitTest) shouldReadStringSetParameterFromMapWhenParameterIsProvided(t *testing.T) {
	emails := []interface{}{"test1", "test2"}
	result := ReadSetParameterFromMap[string](map[string]interface{}{AlertingChannelEmailFieldEmails: schema.NewSet(schema.HashString, emails)}, AlertingChannelEmailFieldEmails)

	require.NotNil(t, result)
	require.Len(t, result, 2)
	require.Contains(t, result, "test1")
	require.Contains(t, result, "test2")
}

func (r *utilsUnitTest) shouldReadStringSetParameterFromMapWhenParameterIsMissing(t *testing.T) {
	result := ReadSetParameterFromMap[string](map[string]interface{}{"foo": "bar"}, AlertingChannelEmailFieldEmails)

	require.Nil(t, result)
}

func (r *utilsUnitTest) shouldShouldReturnStringRepresentationOfSeverityWarning(t *testing.T) {
	r.testShouldReturnStringRepresentationOfSeverity(restapi.SeverityWarning, t)
}

func (r *utilsUnitTest) shouldShouldReturnStringRepresentationOfSeverityCritical(t *testing.T) {
	r.testShouldReturnStringRepresentationOfSeverity(restapi.SeverityCritical, t)
}

func (r *utilsUnitTest) testShouldReturnStringRepresentationOfSeverity(severity restapi.Severity, t *testing.T) {
	result, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(severity.GetAPIRepresentation())

	require.Nil(t, err)
	require.Equal(t, severity.GetTerraformRepresentation(), result)
}

func (r *utilsUnitTest) shouldShouldFailToConvertStringRepresentationForSeverityWhenIntValueIsNotValid(t *testing.T) {
	result, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(1)

	require.NotNil(t, err)
	require.Equal(t, "INVALID", result)
}

func (r *utilsUnitTest) shouldShouldReturnIntRepresentationOfSeverityWarning(t *testing.T) {
	r.testShouldReturnIntRepresentationOfSeverity(restapi.SeverityWarning, t)
}

func (r *utilsUnitTest) shouldShouldReturnIntRepresentationOfSeverityCritical(t *testing.T) {
	r.testShouldReturnIntRepresentationOfSeverity(restapi.SeverityCritical, t)
}

func (r *utilsUnitTest) testShouldReturnIntRepresentationOfSeverity(severity restapi.Severity, t *testing.T) {
	result, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(severity.GetTerraformRepresentation())

	require.Nil(t, err)
	require.Equal(t, severity.GetAPIRepresentation(), result)
}

func (r *utilsUnitTest) shouldShouldFailToConvertIntRepresentationForSeverityWhenStringValueIsNotValid(t *testing.T) {
	result, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation("foo")

	require.NotNil(t, err)
	require.Equal(t, -1, result)
}

func (r *utilsUnitTest) shouldShouldSuccessfullyMergeTheTwoProvidedMapsIntoASingleMap(t *testing.T) {
	mapA := map[string]*schema.Schema{"a": {}, "b": {}}
	mapB := map[string]*schema.Schema{"c": {}, "d": {}}
	mapMerged := map[string]*schema.Schema{"a": {}, "b": {}, "c": {}, "d": {}}

	result := MergeSchemaMap(mapA, mapB)

	require.Equal(t, mapMerged, result)
}

func (r *utilsUnitTest) shouldShouldReturnIntPointerFromResource(t *testing.T) {
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

func (r *utilsUnitTest) shouldShouldReturnNilWhenIntPointerIsRequestedButNotSetInResource(t *testing.T) {
	resourceSchema := make(map[string]*schema.Schema)
	resourceSchema["test"] = &schema.Schema{
		Type:     schema.TypeInt,
		Required: true,
	}

	data := make(map[string]interface{})

	resourceData := schema.TestResourceDataRaw(t, resourceSchema, data)

	require.Nil(t, GetIntPointerFromResourceData(resourceData, "test"))
}

func (r *utilsUnitTest) shouldShouldReturnInt32PointerFromResource(t *testing.T) {
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

func (r *utilsUnitTest) shouldShouldReturnNilWhenInt32PointerIsRequestedButNotSetInResource(t *testing.T) {
	resourceSchema := make(map[string]*schema.Schema)
	resourceSchema["test"] = &schema.Schema{
		Type:     schema.TypeInt,
		Required: true,
	}

	data := make(map[string]interface{})

	resourceData := schema.TestResourceDataRaw(t, resourceSchema, data)

	require.Nil(t, GetInt32PointerFromResourceData(resourceData, "test"))
}

func (r *utilsUnitTest) shouldShouldReturnFloat64PointerFromResource(t *testing.T) {
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

func (r *utilsUnitTest) shouldShouldReturnNilWhenFloat64PointerIsRequestedButNotSetInResource(t *testing.T) {
	resourceSchema := make(map[string]*schema.Schema)
	resourceSchema["test"] = &schema.Schema{
		Type:     schema.TypeFloat,
		Required: true,
	}

	data := make(map[string]interface{})

	resourceData := schema.TestResourceDataRaw(t, resourceSchema, data)

	require.Nil(t, GetFloat64PointerFromResourceData(resourceData, "test"))
}

func (r *utilsUnitTest) shouldShouldReturnFloat32PointerFromResource(t *testing.T) {
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

func (r *utilsUnitTest) shouldShouldReturnNilWhenFloat32PointerIsRequestedButNotSetInResource(t *testing.T) {
	resourceSchema := make(map[string]*schema.Schema)
	resourceSchema["test"] = &schema.Schema{
		Type:     schema.TypeFloat,
		Required: true,
	}

	data := make(map[string]interface{})

	resourceData := schema.TestResourceDataRaw(t, resourceSchema, data)

	require.Nil(t, GetFloat32PointerFromResourceData(resourceData, "test"))
}

func (r *utilsUnitTest) shouldShouldReturnStringPointerFromResource(t *testing.T) {
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

func (r *utilsUnitTest) shouldShouldReturnNilWhenStringPointerIsRequestedButNotSetInResource(t *testing.T) {
	resourceSchema := make(map[string]*schema.Schema)
	resourceSchema["test"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}

	data := make(map[string]interface{})

	resourceData := schema.TestResourceDataRaw(t, resourceSchema, data)

	require.Nil(t, GetStringPointerFromResourceData(resourceData, "test"))
}

func (r *utilsUnitTest) shouldShouldReturnPointerValueFromMap(t *testing.T) {
	stringValue := "myString"
	intValue := 1234
	data := map[string]interface{}{
		"stringValue": stringValue,
		"intValue":    intValue,
	}

	require.Equal(t, &stringValue, GetPointerFromMap[string](data, "stringValue"))
	require.Equal(t, &intValue, GetPointerFromMap[int](data, "intValue"))
}

func (r *utilsUnitTest) shouldShouldReturnNilWhenPointerIsRequestedButNotSetInMap(t *testing.T) {
	stringValue := "myString"
	intValue := 1234
	data := map[string]interface{}{
		"stringValue": stringValue,
		"intValue":    intValue,
	}

	require.Nil(t, GetPointerFromMap[string](data, "otherStringValue"))
	require.Nil(t, GetPointerFromMap[int](data, "otherIntValue"))
}

func (r *utilsUnitTest) shouldShouldConvertInterfaceSliceToTargetIntSlice(t *testing.T) {
	input := []interface{}{12, 34, 56}
	expectedResult := []int{12, 34, 56}

	result := ConvertInterfaceSlice[int](input)

	require.Equal(t, expectedResult, result)
}

func (r *utilsUnitTest) shouldTestGetPointerFromMap(t *testing.T) {
	t.Run("should return string pointer when key is in map and string value is not null", createShouldReturnPointerValueWhenKeyIsAvailableInMapAndValueIsNotNul("test"))
	t.Run("should return int pointer when key is in map and int value is not null", createShouldReturnPointerValueWhenKeyIsAvailableInMapAndValueIsNotNul(1234))
	t.Run("should return bool pointer when key is in map and bool value is not null", createShouldReturnPointerValueWhenKeyIsAvailableInMapAndValueIsNotNul(true))
	t.Run("should return nil when key is not available in map", createShouldReturnNilValueWhenKeyIsNotAvailableInMap())
	t.Run("should return nil when key is available in map and value is nil", createShouldReturnNilValueWhenKeyIsAvailableInMapAndValueIsNil())
}

func createShouldReturnPointerValueWhenKeyIsAvailableInMapAndValueIsNotNul[T any](expectedValue T) func(t *testing.T) {
	return func(t *testing.T) {
		key := "key"
		data := map[string]interface{}{
			key: expectedValue,
		}

		result := GetPointerFromMap[T](data, key)

		require.NotNil(t, result)
		require.Equal(t, expectedValue, *result)
	}
}

func createShouldReturnNilValueWhenKeyIsNotAvailableInMap() func(t *testing.T) {
	return func(t *testing.T) {
		key := "key"
		data := map[string]interface{}{
			"other": "value",
		}

		result := GetPointerFromMap[string](data, key)

		require.Nil(t, result)
	}

}

func createShouldReturnNilValueWhenKeyIsAvailableInMapAndValueIsNil() func(t *testing.T) {
	return func(t *testing.T) {
		key := "key"
		data := map[string]interface{}{
			key: nil,
		}

		result := GetPointerFromMap[string](data, key)

		require.Nil(t, result)
	}

}
