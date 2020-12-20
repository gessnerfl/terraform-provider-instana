package restapi_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

func TestShouldSuccessfullyUnmarshalSliConfig(t *testing.T) {
	sliConfig := SliConfig{
		ID:                         sliConfigID,
		Name:                       sliConfigName,
		InitialEvaluationTimestamp: sliConfigInitialEvaluationTimestamp,
		MetricConfiguration: MetricConfiguration{
			Name:        sliConfigMetricName,
			Aggregation: sliConfigMetricAggregation,
			Threshold:   sliConfigMetricThreshold,
		},
		SliEntity: SliEntity{
			Type:          sliConfigEntityType,
			ApplicationID: sliConfigEntityApplicationID,
			ServiceID:     sliConfigEntityServiceID,
			EndpointID:    sliConfigEntityEndpointID,
			BoundaryScope: sliConfigEntityBoundaryScope,
		},
	}

	serializedJSON, _ := json.Marshal(sliConfig)

	result, err := NewSliConfigUnmarshaller().Unmarshal(serializedJSON)

	assert.Nil(t, err)
	assert.Equal(t, sliConfig, result)
}

func TestShouldFailToUnmarshalSliConfigWhenResponseIsAJsonArray(t *testing.T) {
	response := `["foo","bar"]`

	_, err := NewSliConfigUnmarshaller().Unmarshal([]byte(response))

	assert.NotNil(t, err)
}

func TestShouldFailToUnmarshalSliConfigWhenResponseIsNotAJsonMessage(t *testing.T) {
	response := `foo bar`

	_, err := NewSliConfigUnmarshaller().Unmarshal([]byte(response))

	assert.NotNil(t, err)
}

func TestShouldReturnEmptySliConfigWhenJsonObjectIsReturnWhereNoFiledMatches(t *testing.T) {
	response := `{"foo" : "bar" }`

	result, err := NewSliConfigUnmarshaller().Unmarshal([]byte(response))

	assert.Nil(t, err)
	assert.Equal(t, SliConfig{}, result)
}
