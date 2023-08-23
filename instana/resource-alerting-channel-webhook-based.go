package instana

import (
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	//AlertingChannelWebhookBasedFieldWebhookURL const for the webhookUrl field of the alerting channel
	AlertingChannelWebhookBasedFieldWebhookURL = "webhook_url"
	//ResourceInstanaAlertingChannelOffice365 the name of the terraform-provider-instana resource to manage alerting channels of type Office 365
	ResourceInstanaAlertingChannelOffice365 = "instana_alerting_channel_office_365"
	//ResourceInstanaAlertingChannelGoogleChat the name of the terraform-provider-instana resource to manage alerting channels of type Google Chat
	ResourceInstanaAlertingChannelGoogleChat = "instana_alerting_channel_google_chat"
)

var alertingChannelWebhookBasedSchemaWebhookURL = &schema.Schema{
	Type:        schema.TypeString,
	Required:    true,
	Description: fmt.Sprintf("The webhook URL of the alerting channel"),
}

// NewAlertingChannelGoogleChatResourceHandle creates the terraform resource for Alerting Channels of type Google Chat
func NewAlertingChannelGoogleChatResourceHandle() ResourceHandle[*restapi.AlertingChannel] {
	return newAlertingChannelWebhookBasedResourceHandle(restapi.GoogleChatChannelType, ResourceInstanaAlertingChannelGoogleChat)
}

// NewAlertingChannelOffice365ResourceHandle creates the terraform resource for Alerting Channels of type Office 356
func NewAlertingChannelOffice365ResourceHandle() ResourceHandle[*restapi.AlertingChannel] {
	return newAlertingChannelWebhookBasedResourceHandle(restapi.Office365ChannelType, ResourceInstanaAlertingChannelOffice365)
}

func newAlertingChannelWebhookBasedResourceHandle(channelType restapi.AlertingChannelType, resourceName string) ResourceHandle[*restapi.AlertingChannel] {
	return &alertingChannelWebhookBasedResource{
		metaData: ResourceMetaData{
			ResourceName: resourceName,
			Schema: map[string]*schema.Schema{
				AlertingChannelFieldName:                   alertingChannelNameSchemaField,
				AlertingChannelWebhookBasedFieldWebhookURL: alertingChannelWebhookBasedSchemaWebhookURL,
			},
			SchemaVersion: 1,
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
	return []schema.StateUpgrader{
		{
			Type:    r.schemaV0().CoreConfigSchema().ImpliedType(),
			Upgrade: migrateFullNameToName,
			Version: 0,
		},
	}
}

func (r *alertingChannelWebhookBasedResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.AlertingChannel] {
	return api.AlertingChannels()
}

func (r *alertingChannelWebhookBasedResource) SetComputedFields(_ *schema.ResourceData) error {
	return nil
}

func (r *alertingChannelWebhookBasedResource) UpdateState(d *schema.ResourceData, alertingChannel *restapi.AlertingChannel, _ utils.ResourceNameFormatter) error {
	d.SetId(alertingChannel.ID)
	return tfutils.UpdateState(d, map[string]interface{}{
		AlertingChannelFieldName:                   alertingChannel.Name,
		AlertingChannelWebhookBasedFieldWebhookURL: alertingChannel.WebhookURL,
	})
}

func (r *alertingChannelWebhookBasedResource) MapStateToDataObject(d *schema.ResourceData, _ utils.ResourceNameFormatter) (*restapi.AlertingChannel, error) {
	webhookURL := d.Get(AlertingChannelWebhookBasedFieldWebhookURL).(string)
	return &restapi.AlertingChannel{
		ID:         d.Id(),
		Name:       d.Get(AlertingChannelFieldName).(string),
		Kind:       r.channelType,
		WebhookURL: &webhookURL,
	}, nil
}

func (r *alertingChannelWebhookBasedResource) schemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			AlertingChannelFieldName:                   alertingChannelNameSchemaField,
			AlertingChannelFieldFullName:               alertingChannelFullNameSchemaField,
			AlertingChannelWebhookBasedFieldWebhookURL: alertingChannelWebhookBasedSchemaWebhookURL,
		},
	}
}
