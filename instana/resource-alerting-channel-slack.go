package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
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

var (
	alertingChannelSlackSchemaWebhookURL = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The webhook URL of the Slack alerting channel",
	}
	alertingChannelSlackSchemaIconURL = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The icon URL of the Slack alerting channel",
	}
	alertingChannelSlackSchemaChannel = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The Slack channel of the Slack alerting channel",
	}
)

// NewAlertingChannelSlackResourceHandle creates the resource handle for Alerting Channels of type Email
func NewAlertingChannelSlackResourceHandle() ResourceHandle[*restapi.AlertingChannel] {
	return &alertingChannelSlackResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaAlertingChannelSlack,
			Schema: map[string]*schema.Schema{
				AlertingChannelFieldName:            alertingChannelNameSchemaField,
				AlertingChannelSlackFieldWebhookURL: alertingChannelSlackSchemaWebhookURL,
				AlertingChannelSlackFieldIconURL:    alertingChannelSlackSchemaIconURL,
				AlertingChannelSlackFieldChannel:    alertingChannelSlackSchemaChannel,
			},
			SchemaVersion: 1,
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
	return []schema.StateUpgrader{
		{
			Type:    r.schemaV0().CoreConfigSchema().ImpliedType(),
			Upgrade: migrateFullNameToName,
			Version: 0,
		},
	}
}

func (r *alertingChannelSlackResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.AlertingChannel] {
	return api.AlertingChannels()
}

func (r *alertingChannelSlackResource) SetComputedFields(_ *schema.ResourceData) error {
	return nil
}

func (r *alertingChannelSlackResource) UpdateState(d *schema.ResourceData, alertingChannel *restapi.AlertingChannel) error {
	d.SetId(alertingChannel.ID)
	return tfutils.UpdateState(d, map[string]interface{}{
		AlertingChannelFieldName:            alertingChannel.Name,
		AlertingChannelSlackFieldWebhookURL: alertingChannel.WebhookURL,
		AlertingChannelSlackFieldIconURL:    alertingChannel.IconURL,
		AlertingChannelSlackFieldChannel:    alertingChannel.Channel,
	})
}

func (r *alertingChannelSlackResource) MapStateToDataObject(d *schema.ResourceData) (*restapi.AlertingChannel, error) {
	webhookURL := d.Get(AlertingChannelSlackFieldWebhookURL).(string)
	iconURL := d.Get(AlertingChannelSlackFieldIconURL).(string)
	channel := d.Get(AlertingChannelSlackFieldChannel).(string)
	return &restapi.AlertingChannel{
		ID:         d.Id(),
		Name:       d.Get(AlertingChannelFieldName).(string),
		Kind:       restapi.SlackChannelType,
		WebhookURL: &webhookURL,
		IconURL:    &iconURL,
		Channel:    &channel,
	}, nil
}

func (r *alertingChannelSlackResource) schemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			AlertingChannelFieldName:            alertingChannelNameSchemaField,
			AlertingChannelFieldFullName:        alertingChannelFullNameSchemaField,
			AlertingChannelSlackFieldWebhookURL: alertingChannelSlackSchemaWebhookURL,
			AlertingChannelSlackFieldIconURL:    alertingChannelSlackSchemaIconURL,
			AlertingChannelSlackFieldChannel:    alertingChannelSlackSchemaChannel,
		},
	}
}
