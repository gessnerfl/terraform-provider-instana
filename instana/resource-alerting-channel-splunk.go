package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	//AlertingChannelSplunkFieldURL const for the url field of the Splunk alerting channel
	AlertingChannelSplunkFieldURL = "url"
	//AlertingChannelSplunkFieldToken const for the token field of the Splunk alerting channel
	AlertingChannelSplunkFieldToken = "token"
	//ResourceInstanaAlertingChannelSplunk the name of the terraform-provider-instana resource to manage alerting channels of type Splunk
	ResourceInstanaAlertingChannelSplunk = "instana_alerting_channel_splunk"
)

var (
	alertingChannelSplunkSchemaURL = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The URL of the Splunk alerting channel",
	}
	alertingChannelSplunkSchemaToken = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The token of the Splunk alerting channel",
		Sensitive:   true,
	}
)

// NewAlertingChannelSplunkResourceHandle creates the resource handle for Alerting Channels of type Email
func NewAlertingChannelSplunkResourceHandle() ResourceHandle[*restapi.AlertingChannel] {
	return &alertingChannelSplunkResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaAlertingChannelSplunk,
			DeprecationMessage: "This feature will be removed in version 2.x and should be replaced with instana_alerting_channel",
			Schema: map[string]*schema.Schema{
				AlertingChannelFieldName:        alertingChannelNameSchemaField,
				AlertingChannelSplunkFieldURL:   alertingChannelSplunkSchemaURL,
				AlertingChannelSplunkFieldToken: alertingChannelSplunkSchemaToken,
			},
			SchemaVersion: 1,
		},
	}
}

type alertingChannelSplunkResource struct {
	metaData ResourceMetaData
}

func (r *alertingChannelSplunkResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *alertingChannelSplunkResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{
		{
			Type:    r.schemaV0().CoreConfigSchema().ImpliedType(),
			Upgrade: migrateFullNameToName,
			Version: 0,
		},
	}
}

func (r *alertingChannelSplunkResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.AlertingChannel] {
	return api.AlertingChannels()
}

func (r *alertingChannelSplunkResource) SetComputedFields(_ *schema.ResourceData) error {
	return nil
}

func (r *alertingChannelSplunkResource) UpdateState(d *schema.ResourceData, alertingChannel *restapi.AlertingChannel) error {
	d.SetId(alertingChannel.ID)
	return tfutils.UpdateState(d, map[string]interface{}{
		AlertingChannelFieldName:        alertingChannel.Name,
		AlertingChannelSplunkFieldURL:   alertingChannel.URL,
		AlertingChannelSplunkFieldToken: alertingChannel.Token,
	})
}

func (r *alertingChannelSplunkResource) MapStateToDataObject(d *schema.ResourceData) (*restapi.AlertingChannel, error) {
	url := d.Get(AlertingChannelSplunkFieldURL).(string)
	token := d.Get(AlertingChannelSplunkFieldToken).(string)
	return &restapi.AlertingChannel{
		ID:    d.Id(),
		Name:  d.Get(AlertingChannelFieldName).(string),
		Kind:  restapi.SplunkChannelType,
		URL:   &url,
		Token: &token,
	}, nil
}

func (r *alertingChannelSplunkResource) schemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			AlertingChannelFieldName:        alertingChannelNameSchemaField,
			AlertingChannelFieldFullName:    alertingChannelFullNameSchemaField,
			AlertingChannelSplunkFieldURL:   alertingChannelSplunkSchemaURL,
			AlertingChannelSplunkFieldToken: alertingChannelSplunkSchemaToken,
		},
	}
}
