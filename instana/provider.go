package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//SchemaFieldAPIToken the name of the provider configuration option for the api token
const SchemaFieldAPIToken = "api_token"

//SchemaFieldEndpoint the name of the provider configuration option for the instana endpoint
const SchemaFieldEndpoint = "endpoint"

//SchemaFieldDefaultNamePrefix the default prefix which should be added to all resource names/labels
const SchemaFieldDefaultNamePrefix = "default_name_prefix"

//SchemaFieldDefaultNameSuffix the default prefix which should be added to all resource names/labels
const SchemaFieldDefaultNameSuffix = "default_name_suffix"

//ProviderMeta data structure for the meta data which is configured and provided to the resources by this provider
type ProviderMeta struct {
	InstanaAPI            restapi.InstanaAPI
	ResourceNameFormatter utils.ResourceNameFormatter
}

//Provider interface implementation of hashicorp terraform provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema:         providerSchema(),
		ResourcesMap:   providerResources(),
		DataSourcesMap: providerDataSources(),
		ConfigureFunc:  providerConfigure,
	}
}

func providerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		SchemaFieldAPIToken: {
			Type:        schema.TypeString,
			Sensitive:   true,
			Required:    true,
			Description: "API token used to authenticate with the Instana Backend",
		},
		SchemaFieldEndpoint: {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The DNS Name of the Instana Endpoint (eg. saas-eu-west-1.instana.io)",
		},
		SchemaFieldDefaultNamePrefix: {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "The default prefix which should be added to all resource names/labels",
		},
		SchemaFieldDefaultNameSuffix: {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "(TF managed)",
			Description: "The default suffix which should be added to all resource names/labels - default '(TF managed)'",
		},
	}
}

func providerResources() map[string]*schema.Resource {
	resources := make(map[string]*schema.Resource)
	bindResourceHandle(resources, NewUserRoleResourceHandle())
	bindResourceHandle(resources, NewApplicationConfigResourceHandle())
	bindResourceHandle(resources, NewCustomEventSpecificationWithSystemRuleResourceHandle())
	bindResourceHandle(resources, NewCustomEventSpecificationWithThresholdRuleResourceHandle())
	bindResourceHandle(resources, NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle())
	bindResourceHandle(resources, NewAlertingChannelEmailResourceHandle())
	bindResourceHandle(resources, NewAlertingChannelGoogleChatResourceHandle())
	bindResourceHandle(resources, NewAlertingChannelOffice356ResourceHandle())
	bindResourceHandle(resources, NewAlertingChannelSlackResourceHandle())
	bindResourceHandle(resources, NewAlertingChannelOpsGenieResourceHandle())
	bindResourceHandle(resources, NewAlertingChannelPagerDutyResourceHandle())
	bindResourceHandle(resources, NewAlertingChannelSplunkResourceHandle())
	bindResourceHandle(resources, NewAlertingChannelVictorOpsResourceHandle())
	bindResourceHandle(resources, NewAlertingChannelWebhookResourceHandle())
	bindResourceHandle(resources, NewAlertingConfigResourceHandle())
	bindResourceHandle(resources, NewSliConfigResourceHandle())
	bindResourceHandle(resources, NewWebsiteMonitoringConfigResourceHandle())
	return resources
}

func bindResourceHandle(resources map[string]*schema.Resource, resourceHandle ResourceHandle) {
	resources[resourceHandle.MetaData().ResourceName] = NewTerraformResource(resourceHandle).ToSchemaResource()
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	apiToken := d.Get(SchemaFieldAPIToken).(string)
	endpoint := d.Get(SchemaFieldEndpoint).(string)
	defaultNamePrefix := d.Get(SchemaFieldDefaultNamePrefix).(string)
	defaultNameSuffix := d.Get(SchemaFieldDefaultNameSuffix).(string)
	instanaAPI := restapi.NewInstanaAPI(apiToken, endpoint)
	formatter := utils.NewResourceNameFormatter(defaultNamePrefix, defaultNameSuffix)
	return &ProviderMeta{
		InstanaAPI:            instanaAPI,
		ResourceNameFormatter: formatter,
	}, nil
}

func providerDataSources() map[string]*schema.Resource {
	dataSources := make(map[string]*schema.Resource)
	dataSources[DataSourceBuiltinEvent] = NewBuiltinEventDataSource().CreateResource()
	return dataSources
}
