package instana_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

const sliConfigTerraformTemplate = `
resource "instana_sli_config" "example_sli_config" {
	name = "name %d"
	initial_evaluation_timestamp = 0
	metric_configuration {
		metric_name = "metric_name_example"
		aggregation = "SUM"
		threshold = 1.0
	}
	sli_entity {
		type = "application"
		application_id = "application_id_example"
		service_id = "service_id_example"
		endpoint_id = "endpoint_id_example"
		boundary_scope = "ALL"
	}
}
`

const sliConfigServerResponseTemplate = `
{
	"id"						: "%s",
	"sliName"					: "name %d",
	"initialEvaluationTimestamp": 0,
	"metricConfiguration": {
		"metricName"		: "metric_name_example",
		"metricAggregation"	: "SUM",
		"threshold"			: 1.0
	},
	"sliEntity": {
		"sliType"		: "application",
		"applicationId"	: "application_id_example",
		"serviceId"		: "service_id_example",
		"endpointId"	: "endpoint_id_example",
		"boundaryScope"	: "ALL"
	}
}
`

const (
	sliConfigDefinition = "instana_sli_config.example_sli_config"

	nestedResourceFieldPattern = "%s.0.%s"

	sliConfigID                         = "id"
	sliConfigName                       = resourceName
	sliConfigInitialEvaluationTimestamp = 0
	sliConfigMetricName                 = "metric_name_example"
	sliConfigMetricAggregation          = "SUM"
	sliConfigMetricThreshold            = 1.0
	sliConfigEntityType                 = "application"
	sliConfigEntityApplicationID        = "application_id_example"
	sliConfigEntityServiceID            = "service_id_example"
	sliConfigEntityEndpointID           = "endpoint_id_example"
	sliConfigEntityBoundaryScope        = "ALL"
)

func TestCRUDOfSliConfiguration(t *testing.T) {
	httpServer := createMockHttpServerForResource(restapi.SliConfigResourcePath, sliConfigServerResponseTemplate)
	httpServer.Start()
	defer httpServer.Close()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			createSliConfigTestCheckFunctions(httpServer.GetPort(), 0),
			testStepImport(sliConfigDefinition),
			createSliConfigTestCheckFunctions(httpServer.GetPort(), 1),
			//testStepImport(sliConfigDefinition),
		},
	})
}

func createSliConfigTestCheckFunctions(httpPort int64, iteration int) resource.TestStep {
	return resource.TestStep{
		Config: appendProviderConfig(fmt.Sprintf(sliConfigTerraformTemplate, iteration), httpPort),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(sliConfigDefinition, "id"),
			resource.TestCheckResourceAttr(sliConfigDefinition, SliConfigFieldName, formatResourceName(iteration)),
			resource.TestCheckResourceAttr(sliConfigDefinition, SliConfigFieldInitialEvaluationTimestamp, "0"),
			resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(nestedResourceFieldPattern, SliConfigFieldMetricConfiguration, SliConfigFieldMetricName), sliConfigMetricName),
			resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(nestedResourceFieldPattern, SliConfigFieldMetricConfiguration, SliConfigFieldMetricAggregation), sliConfigMetricAggregation),
			resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(nestedResourceFieldPattern, SliConfigFieldMetricConfiguration, SliConfigFieldMetricThreshold), "1"),
			resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(nestedResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldSliType), sliConfigEntityType),
			resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(nestedResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldApplicationID), sliConfigEntityApplicationID),
			resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(nestedResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldServiceID), sliConfigEntityServiceID),
			resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(nestedResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldEndpointID), sliConfigEntityEndpointID),
			resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(nestedResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldBoundaryScope), sliConfigEntityBoundaryScope),
		),
	}
}

func TestResourceSliConfigDefinition(t *testing.T) {
	resourceHandle := NewSliConfigResourceHandle()

	schemaMap := resourceHandle.MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldName)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeInt(SliConfigFieldInitialEvaluationTimestamp)

	metricConfigurationSchemaMap := schemaMap[SliConfigFieldMetricConfiguration].Elem.(*schema.Resource).Schema

	schemaAssert = testutils.NewTerraformSchemaAssert(metricConfigurationSchemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldMetricName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldMetricAggregation)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeFloat(SliConfigFieldMetricThreshold)

	sliEntitySchemaMap := schemaMap[SliConfigFieldSliEntity].Elem.(*schema.Resource).Schema

	schemaAssert = testutils.NewTerraformSchemaAssert(sliEntitySchemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldSliType)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldBoundaryScope)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(SliConfigFieldApplicationID)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(SliConfigFieldServiceID)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(SliConfigFieldEndpointID)
}

func TestShouldReturnCorrectResourceNameForSliConfigs(t *testing.T) {
	name := NewSliConfigResourceHandle().MetaData().ResourceName

	require.Equal(t, "instana_sli_config", name, "Expected resource name to be instana_sli_config")
}

func TestSliConfigResourceShouldHaveSchemaVersionOne(t *testing.T) {
	require.Equal(t, 1, NewSliConfigResourceHandle().MetaData().SchemaVersion)
}

func TestSliConfigShouldHaveOneStateUpgrader(t *testing.T) {
	require.Equal(t, 1, len(NewSliConfigResourceHandle().StateUpgraders()))
}

func TestSliConfigShouldMigrateFullnameToNameWhenExecutingFirstStateUpgraderAndFullnameIsAvailable(t *testing.T) {
	input := map[string]interface{}{
		"full_name": "test",
	}
	result, err := NewSliConfigResourceHandle().StateUpgraders()[0].Upgrade(nil, input, nil)

	require.NoError(t, err)
	require.Len(t, result, 1)
	require.NotContains(t, result, SliConfigFieldFullName)
	require.Contains(t, result, SliConfigFieldName)
	require.Equal(t, "test", result[SliConfigFieldName])
}

func TestSliConfigShouldDoNothingWhenExecutingFirstStateUpgraderAndFullnameIsNotAvailable(t *testing.T) {
	input := map[string]interface{}{
		"name": "test",
	}
	result, err := NewSliConfigResourceHandle().StateUpgraders()[0].Upgrade(nil, input, nil)

	require.NoError(t, err)
	require.Equal(t, input, result)
}

func TestShouldUpdateResourceStateForSliConfigs(t *testing.T) {
	testHelper := NewTestHelper[*restapi.SliConfig](t)
	resourceHandle := NewSliConfigResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	data := restapi.SliConfig{
		ID:   sliConfigID,
		Name: sliConfigName,
	}

	err := resourceHandle.UpdateState(resourceData, &data, testHelper.ResourceFormatter())

	require.Nil(t, err)

	require.Nil(t, err)
	require.Equal(t, sliConfigID, resourceData.Id())
	require.Equal(t, sliConfigName, resourceData.Get(SliConfigFieldName))
}

func TestShouldConvertStateOfSliConfigsToDataModel(t *testing.T) {
	testHelper := NewTestHelper[*restapi.SliConfig](t)
	resourceHandle := NewSliConfigResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId(sliConfigID)
	setValueOnResourceData(t, resourceData, SliConfigFieldName, sliConfigName)
	setValueOnResourceData(t, resourceData, SliConfigFieldInitialEvaluationTimestamp, 0)

	metricConfigurationStateObject := []map[string]interface{}{
		{
			SliConfigFieldMetricName:        sliConfigMetricName,
			SliConfigFieldMetricAggregation: sliConfigMetricAggregation,
			SliConfigFieldMetricThreshold:   sliConfigMetricThreshold,
		},
	}
	setValueOnResourceData(t, resourceData, SliConfigFieldMetricConfiguration, metricConfigurationStateObject)

	sliEntityStateObject := []map[string]interface{}{
		{
			SliConfigFieldSliType:       sliConfigEntityType,
			SliConfigFieldApplicationID: sliConfigEntityApplicationID,
			SliConfigFieldServiceID:     sliConfigEntityServiceID,
			SliConfigFieldEndpointID:    sliConfigEntityEndpointID,
			SliConfigFieldBoundaryScope: sliConfigEntityBoundaryScope,
		},
	}
	setValueOnResourceData(t, resourceData, SliConfigFieldSliEntity, sliEntityStateObject)

	model, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.IsType(t, &restapi.SliConfig{}, model, "Model should be an sli config")
	require.Equal(t, sliConfigID, model.GetIDForResourcePath())
	require.Equal(t, sliConfigName, model.Name, "name should be equal to name")
	require.Equal(t, sliConfigInitialEvaluationTimestamp, model.InitialEvaluationTimestamp, "initial evaluation timestamp should be 0")
	require.Equal(t, sliConfigMetricName, model.MetricConfiguration.Name)
	require.Equal(t, sliConfigMetricAggregation, model.MetricConfiguration.Aggregation)
	require.Equal(t, sliConfigMetricThreshold, model.MetricConfiguration.Threshold)
	require.Equal(t, sliConfigEntityType, model.SliEntity.Type)
	require.Equal(t, sliConfigEntityApplicationID, model.SliEntity.ApplicationID)
	require.Equal(t, sliConfigEntityServiceID, model.SliEntity.ServiceID)
	require.Equal(t, sliConfigEntityEndpointID, model.SliEntity.EndpointID)
	require.Equal(t, sliConfigEntityBoundaryScope, model.SliEntity.BoundaryScope)
}

func TestShouldRequireMetricConfigurationThresholdToBeHigherThanZero(t *testing.T) {
	testHelper := NewTestHelper[*restapi.SliConfig](t)
	resourceHandle := NewSliConfigResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId(sliConfigID)
	setValueOnResourceData(t, resourceData, SliConfigFieldName, sliConfigName)
	setValueOnResourceData(t, resourceData, SliConfigFieldInitialEvaluationTimestamp, 0)

	metricConfigurationStateObject := []map[string]interface{}{
		{
			SliConfigFieldMetricName:        sliConfigMetricName,
			SliConfigFieldMetricAggregation: sliConfigMetricAggregation,
			SliConfigFieldMetricThreshold:   0.0,
		},
	}
	setValueOnResourceData(t, resourceData, SliConfigFieldMetricConfiguration, metricConfigurationStateObject)

	_, metricThresholdIsOK := resourceData.GetOk("metric_configuration.0.threshold")
	require.False(t, metricThresholdIsOK)
}
