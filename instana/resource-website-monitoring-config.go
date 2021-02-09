package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform/helper/schema"
)

//ResourceInstanaWebsiteMonitoringConfig the name of the terraform-provider-instana resource to manage website monitoring configurations
const ResourceInstanaWebsiteMonitoringConfig = "instana_website_monitoring_config"

const (
	//WebsiteMonitoringConfigFieldName constant value for the schema field name
	WebsiteMonitoringConfigFieldName = "name"
	//WebsiteMonitoringConfigFieldFullName constant value for the schema field full_name
	WebsiteMonitoringConfigFieldFullName = "full_name"
	//WebsiteMonitoringConfigFieldAppName constant value for the schema field app_name
	WebsiteMonitoringConfigFieldAppName = "app_name"
)

//WebsiteMonitoringConfigSchemaName schema field definition of instana_website_monitoring_config field name
var WebsiteMonitoringConfigSchemaName = &schema.Schema{
	Type:        schema.TypeString,
	Required:    true,
	Description: "Configures the name of the website monitoring configuration",
}

//WebsiteMonitoringConfigSchemaFullName schema field definition of instana_website_monitoring_config field full_name
var WebsiteMonitoringConfigSchemaFullName = &schema.Schema{
	Type:        schema.TypeString,
	Required:    false,
	Computed:    true,
	Description: "Configures the full name field of the website monitoring configuration. The field is computed and contains the name which is sent to instana. The computation depends on the configured default_name_prefix and default_name_suffix at provider level",
}

//WebsiteMonitoringConfigSchemaAppName schema field definition of instana_website_monitoring_config field app_name
var WebsiteMonitoringConfigSchemaAppName = &schema.Schema{
	Type:        schema.TypeString,
	Required:    false,
	Computed:    true,
	Description: "Configures the calculated app name of the website monitoring configuration",
}

//NewWebsiteMonitoringConfigResourceHandle creates the resource handle for Alerting Configuration
func NewWebsiteMonitoringConfigResourceHandle() *ResourceHandle {
	return &ResourceHandle{
		ResourceName: ResourceInstanaWebsiteMonitoringConfig,
		Schema: map[string]*schema.Schema{
			WebsiteMonitoringConfigFieldName:     WebsiteMonitoringConfigSchemaName,
			WebsiteMonitoringConfigFieldFullName: WebsiteMonitoringConfigSchemaFullName,
			WebsiteMonitoringConfigFieldAppName:  WebsiteMonitoringConfigSchemaAppName,
		},
		RestResourceFactory:  func(api restapi.InstanaAPI) restapi.RestResource { return api.WebsiteMonitoringConfig() },
		UpdateState:          updateStateForWebsiteMonitoring,
		MapStateToDataObject: mapStateToDataObjectForWebsiteMonitoring,
	}
}

func updateStateForWebsiteMonitoring(d *schema.ResourceData, obj restapi.InstanaDataObject) error {
	config := obj.(*restapi.WebsiteMonitoringConfig)
	d.Set(WebsiteMonitoringConfigFieldFullName, config.Name)
	d.Set(WebsiteMonitoringConfigFieldAppName, config.AppName)
	d.SetId(config.ID)
	return nil
}

func mapStateToDataObjectForWebsiteMonitoring(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	name := computeFullWebsiteMonitoringNameString(d, formatter)

	return &restapi.WebsiteMonitoringConfig{
		ID:   d.Id(),
		Name: name,
	}, nil
}

func computeFullWebsiteMonitoringNameString(d *schema.ResourceData, formatter utils.ResourceNameFormatter) string {
	if d.HasChange(WebsiteMonitoringConfigFieldName) {
		return formatter.Format(d.Get(WebsiteMonitoringConfigFieldName).(string))
	}
	return d.Get(WebsiteMonitoringConfigFieldFullName).(string)
}
