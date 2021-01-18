package restapi_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

func TestShouldSuccessfullyUnmarshalWebsiteMonitoringConfig(t *testing.T) {
	config := WebsiteMonitoringConfig{
		ID:      websiteMonitoringConfigID,
		Name:    websiteMonitoringConfigName,
		AppName: websiteMonitoringConfigAppName,
	}

	serializedJSON, _ := json.Marshal(config)

	result, err := NewWebsiteMonitoringConfigUnmarshaller().Unmarshal(serializedJSON)

	assert.Nil(t, err)
	assert.Equal(t, config, result)
}

func TestShouldFailToUnmarshalWebsiteMonitoringConfigWhenResponseIsAJsonArray(t *testing.T) {
	response := `["foo","bar"]`

	_, err := NewWebsiteMonitoringConfigUnmarshaller().Unmarshal([]byte(response))

	assert.NotNil(t, err)
}

func TestShouldFailToUnmarshalWebsiteMonitoringConfigWhenResponseIsNotAJsonMessage(t *testing.T) {
	response := `foo bar`

	_, err := NewWebsiteMonitoringConfigUnmarshaller().Unmarshal([]byte(response))

	assert.NotNil(t, err)
}

func TestShouldReturnEmptyWebsiteMonitoringConfigWhenJsonObjectIsReturnWhereNoFiledMatches(t *testing.T) {
	response := `{"foo" : "bar" }`

	result, err := NewWebsiteMonitoringConfigUnmarshaller().Unmarshal([]byte(response))

	assert.Nil(t, err)
	assert.Equal(t, WebsiteMonitoringConfig{}, result)
}
