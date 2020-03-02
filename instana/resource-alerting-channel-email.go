package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	//AlertingChannelEmailFieldEmails const for the emails field of the alerting channel
	AlertingChannelEmailFieldEmails = "emails"
	//ResourceInstanaAlertingChannelEmail the name of the terraform-provider-instana resource to manage alerting channels of type email
	ResourceInstanaAlertingChannelEmail = "instana_alerting_channel_email"
)

//AlertingChannelEmailEmailsSchemaField schema definition for instana alerting channel emails field
var AlertingChannelEmailEmailsSchemaField = &schema.Schema{
	Type:     schema.TypeSet,
	MinItems: 1,
	Elem: &schema.Schema{
		Type: schema.TypeString,
	},
	Required:    true,
	Description: "The list of emails of the Email alerting channel",
}

//NewAlertingChannelEmailResourceHandle creates the resource handle for Alerting Channels of type Email
func NewAlertingChannelEmailResourceHandle() *ResourceHandle {
	return &ResourceHandle{
		ResourceName: ResourceInstanaAlertingChannelEmail,
		Schema: map[string]*schema.Schema{
			AlertingChannelFieldName:        alertingChannelNameSchemaField,
			AlertingChannelFieldFullName:    alertingChannelFullNameSchemaField,
			AlertingChannelEmailFieldEmails: AlertingChannelEmailEmailsSchemaField,
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type: alertingChannelEmailSchemaV0().CoreConfigSchema().ImpliedType(),
				Upgrade: func(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
					return rawState, nil
				},
				Version: 0,
			},
		},
		RestResourceFactory:  func(api restapi.InstanaAPI) restapi.RestResource { return api.AlertingChannels() },
		UpdateState:          updateStateForAlertingChannelEmail,
		MapStateToDataObject: mapStateToDataObjectForAlertingChannelEmail,
	}
}

func updateStateForAlertingChannelEmail(d *schema.ResourceData, obj restapi.InstanaDataObject) error {
	alertingChannel := obj.(restapi.AlertingChannel)
	emails := alertingChannel.Emails
	d.Set(AlertingChannelFieldFullName, alertingChannel.Name)
	d.Set(AlertingChannelEmailFieldEmails, emails)
	d.SetId(alertingChannel.ID)
	return nil
}

func mapStateToDataObjectForAlertingChannelEmail(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	name := computeFullAlertingChannelNameString(d, formatter)
	return restapi.AlertingChannel{
		ID:     d.Id(),
		Name:   name,
		Kind:   restapi.EmailChannelType,
		Emails: ReadStringSetParameterFromResource(d, AlertingChannelEmailFieldEmails),
	}, nil
}

func alertingChannelEmailSchemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			AlertingChannelFieldName:     alertingChannelNameSchemaField,
			AlertingChannelFieldFullName: alertingChannelFullNameSchemaField,
			AlertingChannelEmailFieldEmails: &schema.Schema{
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
