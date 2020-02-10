package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	//AlertingChannelSplunkFieldURL const for the url field of the Splunk alerting channel
	AlertingChannelSplunkFieldURL = "url"
	//AlertingChannelSplunkFieldToken const for the token field of the Splunk alerting channel
	AlertingChannelSplunkFieldToken = "token"
	//ResourceInstanaAlertingChannelSplunk the name of the terraform-provider-instana resource to manage alerting channels of type Splunk
	ResourceInstanaAlertingChannelSplunk = "instana_alerting_channel_splunk"
)

//NewAlertingChannelSplunkResourceHandle creates the resource handle for Alerting Channels of type Email
func NewAlertingChannelSplunkResourceHandle() ResourceHandle {
	return &alertingChannelSplunkResourceHandle{}
}

type alertingChannelSplunkResourceHandle struct{}

func (h *alertingChannelSplunkResourceHandle) GetResourceFrom(api restapi.InstanaAPI) restapi.RestResource {
	return api.AlertingChannels()
}

func (h *alertingChannelSplunkResourceHandle) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		AlertingChannelFieldName:     alertingChannelNameSchemaField,
		AlertingChannelFieldFullName: alertingChannelFullNameSchemaField,
		AlertingChannelSplunkFieldURL: &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "The URL of the Splunk alerting channel",
		},
		AlertingChannelSplunkFieldToken: &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "The token of the Splunk alerting channel",
		},
	}
}

func (h *alertingChannelSplunkResourceHandle) SchemaVersion() int {
	return 0
}

func (h *alertingChannelSplunkResourceHandle) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{}
}

func (h *alertingChannelSplunkResourceHandle) ResourceName() string {
	return ResourceInstanaAlertingChannelSplunk
}

func (h *alertingChannelSplunkResourceHandle) UpdateState(d *schema.ResourceData, obj restapi.InstanaDataObject) error {
	alertingChannel := obj.(restapi.AlertingChannel)
	d.Set(AlertingChannelFieldFullName, alertingChannel.Name)
	d.Set(AlertingChannelSplunkFieldURL, alertingChannel.URL)
	d.Set(AlertingChannelSplunkFieldToken, alertingChannel.Token)
	d.SetId(alertingChannel.ID)
	return nil
}

func (h *alertingChannelSplunkResourceHandle) ConvertStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	name := computeFullAlertingChannelNameString(d, formatter)
	url := d.Get(AlertingChannelSplunkFieldURL).(string)
	token := d.Get(AlertingChannelSplunkFieldToken).(string)
	return restapi.AlertingChannel{
		ID:    d.Id(),
		Name:  name,
		Kind:  restapi.SplunkChannelType,
		URL:   &url,
		Token: &token,
	}, nil
}
