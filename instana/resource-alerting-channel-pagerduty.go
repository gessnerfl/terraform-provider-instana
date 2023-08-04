package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	//AlertingChannelPagerDutyFieldServiceIntegrationKey const for the emails field of the alerting channel
	AlertingChannelPagerDutyFieldServiceIntegrationKey = "service_integration_key"
	//ResourceInstanaAlertingChannelPagerDuty the name of the terraform-provider-instana resource to manage alerting channels of type PagerDuty
	ResourceInstanaAlertingChannelPagerDuty = "instana_alerting_channel_pager_duty"
)

// NewAlertingChannelPagerDutyResourceHandle creates the resource handle for Alerting Channels of type PagerDuty
func NewAlertingChannelPagerDutyResourceHandle() ResourceHandle[*restapi.AlertingChannel] {
	return &alertingChannelPagerDutyResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaAlertingChannelPagerDuty,
			Schema: map[string]*schema.Schema{
				AlertingChannelFieldName:     alertingChannelNameSchemaField,
				AlertingChannelFieldFullName: alertingChannelFullNameSchemaField,
				AlertingChannelPagerDutyFieldServiceIntegrationKey: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The Service Integration Key of the PagerDuty alerting channel",
				},
			},
		},
	}
}

type alertingChannelPagerDutyResource struct {
	metaData ResourceMetaData
}

func (r *alertingChannelPagerDutyResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *alertingChannelPagerDutyResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{}
}

func (r *alertingChannelPagerDutyResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.AlertingChannel] {
	return api.AlertingChannels()
}

func (r *alertingChannelPagerDutyResource) SetComputedFields(_ *schema.ResourceData) error {
	return nil
}

func (r *alertingChannelPagerDutyResource) UpdateState(d *schema.ResourceData, alertingChannel *restapi.AlertingChannel, formatter utils.ResourceNameFormatter) error {
	d.SetId(alertingChannel.ID)
	return tfutils.UpdateState(d, map[string]interface{}{
		AlertingChannelFieldName:                           formatter.UndoFormat(alertingChannel.Name),
		AlertingChannelFieldFullName:                       alertingChannel.Name,
		AlertingChannelPagerDutyFieldServiceIntegrationKey: alertingChannel.ServiceIntegrationKey,
	})
}

func (r *alertingChannelPagerDutyResource) MapStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (*restapi.AlertingChannel, error) {
	name := computeFullAlertingChannelNameString(d, formatter)
	serviceIntegrationKey := d.Get(AlertingChannelPagerDutyFieldServiceIntegrationKey).(string)
	return &restapi.AlertingChannel{
		ID:                    d.Id(),
		Name:                  name,
		Kind:                  restapi.PagerDutyChannelType,
		ServiceIntegrationKey: &serviceIntegrationKey,
	}, nil
}
