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
func NewAlertingChannelPagerDutyResourceHandle() *ResourceHandle {
	return &ResourceHandle{
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
		RestResourceFactory:  func(api restapi.InstanaAPI) restapi.RestResource { return api.AlertingChannels() },
		UpdateState:          updateStateForAlertingChannelPagerDuty,
		MapStateToDataObject: mapStateToDataObjectForAlertingChannelPagerDuty,
	}
}

func updateStateForAlertingChannelPagerDuty(d *schema.ResourceData, obj restapi.InstanaDataObject) error {
	alertingChannel := obj.(restapi.AlertingChannel)
	d.Set(AlertingChannelFieldFullName, alertingChannel.Name)
	d.Set(AlertingChannelPagerDutyFieldServiceIntegrationKey, alertingChannel.ServiceIntegrationKey)
	d.SetId(alertingChannel.ID)
	return nil
}

func mapStateToDataObjectForAlertingChannelPagerDuty(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	name := computeFullAlertingChannelNameString(d, formatter)
	serviceIntegrationKey := d.Get(AlertingChannelPagerDutyFieldServiceIntegrationKey).(string)
	return restapi.AlertingChannel{
		ID:                    d.Id(),
		Name:                  name,
		Kind:                  restapi.PagerDutyChannelType,
		ServiceIntegrationKey: &serviceIntegrationKey,
	}, nil
}
