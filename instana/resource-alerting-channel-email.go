package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	//AlertingChannelEmailFieldEmails const for the emails field of the alerting channel
	AlertingChannelEmailFieldEmails = "emails"
)

//NewAlertingChannelEmailResource creates the terraform resource for Alerting Channels of type Email
func NewAlertingChannelEmailResource() TerraformResource {
	return NewTerraformResource(NewAlertingChannelEmailResourceHandle())
}

//NewAlertingChannelEmailResourceHandle creates the resource handle for Alerting Channels of type Email
func NewAlertingChannelEmailResourceHandle() ResourceHandle {
	return &alertingChannelEmailResourceHandle{}
}

type alertingChannelEmailResourceHandle struct {
}

func (h *alertingChannelEmailResourceHandle) GetResource(api restapi.InstanaAPI) restapi.RestResource {
	return api.AlertingChannels()
}

func (h *alertingChannelEmailResourceHandle) GetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		AlertingChannelFieldName:     alertingChannelNameSchemaField,
		AlertingChannelFieldFullName: alertingChannelFullNameSchemaField,
		AlertingChannelEmailFieldEmails: &schema.Schema{
			Type:     schema.TypeList,
			MinItems: 1,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Required:    true,
			Description: "The list of emails of the alerting channel",
		},
	}
}

func (h *alertingChannelEmailResourceHandle) GetResourceName() string {
	return "instana_alerting_channel_email"
}

func (h *alertingChannelEmailResourceHandle) UpdateState(d *schema.ResourceData, obj restapi.InstanaDataObject) {
	alertingChannel := obj.(restapi.AlertingChannel)
	emails := alertingChannel.Emails
	d.Set(AlertingChannelFieldFullName, alertingChannel.Name)
	d.Set(AlertingChannelEmailFieldEmails, emails)
	d.SetId(alertingChannel.ID)
}

func (h *alertingChannelEmailResourceHandle) ConvertStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) restapi.InstanaDataObject {
	name := computeFullAlertingChannelNameString(d, formatter)
	return restapi.AlertingChannel{
		ID:     d.Id(),
		Name:   name,
		Kind:   restapi.EmailChannelType,
		Emails: ReadStringArrayParameterFromResource(d, AlertingChannelEmailFieldEmails),
	}
}
