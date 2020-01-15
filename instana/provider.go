package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi/services"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform/helper/schema"
)

//SchemaFieldAPIToken the name of the provider configuration option for the api token
const SchemaFieldAPIToken = "api_token"

//SchemaFieldEndpoint the name of the provider configuration option for the instana endpoint
const SchemaFieldEndpoint = "endpoint"

//SchemaFieldDefaultNamePrefix the default prefix which should be added to all resource names/labels
const SchemaFieldDefaultNamePrefix = "default_name_prefix"

//SchemaFieldDefaultNameSuffix the default prefix which should be added to all resource names/labels
const SchemaFieldDefaultNameSuffix = "default_name_suffix"

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

//ResourceInstanaCustomEventSpecificationEntityVerificationRule the name of the terraform-provider-instana resource to manage custom event specifications with entity verification rule
const ResourceInstanaCustomEventSpecificationEntityVerificationRule = "instana_custom_event_spec_entity_verification_rule"

//ResourceInstanaAlertingChannelEmail the name of the terraform-provider-instana resource to manage alerting channels of type email
const ResourceInstanaAlertingChannelEmail = "instana_alerting_channel_email"

//ResourceInstanaAlertingChannelGoogleChat the name of the terraform-provider-instana resource to manage alerting channels of type Google Chat
const ResourceInstanaAlertingChannelGoogleChat = "instana_alerting_channel_google_chat"

//ResourceInstanaAlertingChannelSlack the name of the terraform-provider-instana resource to manage alerting channels of type Slack
const ResourceInstanaAlertingChannelSlack = "instana_alerting_channel_slack"

//ResourceInstanaAlertingChannelOffice365 the name of the terraform-provider-instana resource to manage alerting channels of type Office 365
const ResourceInstanaAlertingChannelOffice365 = "instana_alerting_channel_office_365"

//ResourceInstanaAlertingChannelOpsGenie the name of the terraform-provider-instana resource to manage alerting channels of type OpsGenie
const ResourceInstanaAlertingChannelOpsGenie = "instana_alerting_channel_ops_genie"

//ResourceInstanaAlertingChannelPagerDuty the name of the terraform-provider-instana resource to manage alerting channels of type PagerDuty
const ResourceInstanaAlertingChannelPagerDuty = "instana_alerting_channel_pager_duty"

//ResourceInstanaAlertingChannelSplunk the name of the terraform-provider-instana resource to manage alerting channels of type Splunk
const ResourceInstanaAlertingChannelSplunk = "instana_alerting_channel_splunk"

//ResourceInstanaAlertingChannelVictorOps the name of the terraform-provider-instana resource to manage alerting channels of type VictorOps
const ResourceInstanaAlertingChannelVictorOps = "instana_alerting_channel_victor_ops"

//ProviderMeta data structure for the meta data which is configured and provided to the resources by this provider
type ProviderMeta struct {
	InstanaAPI            restapi.InstanaAPI
	ResourceNameFormatter utils.ResourceNameFormatter
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
			Sensitive:   true,
			Required:    true,
			Description: "API token used to authenticate with the Instana Backend",
		},
		SchemaFieldEndpoint: &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "The DNS Name of the Instana Endpoint (eg. saas-eu-west-1.instana.io)",
		},
		SchemaFieldDefaultNamePrefix: &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "The default prefix which should be added to all resource names/labels",
		},
		SchemaFieldDefaultNameSuffix: &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "(TF managed)",
			Description: "The default suffix which should be added to all resource names/labels - default '(TF managed)'",
		},
	}
}

func providerResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		ResourceInstanaUserRole:                                       CreateResourceUserRole(),
		ResourceInstanaApplicationConfig:                              CreateResourceApplicationConfig(),
		ResourceInstanaCustomEventSpecificationSystemRule:             CreateResourceCustomEventSpecificationWithSystemRule(),
		ResourceInstanaCustomEventSpecificationThresholdRule:          CreateResourceCustomEventSpecificationWithThresholdRule(),
		ResourceInstanaCustomEventSpecificationEntityVerificationRule: CreateResourceCustomEventSpecificationWithEntityVerificationRule(),
		ResourceInstanaAlertingChannelEmail:                           NewAlertingChannelEmailResource().ToSchemaResource(),
		ResourceInstanaAlertingChannelGoogleChat:                      NewAlertingChannelGoogleChatResource().ToSchemaResource(),
		ResourceInstanaAlertingChannelOffice365:                       NewAlertingChannelOffice356Resource().ToSchemaResource(),
		ResourceInstanaAlertingChannelSlack:                           NewAlertingChannelSlackResource().ToSchemaResource(),
		ResourceInstanaAlertingChannelOpsGenie:                        NewAlertingChannelOpsGenieResource().ToSchemaResource(),
		ResourceInstanaAlertingChannelPagerDuty:                       NewAlertingChannelPagerDutyResource().ToSchemaResource(),
		ResourceInstanaAlertingChannelSplunk:                          NewAlertingChannelSplunkResource().ToSchemaResource(),
		ResourceInstanaAlertingChannelVictorOps:                       NewAlertingChannelVictorOpsResource().ToSchemaResource(),
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	apiToken := d.Get(SchemaFieldAPIToken).(string)
	endpoint := d.Get(SchemaFieldEndpoint).(string)
	defaultNamePrefix := d.Get(SchemaFieldDefaultNamePrefix).(string)
	defaultNameSuffix := d.Get(SchemaFieldDefaultNameSuffix).(string)
	instanaAPI := services.NewInstanaAPI(apiToken, endpoint)
	formatter := utils.NewResourceNameFormatter(defaultNamePrefix, defaultNameSuffix)
	return &ProviderMeta{
		InstanaAPI:            instanaAPI,
		ResourceNameFormatter: formatter,
	}, nil
}
