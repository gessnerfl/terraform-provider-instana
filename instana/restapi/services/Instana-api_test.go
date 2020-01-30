package services_test

import (
	"github.com/google/go-cmp/cmp"
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi/services"
)

func TestShouldReturnResourcesFromInstanaAPI(t *testing.T) {
	api := NewInstanaAPI("api-token", "endpoint")

	t.Run("Should return CustomEventSpecificationResource instance", func(t *testing.T) {
		customEventSpecificationResource := api.CustomEventSpecifications()
		if customEventSpecificationResource == nil {
			t.Fatal("Expected instance of CustomEventSpecificationResource to be returned")
		}
	})
	t.Run("Should return UserRoleResource instance", func(t *testing.T) {
		userRoleResource := api.UserRoles()
		if userRoleResource == nil {
			t.Fatal("Expected instance of UserRoleResource to be returned")
		}
	})
	t.Run("Should return ApplicationConfigResource instance", func(t *testing.T) {
		applicationConfigResource := api.ApplicationConfigs()
		if applicationConfigResource == nil {
			t.Fatal("Expected instance of ApplicationConfigResource to be returned")
		}
	})
	t.Run("Should return AlertingChannelResource instance", func(t *testing.T) {
		alertingChannelResource := api.AlertingChannels()
		if alertingChannelResource == nil {
			t.Fatal("Expected instance of AlertingChannelResource to be returned")
		}
	})
}

//Add tests for unmarshal

func TestShouldSuccessfullyUnmarshalAlertingChannel(t *testing.T) {
	response := `{
		"id" : "test-id",
		"name" : "test-name",
		"kind" : "EMAIL",
		"emails" : ["test-email1","test-email2"]
	}`

	result, err := UnmarshalAlertingChannel([]byte(response))

	if err != nil {
		t.Fatalf("Expected to successfully unmarshal alerting channel response; %s", err)
	}

	alertingChannel, ok := result.(restapi.AlertingChannel)
	if !ok {
		t.Fatal("Expected result to be a alerting channel")
	}

	if alertingChannel.ID != "test-id" {
		t.Fatal("Expected ID to be properly mapped")
	}
	if alertingChannel.Name != "test-name" {
		t.Fatal("Expected name to be properly mapped")
	}
	if alertingChannel.Kind != restapi.EmailChannelType {
		t.Fatal("Expected kind to be properly mapped")
	}
	if !cmp.Equal(alertingChannel.Emails, []string{"test-email1", "test-email2"}) {
		t.Fatal("Expected emails to be properly mapped")
	}
}

func TestShouldFailToUnmarshalCustomEventSpecificationWhenResponseIsAJsonArray(t *testing.T) {
	response := `["foo","bar"]`

	_, err := UnmarshalCustomEventSpecification([]byte(response))

	if err == nil {
		t.Fatal("Expected unmarshalling to fail")
	}
}

func TestShouldFailToUnmarshalCustomEventSpecificationsWhenResponseIsNotAJsonMessage(t *testing.T) {
	response := `foo bar`

	_, err := UnmarshalCustomEventSpecification([]byte(response))

	if err == nil {
		t.Fatal("Expected unmarshalling to fail")
	}
}

func TestShouldReturnEmptyCustomEventSpecificationWhenJsonObjectIsReturnWhereNoFiledMatches(t *testing.T) {
	response := `{"foo" : "bar" }`

	result, err := UnmarshalCustomEventSpecification([]byte(response))

	if err != nil {
		t.Fatalf("Expected to successfully unmarshal custom event specification response, %s", err)
	}

	if !cmp.Equal(result, restapi.CustomEventSpecification{}) {
		t.Fatal("Expected empty custom event specification")
	}
}

func TestShouldFailToUnmarshalUserRoleWhenResponseIsAJsonArray(t *testing.T) {
	response := `["foo","bar"]`

	_, err := UnmarshalUserRole([]byte(response))

	if err == nil {
		t.Fatal("Expected unmarshalling to fail")
	}
}

func TestShouldFailToUnmarshalUserRoleWhenResponseIsNotAJsonMessage(t *testing.T) {
	response := `foo bar`

	_, err := UnmarshalUserRole([]byte(response))

	if err == nil {
		t.Fatal("Expected unmarshalling to fail")
	}
}

func TestShouldReturnEmptyUserRoleWhenJsonObjectIsReturnWhereNoFiledMatches(t *testing.T) {
	response := `{"foo" : "bar" }`

	result, err := UnmarshalUserRole([]byte(response))

	if err != nil {
		t.Fatalf("Expected to successfully unmarshal user role response, %s", err)
	}

	if !cmp.Equal(result, restapi.UserRole{}) {
		t.Fatal("Expected empty user role")
	}
}

func TestShouldFailToUnmarshalAlertingChannelWhenResponseIsAJsonArray(t *testing.T) {
	response := `["test-email1","test-email2"]`

	_, err := UnmarshalAlertingChannel([]byte(response))

	if err == nil {
		t.Fatal("Expected unmarshalling to fail")
	}
}

func TestShouldFailToUnmarshalAlertingChannelWhenResponseIsNotAJsonMessage(t *testing.T) {
	response := `foo bar`

	_, err := UnmarshalAlertingChannel([]byte(response))

	if err == nil {
		t.Fatal("Expected unmarshalling to fail")
	}
}

func TestShouldReturnEmptyAlertingChannelWhenJsonObjectIsReturnWhereNoFiledMatches(t *testing.T) {
	response := `{"foo" : "bar" }`

	result, err := UnmarshalAlertingChannel([]byte(response))

	if err != nil {
		t.Fatalf("Expected to successfully unmarshal alerting channel response, %s", err)
	}

	if !cmp.Equal(result, restapi.AlertingChannel{}) {
		t.Fatal("Expected empty alerting channel")
	}
}
