package instana

import (
	"context"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceInstanaWebsiteMonitoringConfig the name of the terraform-provider-instana resource to manage website monitoring configurations
const ResourceInstanaWebsiteMonitoringConfig = "instana_website_monitoring_config"

const (
	//WebsiteMonitoringConfigFieldName constant value for the schema field name
	WebsiteMonitoringConfigFieldName = "name"
	//WebsiteMonitoringConfigFieldFullName constant value for the schema field full_name
	//Deprecated: not supported anymore with version 2.0
	WebsiteMonitoringConfigFieldFullName = "full_name"
	//WebsiteMonitoringConfigFieldAppName constant value for the schema field app_name
	WebsiteMonitoringConfigFieldAppName = "app_name"
)

// WebsiteMonitoringConfigSchemaName schema field definition of instana_website_monitoring_config field name
var WebsiteMonitoringConfigSchemaName = &schema.Schema{
	Type:        schema.TypeString,
	Required:    true,
	Description: "Configures the name of the website monitoring configuration",
}

// WebsiteMonitoringConfigSchemaFullName schema field definition of instana_website_monitoring_config field full_name
var WebsiteMonitoringConfigSchemaFullName = &schema.Schema{
	Type:        schema.TypeString,
	Required:    false,
	Computed:    true,
	Description: "Configures the full name field of the website monitoring configuration. The field is computed and contains the name which is sent to instana. The computation depends on the configured default_name_prefix and default_name_suffix at provider level",
}

// WebsiteMonitoringConfigSchemaAppName schema field definition of instana_website_monitoring_config field app_name
var WebsiteMonitoringConfigSchemaAppName = &schema.Schema{
	Type:        schema.TypeString,
	Required:    false,
	Computed:    true,
	Description: "Configures the calculated app name of the website monitoring configuration",
}

// NewWebsiteMonitoringConfigResourceHandle creates the resource handle for Alerting Configuration
func NewWebsiteMonitoringConfigResourceHandle() ResourceHandle[*restapi.WebsiteMonitoringConfig] {
	return &websiteMonitoringConfigResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaWebsiteMonitoringConfig,
			Schema: map[string]*schema.Schema{
				WebsiteMonitoringConfigFieldName:    WebsiteMonitoringConfigSchemaName,
				WebsiteMonitoringConfigFieldAppName: WebsiteMonitoringConfigSchemaAppName,
			},
			SchemaVersion: 1,
		},
	}
}

type websiteMonitoringConfigResource struct {
	metaData ResourceMetaData
}

func (r *websiteMonitoringConfigResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *websiteMonitoringConfigResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{
		{
			Type:    r.websiteMonitoringConfigSchemaV0().CoreConfigSchema().ImpliedType(),
			Upgrade: r.websiteMonitoringConfigStateUpgradeV0,
			Version: 0,
		},
	}
}

func (r *websiteMonitoringConfigResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.WebsiteMonitoringConfig] {
	return api.WebsiteMonitoringConfig()
}

func (r *websiteMonitoringConfigResource) SetComputedFields(_ *schema.ResourceData) error {
	return nil
}

func (r *websiteMonitoringConfigResource) UpdateState(d *schema.ResourceData, config *restapi.WebsiteMonitoringConfig) error {
	d.SetId(config.ID)
	return tfutils.UpdateState(d, map[string]interface{}{
		WebsiteMonitoringConfigFieldName:    config.Name,
		WebsiteMonitoringConfigFieldAppName: config.AppName,
	})
}

func (r *websiteMonitoringConfigResource) MapStateToDataObject(d *schema.ResourceData) (*restapi.WebsiteMonitoringConfig, error) {
	return &restapi.WebsiteMonitoringConfig{
		ID:   d.Id(),
		Name: d.Get(WebsiteMonitoringConfigFieldName).(string),
	}, nil
}

func (r *websiteMonitoringConfigResource) websiteMonitoringConfigStateUpgradeV0(_ context.Context, state map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	if _, ok := state[WebsiteMonitoringConfigFieldFullName]; ok {
		state[WebsiteMonitoringConfigFieldName] = state[WebsiteMonitoringConfigFieldFullName]
		delete(state, WebsiteMonitoringConfigFieldFullName)
	}
	return state, nil
}

func (r *websiteMonitoringConfigResource) websiteMonitoringConfigSchemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			WebsiteMonitoringConfigFieldName:     WebsiteMonitoringConfigSchemaName,
			WebsiteMonitoringConfigFieldFullName: WebsiteMonitoringConfigSchemaFullName,
			WebsiteMonitoringConfigFieldAppName:  WebsiteMonitoringConfigSchemaAppName,
		},
	}
}
