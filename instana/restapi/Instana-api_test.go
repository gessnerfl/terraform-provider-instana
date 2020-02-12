package restapi_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

func TestShouldReturnResourcesFromInstanaAPI(t *testing.T) {
	api := NewInstanaAPI("api-token", "endpoint")

	t.Run("Should return CustomEventSpecification instance", func(t *testing.T) {
		resource := api.CustomEventSpecifications()
		if resource == nil {
			t.Fatal("Expected instance of RestResource to be returned for CustomEventSpecifications")
		}
	})
	t.Run("Should return UserRole instance", func(t *testing.T) {
		resource := api.UserRoles()
		if resource == nil {
			t.Fatal("Expected instance of RestResource to be returned for UserRoles")
		}
	})
	t.Run("Should return ApplicationConfig instance", func(t *testing.T) {
		resource := api.ApplicationConfigs()
		if resource == nil {
			t.Fatal("Expected instance of RestResource to be returned for ApplicationConfigs")
		}
	})
	t.Run("Should return AlertingChannel instance", func(t *testing.T) {
		resource := api.AlertingChannels()
		if resource == nil {
			t.Fatal("Expected instance of RestResource to be returned for AlertingChannels")
		}
	})
	t.Run("Should return AlertingConfiguration instance", func(t *testing.T) {
		resource := api.AlertingConfigurations()
		if resource == nil {
			t.Fatal("Expected instance of RestResource to be returned for AlertingConfigurations")
		}
	})
}
