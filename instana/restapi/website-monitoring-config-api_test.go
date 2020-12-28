package restapi_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
)

const (
	websiteMonitoringConfigID      = "id"
	websiteMonitoringConfigName    = "name"
	websiteMonitoringConfigAppName = "appNme"
)

func TestShouldSuccessfullyValidateWebsiteMonitotingConfigWithNameOnlyForCreateUseCase(t *testing.T) {
	sut := WebsiteMonitoringConfig{Name: websiteMonitoringConfigName}

	err := sut.Validate()

	require.NoError(t, err)
	require.Equal(t, websiteMonitoringConfigName, sut.Name)
}

func TestShouldSuccessfullyValidateFullWebsiteMonitotingConfigAsServerResponseForCreateOrUpdateOperation(t *testing.T) {
	sut := WebsiteMonitoringConfig{
		ID:      websiteMonitoringConfigID,
		Name:    websiteMonitoringConfigName,
		AppName: websiteMonitoringConfigAppName,
	}

	err := sut.Validate()

	require.NoError(t, err)
	require.Equal(t, websiteMonitoringConfigID, sut.ID)
	require.Equal(t, websiteMonitoringConfigID, sut.GetID())
	require.Equal(t, websiteMonitoringConfigName, sut.Name)
	require.Equal(t, websiteMonitoringConfigAppName, sut.AppName)
}

func TestShouldFailToValidateWebsiteMonitoringConfigWhenNameIsNotProvided(t *testing.T) {
	sut := WebsiteMonitoringConfig{}

	err := sut.Validate()

	require.Error(t, err)
	require.Contains(t, err.Error(), "Name is missing")
}
