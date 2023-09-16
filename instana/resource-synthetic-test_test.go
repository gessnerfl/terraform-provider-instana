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

func TestSyntheticTestResource(t *testing.T) {
	ut := &syntheticTestUnitTest{}
	t.Run("CRUD integration test with HTTP Action configuration without application id", syntheticTestHttpActionWithoutApplicationIdIntegrationTest().testCrud)
	t.Run("CRUD integration test with HTTP Action configuration with application id", syntheticTestHttpActionWithApplicationIdIntegrationTest().testCrud)
	t.Run("CRUD integration test with HTTP Script", syntheticTestHttpScriptIntegrationTest().testCrud)
	t.Run("should have valid schema", ut.resourceDefinitionShouldBeValid)
	t.Run("should return correct resource name", ut.shouldReturnCorrectResourceName)
	t.Run("should have schema version zero", ut.shouldHaveSchemaVersionZero)
	t.Run("should have schema no state upgrader", ut.shouldHaveNoStateUpgrader)
	t.Run("should update resource state for http script config", ut.shouldUpdateResourceStateForHttpScript)
	t.Run("should update resource state for http action config", ut.shouldUpdateResourceStateForHttpAction)
	t.Run("should return error when trying to update state and config type is not supported", ut.shouldReturnErrorWhenTryingToUpdateStateAndConfigTypeIsNotSupported)
	t.Run("should map state to data model with http action config", ut.shouldMapStateToDataModelWithConfigOfTypeHttpAction)
	t.Run("should map state to data model with http script config", ut.shouldMapStateToDataModelWithConfigOfTypeHttpScript)
	t.Run("should return errror when trying to map state to model when no configuration is provided", ut.shouldReturnErrorWhenTryingToMapStateToModelWhenNoConfigurationIsProvided)
}

const (
	syntheticTestDefinition    = "instana_synthetic_test.example"
	syntheticTestConfigPattern = "%s.0.%s"

	syntheticTestID        = "id"
	syntheticTestLabel     = "label"
	syntheticTestActive    = true
	syntheticTestUrl       = "https://example.com"
	syntheticTestOperation = "GET"
)

var syntheticTestHttpActionTestCheckFunctions = []resource.TestCheckFunc{
	resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(syntheticTestConfigPattern, SyntheticTestFieldConfigHttpAction, SyntheticTestFieldConfigMarkSyntheticCall), "true"),
	resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(syntheticTestConfigPattern, SyntheticTestFieldConfigHttpAction, SyntheticTestFieldConfigRetries), "0"),
	resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(syntheticTestConfigPattern, SyntheticTestFieldConfigHttpAction, SyntheticTestFieldConfigRetryInterval), "1"),
	resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(syntheticTestConfigPattern, SyntheticTestFieldConfigHttpAction, SyntheticTestFieldConfigTimeout), "3m"),
	resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(syntheticTestConfigPattern, SyntheticTestFieldConfigHttpAction, SyntheticTestFieldConfigUrl), syntheticTestUrl),
	resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(syntheticTestConfigPattern, SyntheticTestFieldConfigHttpAction, SyntheticTestFieldConfigOperation), syntheticTestOperation),
	resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(syntheticTestConfigPattern, SyntheticTestFieldConfigHttpAction, SyntheticTestFieldConfigBody), "expected_body"),
	resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(syntheticTestConfigPattern, SyntheticTestFieldConfigHttpAction, SyntheticTestFieldConfigValidationString), "expected_body"),
	resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(syntheticTestConfigPattern, SyntheticTestFieldConfigHttpAction, SyntheticTestFieldConfigFollowRedirect), "false"),
	resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(syntheticTestConfigPattern, SyntheticTestFieldConfigHttpAction, SyntheticTestFieldConfigAllowInsecure), "true"),
	resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(syntheticTestConfigPattern, SyntheticTestFieldConfigHttpAction, SyntheticTestFieldConfigExpectStatus), "201"),
	resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(syntheticTestConfigPattern, SyntheticTestFieldConfigHttpAction, SyntheticTestFieldConfigExpectMatch), "[a-zA-Z]"),
}

func syntheticTestHttpActionWithoutApplicationIdIntegrationTest() *syntheticTestResourceIntegrationTest {
	const terraformTemplate = `
resource "instana_synthetic_test" "example" {
	label          = "label %d"
	active         = true
	locations      = ["location-id"]
	test_frequency = 10
	playback_mode  = "Staggered"

	http_action {
		mark_synthetic_call = true
		retries             = 0
		retry_interval      = 1
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

	const serverResponseTemplate = `
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
	return newSyntheticTestIntegrationTest(terraformTemplate, serverResponseTemplate, syntheticTestHttpActionTestCheckFunctions)
}

func syntheticTestHttpActionWithApplicationIdIntegrationTest() *syntheticTestResourceIntegrationTest {
	const terraformTemplate = `
resource "instana_synthetic_test" "example" {
	label          = "label %d"
	active         = true
	application_id = "application-id"
	locations      = ["location-id"]
	test_frequency = 10
	playback_mode  = "Staggered"

	http_action {
		mark_synthetic_call = true
		retries             = 0
		retry_interval      = 1
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

	const serverResponseTemplate = `
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
	checks := append(syntheticTestHttpActionTestCheckFunctions, resource.TestCheckResourceAttr(syntheticTestDefinition, SyntheticTestFieldApplicationID, "application-id"))

	return newSyntheticTestIntegrationTest(terraformTemplate, serverResponseTemplate, checks)
}

func syntheticTestHttpScriptIntegrationTest() *syntheticTestResourceIntegrationTest {
	const terraformTemplate = `
resource "instana_synthetic_test" "example" {
	label          = "label %d"
	active         = true
	locations      = ["location-id"]
	test_frequency = 10
	playback_mode  = "Staggered"

	http_script {
		mark_synthetic_call = true
		retries             = 0
		retry_interval      = 1
		timeout             = "3m"
		script              = "my-script"
	}

	custom_properties = {
		"key1" = "val1"
		"key2" = "val2"
	}
}
`

	const serverResponseTemplate = `
{
    "id": "%s",
    "label": "label %d",
    "active": true,
    "locations": ["location-id"],
    "testFrequency": 10,
    "playbackMode": "Staggered",

    "configuration": {
        "syntheticType": "HTTPScript",
        "markSyntheticCall": true,
		"retryInterval": 1,
		"timeout": "3m",
		"script": "my-script"
    },

	"customProperties": {
		"key1": "val1",
		"key2": "val2"
	}
}
`
	var checks = []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(syntheticTestConfigPattern, SyntheticTestFieldConfigHttpScript, SyntheticTestFieldConfigMarkSyntheticCall), "true"),
		resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(syntheticTestConfigPattern, SyntheticTestFieldConfigHttpScript, SyntheticTestFieldConfigRetries), "0"),
		resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(syntheticTestConfigPattern, SyntheticTestFieldConfigHttpScript, SyntheticTestFieldConfigRetryInterval), "1"),
		resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(syntheticTestConfigPattern, SyntheticTestFieldConfigHttpScript, SyntheticTestFieldConfigTimeout), "3m"),
		resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf(syntheticTestConfigPattern, SyntheticTestFieldConfigHttpScript, SyntheticTestFieldConfigScript), "my-script"),
	}
	return newSyntheticTestIntegrationTest(terraformTemplate, serverResponseTemplate, checks)
}

func newSyntheticTestIntegrationTest(resourceTemplate string, serverResponseTemplate string, useCaseSpecificChecks []resource.TestCheckFunc) *syntheticTestResourceIntegrationTest {
	return &syntheticTestResourceIntegrationTest{
		resourceTemplate:       resourceTemplate,
		resourceName:           syntheticTestDefinition,
		serverResponseTemplate: serverResponseTemplate,
		useCaseSpecificChecks:  useCaseSpecificChecks,
	}
}

type syntheticTestResourceIntegrationTest struct {
	resourceTemplate       string
	resourceName           string
	serverResponseTemplate string
	useCaseSpecificChecks  []resource.TestCheckFunc
}

func (it *syntheticTestResourceIntegrationTest) testCrud(t *testing.T) {
	id := RandomID()
	httpServer := it.createTestServerForSyntheticTestResource(id)
	defer httpServer.Close()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			it.createCheckFunctions(httpServer.GetPort(), 0),
			testStepImport(syntheticTestDefinition),
			it.createCheckFunctions(httpServer.GetPort(), 1),
			testStepImport(syntheticTestDefinition),
		},
	})
}

func (it *syntheticTestResourceIntegrationTest) createCheckFunctions(httpPort int64, iteration int) resource.TestStep {
	defaultChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(syntheticTestDefinition, "id"),
		resource.TestCheckResourceAttr(syntheticTestDefinition, SyntheticTestFieldLabel, fmt.Sprintf("%s %d", syntheticTestLabel, iteration)),
		resource.TestCheckResourceAttr(syntheticTestDefinition, SyntheticTestFieldActive, "true"),
		resource.TestCheckResourceAttr(syntheticTestDefinition, fmt.Sprintf("%s.0", SyntheticTestFieldLocations), "location-id"),
		resource.TestCheckResourceAttr(syntheticTestDefinition, SyntheticTestFieldTestFrequency, "10"),
		resource.TestCheckResourceAttr(syntheticTestDefinition, SyntheticTestFieldPlaybackMode, "Staggered"),
	}
	checks := append(defaultChecks, it.useCaseSpecificChecks...)

	return resource.TestStep{
		Config: appendProviderConfig(fmt.Sprintf(it.resourceTemplate, iteration), httpPort),
		Check:  resource.ComposeTestCheckFunc(checks...),
	}
}

func (it *syntheticTestResourceIntegrationTest) createTestServerForSyntheticTestResource(id string) testutils.TestHTTPServer {
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
		jsonData := fmt.Sprintf(it.serverResponseTemplate, id, puts)
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

type syntheticTestUnitTest struct{}

func (ut *syntheticTestUnitTest) resourceDefinitionShouldBeValid(t *testing.T) {
	resourceHandle := NewSyntheticTestResourceHandle()

	schemaMap := resourceHandle.MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	require.Len(t, schemaMap, 10)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SyntheticTestFieldLabel)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(SyntheticTestFieldDescription)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(SyntheticTestFieldActive, true)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(SyntheticTestFieldApplicationID)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeMapOfStrings(SyntheticTestFieldCustomProperties)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeSetOfStrings(SyntheticTestFieldLocations)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(SyntheticTestFieldPlaybackMode)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeInt(SyntheticTestFieldTestFrequency)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfResource(SyntheticTestFieldConfigHttpAction)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeListOfResource(SyntheticTestFieldConfigHttpScript)

	httpActionSchema := schemaMap[SyntheticTestFieldConfigHttpAction].Elem.(*schema.Resource).Schema
	ut.verifyHttpActionSchema(t, httpActionSchema)
	httpScriptSchema := schemaMap[SyntheticTestFieldConfigHttpScript].Elem.(*schema.Resource).Schema
	ut.verifyHttpScriptSchema(t, httpScriptSchema)
}

func (ut *syntheticTestUnitTest) verifyHttpActionSchema(t *testing.T, schemaMap map[string]*schema.Schema) {
	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	require.Len(t, schemaMap, 13)
	ut.verifyCommonConfigurationFields(schemaAssert)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(SyntheticTestFieldConfigUrl)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(SyntheticTestFieldConfigOperation)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeMapOfStrings(SyntheticTestFieldConfigHeaders)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(SyntheticTestFieldConfigBody)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(SyntheticTestFieldConfigValidationString)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(SyntheticTestFieldConfigFollowRedirect, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(SyntheticTestFieldConfigAllowInsecure, false)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeInt(SyntheticTestFieldConfigExpectStatus)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(SyntheticTestFieldConfigExpectMatch)
}

func (ut *syntheticTestUnitTest) verifyHttpScriptSchema(t *testing.T, schemaMap map[string]*schema.Schema) {
	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	require.Len(t, schemaMap, 5)
	ut.verifyCommonConfigurationFields(schemaAssert)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SyntheticTestFieldConfigScript)
}

func (ut *syntheticTestUnitTest) verifyCommonConfigurationFields(schemaAssert testutils.TerraformSchemaAssert) {
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(SyntheticTestFieldConfigMarkSyntheticCall, false)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeInt(SyntheticTestFieldConfigRetries)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeInt(SyntheticTestFieldConfigRetryInterval)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(SyntheticTestFieldConfigTimeout)
}

func (ut *syntheticTestUnitTest) shouldReturnCorrectResourceName(t *testing.T) {
	name := NewSyntheticTestResourceHandle().MetaData().ResourceName

	require.Equal(t, "instana_synthetic_test", name)
}

func (ut *syntheticTestUnitTest) shouldHaveSchemaVersionZero(t *testing.T) {
	require.Equal(t, 0, NewSyntheticTestResourceHandle().MetaData().SchemaVersion)
}

func (ut *syntheticTestUnitTest) shouldHaveNoStateUpgrader(t *testing.T) {
	require.Len(t, NewSyntheticTestResourceHandle().StateUpgraders(), 0)
}

func (ut *syntheticTestUnitTest) shouldUpdateResourceStateForHttpScript(t *testing.T) {
	testHelper := NewTestHelper[*restapi.SyntheticTest](t)
	resourceHandle := NewSyntheticTestResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	description := "my description"
	applicationID := "application-id"
	customProperties := map[string]interface{}{
		"p1": "v1",
		"p2": "v2",
	}
	timeout := "20s"
	script := "my-test-script"
	testFrequency := int32(2)
	data := restapi.SyntheticTest{
		ID:            syntheticTestID,
		Label:         syntheticTestLabel,
		Active:        syntheticTestActive,
		Description:   &description,
		ApplicationID: &applicationID,
		Configuration: restapi.SyntheticTestConfig{
			SyntheticType:     SyntheticCheckTypeHttpScript,
			MarkSyntheticCall: true,
			Retries:           5,
			RetryInterval:     10,
			Timeout:           &timeout,
			Script:            &script,
		},
		CustomProperties: customProperties,
		Locations:        []string{"loc1", "loc2"},
		PlaybackMode:     "Staggered",
		TestFrequency:    &testFrequency,
	}

	err := resourceHandle.UpdateState(resourceData, &data)

	require.Nil(t, err)

	require.Nil(t, err)
	require.Equal(t, syntheticTestID, resourceData.Id())
	require.Equal(t, syntheticTestActive, resourceData.Get(SyntheticTestFieldActive))
	require.Equal(t, syntheticTestLabel, resourceData.Get(SyntheticTestFieldLabel))
	require.Equal(t, description, resourceData.Get(SyntheticTestFieldDescription))
	require.Equal(t, applicationID, resourceData.Get(SyntheticTestFieldApplicationID))
	require.Equal(t, customProperties, resourceData.Get(SyntheticTestFieldCustomProperties))
	require.Equal(t, []interface{}{"loc1", "loc2"}, resourceData.Get(SyntheticTestFieldLocations).(*schema.Set).List())
	require.Equal(t, "Staggered", resourceData.Get(SyntheticTestFieldPlaybackMode))
	require.Equal(t, testFrequency, int32(resourceData.Get(SyntheticTestFieldTestFrequency).(int)))

	require.IsType(t, []interface{}{}, resourceData.Get(SyntheticTestFieldConfigHttpAction))
	require.Len(t, resourceData.Get(SyntheticTestFieldConfigHttpAction).([]interface{}), 0)
	require.IsType(t, []interface{}{}, resourceData.Get(SyntheticTestFieldConfigHttpScript))

	httpScriptConfigs := resourceData.Get(SyntheticTestFieldConfigHttpScript).([]interface{})
	require.Len(t, httpScriptConfigs, 1)
	require.IsType(t, map[string]interface{}{}, httpScriptConfigs[0])

	httpScriptConfig := httpScriptConfigs[0].(map[string]interface{})
	require.Len(t, httpScriptConfig, 5)
	require.Equal(t, true, httpScriptConfig[SyntheticTestFieldConfigMarkSyntheticCall])
	require.Equal(t, 5, httpScriptConfig[SyntheticTestFieldConfigRetries])
	require.Equal(t, 10, httpScriptConfig[SyntheticTestFieldConfigRetryInterval])
	require.Equal(t, timeout, httpScriptConfig[SyntheticTestFieldConfigTimeout])
	require.Equal(t, script, httpScriptConfig[SyntheticTestFieldConfigScript])
}

func (ut *syntheticTestUnitTest) shouldUpdateResourceStateForHttpAction(t *testing.T) {
	testHelper := NewTestHelper[*restapi.SyntheticTest](t)
	resourceHandle := NewSyntheticTestResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	description := "my description"
	applicationID := "application-id"
	customProperties := map[string]interface{}{
		"p1": "v1",
		"p2": "v2",
	}
	timeout := "20s"
	url := "https://app.example.com/health"
	operation := "GET"
	body := "payload"
	validationString := "validation-string"
	followRedirect := false
	allowInsecure := false
	expectStatus := int32(200)
	expectMatch := "expected-match"
	testFrequency := int32(2)
	data := restapi.SyntheticTest{
		ID:            syntheticTestID,
		Label:         syntheticTestLabel,
		Active:        syntheticTestActive,
		Description:   &description,
		ApplicationID: &applicationID,
		Configuration: restapi.SyntheticTestConfig{
			SyntheticType:     SyntheticCheckTypeHttpAction,
			MarkSyntheticCall: true,
			Retries:           5,
			RetryInterval:     10,
			Timeout:           &timeout,
			URL:               &url,
			Operation:         &operation,
			Headers: map[string]interface{}{
				"h1": "v1",
				"h2": "v2",
			},
			Body:             &body,
			ValidationString: &validationString,
			FollowRedirect:   &followRedirect,
			AllowInsecure:    &allowInsecure,
			ExpectStatus:     &expectStatus,
			ExpectMatch:      &expectMatch,
		},
		CustomProperties: customProperties,
		Locations:        []string{"loc1", "loc2"},
		PlaybackMode:     "Staggered",
		TestFrequency:    &testFrequency,
	}

	err := resourceHandle.UpdateState(resourceData, &data)

	require.Nil(t, err)

	require.Nil(t, err)
	require.Equal(t, syntheticTestID, resourceData.Id())
	require.Equal(t, syntheticTestActive, resourceData.Get(SyntheticTestFieldActive))
	require.Equal(t, syntheticTestLabel, resourceData.Get(SyntheticTestFieldLabel))
	require.Equal(t, description, resourceData.Get(SyntheticTestFieldDescription))
	require.Equal(t, applicationID, resourceData.Get(SyntheticTestFieldApplicationID))
	require.Equal(t, customProperties, resourceData.Get(SyntheticTestFieldCustomProperties))
	require.Equal(t, []interface{}{"loc1", "loc2"}, resourceData.Get(SyntheticTestFieldLocations).(*schema.Set).List())
	require.Equal(t, "Staggered", resourceData.Get(SyntheticTestFieldPlaybackMode))
	require.Equal(t, testFrequency, int32(resourceData.Get(SyntheticTestFieldTestFrequency).(int)))

	require.IsType(t, []interface{}{}, resourceData.Get(SyntheticTestFieldConfigHttpScript))
	require.Len(t, resourceData.Get(SyntheticTestFieldConfigHttpScript).([]interface{}), 0)

	require.IsType(t, []interface{}{}, resourceData.Get(SyntheticTestFieldConfigHttpAction))
	httpActionConfigs := resourceData.Get(SyntheticTestFieldConfigHttpAction).([]interface{})
	require.Len(t, httpActionConfigs, 1)
	require.IsType(t, map[string]interface{}{}, httpActionConfigs[0])

	httpActionConfig := httpActionConfigs[0].(map[string]interface{})
	require.Len(t, httpActionConfig, 13)
	require.Equal(t, true, httpActionConfig[SyntheticTestFieldConfigMarkSyntheticCall])
	require.Equal(t, 5, httpActionConfig[SyntheticTestFieldConfigRetries])
	require.Equal(t, 10, httpActionConfig[SyntheticTestFieldConfigRetryInterval])
	require.Equal(t, timeout, httpActionConfig[SyntheticTestFieldConfigTimeout])
	require.Equal(t, url, httpActionConfig[SyntheticTestFieldConfigUrl])
	require.Equal(t, operation, httpActionConfig[SyntheticTestFieldConfigOperation])
	require.Equal(t, map[string]interface{}{
		"h1": "v1",
		"h2": "v2",
	}, httpActionConfig[SyntheticTestFieldConfigHeaders])
	require.Equal(t, body, httpActionConfig[SyntheticTestFieldConfigBody])
	require.Equal(t, validationString, httpActionConfig[SyntheticTestFieldConfigValidationString])
	require.Equal(t, followRedirect, httpActionConfig[SyntheticTestFieldConfigFollowRedirect])
	require.Equal(t, allowInsecure, httpActionConfig[SyntheticTestFieldConfigAllowInsecure])
	require.Equal(t, expectStatus, int32(httpActionConfig[SyntheticTestFieldConfigExpectStatus].(int)))
	require.Equal(t, expectMatch, httpActionConfig[SyntheticTestFieldConfigExpectMatch])
}

func (ut *syntheticTestUnitTest) shouldReturnErrorWhenTryingToUpdateStateAndConfigTypeIsNotSupported(t *testing.T) {
	testHelper := NewTestHelper[*restapi.SyntheticTest](t)
	resourceHandle := NewSyntheticTestResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	data := restapi.SyntheticTest{
		Configuration: restapi.SyntheticTestConfig{
			SyntheticType: "invalid",
		},
	}

	err := resourceHandle.UpdateState(resourceData, &data)

	require.Error(t, err)
	require.ErrorContains(t, err, "unsupported synthetic test of type invalid received")
}

func (ut *syntheticTestUnitTest) shouldMapStateToDataModelWithConfigOfTypeHttpAction(t *testing.T) {
	testHelper := NewTestHelper[*restapi.SyntheticTest](t)
	resourceHandle := NewSyntheticTestResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	description := "my description"
	applicationID := "application-id"
	customProperties := map[string]interface{}{
		"p1": "v1",
		"p2": "v2",
	}
	timeout := "20s"
	url := "https://app.example.com/health"
	operation := "GET"
	headers := map[string]interface{}{
		"h1": "v1",
		"h2": "v2",
	}
	body := "payload"
	validationString := "validation-string"
	followRedirect := false
	allowInsecure := false
	expectStatus := int32(200)
	expectMatch := "expected-match"
	testFrequency := int32(2)
	resourceData.SetId(syntheticTestID)
	setValueOnResourceData(t, resourceData, SyntheticTestFieldActive, syntheticTestActive)
	setValueOnResourceData(t, resourceData, SyntheticTestFieldLabel, syntheticTestLabel)
	setValueOnResourceData(t, resourceData, SyntheticTestFieldDescription, description)
	setValueOnResourceData(t, resourceData, SyntheticTestFieldApplicationID, applicationID)
	setValueOnResourceData(t, resourceData, SyntheticTestFieldCustomProperties, customProperties)
	setValueOnResourceData(t, resourceData, SyntheticTestFieldLocations, []interface{}{"loc1", "loc2"})
	setValueOnResourceData(t, resourceData, SyntheticTestFieldPlaybackMode, "Staggered")
	setValueOnResourceData(t, resourceData, SyntheticTestFieldTestFrequency, testFrequency)
	setValueOnResourceData(t, resourceData, SyntheticTestFieldConfigHttpAction, []interface{}{
		map[string]interface{}{
			SyntheticTestFieldConfigMarkSyntheticCall: true,
			SyntheticTestFieldConfigRetries:           5,
			SyntheticTestFieldConfigRetryInterval:     10,
			SyntheticTestFieldConfigTimeout:           timeout,
			SyntheticTestFieldConfigUrl:               url,
			SyntheticTestFieldConfigOperation:         operation,
			SyntheticTestFieldConfigHeaders:           headers,
			SyntheticTestFieldConfigBody:              body,
			SyntheticTestFieldConfigValidationString:  validationString,
			SyntheticTestFieldConfigFollowRedirect:    false,
			SyntheticTestFieldConfigAllowInsecure:     false,
			SyntheticTestFieldConfigExpectMatch:       expectMatch,
			SyntheticTestFieldConfigExpectStatus:      expectStatus,
		},
	})
	setValueOnResourceData(t, resourceData, SyntheticTestFieldConfigHttpScript, []interface{}{})

	model, err := resourceHandle.MapStateToDataObject(resourceData)

	require.Nil(t, err)
	require.IsType(t, &restapi.SyntheticTest{
		ID:            syntheticTestID,
		Label:         syntheticTestLabel,
		Active:        syntheticTestActive,
		Description:   &description,
		ApplicationID: &applicationID,
		Configuration: restapi.SyntheticTestConfig{
			SyntheticType:     SyntheticCheckTypeHttpAction,
			MarkSyntheticCall: true,
			Retries:           5,
			RetryInterval:     10,
			Timeout:           &timeout,
			URL:               &url,
			Operation:         &operation,
			Headers:           headers,
			Body:              &body,
			ValidationString:  &validationString,
			FollowRedirect:    &followRedirect,
			AllowInsecure:     &allowInsecure,
			ExpectStatus:      &expectStatus,
			ExpectMatch:       &expectMatch,
		},
		CustomProperties: customProperties,
		Locations:        []string{"loc1", "loc2"},
		PlaybackMode:     "Staggered",
		TestFrequency:    &testFrequency,
	}, model)
}

func (ut *syntheticTestUnitTest) shouldMapStateToDataModelWithConfigOfTypeHttpScript(t *testing.T) {
	testHelper := NewTestHelper[*restapi.SyntheticTest](t)
	resourceHandle := NewSyntheticTestResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	description := "my description"
	applicationID := "application-id"
	customProperties := map[string]interface{}{
		"p1": "v1",
		"p2": "v2",
	}
	timeout := "20s"
	testFrequency := int32(2)
	script := "my-script"
	resourceData.SetId(syntheticTestID)
	setValueOnResourceData(t, resourceData, SyntheticTestFieldActive, syntheticTestActive)
	setValueOnResourceData(t, resourceData, SyntheticTestFieldLabel, syntheticTestLabel)
	setValueOnResourceData(t, resourceData, SyntheticTestFieldDescription, description)
	setValueOnResourceData(t, resourceData, SyntheticTestFieldApplicationID, applicationID)
	setValueOnResourceData(t, resourceData, SyntheticTestFieldCustomProperties, customProperties)
	setValueOnResourceData(t, resourceData, SyntheticTestFieldLocations, []interface{}{"loc1", "loc2"})
	setValueOnResourceData(t, resourceData, SyntheticTestFieldPlaybackMode, "Staggered")
	setValueOnResourceData(t, resourceData, SyntheticTestFieldTestFrequency, testFrequency)
	setValueOnResourceData(t, resourceData, SyntheticTestFieldConfigHttpScript, []interface{}{
		map[string]interface{}{
			SyntheticTestFieldConfigMarkSyntheticCall: true,
			SyntheticTestFieldConfigRetries:           5,
			SyntheticTestFieldConfigRetryInterval:     10,
			SyntheticTestFieldConfigTimeout:           timeout,
			SyntheticTestFieldConfigScript:            script,
		},
	})
	setValueOnResourceData(t, resourceData, SyntheticTestFieldConfigHttpAction, []interface{}{})

	model, err := resourceHandle.MapStateToDataObject(resourceData)

	require.Nil(t, err)
	require.IsType(t, &restapi.SyntheticTest{
		ID:            syntheticTestID,
		Label:         syntheticTestLabel,
		Active:        syntheticTestActive,
		Description:   &description,
		ApplicationID: &applicationID,
		Configuration: restapi.SyntheticTestConfig{
			SyntheticType:     SyntheticCheckTypeHttpScript,
			MarkSyntheticCall: true,
			Retries:           5,
			RetryInterval:     10,
			Timeout:           &timeout,
			Script:            &script,
		},
		CustomProperties: customProperties,
		Locations:        []string{"loc1", "loc2"},
		PlaybackMode:     "Staggered",
		TestFrequency:    &testFrequency,
	}, model)
}

func (ut *syntheticTestUnitTest) shouldReturnErrorWhenTryingToMapStateToModelWhenNoConfigurationIsProvided(t *testing.T) {
	testHelper := NewTestHelper[*restapi.SyntheticTest](t)
	resourceHandle := NewSyntheticTestResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	resourceData.SetId(syntheticTestID)
	setValueOnResourceData(t, resourceData, SyntheticTestFieldConfigHttpScript, []interface{}{})
	setValueOnResourceData(t, resourceData, SyntheticTestFieldConfigHttpAction, []interface{}{})

	_, err := resourceHandle.MapStateToDataObject(resourceData)

	require.Error(t, err)
	require.ErrorContains(t, err, "no supported synthetic test configuration provided")
}
