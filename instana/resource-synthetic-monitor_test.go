package instana_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

const syntheticMonitorTerraformTemplate = `
resource "instana_synthetic_monitor" "example" {
	label          = "label"
	active         = true
	locations      = ["location-id"]
	test_frequency = 10
	playback_mode  = "Staggered"
	configuration {
		mark_synthetic_call = true
		retries             = 0
		retry_interval      = 1
		synthetic_type      = "HTTPAction"
		timeout             = "3m"
		url                 = "https://example.com"
		operation           = "GET"
	}
}
`

const syntheticMonitorServerResponseTemplate = `
{
    "id": "%s",
    "label": "label",
    "active": true,
    "locations": ["location-id"],
    "testFrequency": 10,
    "playbackMode": "Staggered",
    "configuration": {
        "syntheticType": "HTTPAction",
        "markSyntheticCall": true,
        "url": "https://example.com",
		"operation": "GET",
		"retryInterval": 1,
		"timeout": "3m"
    }
}
`

const (
	syntheticMonitorDefinition = "instana_synthetic_monitor.example"

	syntheticMonitorID            = "id"
	syntheticMonitorLabel         = "label"
	syntheticMonitorActive        = true
	syntheticMonitorUrl           = "https://example.com"
	syntheticMonitorOperation     = "GET"
	syntheticMonitorSyntheticCall = true
	syntheticMonitorSyntheticType = "HTTPAction"
)

func TestCRUDOfSyntheticMonitorResourceWithMockServer(t *testing.T) {
	id := RandomID()
	resourceRestAPIPath := restapi.SyntheticMonitorResourcePath
	resourceInstanceRestAPIPath := resourceRestAPIPath + "/{id}"

	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPost, resourceRestAPIPath, func(w http.ResponseWriter, r *http.Request) {
		config := &restapi.SyntheticMonitor{}
		err := json.NewDecoder(r.Body).Decode(config)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			r.Write(bytes.NewBufferString("Failed to get request"))
		} else {
			config.ID = id

			w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(config)
			t.Log(config.Configuration.SyntheticType)
		}
	})
	httpServer.AddRoute(http.MethodPut, resourceInstanceRestAPIPath, func(w http.ResponseWriter, r *http.Request) {
		testutils.EchoHandlerFunc(w, r)
	})
	httpServer.AddRoute(http.MethodDelete, resourceInstanceRestAPIPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, resourceInstanceRestAPIPath, func(w http.ResponseWriter, r *http.Request) {
		json := fmt.Sprintf(syntheticMonitorServerResponseTemplate, id)
		w.Header().Set(contentType, r.Header.Get(contentType))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			createSyntheticMonitorTestCheckFunctions(httpServer.GetPort(), 0),
			testStepImport(syntheticMonitorDefinition),
			createSyntheticMonitorTestCheckFunctions(httpServer.GetPort(), 1),
		},
	})
}

func createSyntheticMonitorTestCheckFunctions(httpPort int, iteration int) resource.TestStep {
	nestedConfigPattern := "%s.0.%s"

	return resource.TestStep{
		Config: appendProviderConfig(syntheticMonitorTerraformTemplate, httpPort),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(syntheticMonitorDefinition, "id"),
			resource.TestCheckResourceAttr(syntheticMonitorDefinition, SyntheticMonitorFieldLabel, syntheticMonitorLabel),
			resource.TestCheckResourceAttr(syntheticMonitorDefinition, SyntheticMonitorFieldActive, "true"),
			resource.TestCheckResourceAttr(syntheticMonitorDefinition, fmt.Sprintf(nestedConfigPattern, SyntheticMonitorFieldConfiguration, SyntheticMonitorFieldConfigMarkSyntheticCall), "true"),
			resource.TestCheckResourceAttr(syntheticMonitorDefinition, fmt.Sprintf(nestedConfigPattern, SyntheticMonitorFieldConfiguration, SyntheticMonitorFieldConfigRetries), "0"),
			resource.TestCheckResourceAttr(syntheticMonitorDefinition, fmt.Sprintf(nestedConfigPattern, SyntheticMonitorFieldConfiguration, SyntheticMonitorFieldConfigRetryInterval), "1"),
			resource.TestCheckResourceAttr(syntheticMonitorDefinition, fmt.Sprintf(nestedConfigPattern, SyntheticMonitorFieldConfiguration, SyntheticMonitorFieldConfigSyntheticType), syntheticMonitorSyntheticType),
			resource.TestCheckResourceAttr(syntheticMonitorDefinition, fmt.Sprintf(nestedConfigPattern, SyntheticMonitorFieldConfiguration, SyntheticMonitorFieldConfigTimeout), "3m"),
			resource.TestCheckResourceAttr(syntheticMonitorDefinition, fmt.Sprintf(nestedConfigPattern, SyntheticMonitorFieldConfiguration, SyntheticMonitorFieldConfigOperation), syntheticMonitorOperation),
			resource.TestCheckResourceAttr(syntheticMonitorDefinition, fmt.Sprintf(nestedConfigPattern, SyntheticMonitorFieldConfiguration, SyntheticMonitorFieldConfigUrl), syntheticMonitorUrl),
		),
	}
}

func TestResourceSyntheticMonitorDefinition(t *testing.T) {
	resource := NewSyntheticMonitorResourceHandle()

	schemaMap := resource.MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SyntheticMonitorFieldLabel)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeSetOfStrings(SyntheticMonitorFieldLocations)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(SyntheticMonitorFieldDescription)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(SyntheticMonitorFieldActive, true)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(SyntheticMonitorFieldPlaybackMode)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeInt(SyntheticMonitorFieldTestFrequency)

	syntheticConfigurationSchemaMap := schemaMap[SyntheticMonitorFieldConfiguration].Elem.(*schema.Resource).Schema

	schemaAssert = testutils.NewTerraformSchemaAssert(syntheticConfigurationSchemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SyntheticMonitorFieldConfigSyntheticType)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(SyntheticMonitorFieldConfigMarkSyntheticCall, false)
}

func TestShouldReturnCorrectResourceNameForSyntheticMonitors(t *testing.T) {
	name := NewSyntheticMonitorResourceHandle().MetaData().ResourceName

	require.Equal(t, "instana_synthetic_monitor", name, "Expected resource name to be instana_synthetic_monitor")
}

func TestSyntheticMonitorResourceShouldHaveSchemaVersionZero(t *testing.T) {
	require.Equal(t, 0, NewSyntheticMonitorResourceHandle().MetaData().SchemaVersion)
}

func TestShouldUpdateResourceStateForSyntheticMonitors(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewSyntheticMonitorResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	data := restapi.SyntheticMonitor{
		ID:     syntheticMonitorID,
		Label:  syntheticMonitorLabel,
		Active: syntheticMonitorActive,
	}

	err := resourceHandle.UpdateState(resourceData, &data, testHelper.ResourceFormatter())

	require.Nil(t, err)

	require.Nil(t, err)
	require.Equal(t, syntheticMonitorID, resourceData.Id())
	require.Equal(t, syntheticMonitorActive, resourceData.Get(SyntheticMonitorFieldActive))
}

func TestShouldConvertStateOfSyntheticMonitorsToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewSyntheticMonitorResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId(syntheticMonitorID)
	resourceData.Set(SyntheticMonitorFieldActive, syntheticMonitorActive)
	resourceData.Set(SyntheticMonitorFieldConfigRetries, 2)

	syntheticConfigurationStateObject := []map[string]interface{}{
		{
			SyntheticMonitorFieldConfigMarkSyntheticCall: syntheticMonitorSyntheticCall,
			SyntheticMonitorFieldConfigSyntheticType:     syntheticMonitorSyntheticType,
			SyntheticMonitorFieldConfigRetries:           2,
		},
	}
	resourceData.Set(SyntheticMonitorFieldConfiguration, syntheticConfigurationStateObject)

	model, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.IsType(t, &restapi.SyntheticMonitor{}, model, "Model should be an synthetic monitor")
	require.Equal(t, syntheticMonitorID, model.GetIDForResourcePath())
}
