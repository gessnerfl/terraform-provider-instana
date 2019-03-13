package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform/helper/schema"
)

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
		"api_key": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "API Key used to authenticate with the Instana Backend",
		},
		"endpoint": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "The DNS Name of the Instana Endpoint (eg. saas-eu-west-1.instana.io)",
		},
	}
}

func providerResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"instana_rule": createResourceRule(),
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := restapi.NewClient(d.Get("api_key").(string), d.Get("endpoint").(string))
	return &client, nil
}
