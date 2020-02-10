package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	//AlertingChannelPagerDutyFieldServiceIntegrationKey const for the emails field of the alerting channel
	AlertingChannelPagerDutyFieldServiceIntegrationKey = "service_integration_key"
	//ResourceInstanaAlertingChannelPagerDuty the name of the terraform-provider-instana resource to manage alerting channels of type PagerDuty
	ResourceInstanaAlertingChannelPagerDuty = "instana_alerting_channel_pager_duty"
)

//NewAlertingChannelPagerDutyResourceHandle creates the resource handle for Alerting Channels of type PagerDuty
func NewAlertingChannelPagerDutyResourceHandle() ResourceHandle {
	return &alertingChannelPagerDutyResourceHandle{}
}

type alertingChannelPagerDutyResourceHandle struct {
}

func (h *alertingChannelPagerDutyResourceHandle) GetResourceFrom(api restapi.InstanaAPI) restapi.RestResource {
	return api.AlertingChannels()
}

func (h *alertingChannelPagerDutyResourceHandle) Schema() map[string]*schema.Schema {
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

func (h *alertingChannelPagerDutyResourceHandle) SchemaVersion() int {
	return 0
}

func (h *alertingChannelPagerDutyResourceHandle) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{}
}

func (h *alertingChannelPagerDutyResourceHandle) ResourceName() string {
	return ResourceInstanaAlertingChannelPagerDuty
}

func (h *alertingChannelPagerDutyResourceHandle) UpdateState(d *schema.ResourceData, obj restapi.InstanaDataObject) error {
	alertingChannel := obj.(restapi.AlertingChannel)
	d.Set(AlertingChannelFieldFullName, alertingChannel.Name)
	d.Set(AlertingChannelPagerDutyFieldServiceIntegrationKey, alertingChannel.ServiceIntegrationKey)
	d.SetId(alertingChannel.ID)
	return nil
}

func (h *alertingChannelPagerDutyResourceHandle) ConvertStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	name := computeFullAlertingChannelNameString(d, formatter)
	serviceIntegrationKey := d.Get(AlertingChannelPagerDutyFieldServiceIntegrationKey).(string)
	return restapi.AlertingChannel{
		ID:                    d.Id(),
		Name:                  name,
		Kind:                  restapi.PagerDutyChannelType,
		ServiceIntegrationKey: &serviceIntegrationKey,
	}, nil
}
