package instana

import (
	"fmt"
	"strings"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	//AlertingChannelWebhookBasedFieldWebhookURL const for the webhookUrl field of the alerting channel
	AlertingChannelWebhookBasedFieldWebhookURL = "webhook_url"
)

//NewAlertingChannelGoogleChatResource creates the terraform resource for Alerting Channels of type Google Chat
func NewAlertingChannelGoogleChatResource() TerraformResource {
	return NewTerraformResource(NewAlertingChannelWebhookBasedResourceHandle(restapi.GoogleChatChannelType))
}

//NewAlertingChannelOffice356Resource creates the terraform resource for Alerting Channels of type Office 356
func NewAlertingChannelOffice356Resource() TerraformResource {
	return NewTerraformResource(NewAlertingChannelWebhookBasedResourceHandle(restapi.Office365ChannelType))
}

//NewAlertingChannelSlackResource creates the terraform resource for Alerting Channels of type Slack
func NewAlertingChannelSlackResource() TerraformResource {
	return NewTerraformResource(NewAlertingChannelWebhookBasedResourceHandle(restapi.SlackChannelType))
}

//NewAlertingChannelWebhookBasedResourceHandle creates the resource handle for Alerting Channels of type Email
func NewAlertingChannelWebhookBasedResourceHandle(channelType restapi.AlertingChannelType) ResourceHandle {
	return &alertingChannelWebhookBasedResourceHandle{channelType: channelType}
}

type alertingChannelWebhookBasedResourceHandle struct {
	channelType restapi.AlertingChannelType
}

func (h *alertingChannelWebhookBasedResourceHandle) GetResource(api restapi.InstanaAPI) restapi.RestResource {
	return api.AlertingChannels()
}

func (h *alertingChannelWebhookBasedResourceHandle) GetSchema() map[string]*schema.Schema {
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

func (h *alertingChannelWebhookBasedResourceHandle) GetResourceName() string {
	return "instana_alerting_channel_" + strings.ToLower(string(h.channelType))
}

func (h *alertingChannelWebhookBasedResourceHandle) UpdateState(d *schema.ResourceData, obj restapi.InstanaDataObject) {
	alertingChannel := obj.(restapi.AlertingChannel)
	d.Set(AlertingChannelFieldFullName, alertingChannel.Name)
	d.Set(AlertingChannelWebhookBasedFieldWebhookURL, alertingChannel.WebhookURL)
	d.SetId(alertingChannel.ID)
}

func (h *alertingChannelWebhookBasedResourceHandle) ConvertStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) restapi.InstanaDataObject {
	name := computeFullAlertingChannelNameString(d, formatter)
	webhookURL := d.Get(AlertingChannelWebhookBasedFieldWebhookURL).(string)
	return restapi.AlertingChannel{
		ID:         d.Id(),
		Name:       name,
		Kind:       h.channelType,
		WebhookURL: &webhookURL,
	}
}
