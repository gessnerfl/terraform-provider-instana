package restapi_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

func TestShouldSuccessfullyUnmarshalAlertingConfigWithRuleIds(t *testing.T) {
	response := `{
		"id" : "id",
		"alertName" : "name",
		"customPayload" : "custom",
		"integrationIds" : [ "integrationId-1", "integrationId-2" ],
		"eventFilteringConfiguration" : {
			"query" : "query",
			"ruleIds" : [ "rule-1", "rule-2" ]
		}
	}`

	result, err := NewAlertingConfigurationUnmarshaller().Unmarshal([]byte(response))

	assert.Nil(t, err)
	assert.IsType(t, AlertingConfiguration{}, result)

	config := result.(AlertingConfiguration)
	assert.Equal(t, "id", config.ID)
	assert.Equal(t, "name", config.AlertName)
	assert.Equal(t, "custom", *config.CustomPayload)
	assert.Equal(t, []string{"integrationId-1", "integrationId-2"}, config.IntegrationIDs)
	assert.Equal(t, "query", *config.EventFilteringConfiguration.Query)
	assert.Equal(t, []string{"rule-1", "rule-2"}, config.EventFilteringConfiguration.RuleIDs)
}

func TestShouldSuccessfullyUnmarshalAlertingConfigWithEventTypes(t *testing.T) {
	response := `{
		"id" : "id",
		"alertName" : "name",
		"customPayload" : "custom",
		"integrationIds" : [ "integrationId-1", "integrationId-2" ],
		"eventFilteringConfiguration" : {
			"query" : "query",
			"eventTypes" : [ "INCIDENT", "CRITICAL" ]
		}
	}`

	result, err := NewAlertingConfigurationUnmarshaller().Unmarshal([]byte(response))

	assert.Nil(t, err)
	assert.IsType(t, AlertingConfiguration{}, result)

	config := result.(AlertingConfiguration)
	assert.Equal(t, "id", config.ID)
	assert.Equal(t, "name", config.AlertName)
	assert.Equal(t, "custom", *config.CustomPayload)
	assert.Equal(t, []string{"integrationId-1", "integrationId-2"}, config.IntegrationIDs)
	assert.Equal(t, "query", *config.EventFilteringConfiguration.Query)
	assert.Equal(t, []AlertEventType{IncidentAlertEventType, CriticalAlertEventType}, config.EventFilteringConfiguration.EventTypes)
}

func TestShouldFailToUnmarshalAlertingConfigurationWhenResponseIsAJsonArray(t *testing.T) {
	response := `["test1","test2"]`

	_, err := NewAlertingConfigurationUnmarshaller().Unmarshal([]byte(response))

	assert.NotNil(t, err)
}

func TestShouldFailToUnmarshalAlertingConfigurationWhenResponseIsNotAJsonMessage(t *testing.T) {
	response := `foo bar`

	_, err := NewAlertingConfigurationUnmarshaller().Unmarshal([]byte(response))

	assert.NotNil(t, err)
}

func TestShouldReturnEmptyAlertingConfigurationWhenJsonObjectIsReturnWhereNoFiledMatches(t *testing.T) {
	response := `{"foo" : "bar" }`

	result, err := NewAlertingConfigurationUnmarshaller().Unmarshal([]byte(response))

	assert.Nil(t, err)
	assert.Equal(t, AlertingConfiguration{}, result)
}
