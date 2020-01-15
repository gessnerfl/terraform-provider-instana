package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	//AlertingChannelPagerDutyFieldServiceIntegrationKey const for the emails field of the alerting channel
	AlertingChannelPagerDutyFieldServiceIntegrationKey = "service_integration_key"
)

//NewAlertingChannelPagerDutyResource creates the terraform resource for Alerting Channels of type PagerDuty
func NewAlertingChannelPagerDutyResource() TerraformResource {
	return NewTerraformResource(NewAlertingChannelPagerDutyResourceHandle())
}

//NewAlertingChannelPagerDutyResourceHandle creates the resource handle for Alerting Channels of type PagerDuty
func NewAlertingChannelPagerDutyResourceHandle() ResourceHandle {
	return &alertingChannelPagerDutyResourceHandle{}
}

type alertingChannelPagerDutyResourceHandle struct {
}

func (h *alertingChannelPagerDutyResourceHandle) GetResource(api restapi.InstanaAPI) restapi.RestResource {
	return api.AlertingChannels()
}

func (h *alertingChannelPagerDutyResourceHandle) GetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		AlertingChannelFieldName:     alertingChannelNameSchemaField,
		AlertingChannelFieldFullName: alertingChannelFullNameSchemaField,
		AlertingChannelPagerDutyFieldServiceIntegrationKey: &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "The Service Integration Key of the PagerDuty alerting channel",
		},
	}
}

func (h *alertingChannelPagerDutyResourceHandle) GetResourceName() string {
	return "instana_alerting_channel_pager_duty"
}

func (h *alertingChannelPagerDutyResourceHandle) UpdateState(d *schema.ResourceData, obj restapi.InstanaDataObject) {
	alertingChannel := obj.(restapi.AlertingChannel)
	d.Set(AlertingChannelFieldFullName, alertingChannel.Name)
	d.Set(AlertingChannelPagerDutyFieldServiceIntegrationKey, alertingChannel.ServiceIntegrationKey)
	d.SetId(alertingChannel.ID)
}

func (h *alertingChannelPagerDutyResourceHandle) ConvertStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) restapi.InstanaDataObject {
	name := computeFullAlertingChannelNameString(d, formatter)
	serviceIntegrationKey := d.Get(AlertingChannelPagerDutyFieldServiceIntegrationKey).(string)
	return restapi.AlertingChannel{
		ID:                    d.Id(),
		Name:                  name,
		Kind:                  restapi.PagerDutyChannelType,
		ServiceIntegrationKey: &serviceIntegrationKey,
	}

}
