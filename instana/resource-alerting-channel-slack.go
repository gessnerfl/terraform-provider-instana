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
func NewAlertingChannelSlackResourceHandle() *ResourceHandle {
	return &ResourceHandle{
		ResourceName: ResourceInstanaAlertingChannelSlack,
		Schema: map[string]*schema.Schema{
			AlertingChannelFieldName:     alertingChannelNameSchemaField,
			AlertingChannelFieldFullName: alertingChannelFullNameSchemaField,
			AlertingChannelSlackFieldWebhookURL: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The webhook URL of the Slack alerting channel",
			},
			AlertingChannelSlackFieldIconURL: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The icon URL of the Slack alerting channel",
			},
			AlertingChannelSlackFieldChannel: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Slack channel of the Slack alerting channel",
			},
		},
		RestResourceFactory:  func(api restapi.InstanaAPI) restapi.RestResource { return api.AlertingChannels() },
		UpdateState:          updateStateForAlertingChannelSlack,
		MapStateToDataObject: mapStateToDataObjectForAlertingChannelSlack,
	}
}
func updateStateForAlertingChannelSlack(d *schema.ResourceData, obj restapi.InstanaDataObject) error {
	alertingChannel := obj.(*restapi.AlertingChannel)
	d.Set(AlertingChannelFieldFullName, alertingChannel.Name)
	d.Set(AlertingChannelSlackFieldWebhookURL, alertingChannel.WebhookURL)
	d.Set(AlertingChannelSlackFieldIconURL, alertingChannel.IconURL)
	d.Set(AlertingChannelSlackFieldChannel, alertingChannel.Channel)
	d.SetId(alertingChannel.ID)
	return nil
}

func mapStateToDataObjectForAlertingChannelSlack(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	name := computeFullAlertingChannelNameString(d, formatter)
	webhookURL := d.Get(AlertingChannelSlackFieldWebhookURL).(string)
	iconURL := d.Get(AlertingChannelSlackFieldIconURL).(string)
	channel := d.Get(AlertingChannelSlackFieldChannel).(string)
	return &restapi.AlertingChannel{
		ID:         d.Id(),
		Name:       name,
		Kind:       restapi.SlackChannelType,
		WebhookURL: &webhookURL,
		IconURL:    &iconURL,
		Channel:    &channel,
	}, nil
}
