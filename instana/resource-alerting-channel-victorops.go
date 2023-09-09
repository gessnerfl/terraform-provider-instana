package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	//AlertingChannelVictorOpsFieldAPIKey const for the apiKey field of the VictorOps alerting channel
	AlertingChannelVictorOpsFieldAPIKey = "api_key"
	//AlertingChannelVictorOpsFieldRoutingKey const for the routingKey field of the VictorOps alerting channel
	AlertingChannelVictorOpsFieldRoutingKey = "routing_key"
	//ResourceInstanaAlertingChannelVictorOps the name of the terraform-provider-instana resource to manage alerting channels of type VictorOps
	ResourceInstanaAlertingChannelVictorOps = "instana_alerting_channel_victor_ops"
)

// NewAlertingChannelVictorOpsResourceHandle creates the resource handle for Alerting Channels of type Email
func NewAlertingChannelVictorOpsResourceHandle() ResourceHandle {
	return &alertingChannelVictorOpsResource{
		metaData: ResourceMetaData{
			ResourceName:       ResourceInstanaAlertingChannelVictorOps,
			DeprecationMessage: "This feature will be removed in version 2.x and should be replaced with instana_alerting_channel",
			Schema: map[string]*schema.Schema{
				AlertingChannelFieldName:     alertingChannelNameSchemaField,
				AlertingChannelFieldFullName: alertingChannelFullNameSchemaField,
				AlertingChannelVictorOpsFieldAPIKey: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The API Key of the VictorOps alerting channel",
				},
				AlertingChannelVictorOpsFieldRoutingKey: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The Routing Key of the VictorOps alerting channel",
				},
			},
		},
	}
}

type alertingChannelVictorOpsResource struct {
	metaData ResourceMetaData
}

func (r *alertingChannelVictorOpsResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *alertingChannelVictorOpsResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{}
}

func (r *alertingChannelVictorOpsResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource {
	return api.AlertingChannels()
}

func (r *alertingChannelVictorOpsResource) SetComputedFields(d *schema.ResourceData) {
	//No computed fields defined
}

func (r *alertingChannelVictorOpsResource) UpdateState(d *schema.ResourceData, obj restapi.InstanaDataObject, formatter utils.ResourceNameFormatter) error {
	alertingChannel := obj.(*restapi.AlertingChannel)
	d.Set(AlertingChannelFieldName, formatter.UndoFormat(alertingChannel.Name))
	d.Set(AlertingChannelFieldFullName, alertingChannel.Name)
	d.Set(AlertingChannelVictorOpsFieldAPIKey, alertingChannel.APIKey)
	d.Set(AlertingChannelVictorOpsFieldRoutingKey, alertingChannel.RoutingKey)
	d.SetId(alertingChannel.ID)
	return nil
}

func (r *alertingChannelVictorOpsResource) MapStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	name := computeFullAlertingChannelNameString(d, formatter)
	apiKey := d.Get(AlertingChannelVictorOpsFieldAPIKey).(string)
	routingKey := d.Get(AlertingChannelVictorOpsFieldRoutingKey).(string)
	return &restapi.AlertingChannel{
		ID:         d.Id(),
		Name:       name,
		Kind:       restapi.VictorOpsChannelType,
		APIKey:     &apiKey,
		RoutingKey: &routingKey,
	}, nil
}
