package resources

import (
	"encoding/json"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

//NewAlertingChannelResource constructs a new instance of AlertingChannelResource
func NewAlertingChannelResource(client restapi.RestClient) restapi.AlertingChannelResource {
	return &AlertingChannelResourceImpl{
		client:       client,
		resourcePath: restapi.AlertingChannelsResourcePath,
	}
}

//AlertingChannelResourceImpl is the GO representation of the Alerting Channel Resource of Instana
type AlertingChannelResourceImpl struct {
	client       restapi.RestClient
	resourcePath string
}

//GetOne retrieves a single custom Alerting Channel from Instana API by its ID
func (resource *AlertingChannelResourceImpl) GetOne(id string) (restapi.AlertingChannel, error) {
	data, err := resource.client.GetOne(id, resource.resourcePath)
	if err != nil {
		return restapi.AlertingChannel{}, err
	}
	return resource.validateResponseAndConvertToStruct(data)
}

func (resource *AlertingChannelResourceImpl) validateResponseAndConvertToStruct(data []byte) (restapi.AlertingChannel, error) {
	alertingChannel := restapi.AlertingChannel{}
	if err := json.Unmarshal(data, &alertingChannel); err != nil {
		return alertingChannel, fmt.Errorf("failed to parse json; %s", err)
	}

	if err := alertingChannel.Validate(); err != nil {
		return alertingChannel, err
	}
	return alertingChannel, nil
}

//Upsert creates or updates a user role
func (resource *AlertingChannelResourceImpl) Upsert(alertingChannel restapi.AlertingChannel) (restapi.AlertingChannel, error) {
	if err := alertingChannel.Validate(); err != nil {
		return alertingChannel, err
	}
	data, err := resource.client.Put(alertingChannel, resource.resourcePath)
	if err != nil {
		return alertingChannel, err
	}
	return resource.validateResponseAndConvertToStruct(data)
}

//Delete deletes a user role
func (resource *AlertingChannelResourceImpl) Delete(alertingChannel restapi.AlertingChannel) error {
	return resource.DeleteByID(alertingChannel.ID)
}

//DeleteByID deletes a user role by its ID
func (resource *AlertingChannelResourceImpl) DeleteByID(alertingChannelID string) error {
	return resource.client.Delete(alertingChannelID, resource.resourcePath)
}
