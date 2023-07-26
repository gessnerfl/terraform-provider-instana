package instana

import (
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// NewSyntheticLocationDataSource creates a new DataSource for Synthetic Locations
func NewAlertingChannelOffice365DataSource() DataSource {
	return &alertingChannelOffice365DataSource{}
}

const (
	//DataSourceAlertingChannelOffice365 the name of the terraform-provider-instana resource to manage alerting channels of type Office 365
	DataSourceAlertingChannelOffice365 = "instana_alerting_channel_office_365"
)

type alertingChannelOffice365DataSource struct{}

// CreateResource creates the resource handle Synthetic Locations
func (ds *alertingChannelOffice365DataSource) CreateResource() *schema.Resource {
	return &schema.Resource{
		Read: ds.read,
		Schema: map[string]*schema.Schema{
			AlertingChannelFieldName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the alerting channel",
				// TODO: What is the max length here?
				// ValidateFunc: validation.StringLenBetween(0, 512),
			},
		},
	}
}

func (ds *alertingChannelOffice365DataSource) read(d *schema.ResourceData, meta interface{}) error {
	providerMeta := meta.(*ProviderMeta)
	instanaAPI := providerMeta.InstanaAPI

	name := d.Get(AlertingChannelFieldName).(string)

	data, err := instanaAPI.AlertingChannelsDS().GetAll()
	if err != nil {
		return err
	}

	alertChannel, err := ds.findAlertChannel(name, data)

	if err != nil {
		return err
	}

	return ds.updateState(d, alertChannel)
}

func (ds *alertingChannelOffice365DataSource) findAlertChannel(name string, data *[]restapi.InstanaDataObject) (*restapi.AlertingChannel, error) {
	for _, e := range *data {
		alertingChannel, ok := e.(*restapi.AlertingChannel)
		if ok {
			if alertingChannel.Name == name {
				return alertingChannel, nil
			}
		}
	}
	return nil, fmt.Errorf("no alerting channel found")
}

func (ds *alertingChannelOffice365DataSource) updateState(d *schema.ResourceData, alertingChannel *restapi.AlertingChannel) error {
	d.SetId(alertingChannel.ID)
	d.Set(AlertingChannelFieldName, alertingChannel.Name)
	return nil
}
