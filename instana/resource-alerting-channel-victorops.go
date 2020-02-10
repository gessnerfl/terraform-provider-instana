package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	//AlertingChannelVictorOpsFieldAPIKey const for the apiKey field of the VictorOps alerting channel
	AlertingChannelVictorOpsFieldAPIKey = "api_key"
	//AlertingChannelVictorOpsFieldRoutingKey const for the routingKey field of the VictorOps alerting channel
	AlertingChannelVictorOpsFieldRoutingKey = "routing_key"
	//ResourceInstanaAlertingChannelVictorOps the name of the terraform-provider-instana resource to manage alerting channels of type VictorOps
	ResourceInstanaAlertingChannelVictorOps = "instana_alerting_channel_victor_ops"
)

//NewAlertingChannelVictorOpsResourceHandle creates the resource handle for Alerting Channels of type Email
func NewAlertingChannelVictorOpsResourceHandle() ResourceHandle {
	return &alertingChannelVictorOpsResourceHandle{}
}

type alertingChannelVictorOpsResourceHandle struct{}

func (h *alertingChannelVictorOpsResourceHandle) GetResourceFrom(api restapi.InstanaAPI) restapi.RestResource {
	return api.AlertingChannels()
}

func (h *alertingChannelVictorOpsResourceHandle) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		AlertingChannelFieldName:     alertingChannelNameSchemaField,
		AlertingChannelFieldFullName: alertingChannelFullNameSchemaField,
		AlertingChannelVictorOpsFieldAPIKey: &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "The API Key of the VictorOps alerting channel",
		},
		AlertingChannelVictorOpsFieldRoutingKey: &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "The Routing Key of the VictorOps alerting channel",
		},
	}
}

func (h *alertingChannelVictorOpsResourceHandle) SchemaVersion() int {
	return 0
}

func (h *alertingChannelVictorOpsResourceHandle) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{}
}

func (h *alertingChannelVictorOpsResourceHandle) ResourceName() string {
	return ResourceInstanaAlertingChannelVictorOps
}

func (h *alertingChannelVictorOpsResourceHandle) UpdateState(d *schema.ResourceData, obj restapi.InstanaDataObject) {
	alertingChannel := obj.(restapi.AlertingChannel)
	d.Set(AlertingChannelFieldFullName, alertingChannel.Name)
	d.Set(AlertingChannelVictorOpsFieldAPIKey, alertingChannel.APIKey)
	d.Set(AlertingChannelVictorOpsFieldRoutingKey, alertingChannel.RoutingKey)
	d.SetId(alertingChannel.ID)
}

func (h *alertingChannelVictorOpsResourceHandle) ConvertStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) restapi.InstanaDataObject {
	name := computeFullAlertingChannelNameString(d, formatter)
	apiKey := d.Get(AlertingChannelVictorOpsFieldAPIKey).(string)
	routingKey := d.Get(AlertingChannelVictorOpsFieldRoutingKey).(string)
	return restapi.AlertingChannel{
		ID:         d.Id(),
		Name:       name,
		Kind:       restapi.VictorOpsChannelType,
		APIKey:     &apiKey,
		RoutingKey: &routingKey,
	}
}
