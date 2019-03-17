package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi/services"
	"github.com/hashicorp/terraform/helper/schema"
)

//SchemaFieldAPIToken the name of the provider configuration option for the api token
const SchemaFieldAPIToken = "api_token"

//SchemaFieldEndpoint the name of the provider configuration option for the instana endpoint
const SchemaFieldEndpoint = "endpoint"

//SchemaFieldVerifyServerCertificate activates/deactivates TLS server certificate validation
const SchemaFieldVerifyServerCertificate = "verify_server_certificate"

//ResourceInstanaRule the name of the terraform-provider-instana resource to manage rules
const ResourceInstanaRule = "instana_rule"

//ResourceInstanaRuleBinding the name of the terraform-provider-instana resource to manage rule bindings
const ResourceInstanaRuleBinding = "instana_rule_binding"

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
		SchemaFieldVerifyServerCertificate: &schema.Schema{
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Verify server certificate",
		},
	}
}

func providerResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		ResourceInstanaRule:        CreateResourceRule(),
		ResourceInstanaRuleBinding: createResourceRuleBinding(),
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	apiToken := d.Get(SchemaFieldAPIToken).(string)
	endpoint := d.Get(SchemaFieldEndpoint).(string)
	validateServerCertificate := d.Get(SchemaFieldVerifyServerCertificate).(bool)
	instanaAPI := services.NewInstanaAPI(apiToken, endpoint, validateServerCertificate)
	return instanaAPI, nil
}
