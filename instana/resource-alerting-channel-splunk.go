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
func NewAlertingChannelSplunkResourceHandle() *ResourceHandle {
	return &ResourceHandle{
		ResourceName: ResourceInstanaAlertingChannelSplunk,
		Schema: map[string]*schema.Schema{
			AlertingChannelFieldName:     alertingChannelNameSchemaField,
			AlertingChannelFieldFullName: alertingChannelFullNameSchemaField,
			AlertingChannelSplunkFieldURL: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The URL of the Splunk alerting channel",
			},
			AlertingChannelSplunkFieldToken: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The token of the Splunk alerting channel",
			},
		},
		RestResourceFactory:  func(api restapi.InstanaAPI) restapi.RestResource { return api.AlertingChannels() },
		UpdateState:          updateStateForAlertingChannelSplunk,
		MapStateToDataObject: monvertStateToDataObjectForAlertingChannelSplunk,
	}
}

func updateStateForAlertingChannelSplunk(d *schema.ResourceData, obj restapi.InstanaDataObject) error {
	alertingChannel := obj.(*restapi.AlertingChannel)
	d.Set(AlertingChannelFieldFullName, alertingChannel.Name)
	d.Set(AlertingChannelSplunkFieldURL, alertingChannel.URL)
	d.Set(AlertingChannelSplunkFieldToken, alertingChannel.Token)
	d.SetId(alertingChannel.ID)
	return nil
}

func monvertStateToDataObjectForAlertingChannelSplunk(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	name := computeFullAlertingChannelNameString(d, formatter)
	url := d.Get(AlertingChannelSplunkFieldURL).(string)
	token := d.Get(AlertingChannelSplunkFieldToken).(string)
	return &restapi.AlertingChannel{
		ID:    d.Id(),
		Name:  name,
		Kind:  restapi.SplunkChannelType,
		URL:   &url,
		Token: &token,
	}, nil
}
