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

func createApplicationAlertConfigTestFor(terraformResourceName string, resourceRestAPIPath string, resourceHandle ResourceHandle[*restapi.ApplicationAlertConfig]) *anyApplicationConfigTest {
	terraformResourceInstanceName := terraformResourceName + ".example"
	resourceInstanceRestAPIPath := resourceRestAPIPath + "/{internal-id}"
	inst := &anyApplicationConfigTest{
		terraformResourceName:         terraformResourceName,
		terraformResourceInstanceName: terraformResourceInstanceName,
		resourceRestAPIPath:           resourceRestAPIPath,
		resourceInstanceRestAPIPath:   resourceInstanceRestAPIPath,
		resourceHandle:                resourceHandle,
	}
	return inst
}

type anyApplicationConfigTest struct {
	terraformResourceName         string
	terraformResourceInstanceName string
	resourceRestAPIPath           string
	resourceInstanceRestAPIPath   string
	resourceHandle                ResourceHandle[*restapi.ApplicationAlertConfig]
}

var applicationAlertConfigTerraformTemplate = `
resource "%s" "example" {
	name              = "name %d"
    description       = "test-alert-description"
    boundary_scope    = "ALL"
    severity          = "warning"
    triggering        = false
    include_internal  = false
    include_synthetic = false
    alert_channel_ids = [ "alert-channel-id-1", "alert-channel-id-2" ]
    granularity       = 600000
	evaluation_type   = "PER_AP"

	tag_filter        = "call.type@na EQUALS 'HTTP'"
    
    application {
		application_id = "app-id"
		inclusive 	   = true
		
        service {
			service_id = "service-1-id"
			inclusive  = true

			endpoint {
				endpoint_id = "endpoint-1-1-id"
				inclusive   = true
			}
        }
		
        service {
			service_id = "service-2-id"
			inclusive  = true
        }
	}

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

var applicationAlertConfigServerResponseTemplate = `
	{
    "id": "%s",
    "name": "name %d",
    "description": "test-alert-description",
    "boundaryScope": "ALL",
    "applicationId": "app-id",
    "applications": {
      "app-id": {
        "applicationId": "app-id",
        "inclusive": true,
        "services": {
			"service-1-id": {
				"serviceId": "service-1-id",
				"inclusive": true,
				"endpoints": {
					"endpoint-1-1-id": {
						"endpointId": "endpoint-1-1-id",
					    "inclusive": true
					}
				}
			},
			"service-2-id": {
				"serviceId": "service-2-id",
				"inclusive": true
			}
		}
      }
    },
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
    "includeInternal": false,
    "includeSynthetic": false,
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
    "evaluationType": "PER_AP",
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

func (f *anyApplicationConfigTest) run(t *testing.T) {
	t.Run(fmt.Sprintf("CRUDD integration test of %s", f.terraformResourceName), f.createIntegrationTest())
	t.Run(fmt.Sprintf("DiffSuppressFunc of TagFilter of %s should return true when value can be normalized and old and new normalized value are equal", f.terraformResourceName), f.createTestOfDiffSuppressFuncOfTagFilterShouldReturnTrueWhenValueCanBeNormalizedAndOldAndNewNormalizedValueAreEqual())
	t.Run(fmt.Sprintf("DiffSuppressFunc of TagFilter of %s should return false when value can be normalized and old and new normalized value are not equal", f.terraformResourceName), f.createTestOfDiffSuppressFuncOfTagFilterShouldReturnFalseWhenValueCanBeNormalizedAndOldAndNewNormalizedValueAreNotEqual())
	t.Run(fmt.Sprintf("DiffSuppressFunc of TagFilter of %s should return true when value can be normalized and old and new value are equal", f.terraformResourceName), f.createTestOfDiffSuppressFuncOfTagFilterShouldReturnTrueWhenValueCannotBeNormalizedAndOldAndNewValueAreEqual())
	t.Run(fmt.Sprintf("DiffSuppressFunc of TagFilter of %s should return false when value cannot be normalized and old and new value are not equal", f.terraformResourceName), f.createTestOfDiffSuppressFuncOfTagFilterShouldReturnFalseWhenValueCannotBeNormalizedAndOldAndNewValueAreNotEqual())
	t.Run(fmt.Sprintf("StateFunc of TagFilter of %s should return normalized value when value can be normalized", f.terraformResourceName), f.createTestOfStateFuncOfTagFilterShouldReturnNormalizedValueWhenValueCanBeNormalized())
	t.Run(fmt.Sprintf("StateFunc of TagFilter of %s should return provided value when value cannot be normalized", f.terraformResourceName), f.createTestOfStateFuncOfTagFilterShouldReturnProvidedValueWhenValueCannotBeNormalized())
	t.Run(fmt.Sprintf("ValidateFunc of TagFilter of %s should return no errors and warnings when value can be parsed", f.terraformResourceName), f.createTestOfValidateFuncOfTagFilterShouldReturnNoErrorsAndWarningsWhenValueCanBeParsed())
	t.Run(fmt.Sprintf("ValidateFunc of TagFilter of %s should return one error and no warnings when value can be parsed", f.terraformResourceName), f.createTestOfValidateFuncOfTagFilterShouldReturnOneErrorAndNoWarningsWhenValueCannotBeParsed())
	t.Run(fmt.Sprintf("%s should have schema version zero", f.terraformResourceName), f.createTetResourceShouldHaveSchemaVersionOne())
	t.Run(fmt.Sprintf("%s should have no state upgrader", f.terraformResourceName), f.createTetResourceShouldHaveOneStateUpgrader())
	t.Run(fmt.Sprintf("%s should migrate full_name to name when executing first state upgrader and full_name is available", ResourceInstanaCustomDashboard), f.createTestResourceShouldMigrateFullNameToNameWhenExecutingFirstStateUpgraderAndFullNameIsAvailable())
	t.Run(fmt.Sprintf("%s should do nothing when executing first state upgrader and full_name is not available", ResourceInstanaCustomDashboard), f.createTestResourceShouldDoNothingWhenExecutingFirstStateUpgraderAndFullNameIsNotAvailable())
	t.Run(fmt.Sprintf("%s should have correct resouce name", f.terraformResourceName), f.createTetResourceShouldHaveCorrectResourceName())
	f.createTestCasesForUpdatesOfTerraformResourceStateFromModel(t)
	f.createTestCasesForMappingOfTerraformResourceStateToModel(t)
}

func (f *anyApplicationConfigTest) createIntegrationTest() func(t *testing.T) {
	return func(t *testing.T) {
		id := RandomID()
		httpServer := testutils.NewTestHTTPServer()
		httpServer.AddRoute(http.MethodPost, f.resourceRestAPIPath, func(w http.ResponseWriter, r *http.Request) {
			config := &restapi.ApplicationAlertConfig{}
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
		httpServer.AddRoute(http.MethodPost, f.resourceInstanceRestAPIPath, func(w http.ResponseWriter, r *http.Request) {
			testutils.EchoHandlerFunc(w, r)
		})
		httpServer.AddRoute(http.MethodDelete, f.resourceInstanceRestAPIPath, testutils.EchoHandlerFunc)
		httpServer.AddRoute(http.MethodGet, f.resourceInstanceRestAPIPath, func(w http.ResponseWriter, r *http.Request) {
			modCount := httpServer.GetCallCount(http.MethodPost, f.resourceRestAPIPath+"/"+id)
			jsonData := fmt.Sprintf(applicationAlertConfigServerResponseTemplate, id, modCount)
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
				f.createIntegrationTestStep(httpServer.GetPort(), 0, id),
				testStepImportWithCustomID(f.terraformResourceInstanceName, id),
				f.createIntegrationTestStep(httpServer.GetPort(), 1, id),
				testStepImportWithCustomID(f.terraformResourceInstanceName, id),
			},
		})
	}
}

func (f *anyApplicationConfigTest) createIntegrationTestStep(httpPort int64, iteration int, id string) resource.TestStep {
	application1ApplicationId := fmt.Sprintf("%s.%d.%s", ApplicationAlertConfigFieldApplications, 0, ApplicationAlertConfigFieldApplicationsApplicationID)
	application1Inclusive := fmt.Sprintf("%s.%d.%s", ApplicationAlertConfigFieldApplications, 0, ApplicationAlertConfigFieldApplicationsInclusive)
	application1Service1ServiceId := fmt.Sprintf("%s.%d.%s.%d.%s", ApplicationAlertConfigFieldApplications, 0, ApplicationAlertConfigFieldApplicationsServices, 0, ApplicationAlertConfigFieldApplicationsServicesServiceID)
	application1Service1Inclusive := fmt.Sprintf("%s.%d.%s.%d.%s", ApplicationAlertConfigFieldApplications, 0, ApplicationAlertConfigFieldApplicationsServices, 0, ApplicationAlertConfigFieldApplicationsInclusive)
	application1Service2ServiceId := fmt.Sprintf("%s.%d.%s.%d.%s", ApplicationAlertConfigFieldApplications, 0, ApplicationAlertConfigFieldApplicationsServices, 1, ApplicationAlertConfigFieldApplicationsServicesServiceID)
	application1Service2Inclusive := fmt.Sprintf("%s.%d.%s.%d.%s", ApplicationAlertConfigFieldApplications, 0, ApplicationAlertConfigFieldApplicationsServices, 1, ApplicationAlertConfigFieldApplicationsInclusive)
	application1Service1Endpoint1EndpointId := fmt.Sprintf("%s.%d.%s.%d.%s.%d.%s", ApplicationAlertConfigFieldApplications, 0, ApplicationAlertConfigFieldApplicationsServices, 0, ApplicationAlertConfigFieldApplicationsServicesEndpoints, 0, ApplicationAlertConfigFieldApplicationsServicesEndpointsEndpointID)
	application1Service1Endpoint1Inclusive := fmt.Sprintf("%s.%d.%s.%d.%s.%d.%s", ApplicationAlertConfigFieldApplications, 0, ApplicationAlertConfigFieldApplicationsServices, 0, ApplicationAlertConfigFieldApplicationsServicesEndpoints, 0, ApplicationAlertConfigFieldApplicationsInclusive)
	ruleSlownessMetricName := fmt.Sprintf("%s.%d.%s.%d.%s", ApplicationAlertConfigFieldRule, 0, ApplicationAlertConfigFieldRuleSlowness, 0, ApplicationAlertConfigFieldRuleMetricName)
	ruleSlownessAggregation := fmt.Sprintf("%s.%d.%s.%d.%s", ApplicationAlertConfigFieldRule, 0, ApplicationAlertConfigFieldRuleSlowness, 0, ApplicationAlertConfigFieldRuleAggregation)
	thresholdStaticOperator := fmt.Sprintf("%s.%d.%s.%d.%s", ResourceFieldThreshold, 0, ResourceFieldThresholdStatic, 0, ResourceFieldThresholdOperator)
	thresholdStaticValue := fmt.Sprintf("%s.%d.%s.%d.%s", ResourceFieldThreshold, 0, ResourceFieldThresholdStatic, 0, ResourceFieldThresholdStaticValue)
	timeThresholdViolationsInSequence := fmt.Sprintf("%s.%d.%s.%d.%s", ApplicationAlertConfigFieldTimeThreshold, 0, ApplicationAlertConfigFieldTimeThresholdViolationsInSequence, 0, ApplicationAlertConfigFieldTimeThresholdTimeWindow)
	customPayloadFieldStaticKey := fmt.Sprintf("%s.%d.%s", ApplicationAlertConfigFieldCustomPayloadFields, 0, ApplicationAlertConfigFieldCustomPayloadFieldsKey)
	customPayloadFieldStaticValue := fmt.Sprintf("%s.%d.%s", ApplicationAlertConfigFieldCustomPayloadFields, 0, ApplicationAlertConfigFieldCustomPayloadFieldsValue)
	return resource.TestStep{
		Config: appendProviderConfig(fmt.Sprintf(applicationAlertConfigTerraformTemplate, f.terraformResourceName, iteration), httpPort),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(f.terraformResourceInstanceName, "id", id),
			resource.TestCheckResourceAttr(f.terraformResourceInstanceName, ApplicationAlertConfigFieldName, formatResourceName(iteration)),
			resource.TestCheckResourceAttr(f.terraformResourceInstanceName, ApplicationAlertConfigFieldDescription, "test-alert-description"),
			resource.TestCheckResourceAttr(f.terraformResourceInstanceName, ApplicationAlertConfigFieldBoundaryScope, string(restapi.BoundaryScopeAll)),
			resource.TestCheckResourceAttr(f.terraformResourceInstanceName, ApplicationAlertConfigFieldSeverity, restapi.SeverityWarning.GetTerraformRepresentation()),
			resource.TestCheckResourceAttr(f.terraformResourceInstanceName, ApplicationAlertConfigFieldTriggering, falseAsString),
			resource.TestCheckResourceAttr(f.terraformResourceInstanceName, ApplicationAlertConfigFieldIncludeInternal, falseAsString),
			resource.TestCheckResourceAttr(f.terraformResourceInstanceName, ApplicationAlertConfigFieldIncludeSynthetic, falseAsString),
			resource.TestCheckResourceAttr(f.terraformResourceInstanceName, ApplicationAlertConfigFieldAlertChannelIDs+".0", "alert-channel-id-1"),
			resource.TestCheckResourceAttr(f.terraformResourceInstanceName, ApplicationAlertConfigFieldAlertChannelIDs+".1", "alert-channel-id-2"),
			resource.TestCheckResourceAttr(f.terraformResourceInstanceName, ApplicationAlertConfigFieldGranularity, "600000"),
			resource.TestCheckResourceAttr(f.terraformResourceInstanceName, ApplicationAlertConfigFieldEvaluationType, string(restapi.EvaluationTypePerApplication)),
			resource.TestCheckResourceAttr(f.terraformResourceInstanceName, ApplicationAlertConfigFieldTagFilter, "call.type@na EQUALS 'HTTP'"),
			resource.TestCheckResourceAttr(f.terraformResourceInstanceName, application1ApplicationId, "app-id"),
			resource.TestCheckResourceAttr(f.terraformResourceInstanceName, application1Inclusive, trueAsString),
			resource.TestCheckResourceAttr(f.terraformResourceInstanceName, application1Service1ServiceId, "service-1-id"),
			resource.TestCheckResourceAttr(f.terraformResourceInstanceName, application1Service1Inclusive, trueAsString),
			resource.TestCheckResourceAttr(f.terraformResourceInstanceName, application1Service2ServiceId, "service-2-id"),
			resource.TestCheckResourceAttr(f.terraformResourceInstanceName, application1Service2Inclusive, trueAsString),
			resource.TestCheckResourceAttr(f.terraformResourceInstanceName, application1Service1Endpoint1EndpointId, "endpoint-1-1-id"),
			resource.TestCheckResourceAttr(f.terraformResourceInstanceName, application1Service1Endpoint1Inclusive, trueAsString),
			resource.TestCheckResourceAttr(f.terraformResourceInstanceName, ruleSlownessMetricName, "latency"),
			resource.TestCheckResourceAttr(f.terraformResourceInstanceName, ruleSlownessAggregation, "P90"),
			resource.TestCheckResourceAttr(f.terraformResourceInstanceName, thresholdStaticOperator, ">="),
			resource.TestCheckResourceAttr(f.terraformResourceInstanceName, thresholdStaticValue, "5"),
			resource.TestCheckResourceAttr(f.terraformResourceInstanceName, timeThresholdViolationsInSequence, "600000"),
			resource.TestCheckResourceAttr(f.terraformResourceInstanceName, customPayloadFieldStaticKey, "test"),
			resource.TestCheckResourceAttr(f.terraformResourceInstanceName, customPayloadFieldStaticValue, "test123"),
		),
	}
}

func (f *anyApplicationConfigTest) createTestOfDiffSuppressFuncOfTagFilterShouldReturnTrueWhenValueCanBeNormalizedAndOldAndNewNormalizedValueAreEqual() func(t *testing.T) {
	return func(t *testing.T) {
		resourceSchema := f.resourceHandle.MetaData().Schema
		oldValue := expressionEntityTypeDestEqValue
		newValue := "entity.type  EQUALS    'foo'"

		require.True(t, resourceSchema[ApplicationAlertConfigFieldTagFilter].DiffSuppressFunc(ApplicationAlertConfigFieldTagFilter, oldValue, newValue, nil))
	}
}

func (f *anyApplicationConfigTest) createTestOfDiffSuppressFuncOfTagFilterShouldReturnFalseWhenValueCanBeNormalizedAndOldAndNewNormalizedValueAreNotEqual() func(t *testing.T) {
	return func(t *testing.T) {
		resourceSchema := f.resourceHandle.MetaData().Schema
		oldValue := expressionEntityTypeSrcEqValue
		newValue := validTagFilter

		require.False(t, resourceSchema[ApplicationAlertConfigFieldTagFilter].DiffSuppressFunc(ApplicationAlertConfigFieldTagFilter, oldValue, newValue, nil))
	}
}

func (f *anyApplicationConfigTest) createTestOfDiffSuppressFuncOfTagFilterShouldReturnTrueWhenValueCannotBeNormalizedAndOldAndNewValueAreEqual() func(t *testing.T) {
	return func(t *testing.T) {
		resourceSchema := f.resourceHandle.MetaData().Schema
		invalidValue := invalidTagFilter

		require.True(t, resourceSchema[ApplicationAlertConfigFieldTagFilter].DiffSuppressFunc(ApplicationAlertConfigFieldTagFilter, invalidValue, invalidValue, nil))
	}
}

func (f *anyApplicationConfigTest) createTestOfDiffSuppressFuncOfTagFilterShouldReturnFalseWhenValueCannotBeNormalizedAndOldAndNewValueAreNotEqual() func(t *testing.T) {
	return func(t *testing.T) {
		resourceSchema := f.resourceHandle.MetaData().Schema
		oldValue := invalidTagFilter
		newValue := "entity.type foo foo foo"

		require.False(t, resourceSchema[ApplicationAlertConfigFieldTagFilter].DiffSuppressFunc(ApplicationAlertConfigFieldTagFilter, oldValue, newValue, nil))
	}
}

func (f *anyApplicationConfigTest) createTestOfStateFuncOfTagFilterShouldReturnNormalizedValueWhenValueCanBeNormalized() func(t *testing.T) {
	return func(t *testing.T) {
		resourceSchema := f.resourceHandle.MetaData().Schema
		expectedValue := expressionEntityTypeDestEqValue
		newValue := validTagFilter

		require.Equal(t, expectedValue, resourceSchema[ApplicationAlertConfigFieldTagFilter].StateFunc(newValue))
	}
}

func (f *anyApplicationConfigTest) createTestOfStateFuncOfTagFilterShouldReturnProvidedValueWhenValueCannotBeNormalized() func(t *testing.T) {
	return func(t *testing.T) {
		resourceSchema := f.resourceHandle.MetaData().Schema
		value := invalidTagFilter

		require.Equal(t, value, resourceSchema[ApplicationAlertConfigFieldTagFilter].StateFunc(value))
	}
}

func (f *anyApplicationConfigTest) createTestOfValidateFuncOfTagFilterShouldReturnNoErrorsAndWarningsWhenValueCanBeParsed() func(t *testing.T) {
	return func(t *testing.T) {
		resourceSchema := f.resourceHandle.MetaData().Schema
		value := validTagFilter

		warns, errs := resourceSchema[ApplicationAlertConfigFieldTagFilter].ValidateFunc(value, ApplicationAlertConfigFieldTagFilter)
		require.Empty(t, warns)
		require.Empty(t, errs)
	}
}

func (f *anyApplicationConfigTest) createTestOfValidateFuncOfTagFilterShouldReturnOneErrorAndNoWarningsWhenValueCannotBeParsed() func(t *testing.T) {
	return func(t *testing.T) {
		resourceSchema := f.resourceHandle.MetaData().Schema
		value := invalidTagFilter

		warns, errs := resourceSchema[ApplicationAlertConfigFieldTagFilter].ValidateFunc(value, ApplicationAlertConfigFieldTagFilter)
		require.Empty(t, warns)
		require.Len(t, errs, 1)
	}
}

func (f *anyApplicationConfigTest) createTetResourceShouldHaveSchemaVersionOne() func(t *testing.T) {
	return func(t *testing.T) {
		require.Equal(t, 1, f.resourceHandle.MetaData().SchemaVersion)
	}
}

func (f *anyApplicationConfigTest) createTetResourceShouldHaveOneStateUpgrader() func(t *testing.T) {
	return func(t *testing.T) {
		require.Len(t, f.resourceHandle.StateUpgraders(), 1)
	}
}

func (f *anyApplicationConfigTest) createTestResourceShouldMigrateFullNameToNameWhenExecutingFirstStateUpgraderAndFullNameIsAvailable() func(t *testing.T) {
	return func(t *testing.T) {
		input := map[string]interface{}{
			"full_name": "test",
		}
		result, err := f.resourceHandle.StateUpgraders()[0].Upgrade(nil, input, nil)

		require.NoError(t, err)
		require.Len(t, result, 1)
		require.NotContains(t, result, ApplicationAlertConfigFieldFullName)
		require.Contains(t, result, ApplicationAlertConfigFieldName)
		require.Equal(t, "test", result[ApplicationAlertConfigFieldName])
	}
}

func (f *anyApplicationConfigTest) createTestResourceShouldDoNothingWhenExecutingFirstStateUpgraderAndFullNameIsNotAvailable() func(t *testing.T) {
	return func(t *testing.T) {
		input := map[string]interface{}{
			"name": "test",
		}
		result, err := f.resourceHandle.StateUpgraders()[0].Upgrade(nil, input, nil)

		require.NoError(t, err)
		require.Equal(t, input, result)
	}
}

func (f *anyApplicationConfigTest) createTetResourceShouldHaveCorrectResourceName() func(t *testing.T) {
	return func(t *testing.T) {
		require.Equal(t, f.resourceHandle.MetaData().ResourceName, f.terraformResourceName)
	}
}

func (f *anyApplicationConfigTest) createTestCasesForUpdatesOfTerraformResourceStateFromModel(t *testing.T) {
	metricName := "test-metric"
	stableHash := int32(1234)
	statusCodeStart := int32(200)
	statusCodeEnd := int32(300)
	logMessage := "test-log-message"
	logLevel := restapi.LogLevelError
	logOperator := restapi.EqualsOperator
	rules := []testPair[restapi.ApplicationAlertRule, []interface{}]{
		{
			name: ApplicationAlertConfigFieldRuleThroughput,
			input: restapi.ApplicationAlertRule{
				AlertType:   ApplicationAlertConfigFieldRuleThroughput,
				Aggregation: restapi.MinAggregation,
				MetricName:  metricName,
				StableHash:  &stableHash,
			},
			expected: []interface{}{
				map[string]interface{}{
					ApplicationAlertConfigFieldRuleErrorRate:  []interface{}{},
					ApplicationAlertConfigFieldRuleLogs:       []interface{}{},
					ApplicationAlertConfigFieldRuleSlowness:   []interface{}{},
					ApplicationAlertConfigFieldRuleStatusCode: []interface{}{},
					ApplicationAlertConfigFieldRuleThroughput: []interface{}{
						map[string]interface{}{
							ApplicationAlertConfigFieldRuleAggregation: string(restapi.MinAggregation),
							ApplicationAlertConfigFieldRuleMetricName:  metricName,
							ApplicationAlertConfigFieldRuleStableHash:  int(stableHash),
						},
					},
				},
			},
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
			expected: []interface{}{
				map[string]interface{}{
					ApplicationAlertConfigFieldRuleErrorRate: []interface{}{},
					ApplicationAlertConfigFieldRuleLogs:      []interface{}{},
					ApplicationAlertConfigFieldRuleSlowness:  []interface{}{},
					ApplicationAlertConfigFieldRuleStatusCode: []interface{}{
						map[string]interface{}{
							ApplicationAlertConfigFieldRuleAggregation:     string(restapi.MinAggregation),
							ApplicationAlertConfigFieldRuleMetricName:      metricName,
							ApplicationAlertConfigFieldRuleStableHash:      int(stableHash),
							ApplicationAlertConfigFieldRuleStatusCodeStart: int(statusCodeStart),
							ApplicationAlertConfigFieldRuleStatusCodeEnd:   int(statusCodeEnd),
						},
					},
					ApplicationAlertConfigFieldRuleThroughput: []interface{}{},
				},
			},
		},
		{
			name: ApplicationAlertConfigFieldRuleSlowness,
			input: restapi.ApplicationAlertRule{
				AlertType:   ApplicationAlertConfigFieldRuleSlowness,
				Aggregation: restapi.MinAggregation,
				MetricName:  metricName,
				StableHash:  &stableHash,
			},
			expected: []interface{}{
				map[string]interface{}{
					ApplicationAlertConfigFieldRuleErrorRate: []interface{}{},
					ApplicationAlertConfigFieldRuleLogs:      []interface{}{},
					ApplicationAlertConfigFieldRuleSlowness: []interface{}{
						map[string]interface{}{
							ApplicationAlertConfigFieldRuleAggregation: string(restapi.MinAggregation),
							ApplicationAlertConfigFieldRuleMetricName:  metricName,
							ApplicationAlertConfigFieldRuleStableHash:  int(stableHash),
						},
					},
					ApplicationAlertConfigFieldRuleStatusCode: []interface{}{},
					ApplicationAlertConfigFieldRuleThroughput: []interface{}{},
				},
			},
		},
		{
			name: ApplicationAlertConfigFieldRuleLogs,
			input: restapi.ApplicationAlertRule{
				AlertType:   ApplicationAlertConfigFieldRuleLogs,
				Aggregation: restapi.MinAggregation,
				MetricName:  metricName,
				StableHash:  &stableHash,
				Message:     &logMessage,
				Operator:    &logOperator,
				Level:       &logLevel,
			},
			expected: []interface{}{
				map[string]interface{}{
					ApplicationAlertConfigFieldRuleErrorRate: []interface{}{},
					ApplicationAlertConfigFieldRuleLogs: []interface{}{
						map[string]interface{}{
							ApplicationAlertConfigFieldRuleAggregation:  string(restapi.MinAggregation),
							ApplicationAlertConfigFieldRuleMetricName:   metricName,
							ApplicationAlertConfigFieldRuleStableHash:   int(stableHash),
							ApplicationAlertConfigFieldRuleLogsLevel:    string(logLevel),
							ApplicationAlertConfigFieldRuleLogsMessage:  logMessage,
							ApplicationAlertConfigFieldRuleLogsOperator: string(logOperator),
						},
					},
					ApplicationAlertConfigFieldRuleSlowness:   []interface{}{},
					ApplicationAlertConfigFieldRuleStatusCode: []interface{}{},
					ApplicationAlertConfigFieldRuleThroughput: []interface{}{},
				},
			},
		},
		{
			name: "ErrorRate",
			input: restapi.ApplicationAlertRule{
				AlertType:   "errorRate",
				Aggregation: restapi.MinAggregation,
				MetricName:  metricName,
				StableHash:  &stableHash,
			},
			expected: []interface{}{
				map[string]interface{}{
					ApplicationAlertConfigFieldRuleErrorRate: []interface{}{
						map[string]interface{}{
							ApplicationAlertConfigFieldRuleAggregation: string(restapi.MinAggregation),
							ApplicationAlertConfigFieldRuleMetricName:  metricName,
							ApplicationAlertConfigFieldRuleStableHash:  int(stableHash),
						},
					},
					ApplicationAlertConfigFieldRuleLogs:       []interface{}{},
					ApplicationAlertConfigFieldRuleSlowness:   []interface{}{},
					ApplicationAlertConfigFieldRuleStatusCode: []interface{}{},
					ApplicationAlertConfigFieldRuleThroughput: []interface{}{},
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
			expected: []interface{}{
				map[string]interface{}{
					ApplicationAlertConfigFieldTimeThresholdRequestImpact: []interface{}{
						map[string]interface{}{
							ApplicationAlertConfigFieldTimeThresholdTimeWindow:            int(timeThresholdWindow),
							ApplicationAlertConfigFieldTimeThresholdRequestImpactRequests: int(timeThresholdRequests),
						},
					},
					ApplicationAlertConfigFieldTimeThresholdViolationsInPeriod:   []interface{}{},
					ApplicationAlertConfigFieldTimeThresholdViolationsInSequence: []interface{}{},
				},
			},
		},
		{
			name: "ViolationsInPeriod",
			input: restapi.TimeThreshold{
				Type:       "violationsInPeriod",
				TimeWindow: timeThresholdWindow,
				Violations: &timeThresholdViolations,
			},
			expected: []interface{}{
				map[string]interface{}{
					ApplicationAlertConfigFieldTimeThresholdRequestImpact: []interface{}{},
					ApplicationAlertConfigFieldTimeThresholdViolationsInPeriod: []interface{}{
						map[string]interface{}{
							ApplicationAlertConfigFieldTimeThresholdTimeWindow:                   int(timeThresholdWindow),
							ApplicationAlertConfigFieldTimeThresholdViolationsInPeriodViolations: int(timeThresholdViolations),
						},
					},
					ApplicationAlertConfigFieldTimeThresholdViolationsInSequence: []interface{}{},
				},
			},
		},
		{
			name: "ViolationsInSequence",
			input: restapi.TimeThreshold{
				Type:       "violationsInSequence",
				TimeWindow: timeThresholdWindow,
			},
			expected: []interface{}{
				map[string]interface{}{
					ApplicationAlertConfigFieldTimeThresholdRequestImpact:      []interface{}{},
					ApplicationAlertConfigFieldTimeThresholdViolationsInPeriod: []interface{}{},
					ApplicationAlertConfigFieldTimeThresholdViolationsInSequence: []interface{}{
						map[string]interface{}{
							ApplicationAlertConfigFieldTimeThresholdTimeWindow: int(timeThresholdWindow),
						},
					},
				},
			},
		},
	}

	for _, rule := range rules {
		for _, threshold := range thresholds {
			for _, timeThreshold := range timeThresholds {
				t.Run(fmt.Sprintf("Should update terraform state of %s from REST response with %s and %s and %s", f.terraformResourceName, rule.name, threshold.name, timeThreshold.name), f.createTestShouldUpdateTerraformResourceStateFromModelCase(rule, threshold, timeThreshold))
			}
		}
	}
}

func (f *anyApplicationConfigTest) createTestShouldUpdateTerraformResourceStateFromModelCase(ruleTestPair testPair[restapi.ApplicationAlertRule, []interface{}],
	thresholdTestPair testPair[restapi.Threshold, []interface{}],
	timeThresholdTestPair testPair[restapi.TimeThreshold, []interface{}]) func(t *testing.T) {
	return func(t *testing.T) {
		name := "application-alert-config-name"
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
			BoundaryScope:    restapi.BoundaryScopeInbound,
			Description:      "application-alert-config-description",
			EvaluationType:   restapi.EvaluationTypePerApplication,
			Granularity:      restapi.Granularity600000,
			IncludeInternal:  false,
			IncludeSynthetic: false,
			CustomerPayloadFields: []restapi.CustomPayloadField[any]{
				{
					Type:  restapi.StaticCustomPayloadType,
					Key:   "static-key",
					Value: restapi.StaticStringCustomPayloadFieldValue("static-value"),
				},
			},
			Name:                name,
			Rule:                ruleTestPair.input,
			Severity:            restapi.SeverityCritical.GetAPIRepresentation(),
			TagFilterExpression: restapi.NewStringTagFilter(restapi.TagFilterEntitySource, "service.name", restapi.EqualsOperator, "test"),
			Threshold:           thresholdTestPair.input,
			TimeThreshold:       timeThresholdTestPair.input,
			Triggering:          true,
		}

		testHelper := NewTestHelper[*restapi.ApplicationAlertConfig](t)
		sut := f.resourceHandle
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

		err := sut.UpdateState(resourceData, &applicationConfig)

		require.NoError(t, err)
		require.Equal(t, applicationAlertConfigID, resourceData.Id())
		require.Equal(t, []interface{}{"channel-2", "channel-1"}, (resourceData.Get(ApplicationAlertConfigFieldAlertChannelIDs).(*schema.Set)).List())
		f.requireApplicationAlertConfigApplicationSetOnSchema(t, resourceData)
		require.Equal(t, string(restapi.BoundaryScopeInbound), resourceData.Get(ApplicationAlertConfigFieldBoundaryScope))
		require.Equal(t, "application-alert-config-description", resourceData.Get(ApplicationAlertConfigFieldDescription))
		require.Equal(t, name, resourceData.Get(ApplicationAlertConfigFieldName))
		require.Equal(t, string(restapi.EvaluationTypePerApplication), resourceData.Get(ApplicationAlertConfigFieldEvaluationType))
		require.False(t, resourceData.Get(ApplicationAlertConfigFieldIncludeInternal).(bool))
		require.False(t, resourceData.Get(ApplicationAlertConfigFieldIncludeSynthetic).(bool))
		require.Equal(t, []interface{}{
			map[string]interface{}{ApplicationAlertConfigFieldCustomPayloadFieldsKey: "static-key", ApplicationAlertConfigFieldCustomPayloadFieldsValue: "static-value"},
		}, resourceData.Get(ApplicationAlertConfigFieldCustomPayloadFields).(*schema.Set).List())
		require.Equal(t, ruleTestPair.expected, resourceData.Get(ApplicationAlertConfigFieldRule))
		require.Equal(t, restapi.SeverityCritical.GetTerraformRepresentation(), resourceData.Get(ApplicationAlertConfigFieldSeverity))
		require.Equal(t, "service.name@src EQUALS 'test'", resourceData.Get(ApplicationAlertConfigFieldTagFilter))
		f.requireApplicationAlertConfigThresholdSetOnSchema(t, thresholdTestPair.expected, resourceData)
		require.Equal(t, timeThresholdTestPair.expected, resourceData.Get(ApplicationAlertConfigFieldTimeThreshold))
		require.True(t, resourceData.Get(ApplicationAlertConfigFieldTriggering).(bool))
	}
}

func (f *anyApplicationConfigTest) requireApplicationAlertConfigApplicationSetOnSchema(t *testing.T, resourceData *schema.ResourceData) {
	actualValues := resourceData.Get(ApplicationAlertConfigFieldApplications).(*schema.Set)
	require.Equal(t, 1, actualValues.Len())
	app := actualValues.List()[0].(map[string]interface{})
	require.Equal(t, "app-1", app[ApplicationAlertConfigFieldApplicationsApplicationID])
	require.True(t, app[ApplicationAlertConfigFieldApplicationsInclusive].(bool))
	services := app[ApplicationAlertConfigFieldApplicationsServices].(*schema.Set)
	require.Equal(t, 1, services.Len())
	service := services.List()[0].(map[string]interface{})
	require.Equal(t, "srv-1", service[ApplicationAlertConfigFieldApplicationsServicesServiceID])
	require.True(t, service[ApplicationAlertConfigFieldApplicationsInclusive].(bool))
	endpoints := service[ApplicationAlertConfigFieldApplicationsServicesEndpoints].(*schema.Set)
	require.Equal(t, 1, endpoints.Len())
	endpoint := endpoints.List()[0].(map[string]interface{})
	require.Equal(t, "edp-1", endpoint[ApplicationAlertConfigFieldApplicationsServicesEndpointsEndpointID])
	require.True(t, endpoint[ApplicationAlertConfigFieldApplicationsInclusive].(bool))
}

func (f *anyApplicationConfigTest) requireApplicationAlertConfigThresholdSetOnSchema(t *testing.T, expected []interface{}, resourceData *schema.ResourceData) {
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

func (f *anyApplicationConfigTest) createTestCasesForMappingOfTerraformResourceStateToModel(t *testing.T) {
	metricName := "test-metric"
	stableHash := int32(1234)
	statusCodeStart := int32(200)
	statusCodeEnd := int32(300)
	logMessage := "test-log-message"
	logLevel := restapi.LogLevelError
	logOperator := restapi.EqualsOperator
	rules := []testPair[[]map[string]interface{}, restapi.ApplicationAlertRule]{
		{
			name: ApplicationAlertConfigFieldRuleThroughput,
			expected: restapi.ApplicationAlertRule{
				AlertType:   ApplicationAlertConfigFieldRuleThroughput,
				Aggregation: restapi.MinAggregation,
				MetricName:  metricName,
				StableHash:  &stableHash,
			},
			input: []map[string]interface{}{
				{
					ApplicationAlertConfigFieldRuleErrorRate:  []interface{}{},
					ApplicationAlertConfigFieldRuleLogs:       []interface{}{},
					ApplicationAlertConfigFieldRuleSlowness:   []interface{}{},
					ApplicationAlertConfigFieldRuleStatusCode: []interface{}{},
					ApplicationAlertConfigFieldRuleThroughput: []interface{}{
						map[string]interface{}{
							ApplicationAlertConfigFieldRuleAggregation: string(restapi.MinAggregation),
							ApplicationAlertConfigFieldRuleMetricName:  metricName,
							ApplicationAlertConfigFieldRuleStableHash:  int(stableHash),
						},
					},
				},
			},
		},
		{
			name: "StatusCode",
			expected: restapi.ApplicationAlertRule{
				AlertType:       "statusCode",
				Aggregation:     restapi.MinAggregation,
				MetricName:      metricName,
				StableHash:      &stableHash,
				StatusCodeStart: &statusCodeStart,
				StatusCodeEnd:   &statusCodeEnd,
			},
			input: []map[string]interface{}{
				{
					ApplicationAlertConfigFieldRuleErrorRate: []interface{}{},
					ApplicationAlertConfigFieldRuleLogs:      []interface{}{},
					ApplicationAlertConfigFieldRuleSlowness:  []interface{}{},
					ApplicationAlertConfigFieldRuleStatusCode: []interface{}{
						map[string]interface{}{
							ApplicationAlertConfigFieldRuleAggregation:     string(restapi.MinAggregation),
							ApplicationAlertConfigFieldRuleMetricName:      metricName,
							ApplicationAlertConfigFieldRuleStableHash:      int(stableHash),
							ApplicationAlertConfigFieldRuleStatusCodeStart: int(statusCodeStart),
							ApplicationAlertConfigFieldRuleStatusCodeEnd:   int(statusCodeEnd),
						},
					},
					ApplicationAlertConfigFieldRuleThroughput: []interface{}{},
				},
			},
		},
		{
			name: ApplicationAlertConfigFieldRuleSlowness,
			expected: restapi.ApplicationAlertRule{
				AlertType:   ApplicationAlertConfigFieldRuleSlowness,
				Aggregation: restapi.MinAggregation,
				MetricName:  metricName,
				StableHash:  &stableHash,
			},
			input: []map[string]interface{}{
				{
					ApplicationAlertConfigFieldRuleErrorRate: []interface{}{},
					ApplicationAlertConfigFieldRuleLogs:      []interface{}{},
					ApplicationAlertConfigFieldRuleSlowness: []interface{}{
						map[string]interface{}{
							ApplicationAlertConfigFieldRuleAggregation: string(restapi.MinAggregation),
							ApplicationAlertConfigFieldRuleMetricName:  metricName,
							ApplicationAlertConfigFieldRuleStableHash:  int(stableHash),
						},
					},
					ApplicationAlertConfigFieldRuleStatusCode: []interface{}{},
					ApplicationAlertConfigFieldRuleThroughput: []interface{}{},
				},
			},
		},
		{
			name: ApplicationAlertConfigFieldRuleLogs,
			expected: restapi.ApplicationAlertRule{
				AlertType:   ApplicationAlertConfigFieldRuleLogs,
				Aggregation: restapi.MinAggregation,
				MetricName:  metricName,
				StableHash:  &stableHash,
				Message:     &logMessage,
				Operator:    &logOperator,
				Level:       &logLevel,
			},
			input: []map[string]interface{}{
				{
					ApplicationAlertConfigFieldRuleErrorRate: []interface{}{},
					ApplicationAlertConfigFieldRuleLogs: []interface{}{
						map[string]interface{}{
							ApplicationAlertConfigFieldRuleAggregation:  string(restapi.MinAggregation),
							ApplicationAlertConfigFieldRuleMetricName:   metricName,
							ApplicationAlertConfigFieldRuleStableHash:   int(stableHash),
							ApplicationAlertConfigFieldRuleLogsLevel:    string(logLevel),
							ApplicationAlertConfigFieldRuleLogsMessage:  logMessage,
							ApplicationAlertConfigFieldRuleLogsOperator: string(logOperator),
						},
					},
					ApplicationAlertConfigFieldRuleSlowness:   []interface{}{},
					ApplicationAlertConfigFieldRuleStatusCode: []interface{}{},
					ApplicationAlertConfigFieldRuleThroughput: []interface{}{},
				},
			},
		},
		{
			name: "ErrorRate",
			expected: restapi.ApplicationAlertRule{
				AlertType:   "errorRate",
				Aggregation: restapi.MinAggregation,
				MetricName:  metricName,
				StableHash:  &stableHash,
			},
			input: []map[string]interface{}{
				{
					ApplicationAlertConfigFieldRuleErrorRate: []interface{}{
						map[string]interface{}{
							ApplicationAlertConfigFieldRuleAggregation: string(restapi.MinAggregation),
							ApplicationAlertConfigFieldRuleMetricName:  metricName,
							ApplicationAlertConfigFieldRuleStableHash:  int(stableHash),
						},
					},
					ApplicationAlertConfigFieldRuleLogs:       []interface{}{},
					ApplicationAlertConfigFieldRuleSlowness:   []interface{}{},
					ApplicationAlertConfigFieldRuleStatusCode: []interface{}{},
					ApplicationAlertConfigFieldRuleThroughput: []interface{}{},
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
	timeThresholdRequests := int32(5)
	timeThresholdViolations := int32(3)
	timeThresholds := []testPair[[]map[string]interface{}, restapi.TimeThreshold]{
		{
			name: "RequestImpact",
			expected: restapi.TimeThreshold{
				Type:       "requestImpact",
				TimeWindow: timeThresholdWindow,
				Requests:   &timeThresholdRequests,
			},
			input: []map[string]interface{}{
				{
					ApplicationAlertConfigFieldTimeThresholdRequestImpact: []interface{}{
						map[string]interface{}{
							ApplicationAlertConfigFieldTimeThresholdTimeWindow:            int(timeThresholdWindow),
							ApplicationAlertConfigFieldTimeThresholdRequestImpactRequests: int(timeThresholdRequests),
						},
					},
					ApplicationAlertConfigFieldTimeThresholdViolationsInPeriod:   []interface{}{},
					ApplicationAlertConfigFieldTimeThresholdViolationsInSequence: []interface{}{},
				},
			},
		},
		{
			name: "ViolationsInPeriod",
			expected: restapi.TimeThreshold{
				Type:       "violationsInPeriod",
				TimeWindow: timeThresholdWindow,
				Violations: &timeThresholdViolations,
			},
			input: []map[string]interface{}{
				{
					ApplicationAlertConfigFieldTimeThresholdRequestImpact: []interface{}{},
					ApplicationAlertConfigFieldTimeThresholdViolationsInPeriod: []interface{}{
						map[string]interface{}{
							ApplicationAlertConfigFieldTimeThresholdTimeWindow:                   int(timeThresholdWindow),
							ApplicationAlertConfigFieldTimeThresholdViolationsInPeriodViolations: int(timeThresholdViolations),
						},
					},
					ApplicationAlertConfigFieldTimeThresholdViolationsInSequence: []interface{}{},
				},
			},
		},
		{
			name: "ViolationsInSequence",
			expected: restapi.TimeThreshold{
				Type:       "violationsInSequence",
				TimeWindow: timeThresholdWindow,
			},
			input: []map[string]interface{}{
				{
					ApplicationAlertConfigFieldTimeThresholdRequestImpact:      []interface{}{},
					ApplicationAlertConfigFieldTimeThresholdViolationsInPeriod: []interface{}{},
					ApplicationAlertConfigFieldTimeThresholdViolationsInSequence: []interface{}{
						map[string]interface{}{
							ApplicationAlertConfigFieldTimeThresholdTimeWindow: int(timeThresholdWindow),
						},
					},
				},
			},
		},
	}

	for _, rule := range rules {
		for _, threshold := range thresholds {
			for _, timeThreshold := range timeThresholds {
				t.Run(fmt.Sprintf("Should map terraform state of %s to REST model with %s and %s and %s", f.terraformResourceName, rule.name, threshold.name, timeThreshold.name), f.createTestShouldMapTerraformResourceStateToModelCase(rule, threshold, timeThreshold))
			}
		}
	}
}

func (f *anyApplicationConfigTest) createTestShouldMapTerraformResourceStateToModelCase(ruleTestPair testPair[[]map[string]interface{}, restapi.ApplicationAlertRule],
	thresholdTestPair testPair[[]map[string]interface{}, restapi.Threshold],
	timeThresholdTestPair testPair[[]map[string]interface{}, restapi.TimeThreshold]) func(t *testing.T) {
	return func(t *testing.T) {
		applicationAlertConfigID := "application-alert-config-id"
		name := "application-alert-config-name"
		expectedApplicationConfig := restapi.ApplicationAlertConfig{
			ID:              applicationAlertConfigID,
			AlertChannelIDs: []string{"channel-2", "channel-1"},
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
			CustomerPayloadFields: []restapi.CustomPayloadField[any]{
				{
					Type:  restapi.StaticCustomPayloadType,
					Key:   "static-key",
					Value: restapi.StaticStringCustomPayloadFieldValue("static-value"),
				},
			},
			Name:                name,
			Rule:                ruleTestPair.expected,
			Severity:            restapi.SeverityCritical.GetAPIRepresentation(),
			TagFilterExpression: restapi.NewStringTagFilter(restapi.TagFilterEntitySource, "service.name", restapi.EqualsOperator, "test"),
			Threshold:           thresholdTestPair.expected,
			TimeThreshold:       timeThresholdTestPair.expected,
			Triggering:          true,
		}

		testHelper := NewTestHelper[*restapi.ApplicationAlertConfig](t)
		sut := f.resourceHandle
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)
		setValueOnResourceData(t, resourceData, ApplicationAlertConfigFieldAlertChannelIDs, []interface{}{"channel-2", "channel-1"})
		setValueOnResourceData(t, resourceData, ApplicationAlertConfigFieldApplications, []interface{}{
			map[string]interface{}{
				ApplicationAlertConfigFieldApplicationsApplicationID: "app-1",
				ApplicationAlertConfigFieldApplicationsInclusive:     true,
				ApplicationAlertConfigFieldApplicationsServices: []interface{}{
					map[string]interface{}{
						ApplicationAlertConfigFieldApplicationsServicesServiceID: "srv-1",
						ApplicationAlertConfigFieldApplicationsInclusive:         true,
						ApplicationAlertConfigFieldApplicationsServicesEndpoints: []interface{}{
							map[string]interface{}{
								ApplicationAlertConfigFieldApplicationsServicesEndpointsEndpointID: "edp-1",
								ApplicationAlertConfigFieldApplicationsInclusive:                   true,
							},
						},
					},
				},
			},
		})
		setValueOnResourceData(t, resourceData, ApplicationAlertConfigFieldBoundaryScope, restapi.BoundaryScopeInbound)
		setValueOnResourceData(t, resourceData, ApplicationAlertConfigFieldCustomPayloadFields, []interface{}{
			map[string]interface{}{ApplicationAlertConfigFieldCustomPayloadFieldsKey: "static-key", ApplicationAlertConfigFieldCustomPayloadFieldsValue: "static-value"},
		})
		setValueOnResourceData(t, resourceData, ApplicationAlertConfigFieldDescription, "application-alert-config-description")
		setValueOnResourceData(t, resourceData, ApplicationAlertConfigFieldEvaluationType, restapi.EvaluationTypePerApplication)
		setValueOnResourceData(t, resourceData, ApplicationAlertConfigFieldGranularity, restapi.Granularity600000)
		setValueOnResourceData(t, resourceData, ApplicationAlertConfigFieldIncludeInternal, false)
		setValueOnResourceData(t, resourceData, ApplicationAlertConfigFieldIncludeSynthetic, false)
		setValueOnResourceData(t, resourceData, ApplicationAlertConfigFieldName, name)
		setValueOnResourceData(t, resourceData, ApplicationAlertConfigFieldRule, ruleTestPair.input)
		setValueOnResourceData(t, resourceData, ApplicationAlertConfigFieldSeverity, restapi.SeverityCritical.GetTerraformRepresentation())
		setValueOnResourceData(t, resourceData, ApplicationAlertConfigFieldTagFilter, "service.name@src EQUALS 'test'")
		setValueOnResourceData(t, resourceData, ResourceFieldThreshold, thresholdTestPair.input)
		setValueOnResourceData(t, resourceData, ApplicationAlertConfigFieldTimeThreshold, timeThresholdTestPair.input)
		setValueOnResourceData(t, resourceData, ApplicationAlertConfigFieldTriggering, true)
		resourceData.SetId(applicationAlertConfigID)

		result, err := sut.MapStateToDataObject(resourceData)

		require.NoError(t, err)
		require.Equal(t, &expectedApplicationConfig, result)
	}
}
