package instana_test

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/gessnerfl/terraform-provider-instana/utils"
)

const sliConfigTerraformTemplate = `
provider "instana" {
	api_token = "test-token"
	endpoint = "localhost:{{PORT}}"
	default_name_prefix = "prefix"
	default_name_suffix = "suffix"
}

resource "instana_sli_config" "example_sli_config" {
	name = "name {{ITERATOR}}"
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
	"id"						: "{{id}}",
	"sliName"					: "prefix name suffix",
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
	sliConfigApiPath    = restapi.SliConfigResourcePath + "/{id}"
	sliConfigDefinition = "instana_sli_config.example_sli_config"

	nestedResourceFieldPattern = "%s.0.%s"

	sliConfigID                         = "id"
	sliConfigName                       = "name"
	sliConfigFullName                   = "prefix name suffix"
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
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, sliConfigApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, sliConfigApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, sliConfigApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		json := strings.ReplaceAll(sliConfigServerResponseTemplate, "{{id}}", vars["id"])
		w.Header().Set(contentType, r.Header.Get(contentType))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceDefinitionWithoutName := strings.ReplaceAll(sliConfigTerraformTemplate, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))
	resourceDefinitionWithName0 := strings.ReplaceAll(resourceDefinitionWithoutName, iteratorPlaceholder, "0")
	resourceDefinitionWithName1 := strings.ReplaceAll(resourceDefinitionWithoutName, iteratorPlaceholder, "1")

	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceDefinitionWithName0,
				Check:  resource.ComposeTestCheckFunc(createSliConfigTestCheckFunctions(0)...),
			},
			{
				Config: resourceDefinitionWithName1,
				Check:  resource.ComposeTestCheckFunc(createSliConfigTestCheckFunctions(1)...),
			},
		},
	})
}

func createSliConfigTestCheckFunctions(iteration int) []resource.TestCheckFunc {
	testCheckFunctions := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(sliConfigDefinition, "id"),
		resource.TestCheckResourceAttr(sliConfigDefinition, SliConfigFieldName, fmt.Sprintf("name %d", iteration)),
		resource.TestCheckResourceAttr(sliConfigDefinition, SliConfigFieldFullName, fmt.Sprintf("prefix name %d suffix", iteration)),
		resource.TestCheckResourceAttr(sliConfigDefinition, SliConfigFieldInitialEvaluationTimestamp, "0"),
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(nestedResourceFieldPattern, SliConfigFieldMetricConfiguration, SliConfigFieldMetricName), sliConfigMetricName),
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(nestedResourceFieldPattern, SliConfigFieldMetricConfiguration, SliConfigFieldMetricAggregation), sliConfigMetricAggregation),
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(nestedResourceFieldPattern, SliConfigFieldMetricConfiguration, SliConfigFieldMetricThreshold), "1"),
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(nestedResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldSliType), sliConfigEntityType),
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(nestedResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldApplicationID), sliConfigEntityApplicationID),
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(nestedResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldServiceID), sliConfigEntityServiceID),
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(nestedResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldEndpointID), sliConfigEntityEndpointID),
		resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf(nestedResourceFieldPattern, SliConfigFieldSliEntity, SliConfigFieldBoundaryScope), sliConfigEntityBoundaryScope),
	}
	return testCheckFunctions
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

	assert.Equal(t, "instana_sli_config", name, "Expected resource name to be instana_sli_config")
}

func TestSliConfigResourceShouldHaveSchemaVersionZero(t *testing.T) {
	assert.Equal(t, 0, NewSliConfigResourceHandle().MetaData().SchemaVersion)
}

func TestShouldUpdateResourceStateForSliConfigs(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewSliConfigResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	data := restapi.SliConfig{
		ID:   sliConfigID,
		Name: sliConfigName,
	}

	err := resourceHandle.UpdateState(resourceData, &data)

	assert.Nil(t, err)

	assert.Nil(t, err)
	assert.Equal(t, sliConfigID, resourceData.Id(), "id should be equal")
	assert.Equal(t, sliConfigName, resourceData.Get(SliConfigFieldFullName), "name should be equal to full name")
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

	model, err := resourceHandle.MapStateToDataObject(resourceData, utils.NewResourceNameFormatter("prefix ", " suffix"))

	assert.Nil(t, err)
	assert.IsType(t, &restapi.SliConfig{}, model, "Model should be an sli config")
	assert.Equal(t, sliConfigID, model.GetID())
	assert.Equal(t, sliConfigFullName, model.(*restapi.SliConfig).Name, "name should be equal to full name")
	assert.Equal(t, sliConfigInitialEvaluationTimestamp, model.(*restapi.SliConfig).InitialEvaluationTimestamp, "initial evaluation timestamp should be 0")
	assert.Equal(t, sliConfigMetricName, model.(*restapi.SliConfig).MetricConfiguration.Name)
	assert.Equal(t, sliConfigMetricAggregation, model.(*restapi.SliConfig).MetricConfiguration.Aggregation)
	assert.Equal(t, sliConfigMetricThreshold, model.(*restapi.SliConfig).MetricConfiguration.Threshold)
	assert.Equal(t, sliConfigEntityType, model.(*restapi.SliConfig).SliEntity.Type)
	assert.Equal(t, sliConfigEntityApplicationID, model.(*restapi.SliConfig).SliEntity.ApplicationID)
	assert.Equal(t, sliConfigEntityServiceID, model.(*restapi.SliConfig).SliEntity.ServiceID)
	assert.Equal(t, sliConfigEntityEndpointID, model.(*restapi.SliConfig).SliEntity.EndpointID)
	assert.Equal(t, sliConfigEntityBoundaryScope, model.(*restapi.SliConfig).SliEntity.BoundaryScope)
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
	assert.False(t, metricThresholdIsOK)
}
