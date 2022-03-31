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
	"sliName"					: "prefix name %d suffix",
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
	sliConfigFullName                   = resourceFullName
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

func createSliConfigTestCheckFunctions(httpPort int, iteration int) resource.TestStep {
	return resource.TestStep{
		Config: appendProviderConfig(fmt.Sprintf(sliConfigTerraformTemplate, iteration), httpPort),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(sliConfigDefinition, "id"),
			resource.TestCheckResourceAttr(sliConfigDefinition, SliConfigFieldName, formatResourceName(iteration)),
			resource.TestCheckResourceAttr(sliConfigDefinition, SliConfigFieldFullName, formatResourceFullName(iteration)),
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
	resource := NewSliConfigResourceHandle()

	schemaMap := resource.MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(SliConfigFieldFullName)
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

func TestSliConfigResourceShouldHaveSchemaVersionZero(t *testing.T) {
	require.Equal(t, 0, NewSliConfigResourceHandle().MetaData().SchemaVersion)
}

func TestShouldUpdateResourceStateForSliConfigs(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewSliConfigResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	data := restapi.SliConfig{
		ID:   sliConfigID,
		Name: sliConfigFullName,
	}

	err := resourceHandle.UpdateState(resourceData, &data, testHelper.ResourceFormatter())

	require.Nil(t, err)

	require.Nil(t, err)
	require.Equal(t, sliConfigID, resourceData.Id())
	require.Equal(t, sliConfigName, resourceData.Get(SliConfigFieldName))
	require.Equal(t, sliConfigFullName, resourceData.Get(SliConfigFieldFullName))
}

func TestShouldConvertStateOfSliConfigsToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewSliConfigResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId(sliConfigID)
	resourceData.Set(SliConfigFieldName, sliConfigName)
	resourceData.Set(SliConfigFieldFullName, sliConfigFullName)
	resourceData.Set(SliConfigFieldInitialEvaluationTimestamp, 0)

	metricConfigurationStateObject := []map[string]interface{}{
		{
			SliConfigFieldMetricName:        sliConfigMetricName,
			SliConfigFieldMetricAggregation: sliConfigMetricAggregation,
			SliConfigFieldMetricThreshold:   sliConfigMetricThreshold,
		},
	}
	resourceData.Set(SliConfigFieldMetricConfiguration, metricConfigurationStateObject)

	sliEntityStateObject := []map[string]interface{}{
		{
			SliConfigFieldSliType:       sliConfigEntityType,
			SliConfigFieldApplicationID: sliConfigEntityApplicationID,
			SliConfigFieldServiceID:     sliConfigEntityServiceID,
			SliConfigFieldEndpointID:    sliConfigEntityEndpointID,
			SliConfigFieldBoundaryScope: sliConfigEntityBoundaryScope,
		},
	}
	resourceData.Set(SliConfigFieldSliEntity, sliEntityStateObject)

	model, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.IsType(t, &restapi.SliConfig{}, model, "Model should be an sli config")
	require.Equal(t, sliConfigID, model.GetIDForResourcePath())
	require.Equal(t, sliConfigFullName, model.(*restapi.SliConfig).Name, "name should be equal to full name")
	require.Equal(t, sliConfigInitialEvaluationTimestamp, model.(*restapi.SliConfig).InitialEvaluationTimestamp, "initial evaluation timestamp should be 0")
	require.Equal(t, sliConfigMetricName, model.(*restapi.SliConfig).MetricConfiguration.Name)
	require.Equal(t, sliConfigMetricAggregation, model.(*restapi.SliConfig).MetricConfiguration.Aggregation)
	require.Equal(t, sliConfigMetricThreshold, model.(*restapi.SliConfig).MetricConfiguration.Threshold)
	require.Equal(t, sliConfigEntityType, model.(*restapi.SliConfig).SliEntity.Type)
	require.Equal(t, sliConfigEntityApplicationID, model.(*restapi.SliConfig).SliEntity.ApplicationID)
	require.Equal(t, sliConfigEntityServiceID, model.(*restapi.SliConfig).SliEntity.ServiceID)
	require.Equal(t, sliConfigEntityEndpointID, model.(*restapi.SliConfig).SliEntity.EndpointID)
	require.Equal(t, sliConfigEntityBoundaryScope, model.(*restapi.SliConfig).SliEntity.BoundaryScope)
}

func TestShouldRequireMetricConfigurationThresholdToBeHigherThanZero(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewSliConfigResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId(sliConfigID)
	resourceData.Set(SliConfigFieldName, sliConfigName)
	resourceData.Set(SliConfigFieldFullName, sliConfigFullName)
	resourceData.Set(SliConfigFieldInitialEvaluationTimestamp, 0)

	metricConfigurationStateObject := []map[string]interface{}{
		{
			SliConfigFieldMetricName:        sliConfigMetricName,
			SliConfigFieldMetricAggregation: sliConfigMetricAggregation,
			SliConfigFieldMetricThreshold:   0.0,
		},
	}
	resourceData.Set(SliConfigFieldMetricConfiguration, metricConfigurationStateObject)

	_, metricThresholdIsOK := resourceData.GetOk("metric_configuration.0.threshold")
	require.False(t, metricThresholdIsOK)
}
