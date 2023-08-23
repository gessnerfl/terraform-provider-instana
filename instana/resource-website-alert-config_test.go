package instana_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
)

func TestWebsiteAlertConfig(t *testing.T) {
	terraformResourceInstanceName := ResourceInstanaWebsiteAlertConfig + ".example"
	inst := &websiteAlertConfigTest{
		terraformResourceInstanceName: terraformResourceInstanceName,
		resourceHandle:                NewWebsiteAlertConfigResourceHandle(),
	}
	inst.run(t)
}

type websiteAlertConfigTest struct {
	terraformResourceInstanceName string
	resourceHandle                ResourceHandle[*restapi.WebsiteAlertConfig]
}

var websiteAlertConfigTerraformTemplate = `
resource "instana_website_alert_config" "example" {
	name              = "name %d"
    description       = "test-alert-description"
    severity          = "warning"
    triggering        = false
    alert_channel_ids = [ "alert-channel-id-1", "alert-channel-id-2" ]
    granularity       = 600000
	tag_filter        = "call.type@na EQUALS 'HTTP'"
    website_id        = "website-id"

	rule {
		slowness {
			metric_name = "latency"
			aggregation = "P90"
		}
    }

	threshold {
		static {
			operator = ">="
			value    = 5.0
		}
    }

	time_threshold {
		violations_in_sequence {
			time_window = 600000
		}
    }

	custom_payload_field {
		key   = "test"
		value = "test123"
	}
}
`

var websiteAlertConfigServerResponseTemplate = `
	{
    "id": "%s",
    "name": "name %d",
    "description": "test-alert-description",
    "websiteId": "website-id",
    "severity": 5,
    "triggering": false,
    "tagFilters": [],
    "tagFilterExpression": {
      "type": "TAG_FILTER",
      "name": "call.type",
      "stringValue": "HTTP",
      "numberValue": null,
      "booleanValue": null,
      "key": null,
      "value": "HTTP",
      "operator": "EQUALS",
      "entity": "NOT_APPLICABLE"
    },
    "rule": {
      "alertType": "slowness",
      "aggregation": "P90",
      "metricName": "latency"
    },
    "threshold": {
      "type": "staticThreshold",
      "operator": ">=",
      "value": 5.0,
      "lastUpdated": 0
    },
    "alertChannelIds": [ "alert-channel-id-1", "alert-channel-id-2" ],
    "granularity": 600000,
    "timeThreshold": {
      "type": "violationsInSequence",
      "timeWindow": 600000
    },
    "customPayloadFields": [
		{
			"type": "staticString",
			"key": "test",
			"value": "test123"
      	}
	],
    "created": 1647679325301,
    "readOnly": false,
    "enabled": true,
    "derivedFromGlobalAlert": false
  }
`

func (test *websiteAlertConfigTest) run(t *testing.T) {
	t.Run(fmt.Sprintf("CRUD integration test of %s", ResourceInstanaWebsiteAlertConfig), test.createIntegrationTest())
	t.Run(fmt.Sprintf("DiffSuppressFunc of TagFilter of %s should return true when value can be normalized and old and new normalized value are equal", ResourceInstanaWebsiteAlertConfig), test.createTestOfDiffSuppressFuncOfTagFilterShouldReturnTrueWhenValueCanBeNormalizedAndOldAndNewNormalizedValueAreEqual())
	t.Run(fmt.Sprintf("DiffSuppressFunc of TagFilter of %s should return false when value can be normalized and old and new normalized value are not equal", ResourceInstanaWebsiteAlertConfig), test.createTestOfDiffSuppressFuncOfTagFilterShouldReturnFalseWhenValueCanBeNormalizedAndOldAndNewNormalizedValueAreNotEqual())
	t.Run(fmt.Sprintf("DiffSuppressFunc of TagFilter of %s should return true when value can be normalized and old and new value are equal", ResourceInstanaWebsiteAlertConfig), test.createTestOfDiffSuppressFuncOfTagFilterShouldReturnTrueWhenValueCannotBeNormalizedAndOldAndNewValueAreEqual())
	t.Run(fmt.Sprintf("DiffSuppressFunc of TagFilter of %s should return false when value cannot be normalized and old and new value are not equal", ResourceInstanaWebsiteAlertConfig), test.createTestOfDiffSuppressFuncOfTagFilterShouldReturnFalseWhenValueCannotBeNormalizedAndOldAndNewValueAreNotEqual())
	t.Run(fmt.Sprintf("StateFunc of TagFilter of %s should return normalized value when value can be normalized", ResourceInstanaWebsiteAlertConfig), test.createTestOfStateFuncOfTagFilterShouldReturnNormalizedValueWhenValueCanBeNormalized())
	t.Run(fmt.Sprintf("StateFunc of TagFilter of %s should return provided value when value cannot be normalized", ResourceInstanaWebsiteAlertConfig), test.createTestOfStateFuncOfTagFilterShouldReturnProvidedValueWhenValueCannotBeNormalized())
	t.Run(fmt.Sprintf("ValidateFunc of TagFilter of %s should return no errors and warnings when value can be parsed", ResourceInstanaWebsiteAlertConfig), test.createTestOfValidateFuncOfTagFilterShouldReturnNoErrorsAndWarningsWhenValueCanBeParsed())
	t.Run(fmt.Sprintf("ValidateFunc of TagFilter of %s should return one error and no warnings when value can be parsed", ResourceInstanaWebsiteAlertConfig), test.createTestOfValidateFuncOfTagFilterShouldReturnOneErrorAndNoWarningsWhenValueCannotBeParsed())
	t.Run(fmt.Sprintf("%s should have schema version one", ResourceInstanaWebsiteAlertConfig), test.createTestResourceShouldHaveSchemaVersionOne())
	t.Run(fmt.Sprintf("%s should have one state upgrader", ResourceInstanaWebsiteAlertConfig), test.createTestResourceShouldHaveOneStateUpgrader())
	t.Run(fmt.Sprintf("%s should migrate fullname to name when executing first state migration and fullname is available", ResourceInstanaWebsiteAlertConfig), test.createTestWebsiteAlertConfigShouldMigrateFullnameToNameWhenExecutingFirstStateUpgraderAndFullnameIsAvailable())
	t.Run(fmt.Sprintf("%s should do nothing when executing first state migration and fullname is not available", ResourceInstanaWebsiteAlertConfig), test.createTestWebsiteAlertConfigShouldDoNothingWhenExecutingFirstStateUpgraderAndFullnameIsAvailable())
	t.Run(fmt.Sprintf("%s should have correct resouce name", ResourceInstanaWebsiteAlertConfig), test.createTestResourceShouldHaveCorrectResourceName())
	test.createTestCasesForUpdatesOfTerraformResourceStateFromModel(t)
	t.Run(fmt.Sprintf("%s should fail to update state from model when severity is invalid", ResourceInstanaWebsiteAlertConfig), test.createTestCasesShouldFailToUpdateTerraformResourceStateFromModeWhenSeverityIsNotValid())
	t.Run(fmt.Sprintf("%s should fail to update state from model when tag filter expression is invalid", ResourceInstanaWebsiteAlertConfig), test.createTestCasesShouldFailToUpdateTerraformResourceStateFromModeWhenTagFilterExpressionIsNotValid())
	test.createTestCasesForMappingOfTerraformResourceStateToModel(t)
	t.Run(fmt.Sprintf("%s should fail to map state to model when severity is invalid", ResourceInstanaWebsiteAlertConfig), test.createTestCaseShouldFailToMapTerraformResourceStateToModelWhenSeverityIsNotValid())
	t.Run(fmt.Sprintf("%s should fail to map state to model when tag filter expression is invalid", ResourceInstanaWebsiteAlertConfig), test.createTestCaseShouldFailToMapTerraformResourceStateToModelWhenTagFilterIsNotValid())
}

func (test *websiteAlertConfigTest) createIntegrationTest() func(t *testing.T) {
	return func(t *testing.T) {
		id := RandomID()
		resourceRestAPIPath := restapi.WebsiteAlertConfigResourcePath
		resourceInstanceRestAPIPath := resourceRestAPIPath + "/{internal-id}"

		httpServer := testutils.NewTestHTTPServer()
		httpServer.AddRoute(http.MethodPost, resourceRestAPIPath, func(w http.ResponseWriter, r *http.Request) {
			config := &restapi.WebsiteAlertConfig{}
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
		httpServer.AddRoute(http.MethodPost, resourceInstanceRestAPIPath, func(w http.ResponseWriter, r *http.Request) {
			testutils.EchoHandlerFunc(w, r)
		})
		httpServer.AddRoute(http.MethodDelete, resourceInstanceRestAPIPath, testutils.EchoHandlerFunc)
		httpServer.AddRoute(http.MethodGet, resourceInstanceRestAPIPath, func(w http.ResponseWriter, r *http.Request) {
			modCount := httpServer.GetCallCount(http.MethodPost, resourceRestAPIPath+"/"+id)
			jsonData := fmt.Sprintf(websiteAlertConfigServerResponseTemplate, id, modCount)
			w.Header().Set(contentType, r.Header.Get(contentType))
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(jsonData))
			if err != nil {
				fmt.Printf("failed to write response; %s\n", err)
			}
		})
		httpServer.Start()
		defer httpServer.Close()

		resource.UnitTest(t, resource.TestCase{
			ProviderFactories: testProviderFactory,
			Steps: []resource.TestStep{
				test.createIntegrationTestStep(httpServer.GetPort(), 0, id),
				testStepImportWithCustomID(test.terraformResourceInstanceName, id),
				test.createIntegrationTestStep(httpServer.GetPort(), 1, id),
				testStepImportWithCustomID(test.terraformResourceInstanceName, id),
			},
		})
	}
}

func (test *websiteAlertConfigTest) createIntegrationTestStep(httpPort int64, iteration int, id string) resource.TestStep {
	ruleSlownessMetricName := fmt.Sprintf("%s.%d.%s.%d.%s", WebsiteAlertConfigFieldRule, 0, WebsiteAlertConfigFieldRuleSlowness, 0, WebsiteAlertConfigFieldRuleMetricName)
	ruleSlownessAggregation := fmt.Sprintf("%s.%d.%s.%d.%s", WebsiteAlertConfigFieldRule, 0, WebsiteAlertConfigFieldRuleSlowness, 0, WebsiteAlertConfigFieldRuleAggregation)
	thresholdStaticOperator := fmt.Sprintf("%s.%d.%s.%d.%s", ResourceFieldThreshold, 0, ResourceFieldThresholdStatic, 0, ResourceFieldThresholdOperator)
	thresholdStaticValue := fmt.Sprintf("%s.%d.%s.%d.%s", ResourceFieldThreshold, 0, ResourceFieldThresholdStatic, 0, ResourceFieldThresholdStaticValue)
	timeThresholdViolationsInSequence := fmt.Sprintf("%s.%d.%s.%d.%s", WebsiteAlertConfigFieldTimeThreshold, 0, WebsiteAlertConfigFieldTimeThresholdViolationsInSequence, 0, WebsiteAlertConfigFieldTimeThresholdTimeWindow)
	customPayloadFieldStaticKey := fmt.Sprintf("%s.%d.%s", WebsiteAlertConfigFieldCustomPayloadFields, 0, WebsiteAlertConfigFieldCustomPayloadFieldsKey)
	customPayloadFieldStaticValue := fmt.Sprintf("%s.%d.%s", WebsiteAlertConfigFieldCustomPayloadFields, 0, WebsiteAlertConfigFieldCustomPayloadFieldsValue)
	return resource.TestStep{
		Config: appendProviderConfig(fmt.Sprintf(websiteAlertConfigTerraformTemplate, iteration), httpPort),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, "id", id),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, WebsiteAlertConfigFieldName, formatResourceName(iteration)),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, WebsiteAlertConfigFieldDescription, "test-alert-description"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, WebsiteAlertConfigFieldSeverity, restapi.SeverityWarning.GetTerraformRepresentation()),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, WebsiteAlertConfigFieldTriggering, falseAsString),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, WebsiteAlertConfigFieldAlertChannelIDs+".0", "alert-channel-id-1"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, WebsiteAlertConfigFieldAlertChannelIDs+".1", "alert-channel-id-2"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, WebsiteAlertConfigFieldGranularity, "600000"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, WebsiteAlertConfigFieldTagFilter, "call.type@na EQUALS 'HTTP'"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, WebsiteAlertConfigFieldWebsiteID, "website-id"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, ruleSlownessMetricName, "latency"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, ruleSlownessAggregation, "P90"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, thresholdStaticOperator, ">="),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, thresholdStaticValue, "5"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, timeThresholdViolationsInSequence, "600000"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, customPayloadFieldStaticKey, "test"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, customPayloadFieldStaticValue, "test123"),
		),
	}
}

func (test *websiteAlertConfigTest) createTestOfDiffSuppressFuncOfTagFilterShouldReturnTrueWhenValueCanBeNormalizedAndOldAndNewNormalizedValueAreEqual() func(t *testing.T) {
	return func(t *testing.T) {
		resourceSchema := test.resourceHandle.MetaData().Schema
		oldValue := expressionEntityTypeDestEqValue
		newValue := "entity.type  EQUALS    'foo'"

		require.True(t, resourceSchema[WebsiteAlertConfigFieldTagFilter].DiffSuppressFunc(WebsiteAlertConfigFieldTagFilter, oldValue, newValue, nil))
	}
}

func (test *websiteAlertConfigTest) createTestOfDiffSuppressFuncOfTagFilterShouldReturnFalseWhenValueCanBeNormalizedAndOldAndNewNormalizedValueAreNotEqual() func(t *testing.T) {
	return func(t *testing.T) {
		resourceSchema := test.resourceHandle.MetaData().Schema
		oldValue := expressionEntityTypeSrcEqValue
		newValue := validTagFilter

		require.False(t, resourceSchema[WebsiteAlertConfigFieldTagFilter].DiffSuppressFunc(WebsiteAlertConfigFieldTagFilter, oldValue, newValue, nil))
	}
}

func (test *websiteAlertConfigTest) createTestOfDiffSuppressFuncOfTagFilterShouldReturnTrueWhenValueCannotBeNormalizedAndOldAndNewValueAreEqual() func(t *testing.T) {
	return func(t *testing.T) {
		resourceSchema := test.resourceHandle.MetaData().Schema
		invalidValue := invalidTagFilter

		require.True(t, resourceSchema[WebsiteAlertConfigFieldTagFilter].DiffSuppressFunc(WebsiteAlertConfigFieldTagFilter, invalidValue, invalidValue, nil))
	}
}

func (test *websiteAlertConfigTest) createTestOfDiffSuppressFuncOfTagFilterShouldReturnFalseWhenValueCannotBeNormalizedAndOldAndNewValueAreNotEqual() func(t *testing.T) {
	return func(t *testing.T) {
		resourceSchema := test.resourceHandle.MetaData().Schema
		oldValue := invalidTagFilter
		newValue := "entity.type foo foo foo"

		require.False(t, resourceSchema[WebsiteAlertConfigFieldTagFilter].DiffSuppressFunc(WebsiteAlertConfigFieldTagFilter, oldValue, newValue, nil))
	}
}

func (test *websiteAlertConfigTest) createTestOfStateFuncOfTagFilterShouldReturnNormalizedValueWhenValueCanBeNormalized() func(t *testing.T) {
	return func(t *testing.T) {
		resourceSchema := test.resourceHandle.MetaData().Schema
		expectedValue := expressionEntityTypeDestEqValue
		newValue := validTagFilter

		require.Equal(t, expectedValue, resourceSchema[WebsiteAlertConfigFieldTagFilter].StateFunc(newValue))
	}
}

func (test *websiteAlertConfigTest) createTestOfStateFuncOfTagFilterShouldReturnProvidedValueWhenValueCannotBeNormalized() func(t *testing.T) {
	return func(t *testing.T) {
		resourceSchema := test.resourceHandle.MetaData().Schema
		value := invalidTagFilter

		require.Equal(t, value, resourceSchema[WebsiteAlertConfigFieldTagFilter].StateFunc(value))
	}
}

func (test *websiteAlertConfigTest) createTestOfValidateFuncOfTagFilterShouldReturnNoErrorsAndWarningsWhenValueCanBeParsed() func(t *testing.T) {
	return func(t *testing.T) {
		resourceSchema := test.resourceHandle.MetaData().Schema
		value := validTagFilter

		warns, errs := resourceSchema[WebsiteAlertConfigFieldTagFilter].ValidateFunc(value, WebsiteAlertConfigFieldTagFilter)
		require.Empty(t, warns)
		require.Empty(t, errs)
	}
}

func (test *websiteAlertConfigTest) createTestOfValidateFuncOfTagFilterShouldReturnOneErrorAndNoWarningsWhenValueCannotBeParsed() func(t *testing.T) {
	return func(t *testing.T) {
		resourceSchema := test.resourceHandle.MetaData().Schema
		value := invalidTagFilter

		warns, errs := resourceSchema[WebsiteAlertConfigFieldTagFilter].ValidateFunc(value, WebsiteAlertConfigFieldTagFilter)
		require.Empty(t, warns)
		require.Len(t, errs, 1)
	}
}

func (test *websiteAlertConfigTest) createTestResourceShouldHaveSchemaVersionOne() func(t *testing.T) {
	return func(t *testing.T) {
		require.Equal(t, 1, test.resourceHandle.MetaData().SchemaVersion)
	}
}

func (test *websiteAlertConfigTest) createTestResourceShouldHaveOneStateUpgrader() func(t *testing.T) {
	return func(t *testing.T) {
		require.Len(t, test.resourceHandle.StateUpgraders(), 1)
	}
}

func (test *websiteAlertConfigTest) createTestWebsiteAlertConfigShouldMigrateFullnameToNameWhenExecutingFirstStateUpgraderAndFullnameIsAvailable() func(t *testing.T) {
	return func(t *testing.T) {
		input := map[string]interface{}{
			"full_name": "test",
		}
		result, err := NewWebsiteAlertConfigResourceHandle().StateUpgraders()[0].Upgrade(nil, input, nil)

		require.NoError(t, err)
		require.Len(t, result, 1)
		require.NotContains(t, result, WebsiteAlertConfigFieldFullName)
		require.Contains(t, result, WebsiteAlertConfigFieldName)
		require.Equal(t, "test", result[WebsiteAlertConfigFieldName])
	}
}

func (test *websiteAlertConfigTest) createTestWebsiteAlertConfigShouldDoNothingWhenExecutingFirstStateUpgraderAndFullnameIsAvailable() func(t *testing.T) {
	return func(t *testing.T) {
		input := map[string]interface{}{
			"name": "test",
		}
		result, err := NewWebsiteAlertConfigResourceHandle().StateUpgraders()[0].Upgrade(nil, input, nil)

		require.NoError(t, err)
		require.Equal(t, input, result)
	}
}

func (test *websiteAlertConfigTest) createTestResourceShouldHaveCorrectResourceName() func(t *testing.T) {
	return func(t *testing.T) {
		require.Equal(t, test.resourceHandle.MetaData().ResourceName, "instana_website_alert_config")
	}
}

func (test *websiteAlertConfigTest) createTestCasesForUpdatesOfTerraformResourceStateFromModel(t *testing.T) {
	metricName := "test-metric"
	equalsOperator := restapi.EqualsOperator
	minAggregation := restapi.MinAggregation
	statusCodeValue := "200"
	jsErrorValue := "jsErrorValue"
	rules := []testPair[restapi.WebsiteAlertRule, []interface{}]{
		{
			name: WebsiteAlertConfigFieldRuleThroughput,
			input: restapi.WebsiteAlertRule{
				AlertType:   WebsiteAlertConfigFieldRuleThroughput,
				Aggregation: &minAggregation,
				MetricName:  metricName,
			},
			expected: []interface{}{
				map[string]interface{}{
					WebsiteAlertConfigFieldRuleSpecificJsError: []interface{}{},
					WebsiteAlertConfigFieldRuleSlowness:        []interface{}{},
					WebsiteAlertConfigFieldRuleStatusCode:      []interface{}{},
					WebsiteAlertConfigFieldRuleThroughput: []interface{}{
						map[string]interface{}{
							WebsiteAlertConfigFieldRuleAggregation: string(minAggregation),
							WebsiteAlertConfigFieldRuleMetricName:  metricName,
						},
					},
				},
			},
		},
		{
			name: "StatusCode",
			input: restapi.WebsiteAlertRule{
				AlertType:   "statusCode",
				Aggregation: &minAggregation,
				MetricName:  metricName,
				Operator:    &equalsOperator,
				Value:       &statusCodeValue,
			},
			expected: []interface{}{
				map[string]interface{}{
					WebsiteAlertConfigFieldRuleSpecificJsError: []interface{}{},
					WebsiteAlertConfigFieldRuleSlowness:        []interface{}{},
					WebsiteAlertConfigFieldRuleStatusCode: []interface{}{
						map[string]interface{}{
							WebsiteAlertConfigFieldRuleAggregation: string(minAggregation),
							WebsiteAlertConfigFieldRuleMetricName:  metricName,
							WebsiteAlertConfigFieldRuleOperator:    string(equalsOperator),
							WebsiteAlertConfigFieldRuleValue:       statusCodeValue,
						},
					},
					WebsiteAlertConfigFieldRuleThroughput: []interface{}{},
				},
			},
		},
		{
			name: WebsiteAlertConfigFieldRuleSlowness,
			input: restapi.WebsiteAlertRule{
				AlertType:   WebsiteAlertConfigFieldRuleSlowness,
				Aggregation: &minAggregation,
				MetricName:  metricName,
			},
			expected: []interface{}{
				map[string]interface{}{
					WebsiteAlertConfigFieldRuleSpecificJsError: []interface{}{},
					WebsiteAlertConfigFieldRuleSlowness: []interface{}{
						map[string]interface{}{
							WebsiteAlertConfigFieldRuleAggregation: string(minAggregation),
							WebsiteAlertConfigFieldRuleMetricName:  metricName,
						},
					},
					WebsiteAlertConfigFieldRuleStatusCode: []interface{}{},
					WebsiteAlertConfigFieldRuleThroughput: []interface{}{},
				},
			},
		},
		{
			name: "SpecificJsError",
			input: restapi.WebsiteAlertRule{
				AlertType:   "specificJsError",
				Aggregation: &minAggregation,
				MetricName:  metricName,
				Operator:    &equalsOperator,
				Value:       &jsErrorValue,
			},
			expected: []interface{}{
				map[string]interface{}{
					WebsiteAlertConfigFieldRuleSpecificJsError: []interface{}{
						map[string]interface{}{
							WebsiteAlertConfigFieldRuleAggregation: string(minAggregation),
							WebsiteAlertConfigFieldRuleMetricName:  metricName,
							WebsiteAlertConfigFieldRuleOperator:    string(equalsOperator),
							WebsiteAlertConfigFieldRuleValue:       jsErrorValue,
						},
					},
					WebsiteAlertConfigFieldRuleSlowness:   []interface{}{},
					WebsiteAlertConfigFieldRuleStatusCode: []interface{}{},
					WebsiteAlertConfigFieldRuleThroughput: []interface{}{},
				},
			},
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
			expected: []interface{}{
				map[string]interface{}{
					ResourceFieldThresholdHistoricBaseline: []interface{}{},
					ResourceFieldThresholdStatic: []interface{}{
						map[string]interface{}{
							ResourceFieldThresholdLastUpdated: int(thresholdLastUpdate),
							ResourceFieldThresholdOperator:    string(restapi.ThresholdOperatorGreaterThan),
							ResourceFieldThresholdStaticValue: thresholdValue,
						},
					},
				},
			},
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
			expected: []interface{}{
				map[string]interface{}{
					ResourceFieldThresholdHistoricBaseline: []interface{}{
						map[string]interface{}{
							ResourceFieldThresholdLastUpdated:                     int(thresholdLastUpdate),
							ResourceFieldThresholdOperator:                        string(restapi.ThresholdOperatorGreaterThan),
							ResourceFieldThresholdHistoricBaselineSeasonality:     string(thresholdSeasonality),
							ResourceFieldThresholdHistoricBaselineBaseline:        thresholdBaseline,
							ResourceFieldThresholdHistoricBaselineDeviationFactor: float64(thresholdDeviationFactor),
						},
					},
					ResourceFieldThresholdStatic: []interface{}{},
				},
			},
		},
	}

	timeThresholdWindow := int64(12345)
	timeThresholdImpactMeasurementMethod := restapi.WebsiteImpactMeasurementMethodAggregated
	timeThresholdUsers := int32(5)
	timeThresholdUserPercentage := 0.8
	timeThresholdViolations := int32(3)
	timeThresholds := []testPair[restapi.WebsiteTimeThreshold, []interface{}]{
		{
			name: "UserImpactOfViolationsInSequence",
			input: restapi.WebsiteTimeThreshold{
				Type:                    "userImpactOfViolationsInSequence",
				TimeWindow:              &timeThresholdWindow,
				ImpactMeasurementMethod: &timeThresholdImpactMeasurementMethod,
				Users:                   &timeThresholdUsers,
				UserPercentage:          &timeThresholdUserPercentage,
			},
			expected: []interface{}{
				map[string]interface{}{
					WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequence: []interface{}{
						map[string]interface{}{
							WebsiteAlertConfigFieldTimeThresholdTimeWindow:                                              int(timeThresholdWindow),
							WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceImpactMeasurementMethod: string(timeThresholdImpactMeasurementMethod),
							WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceUserPercentage:          timeThresholdUserPercentage,
							WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceUsers:                   int(timeThresholdUsers),
						},
					},
					WebsiteAlertConfigFieldTimeThresholdViolationsInPeriod:   []interface{}{},
					WebsiteAlertConfigFieldTimeThresholdViolationsInSequence: []interface{}{},
				},
			},
		},
		{
			name: "ViolationsInPeriod",
			input: restapi.WebsiteTimeThreshold{
				Type:       "violationsInPeriod",
				TimeWindow: &timeThresholdWindow,
				Violations: &timeThresholdViolations,
			},
			expected: []interface{}{
				map[string]interface{}{
					WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequence: []interface{}{},
					WebsiteAlertConfigFieldTimeThresholdViolationsInPeriod: []interface{}{
						map[string]interface{}{
							WebsiteAlertConfigFieldTimeThresholdTimeWindow:                   int(timeThresholdWindow),
							WebsiteAlertConfigFieldTimeThresholdViolationsInPeriodViolations: int(timeThresholdViolations),
						},
					},
					WebsiteAlertConfigFieldTimeThresholdViolationsInSequence: []interface{}{},
				},
			},
		},
		{
			name: "ViolationsInSequence",
			input: restapi.WebsiteTimeThreshold{
				Type:       "violationsInSequence",
				TimeWindow: &timeThresholdWindow,
			},
			expected: []interface{}{
				map[string]interface{}{
					WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequence: []interface{}{},
					WebsiteAlertConfigFieldTimeThresholdViolationsInPeriod:               []interface{}{},
					WebsiteAlertConfigFieldTimeThresholdViolationsInSequence: []interface{}{
						map[string]interface{}{
							WebsiteAlertConfigFieldTimeThresholdTimeWindow: int(timeThresholdWindow),
						},
					},
				},
			},
		},
	}

	for _, rule := range rules {
		for _, threshold := range thresholds {
			for _, timeThreshold := range timeThresholds {
				t.Run(fmt.Sprintf("Should update terraform state of %s from REST response with %s and %s and %s", ResourceInstanaWebsiteAlertConfig, rule.name, threshold.name, timeThreshold.name), test.createTestShouldUpdateTerraformResourceStateFromModelCase(rule, threshold, timeThreshold))
			}
		}
	}
}

func (test *websiteAlertConfigTest) createTestShouldUpdateTerraformResourceStateFromModelCase(ruleTestPair testPair[restapi.WebsiteAlertRule, []interface{}],
	thresholdTestPair testPair[restapi.Threshold, []interface{}],
	timeThresholdTestPair testPair[restapi.WebsiteTimeThreshold, []interface{}]) func(t *testing.T) {
	return func(t *testing.T) {
		fullName := "website-alert-config-name"
		websiteAlertConfigID := "website-alert-config-id"
		websiteID := "website-id"
		websiteConfig := restapi.WebsiteAlertConfig{
			ID:              websiteAlertConfigID,
			AlertChannelIDs: []string{"channel-1", "channel-2"},
			WebsiteID:       websiteID,
			Description:     "website-alert-config-description",
			Granularity:     restapi.Granularity600000,
			CustomerPayloadFields: []restapi.CustomPayloadField[restapi.StaticStringCustomPayloadFieldValue]{
				{
					Type:  restapi.StaticCustomPayloadType,
					Key:   "static-key",
					Value: restapi.StaticStringCustomPayloadFieldValue("static-value"),
				},
			},
			Name:                fullName,
			Rule:                ruleTestPair.input,
			Severity:            restapi.SeverityCritical.GetAPIRepresentation(),
			TagFilterExpression: restapi.NewStringTagFilter(restapi.TagFilterEntitySource, "service.name", restapi.EqualsOperator, "test"),
			Threshold:           thresholdTestPair.input,
			TimeThreshold:       timeThresholdTestPair.input,
			Triggering:          true,
		}

		testHelper := NewTestHelper[*restapi.WebsiteAlertConfig](t)
		sut := test.resourceHandle
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

		err := sut.UpdateState(resourceData, &websiteConfig)

		require.NoError(t, err)
		require.Equal(t, websiteAlertConfigID, resourceData.Id())
		require.Equal(t, []interface{}{"channel-2", "channel-1"}, (resourceData.Get(WebsiteAlertConfigFieldAlertChannelIDs).(*schema.Set)).List())
		require.Equal(t, websiteID, resourceData.Get(WebsiteAlertConfigFieldWebsiteID))
		require.Equal(t, "website-alert-config-description", resourceData.Get(WebsiteAlertConfigFieldDescription))
		require.Equal(t, "website-alert-config-name", resourceData.Get(WebsiteAlertConfigFieldName))
		require.Equal(t, []interface{}{
			map[string]interface{}{WebsiteAlertConfigFieldCustomPayloadFieldsKey: "static-key", WebsiteAlertConfigFieldCustomPayloadFieldsValue: "static-value"},
		}, resourceData.Get(WebsiteAlertConfigFieldCustomPayloadFields).(*schema.Set).List())
		require.Equal(t, ruleTestPair.expected, resourceData.Get(WebsiteAlertConfigFieldRule))
		require.Equal(t, restapi.SeverityCritical.GetTerraformRepresentation(), resourceData.Get(WebsiteAlertConfigFieldSeverity))
		require.Equal(t, "service.name@src EQUALS 'test'", resourceData.Get(WebsiteAlertConfigFieldTagFilter))
		test.requireWebsiteAlertConfigThresholdSetOnSchema(t, thresholdTestPair.expected, resourceData)
		require.Equal(t, timeThresholdTestPair.expected, resourceData.Get(WebsiteAlertConfigFieldTimeThreshold))
		require.True(t, resourceData.Get(WebsiteAlertConfigFieldTriggering).(bool))
	}
}

func (test *websiteAlertConfigTest) requireWebsiteAlertConfigThresholdSetOnSchema(t *testing.T, expected []interface{}, resourceData *schema.ResourceData) {
	actual := resourceData.Get(ResourceFieldThreshold).([]interface{})
	require.Equal(t, 1, len(expected))
	require.Equal(t, len(expected), len(actual))
	expectedEntries := expected[0].(map[string]interface{})
	actualEntries := actual[0].(map[string]interface{})

	for k, e := range expectedEntries {
		if k == ResourceFieldThresholdHistoricBaseline && len(e.([]interface{})) > 0 {
			expectedHistoricBaselineSlice := e.([]interface{})
			actualHistoricBaselineSlice := actualEntries[k].([]interface{})
			require.Equal(t, 1, len(expectedHistoricBaselineSlice))
			require.Equal(t, len(expected), len(actual))
			expectedHistoricBaseline := expectedHistoricBaselineSlice[0].(map[string]interface{})
			actualHistoricBaseline := actualHistoricBaselineSlice[0].(map[string]interface{})
			for key, expectedBaselineValue := range expectedHistoricBaseline {
				if key == ResourceFieldThresholdHistoricBaselineBaseline {
					actualBaseline := actualHistoricBaseline[key].(*schema.Set)
					actualBaselineSlice := make([][]float64, actualBaseline.Len())
					for i, v := range actualBaseline.List() {
						values := v.(*schema.Set).List()
						actualBaselineSlice[i] = ConvertInterfaceSlice[float64](values)
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

func (test *websiteAlertConfigTest) createTestCasesShouldFailToUpdateTerraformResourceStateFromModeWhenSeverityIsNotValid() func(t *testing.T) {
	return func(t *testing.T) {
		websiteConfig := restapi.WebsiteAlertConfig{
			Name:     "prefix test suffix",
			Severity: -1,
		}

		testHelper := NewTestHelper[*restapi.WebsiteAlertConfig](t)
		sut := test.resourceHandle
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

		err := sut.UpdateState(resourceData, &websiteConfig)

		require.Error(t, err)
		require.Equal(t, "-1 is not a valid severity", err.Error())
	}
}

func (test *websiteAlertConfigTest) createTestCasesShouldFailToUpdateTerraformResourceStateFromModeWhenTagFilterExpressionIsNotValid() func(t *testing.T) {
	return func(t *testing.T) {
		value := "test"
		websiteConfig := restapi.WebsiteAlertConfig{
			Name:     "prefix test suffix",
			Severity: restapi.SeverityWarning.GetAPIRepresentation(),
			TagFilterExpression: &restapi.TagFilter{
				Entity:      restapi.TagFilterEntitySource,
				Name:        "service.name",
				Operator:    restapi.EqualsOperator,
				StringValue: &value,
				Value:       value,
				Type:        restapi.TagFilterExpressionElementType("invalid"),
			},
		}

		testHelper := NewTestHelper[*restapi.WebsiteAlertConfig](t)
		sut := test.resourceHandle
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

		err := sut.UpdateState(resourceData, &websiteConfig)

		require.Error(t, err)
		require.Equal(t, "unsupported tag filter expression of type invalid", err.Error())
	}
}

func (test *websiteAlertConfigTest) createTestCasesForMappingOfTerraformResourceStateToModel(t *testing.T) {
	metricName := "test-metric"
	equalsOperator := restapi.EqualsOperator
	minAggregation := restapi.MinAggregation
	statusCodeValue := "200"
	jsErrorValue := "jsErrorValue"
	rules := []testPair[[]map[string]interface{}, restapi.WebsiteAlertRule]{
		{
			name: WebsiteAlertConfigFieldRuleThroughput,
			expected: restapi.WebsiteAlertRule{
				AlertType:   WebsiteAlertConfigFieldRuleThroughput,
				Aggregation: &minAggregation,
				MetricName:  metricName,
			},
			input: []map[string]interface{}{
				{
					WebsiteAlertConfigFieldRuleSpecificJsError: []interface{}{},
					WebsiteAlertConfigFieldRuleSlowness:        []interface{}{},
					WebsiteAlertConfigFieldRuleStatusCode:      []interface{}{},
					WebsiteAlertConfigFieldRuleThroughput: []interface{}{
						map[string]interface{}{
							WebsiteAlertConfigFieldRuleAggregation: string(minAggregation),
							WebsiteAlertConfigFieldRuleMetricName:  metricName,
						},
					},
				},
			},
		},
		{
			name: "StatusCode",
			expected: restapi.WebsiteAlertRule{
				AlertType:   "statusCode",
				Aggregation: &minAggregation,
				MetricName:  metricName,
				Operator:    &equalsOperator,
				Value:       &statusCodeValue,
			},
			input: []map[string]interface{}{
				{
					WebsiteAlertConfigFieldRuleSpecificJsError: []interface{}{},
					WebsiteAlertConfigFieldRuleSlowness:        []interface{}{},
					WebsiteAlertConfigFieldRuleStatusCode: []interface{}{
						map[string]interface{}{
							WebsiteAlertConfigFieldRuleAggregation: string(minAggregation),
							WebsiteAlertConfigFieldRuleMetricName:  metricName,
							WebsiteAlertConfigFieldRuleOperator:    string(equalsOperator),
							WebsiteAlertConfigFieldRuleValue:       statusCodeValue,
						},
					},
					WebsiteAlertConfigFieldRuleThroughput: []interface{}{},
				},
			},
		},
		{
			name: WebsiteAlertConfigFieldRuleSlowness,
			expected: restapi.WebsiteAlertRule{
				AlertType:   WebsiteAlertConfigFieldRuleSlowness,
				Aggregation: &minAggregation,
				MetricName:  metricName,
			},
			input: []map[string]interface{}{
				{
					WebsiteAlertConfigFieldRuleSpecificJsError: []interface{}{},
					WebsiteAlertConfigFieldRuleSlowness: []interface{}{
						map[string]interface{}{
							WebsiteAlertConfigFieldRuleAggregation: string(minAggregation),
							WebsiteAlertConfigFieldRuleMetricName:  metricName,
						},
					},
					WebsiteAlertConfigFieldRuleStatusCode: []interface{}{},
					WebsiteAlertConfigFieldRuleThroughput: []interface{}{},
				},
			},
		},
		{
			name: "SpecificJsError",
			expected: restapi.WebsiteAlertRule{
				AlertType:   "specificJsError",
				Aggregation: &minAggregation,
				MetricName:  metricName,
				Operator:    &equalsOperator,
				Value:       &jsErrorValue,
			},
			input: []map[string]interface{}{
				{
					WebsiteAlertConfigFieldRuleSpecificJsError: []interface{}{
						map[string]interface{}{
							WebsiteAlertConfigFieldRuleAggregation: string(minAggregation),
							WebsiteAlertConfigFieldRuleMetricName:  metricName,
							WebsiteAlertConfigFieldRuleOperator:    string(equalsOperator),
							WebsiteAlertConfigFieldRuleValue:       jsErrorValue,
						},
					},
					WebsiteAlertConfigFieldRuleSlowness:   []interface{}{},
					WebsiteAlertConfigFieldRuleStatusCode: []interface{}{},
					WebsiteAlertConfigFieldRuleThroughput: []interface{}{},
				},
			},
		},
	}

	thresholdValue := 123.3
	thresholdLastUpdate := int64(12345)
	thresholdSeasonality := restapi.ThresholdSeasonalityDaily
	thresholdBaseline := [][]float64{{1.23, 4.56}, {1.23, 7.89}}
	thresholdDeviationFactor := float32(1.23)
	thresholds := []testPair[[]map[string]interface{}, restapi.Threshold]{
		{
			name: "StaticThreshold",
			expected: restapi.Threshold{
				Type:        "staticThreshold",
				Operator:    restapi.ThresholdOperatorGreaterThan,
				LastUpdated: &thresholdLastUpdate,
				Value:       &thresholdValue,
			},
			input: []map[string]interface{}{
				{
					ResourceFieldThresholdHistoricBaseline: []interface{}{},
					ResourceFieldThresholdStatic: []interface{}{
						map[string]interface{}{
							ResourceFieldThresholdLastUpdated: int(thresholdLastUpdate),
							ResourceFieldThresholdOperator:    string(restapi.ThresholdOperatorGreaterThan),
							ResourceFieldThresholdStaticValue: thresholdValue,
						},
					},
				},
			},
		},
		{
			name: "HistoricBaseLine",
			expected: restapi.Threshold{
				Type:            "historicBaseline",
				Operator:        restapi.ThresholdOperatorGreaterThan,
				LastUpdated:     &thresholdLastUpdate,
				Seasonality:     &thresholdSeasonality,
				Baseline:        &thresholdBaseline,
				DeviationFactor: &thresholdDeviationFactor,
			},
			input: []map[string]interface{}{
				{
					ResourceFieldThresholdHistoricBaseline: []interface{}{
						map[string]interface{}{
							ResourceFieldThresholdLastUpdated:                     int(thresholdLastUpdate),
							ResourceFieldThresholdOperator:                        string(restapi.ThresholdOperatorGreaterThan),
							ResourceFieldThresholdHistoricBaselineSeasonality:     string(thresholdSeasonality),
							ResourceFieldThresholdHistoricBaselineBaseline:        thresholdBaseline,
							ResourceFieldThresholdHistoricBaselineDeviationFactor: float64(thresholdDeviationFactor),
						},
					},
					ResourceFieldThresholdStatic: []interface{}{},
				},
			},
		},
	}

	timeThresholdWindow := int64(12345)
	timeThresholdImpactMeasurementMethod := restapi.WebsiteImpactMeasurementMethodAggregated
	timeThresholdUsers := int32(5)
	timeThresholdUserPercentage := 0.8
	timeThresholdViolations := int32(3)
	timeThresholds := []testPair[[]map[string]interface{}, restapi.WebsiteTimeThreshold]{
		{
			name: "UserImpactOfViolationsInSequence",
			expected: restapi.WebsiteTimeThreshold{
				Type:                    "userImpactOfViolationsInSequence",
				TimeWindow:              &timeThresholdWindow,
				ImpactMeasurementMethod: &timeThresholdImpactMeasurementMethod,
				Users:                   &timeThresholdUsers,
				UserPercentage:          &timeThresholdUserPercentage,
			},
			input: []map[string]interface{}{
				{
					WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequence: []interface{}{
						map[string]interface{}{
							WebsiteAlertConfigFieldTimeThresholdTimeWindow:                                              int(timeThresholdWindow),
							WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceImpactMeasurementMethod: string(timeThresholdImpactMeasurementMethod),
							WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceUserPercentage:          timeThresholdUserPercentage,
							WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceUsers:                   int(timeThresholdUsers),
						},
					},
					WebsiteAlertConfigFieldTimeThresholdViolationsInPeriod:   []interface{}{},
					WebsiteAlertConfigFieldTimeThresholdViolationsInSequence: []interface{}{},
				},
			},
		},
		{
			name: "ViolationsInPeriod",
			expected: restapi.WebsiteTimeThreshold{
				Type:       "violationsInPeriod",
				TimeWindow: &timeThresholdWindow,
				Violations: &timeThresholdViolations,
			},
			input: []map[string]interface{}{
				{
					WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequence: []interface{}{},
					WebsiteAlertConfigFieldTimeThresholdViolationsInPeriod: []interface{}{
						map[string]interface{}{
							WebsiteAlertConfigFieldTimeThresholdTimeWindow:                   int(timeThresholdWindow),
							WebsiteAlertConfigFieldTimeThresholdViolationsInPeriodViolations: int(timeThresholdViolations),
						},
					},
					WebsiteAlertConfigFieldTimeThresholdViolationsInSequence: []interface{}{},
				},
			},
		},
		{
			name: "ViolationsInSequence",
			expected: restapi.WebsiteTimeThreshold{
				Type:       "violationsInSequence",
				TimeWindow: &timeThresholdWindow,
			},
			input: []map[string]interface{}{
				{
					WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequence: []interface{}{},
					WebsiteAlertConfigFieldTimeThresholdViolationsInPeriod:               []interface{}{},
					WebsiteAlertConfigFieldTimeThresholdViolationsInSequence: []interface{}{
						map[string]interface{}{
							WebsiteAlertConfigFieldTimeThresholdTimeWindow: int(timeThresholdWindow),
						},
					},
				},
			},
		},
	}

	for _, rule := range rules {
		for _, threshold := range thresholds {
			for _, timeThreshold := range timeThresholds {
				t.Run(fmt.Sprintf("Should map terraform state of %s to REST model with %s and %s and %s", ResourceInstanaWebsiteAlertConfig, rule.name, threshold.name, timeThreshold.name), test.createTestShouldMapTerraformResourceStateToModelCase(rule, threshold, timeThreshold))
			}
		}
	}
}

func (test *websiteAlertConfigTest) createTestShouldMapTerraformResourceStateToModelCase(ruleTestPair testPair[[]map[string]interface{}, restapi.WebsiteAlertRule],
	thresholdTestPair testPair[[]map[string]interface{}, restapi.Threshold],
	timeThresholdTestPair testPair[[]map[string]interface{}, restapi.WebsiteTimeThreshold]) func(t *testing.T) {
	return func(t *testing.T) {
		websiteAlertConfigID := "website-alert-config-id"
		websiteID := "website-id"
		expectedWebsiteConfig := restapi.WebsiteAlertConfig{
			ID:              websiteAlertConfigID,
			AlertChannelIDs: []string{"channel-2", "channel-1"},
			WebsiteID:       websiteID,
			Description:     "website-alert-config-description",
			Granularity:     restapi.Granularity600000,
			CustomerPayloadFields: []restapi.CustomPayloadField[restapi.StaticStringCustomPayloadFieldValue]{
				{
					Type:  restapi.StaticCustomPayloadType,
					Key:   "static-key",
					Value: restapi.StaticStringCustomPayloadFieldValue("static-value"),
				},
			},
			Name:                "website-alert-config-name",
			Rule:                ruleTestPair.expected,
			Severity:            restapi.SeverityCritical.GetAPIRepresentation(),
			TagFilterExpression: restapi.NewStringTagFilter(restapi.TagFilterEntitySource, "service.name", restapi.EqualsOperator, "test"),
			Threshold:           thresholdTestPair.expected,
			TimeThreshold:       timeThresholdTestPair.expected,
			Triggering:          true,
		}

		testHelper := NewTestHelper[*restapi.WebsiteAlertConfig](t)
		sut := test.resourceHandle
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)
		setValueOnResourceData(t, resourceData, WebsiteAlertConfigFieldAlertChannelIDs, []interface{}{"channel-2", "channel-1"})
		setValueOnResourceData(t, resourceData, WebsiteAlertConfigFieldWebsiteID, websiteID)
		setValueOnResourceData(t, resourceData, WebsiteAlertConfigFieldCustomPayloadFields, []interface{}{
			map[string]interface{}{WebsiteAlertConfigFieldCustomPayloadFieldsKey: "static-key", WebsiteAlertConfigFieldCustomPayloadFieldsValue: "static-value"},
		})
		setValueOnResourceData(t, resourceData, WebsiteAlertConfigFieldDescription, "website-alert-config-description")
		setValueOnResourceData(t, resourceData, WebsiteAlertConfigFieldGranularity, restapi.Granularity600000)
		setValueOnResourceData(t, resourceData, WebsiteAlertConfigFieldName, "website-alert-config-name")
		setValueOnResourceData(t, resourceData, WebsiteAlertConfigFieldRule, ruleTestPair.input)
		setValueOnResourceData(t, resourceData, WebsiteAlertConfigFieldSeverity, restapi.SeverityCritical.GetTerraformRepresentation())
		setValueOnResourceData(t, resourceData, WebsiteAlertConfigFieldTagFilter, "service.name@src EQUALS 'test'")
		setValueOnResourceData(t, resourceData, ResourceFieldThreshold, thresholdTestPair.input)
		setValueOnResourceData(t, resourceData, WebsiteAlertConfigFieldTimeThreshold, timeThresholdTestPair.input)
		setValueOnResourceData(t, resourceData, WebsiteAlertConfigFieldTriggering, true)
		resourceData.SetId(websiteAlertConfigID)

		result, err := sut.MapStateToDataObject(resourceData)

		require.NoError(t, err)
		require.Equal(t, &expectedWebsiteConfig, result)
	}
}

func (test *websiteAlertConfigTest) createTestCaseShouldFailToMapTerraformResourceStateToModelWhenSeverityIsNotValid() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.WebsiteAlertConfig](t)
		sut := test.resourceHandle
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)
		setValueOnResourceData(t, resourceData, WebsiteAlertConfigFieldName, "website-alert-config-name")
		setValueOnResourceData(t, resourceData, WebsiteAlertConfigFieldSeverity, "invalid")

		_, err := sut.MapStateToDataObject(resourceData)

		require.Error(t, err)
		require.Equal(t, "invalid is not a valid severity", err.Error())
	}
}

func (test *websiteAlertConfigTest) createTestCaseShouldFailToMapTerraformResourceStateToModelWhenTagFilterIsNotValid() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper[*restapi.WebsiteAlertConfig](t)
		sut := test.resourceHandle
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)
		setValueOnResourceData(t, resourceData, WebsiteAlertConfigFieldName, "website-alert-config-name")
		setValueOnResourceData(t, resourceData, WebsiteAlertConfigFieldSeverity, restapi.SeverityWarning.GetTerraformRepresentation())
		setValueOnResourceData(t, resourceData, WebsiteAlertConfigFieldTagFilter, "invalid invalid invalid")

		_, err := sut.MapStateToDataObject(resourceData)

		require.Error(t, err)
		require.Contains(t, err.Error(), "unexpected token")
	}
}
