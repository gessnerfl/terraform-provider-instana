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

func TestSliConfigTest(t *testing.T) {
	integrationTest := &sliConfigIntegrationTest{}
	unitTest := &sliConfigUnitTest{}
	t.Run("CRUD integration test of with SLI Entity of type application", integrationTest.testCRUDOfSliConfigurationWithSliEntityOfTypeApplication())
	t.Run("should have valid resource schema", unitTest.shouldHaveValidResourceSchema())
	t.Run("should return correct resource name", unitTest.shouldReturnCorrectResourceNameForSliConfigs())
	t.Run("should have schema version one", unitTest.shouldHaveSchemaVersionOne())
	t.Run("should have on schema state upgrader", unitTest.shouldHaveOneStateUpgrader())
	t.Run("should migrate full name to name when executing first state upgrader and full name is available", unitTest.shouldMigrateFullnameToNameWhenExecutingFirstStateUpgraderAndFullnameIsAvailable())
	t.Run("should migrate do nothing when executing first state upgrader and full name is not available", unitTest.shouldDoNothingWhenExecutingFirstStateUpgraderAndFullnameIsNotAvailable())
	t.Run("should update resource state for SLI Config with SLI Entity of type Application", unitTest.shouldUpdateResourceStateForSliConfigWithSliEntityOfTypeApplication())
	t.Run("should convert state of SLI Config with SLI Entity of type Application to Data Model", unitTest.shouldConvertStateOfSliConfigWithEntityOfTypeApplicationToDataModel())
	t.Run("should require metric threshold to be greater than 0", unitTest.shouldRequireMetricConfigurationThresholdToBeGreaterThanZero())
}

type sliConfigIntegrationTest struct{}

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
		application {
			application_id = "application_id_example"
			service_id = "service_id_example"
			endpoint_id = "endpoint_id_example"
			boundary_scope = "ALL"
		}
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

	sliMetricResourceFieldPattern = "%s.0.%s"
	sliEntityResourceFieldPattern = "%s.0.%s.0.%s"

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

func (r *sliConfigIntegrationTest) testCRUDOfSliConfigurationWithSliEntityOfTypeApplication() func(t *testing.T) {
	return func(t *testing.T) {
		httpServer := createMockHttpServerForResource(restapi.SliConfigResourcePath, sliConfigServerResponseTemplate)
		httpServer.Start()
		defer httpServer.Close()

		resource.UnitTest(t, resource.TestCase{
			ProviderFactories: testProviderFactory,
			Steps: []resource.TestStep{
				r.createSliConfigTestCheckFunctions(httpServer.GetPort(), 0),
				testStepImport(sliConfigDefinition),
				r.createSliConfigTestCheckFunctions(httpServer.GetPort(), 1),
				//testStepImport(sliConfigDefinition),
			},
		})
	}
}

func (r *sliConfigIntegrationTest) createSliConfigTestCheckFunctions(httpPort int64, iteration int) resource.TestStep {
	return resource.TestStep{
		Config: appendProviderConfig(fmt.Sprintf(sliConfigTerraformTemplate, iteration), httpPort),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(sliConfigDefinition, "id"),
			resource.TestCheckResourceAttr(sliConfigDefinition, SliConfigFieldName, formatResourceName(iteration)),
			resource.TestCheckResourceAttr(sliConfigDefinition, SliConfigFieldInitialEvaluationTimestamp, "0"),
			resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliMetricResourceFieldPattern, SliConfigFieldMetricConfiguration, SliConfigFieldMetricName), sliConfigMetricName),
			resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliMetricResourceFieldPattern, SliConfigFieldMetricConfiguration, SliConfigFieldMetricAggregation), sliConfigMetricAggregation),
			resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliMetricResourceFieldPattern, SliConfigFieldMetricConfiguration, SliConfigFieldMetricThreshold), "1"),
			resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliEntityResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldSliEntityApplication, SliConfigFieldApplicationID), sliConfigEntityApplicationID),
			resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliEntityResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldSliEntityApplication, SliConfigFieldServiceID), sliConfigEntityServiceID),
			resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliEntityResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldSliEntityApplication, SliConfigFieldEndpointID), sliConfigEntityEndpointID),
			resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(sliEntityResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldSliEntityApplication, SliConfigFieldBoundaryScope), sliConfigEntityBoundaryScope),
		),
	}
}

type sliConfigUnitTest struct{}

func (r *sliConfigUnitTest) shouldHaveValidResourceSchema() func(t *testing.T) {
	return func(t *testing.T) {
		resourceHandle := NewSliConfigResourceHandle()

		schemaMap := resourceHandle.MetaData().Schema

		schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
		schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldName)
		schemaAssert.AssertSchemaIsOptionalAndOfTypeInt(SliConfigFieldInitialEvaluationTimestamp)
		schemaAssert.AssertSchemaIsRequiredAndOfTypeListOfResource(SliConfigFieldSliEntity)

		r.validateMetricsConfig(t, schemaMap)

		r.validateSliEntity(t, schemaMap)
	}
}

func (r *sliConfigUnitTest) validateSliEntity(t *testing.T, schemaMap map[string]*schema.Schema) {
	sliEntitySchemaMap := schemaMap[SliConfigFieldSliEntity].Elem.(*schema.Resource).Schema
	require.Len(t, sliEntitySchemaMap, 4)
	schemaAssert := testutils.NewTerraformSchemaAssert(sliEntitySchemaMap, t)

	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfResource(SliConfigFieldSliEntityApplication)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfResource(SliConfigFieldSliEntityAvailability)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfResource(SliConfigFieldSliEntityWebsiteEventBased)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfResource(SliConfigFieldSliEntityWebsiteTimeBased)

	sliEntityApplicationSchemaMap := sliEntitySchemaMap[SliConfigFieldSliEntityApplication].Elem.(*schema.Resource).Schema
	require.Len(t, sliEntityApplicationSchemaMap, 4)
	schemaAssert = testutils.NewTerraformSchemaAssert(sliEntityApplicationSchemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldBoundaryScope)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldApplicationID)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(SliConfigFieldServiceID)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(SliConfigFieldEndpointID)

	sliEntityAvailabilitySchemaMap := sliEntitySchemaMap[SliConfigFieldSliEntityAvailability].Elem.(*schema.Resource).Schema
	require.Len(t, sliEntityAvailabilitySchemaMap, 6)
	schemaAssert = testutils.NewTerraformSchemaAssert(sliEntityAvailabilitySchemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldApplicationID)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldBoundaryScope)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldBadEventFilterExpression)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldGoodEventFilterExpression)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(SliConfigFieldIncludeInternal, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(SliConfigFieldIncludeSynthetic, false)

	sliEntityWebsiteEventBasedSchemaMap := sliEntitySchemaMap[SliConfigFieldSliEntityWebsiteEventBased].Elem.(*schema.Resource).Schema
	require.Len(t, sliEntityWebsiteEventBasedSchemaMap, 4)
	schemaAssert = testutils.NewTerraformSchemaAssert(sliEntityWebsiteEventBasedSchemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldWebsiteID)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldBeaconType)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldBadEventFilterExpression)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldGoodEventFilterExpression)

	sliEntityWebsiteTimeBasedSchemaMap := sliEntitySchemaMap[SliConfigFieldSliEntityWebsiteTimeBased].Elem.(*schema.Resource).Schema
	require.Len(t, sliEntityWebsiteTimeBasedSchemaMap, 3)
	schemaAssert = testutils.NewTerraformSchemaAssert(sliEntityWebsiteTimeBasedSchemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldWebsiteID)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldBeaconType)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(SliConfigFieldFilterExpression)
}

func (r *sliConfigUnitTest) validateMetricsConfig(t *testing.T, schemaMap map[string]*schema.Schema) {
	metricConfigurationSchemaMap := schemaMap[SliConfigFieldMetricConfiguration].Elem.(*schema.Resource).Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(metricConfigurationSchemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldMetricName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SliConfigFieldMetricAggregation)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeFloat(SliConfigFieldMetricThreshold)
}

func (r *sliConfigUnitTest) shouldReturnCorrectResourceNameForSliConfigs() func(t *testing.T) {
	return func(t *testing.T) {
		name := NewSliConfigResourceHandle().MetaData().ResourceName

		require.Equal(t, "instana_sli_config", name, "Expected resource name to be instana_sli_config")
	}
}

func (r *sliConfigUnitTest) shouldHaveSchemaVersionOne() func(t *testing.T) {
	return func(t *testing.T) {
		require.Equal(t, 1, NewSliConfigResourceHandle().MetaData().SchemaVersion)
	}
}

func (r *sliConfigUnitTest) shouldHaveOneStateUpgrader() func(t *testing.T) {
	return func(t *testing.T) {
		require.Equal(t, 1, len(NewSliConfigResourceHandle().StateUpgraders()))
	}
}

func (r *sliConfigUnitTest) shouldMigrateFullnameToNameWhenExecutingFirstStateUpgraderAndFullnameIsAvailable() func(t *testing.T) {
	return func(t *testing.T) {
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
}

func (r *sliConfigUnitTest) shouldDoNothingWhenExecutingFirstStateUpgraderAndFullnameIsNotAvailable() func(t *testing.T) {
	return func(t *testing.T) {
		input := map[string]interface{}{
			"name": "test",
		}
		result, err := NewSliConfigResourceHandle().StateUpgraders()[0].Upgrade(nil, input, nil)

		require.NoError(t, err)
		require.Equal(t, input, result)
	}
}

func (r *sliConfigUnitTest) shouldUpdateResourceStateForSliConfigWithSliEntityOfTypeApplication() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.SliConfig](t)
		resourceHandle := NewSliConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
		applicationId := "my-application"
		serviceId := "my-service"
		endpointId := "my-endpint"
		boundaryScope := "INBOUND"
		data := restapi.SliConfig{
			ID:   sliConfigID,
			Name: sliConfigName,
			SliEntity: restapi.SliEntity{
				Type:          "application",
				ApplicationID: &applicationId,
				ServiceID:     &serviceId,
				EndpointID:    &endpointId,
				BoundaryScope: &boundaryScope,
			},
		}

		err := resourceHandle.UpdateState(resourceData, &data)

		require.Nil(t, err)

		require.Nil(t, err)
		require.Equal(t, sliConfigID, resourceData.Id())
		require.Equal(t, sliConfigName, resourceData.Get(SliConfigFieldName))
		require.IsType(t, []interface{}{}, resourceData.Get(SliConfigFieldSliEntity))
		sliEntitySlice := resourceData.Get(SliConfigFieldSliEntity).([]interface{})
		require.IsType(t, map[string]interface{}{}, sliEntitySlice[0])
		sliEntityData := sliEntitySlice[0].(map[string]interface{})
		require.IsType(t, []interface{}{}, sliEntityData[SliConfigFieldSliEntityApplication])
		sliEntityApplicationSlice := sliEntityData[SliConfigFieldSliEntityApplication].([]interface{})
		require.IsType(t, map[string]interface{}{}, sliEntityApplicationSlice[0])
		sliEntityApplicationData := sliEntityApplicationSlice[0].(map[string]interface{})
		require.Equal(t, applicationId, sliEntityApplicationData[SliConfigFieldApplicationID])
		require.Equal(t, serviceId, sliEntityApplicationData[SliConfigFieldServiceID])
		require.Equal(t, endpointId, sliEntityApplicationData[SliConfigFieldEndpointID])
		require.Equal(t, boundaryScope, sliEntityApplicationData[SliConfigFieldBoundaryScope])
	}
}

func (r *sliConfigUnitTest) shouldConvertStateOfSliConfigWithEntityOfTypeApplicationToDataModel() func(t *testing.T) {
	return func(t *testing.T) {
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

		sliEntityStateObject := []interface{}{
			map[string]interface{}{
				SliConfigFieldSliEntityApplication: []interface{}{
					map[string]interface{}{
						SliConfigFieldApplicationID: sliConfigEntityApplicationID,
						SliConfigFieldServiceID:     sliConfigEntityServiceID,
						SliConfigFieldEndpointID:    sliConfigEntityEndpointID,
						SliConfigFieldBoundaryScope: sliConfigEntityBoundaryScope,
					},
				},
			},
		}
		setValueOnResourceData(t, resourceData, SliConfigFieldSliEntity, sliEntityStateObject)

		model, err := resourceHandle.MapStateToDataObject(resourceData)

		require.Nil(t, err)
		require.IsType(t, &restapi.SliConfig{}, model, "Model should be an sli config")
		require.Equal(t, sliConfigID, model.GetIDForResourcePath())
		require.Equal(t, sliConfigName, model.Name, "name should be equal to name")
		require.Equal(t, sliConfigInitialEvaluationTimestamp, model.InitialEvaluationTimestamp, "initial evaluation timestamp should be 0")
		require.Equal(t, sliConfigMetricName, model.MetricConfiguration.Name)
		require.Equal(t, sliConfigMetricAggregation, model.MetricConfiguration.Aggregation)
		require.Equal(t, sliConfigMetricThreshold, model.MetricConfiguration.Threshold)
		require.Equal(t, sliConfigEntityType, model.SliEntity.Type)
		require.Equal(t, sliConfigEntityApplicationID, *model.SliEntity.ApplicationID)
		require.Equal(t, sliConfigEntityServiceID, *model.SliEntity.ServiceID)
		require.Equal(t, sliConfigEntityEndpointID, *model.SliEntity.EndpointID)
		require.Equal(t, sliConfigEntityBoundaryScope, *model.SliEntity.BoundaryScope)
	}
}

func (r *sliConfigUnitTest) shouldRequireMetricConfigurationThresholdToBeGreaterThanZero() func(t *testing.T) {
	return func(t *testing.T) {
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
}
