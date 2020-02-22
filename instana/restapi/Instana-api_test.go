package restapi_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnResourcesFromInstanaAPI(t *testing.T) {
	api := NewInstanaAPI("api-token", "endpoint")

	t.Run("Should return CustomEventSpecification instance", func(t *testing.T) {
		resource := api.CustomEventSpecifications()

		assert.NotNil(t, resource)
	})
	t.Run("Should return UserRole instance", func(t *testing.T) {
		resource := api.UserRoles()

		assert.NotNil(t, resource)
	})
	t.Run("Should return ApplicationConfig instance", func(t *testing.T) {
		resource := api.ApplicationConfigs()

		assert.NotNil(t, resource)
	})
	t.Run("Should return AlertingChannel instance", func(t *testing.T) {
		resource := api.AlertingChannels()

		assert.NotNil(t, resource)
	})
	t.Run("Should return AlertingConfiguration instance", func(t *testing.T) {
		resource := api.AlertingConfigurations()

		assert.NotNil(t, resource)
	})
}
