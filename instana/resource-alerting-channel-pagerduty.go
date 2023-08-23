package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	//AlertingChannelPagerDutyFieldServiceIntegrationKey const for the emails field of the alerting channel
	AlertingChannelPagerDutyFieldServiceIntegrationKey = "service_integration_key"
	//ResourceInstanaAlertingChannelPagerDuty the name of the terraform-provider-instana resource to manage alerting channels of type PagerDuty
	ResourceInstanaAlertingChannelPagerDuty = "instana_alerting_channel_pager_duty"
)

var alertingChannelPagerDutySchemaServiceIntegrationKey = &schema.Schema{
	Type:        schema.TypeString,
	Required:    true,
	Description: "The Service Integration Key of the PagerDuty alerting channel",
}

// NewAlertingChannelPagerDutyResourceHandle creates the resource handle for Alerting Channels of type PagerDuty
func NewAlertingChannelPagerDutyResourceHandle() ResourceHandle[*restapi.AlertingChannel] {
	return &alertingChannelPagerDutyResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaAlertingChannelPagerDuty,
			Schema: map[string]*schema.Schema{
				AlertingChannelFieldName:                           alertingChannelNameSchemaField,
				AlertingChannelPagerDutyFieldServiceIntegrationKey: alertingChannelPagerDutySchemaServiceIntegrationKey,
			},
			SchemaVersion: 1,
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
	return []schema.StateUpgrader{
		{
			Type:    r.schemaV0().CoreConfigSchema().ImpliedType(),
			Upgrade: migrateFullNameToName,
			Version: 0,
		},
	}
}

func (r *alertingChannelPagerDutyResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.AlertingChannel] {
	return api.AlertingChannels()
}

func (r *alertingChannelPagerDutyResource) SetComputedFields(_ *schema.ResourceData) error {
	return nil
}

func (r *alertingChannelPagerDutyResource) UpdateState(d *schema.ResourceData, alertingChannel *restapi.AlertingChannel) error {
	d.SetId(alertingChannel.ID)
	return tfutils.UpdateState(d, map[string]interface{}{
		AlertingChannelFieldName:                           alertingChannel.Name,
		AlertingChannelPagerDutyFieldServiceIntegrationKey: alertingChannel.ServiceIntegrationKey,
	})
}

func (r *alertingChannelPagerDutyResource) MapStateToDataObject(d *schema.ResourceData) (*restapi.AlertingChannel, error) {
	serviceIntegrationKey := d.Get(AlertingChannelPagerDutyFieldServiceIntegrationKey).(string)
	return &restapi.AlertingChannel{
		ID:                    d.Id(),
		Name:                  d.Get(AlertingChannelFieldName).(string),
		Kind:                  restapi.PagerDutyChannelType,
		ServiceIntegrationKey: &serviceIntegrationKey,
	}, nil
}

func (r *alertingChannelPagerDutyResource) schemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			AlertingChannelFieldName:                           alertingChannelNameSchemaField,
			AlertingChannelFieldFullName:                       alertingChannelFullNameSchemaField,
			AlertingChannelPagerDutyFieldServiceIntegrationKey: alertingChannelPagerDutySchemaServiceIntegrationKey,
		},
	}
}
