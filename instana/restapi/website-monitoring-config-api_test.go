package restapi_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
)

const (
	WebsiteMonitoringConfigID      = "id"
	WebsiteMonitoringConfigName    = "name"
	WebsiteMonitoringConfigAppName = "appNme"
)

func TestShouldSuccessfullyValidateWebsiteMonitotingConfigWithNameOnlyForCreateUseCase(t *testing.T) {
	sut := WebsiteMonitoringConfig{Name: WebsiteMonitoringConfigName}

	err := sut.Validate()

	require.NoError(t, err)
	require.Equal(t, WebsiteMonitoringConfigName, sut.Name)
}

func TestShouldSuccessfullyValidateFullWebsiteMonitotingConfigAsServerResponseForCreateOrUpdateOperation(t *testing.T) {
	sut := WebsiteMonitoringConfig{
		ID:      WebsiteMonitoringConfigID,
		Name:    WebsiteMonitoringConfigName,
		AppName: WebsiteMonitoringConfigAppName,
	}

	err := sut.Validate()

	require.NoError(t, err)
	require.Equal(t, WebsiteMonitoringConfigID, sut.ID)
	require.Equal(t, WebsiteMonitoringConfigID, sut.GetID())
	require.Equal(t, WebsiteMonitoringConfigName, sut.Name)
	require.Equal(t, WebsiteMonitoringConfigAppName, sut.AppName)
}

func TestShouldFailToValidateWebsiteMonitoringConfigWhenNameIsNotProvided(t *testing.T) {
	sut := WebsiteMonitoringConfig{}

	err := sut.Validate()

	require.Error(t, err)
	require.Contains(t, err.Error(), "Name is missing")
}
