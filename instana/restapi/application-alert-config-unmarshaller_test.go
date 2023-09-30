package restapi_test

import (
	"encoding/json"
	"fmt"
	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestApplicationAlertConfigUnmarshaller(t *testing.T) {
	ut := &applicationAlertConfigUnmarshallerTest{}
	t.Run("should successfully unmarshal application alert config", ut.shouldSuccessfullyUnmarshalApplicationAlertConfig)
	t.Run("should fail to unmarshal application alert config when response is a json array", ut.shouldFailToUnmarshalApplicationAlertConfigWhenResponseIsAJsonArray)
	t.Run("should fail to unmarshal application alert config when no field of response matches model", ut.shouldReturnEmptyApplicationAlertConfigWhenNoFieldOfResponseMatchesToModel)
	t.Run("should fail to unmarshal application alert config when no response is not a valid json", ut.shouldFailToUnmarshalApplicationAlertConfigWhenResponseIsNotAValidJson)
	t.Run("should successfully unmarshal application alert config array", ut.shouldSuccessfullyUnmarshalApplicationAlertConfigArray)
	t.Run("should fail to unmarshal application alert config array containing at least on invalid application alert config", ut.shouldFailToUnmarshalApplicationAlertConfigArrayContainingAtLeastOneInvalidApplicationAlertConfig)
	t.Run("should fail to unmarshal application alert config array when no valid json is provided", ut.shouldFailToUnmarshalApplicationAlertConfigArrayWhenNoValidJsonIsProvided)
}

type applicationAlertConfigUnmarshallerTest struct {
}

func (ut *applicationAlertConfigUnmarshallerTest) shouldSuccessfullyUnmarshalApplicationAlertConfig(t *testing.T) {
	applicationAlertConfig := ut.createTestApplicationAlertConfig()

	serializedJSON, _ := json.Marshal(applicationAlertConfig)

	result, err := NewApplicationAlertConfigUnmarshaller().Unmarshal(serializedJSON)

	require.NoError(t, err)
	require.Equal(t, applicationAlertConfig, result)
}

func (ut *applicationAlertConfigUnmarshallerTest) shouldFailToUnmarshalApplicationAlertConfigWhenResponseIsAJsonArray(t *testing.T) {
	response := `["foo","bar"]`

	_, err := NewApplicationAlertConfigUnmarshaller().Unmarshal([]byte(response))

	require.Error(t, err)
}

func (ut *applicationAlertConfigUnmarshallerTest) shouldReturnEmptyApplicationAlertConfigWhenNoFieldOfResponseMatchesToModel(t *testing.T) {
	response := `{"foo" : "bar"}`
	config, err := NewApplicationAlertConfigUnmarshaller().Unmarshal([]byte(response))

	require.NoError(t, err)
	require.Equal(t, &ApplicationAlertConfig{}, config)
}

func (ut *applicationAlertConfigUnmarshallerTest) shouldFailToUnmarshalApplicationAlertConfigWhenResponseIsNotAValidJson(t *testing.T) {
	response := `Invalid Data`

	_, err := NewApplicationAlertConfigUnmarshaller().Unmarshal([]byte(response))

	require.Error(t, err)
}

func (ut *applicationAlertConfigUnmarshallerTest) shouldSuccessfullyUnmarshalApplicationAlertConfigArray(t *testing.T) {
	applicationAlertConfig := ut.createTestApplicationAlertConfig()
	input := []*ApplicationAlertConfig{applicationAlertConfig}

	serializedJSON, _ := json.Marshal(&input)

	result, err := NewApplicationAlertConfigUnmarshaller().UnmarshalArray(serializedJSON)

	require.NoError(t, err)
	require.Equal(t, &input, result)
}

func (ut *applicationAlertConfigUnmarshallerTest) shouldFailToUnmarshalApplicationAlertConfigArrayContainingAtLeastOneInvalidApplicationAlertConfig(t *testing.T) {
	applicationAlertConfig := ut.createTestApplicationAlertConfig()
	objectJson, _ := json.Marshal(applicationAlertConfig)

	serializedJSON := fmt.Sprintf("[%s,[\"foo\",\"bar\"]]", objectJson)

	_, err := NewApplicationAlertConfigUnmarshaller().UnmarshalArray([]byte(serializedJSON))

	require.Error(t, err)
}

func (ut *applicationAlertConfigUnmarshallerTest) shouldFailToUnmarshalApplicationAlertConfigArrayWhenNoValidJsonIsProvided(t *testing.T) {
	_, err := NewApplicationAlertConfigUnmarshaller().UnmarshalArray([]byte("invalid json data"))

	require.Error(t, err)
}

func (ut *applicationAlertConfigUnmarshallerTest) createTestApplicationAlertConfig() *ApplicationAlertConfig {
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
		TagFilterExpression: NewLogicalAndTagFilter([]*TagFilter{NewStringTagFilter(TagFilterEntitySource, "service.name", EqualsOperator, "test"), NewStringTagFilter(TagFilterEntitySource, "entity.type", EqualsOperator, "host")}),
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
	return &applicationAlertConfig
}
