package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//ResourceInstanaApplicationAlertConfig the name of the terraform-provider-instana resource to manage application alert config
const ResourceInstanaApplicationAlertConfig = "instana_application_alert_config"

//NewApplicationAlertConfigResourceHandle creates a new instance of the ResourceHandle for application alert configs
func NewApplicationAlertConfigResourceHandle() ResourceHandle {
	return &applicationAlertConfigResource{
		metaData: ResourceMetaData{
			ResourceName:  ResourceInstanaApplicationAlertConfig,
			Schema:        map[string]*schema.Schema{},
			SchemaVersion: 3,
		},
	}
}

type applicationAlertConfigResource struct {
	metaData ResourceMetaData
}

func (r *applicationAlertConfigResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *applicationAlertConfigResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{}
}

func (r *applicationAlertConfigResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource {
	return api.ApplicationConfigs()
}

func (r *applicationAlertConfigResource) SetComputedFields(d *schema.ResourceData) {
	//No computed fields defined
}

func (r *applicationAlertConfigResource) UpdateState(d *schema.ResourceData, obj restapi.InstanaDataObject, formatter utils.ResourceNameFormatter) error {
	return nil
}

func (r *applicationAlertConfigResource) MapStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	return nil, nil
}
