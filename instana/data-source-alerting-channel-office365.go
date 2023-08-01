package instana

import (
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// NewAlertingChannelOffice365DataSource creates a new DataSource for Office 365 alerting channel
func NewAlertingChannelOffice365DataSource() DataSource {
	return &alertingChannelOffice365DataSource{}
}

const (
	//DataSourceAlertingChannelOffice365 the name of the terraform-provider-instana resource to manage alerting channels of type Office 365
	DataSourceAlertingChannelOffice365 = "instana_alerting_channel_office_365"
)

type alertingChannelOffice365DataSource struct{}

// CreateResource creates the resource handle for Office 365 alerting channel
func (ds *alertingChannelOffice365DataSource) CreateResource() *schema.Resource {
	return &schema.Resource{
		Read: ds.read,
		Schema: map[string]*schema.Schema{
			AlertingChannelFieldName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the alerting channel",
			},
		},
	}
}

func (ds *alertingChannelOffice365DataSource) read(d *schema.ResourceData, meta interface{}) error {
	providerMeta := meta.(*ProviderMeta)
	instanaAPI := providerMeta.InstanaAPI

	name := d.Get(AlertingChannelFieldName).(string)

	data, err := instanaAPI.AlertingChannels().GetAll()
	if err != nil {
		return err
	}

	alertChannel, err := ds.findAlertChannel(name, data)

	if err != nil {
		return err
	}

	return ds.updateState(d, alertChannel)
}

func (ds *alertingChannelOffice365DataSource) findAlertChannel(name string, data *[]*restapi.AlertingChannel) (*restapi.AlertingChannel, error) {
	for _, alertingChannel := range *data {
		if alertingChannel.Name == name {
			return alertingChannel, nil
		}
	}
	return nil, fmt.Errorf("no alerting channel found")
}

func (ds *alertingChannelOffice365DataSource) updateState(d *schema.ResourceData, alertingChannel *restapi.AlertingChannel) error {
	d.SetId(alertingChannel.ID)
	return d.Set(AlertingChannelFieldName, alertingChannel.Name)
}
