package instana_test

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/stretchr/testify/assert"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/gessnerfl/terraform-provider-instana/utils"
)

var testSliConfigProviders = map[string]terraform.ResourceProvider{
	"instana": Provider(),
}

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
		threshold = 1
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
		"threshold"			: 1
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

const sliConfigApiPath = restapi.SliConfigResourcePath + "/{id}"
const sliConfigDefinition = "instana_sli_config.example_sli_config"

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
		Providers: testSliConfigProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceDefinitionWithName0,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(sliConfigDefinition, "id"),
					resource.TestCheckResourceAttr(sliConfigDefinition, SliConfigFieldName, "name 0"),
					resource.TestCheckResourceAttr(sliConfigDefinition, SliConfigFieldFullName, "prefix name 0 suffix"),
					resource.TestCheckResourceAttr(sliConfigDefinition, SliConfigFieldInitialEvaluationTimestamp, "0"),
					resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf("%s.0.%s", SliConfigFieldMetricConfiguration, SliConfigFieldMetricName), "metric_name_example"),
					resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf("%s.0.%s", SliConfigFieldMetricConfiguration, SliConfigFieldMetricAggregation), "SUM"),
					resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf("%s.0.%s", SliConfigFieldMetricConfiguration, SliConfigFieldMetricThreshold), "1"),
					resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf("%s.0.%s", SliConfigFieldSliEntity, SliConfigFieldSliType), "application"),
					resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf("%s.0.%s", SliConfigFieldSliEntity, SliConfigFieldApplicationID), "application_id_example"),
					resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf("%s.0.%s", SliConfigFieldSliEntity, SliConfigFieldServiceID), "service_id_example"),
					resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf("%s.0.%s", SliConfigFieldSliEntity, SliConfigFieldEndpointID), "endpoint_id_example"),
					resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf("%s.0.%s", SliConfigFieldSliEntity, SliConfigFieldBoundaryScope), "ALL"),
				),
			},
			{
				Config: resourceDefinitionWithName1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(sliConfigDefinition, "id"),
					resource.TestCheckResourceAttr(sliConfigDefinition, SliConfigFieldName, "name 1"),
					resource.TestCheckResourceAttr(sliConfigDefinition, SliConfigFieldFullName, "prefix name 1 suffix"),
					resource.TestCheckResourceAttr(sliConfigDefinition, SliConfigFieldInitialEvaluationTimestamp, "0"),
					resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf("%s.0.%s", SliConfigFieldMetricConfiguration, SliConfigFieldMetricName), "metric_name_example"),
					resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf("%s.0.%s", SliConfigFieldMetricConfiguration, SliConfigFieldMetricAggregation), "SUM"),
					resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf("%s.0.%s", SliConfigFieldMetricConfiguration, SliConfigFieldMetricThreshold), "1"),
					resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf("%s.0.%s", SliConfigFieldSliEntity, SliConfigFieldSliType), "application"),
					resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf("%s.0.%s", SliConfigFieldSliEntity, SliConfigFieldApplicationID), "application_id_example"),
					resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf("%s.0.%s", SliConfigFieldSliEntity, SliConfigFieldServiceID), "service_id_example"),
					resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf("%s.0.%s", SliConfigFieldSliEntity, SliConfigFieldEndpointID), "endpoint_id_example"),
					resource.TestCheckResourceAttr(sliConfigDefinition, fmt.Sprintf("%s.0.%s", SliConfigFieldSliEntity, SliConfigFieldBoundaryScope), "ALL"),
				),
			},
		},
	})
}

func TestResourceSliConfigDefinition(t *testing.T) {
	resource := NewSliConfigResourceHandle()

	schemaMap := resource.Schema

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
	name := NewSliConfigResourceHandle().ResourceName

	assert.Equal(t, "instana_sli_config", name, "Expected resource name to be instana_sli_config")
}

func TestSliConfigResourceShouldHaveSchemaVersionZero(t *testing.T) {
	assert.Equal(t, 0, NewSliConfigResourceHandle().SchemaVersion)
}

func TestShouldUpdateResourceStateForSliConfigs(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewSliConfigResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	data := restapi.SliConfig{
		ID:   "id",
		Name: "name",
	}

	err := resourceHandle.UpdateState(resourceData, data)

	assert.Nil(t, err)

	assert.Nil(t, err)
	assert.Equal(t, "id", resourceData.Id(), "id should be equal")
	assert.Equal(t, "name", resourceData.Get(SliConfigFieldFullName), "name should be equal to full name")
}

func TestShouldConvertStateOfSliConfigsToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewSliConfigResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId("id")
	resourceData.Set(SliConfigFieldName, "name")
	resourceData.Set(SliConfigFieldFullName, "prefix name suffix")
	resourceData.Set(SliConfigFieldInitialEvaluationTimestamp, 0)

	metricConfigurationStateObject := []map[string]interface{}{
		{
			"metric_name": "test",
			"aggregation": "SUM",
			"threshold":   1.0,
		},
	}
	resourceData.Set(SliConfigFieldMetricConfiguration, metricConfigurationStateObject)

	sliEntityStateObject := []map[string]interface{}{
		{
			"type":           "test_sli_type",
			"application_id": "test_application_id",
			"service_id":     "test_service_id",
			"endpoint_id":    "test_endpoint_id",
			"boundary_scope": "ALL",
		},
	}
	resourceData.Set(SliConfigFieldSliEntity, sliEntityStateObject)

	model, err := resourceHandle.MapStateToDataObject(resourceData, utils.NewResourceNameFormatter("prefix ", " suffix"))

	assert.Nil(t, err)
	assert.IsType(t, restapi.SliConfig{}, model, "Model should be an sli config")
	assert.Equal(t, "id", model.GetID())
	assert.Equal(t, "prefix name suffix", model.(restapi.SliConfig).Name, "name should be equal to full name")
	assert.Equal(t, 0, model.(restapi.SliConfig).InitialEvaluationTimestamp, "initial evaluation timestamp should be 0")
	assert.Equal(t, "test", model.(restapi.SliConfig).MetricConfiguration.Name)
	assert.Equal(t, "SUM", model.(restapi.SliConfig).MetricConfiguration.Aggregation)
	assert.Equal(t, 1.0, model.(restapi.SliConfig).MetricConfiguration.Threshold)
	assert.Equal(t, "test_sli_type", model.(restapi.SliConfig).SliEntity.Type)
	assert.Equal(t, "test_application_id", model.(restapi.SliConfig).SliEntity.ApplicationID)
	assert.Equal(t, "test_service_id", model.(restapi.SliConfig).SliEntity.ServiceID)
	assert.Equal(t, "test_endpoint_id", model.(restapi.SliConfig).SliEntity.EndpointID)
	assert.Equal(t, "ALL", model.(restapi.SliConfig).SliEntity.BoundaryScope)
}

func TestShouldRequireMetricConfigurationThresholdToBeHigherThanZero(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewSliConfigResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId("id")
	resourceData.Set(SliConfigFieldName, "name")
	resourceData.Set(SliConfigFieldFullName, "prefix name suffix")
	resourceData.Set(SliConfigFieldInitialEvaluationTimestamp, 0)

	metricConfigurationStateObject := []map[string]interface{}{
		{
			"metric_name": "test",
			"aggregation": "SUM",
			"threshold":   0.0,
		},
	}
	resourceData.Set(SliConfigFieldMetricConfiguration, metricConfigurationStateObject)

	_, metricThresholdIsOK := resourceData.GetOk("metric_configuration.0.threshold")
	assert.False(t, metricThresholdIsOK)
}
