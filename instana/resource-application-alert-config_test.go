package instana_test

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
)

func TestShouldReturnTrueWhenCheckingForSchemaDiffSuppressForTagFilterOfApplicationAlertConfigAndValueCanBeNormalizedAndOldAndNewNormalizedValueAreEqual(t *testing.T) {
	resourceHandle := NewApplicationAlertConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	oldValue := expressionEntityTypeDestEqValue
	newValue := "entity.type  EQUALS    'foo'"

	require.True(t, schema[ApplicationAlertConfigFieldTagFilter].DiffSuppressFunc(ApplicationAlertConfigFieldTagFilter, oldValue, newValue, nil))
}

func TestShouldReturnFalseWhenCheckingForSchemaDiffSuppressForTagFilterOfApplicationAlertConfigAndValueCanBeNormalizedAndOldAndNewNormalizedValueAreNotEqual(t *testing.T) {
	resourceHandle := NewApplicationAlertConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	oldValue := expressionEntityTypeSrcEqValue
	newValue := validTagFilter

	require.False(t, schema[ApplicationAlertConfigFieldTagFilter].DiffSuppressFunc(ApplicationAlertConfigFieldTagFilter, oldValue, newValue, nil))
}

func TestShouldReturnTrueWhenCheckingForSchemaDiffSuppressForTagFilterOfApplicationAlertConfigAndValueCannotBeNormalizedAndOldAndNewValueAreEqual(t *testing.T) {
	resourceHandle := NewApplicationAlertConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	invalidValue := invalidTagFilter

	require.True(t, schema[ApplicationAlertConfigFieldTagFilter].DiffSuppressFunc(ApplicationAlertConfigFieldTagFilter, invalidValue, invalidValue, nil))
}

func TestShouldReturnFalseWhenCheckingForSchemaDiffSuppressForTagFilterOfApplicationAlertConfigAndValueCannotBeNormalizedAndOldAndNewValueAreNotEqual(t *testing.T) {
	resourceHandle := NewApplicationAlertConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	oldValue := invalidTagFilter
	newValue := "entity.type foo foo foo"

	require.False(t, schema[ApplicationAlertConfigFieldTagFilter].DiffSuppressFunc(ApplicationAlertConfigFieldTagFilter, oldValue, newValue, nil))
}

func TestShouldReturnNormalizedValueForTagFilterOfApplicationAlertConfigWhenStateFuncIsCalledAndValueCanBeNormalized(t *testing.T) {
	resourceHandle := NewApplicationAlertConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	expectedValue := expressionEntityTypeDestEqValue
	newValue := validTagFilter

	require.Equal(t, expectedValue, schema[ApplicationAlertConfigFieldTagFilter].StateFunc(newValue))
}

func TestShouldReturnProvidedValueForTagFilterOfApplicationAlertConfigWhenStateFuncIsCalledAndValueCannotBeNormalized(t *testing.T) {
	resourceHandle := NewApplicationAlertConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	value := invalidTagFilter

	require.Equal(t, value, schema[ApplicationAlertConfigFieldTagFilter].StateFunc(value))
}

func TestShouldReturnNoErrorsAndWarningsWhenValidationOfTagFilterOfApplicationAlertConfigIsCalledAndValueCanBeParsed(t *testing.T) {
	resourceHandle := NewApplicationAlertConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	value := validTagFilter

	warns, errs := schema[ApplicationAlertConfigFieldTagFilter].ValidateFunc(value, ApplicationAlertConfigFieldTagFilter)
	require.Empty(t, warns)
	require.Empty(t, errs)
}

func TestShouldReturnOneErrorAndNoWarningsWhenValidationOfTagFilterOfApplicationAlertConfigIsCalledAndValueCannotBeParsed(t *testing.T) {
	resourceHandle := NewApplicationAlertConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	value := invalidTagFilter

	warns, errs := schema[ApplicationAlertConfigFieldTagFilter].ValidateFunc(value, ApplicationAlertConfigFieldTagFilter)
	require.Empty(t, warns)
	require.Len(t, errs, 1)
}

func TestApplicationAlertConfigResourceShouldHaveSchemaVersionZero(t *testing.T) {
	require.Equal(t, 0, NewApplicationAlertConfigResourceHandle().MetaData().SchemaVersion)
}

func TestApplicationConfigResourceShouldHaveNoStateUpgrader(t *testing.T) {
	resourceHandler := NewApplicationAlertConfigResourceHandle()

	require.Empty(t, resourceHandler.StateUpgraders())
}

func TestShouldReturnCorrectResourceNameForApplicationAlertConfigResource(t *testing.T) {
	name := NewApplicationAlertConfigResourceHandle().MetaData().ResourceName

	require.Equal(t, name, "instana_application_alert_config")
}

func TestShouldUpdateApplicationConfigTerraformResourceStateFromModelWhenThroughputRuleIsProvided(t *testing.T) {
	fullName := "prefix application-alert-config-name suffix"
	applicationAlertConfigID := "application-alert-config-id"
	stableHash := int32(1234)
	thresholdValue := 123.3
	thresholdLastUpdate := int64(12345)
	applicationConfig := restapi.ApplicationAlertConfig{
		ID:              applicationAlertConfigID,
		AlertChannelIDs: []string{"channel-1", "channel-2"},
		Applications: map[string]restapi.IncludedApplication{
			"app-1": {
				ApplicationID: "app-1",
				Inclusive:     true,
				Services: map[string]restapi.IncludedService{
					"srv-1": {
						ServiceID: "srv-1",
						Inclusive: true,
						Endpoints: map[string]restapi.IncludedEndpoint{
							"edp-1": {
								EndpointID: "edp-1",
								Inclusive:  true,
							},
						},
					},
				},
			},
		},
		BoundaryScope:    restapi.BoundaryScopeInbound,
		Description:      "application-alert-config-description",
		EvaluationType:   restapi.EvaluationTypePerApplication,
		Granularity:      restapi.Granularity600000,
		IncludeInternal:  false,
		IncludeSynthetic: false,
		Name:             fullName,
		Rule: restapi.ApplicationAlertRule{
			AlertType:   "throughput",
			Aggregation: restapi.MinAggregation,
			MetricName:  "test-metric",
			StableHash:  &stableHash,
		},
		Severity:            restapi.SeverityCritical.GetAPIRepresentation(),
		TagFilterExpression: restapi.NewStringTagFilter(restapi.TagFilterEntitySource, "service.name", restapi.EqualsOperator, "test"),
		Threshold: restapi.Threshold{
			Type:        "staticThreshold",
			Operator:    restapi.ThresholdOperatorGreaterThan,
			Value:       &thresholdValue,
			LastUpdated: &thresholdLastUpdate,
		},
		TimeThreshold: restapi.TimeThreshold{
			Type:       "violationsInSequence",
			TimeWindow: 12345,
		},
		Triggering: true,
	}

	testHelper := NewTestHelper(t)
	sut := NewApplicationAlertConfigResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, &applicationConfig, testHelper.ResourceFormatter())

	require.NoError(t, err)
	require.Equal(t, applicationAlertConfigID, resourceData.Id())
	require.Equal(t, []interface{}{"channel-2", "channel-1"}, (resourceData.Get(ApplicationAlertConfigFieldAlertChannelIDs).(*schema.Set)).List())
	require.Equal(t, []interface{}{map[string]interface{}{"application_id": "app-1", "inclusive": true, "services": []interface{}{map[string]interface{}{"endpoints": []interface{}{map[string]interface{}{"endpoint_id": "edp-1", "inclusive": true}}, "inclusive": true, "service_id": "srv-1"}}}}, resourceData.Get(ApplicationAlertConfigFieldApplications))
	require.Equal(t, string(restapi.BoundaryScopeInbound), resourceData.Get(ApplicationAlertConfigFieldBoundaryScope))
	require.Equal(t, "application-alert-config-description", resourceData.Get(ApplicationAlertConfigFieldDescription))
	require.Equal(t, "application-alert-config-name", resourceData.Get(ApplicationAlertConfigFieldName))
	require.Equal(t, fullName, resourceData.Get(ApplicationAlertConfigFieldFullName))
	require.Equal(t, string(restapi.EvaluationTypePerApplication), resourceData.Get(ApplicationAlertConfigFieldEvaluationType))
	require.False(t, resourceData.Get(ApplicationAlertConfigFieldIncludeInternal).(bool))
	require.False(t, resourceData.Get(ApplicationAlertConfigFieldIncludeSynthetic).(bool))
	require.Empty(t, resourceData.Get(ApplicationAlertConfigFieldCustomPayloadFields))
	require.Equal(t, []interface{}{map[string]interface{}{"error_rate": []interface{}{}, "logs": []interface{}{}, "slowness": []interface{}{}, "status_code": []interface{}{}, "throughput": []interface{}{map[string]interface{}{"aggregation": "MIN", "metric_name": "test-metric", "stable_hash": 1234}}}}, resourceData.Get(ApplicationAlertConfigFieldRule))
	require.Equal(t, restapi.SeverityCritical.GetTerraformRepresentation(), resourceData.Get(ApplicationAlertConfigFieldSeverity))
	require.Equal(t, "service.name@src EQUALS 'test'", resourceData.Get(ApplicationAlertConfigFieldTagFilter))
	require.Equal(t, []interface{}{map[string]interface{}{"historic_baseline": []interface{}{}, "static": []interface{}{map[string]interface{}{"last_updated": 12345, "operator": ">", "value": 123.3}}}}, resourceData.Get(ApplicationAlertConfigFieldThreshold))
	require.Equal(t, []interface{}{map[string]interface{}{"request_impact": []interface{}{}, "violations_in_period": []interface{}{}, "violations_in_sequence": []interface{}{map[string]interface{}{"time_window": 12345}}}}, resourceData.Get(ApplicationAlertConfigFieldTimeThreshold))
	require.True(t, resourceData.Get(ApplicationAlertConfigFieldTriggering).(bool))
}
