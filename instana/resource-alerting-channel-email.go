package instana

import (
	"context"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	//AlertingChannelEmailFieldEmails const for the emails field of the alerting channel
	AlertingChannelEmailFieldEmails = "emails"
	//ResourceInstanaAlertingChannelEmail the name of the terraform-provider-instana resource to manage alerting channels of type email
	ResourceInstanaAlertingChannelEmail = "instana_alerting_channel_email"
)

// AlertingChannelEmailEmailsSchemaField schema definition for instana alerting channel emails field
var AlertingChannelEmailEmailsSchemaField = &schema.Schema{
	Type:     schema.TypeSet,
	MinItems: 1,
	Elem: &schema.Schema{
		Type: schema.TypeString,
	},
	Required:    true,
	Description: "The list of emails of the Email alerting channel",
}

// NewAlertingChannelEmailResourceHandle creates the resource handle for Alerting Channels of type Email
func NewAlertingChannelEmailResourceHandle() ResourceHandle[*restapi.AlertingChannel] {
	return &alertingChannelEmailResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaAlertingChannelEmail,
			Schema: map[string]*schema.Schema{
				AlertingChannelFieldName:        alertingChannelNameSchemaField,
				AlertingChannelEmailFieldEmails: AlertingChannelEmailEmailsSchemaField,
			},
			SchemaVersion: 2,
		},
	}
}

type alertingChannelEmailResource struct {
	metaData ResourceMetaData
}

func (r *alertingChannelEmailResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *alertingChannelEmailResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{
		{
			Type: r.schemaV0().CoreConfigSchema().ImpliedType(),
			Upgrade: func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
				return rawState, nil
			},
			Version: 0,
		},
		{
			Type:    r.schemaV1().CoreConfigSchema().ImpliedType(),
			Upgrade: r.stateUpgradeV1,
			Version: 1,
		},
	}
}

func (r *alertingChannelEmailResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.AlertingChannel] {
	return api.AlertingChannels()
}

func (r *alertingChannelEmailResource) SetComputedFields(_ *schema.ResourceData) error {
	return nil
}

func (r *alertingChannelEmailResource) UpdateState(d *schema.ResourceData, alertingChannel *restapi.AlertingChannel, _ utils.ResourceNameFormatter) error {
	emails := alertingChannel.Emails
	d.SetId(alertingChannel.ID)
	return tfutils.UpdateState(d, map[string]interface{}{
		AlertingChannelFieldName:        alertingChannel.Name,
		AlertingChannelEmailFieldEmails: emails,
	})
}

func (r *alertingChannelEmailResource) MapStateToDataObject(d *schema.ResourceData, _ utils.ResourceNameFormatter) (*restapi.AlertingChannel, error) {
	return &restapi.AlertingChannel{
		ID:     d.Id(),
		Name:   d.Get(AlertingChannelFieldName).(string),
		Kind:   restapi.EmailChannelType,
		Emails: ReadStringSetParameterFromResource(d, AlertingChannelEmailFieldEmails),
	}, nil
}

func (r *alertingChannelEmailResource) schemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			AlertingChannelFieldName:     alertingChannelNameSchemaField,
			AlertingChannelFieldFullName: alertingChannelFullNameSchemaField,
			AlertingChannelEmailFieldEmails: {
				Type:     schema.TypeList,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "The list of emails of the Email alerting channel",
			},
		},
	}
}

func (r *alertingChannelEmailResource) stateUpgradeV1(_ context.Context, state map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	if _, ok := state[AlertingChannelFieldFullName]; ok {
		state[AlertingChannelFieldName] = state[AlertingChannelFieldFullName]
		delete(state, AlertingChannelFieldFullName)
	}
	return state, nil
}

func (r *alertingChannelEmailResource) schemaV1() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			AlertingChannelFieldName:        alertingChannelNameSchemaField,
			AlertingChannelFieldFullName:    alertingChannelFullNameSchemaField,
			AlertingChannelEmailFieldEmails: AlertingChannelEmailEmailsSchemaField,
		},
	}
}
