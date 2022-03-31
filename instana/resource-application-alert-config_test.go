package instana_test

import (
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
)

/*
func TestCRUDOfApplicationAlertConfig(t *testing.T) {
	id := RandomID()
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPost, apiTokenApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, apiTokenApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, apiTokenApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		modCount := httpServer.GetCallCount(http.MethodPut, restapi.APITokensResourcePath+"/"+internalID)
		json := fmt.Sprintf(`
		{
			"id" : "%s",
			"accessGrantingToken": "%s",
			"internalId" : "%s",
			"name" : "name %d",
			"canConfigureServiceMapping" : true,
			"canConfigureEumApplications" : true,
			"canConfigureMobileAppMonitoring" : true,
			"canConfigureUsers" : true,
			"canInstallNewAgents" : true,
			"canSeeUsageInformation" : true,
			"canConfigureIntegrations" : true,
			"canSeeOnPremLicenseInformation" : true,
			"canConfigureCustomAlerts" : true,
			"canConfigureApiTokens" : true,
			"canConfigureAgentRunMode" : true,
			"canViewAuditLog" : true,
			"canConfigureAgents" : true,
			"canConfigureAuthenticationMethods" : true,
			"canConfigureApplications" : true,
			"canConfigureTeams" : true,
			"canConfigureReleases" : true,
			"canConfigureLogManagement" : true,
			"canCreatePublicCustomDashboards" : true,
			"canViewLogs" : true,
			"canViewTraceDetails" : true,
			"canConfigureSessionSettings" : true,
			"canConfigureServiceLevelIndicators" : true,
			"canConfigureGlobalAlertPayload" : true,
			"canConfigureGlobalAlertConfigs" : true,
			"canViewAccountAndBillingInformation" : true,
			"canEditAllAccessibleCustomDashboards" : true
		}
		`, id, accessGrantingToken, vars["internal-id"], modCount)
		w.Header().Set(contentType, r.Header.Get(contentType))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			createAPITokenConfigResourceTestStep(httpServer.GetPort(), 0, id, accessGrantingToken, internalID),
			testStepImportWithCustomID(testAPITokenDefinition, internalID),
			createAPITokenConfigResourceTestStep(httpServer.GetPort(), 1, id, accessGrantingToken, internalID),
			testStepImportWithCustomID(testAPITokenDefinition, internalID),
		},
	})
}

func createApplicationAlertConfigResourceTestStep(httpPort int, iteration int, id string, accessGrantingToken string, internalID string) resource.TestStep {
	return resource.TestStep{
		Config: fmt.Sprintf(resourceAPITokenDefinitionTemplate, httpPort, iteration),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(testAPITokenDefinition, "id", id),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldAccessGrantingToken, accessGrantingToken),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldInternalID, internalID),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldName, formatResourceName(iteration)),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldFullName, formatResourceFullName(iteration)),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureServiceMapping, valueTrue),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureEumApplications, valueTrue),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureMobileAppMonitoring, valueTrue),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureUsers, valueTrue),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanInstallNewAgents, valueTrue),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanSeeUsageInformation, valueTrue),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureIntegrations, valueTrue),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanSeeOnPremiseLicenseInformation, valueTrue),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureCustomAlerts, valueTrue),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureAPITokens, valueTrue),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureAgentRunMode, valueTrue),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanViewAuditLog, valueTrue),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureAgents, valueTrue),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureAuthenticationMethods, valueTrue),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureApplications, valueTrue),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureTeams, valueTrue),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureReleases, valueTrue),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureLogManagement, valueTrue),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanCreatePublicCustomDashboards, valueTrue),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanViewLogs, valueTrue),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanViewTraceDetails, valueTrue),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureSessionSettings, valueTrue),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureServiceLevelIndicators, valueTrue),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureGlobalAlertPayload, valueTrue),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanConfigureGlobalAlertConfigs, valueTrue),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanViewAccountAndBillingInformation, valueTrue),
			resource.TestCheckResourceAttr(testAPITokenDefinition, APITokenFieldCanEditAllAccessibleCustomDashboards, valueTrue),
		),
	}
}
*/

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
							ApplicationAlertConfigFieldRuleAggregation:   string(restapi.MinAggregation),
							ApplicationAlertConfigFieldRuleMetricName:    metricName,
							ApplicationAlertConfigFieldRuleStableHash:    int(stableHash),
							ApplicationAlertConfigFieldRuleLogsLevel:     string(logLevel),
							ApplicationAlertConfigFieldRuleLogsMessage:   logMessage,
							ApplicationAlertConfigFieldThresholdOperator: string(logOperator),
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
					ApplicationAlertConfigFieldThresholdHistoricBaseline: []interface{}{},
					ApplicationAlertConfigFieldThresholdStatic: []interface{}{
						map[string]interface{}{
							ApplicationAlertConfigFieldThresholdLastUpdated: int(thresholdLastUpdate),
							ApplicationAlertConfigFieldThresholdOperator:    string(restapi.ThresholdOperatorGreaterThan),
							ApplicationAlertConfigFieldThresholdStaticValue: thresholdValue,
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
					ApplicationAlertConfigFieldThresholdHistoricBaseline: []interface{}{
						map[string]interface{}{
							ApplicationAlertConfigFieldThresholdLastUpdated:                     int(thresholdLastUpdate),
							ApplicationAlertConfigFieldThresholdOperator:                        string(restapi.ThresholdOperatorGreaterThan),
							ApplicationAlertConfigFieldThresholdHistoricBaselineSeasonality:     string(thresholdSeasonality),
							ApplicationAlertConfigFieldThresholdHistoricBaselineBaseline:        thresholdBaseline,
							ApplicationAlertConfigFieldThresholdHistoricBaselineDeviationFactor: float64(thresholdDeviationFactor),
						},
					},
					ApplicationAlertConfigFieldThresholdStatic: []interface{}{},
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
			BoundaryScope:    restapi.BoundaryScopeInbound,
			Description:      "application-alert-config-description",
			EvaluationType:   restapi.EvaluationTypePerApplication,
			Granularity:      restapi.Granularity600000,
			IncludeInternal:  false,
			IncludeSynthetic: false,
			CustomerPayloadFields: []restapi.CustomPayloadField{
				{
					Type:  restapi.StaticCustomPayloadType,
					Key:   "static-key",
					Value: "static-value",
				},
				{
					Type:  restapi.DynamicCustomPayloadType,
					Key:   "dynamic-key",
					Value: "dynamic-value",
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

		testHelper := NewTestHelper(t)
		sut := NewApplicationAlertConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

		err := sut.UpdateState(resourceData, &applicationConfig, testHelper.ResourceFormatter())

		require.NoError(t, err)
		require.Equal(t, applicationAlertConfigID, resourceData.Id())
		require.Equal(t, []interface{}{"channel-2", "channel-1"}, (resourceData.Get(ApplicationAlertConfigFieldAlertChannelIDs).(*schema.Set)).List())
		require.Equal(t, []interface{}{
			map[string]interface{}{
				"application_id": "app-1",
				"inclusive":      true,
				"services": []interface{}{
					map[string]interface{}{
						"service_id": "srv-1",
						"inclusive":  true,
						"endpoints": []interface{}{
							map[string]interface{}{
								"endpoint_id": "edp-1",
								"inclusive":   true,
							},
						},
					},
				},
			},
		}, resourceData.Get(ApplicationAlertConfigFieldApplications))
		require.Equal(t, string(restapi.BoundaryScopeInbound), resourceData.Get(ApplicationAlertConfigFieldBoundaryScope))
		require.Equal(t, "application-alert-config-description", resourceData.Get(ApplicationAlertConfigFieldDescription))
		require.Equal(t, "application-alert-config-name", resourceData.Get(ApplicationAlertConfigFieldName))
		require.Equal(t, fullName, resourceData.Get(ApplicationAlertConfigFieldFullName))
		require.Equal(t, string(restapi.EvaluationTypePerApplication), resourceData.Get(ApplicationAlertConfigFieldEvaluationType))
		require.False(t, resourceData.Get(ApplicationAlertConfigFieldIncludeInternal).(bool))
		require.False(t, resourceData.Get(ApplicationAlertConfigFieldIncludeSynthetic).(bool))
		require.Equal(t, []interface{}{
			map[string]interface{}{ApplicationAlertConfigFieldCustomPayloadFieldsType: ApplicationAlertConfigFieldThresholdStatic, ApplicationAlertConfigFieldCustomPayloadFieldsKey: "static-key", ApplicationAlertConfigFieldCustomPayloadFieldsValue: "static-value"},
			map[string]interface{}{ApplicationAlertConfigFieldCustomPayloadFieldsType: "dynamic", ApplicationAlertConfigFieldCustomPayloadFieldsKey: "dynamic-key", ApplicationAlertConfigFieldCustomPayloadFieldsValue: "dynamic-value"},
		}, resourceData.Get(ApplicationAlertConfigFieldCustomPayloadFields))
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
		if k == ApplicationAlertConfigFieldThresholdHistoricBaseline && len(e.([]interface{})) > 0 {
			expectedHistoricBaselineSlice := e.([]interface{})
			actualHistoricBaselineSlice := actualEntries[k].([]interface{})
			require.Equal(t, 1, len(expectedHistoricBaselineSlice))
			require.Equal(t, len(expected), len(actual))
			expectedHistoricBaseline := expectedHistoricBaselineSlice[0].(map[string]interface{})
			actualHistoricBaseline := actualHistoricBaselineSlice[0].(map[string]interface{})
			for key, expectedBaselineValue := range expectedHistoricBaseline {
				if key == ApplicationAlertConfigFieldThresholdHistoricBaselineBaseline {
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

func TestShouldMapApplicationConfigTerraformResourceStateToModel(t *testing.T) {
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
							ApplicationAlertConfigFieldRuleAggregation:   string(restapi.MinAggregation),
							ApplicationAlertConfigFieldRuleMetricName:    metricName,
							ApplicationAlertConfigFieldRuleStableHash:    int(stableHash),
							ApplicationAlertConfigFieldRuleLogsLevel:     string(logLevel),
							ApplicationAlertConfigFieldRuleLogsMessage:   logMessage,
							ApplicationAlertConfigFieldThresholdOperator: string(logOperator),
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
					ApplicationAlertConfigFieldThresholdHistoricBaseline: []interface{}{},
					ApplicationAlertConfigFieldThresholdStatic: []interface{}{
						map[string]interface{}{
							ApplicationAlertConfigFieldThresholdLastUpdated: int(thresholdLastUpdate),
							ApplicationAlertConfigFieldThresholdOperator:    string(restapi.ThresholdOperatorGreaterThan),
							ApplicationAlertConfigFieldThresholdStaticValue: thresholdValue,
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
					ApplicationAlertConfigFieldThresholdHistoricBaseline: []interface{}{
						map[string]interface{}{
							ApplicationAlertConfigFieldThresholdLastUpdated:                     int(thresholdLastUpdate),
							ApplicationAlertConfigFieldThresholdOperator:                        string(restapi.ThresholdOperatorGreaterThan),
							ApplicationAlertConfigFieldThresholdHistoricBaselineSeasonality:     string(thresholdSeasonality),
							ApplicationAlertConfigFieldThresholdHistoricBaselineBaseline:        thresholdBaseline,
							ApplicationAlertConfigFieldThresholdHistoricBaselineDeviationFactor: float64(thresholdDeviationFactor),
						},
					},
					ApplicationAlertConfigFieldThresholdStatic: []interface{}{},
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
				t.Run(fmt.Sprintf("TestShouldMapApplicationConfigTerraformResourceStateToModelWith%sAnd%sAnd%s", rule.name, threshold.name, timeThreshold.name), createTestShouldMapApplicationConfigTerraformResourceStateToModelCase(rule, threshold, timeThreshold))
			}
		}
	}
}

func createTestShouldMapApplicationConfigTerraformResourceStateToModelCase(ruleTestPair testPair[[]map[string]interface{}, restapi.ApplicationAlertRule],
	thresholdTestPair testPair[[]map[string]interface{}, restapi.Threshold],
	timeThresholdTestPair testPair[[]map[string]interface{}, restapi.TimeThreshold]) func(t *testing.T) {
	return func(t *testing.T) {
		fullName := "prefix application-alert-config-name suffix"
		applicationAlertConfigID := "application-alert-config-id"

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
			CustomerPayloadFields: []restapi.CustomPayloadField{
				{
					Type:  restapi.StaticCustomPayloadType,
					Key:   "static-key",
					Value: "static-value",
				},
				{
					Type:  restapi.DynamicCustomPayloadType,
					Key:   "dynamic-key",
					Value: "dynamic-value",
				},
			},
			Name:                fullName,
			Rule:                ruleTestPair.expected,
			Severity:            restapi.SeverityCritical.GetAPIRepresentation(),
			TagFilterExpression: restapi.NewStringTagFilter(restapi.TagFilterEntitySource, "service.name", restapi.EqualsOperator, "test"),
			Threshold:           thresholdTestPair.expected,
			TimeThreshold:       timeThresholdTestPair.expected,
			Triggering:          true,
		}

		testHelper := NewTestHelper(t)
		sut := NewApplicationAlertConfigResourceHandle()
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)
		resourceData.Set(ApplicationAlertConfigFieldAlertChannelIDs, []interface{}{"channel-2", "channel-1"})
		resourceData.Set(ApplicationAlertConfigFieldApplications, []interface{}{
			map[string]interface{}{
				"application_id": "app-1",
				"inclusive":      true,
				"services": []interface{}{
					map[string]interface{}{
						"service_id": "srv-1",
						"inclusive":  true,
						"endpoints": []interface{}{
							map[string]interface{}{
								"endpoint_id": "edp-1",
								"inclusive":   true,
							},
						},
					},
				},
			},
		})
		resourceData.Set(ApplicationAlertConfigFieldBoundaryScope, restapi.BoundaryScopeInbound)
		resourceData.Set(ApplicationAlertConfigFieldCustomPayloadFields, []interface{}{
			map[string]interface{}{ApplicationAlertConfigFieldCustomPayloadFieldsType: ApplicationAlertConfigFieldThresholdStatic, ApplicationAlertConfigFieldCustomPayloadFieldsKey: "static-key", ApplicationAlertConfigFieldCustomPayloadFieldsValue: "static-value"},
			map[string]interface{}{ApplicationAlertConfigFieldCustomPayloadFieldsType: "dynamic", ApplicationAlertConfigFieldCustomPayloadFieldsKey: "dynamic-key", ApplicationAlertConfigFieldCustomPayloadFieldsValue: "dynamic-value"},
		})
		resourceData.Set(ApplicationAlertConfigFieldDescription, "application-alert-config-description")
		resourceData.Set(ApplicationAlertConfigFieldEvaluationType, restapi.EvaluationTypePerApplication)
		resourceData.Set(ApplicationAlertConfigFieldGranularity, restapi.Granularity600000)
		resourceData.Set(ApplicationAlertConfigFieldIncludeInternal, false)
		resourceData.Set(ApplicationAlertConfigFieldIncludeSynthetic, false)
		resourceData.Set(ApplicationAlertConfigFieldName, "application-alert-config-name")
		resourceData.Set(ApplicationAlertConfigFieldFullName, fullName)
		resourceData.Set(ApplicationAlertConfigFieldRule, ruleTestPair.input)
		resourceData.Set(ApplicationAlertConfigFieldSeverity, restapi.SeverityCritical.GetTerraformRepresentation())
		resourceData.Set(ApplicationAlertConfigFieldTagFilter, "service.name@src EQUALS 'test'")
		resourceData.Set(ApplicationAlertConfigFieldThreshold, thresholdTestPair.input)
		resourceData.Set(ApplicationAlertConfigFieldTimeThreshold, timeThresholdTestPair.input)
		resourceData.Set(ApplicationAlertConfigFieldTriggering, true)
		resourceData.SetId(applicationAlertConfigID)

		result, err := sut.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

		require.NoError(t, err)
		require.Equal(t, &expectedApplicationConfig, result)
	}
}
