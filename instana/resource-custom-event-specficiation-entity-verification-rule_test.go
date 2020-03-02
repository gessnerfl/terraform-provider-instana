package instana_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/stretchr/testify/assert"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/gessnerfl/terraform-provider-instana/utils"
)

var testCustomEventSpecificationWithEntityVerificationRuleProviders = map[string]terraform.ResourceProvider{
	"instana": Provider(),
}

const resourceCustomEventSpecificationWithEntityVerificationRuleDefinitionTemplate = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:{{PORT}}"
}

resource "instana_custom_event_spec_entity_verification_rule" "example" {
  name = "name"
  query = "query"
  enabled = true
  triggering = true
  description = "description"
  expiration_time = 60000
  rule_severity = "warning"
  rule_matching_entity_type = "matching-entity-type"
  rule_matching_operator = "is"
  rule_matching_entity_label = "matching-entity-label"
  rule_offline_duration = 60000
}
`

const (
	customEntityVerificationEventApiPath                             = restapi.CustomEventSpecificationResourcePath + "/{id}"
	testCustomEventSpecificationWithEntityVerificationRuleDefinition = ResourceInstanaCustomEventSpecificationEntityVerificationRule + ".example"

	customEntityVerificationEventID                      = "custom-entity-verification-event-id"
	customEntityVerificationEventName                    = "name"
	customEntityVerificationEventQuery                   = "query"
	customEntityVerificationEventExpirationTime          = 60000
	customEntityVerificationEventDescription             = "description"
	customEntityVerificationEventRuleMatchingEntityLabel = "matching-entity-label"
	customEntityVerificationEventRuleMatchingEntityType  = "matching-entity-type"
	customEntityVerificationEventRuleMatchingOperator    = restapi.MatchingOperatorIs
	customEntityVerificationEventRuleOfflineDuration     = 60000

	customEntityVerificationEventMessageNotAValidSeverity           = "not a valid severity"
	customEntityVerificationEventTestMessageExpectedInvalidSeverity = "Expected to get error that the provided severity is not valid"

	constEntityVerificationEventContentType = "Content-Type"
)

var customEntityVerificationEventRuleSeverity = restapi.SeverityWarning.GetTerraformRepresentation()

func TestCRUDOfCreateResourceCustomEventSpecificationWithEntityVerificationRuleResourceWithMockServer(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, customEntityVerificationEventApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, customEntityVerificationEventApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, customEntityVerificationEventApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		json := strings.ReplaceAll(`
		{
			"id" : "{{id}}",
			"name" : "name",
			"query" : "query",
			"entityType" : "host",
			"enabled" : true,
			"triggering" : true,
			"description" : "description",
			"expirationTime" : 60000,
			"rules" : [ { "ruleType" : "entity_verification", "severity" : 5, "matchingEntityLabel" : "matching-entity-label", "matchingEntityType" : "matching-entity-type", "matchingOperator" : "is", "offlineDuration" : 60000 } ]
		}
		`, "{{id}}", vars["id"])
		w.Header().Set(constEntityVerificationEventContentType, r.Header.Get(constEntityVerificationEventContentType))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceCustomEventSpecificationWithEntityVerificationRuleDefinition := strings.ReplaceAll(resourceCustomEventSpecificationWithEntityVerificationRuleDefinitionTemplate, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))

	resource.UnitTest(t, resource.TestCase{
		Providers: testCustomEventSpecificationWithEntityVerificationRuleProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: resourceCustomEventSpecificationWithEntityVerificationRuleDefinition,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testCustomEventSpecificationWithEntityVerificationRuleDefinition, "id"),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationFieldName, customEntityVerificationEventName),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationFieldEntityType, EntityVerificationRuleEntityType),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationFieldQuery, customEntityVerificationEventQuery),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationFieldTriggering, "true"),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationFieldDescription, customEntityVerificationEventDescription),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationFieldExpirationTime, strconv.Itoa(customEntityVerificationEventExpirationTime)),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationFieldEnabled, "true"),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, CustomEventSpecificationRuleSeverity, customEntityVerificationEventRuleSeverity),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, EntityVerificationRuleFieldMatchingEntityLabel, customEntityVerificationEventRuleMatchingEntityLabel),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, EntityVerificationRuleFieldMatchingEntityType, customEntityVerificationEventRuleMatchingEntityType),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, EntityVerificationRuleFieldMatchingOperator, string(customEntityVerificationEventRuleMatchingOperator)),
					resource.TestCheckResourceAttr(testCustomEventSpecificationWithEntityVerificationRuleDefinition, EntityVerificationRuleFieldOfflineDuration, strconv.Itoa(customEntityVerificationEventRuleOfflineDuration)),
				),
			},
		},
	})
}

func TestCustomEventSpecificationWithEntityVerificationRuleSchemaDefinitionIsValid(t *testing.T) {
	schema := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle().Schema

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
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(EntityVerificationRuleFieldMatchingEntityLabel)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(EntityVerificationRuleFieldMatchingEntityType)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(EntityVerificationRuleFieldMatchingOperator)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeInt(EntityVerificationRuleFieldOfflineDuration)
}

func TestCustomEventSpecificationWithEntityVerificationRuleResourceShouldHaveSchemaVersionTwo(t *testing.T) {
	assert.Equal(t, 2, NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle().SchemaVersion)
}

func TestCustomEventSpecificationWithEntityVerificationRuleShouldHaveTwoStateUpgraderForVersionZeroAndOne(t *testing.T) {
	resourceHandler := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle()

	assert.Equal(t, 2, len(resourceHandler.StateUpgraders))
	assert.Equal(t, 0, resourceHandler.StateUpgraders[0].Version)
	assert.Equal(t, 1, resourceHandler.StateUpgraders[1].Version)
}

func TestShouldMigrateCustomEventSpecificationWithEntityVerificationRuleStateAndAddFullNameWithSameValueAsNameWhenMigratingFromVersion0To1(t *testing.T) {
	name := "Test Name"
	rawData := make(map[string]interface{})
	rawData[CustomEventSpecificationFieldName] = name
	meta := "dummy"

	result, err := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle().StateUpgraders[0].Upgrade(rawData, meta)

	assert.Nil(t, err)
	assert.Equal(t, name, result[CustomEventSpecificationFieldFullName])
}

func TestShouldMigrateEmptyCustomEventSpecificationWithEntityVerificationRuleStateFromVersion0To1(t *testing.T) {
	rawData := make(map[string]interface{})
	meta := "dummy"

	result, err := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle().StateUpgraders[1].Upgrade(rawData, meta)

	assert.Nil(t, err)
	assert.Nil(t, result[CustomEventSpecificationFieldFullName])
}

func TestShouldMigrateCustomEventSpecificationWithEntityVerificationRuleStateToVersion2WhenDownstreamConfigurationIsProvided(t *testing.T) {
	rawData := make(map[string]interface{})
	rawData["downstream_integration_ids"] = []interface{}{"id1", "id2"}
	rawData["downstream_broadcast_to_all_alerting_configs"] = true
	meta := "dummy"

	result, err := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle().StateUpgraders[1].Upgrade(rawData, meta)

	assert.Nil(t, err)
	assert.Nil(t, result["downstream_integration_ids"])
	assert.Nil(t, result["downstream_broadcast_to_all_alerting_configs"])
}

func TestShouldMigrateCustomEventSpecificationWithEntityVerificationRuleStateToVersion2WhenNoDownstreamConfigurationIsProvided(t *testing.T) {
	rawData := make(map[string]interface{})
	meta := "dummy"

	result, err := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle().StateUpgraders[0].Upgrade(rawData, meta)

	assert.Nil(t, err)
	assert.Nil(t, result["downstream_integration_ids"])
	assert.Nil(t, result["downstream_broadcast_to_all_alerting_configs"])
}

func TestShouldReturnCorrectResourceNameForCustomEventSpecificationWithEntityVerificationRuleResource(t *testing.T) {
	name := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle().ResourceName

	assert.Equal(t, name, "instana_custom_event_spec_entity_verification_rule")
}

func TestShouldUpdateCustomEventSpecificationWithEntityVerificationRuleTerraformStateFromApiObject(t *testing.T) {
	description := customEntityVerificationEventDescription
	expirationTime := customEntityVerificationEventExpirationTime
	query := customEntityVerificationEventQuery
	spec := restapi.CustomEventSpecification{
		ID:             customEntityVerificationEventID,
		Name:           customEntityVerificationEventName,
		EntityType:     EntityVerificationRuleEntityType,
		Query:          &query,
		Description:    &description,
		ExpirationTime: &expirationTime,
		Triggering:     true,
		Enabled:        true,
		Rules: []restapi.RuleSpecification{
			restapi.NewEntityVerificationRuleSpecification(customEntityVerificationEventRuleMatchingEntityLabel,
				customEntityVerificationEventRuleMatchingEntityType,
				customEntityVerificationEventRuleMatchingOperator,
				customEntityVerificationEventRuleOfflineDuration,
				restapi.SeverityWarning.GetAPIRepresentation()),
		},
	}

	testHelper := NewTestHelper(t)
	sut := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, spec)

	assert.Nil(t, err)
	assert.Equal(t, customEntityVerificationEventID, resourceData.Id())
	assert.Equal(t, customEntityVerificationEventName, resourceData.Get(CustomEventSpecificationFieldFullName))
	assert.Equal(t, EntityVerificationRuleEntityType, resourceData.Get(CustomEventSpecificationFieldEntityType))
	assert.Equal(t, customEntityVerificationEventQuery, resourceData.Get(CustomEventSpecificationFieldQuery))
	assert.Equal(t, description, resourceData.Get(CustomEventSpecificationFieldDescription))
	assert.True(t, resourceData.Get(CustomEventSpecificationFieldTriggering).(bool))
	assert.True(t, resourceData.Get(CustomEventSpecificationFieldEnabled).(bool))

	assert.Equal(t, customEntityVerificationEventRuleMatchingEntityLabel, resourceData.Get(EntityVerificationRuleFieldMatchingEntityLabel))
	assert.Equal(t, customEntityVerificationEventRuleMatchingEntityType, resourceData.Get(EntityVerificationRuleFieldMatchingEntityType))
	assert.Equal(t, string(customEntityVerificationEventRuleMatchingOperator), resourceData.Get(EntityVerificationRuleFieldMatchingOperator))
	assert.Equal(t, customEntityVerificationEventRuleOfflineDuration, resourceData.Get(EntityVerificationRuleFieldOfflineDuration))
	assert.Equal(t, restapi.SeverityWarning.GetTerraformRepresentation(), resourceData.Get(CustomEventSpecificationRuleSeverity))
}

func TestShouldSuccessfullyConvertCustomEventSpecificationWithEntityVerificationRuleStateToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	resourceData.SetId(customEntityVerificationEventID)
	resourceData.Set(CustomEventSpecificationFieldFullName, customEntityVerificationEventName)
	resourceData.Set(CustomEventSpecificationFieldEntityType, EntityVerificationRuleEntityType)
	resourceData.Set(CustomEventSpecificationFieldQuery, customEntityVerificationEventQuery)
	resourceData.Set(CustomEventSpecificationFieldTriggering, true)
	resourceData.Set(CustomEventSpecificationFieldDescription, customEntityVerificationEventDescription)
	resourceData.Set(CustomEventSpecificationFieldExpirationTime, customEntityVerificationEventExpirationTime)
	resourceData.Set(CustomEventSpecificationFieldEnabled, true)
	resourceData.Set(CustomEventSpecificationRuleSeverity, customEntityVerificationEventRuleSeverity)
	resourceData.Set(EntityVerificationRuleFieldMatchingEntityLabel, customEntityVerificationEventRuleMatchingEntityLabel)
	resourceData.Set(EntityVerificationRuleFieldMatchingEntityType, customEntityVerificationEventRuleMatchingEntityType)
	resourceData.Set(EntityVerificationRuleFieldMatchingOperator, string(customEntityVerificationEventRuleMatchingOperator))
	resourceData.Set(EntityVerificationRuleFieldOfflineDuration, customEntityVerificationEventRuleOfflineDuration)

	result, err := resourceHandle.MapStateToDataObject(resourceData, utils.NewResourceNameFormatter("prefix ", " suffix"))

	assert.Nil(t, err)
	assert.IsType(t, restapi.CustomEventSpecification{}, result)
	customEventSpec := result.(restapi.CustomEventSpecification)
	assert.Equal(t, customEntityVerificationEventID, customEventSpec.GetID())
	assert.Equal(t, customEntityVerificationEventName, customEventSpec.Name)
	assert.Equal(t, EntityVerificationRuleEntityType, customEventSpec.EntityType)
	assert.Equal(t, customEntityVerificationEventQuery, *customEventSpec.Query)
	assert.Equal(t, customEntityVerificationEventDescription, *customEventSpec.Description)
	assert.Equal(t, customEntityVerificationEventExpirationTime, *customEventSpec.ExpirationTime)
	assert.True(t, customEventSpec.Triggering)
	assert.True(t, customEventSpec.Enabled)

	assert.Equal(t, 1, len(customEventSpec.Rules))
	assert.Equal(t, customEntityVerificationEventRuleMatchingEntityLabel, *customEventSpec.Rules[0].MatchingEntityLabel)
	assert.Equal(t, customEntityVerificationEventRuleMatchingEntityType, *customEventSpec.Rules[0].MatchingEntityType)
	assert.Equal(t, customEntityVerificationEventRuleMatchingOperator, *customEventSpec.Rules[0].MatchingOperator)
	assert.Equal(t, customEntityVerificationEventRuleOfflineDuration, *customEventSpec.Rules[0].OfflineDuration)
	assert.Equal(t, restapi.SeverityWarning.GetAPIRepresentation(), customEventSpec.Rules[0].Severity)
}

func TestShouldFailToConvertCustomEventSpecificationWithEntityVerificationRuleStateToDataModelWhenSeverityIsNotValid(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.Set(CustomEventSpecificationRuleSeverity, "INVALID")

	_, err := resourceHandle.MapStateToDataObject(resourceData, utils.NewResourceNameFormatter("prefix ", " suffix"))

	assert.NotNil(t, err)
}
