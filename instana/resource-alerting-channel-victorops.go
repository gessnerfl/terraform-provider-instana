package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
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

var (
	alertingChannelVictorOpsSchemaAPIKey = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The API Key of the VictorOps alerting channel",
	}
	alertingChannelVictorOpsSchemaRoutingKey = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The Routing Key of the VictorOps alerting channel",
	}
)

// NewAlertingChannelVictorOpsResourceHandle creates the resource handle for Alerting Channels of type Email
func NewAlertingChannelVictorOpsResourceHandle() ResourceHandle[*restapi.AlertingChannel] {
	return &alertingChannelVictorOpsResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaAlertingChannelVictorOps,
			Schema: map[string]*schema.Schema{
				AlertingChannelFieldName:                alertingChannelNameSchemaField,
				AlertingChannelVictorOpsFieldAPIKey:     alertingChannelVictorOpsSchemaAPIKey,
				AlertingChannelVictorOpsFieldRoutingKey: alertingChannelVictorOpsSchemaRoutingKey,
			},
			SchemaVersion: 1,
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
	return []schema.StateUpgrader{
		{
			Type:    r.schemaV0().CoreConfigSchema().ImpliedType(),
			Upgrade: migrateFullNameToName,
			Version: 0,
		},
	}
}

func (r *alertingChannelVictorOpsResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.AlertingChannel] {
	return api.AlertingChannels()
}

func (r *alertingChannelVictorOpsResource) SetComputedFields(_ *schema.ResourceData) error {
	return nil
}

func (r *alertingChannelVictorOpsResource) UpdateState(d *schema.ResourceData, alertingChannel *restapi.AlertingChannel, _ utils.ResourceNameFormatter) error {
	d.SetId(alertingChannel.ID)
	return tfutils.UpdateState(d, map[string]interface{}{
		AlertingChannelFieldName:                alertingChannel.Name,
		AlertingChannelVictorOpsFieldAPIKey:     alertingChannel.APIKey,
		AlertingChannelVictorOpsFieldRoutingKey: alertingChannel.RoutingKey,
	})
}

func (r *alertingChannelVictorOpsResource) MapStateToDataObject(d *schema.ResourceData, _ utils.ResourceNameFormatter) (*restapi.AlertingChannel, error) {
	apiKey := d.Get(AlertingChannelVictorOpsFieldAPIKey).(string)
	routingKey := d.Get(AlertingChannelVictorOpsFieldRoutingKey).(string)
	return &restapi.AlertingChannel{
		ID:         d.Id(),
		Name:       d.Get(AlertingChannelFieldName).(string),
		Kind:       restapi.VictorOpsChannelType,
		APIKey:     &apiKey,
		RoutingKey: &routingKey,
	}, nil
}

func (r *alertingChannelVictorOpsResource) schemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			AlertingChannelFieldName:                alertingChannelNameSchemaField,
			AlertingChannelFieldFullName:            alertingChannelFullNameSchemaField,
			AlertingChannelVictorOpsFieldAPIKey:     alertingChannelVictorOpsSchemaAPIKey,
			AlertingChannelVictorOpsFieldRoutingKey: alertingChannelVictorOpsSchemaRoutingKey,
		},
	}
}
