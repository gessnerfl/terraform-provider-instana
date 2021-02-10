package instana

import (
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	//AlertingChannelWebhookBasedFieldWebhookURL const for the webhookUrl field of the alerting channel
	AlertingChannelWebhookBasedFieldWebhookURL = "webhook_url"
	//ResourceInstanaAlertingChannelOffice365 the name of the terraform-provider-instana resource to manage alerting channels of type Office 365
	ResourceInstanaAlertingChannelOffice365 = "instana_alerting_channel_office_365"
	//ResourceInstanaAlertingChannelGoogleChat the name of the terraform-provider-instana resource to manage alerting channels of type Google Chat
	ResourceInstanaAlertingChannelGoogleChat = "instana_alerting_channel_google_chat"
)

//NewAlertingChannelGoogleChatResourceHandle creates the terraform resource for Alerting Channels of type Google Chat
func NewAlertingChannelGoogleChatResourceHandle() ResourceHandle {
	return newAlertingChannelWebhookBasedResourceHandle(restapi.GoogleChatChannelType, ResourceInstanaAlertingChannelGoogleChat)
}

//NewAlertingChannelOffice356ResourceHandle creates the terraform resource for Alerting Channels of type Office 356
func NewAlertingChannelOffice356ResourceHandle() ResourceHandle {
	return newAlertingChannelWebhookBasedResourceHandle(restapi.Office365ChannelType, ResourceInstanaAlertingChannelOffice365)
}

func newAlertingChannelWebhookBasedResourceHandle(channelType restapi.AlertingChannelType, resourceName string) ResourceHandle {
	return &alertingChannelWebhookBasedResource{
		metaData: ResourceMetaData{
			ResourceName: resourceName,
			Schema: map[string]*schema.Schema{
				AlertingChannelFieldName:     alertingChannelNameSchemaField,
				AlertingChannelFieldFullName: alertingChannelFullNameSchemaField,
				AlertingChannelWebhookBasedFieldWebhookURL: {
					Type:        schema.TypeString,
					Required:    true,
					Description: fmt.Sprintf("The webhook URL of the %s alerting channel", channelType),
				},
			},
		},
		channelType: channelType,
	}
}

type alertingChannelWebhookBasedResource struct {
	metaData    ResourceMetaData
	channelType restapi.AlertingChannelType
}

func (r *alertingChannelWebhookBasedResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *alertingChannelWebhookBasedResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{}
}

func (r *alertingChannelWebhookBasedResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource {
	return api.AlertingChannels()
}

func (r *alertingChannelWebhookBasedResource) SetComputedFields(d *schema.ResourceData) {
	//No computed fields defined
}

func (r *alertingChannelWebhookBasedResource) UpdateState(d *schema.ResourceData, obj restapi.InstanaDataObject) error {
	alertingChannel := obj.(*restapi.AlertingChannel)
	d.Set(AlertingChannelFieldFullName, alertingChannel.Name)
	d.Set(AlertingChannelWebhookBasedFieldWebhookURL, alertingChannel.WebhookURL)
	d.SetId(alertingChannel.ID)
	return nil
}

func (r *alertingChannelWebhookBasedResource) MapStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	name := computeFullAlertingChannelNameString(d, formatter)
	webhookURL := d.Get(AlertingChannelWebhookBasedFieldWebhookURL).(string)
	return &restapi.AlertingChannel{
		ID:         d.Id(),
		Name:       name,
		Kind:       r.channelType,
		WebhookURL: &webhookURL,
	}, nil
}
