package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi/services"
	"github.com/hashicorp/terraform/helper/schema"
)

//SchemaFieldAPIToken the name of the provider configuration option for the api token
const SchemaFieldAPIToken = "api_token"

//SchemaFieldEndpoint the name of the provider configuration option for the instana endpoint
const SchemaFieldEndpoint = "endpoint"

//SchemaFieldAddTerraformManagedString the indicator if a string should be appended to a name or description that the resource is managed by terraform
const SchemaFieldAddTerraformManagedString = "add_terraform_managed_string"

//ResourceInstanaRule the name of the terraform-provider-instana resource to manage rules
const ResourceInstanaRule = "instana_rule"

//ResourceInstanaRuleBinding the name of the terraform-provider-instana resource to manage rule bindings
const ResourceInstanaRuleBinding = "instana_rule_binding"

//ResourceInstanaUserRole the name of the terraform-provider-instana resource to manage user roles
const ResourceInstanaUserRole = "instana_user_role"

//ResourceInstanaApplicationConfig the name of the terraform-provider-instana resource to manage application config
const ResourceInstanaApplicationConfig = "instana_application_config"

//ResourceInstanaCustomEventSpecificationSystemRule the name of the terraform-provider-instana resource to manage custom event specifications with system rule
const ResourceInstanaCustomEventSpecificationSystemRule = "instana_custom_event_spec_system_rule"

//ResourceInstanaCustomEventSpecificationThresholdRule the name of the terraform-provider-instana resource to manage custom event specifications with threshold rule
const ResourceInstanaCustomEventSpecificationThresholdRule = "instana_custom_event_spec_threshold_rule"

//ProviderMeta data structure for the meta data which is configured and provided to the resources by this provider
type ProviderMeta struct {
	InstanaAPI              restapi.InstanaAPI
	ResourceStringFormatter ResourceStringFormatter
}

//Provider interface implementation of hashicorp terraform provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema:        providerSchema(),
		ResourcesMap:  providerResources(),
		ConfigureFunc: providerConfigure,
	}
}

func providerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		SchemaFieldAPIToken: &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "API token used to authenticate with the Instana Backend",
		},
		SchemaFieldEndpoint: &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "The DNS Name of the Instana Endpoint (eg. saas-eu-west-1.instana.io)",
		},
		SchemaFieldAddTerraformManagedString: &schema.Schema{
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "The indicator if a string should be appended to a name or description that the resource is managed by terraform",
		},
	}
}

func providerResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		ResourceInstanaUserRole:                              CreateResourceUserRole(),
		ResourceInstanaApplicationConfig:                     CreateResourceApplicationConfig(),
		ResourceInstanaCustomEventSpecificationSystemRule:    CreateResourceCustomEventSpecificationWithSystemRule(),
		ResourceInstanaCustomEventSpecificationThresholdRule: CreateResourceCustomEventSpecificationWithThresholdRule(),
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	apiToken := d.Get(SchemaFieldAPIToken).(string)
	endpoint := d.Get(SchemaFieldEndpoint).(string)
	appendTerraformManagedString := d.Get(SchemaFieldAddTerraformManagedString).(bool)
	instanaAPI := services.NewInstanaAPI(apiToken, endpoint)
	formatter := NewResourceStringFormatter(appendTerraformManagedString)
	return &ProviderMeta{
		InstanaAPI:              instanaAPI,
		ResourceStringFormatter: formatter,
	}, nil
}
