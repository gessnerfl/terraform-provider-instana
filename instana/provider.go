package instana

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SchemaFieldAPIToken the name of the provider configuration option for the api token
const SchemaFieldAPIToken = "api_token"

// SchemaFieldEndpoint the name of the provider configuration option for the instana endpoint
const SchemaFieldEndpoint = "endpoint"

// SchemaFieldTlsSkipVerify flag to deactivate skip tls verification
const SchemaFieldTlsSkipVerify = "tls_skip_verify"

// ProviderMeta data structure for the metadata which is configured and provided to the resources by this provider
type ProviderMeta struct {
	InstanaAPI restapi.InstanaAPI
}

// Provider interface implementation of hashicorp terraform provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema:               providerSchema(),
		ResourcesMap:         providerResources(),
		DataSourcesMap:       providerDataSources(),
		ConfigureContextFunc: providerConfigure,
	}
}

func providerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		SchemaFieldAPIToken: {
			Type:        schema.TypeString,
			Sensitive:   true,
			Required:    true,
			DefaultFunc: schema.EnvDefaultFunc("INSTANA_API_TOKEN", nil),
			Description: "API token used to authenticate with the Instana Backend",
			Deprecated:  "This project has been handed over to and is maintained under IBM's offical Instana org. Please use the official IBM Instana Terraform provider instana/instana (https://registry.terraform.io/providers/instana/instana/latest/) instead",
		},
		SchemaFieldEndpoint: {
			Type:        schema.TypeString,
			Required:    true,
			DefaultFunc: schema.EnvDefaultFunc("INSTANA_ENDPOINT", nil),
			Description: "The DNS Name of the Instana Endpoint (eg. saas-eu-west-1.instana.io)",
		},
		SchemaFieldTlsSkipVerify: {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "If set to true, TLS verification will be skipped when calling Instana API",
		},
	}
}

func providerResources() map[string]*schema.Resource {
	resources := make(map[string]*schema.Resource)
	bindResourceHandle(resources, NewAPITokenResourceHandle())
	bindResourceHandle(resources, NewApplicationConfigResourceHandle())
	bindResourceHandle(resources, NewApplicationAlertConfigResourceHandle())
	bindResourceHandle(resources, NewGlobalApplicationAlertConfigResourceHandle())
	bindResourceHandle(resources, NewCustomEventSpecificationResourceHandle())
	bindResourceHandle(resources, NewAlertingChannelResourceHandle())
	bindResourceHandle(resources, NewAlertingConfigResourceHandle())
	bindResourceHandle(resources, NewSliConfigResourceHandle())
	bindResourceHandle(resources, NewWebsiteMonitoringConfigResourceHandle())
	bindResourceHandle(resources, NewWebsiteAlertConfigResourceHandle())
	bindResourceHandle(resources, NewGroupResourceHandle())
	bindResourceHandle(resources, NewCustomDashboardResourceHandle())
	bindResourceHandle(resources, NewSyntheticTestResourceHandle())
	return resources
}

func bindResourceHandle[T restapi.InstanaDataObject](resources map[string]*schema.Resource, resourceHandle ResourceHandle[T]) {
	resources[resourceHandle.MetaData().ResourceName] = NewTerraformResource(resourceHandle).ToSchemaResource()
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiToken := strings.TrimSpace(d.Get(SchemaFieldAPIToken).(string))
	endpoint := strings.TrimSpace(d.Get(SchemaFieldEndpoint).(string))
	skipTlsVerify := d.Get(SchemaFieldTlsSkipVerify).(bool)
	instanaAPI := restapi.NewInstanaAPI(apiToken, endpoint, skipTlsVerify)
	return &ProviderMeta{
		InstanaAPI: instanaAPI,
	}, nil
}

func providerDataSources() map[string]*schema.Resource {
	dataSources := make(map[string]*schema.Resource)
	dataSources[DataSourceBuiltinEvent] = NewBuiltinEventDataSource().CreateResource()
	dataSources[DataSourceSyntheticLocation] = NewSyntheticLocationDataSource().CreateResource()
	dataSources[DataSourceAlertingChannel] = NewAlertingChannelDataSource().CreateResource()
	return dataSources
}
