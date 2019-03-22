package services

import "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
import "github.com/gessnerfl/terraform-provider-instana/instana/restapi/resources"

//NewInstanaAPI creates a new instance of the instana API
func NewInstanaAPI(apiToken string, endpoint string) restapi.InstanaAPI {
	client := NewClient(apiToken, endpoint)
	return &baseInstanaAPI{client: client}
}

type baseInstanaAPI struct {
	client restapi.RestClient
}

//Rules implementation of InstanaAPI interface
func (api baseInstanaAPI) Rules() restapi.RuleResource {
	return resources.NewRuleResource(api.client)
}

//RuleBindings implementation of InstanaAPI interface
func (api baseInstanaAPI) RuleBindings() restapi.RuleBindingResource {
	return resources.NewRuleBindingResource(api.client)
}

//UserRoles implementation of InstanaAPI interface
func (api baseInstanaAPI) UserRoles() restapi.UserRoleResource {
	return resources.NewUserRoleResource(api.client)
}
