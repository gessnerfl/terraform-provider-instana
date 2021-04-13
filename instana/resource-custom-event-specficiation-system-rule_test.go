package instana_test

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

const resourceCustomEventSpecificationWithSystemRuleDefinitionTemplate = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:{{PORT}}"
}

resource "instana_custom_event_spec_system_rule" "example" {
  name = "name"
  query = "query"
  enabled = true
  triggering = true
  description = "description"
  expiration_time = 60000
  rule_severity = "warning"
  rule_system_rule_id = "system-rule-id"
}
`

const (
	customSystemEventApiPath                             = restapi.CustomEventSpecificationResourcePath + "/{id}"
	testCustomEventSpecificationWithSystemRuleDefinition = "instana_custom_event_spec_system_rule.example"

	customSystemEventID               = "custom-system-event-id"
	customSystemEventName             = "name"
	customSystemEventFullName         = "prefix name suffix"
	customSystemEventQuery            = "query"
	customSystemEventExpirationTime   = 60000
	customSystemEventDescription      = "description"
	customSystemEventRuleSystemRuleId = "system-rule-id"

	constSystemEventContentType = "Content-Type"
)

var customSystemEventRuleSeverity = restapi.SeverityWarning.GetTerraformRepresentation()

func TestCRUDOfCreateResourceCustomEventSpecificationWithSystemdRuleResourceWithMockServer(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, customSystemEventApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, customSystemEventApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, customSystemEventApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		json := strings.ReplaceAll(`
		{
			"id" : "{{id}}",
			"name" : "name",
			"entityType" : "entity_type",
			"query" : "query",
			"enabled" : true,
			"triggering" : true,
			"description" : "description",
			"expirationTime" : 60000,
			"rules" : [ { "ruleType" : "system", "severity" : 5, "systemRuleId" : "system-rule-id" } ]
		}
		`, "{{id}}", vars["id"])
		w.Header().Set(constSystemEventContentType, r.Header.Get(constSystemEventContentType))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceCustomEventSpecificationWithSystemRuleDefinition := strings.ReplaceAll(resourceCustomEventSpecificationWithSystemRuleDefinitionTemplate, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))

	resource.ParallelTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceCustomEventSpecificationWithSystemRuleDefinition,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testCustomEventSpecificationWithSystemRuleDefinition, "id"),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithSystemRuleDefinition, CustomEventSpecificationFieldName, customSystemEventName),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithSystemRuleDefinition, CustomEventSpecificationFieldEntityType, SystemRuleEntityType),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithSystemRuleDefinition, CustomEventSpecificationFieldQuery, customSystemEventQuery),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithSystemRuleDefinition, CustomEventSpecificationFieldTriggering, "true"),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithSystemRuleDefinition, CustomEventSpecificationFieldDescription, customSystemEventDescription),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithSystemRuleDefinition, CustomEventSpecificationFieldExpirationTime, strconv.Itoa(customSystemEventExpirationTime)),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithSystemRuleDefinition, CustomEventSpecificationFieldEnabled, "true"),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithSystemRuleDefinition, CustomEventSpecificationRuleSeverity, customSystemEventRuleSeverity),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithSystemRuleDefinition, SystemRuleSpecificationSystemRuleID, customSystemEventRuleSystemRuleId),
				),
			},
			{
				ResourceName:      testApplicationConfigDefinition,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestCustomEventSpecificationWithSystemRuleSchemaDefinitionIsValid(t *testing.T) {
	schema := NewCustomEventSpecificationWithSystemRuleResourceHandle().MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationFieldName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(CustomEventSpecificationFieldFullName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(CustomEventSpecificationFieldEntityType)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(CustomEventSpecificationFieldQuery)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(CustomEventSpecificationFieldTriggering, false)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(CustomEventSpecificationFieldDescription)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeInt(CustomEventSpecificationFieldExpirationTime)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(CustomEventSpecificationFieldEnabled, true)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(CustomEventSpecificationRuleSeverity)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(SystemRuleSpecificationSystemRuleID)
}

func TestCustomEventSpecificationWithSystemRuleResourceShouldHaveSchemaVersionOne(t *testing.T) {
	require.Equal(t, 2, NewCustomEventSpecificationWithSystemRuleResourceHandle().MetaData().SchemaVersion)
}

func TestCustomEventSpecificationWithSystemRuleShouldHaveOneStateUpgraderForVersionZero(t *testing.T) {
	resourceHandler := NewCustomEventSpecificationWithSystemRuleResourceHandle()

	require.Equal(t, 2, len(resourceHandler.StateUpgraders()))
	require.Equal(t, 0, resourceHandler.StateUpgraders()[0].Version)
	require.Equal(t, 1, resourceHandler.StateUpgraders()[1].Version)
}

func TestShouldMigrateCustomEventSpecificationWithSystemRuleStateAndAddFullNameWithSameValueAsNameWhenMigratingFromVersion0To1(t *testing.T) {
	name := "Test Name"
	rawData := make(map[string]interface{})
	rawData[CustomEventSpecificationFieldName] = name
	meta := "dummy"
	ctx := context.Background()

	result, err := NewCustomEventSpecificationWithSystemRuleResourceHandle().StateUpgraders()[0].Upgrade(ctx, rawData, meta)

	require.Nil(t, err)
	require.Equal(t, name, result[CustomEventSpecificationFieldFullName])
}

func TestShouldMigrateEmptyCustomEventSpecificationWithSystemRuleStateFromVersion0To1(t *testing.T) {
	rawData := make(map[string]interface{})
	meta := "dummy"
	ctx := context.Background()

	result, err := NewCustomEventSpecificationWithSystemRuleResourceHandle().StateUpgraders()[0].Upgrade(ctx, rawData, meta)

	require.Nil(t, err)
	require.Nil(t, result[CustomEventSpecificationFieldFullName])
}

func TestShouldMigrateCustomEventSpecificationWithSystemRuleStateToVersion2WhenDownstreamConfigurationIsProvided(t *testing.T) {
	rawData := make(map[string]interface{})
	rawData["downstream_integration_ids"] = []interface{}{"id1", "id2"}
	rawData["downstream_broadcast_to_all_alerting_configs"] = true
	meta := "dummy"
	ctx := context.Background()

	result, err := NewCustomEventSpecificationWithSystemRuleResourceHandle().StateUpgraders()[1].Upgrade(ctx, rawData, meta)

	require.Nil(t, err)
	require.Nil(t, result["downstream_integration_ids"])
	require.Nil(t, result["downstream_broadcast_to_all_alerting_configs"])
}

func TestShouldMigrateCustomEventSpecificationWithSystemRuleStateToVersion2WhenNoDownstreamConfigurationIsProvided(t *testing.T) {
	rawData := make(map[string]interface{})
	meta := "dummy"
	ctx := context.Background()

	result, err := NewCustomEventSpecificationWithSystemRuleResourceHandle().StateUpgraders()[0].Upgrade(ctx, rawData, meta)

	require.Nil(t, err)
	require.Nil(t, result["downstream_integration_ids"])
	require.Nil(t, result["downstream_broadcast_to_all_alerting_configs"])
}

func TestShouldReturnCorrectResourceNameForCustomEventSpecificationWithSystemRuleResource(t *testing.T) {
	name := NewCustomEventSpecificationWithSystemRuleResourceHandle().MetaData().ResourceName

	require.Equal(t, name, "instana_custom_event_spec_system_rule")
}

func TestShouldUpdateCustomEventSpecificationWithSystemRuleTerraformStateFromApiObject(t *testing.T) {
	description := customSystemEventDescription
	expirationTime := customSystemEventExpirationTime
	query := customSystemEventQuery
	spec := restapi.CustomEventSpecification{
		ID:             customSystemEventID,
		Name:           customSystemEventFullName,
		EntityType:     SystemRuleEntityType,
		Query:          &query,
		Description:    &description,
		ExpirationTime: &expirationTime,
		Triggering:     true,
		Enabled:        true,
		Rules: []restapi.RuleSpecification{
			restapi.NewSystemRuleSpecification(customSystemEventRuleSystemRuleId, restapi.SeverityWarning.GetAPIRepresentation()),
		},
	}

	testHelper := NewTestHelper(t)
	sut := NewCustomEventSpecificationWithSystemRuleResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, &spec, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.Equal(t, customSystemEventID, resourceData.Id())
	require.Equal(t, customSystemEventName, resourceData.Get(CustomEventSpecificationFieldName))
	require.Equal(t, customSystemEventFullName, resourceData.Get(CustomEventSpecificationFieldFullName))
	require.Equal(t, SystemRuleEntityType, resourceData.Get(CustomEventSpecificationFieldEntityType))
	require.Equal(t, customSystemEventQuery, resourceData.Get(CustomEventSpecificationFieldQuery))
	require.Equal(t, description, resourceData.Get(CustomEventSpecificationFieldDescription))
	require.True(t, resourceData.Get(CustomEventSpecificationFieldTriggering).(bool))
	require.True(t, resourceData.Get(CustomEventSpecificationFieldEnabled).(bool))

	require.Equal(t, customSystemEventRuleSystemRuleId, resourceData.Get(SystemRuleSpecificationSystemRuleID))
	require.Equal(t, restapi.SeverityWarning.GetTerraformRepresentation(), resourceData.Get(CustomEventSpecificationRuleSeverity))
}

func TestShouldSuccessfullyConvertCustomEventSpecificationWithSystemRuleStateToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewCustomEventSpecificationWithSystemRuleResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	resourceData.SetId(customSystemEventID)
	resourceData.Set(CustomEventSpecificationFieldFullName, customSystemEventName)
	resourceData.Set(CustomEventSpecificationFieldEntityType, SystemRuleEntityType)
	resourceData.Set(CustomEventSpecificationFieldQuery, customSystemEventQuery)
	resourceData.Set(CustomEventSpecificationFieldTriggering, true)
	resourceData.Set(CustomEventSpecificationFieldDescription, customSystemEventDescription)
	resourceData.Set(CustomEventSpecificationFieldExpirationTime, customSystemEventExpirationTime)
	resourceData.Set(CustomEventSpecificationFieldEnabled, true)
	resourceData.Set(CustomEventSpecificationRuleSeverity, customSystemEventRuleSeverity)
	resourceData.Set(SystemRuleSpecificationSystemRuleID, customSystemEventRuleSystemRuleId)

	result, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.IsType(t, &restapi.CustomEventSpecification{}, result)
	customEventSpec := result.(*restapi.CustomEventSpecification)
	require.Equal(t, customSystemEventID, customEventSpec.GetIDForResourcePath())
	require.Equal(t, customSystemEventName, customEventSpec.Name)
	require.Equal(t, SystemRuleEntityType, customEventSpec.EntityType)
	require.Equal(t, customSystemEventQuery, *customEventSpec.Query)
	require.Equal(t, customSystemEventDescription, *customEventSpec.Description)
	require.Equal(t, customSystemEventExpirationTime, *customEventSpec.ExpirationTime)
	require.True(t, customEventSpec.Triggering)
	require.True(t, customEventSpec.Enabled)

	require.Equal(t, 1, len(customEventSpec.Rules))
	require.Equal(t, customSystemEventRuleSystemRuleId, *customEventSpec.Rules[0].SystemRuleID)
	require.Equal(t, restapi.SeverityWarning.GetAPIRepresentation(), customEventSpec.Rules[0].Severity)
}

func TestShouldFailToConvertCustomEventSpecificationWithSystemRuleStateToDataModelWhenSeverityIsNotValid(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewCustomEventSpecificationWithSystemRuleResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.Set(CustomEventSpecificationRuleSeverity, "INVALID")

	_, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.NotNil(t, err)
}
