package restapi_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
)

func TestShouldReturnResourcesFromInstanaAPI(t *testing.T) {
	api := NewInstanaAPI("api-token", "endpoint", false)

	t.Run("Should return CustomEventSpecification instance", func(t *testing.T) {
		resource := api.CustomEventSpecifications()

		require.NotNil(t, resource)
	})

	t.Run("Should return BuiltinEventSpecifications instance", func(t *testing.T) {
		resource := api.BuiltinEventSpecifications()

		require.NotNil(t, resource)
	})
	t.Run("Should return APITokens instance", func(t *testing.T) {
		resource := api.APITokens()

		require.NotNil(t, resource)
	})
	t.Run("Should return ApplicationConfig instance", func(t *testing.T) {
		resource := api.ApplicationConfigs()

		require.NotNil(t, resource)
	})
	t.Run("Should return ApplicationAlertConfig instance", func(t *testing.T) {
		resource := api.ApplicationAlertConfigs()

		require.NotNil(t, resource)
	})
	t.Run("Should return GlobalApplicationAlertConfig instance", func(t *testing.T) {
		resource := api.GlobalApplicationAlertConfigs()

		require.NotNil(t, resource)
	})
	t.Run("Should return AlertingChannel instance", func(t *testing.T) {
		resource := api.AlertingChannels()

		require.NotNil(t, resource)
	})
	t.Run("Should return AlertingConfiguration instance", func(t *testing.T) {
		resource := api.AlertingConfigurations()

		require.NotNil(t, resource)
	})
	t.Run("Should return SliConfig instance", func(t *testing.T) {
		resource := api.SliConfigs()

		require.NotNil(t, resource)
	})
	t.Run("Should return WebsiteMonitoringConfig instance", func(t *testing.T) {
		resource := api.WebsiteMonitoringConfig()

		require.NotNil(t, resource)
	})
	t.Run("Should return WebsiteAlertConfig instance", func(t *testing.T) {
		resource := api.WebsiteAlertConfig()

		require.NotNil(t, resource)
	})
	t.Run("Should return Groups instance", func(t *testing.T) {
		resource := api.Groups()

		require.NotNil(t, resource)
	})
	t.Run("Should return Custom Dashboard instance", func(t *testing.T) {
		resource := api.CustomDashboards()

		require.NotNil(t, resource)
	})
	t.Run("Should return Synthetic test instance", func(t *testing.T) {
		resource := api.SyntheticTest()

		require.NotNil(t, resource)
	})
	t.Run("Should return Synthetic location instance", func(t *testing.T) {
		resource := api.SyntheticLocation()

		require.NotNil(t, resource)
	})

}
