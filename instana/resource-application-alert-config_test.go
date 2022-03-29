package instana_test

import (
	"fmt"
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

func TestShouldUpdateApplicationConfigTerraformResourceStateFromModel(t *testing.T) {
	metricName := "test-metric"
	stableHash := int32(1234)
	statusCodeStart := int32(200)
	statusCodeEnd := int32(300)
	logMessage := "test-log-message"
	logLevel := restapi.LogLevelError
	logOperator := restapi.EqualsOperator
	rules := []testPair[restapi.ApplicationAlertRule, []interface{}]{
		{
			name: "Throughput",
			input: restapi.ApplicationAlertRule{
				AlertType:   "throughput",
				Aggregation: restapi.MinAggregation,
				MetricName:  metricName,
				StableHash:  &stableHash,
			},
			expected: []interface{}{map[string]interface{}{"error_rate": []interface{}{}, "logs": []interface{}{}, "slowness": []interface{}{}, "status_code": []interface{}{}, "throughput": []interface{}{map[string]interface{}{"aggregation": string(restapi.MinAggregation), "metric_name": metricName, "stable_hash": int(stableHash)}}}},
		},
		{
			name: "StatusCode",
			input: restapi.ApplicationAlertRule{
				AlertType:       "statusCode",
				Aggregation:     restapi.MinAggregation,
				MetricName:      metricName,
				StableHash:      &stableHash,
				StatusCodeStart: &statusCodeStart,
				StatusCodeEnd:   &statusCodeEnd,
			},
			expected: []interface{}{map[string]interface{}{"error_rate": []interface{}{}, "logs": []interface{}{}, "slowness": []interface{}{}, "status_code": []interface{}{map[string]interface{}{"aggregation": string(restapi.MinAggregation), "metric_name": metricName, "stable_hash": int(stableHash), "status_code_start": int(statusCodeStart), "status_code_end": int(statusCodeEnd)}}, "throughput": []interface{}{}}},
		},
		{
			name: "Slowness",
			input: restapi.ApplicationAlertRule{
				AlertType:   "slowness",
				Aggregation: restapi.MinAggregation,
				MetricName:  metricName,
				StableHash:  &stableHash,
			},
			expected: []interface{}{map[string]interface{}{"error_rate": []interface{}{}, "logs": []interface{}{}, "slowness": []interface{}{map[string]interface{}{"aggregation": string(restapi.MinAggregation), "metric_name": metricName, "stable_hash": int(stableHash)}}, "status_code": []interface{}{}, "throughput": []interface{}{}}},
		},
		{
			name: "Logs",
			input: restapi.ApplicationAlertRule{
				AlertType:   "logs",
				Aggregation: restapi.MinAggregation,
				MetricName:  metricName,
				StableHash:  &stableHash,
				Message:     &logMessage,
				Operator:    &logOperator,
				Level:       &logLevel,
			},
			expected: []interface{}{map[string]interface{}{"error_rate": []interface{}{}, "logs": []interface{}{map[string]interface{}{"aggregation": string(restapi.MinAggregation), "metric_name": metricName, "stable_hash": int(stableHash), "level": string(logLevel), "message": logMessage, "operator": string(logOperator)}}, "slowness": []interface{}{}, "status_code": []interface{}{}, "throughput": []interface{}{}}},
		},
		{
			name: "ErrorRate",
			input: restapi.ApplicationAlertRule{
				AlertType:   "errorRate",
				Aggregation: restapi.MinAggregation,
				MetricName:  metricName,
				StableHash:  &stableHash,
			},
			expected: []interface{}{map[string]interface{}{"error_rate": []interface{}{map[string]interface{}{"aggregation": string(restapi.MinAggregation), "metric_name": metricName, "stable_hash": int(stableHash)}}, "logs": []interface{}{}, "slowness": []interface{}{}, "status_code": []interface{}{}, "throughput": []interface{}{}}},
		},
	}

	thresholdValue := 123.3
	thresholdLastUpdate := int64(12345)
	thresholdSeasonality := restapi.ThresholdSeasonalityDaily
	thresholdBaseline := [][]float64{{1.23, 4.56}, {1.23, 7.89}}
	thresholdDeviationFactor := float32(1.23)
	thresholds := []testPair[restapi.Threshold, []interface{}]{
		{
			name: "StaticThreshold",
			input: restapi.Threshold{
				Type:        "staticThreshold",
				Operator:    restapi.ThresholdOperatorGreaterThan,
				LastUpdated: &thresholdLastUpdate,
				Value:       &thresholdValue,
			},
			expected: []interface{}{map[string]interface{}{"historic_baseline": []interface{}{}, "static": []interface{}{map[string]interface{}{"last_updated": int(thresholdLastUpdate), "operator": string(restapi.ThresholdOperatorGreaterThan), "value": thresholdValue}}}},
		},
		{
			name: "HistoricBaseLine",
			input: restapi.Threshold{
				Type:            "historicBaseline",
				Operator:        restapi.ThresholdOperatorGreaterThan,
				LastUpdated:     &thresholdLastUpdate,
				Seasonality:     &thresholdSeasonality,
				Baseline:        &thresholdBaseline,
				DeviationFactor: &thresholdDeviationFactor,
			},
			expected: []interface{}{map[string]interface{}{"historic_baseline": []interface{}{map[string]interface{}{"last_updated": int(thresholdLastUpdate), "operator": string(restapi.ThresholdOperatorGreaterThan), "seasonality": string(thresholdSeasonality), "baseline": thresholdBaseline, "deviation_factor": float64(thresholdDeviationFactor)}}, "static": []interface{}{}}},
		},
	}

	timeThresholdWindow := int64(12345)
	timeThresholdRequests := int32(5)
	timeThresholdViolations := int32(3)
	timeThresholds := []testPair[restapi.TimeThreshold, []interface{}]{
		{
			name: "RequestImpact",
			input: restapi.TimeThreshold{
				Type:       "requestImpact",
				TimeWindow: timeThresholdWindow,
				Requests:   &timeThresholdRequests,
			},
			expected: []interface{}{map[string]interface{}{"request_impact": []interface{}{map[string]interface{}{"time_window": int(timeThresholdWindow), "requests": int(timeThresholdRequests)}}, "violations_in_period": []interface{}{}, "violations_in_sequence": []interface{}{}}},
		},
		{
			name: "ViolationsInPeriod",
			input: restapi.TimeThreshold{
				Type:       "violationsInPeriod",
				TimeWindow: timeThresholdWindow,
				Violations: &timeThresholdViolations,
			},
			expected: []interface{}{map[string]interface{}{"request_impact": []interface{}{}, "violations_in_period": []interface{}{map[string]interface{}{"time_window": int(timeThresholdWindow), "violations": int(timeThresholdViolations)}}, "violations_in_sequence": []interface{}{}}},
		},
		{
			name: "ViolationsInSequence",
			input: restapi.TimeThreshold{
				Type:       "violationsInSequence",
				TimeWindow: timeThresholdWindow,
			},
			expected: []interface{}{map[string]interface{}{"request_impact": []interface{}{}, "violations_in_period": []interface{}{}, "violations_in_sequence": []interface{}{map[string]interface{}{"time_window": int(timeThresholdWindow)}}}},
		},
	}

	for _, rule := range rules {
		for _, threshold := range thresholds {
			for _, timeThreshold := range timeThresholds {
				t.Run(fmt.Sprintf("TestShouldUpdateApplicationConfigTerraformResourceStateFromModelWith%sAnd%sAnd%s", rule.name, threshold.name, timeThreshold.name), createTestShouldUpdateApplicationConfigTerraformResourceStateFromModelCase(rule, threshold, timeThreshold))
			}
		}
	}
}

func createTestShouldUpdateApplicationConfigTerraformResourceStateFromModelCase(ruleTestPair testPair[restapi.ApplicationAlertRule, []interface{}],
	thresholdTestPair testPair[restapi.Threshold, []interface{}],
	timeThresholdTestPair testPair[restapi.TimeThreshold, []interface{}]) func(t *testing.T) {
	return func(t *testing.T) {
		fullName := "prefix application-alert-config-name suffix"
		applicationAlertConfigID := "application-alert-config-id"
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
			BoundaryScope:       restapi.BoundaryScopeInbound,
			Description:         "application-alert-config-description",
			EvaluationType:      restapi.EvaluationTypePerApplication,
			Granularity:         restapi.Granularity600000,
			IncludeInternal:     false,
			IncludeSynthetic:    false,
			Name:                fullName,
			Rule:                ruleTestPair.input,
			Severity:            restapi.SeverityCritical.GetAPIRepresentation(),
			TagFilterExpression: restapi.NewStringTagFilter(restapi.TagFilterEntitySource, "service.name", restapi.EqualsOperator, "test"),
			Threshold:           thresholdTestPair.input,
			TimeThreshold:       timeThresholdTestPair.input,
			Triggering:          true,
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
		require.Equal(t, ruleTestPair.expected, resourceData.Get(ApplicationAlertConfigFieldRule))
		require.Equal(t, restapi.SeverityCritical.GetTerraformRepresentation(), resourceData.Get(ApplicationAlertConfigFieldSeverity))
		require.Equal(t, "service.name@src EQUALS 'test'", resourceData.Get(ApplicationAlertConfigFieldTagFilter))
		requireApplicationAlertConfigThresholdSetOnSchema(t, thresholdTestPair.expected, resourceData)
		require.Equal(t, timeThresholdTestPair.expected, resourceData.Get(ApplicationAlertConfigFieldTimeThreshold))
		require.True(t, resourceData.Get(ApplicationAlertConfigFieldTriggering).(bool))
	}
}

func requireApplicationAlertConfigThresholdSetOnSchema(t *testing.T, expected []interface{}, resourceData *schema.ResourceData) {
	actual := resourceData.Get(ApplicationAlertConfigFieldThreshold).([]interface{})
	require.Equal(t, 1, len(expected))
	require.Equal(t, len(expected), len(actual))
	expectedEntries := expected[0].(map[string]interface{})
	actualEntries := actual[0].(map[string]interface{})

	for k, e := range expectedEntries {
		if k == "historic_baseline" && len(e.([]interface{})) > 0 {
			expectedHistoricBaselineSlice := e.([]interface{})
			actualHistoricBaselineSlice := actualEntries[k].([]interface{})
			require.Equal(t, 1, len(expectedHistoricBaselineSlice))
			require.Equal(t, len(expected), len(actual))
			expectedHistoricBaseline := expectedHistoricBaselineSlice[0].(map[string]interface{})
			actualHistoricBaseline := actualHistoricBaselineSlice[0].(map[string]interface{})
			for key, expectedBaselineValue := range expectedHistoricBaseline {
				if key == "baseline" {
					actualBaseline := actualHistoricBaseline[key].(*schema.Set)
					actualBaselineSlice := make([][]float64, actualBaseline.Len())
					for i, v := range actualBaseline.List() {
						values := v.(*schema.Set).List()
						actualBaselineSlice[i] = toFloat64Slice(values)
					}
					require.Equal(t, expectedBaselineValue, actualBaselineSlice)
				} else {
					require.Equal(t, expectedBaselineValue, actualHistoricBaseline[key])
				}
			}
		} else {
			require.Equal(t, e, actualEntries[k])
		}
	}
}

func toFloat64Slice(input []interface{}) []float64 {
	result := make([]float64, len(input))
	for i, v := range input {
		result[i] = v.(float64)
	}
	return result
}
