package instana

import (
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform/helper/schema"
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
	return &alertingChannelWebhookBasedResourceHandle{channelType: channelType, resourceName: resourceName}
}

type alertingChannelWebhookBasedResourceHandle struct {
	channelType  restapi.AlertingChannelType
	resourceName string
}

func (h *alertingChannelWebhookBasedResourceHandle) GetResourceFrom(api restapi.InstanaAPI) restapi.RestResource {
	return api.AlertingChannels()
}

func (h *alertingChannelWebhookBasedResourceHandle) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		AlertingChannelFieldName:     alertingChannelNameSchemaField,
		AlertingChannelFieldFullName: alertingChannelFullNameSchemaField,
		AlertingChannelWebhookBasedFieldWebhookURL: &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: fmt.Sprintf("The webhook URL of the %s alerting channel", h.channelType),
		},
	}
}

func (h *alertingChannelWebhookBasedResourceHandle) SchemaVersion() int {
	return 0
}

func (h *alertingChannelWebhookBasedResourceHandle) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{}
}

func (h *alertingChannelWebhookBasedResourceHandle) ResourceName() string {
	return h.resourceName
}

func (h *alertingChannelWebhookBasedResourceHandle) UpdateState(d *schema.ResourceData, obj restapi.InstanaDataObject) error {
	alertingChannel := obj.(restapi.AlertingChannel)
	d.Set(AlertingChannelFieldFullName, alertingChannel.Name)
	d.Set(AlertingChannelWebhookBasedFieldWebhookURL, alertingChannel.WebhookURL)
	d.SetId(alertingChannel.ID)
	return nil
}

func (h *alertingChannelWebhookBasedResourceHandle) ConvertStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	name := computeFullAlertingChannelNameString(d, formatter)
	webhookURL := d.Get(AlertingChannelWebhookBasedFieldWebhookURL).(string)
	return restapi.AlertingChannel{
		ID:         d.Id(),
		Name:       name,
		Kind:       h.channelType,
		WebhookURL: &webhookURL,
	}, nil
}
