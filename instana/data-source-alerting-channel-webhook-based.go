package instana

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// NewSyntheticLocationDataSource creates a new DataSource for Generic Webhook alerting channel
func NewAlertingChannelWebhookBasedFieldDataSource() DataSource {
	return &alertingChannelWebhookBasedFieldDataSource{}
}

const (
	//DataSourceAlertingChannelWebhookBasedFieldWebhookURL const for the webhookUrl field of the alerting channel
	DataSourceAlertingChannelWebhookBasedFieldWebhookURL = "webhook_url"
)

type alertingChannelWebhookBasedFieldDataSource struct{}

// CreateResource creates the resource handle for Generic Webhook alerting channel
func (ds *alertingChannelWebhookBasedFieldDataSource) CreateResource() *schema.Resource {
	// unimplemented
	return nil
}