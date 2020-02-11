package restapi_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

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

	if err != nil {
		t.Fatalf("Expected to successfully unmarshal alerting channel response; %s", err)
	}

	alertingChannel, ok := result.(AlertingChannel)
	if !ok {
		t.Fatal("Expected result to be a alerting channel")
	}

	if alertingChannel.ID != "test-id" {
		t.Fatal("Expected ID to be properly mapped")
	}
	if alertingChannel.Name != "test-name" {
		t.Fatal("Expected name to be properly mapped")
	}
	if alertingChannel.Kind != EmailChannelType {
		t.Fatal("Expected kind to be properly mapped")
	}
	if !cmp.Equal(alertingChannel.Emails, []string{"test-email1", "test-email2"}) {
		t.Fatal("Expected emails to be properly mapped")
	}
}

func TestShouldFailToUnmarshalAlertingChannelWhenResponseIsAJsonArray(t *testing.T) {
	response := `["test-email1","test-email2"]`

	_, err := NewAlertingChannelUnmarshaller().Unmarshal([]byte(response))

	if err == nil {
		t.Fatal("Expected unmarshalling to fail")
	}
}

func TestShouldFailToUnmarshalAlertingChannelWhenResponseIsNotAJsonMessage(t *testing.T) {
	response := `foo bar`

	_, err := NewAlertingChannelUnmarshaller().Unmarshal([]byte(response))

	if err == nil {
		t.Fatal("Expected unmarshalling to fail")
	}
}

func TestShouldReturnEmptyAlertingChannelWhenJsonObjectIsReturnWhereNoFiledMatches(t *testing.T) {
	response := `{"foo" : "bar" }`

	result, err := NewAlertingChannelUnmarshaller().Unmarshal([]byte(response))

	if err != nil {
		t.Fatalf("Expected to successfully unmarshal alerting channel response, %s", err)
	}

	if !cmp.Equal(result, AlertingChannel{}) {
		t.Fatal("Expected empty alerting channel")
	}
}
