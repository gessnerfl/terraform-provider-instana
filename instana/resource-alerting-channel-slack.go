package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	//AlertingChannelSlackFieldWebhookURL const for the webhookUrl field of the Slack alerting channel
	AlertingChannelSlackFieldWebhookURL = "webhook_url"
	//AlertingChannelSlackFieldIconURL const for the iconURL field of the Slack alerting channel
	AlertingChannelSlackFieldIconURL = "icon_url"
	//AlertingChannelSlackFieldChannel const for the channel field of the Slack alerting channel
	AlertingChannelSlackFieldChannel = "channel"
	//ResourceInstanaAlertingChannelSlack the name of the terraform-provider-instana resource to manage alerting channels of type Slack
	ResourceInstanaAlertingChannelSlack = "instana_alerting_channel_slack"
)

//NewAlertingChannelSlackResourceHandle creates the resource handle for Alerting Channels of type Email
func NewAlertingChannelSlackResourceHandle() ResourceHandle {
	return &alertingChannelSlackResourceHandle{}
}

type alertingChannelSlackResourceHandle struct{}

func (h *alertingChannelSlackResourceHandle) GetResource(api restapi.InstanaAPI) restapi.RestResource {
	return api.AlertingChannels()
}

func (h *alertingChannelSlackResourceHandle) GetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		AlertingChannelFieldName:     alertingChannelNameSchemaField,
		AlertingChannelFieldFullName: alertingChannelFullNameSchemaField,
		AlertingChannelSlackFieldWebhookURL: &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "The webhook URL of the Slack alerting channel",
		},
		AlertingChannelSlackFieldIconURL: &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The icon URL of the Slack alerting channel",
		},
		AlertingChannelSlackFieldChannel: &schema.Schema{
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The Slack channel of the Slack alerting channel",
		},
	}
}

func (h *alertingChannelSlackResourceHandle) GetResourceName() string {
	return ResourceInstanaAlertingChannelSlack
}

func (h *alertingChannelSlackResourceHandle) UpdateState(d *schema.ResourceData, obj restapi.InstanaDataObject) {
	alertingChannel := obj.(restapi.AlertingChannel)
	d.Set(AlertingChannelFieldFullName, alertingChannel.Name)
	d.Set(AlertingChannelSlackFieldWebhookURL, alertingChannel.WebhookURL)
	d.Set(AlertingChannelSlackFieldIconURL, alertingChannel.IconURL)
	d.Set(AlertingChannelSlackFieldChannel, alertingChannel.Channel)
	d.SetId(alertingChannel.ID)
}

func (h *alertingChannelSlackResourceHandle) ConvertStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) restapi.InstanaDataObject {
	name := computeFullAlertingChannelNameString(d, formatter)
	webhookURL := d.Get(AlertingChannelSlackFieldWebhookURL).(string)
	iconURL := d.Get(AlertingChannelSlackFieldIconURL).(string)
	channel := d.Get(AlertingChannelSlackFieldChannel).(string)
	return restapi.AlertingChannel{
		ID:         d.Id(),
		Name:       name,
		Kind:       restapi.SlackChannelType,
		WebhookURL: &webhookURL,
		IconURL:    &iconURL,
		Channel:    &channel,
	}
}
