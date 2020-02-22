package restapi_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

func TestShouldSuccessfullyUnmarshalAlertingChannel(t *testing.T) {
	response := `{
		"id" : "test-id",
		"name" : "test-name",
		"kind" : "EMAIL",
		"emails" : ["test-email1","test-email2"]
	}`

	result, err := NewAlertingChannelUnmarshaller().Unmarshal([]byte(response))

	assert.Nil(t, err)

	alertingChannel, ok := result.(AlertingChannel)
	assert.True(t, ok)

	assert.Equal(t, "test-id", alertingChannel.ID)
	assert.Equal(t, "test-name", alertingChannel.Name)
	assert.Equal(t, EmailChannelType, alertingChannel.Kind)
	assert.Equal(t, []string{"test-email1", "test-email2"}, alertingChannel.Emails)
}

func TestShouldFailToUnmarshalAlertingChannelWhenResponseIsAJsonArray(t *testing.T) {
	response := `["test-email1","test-email2"]`

	_, err := NewAlertingChannelUnmarshaller().Unmarshal([]byte(response))

	assert.NotNil(t, err)
}

func TestShouldFailToUnmarshalAlertingChannelWhenResponseIsNotAJsonMessage(t *testing.T) {
	response := `foo bar`

	_, err := NewAlertingChannelUnmarshaller().Unmarshal([]byte(response))

	assert.NotNil(t, err)
}

func TestShouldReturnEmptyAlertingChannelWhenJsonObjectIsReturnWhereNoFiledMatches(t *testing.T) {
	response := `{"foo" : "bar" }`

	result, err := NewAlertingChannelUnmarshaller().Unmarshal([]byte(response))

	assert.Nil(t, err)
	assert.Equal(t, AlertingChannel{}, result)
}
