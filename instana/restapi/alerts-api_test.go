package restapi_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/stretchr/testify/assert"
)

const (
	alertingConfigID             = "alerting-id"
	alertingConfigName           = "alerting-name"
	alertingConfigIntegrationId1 = "alerting-integration-id1"
	alertingConfigIntegrationId2 = "alerting-integration-id2"
	alertingConfigRuleId1        = "alerting-rule-id1"
	alertingConfigRuleId2        = "alerting-rule-id2"
	alertingConfigQuery          = "alerting-query"
)

func TestReturnIdOfAlertingConfig(t *testing.T) {
	config := AlertingConfiguration{
		ID:             alertingConfigID,
		AlertName:      alertingConfigName,
		IntegrationIDs: []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2},
		EventFilteringConfiguration: EventFilteringConfiguration{
			RuleIDs: []string{alertingConfigRuleId1, alertingConfigRuleId2},
		},
	}

	assert.Equal(t, alertingConfigID, config.GetID())
}

func TestShouldSuccessFullyValidateAlertingConfigurationWhenRuleIdsAreConfigured(t *testing.T) {
	config := AlertingConfiguration{
		ID:             alertingConfigID,
		AlertName:      alertingConfigName,
		IntegrationIDs: []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2},
		EventFilteringConfiguration: EventFilteringConfiguration{
			RuleIDs: []string{alertingConfigRuleId1, alertingConfigRuleId2},
		},
	}

	assert.Nil(t, config.Validate())
}

func TestShouldSuccessFullyValidateAlertingConfigurationWhenEventTypessAreConfigured(t *testing.T) {
	config := AlertingConfiguration{
		ID:             alertingConfigID,
		AlertName:      alertingConfigName,
		IntegrationIDs: []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2},
		EventFilteringConfiguration: EventFilteringConfiguration{
			EventTypes: []AlertEventType{WarningAlertEventType, IncidentAlertEventType},
		},
	}

	assert.Nil(t, config.Validate())
}

func TestShouldSuccessFullyValidateAlertingConfigurationWhenAnAdditionalQueryIsConfiguredForTheEventFilterConfig(t *testing.T) {
	query := alertingConfigQuery
	config := AlertingConfiguration{
		ID:             alertingConfigID,
		AlertName:      alertingConfigName,
		IntegrationIDs: []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2},
		EventFilteringConfiguration: EventFilteringConfiguration{
			Query:      &query,
			EventTypes: []AlertEventType{WarningAlertEventType, IncidentAlertEventType},
		},
	}

	assert.Nil(t, config.Validate())
}

func TestShouldSuccessFullyValidateAlertingConfigurationWhenCustomPayloadIsDefined(t *testing.T) {
	customPayload := "foo"
	config := AlertingConfiguration{
		ID:             alertingConfigID,
		AlertName:      alertingConfigName,
		CustomPayload:  &customPayload,
		IntegrationIDs: []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2},
		EventFilteringConfiguration: EventFilteringConfiguration{
			EventTypes: []AlertEventType{WarningAlertEventType, IncidentAlertEventType},
		},
	}

	assert.Nil(t, config.Validate())
}

func TestShouldFailToValidateAlertingChannelConfigurationWhenIDIsMissing(t *testing.T) {
	config := AlertingConfiguration{
		AlertName:      alertingConfigName,
		IntegrationIDs: []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2},
		EventFilteringConfiguration: EventFilteringConfiguration{
			EventTypes: []AlertEventType{WarningAlertEventType, IncidentAlertEventType},
		},
	}

	err := config.Validate()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "ID")
	assert.Contains(t, err.Error(), "missing")
}

func TestShouldFailToValidateAlertingChannelConfigurationWhenIDIsBlank(t *testing.T) {
	config := AlertingConfiguration{
		ID:             "",
		AlertName:      alertingConfigName,
		IntegrationIDs: []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2},
		EventFilteringConfiguration: EventFilteringConfiguration{
			EventTypes: []AlertEventType{WarningAlertEventType, IncidentAlertEventType},
		},
	}

	err := config.Validate()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "ID")
	assert.Contains(t, err.Error(), "missing")
}

func TestShouldFailToValidateAlertingChannelConfigurationWhenAlertNameIsMissing(t *testing.T) {
	config := AlertingConfiguration{
		ID:             alertingConfigID,
		IntegrationIDs: []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2},
		EventFilteringConfiguration: EventFilteringConfiguration{
			EventTypes: []AlertEventType{WarningAlertEventType, IncidentAlertEventType},
		},
	}

	err := config.Validate()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "AlertName")
	assert.Contains(t, err.Error(), "missing")
}

func TestShouldFailToValidateAlertingChannelConfigurationWhenAlertNameIsBlank(t *testing.T) {
	config := AlertingConfiguration{
		ID:             alertingConfigID,
		AlertName:      "",
		IntegrationIDs: []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2},
		EventFilteringConfiguration: EventFilteringConfiguration{
			EventTypes: []AlertEventType{WarningAlertEventType, IncidentAlertEventType},
		},
	}

	err := config.Validate()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "AlertName")
	assert.Contains(t, err.Error(), "missing")
}

func TestShouldFailToValidateAlertingChannelConfigurationWhenAlertNameExceedsMaxLength(t *testing.T) {
	config := AlertingConfiguration{
		ID:             alertingConfigID,
		AlertName:      utils.RandomString(1025),
		IntegrationIDs: []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2},
		EventFilteringConfiguration: EventFilteringConfiguration{
			EventTypes: []AlertEventType{WarningAlertEventType, IncidentAlertEventType},
		},
	}

	err := config.Validate()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "AlertName")
	assert.Contains(t, err.Error(), "length")
}

func TestShouldFailToValidateAlertingChannelConfigurationWhenNoIntegrationIdIsProvided(t *testing.T) {
	config := AlertingConfiguration{
		ID:        alertingConfigID,
		AlertName: alertingConfigName,
		EventFilteringConfiguration: EventFilteringConfiguration{
			EventTypes: []AlertEventType{WarningAlertEventType, IncidentAlertEventType},
		},
	}

	err := config.Validate()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "IntegrationID")
	assert.Contains(t, err.Error(), "provided")
}

func TestShouldFailToValidateAlertingChannelConfigurationWhenTooManyIntegrationIDsAreProvided(t *testing.T) {
	integrationIDs := make([]string, 1025)
	for i, _ := range integrationIDs {
		integrationIDs[i] = utils.RandomString(10)
	}
	config := AlertingConfiguration{
		ID:             alertingConfigID,
		AlertName:      alertingConfigName,
		IntegrationIDs: integrationIDs,
		EventFilteringConfiguration: EventFilteringConfiguration{
			EventTypes: []AlertEventType{WarningAlertEventType, IncidentAlertEventType},
		},
	}

	err := config.Validate()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "IntegrationID")
	assert.Contains(t, err.Error(), "number")
}

func TestShouldFailToValidateAlertingChannelConfigurationWhenTooManyIntegrationIDsAreNotUnique(t *testing.T) {
	config := AlertingConfiguration{
		ID:             alertingConfigID,
		AlertName:      alertingConfigName,
		IntegrationIDs: []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2, alertingConfigIntegrationId1},
		EventFilteringConfiguration: EventFilteringConfiguration{
			EventTypes: []AlertEventType{WarningAlertEventType, IncidentAlertEventType},
		},
	}

	err := config.Validate()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "IntegrationID")
	assert.Contains(t, err.Error(), "unique")
}

func TestShouldFailToValidateAlertingChannelConfigurationWhenCustomPayloadExceedsMaximumLength(t *testing.T) {
	customPayload := utils.RandomString(65537)
	config := AlertingConfiguration{
		ID:             alertingConfigID,
		AlertName:      alertingConfigName,
		IntegrationIDs: []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2},
		EventFilteringConfiguration: EventFilteringConfiguration{
			EventTypes: []AlertEventType{WarningAlertEventType, IncidentAlertEventType},
		},
		CustomPayload: &customPayload,
	}

	err := config.Validate()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "CustomPayload")
	assert.Contains(t, err.Error(), "length")
}

func TestShouldFailToValidateAlertingConfigurationWhenRuleIdsAndEventTypesAreConfigured(t *testing.T) {
	config := AlertingConfiguration{
		ID:             alertingConfigID,
		AlertName:      alertingConfigName,
		IntegrationIDs: []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2},
		EventFilteringConfiguration: EventFilteringConfiguration{
			RuleIDs:    []string{alertingConfigRuleId1, alertingConfigRuleId2},
			EventTypes: []AlertEventType{WarningAlertEventType, IncidentAlertEventType},
		},
	}

	err := config.Validate()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Either")
	assert.Contains(t, err.Error(), "RuleIDs")
	assert.Contains(t, err.Error(), "EventTypes")
}

func TestShouldFailToValidateAlertingConfigurationWhenNeitherRuleIdsNorEventTypesAreConfigured(t *testing.T) {
	config := AlertingConfiguration{
		ID:                          alertingConfigID,
		AlertName:                   alertingConfigName,
		IntegrationIDs:              []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2},
		EventFilteringConfiguration: EventFilteringConfiguration{},
	}

	err := config.Validate()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Either")
	assert.Contains(t, err.Error(), "RuleIDs")
	assert.Contains(t, err.Error(), "EventTypes")
}

func TestShouldFailToValidateAlertingConfigurationWhenRuleIdsExceedTheMaximumNumberOfAllowedRuleIds(t *testing.T) {
	ruleIDs := make([]string, 1025)
	for i, _ := range ruleIDs {
		ruleIDs[i] = utils.RandomString(10)
	}
	config := AlertingConfiguration{
		ID:             alertingConfigID,
		AlertName:      alertingConfigName,
		IntegrationIDs: []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2},
		EventFilteringConfiguration: EventFilteringConfiguration{
			RuleIDs: ruleIDs,
		},
	}

	err := config.Validate()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "RuleIDs")
	assert.Contains(t, err.Error(), "number")
}

func TestShouldFailToValidateAlertingConfigurationWhenRuleIdsAreNotUnique(t *testing.T) {
	config := AlertingConfiguration{
		ID:             alertingConfigID,
		AlertName:      alertingConfigName,
		IntegrationIDs: []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2},
		EventFilteringConfiguration: EventFilteringConfiguration{
			RuleIDs: []string{alertingConfigRuleId1, alertingConfigRuleId2, alertingConfigRuleId1},
		},
	}

	err := config.Validate()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "RuleIDs")
	assert.Contains(t, err.Error(), "unique")
}

func TestShouldFailToValidateAlertingConfigurationWhenEventTypesExceedTheNumberOfAllowedEventTypes(t *testing.T) {
	config := AlertingConfiguration{
		ID:             alertingConfigID,
		AlertName:      alertingConfigName,
		IntegrationIDs: []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2},
		EventFilteringConfiguration: EventFilteringConfiguration{
			EventTypes: SupportedAlertEventTypes,
		},
	}

	err := config.Validate()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "EventTypes")
	assert.Contains(t, err.Error(), "number")
}

func TestShouldFailToValidateAlertingConfigurationWhenEventTypesAreNotUnique(t *testing.T) {
	config := AlertingConfiguration{
		ID:             alertingConfigID,
		AlertName:      alertingConfigName,
		IntegrationIDs: []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2},
		EventFilteringConfiguration: EventFilteringConfiguration{
			EventTypes: []AlertEventType{WarningAlertEventType, CriticalAlertEventType, WarningAlertEventType},
		},
	}

	err := config.Validate()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "EventTypes")
	assert.Contains(t, err.Error(), "unique")
}

func TestShouldFailToValidateAlertingConfigurationWhenEventTypeIsNotSupported(t *testing.T) {
	config := AlertingConfiguration{
		ID:             alertingConfigID,
		AlertName:      alertingConfigName,
		IntegrationIDs: []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2},
		EventFilteringConfiguration: EventFilteringConfiguration{
			EventTypes: []AlertEventType{AlertEventType("INVALID")},
		},
	}

	err := config.Validate()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "INVALID")
	assert.Contains(t, err.Error(), "EventType")
	assert.Contains(t, err.Error(), "supported")
}

func TestShouldFailToValidateAlertingConfigurationWhenQueryExceedsTheMaximumNumberOfCharacters(t *testing.T) {
	query := utils.RandomString(2049)
	config := AlertingConfiguration{
		ID:             alertingConfigID,
		AlertName:      alertingConfigName,
		IntegrationIDs: []string{alertingConfigIntegrationId1, alertingConfigIntegrationId2},
		EventFilteringConfiguration: EventFilteringConfiguration{
			Query:   &query,
			RuleIDs: []string{alertingConfigRuleId1, alertingConfigRuleId2},
		},
	}

	err := config.Validate()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Query")
	assert.Contains(t, err.Error(), "length")
}
