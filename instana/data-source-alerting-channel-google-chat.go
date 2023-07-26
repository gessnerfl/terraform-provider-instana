package instana

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// NewAlertingChannelGoogleChatDataSource creates a new DataSource for Google Chat alerting channel
func NewAlertingChannelGoogleChatDataSource() DataSource {
	return &alertingChannelGoogleChatDataSource{}
}

const (
	//DataSourceAlertingChannelGoogleChat the name of the terraform-provider-instana resource to manage alerting channels of type Google Chat
	DataSourceAlertingChannelGoogleChat = "instana_alerting_channel_google_chat"
)

type alertingChannelGoogleChatDataSource struct{}

// CreateResource creates the resource handle for Google Chat alerting channel
func (ds *alertingChannelGoogleChatDataSource) CreateResource() *schema.Resource {
	// unimplemented
	return nil
}
