package restapi_test

import (
	"encoding/json"
	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestShouldSuccessfullyUnmarshalApplicationAlertConfig(t *testing.T) {
	thresholdValue := 123.3
	thresholdLastUpdate := int64(12345)
	dynamicValueKey := "dynamic-value-key"
	applicationAlertConfig := ApplicationAlertConfig{
		ID:              "test-application-alert-config-id",
		AlertChannelIDs: []string{"channel-2", "channel-1"},
		Applications: map[string]IncludedApplication{
			"app-1": {
				ApplicationID: "app-1",
				Inclusive:     true,
				Services: map[string]IncludedService{
					"srv-1": {
						ServiceID: "srv-1",
						Inclusive: true,
						Endpoints: map[string]IncludedEndpoint{
							"edp-1": {
								EndpointID: "edp-1",
								Inclusive:  true,
							},
						},
					},
				},
			},
		},
		BoundaryScope:    BoundaryScopeInbound,
		Description:      "application-alert-config-description",
		EvaluationType:   EvaluationTypePerApplication,
		Granularity:      Granularity600000,
		IncludeInternal:  false,
		IncludeSynthetic: false,
		CustomerPayloadFields: []CustomPayloadField[any]{
			{
				Type:  StaticCustomPayloadType,
				Key:   "static-key",
				Value: StaticStringCustomPayloadFieldValue("static-value"),
			},
			{
				Type: DynamicCustomPayloadType,
				Key:  "dynamic-key",
				Value: DynamicCustomPayloadFieldValue{
					TagName: "dynamic-value-tag",
					Key:     &dynamicValueKey,
				},
			},
		},
		Name: "full-name",
		Rule: ApplicationAlertRule{
			AlertType:   "errorRate",
			Aggregation: MinAggregation,
			MetricName:  "metric-name",
		},
		Severity:            SeverityCritical.GetAPIRepresentation(),
		TagFilterExpression: NewStringTagFilter(TagFilterEntitySource, "service.name", EqualsOperator, "test"),
		Threshold: Threshold{
			Type:        "staticThreshold",
			Operator:    ThresholdOperatorGreaterThan,
			LastUpdated: &thresholdLastUpdate,
			Value:       &thresholdValue,
		},
		TimeThreshold: TimeThreshold{
			Type:       "violationsInSequence",
			TimeWindow: 1234,
		},
		Triggering: true,
	}

	serializedJSON, _ := json.Marshal(applicationAlertConfig)

	result, err := NewApplicationAlertConfigUnmarshaller().Unmarshal(serializedJSON)

	require.NoError(t, err)
	require.Equal(t, &applicationAlertConfig, result)
}

func TestShouldFailToUnmarshalApplicationAlertConfigWhenResponseIsAJsonArray(t *testing.T) {
	response := `["foo","bar"]`

	_, err := NewApplicationAlertConfigUnmarshaller().Unmarshal([]byte(response))

	require.Error(t, err)
}

func TestShouldReturnEmptyApplicationAlertConfigWhenNoFieldOfResponseMatchesToModel(t *testing.T) {
	response := `{"foo" : "bar"}`
	config, err := NewApplicationAlertConfigUnmarshaller().Unmarshal([]byte(response))

	require.NoError(t, err)
	require.Equal(t, &ApplicationAlertConfig{}, config)
}

func TestShouldFailToUnmarshalApplicationAlertConfigWhenResponseIsNotAValidJson(t *testing.T) {
	response := `Invalid Data`

	_, err := NewApplicationAlertConfigUnmarshaller().Unmarshal([]byte(response))

	require.Error(t, err)
}
