package instana

import (
	"errors"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	//AlertingChannelEmailFieldEmails const for the emails field of thealerting channel
	AlertingChannelEmailFieldEmails = "emails"
)

//CreateResourceAlertingChannelEmail creates the resource definition for the resource instana_alerting_channel_email
func CreateResourceAlertingChannelEmail() *schema.Resource {
	return &schema.Resource{
		Create: CreateAlertingChannelEmail,
		Read:   ReadAlertingChannelEmail,
		Update: UpdateAlertingChannelEmail,
		Delete: DeleteAlertingChannelEmail,

		Schema: map[string]*schema.Schema{
			AlertingChannelFieldName:     alertingChannelNameSchemaField,
			AlertingChannelFieldFullName: alertingChannelFullNameSchemaField,
			AlertingChannelEmailFieldEmails: &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "The list of emails of the alerting channel",
			},
		},
	}
}

//CreateAlertingChannelEmail defines the create operation for the resource instana_alerting_channel_email
func CreateAlertingChannelEmail(d *schema.ResourceData, meta interface{}) error {
	d.SetId(RandomID())
	return UpdateAlertingChannelEmail(d, meta)
}

//ReadAlertingChannelEmail defines the read operation for the resource instana_alerting_channel_email
func ReadAlertingChannelEmail(d *schema.ResourceData, meta interface{}) error {
	providerMeta := meta.(*ProviderMeta)
	instanaAPI := providerMeta.InstanaAPI
	alertingChannelID := d.Id()
	if len(alertingChannelID) == 0 {
		return errors.New("ID of alerting channel email is missing")
	}
	alertingChannel, err := instanaAPI.AlertingChannels().GetOne(alertingChannelID)
	if err != nil {
		if err == restapi.ErrEntityNotFound {
			d.SetId("")
			return nil
		}
		return err
	}
	updateAlertingChannelEmailState(d, alertingChannel, providerMeta.ResourceNameFormatter)
	return nil
}

//UpdateAlertingChannelEmail defines the update operation for the resource instana_alerting_channel_email
func UpdateAlertingChannelEmail(d *schema.ResourceData, meta interface{}) error {
	providerMeta := meta.(*ProviderMeta)
	instanaAPI := providerMeta.InstanaAPI
	alertingChannel := createAlertingChannelEmailFromResourceData(d, providerMeta.ResourceNameFormatter)
	updatedAlertingChannelEmail, err := instanaAPI.AlertingChannels().Upsert(alertingChannel)
	if err != nil {
		return err
	}
	updateAlertingChannelEmailState(d, updatedAlertingChannelEmail, providerMeta.ResourceNameFormatter)
	return nil
}

//DeleteAlertingChannelEmail defines the delete operation for the resource instana_alerting_channel_email
func DeleteAlertingChannelEmail(d *schema.ResourceData, meta interface{}) error {
	providerMeta := meta.(*ProviderMeta)
	instanaAPI := providerMeta.InstanaAPI
	alertingChannel := createAlertingChannelEmailFromResourceData(d, providerMeta.ResourceNameFormatter)
	err := instanaAPI.AlertingChannels().DeleteByID(alertingChannel.ID)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func createAlertingChannelEmailFromResourceData(d *schema.ResourceData, formatter ResourceNameFormatter) restapi.AlertingChannel {
	name := computeFullAlertingChannelNameString(d, formatter)
	return restapi.AlertingChannel{
		ID:     d.Id(),
		Name:   name,
		Kind:   restapi.EmailChannelType,
		Emails: ReadStringArrayParameterFromResource(d, AlertingChannelEmailFieldEmails),
	}
}

func updateAlertingChannelEmailState(d *schema.ResourceData, alertingChannel restapi.AlertingChannel, formatter ResourceNameFormatter) {
	d.Set(AlertingChannelFieldFullName, alertingChannel.Name)
	d.Set(AlertingChannelEmailFieldEmails, alertingChannel.Emails)
	d.SetId(alertingChannel.ID)
}
