package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
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
func NewAlertingChannelPagerDutyResourceHandle() ResourceHandle {
	return &alertingChannelPagerDutyResource{
		metaData: ResourceMetaData{
			ResourceName:       ResourceInstanaAlertingChannelPagerDuty,
			DeprecationMessage: "This feature will be removed in version 2.x and should be replaced with instana_alerting_channel",
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

func (r *alertingChannelPagerDutyResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource {
	return api.AlertingChannels()
}

func (r *alertingChannelPagerDutyResource) SetComputedFields(d *schema.ResourceData) {
	//No computed fields defined
}

func (r *alertingChannelPagerDutyResource) UpdateState(d *schema.ResourceData, obj restapi.InstanaDataObject, formatter utils.ResourceNameFormatter) error {
	alertingChannel := obj.(*restapi.AlertingChannel)
	d.Set(AlertingChannelFieldName, formatter.UndoFormat(alertingChannel.Name))
	d.Set(AlertingChannelFieldFullName, alertingChannel.Name)
	d.Set(AlertingChannelPagerDutyFieldServiceIntegrationKey, alertingChannel.ServiceIntegrationKey)
	d.SetId(alertingChannel.ID)
	return nil
}

func (r *alertingChannelPagerDutyResource) MapStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	name := computeFullAlertingChannelNameString(d, formatter)
	serviceIntegrationKey := d.Get(AlertingChannelPagerDutyFieldServiceIntegrationKey).(string)
	return &restapi.AlertingChannel{
		ID:                    d.Id(),
		Name:                  name,
		Kind:                  restapi.PagerDutyChannelType,
		ServiceIntegrationKey: &serviceIntegrationKey,
	}, nil
}
