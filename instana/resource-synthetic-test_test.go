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
	label          = "label %d"
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
		body 				= "expected_body"
		validation_string   = "expected_body"
		follow_redirect     = false
		allow_insecure      = true
		expect_status       = 201
		expect_match        = "[a-zA-Z]"
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
    "label": "label %d",
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
		"timeout": "3m",
		"body": "expected_body",
		"validationString": "expected_body",
		"followRedirect": false,
		"allowInsecure": true,
		"expectStatus": 201,
		"expectMatch": "[a-zA-Z]"
    },
	"customProperties": {
		"key1": "val1",
		"key2": "val2"
	}
}
`

const syntheticTestWithApplicationIdTerraformTemplate = `
resource "instana_synthetic_test" "example" {
	label          = "label %d"
	active         = true
	application_id = "application-id"
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
		body 				= "expected_body"
		validation_string   = "expected_body"
		follow_redirect     = false
		allow_insecure      = true
		expect_status       = 201
		expect_match        = "[a-zA-Z]"
	}
	custom_properties = {
		"key1" = "val1"
		"key2" = "val2"
	}
}
`

const syntheticTestWithApplicationIdServerResponseTemplate = `
{
    "id": "%s",
    "label": "label %d",
    "active": true,
    "applicationId": "application-id",
    "locations": ["location-id"],
    "testFrequency": 10,
    "playbackMode": "Staggered",
    "configuration": {
        "syntheticType": "HTTPAction",
        "markSyntheticCall": true,
        "url": "https://example.com",
		"operation": "GET",
		"retryInterval": 1,
		"timeout": "3m",
		"body": "expected_body",
		"validationString": "expected_body",
		"followRedirect": false,
		"allowInsecure": true,
		"expectStatus": 201,
		"expectMatch": "[a-zA-Z]"
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

func TestCRUDOfSyntheticTestResourceUsingMockServer(t *testing.T) {
	id := RandomID()
	httpServer := createTestServerForSyntheticTestResource(id, syntheticTestServerResponseTemplate)
	defer httpServer.Close()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			createSyntheticTestTestCheckFunctions(syntheticTestTerraformTemplate, httpServer.GetPort(), 0),
			testStepImport(syntheticTestDefinition),
			createSyntheticTestTestCheckFunctions(syntheticTestTerraformTemplate, httpServer.GetPort(), 1),
			testStepImport(syntheticTestDefinition),
		},
	})
}

func TestCRUDOfSyntheticTestResourceWithApplicationIdUsingMockServer(t *testing.T) {
	id := RandomID()
	httpServer := createTestServerForSyntheticTestResource(id, syntheticTestWithApplicationIdServerResponseTemplate)
	defer httpServer.Close()

	additionalCheck := resource.TestCheckResourceAttr(syntheticTestDefinition, SyntheticTestFieldApplicationID, "application-id")
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			createSyntheticTestTestCheckFunctions(syntheticTestWithApplicationIdTerraformTemplate, httpServer.GetPort(), 0, additionalCheck),
			testStepImport(syntheticTestDefinition),
			createSyntheticTestTestCheckFunctions(syntheticTestWithApplicationIdTerraformTemplate, httpServer.GetPort(), 1, additionalCheck),
			testStepImport(syntheticTestDefinition),
		},
	})
}

func createSyntheticTestTestCheckFunctions(template string, httpPort int64, iteration int, additionalChecks ...resource.TestCheckFunc) resource.TestStep {
	nestedConfigPattern := "%s.0.%s"

	defaultChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(syntheticTestDefinition, "id"),
		resource.TestCheckResourceAttr(syntheticTestDefinition, SyntheticTestFieldLabel, fmt.Sprintf("%s %d", syntheticTestLabel, iteration)),
		resource.TestCheckResourceAttr(syntheticTestDefinition, SyntheticTestFieldActive, "true"),
		resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(nestedConfigPattern, SyntheticTestFieldConfiguration, SyntheticTestFieldConfigMarkSyntheticCall), "true"),
		resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(nestedConfigPattern, SyntheticTestFieldConfiguration, SyntheticTestFieldConfigRetries), "0"),
		resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(nestedConfigPattern, SyntheticTestFieldConfiguration, SyntheticTestFieldConfigRetryInterval), "1"),
		resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(nestedConfigPattern, SyntheticTestFieldConfiguration, SyntheticTestFieldConfigSyntheticType), syntheticTestSyntheticType),
		resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(nestedConfigPattern, SyntheticTestFieldConfiguration, SyntheticTestFieldConfigTimeout), "3m"),
		resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(nestedConfigPattern, SyntheticTestFieldConfiguration, SyntheticTestFieldConfigOperation), syntheticTestOperation),
		resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(nestedConfigPattern, SyntheticTestFieldConfiguration, SyntheticTestFieldConfigUrl), syntheticTestUrl),
	}
	checks := append(defaultChecks, additionalChecks...)

	return resource.TestStep{
		Config: appendProviderConfig(fmt.Sprintf(template, iteration), httpPort),
		Check:  resource.ComposeTestCheckFunc(checks...),
	}
}

func createTestServerForSyntheticTestResource(id string, responseTemplate string) testutils.TestHTTPServer {
	resourceRestAPIPath := restapi.SyntheticTestResourcePath
	resourceInstanceRestAPIPath := resourceRestAPIPath + "/{id}"

	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPost, resourceRestAPIPath, func(w http.ResponseWriter, r *http.Request) {
		config := &restapi.SyntheticTest{}
		err := json.NewDecoder(r.Body).Decode(config)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = r.Write(bytes.NewBufferString("Failed to get request"))
			if err != nil {
				fmt.Printf("failed to write response; %s\n", err)
			}
		} else {
			config.ID = id

			w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(config)
			if err != nil {
				fmt.Printf("failed to encode json; %s\n", err)
			}
		}
	})
	httpServer.AddRoute(http.MethodPut, resourceInstanceRestAPIPath, func(w http.ResponseWriter, r *http.Request) {
		testutils.EchoHandlerFunc(w, r)
	})
	httpServer.AddRoute(http.MethodDelete, resourceInstanceRestAPIPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, resourceInstanceRestAPIPath, func(w http.ResponseWriter, r *http.Request) {
		puts := httpServer.GetCallCount(http.MethodPut, resourceRestAPIPath+"/"+id)
		jsonData := fmt.Sprintf(responseTemplate, id, puts)
		w.Header().Set(contentType, r.Header.Get(contentType))
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(jsonData))
		if err != nil {
			fmt.Printf("failed to write response; %s\n", err)
		}
	})
	httpServer.Start()
	return httpServer
}

func TestResourceSyntheticTestDefinition(t *testing.T) {
	resourceHandle := NewSyntheticTestResourceHandle()

	schemaMap := resourceHandle.MetaData().Schema

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
	testHelper := NewTestHelper[*restapi.SyntheticTest](t)
	resourceHandle := NewSyntheticTestResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	data := restapi.SyntheticTest{
		ID:     syntheticTestID,
		Label:  syntheticTestLabel,
		Active: syntheticTestActive,
	}

	err := resourceHandle.UpdateState(resourceData, &data)

	require.Nil(t, err)

	require.Nil(t, err)
	require.Equal(t, syntheticTestID, resourceData.Id())
	require.Equal(t, syntheticTestActive, resourceData.Get(SyntheticTestFieldActive))
}

func TestShouldConvertStateOfSyntheticTestsToDataModel(t *testing.T) {
	testHelper := NewTestHelper[*restapi.SyntheticTest](t)
	resourceHandle := NewSyntheticTestResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId(syntheticTestID)
	setValueOnResourceData(t, resourceData, SyntheticTestFieldActive, syntheticTestActive)

	syntheticConfigurationStateObject := []map[string]interface{}{
		{
			SyntheticTestFieldConfigMarkSyntheticCall: syntheticTestSyntheticCall,
			SyntheticTestFieldConfigSyntheticType:     syntheticTestSyntheticType,
			SyntheticTestFieldConfigRetries:           2,
		},
	}
	setValueOnResourceData(t, resourceData, SyntheticTestFieldConfiguration, syntheticConfigurationStateObject)

	model, err := resourceHandle.MapStateToDataObject(resourceData)

	require.Nil(t, err)
	require.IsType(t, &restapi.SyntheticTest{}, model, "Model should be an synthetic test")
	require.Equal(t, syntheticTestID, model.GetIDForResourcePath())
}
