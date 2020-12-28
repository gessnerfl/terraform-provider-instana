package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform/helper/schema"
)

//ResourceInstanaWebsiteMonitoring the name of the terraform-provider-instana resource to manage website monitoring configurations
const ResourceInstanaWebsiteMonitoring = "instana_website_monitoring_config"

const (
	//WebsiteMonitoringFieldName constant value for the schema field name
	WebsiteMonitoringFieldName = "name"
	//WebsiteMonitoringFieldFullName constant value for the schema field full_name
	WebsiteMonitoringFieldFullName = "full_name"
	//WebsiteMonitoringFieldAppName constant value for the schema field app_name
	WebsiteMonitoringFieldAppName = "app_name"
)

//WebsiteMonitoringSchemaName schema field definition of instana_website_monitoring_config field name
var WebsiteMonitoringSchemaName = &schema.Schema{
	Type:        schema.TypeString,
	Required:    true,
	Description: "Configures the name of the website monitoring configuration",
}

//WebsiteMonitoringSchemaFullName schema field definition of instana_website_monitoring_config field full_name
var WebsiteMonitoringSchemaFullName = &schema.Schema{
	Type:        schema.TypeString,
	Required:    false,
	Computed:    true,
	Description: "Configures the full name field of the website monitoring configuration. The field is computed and contains the name which is sent to instana. The computation depends on the configured default_name_prefix and default_name_suffix at provider level",
}

//WebsiteMonitoringSchemaAppName schema field definition of instana_website_monitoring_config field app_name
var WebsiteMonitoringSchemaAppName = &schema.Schema{
	Type:        schema.TypeString,
	Required:    false,
	Computed:    true,
	Description: "Configures the calculated app name of the website monitoring configuration",
}

//NewWebsiteMonitoringResourceHandle creates the resource handle for Alerting Configuration
func NewWebsiteMonitoringResourceHandle() *ResourceHandle {
	return &ResourceHandle{
		ResourceName: ResourceInstanaWebsiteMonitoring,
		Schema: map[string]*schema.Schema{
			WebsiteMonitoringFieldName:     WebsiteMonitoringSchemaName,
			WebsiteMonitoringFieldFullName: WebsiteMonitoringSchemaFullName,
			WebsiteMonitoringFieldAppName:  WebsiteMonitoringSchemaAppName,
		},
		SchemaVersion:        1,
		RestResourceFactory:  func(api restapi.InstanaAPI) restapi.RestResource { return api.WebsiteMonitoringConfig() },
		UpdateState:          updateStateForWebsiteMonitoring,
		MapStateToDataObject: mapStateToDataObjectForWebsiteMonitoring,
	}
}

func updateStateForWebsiteMonitoring(d *schema.ResourceData, obj restapi.InstanaDataObject) error {
	config := obj.(restapi.WebsiteMonitoringConfig)
	d.Set(WebsiteMonitoringFieldFullName, config.Name)
	d.Set(WebsiteMonitoringFieldAppName, config.AppName)
	d.SetId(config.ID)
	return nil
}

func mapStateToDataObjectForWebsiteMonitoring(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	name := computeFullWebsiteMonitoringNameString(d, formatter)

	return restapi.WebsiteMonitoringConfig{
		ID:   d.Id(),
		Name: name,
	}, nil
}

func computeFullWebsiteMonitoringNameString(d *schema.ResourceData, formatter utils.ResourceNameFormatter) string {
	if d.HasChange(WebsiteMonitoringFieldName) {
		return formatter.Format(d.Get(WebsiteMonitoringFieldName).(string))
	}
	return d.Get(WebsiteMonitoringFieldFullName).(string)
}
