package services

import (
	"encoding/json"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi/resources"
)

//NewInstanaAPI creates a new instance of the instana API
func NewInstanaAPI(apiToken string, endpoint string) restapi.InstanaAPI {
	client := NewClient(apiToken, endpoint)
	return &baseInstanaAPI{client: client}
}

type baseInstanaAPI struct {
	client restapi.RestClient
}

//CustomEventSpecifications implementation of InstanaAPI interface
func (api *baseInstanaAPI) CustomEventSpecifications() restapi.CustomEventSpecificationResource {
	return resources.NewCustomEventSpecificationResource(api.client)
}

//UserRoles implementation of InstanaAPI interface
func (api *baseInstanaAPI) UserRoles() restapi.UserRoleResource {
	return resources.NewUserRoleResource(api.client)
}

//ApplicationConfigs implementation of InstanaAPI interface
func (api *baseInstanaAPI) ApplicationConfigs() restapi.ApplicationConfigResource {
	return resources.NewApplicationConfigResource(api.client)
}

//AlertingChannels implementation of InstanaAPI interface
func (api *baseInstanaAPI) AlertingChannels() restapi.RestResource {
	return NewRestResource(restapi.AlertingChannelsResourcePath, UnmarshalAlertingChannel, api.client)
}

//UnmarshalAlertingChannel unmarshal the JSON response of an alerting channel to an AlertingChannel struct
func UnmarshalAlertingChannel(data []byte) (restapi.InstanaDataObject, error) {
	alertingChannel := restapi.AlertingChannel{}
	if err := json.Unmarshal(data, &alertingChannel); err != nil {
		return alertingChannel, fmt.Errorf("failed to parse json; %s", err)
	}
	return alertingChannel, nil
}
