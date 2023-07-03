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

const syntheticTestTerraformTemplate = `
resource "instana_synthetic_test" "example" {
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
	custom_properties = {
		"key1" = "val1"
		"key2" = "val2"
	}
}
`

const syntheticTestServerResponseTemplate = `
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
    },
	"customProperties": {
		"key1": "val1",
		"key2": "val2"
	}
}
`

const (
	syntheticTestDefinition = "instana_synthetic_test.example"

	syntheticTestID            = "id"
	syntheticTestLabel         = "label"
	syntheticTestActive        = true
	syntheticTestUrl           = "https://example.com"
	syntheticTestOperation     = "GET"
	syntheticTestSyntheticCall = true
	syntheticTestSyntheticType = "HTTPAction"
)

func TestCRUDOfSyntheticTestResourceWithMockServer(t *testing.T) {
	id := RandomID()
	resourceRestAPIPath := restapi.SyntheticTestResourcePath
	resourceInstanceRestAPIPath := resourceRestAPIPath + "/{id}"

	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPost, resourceRestAPIPath, func(w http.ResponseWriter, r *http.Request) {
		config := &restapi.SyntheticTest{}
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
		json := fmt.Sprintf(syntheticTestServerResponseTemplate, id)
		w.Header().Set(contentType, r.Header.Get(contentType))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			createSyntheticTestTestCheckFunctions(httpServer.GetPort(), 0),
			testStepImport(syntheticTestDefinition),
			createSyntheticTestTestCheckFunctions(httpServer.GetPort(), 1),
		},
	})
}

func createSyntheticTestTestCheckFunctions(httpPort int, iteration int) resource.TestStep {
	nestedConfigPattern := "%s.0.%s"

	return resource.TestStep{
		Config: appendProviderConfig(syntheticTestTerraformTemplate, httpPort),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(syntheticTestDefinition, "id"),
			resource.TestCheckResourceAttr(syntheticTestDefinition, SyntheticTestFieldLabel, syntheticTestLabel),
			resource.TestCheckResourceAttr(syntheticTestDefinition, SyntheticTestFieldActive, "true"),
			resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(nestedConfigPattern, SyntheticTestFieldConfiguration, SyntheticTestFieldConfigMarkSyntheticCall), "true"),
			resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(nestedConfigPattern, SyntheticTestFieldConfiguration, SyntheticTestFieldConfigRetries), "0"),
			resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(nestedConfigPattern, SyntheticTestFieldConfiguration, SyntheticTestFieldConfigRetryInterval), "1"),
			resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(nestedConfigPattern, SyntheticTestFieldConfiguration, SyntheticTestFieldConfigSyntheticType), syntheticTestSyntheticType),
			resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(nestedConfigPattern, SyntheticTestFieldConfiguration, SyntheticTestFieldConfigTimeout), "3m"),
			resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(nestedConfigPattern, SyntheticTestFieldConfiguration, SyntheticTestFieldConfigOperation), syntheticTestOperation),
			resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(nestedConfigPattern, SyntheticTestFieldConfiguration, SyntheticTestFieldConfigUrl), syntheticTestUrl),
			resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(nestedConfigPattern, SyntheticTestFieldConfiguration, SyntheticTestFieldConfigUrl), syntheticTestUrl),
		),
	}
}

func TestResourceSyntheticTestDefinition(t *testing.T) {
	resource := NewSyntheticTestResourceHandle()

	schemaMap := resource.MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SyntheticTestFieldLabel)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeSetOfStrings(SyntheticTestFieldLocations)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(SyntheticTestFieldDescription)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(SyntheticTestFieldActive, true)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(SyntheticTestFieldPlaybackMode)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeInt(SyntheticTestFieldTestFrequency)

	syntheticConfigurationSchemaMap := schemaMap[SyntheticTestFieldConfiguration].Elem.(*schema.Resource).Schema

	schemaAssert = testutils.NewTerraformSchemaAssert(syntheticConfigurationSchemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SyntheticTestFieldConfigSyntheticType)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(SyntheticTestFieldConfigMarkSyntheticCall, false)
}

func TestShouldReturnCorrectResourceNameForSyntheticTests(t *testing.T) {
	name := NewSyntheticTestResourceHandle().MetaData().ResourceName

	require.Equal(t, "instana_synthetic_test", name, "Expected resource name to be instana_synthetic_test")
}

func TestSyntheticTestResourceShouldHaveSchemaVersionZero(t *testing.T) {
	require.Equal(t, 0, NewSyntheticTestResourceHandle().MetaData().SchemaVersion)
}

func TestShouldUpdateResourceStateForSyntheticTests(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewSyntheticTestResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	data := restapi.SyntheticTest{
		ID:     syntheticTestID,
		Label:  syntheticTestLabel,
		Active: syntheticTestActive,
	}

	err := resourceHandle.UpdateState(resourceData, &data, testHelper.ResourceFormatter())

	require.Nil(t, err)

	require.Nil(t, err)
	require.Equal(t, syntheticTestID, resourceData.Id())
	require.Equal(t, syntheticTestActive, resourceData.Get(SyntheticTestFieldActive))
}

func TestShouldConvertStateOfSyntheticTestsToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewSyntheticTestResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId(syntheticTestID)
	resourceData.Set(SyntheticTestFieldActive, syntheticTestActive)
	resourceData.Set(SyntheticTestFieldConfigRetries, 2)

	syntheticConfigurationStateObject := []map[string]interface{}{
		{
			SyntheticTestFieldConfigMarkSyntheticCall: syntheticTestSyntheticCall,
			SyntheticTestFieldConfigSyntheticType:     syntheticTestSyntheticType,
			SyntheticTestFieldConfigRetries:           2,
		},
	}
	resourceData.Set(SyntheticTestFieldConfiguration, syntheticConfigurationStateObject)

	model, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.IsType(t, &restapi.SyntheticTest{}, model, "Model should be an synthetic test")
	require.Equal(t, syntheticTestID, model.GetIDForResourcePath())
}
