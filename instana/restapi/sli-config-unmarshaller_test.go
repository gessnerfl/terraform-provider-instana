package restapi_test

import (
	"encoding/json"
	"fmt"
	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestShouldSuccessfullyUnmarshalSliConfigWithSliEntityOfTypeApplication(t *testing.T) {
	sliConfig := createTestSliConfigWithSliEntityOfTypeApplication()

	serializedJSON, _ := json.Marshal(sliConfig)

	result, err := NewSliConfigUnmarshaller().Unmarshal(serializedJSON)

	require.NoError(t, err)
	require.Equal(t, sliConfig, result)
}

func TestShouldSuccessfullyUnmarshalSliConfigWithSliEntityOfTypeAvailability(t *testing.T) {
	sliConfig := createTestSliConfigWithSliEntityOfTypeAvailability()

	serializedJSON, _ := json.Marshal(sliConfig)

	result, err := NewSliConfigUnmarshaller().Unmarshal(serializedJSON)

	require.NoError(t, err)
	require.Equal(t, sliConfig, result)
}

func TestShouldSuccessfullyUnmarshalSliConfigWithSliEntityOfTypeWebsiteTimeBased(t *testing.T) {
	sliConfig := createTestSliConfigWithSliEntityOfTypeWebsiteTimeBased()

	serializedJSON, _ := json.Marshal(sliConfig)

	result, err := NewSliConfigUnmarshaller().Unmarshal(serializedJSON)

	require.NoError(t, err)
	require.Equal(t, sliConfig, result)
}

func TestShouldSuccessfullyUnmarshalSliConfigWithSliEntityOfTypeWebsiteEventBased(t *testing.T) {
	sliConfig := createTestSliConfigWithSliEntityOfTypeWebsiteEventBased()

	serializedJSON, _ := json.Marshal(sliConfig)

	result, err := NewSliConfigUnmarshaller().Unmarshal(serializedJSON)

	require.NoError(t, err)
	require.Equal(t, sliConfig, result)
}

func TestShouldFailToUnmarshalSliConfigWhenResponseIsAJsonArray(t *testing.T) {
	response := `["foo","bar"]`

	_, err := NewSliConfigUnmarshaller().Unmarshal([]byte(response))

	require.Error(t, err)
}

func TestShouldReturnEmptySliConfigWhenNoFieldOfResponseMatchesToModel(t *testing.T) {
	response := `{"foo" : "bar"}`
	config, err := NewSliConfigUnmarshaller().Unmarshal([]byte(response))

	require.NoError(t, err)
	require.Equal(t, &SliConfig{}, config)
}

func TestShouldFailToUnmarshalSliConfigWhenResponseIsNotAValidJson(t *testing.T) {
	response := `Invalid Data`

	_, err := NewSliConfigUnmarshaller().Unmarshal([]byte(response))

	require.Error(t, err)
}

func TestShouldFailToUnmarshalSliConfigWhenGoodEventFilterExpressionIsNotValid(t *testing.T) {
	response := `
{
	"id" : "sli-config.i",
	"sliName" : "sli-config-name",
	"initialEvaluationTimestamp": 0,
	"metricConfiguration": {
		"metricName" : "metric-name",
		"metricAggregation"	: "SUM",
		"threshold" : 1.0
	},
	"sliEntity": {
		"sliType" : "websiteEventBased",
		"websiteId" : "website_id",
		"beaconType" : "pageLoad",
		"goodEventFilterExpression" : ["foo", "bar"],
		"badEventFilterExpression" : {
			"type" : "TAG_FILTER",
			"name" : "request.path",
			"entity" : "DESTINATION",
			"operator" : "EQUALS",
			"stringValue" : "/404",
			"value" : "/404"
		}
	}
}
`
	_, err := NewSliConfigUnmarshaller().Unmarshal([]byte(response))

	require.Error(t, err)
}

func TestShouldFailToUnmarshalSliConfigWhenBadEventFilterExpressionIsNotValid(t *testing.T) {
	response := `
{
	"id" : "sli-config.i",
	"sliName" : "sli-config-name",
	"initialEvaluationTimestamp": 0,
	"metricConfiguration": {
		"metricName" : "metric-name",
		"metricAggregation"	: "SUM",
		"threshold" : 1.0
	},
	"sliEntity": {
		"sliType" : "websiteEventBased",
		"websiteId" : "website_id",
		"beaconType" : "pageLoad",
		"goodEventFilterExpression" : {
			"type" : "TAG_FILTER",
			"name" : "request.path",
			"entity" : "DESTINATION",
			"operator" : "EQUALS",
			"stringValue" : "/404",
			"value" : "/404"
		},
		"badEventFilterExpression" : ["foo", "bar"]
	}
}
`
	_, err := NewSliConfigUnmarshaller().Unmarshal([]byte(response))

	require.Error(t, err)
}

func TestShouldFailToUnmarshalSliConfigWhenFilterExpressionIsNotValid(t *testing.T) {
	response := `
{
	"id" : "sli-config.i",
	"sliName" : "sli-config-name",
	"initialEvaluationTimestamp": 0,
	"metricConfiguration": {
		"metricName" : "metric-name",
		"metricAggregation"	: "SUM",
		"threshold" : 1.0
	},
	"sliEntity": {
		"sliType" : "websiteTimeBased",
		"websiteId" : "website_id",
		"beaconType" : "pageLoad",
		"filterExpression" : ["foo", "bar"]
	}
}
`
	_, err := NewSliConfigUnmarshaller().Unmarshal([]byte(response))

	require.Error(t, err)
}

func TestShouldSuccessfullyUnmarshalSliConfigArray(t *testing.T) {
	sliConfig := createTestSliConfigWithSliEntityOfTypeApplication()
	input := []*SliConfig{sliConfig}

	serializedJSON, _ := json.Marshal(&input)

	result, err := NewSliConfigUnmarshaller().UnmarshalArray(serializedJSON)

	require.NoError(t, err)
	require.Equal(t, &input, result)
}

func TestShouldFailToUnmarshalSliConfigArrayContainingAtLeastOneInvalidSliConfig(t *testing.T) {
	sliConfig := createTestSliConfigWithSliEntityOfTypeApplication()
	objectJson, _ := json.Marshal(sliConfig)

	serializedJSON := fmt.Sprintf("[%s,[\"foo\",\"bar\"]]", objectJson)

	_, err := NewSliConfigUnmarshaller().UnmarshalArray([]byte(serializedJSON))

	require.Error(t, err)
}

func TestShouldFailToUnmarshalSliConfigArrayyWhenNoValidJsonIsProvided(t *testing.T) {
	_, err := NewSliConfigUnmarshaller().UnmarshalArray([]byte("invalid json data"))

	require.Error(t, err)
}

func createTestSliConfigWithSliEntityOfTypeApplication() *SliConfig {
	applicationID := "application-id"
	serviceID := "service-id"
	endpointID := "endpoint-id"
	boundaryScope := string(BoundaryScopeAll)
	sliConfig := SliConfig{
		ID:                         "sli-config-id",
		Name:                       "sli-config-name",
		InitialEvaluationTimestamp: 10,
		MetricConfiguration: MetricConfiguration{
			Name:        "metric-config-name",
			Aggregation: "SUM",
			Threshold:   0.5,
		},
		SliEntity: SliEntity{
			Type:          "application",
			ApplicationID: &applicationID,
			ServiceID:     &serviceID,
			EndpointID:    &endpointID,
			BoundaryScope: &boundaryScope,
		},
	}
	return &sliConfig
}

func createTestSliConfigWithSliEntityOfTypeAvailability() *SliConfig {
	applicationID := "application-id"
	goodEventFilterExpression := NewStringTagFilter(TagFilterEntityDestination, "entity", EqualsOperator, "good")
	badEventFilterExpression := NewStringTagFilter(TagFilterEntityDestination, "entity", EqualsOperator, "bad")
	boundaryScope := string(BoundaryScopeAll)
	sliConfig := SliConfig{
		ID:                         "sli-config-id",
		Name:                       "sli-config-name",
		InitialEvaluationTimestamp: 10,
		MetricConfiguration: MetricConfiguration{
			Name:        "metric-config-name",
			Aggregation: "SUM",
			Threshold:   0.5,
		},
		SliEntity: SliEntity{
			Type:                      "availability",
			ApplicationID:             &applicationID,
			BoundaryScope:             &boundaryScope,
			GoodEventFilterExpression: goodEventFilterExpression,
			BadEventFilterExpression:  badEventFilterExpression,
		},
	}
	return &sliConfig
}

func createTestSliConfigWithSliEntityOfTypeWebsiteEventBased() *SliConfig {
	WebsiteID := "website-id"
	goodEventFilterExpression := NewStringTagFilter(TagFilterEntityDestination, "entity", EqualsOperator, "good")
	badEventFilterExpression := NewStringTagFilter(TagFilterEntityDestination, "entity", EqualsOperator, "bad")
	beaconType := "pageLoad"
	sliConfig := SliConfig{
		ID:                         "sli-config-id",
		Name:                       "sli-config-name",
		InitialEvaluationTimestamp: 10,
		MetricConfiguration: MetricConfiguration{
			Name:        "metric-config-name",
			Aggregation: "SUM",
			Threshold:   0.5,
		},
		SliEntity: SliEntity{
			Type:                      "availability",
			WebsiteId:                 &WebsiteID,
			BeaconType:                &beaconType,
			GoodEventFilterExpression: goodEventFilterExpression,
			BadEventFilterExpression:  badEventFilterExpression,
		},
	}
	return &sliConfig
}

func createTestSliConfigWithSliEntityOfTypeWebsiteTimeBased() *SliConfig {
	WebsiteID := "website-id"
	filterExpression := NewStringTagFilter(TagFilterEntityDestination, "entity", EqualsOperator, "value")
	beaconType := "pageLoad"
	sliConfig := SliConfig{
		ID:                         "sli-config-id",
		Name:                       "sli-config-name",
		InitialEvaluationTimestamp: 10,
		MetricConfiguration: MetricConfiguration{
			Name:        "metric-config-name",
			Aggregation: "SUM",
			Threshold:   0.5,
		},
		SliEntity: SliEntity{
			Type:             "availability",
			WebsiteId:        &WebsiteID,
			BeaconType:       &beaconType,
			FilterExpression: filterExpression,
		},
	}
	return &sliConfig
}
