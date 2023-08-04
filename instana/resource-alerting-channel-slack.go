package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

// NewAlertingChannelSlackResourceHandle creates the resource handle for Alerting Channels of type Email
func NewAlertingChannelSlackResourceHandle() ResourceHandle[*restapi.AlertingChannel] {
	return &alertingChannelSlackResource{
		metaData: ResourceMetaData{
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
		},
	}
}

type alertingChannelSlackResource struct {
	metaData ResourceMetaData
}

func (r *alertingChannelSlackResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *alertingChannelSlackResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{}
}

func (r *alertingChannelSlackResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.AlertingChannel] {
	return api.AlertingChannels()
}

func (r *alertingChannelSlackResource) SetComputedFields(_ *schema.ResourceData) error {
	return nil
}

func (r *alertingChannelSlackResource) UpdateState(d *schema.ResourceData, alertingChannel *restapi.AlertingChannel, formatter utils.ResourceNameFormatter) error {
	d.SetId(alertingChannel.ID)
	return tfutils.UpdateState(d, map[string]interface{}{
		AlertingChannelFieldName:            formatter.UndoFormat(alertingChannel.Name),
		AlertingChannelFieldFullName:        alertingChannel.Name,
		AlertingChannelSlackFieldWebhookURL: alertingChannel.WebhookURL,
		AlertingChannelSlackFieldIconURL:    alertingChannel.IconURL,
		AlertingChannelSlackFieldChannel:    alertingChannel.Channel,
	})
}

func (r *alertingChannelSlackResource) MapStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (*restapi.AlertingChannel, error) {
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
