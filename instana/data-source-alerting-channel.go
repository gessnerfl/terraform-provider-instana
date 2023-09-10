package instana

import (
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
	"strings"
)

// NewAlertingChannelDataSource creates a new DataSource for alerting channel
func NewAlertingChannelDataSource() DataSource {
	return &alertingChannelDataSource{}
}

const (
	//DataSourceAlertingChannel the name of the terraform-provider-instana data source to read alerting channel
	DataSourceAlertingChannel = "instana_alerting_channel"
)

type alertingChannelDataSource struct{}

// CreateResource creates the resource handle for Office 365 alerting channel
func (ds *alertingChannelDataSource) CreateResource() *schema.Resource {
	return &schema.Resource{
		Read:   ds.read,
		Schema: ds.convertResourceSchema(),
	}
}

func (ds *alertingChannelDataSource) convertResourceSchema() map[string]*schema.Schema {
	resourceSchema := NewAlertingChannelResourceHandle().MetaData().Schema

	return ds.convertSchemaMap(resourceSchema)
}

func (ds *alertingChannelDataSource) convertSchemaMap(schemaMap map[string]*schema.Schema) map[string]*schema.Schema {
	result := make(map[string]*schema.Schema)

	for k, v := range schemaMap {
		s := *v
		if k == AlertingChannelFieldName {
			s.Required = true
			s.Optional = false
			s.Computed = false
		} else {
			s.Required = false
			s.Optional = false
			s.Computed = true
			s.ExactlyOneOf = []string{}
			s.RequiredWith = []string{}
			s.ValidateFunc = nil
			s.ConflictsWith = []string{}
			s.MinItems = 0
			s.MaxItems = 0
		}

		if s.Type == schema.TypeList && reflect.TypeOf(s.Elem) == reflect.TypeOf(&schema.Resource{}) {
			nestedSchema := s.Elem.(*schema.Resource).Schema
			convertedNestedSchema := ds.convertSchemaMap(nestedSchema)
			s.Elem = &schema.Resource{
				Schema: convertedNestedSchema,
			}
		} else if (s.Type == schema.TypeList || s.Type == schema.TypeSet || s.Type == schema.TypeMap) && reflect.TypeOf(s.Elem) == reflect.TypeOf(&schema.Schema{}) {
			nestedSchema := *s.Elem.(*schema.Schema)
			s.Elem = &nestedSchema
		}

		result[k] = &s
	}

	return result
}

func (ds *alertingChannelDataSource) read(d *schema.ResourceData, meta interface{}) error {
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

func (ds *alertingChannelDataSource) findAlertChannel(name string, data *[]restapi.InstanaDataObject) (*restapi.AlertingChannelDS, error) {
	for _, e := range *data {
		alertingChannel, ok := e.(restapi.AlertingChannelDS)
		if ok {
			if alertingChannel.Name == name {
				return &alertingChannel, nil
			}
		}
	}
	return nil, fmt.Errorf("no alerting channel found")
}

func (ds *alertingChannelDataSource) updateState(d *schema.ResourceData, alertingChannel *restapi.AlertingChannelDS) error {
	data, err := ds.mapChannelToState(alertingChannel)
	if err != nil {
		return err
	}

	d.SetId(alertingChannel.ID)
	return tfutils.UpdateState(d, data)
}

func (ds *alertingChannelDataSource) mapChannelToState(channel *restapi.AlertingChannelDS) (map[string]interface{}, error) {
	if channel.Kind == restapi.EmailChannelType {
		return ds.mapEmailChannelToState(channel), nil
	}
	if channel.Kind == restapi.OpsGenieChannelType {
		return ds.mapOpsGenieChannelToState(channel), nil
	}
	if channel.Kind == restapi.PagerDutyChannelType {
		return ds.mapPagerDutyChannelToState(channel), nil
	}
	if channel.Kind == restapi.SlackChannelType {
		return ds.mapSlackChannelToState(channel), nil
	}
	if channel.Kind == restapi.SplunkChannelType {
		return ds.mapSplunkChannelToState(channel), nil
	}
	if channel.Kind == restapi.VictorOpsChannelType {
		return ds.mapVictorOpsChannelToState(channel), nil
	}
	if channel.Kind == restapi.WebhookChannelType {
		return ds.mapWebhookChannelToState(channel), nil
	}
	if channel.Kind == restapi.Office365ChannelType {
		return ds.mapOffice365ChannelToState(channel), nil
	}
	if channel.Kind == restapi.GoogleChatChannelType {
		return ds.mapGoogleChatChannelToState(channel), nil
	}
	return nil, fmt.Errorf("received unsupported alerting channel of type %s", channel.Kind)
}

func (ds *alertingChannelDataSource) mapEmailChannelToState(channel *restapi.AlertingChannelDS) map[string]interface{} {
	return map[string]interface{}{
		AlertingChannelFieldName: channel.Name,
		AlertingChannelFieldChannelEmail: []interface{}{
			map[string]interface{}{
				AlertingChannelEmailFieldEmails: channel.Emails,
			},
		},
	}
}

func (ds *alertingChannelDataSource) mapOpsGenieChannelToState(channel *restapi.AlertingChannelDS) map[string]interface{} {
	tags := ds.convertCommaSeparatedListToSlice(*channel.Tags)
	return map[string]interface{}{
		AlertingChannelFieldName: channel.Name,
		AlertingChannelFieldChannelOpsGenie: []interface{}{
			map[string]interface{}{
				AlertingChannelOpsGenieFieldAPIKey: channel.APIKey,
				AlertingChannelOpsGenieFieldRegion: channel.Region,
				AlertingChannelOpsGenieFieldTags:   tags,
			},
		},
	}
}

func (ds *alertingChannelDataSource) convertCommaSeparatedListToSlice(csv string) []string {
	entries := strings.Split(csv, ",")
	result := make([]string, len(entries))
	for i, e := range entries {
		result[i] = strings.TrimSpace(e)
	}
	return result
}

func (ds *alertingChannelDataSource) mapPagerDutyChannelToState(channel *restapi.AlertingChannelDS) map[string]interface{} {
	return map[string]interface{}{
		AlertingChannelFieldName: channel.Name,
		AlertingChannelFieldChannelPageDuty: []interface{}{
			map[string]interface{}{
				AlertingChannelPagerDutyFieldServiceIntegrationKey: channel.ServiceIntegrationKey,
			},
		},
	}
}

func (ds *alertingChannelDataSource) mapSlackChannelToState(channel *restapi.AlertingChannelDS) map[string]interface{} {
	return map[string]interface{}{
		AlertingChannelFieldName: channel.Name,
		AlertingChannelFieldChannelSlack: []interface{}{
			map[string]interface{}{
				AlertingChannelSlackFieldWebhookURL: channel.WebhookURL,
				AlertingChannelSlackFieldIconURL:    channel.IconURL,
				AlertingChannelSlackFieldChannel:    channel.Channel,
			},
		},
	}
}

func (ds *alertingChannelDataSource) mapSplunkChannelToState(channel *restapi.AlertingChannelDS) map[string]interface{} {
	return map[string]interface{}{
		AlertingChannelFieldName: channel.Name,
		AlertingChannelFieldChannelSplunk: []interface{}{
			map[string]interface{}{
				AlertingChannelSplunkFieldURL:   channel.URL,
				AlertingChannelSplunkFieldToken: channel.Token,
			},
		},
	}
}

func (ds *alertingChannelDataSource) mapVictorOpsChannelToState(channel *restapi.AlertingChannelDS) map[string]interface{} {
	return map[string]interface{}{
		AlertingChannelFieldName: channel.Name,
		AlertingChannelFieldChannelVictorOps: []interface{}{
			map[string]interface{}{
				AlertingChannelVictorOpsFieldAPIKey:     channel.APIKey,
				AlertingChannelVictorOpsFieldRoutingKey: channel.RoutingKey,
			},
		},
	}
}

func (ds *alertingChannelDataSource) mapWebhookChannelToState(channel *restapi.AlertingChannelDS) map[string]interface{} {
	headers := ds.createHTTPHeaderMapFromList(channel.Headers)
	return map[string]interface{}{
		AlertingChannelFieldName: channel.Name,
		AlertingChannelFieldChannelWebhook: []interface{}{
			map[string]interface{}{
				AlertingChannelWebhookFieldWebhookURLs: channel.WebhookURLs,
				AlertingChannelWebhookFieldHTTPHeaders: headers,
			},
		},
	}
}

func (ds *alertingChannelDataSource) createHTTPHeaderMapFromList(headers []string) map[string]interface{} {
	result := make(map[string]interface{})
	for _, header := range headers {
		keyValue := strings.Split(header, ":")
		if len(keyValue) == 2 {
			result[strings.TrimSpace(keyValue[0])] = strings.TrimSpace(keyValue[1])
		} else {
			result[strings.TrimSpace(keyValue[0])] = ""
		}
	}
	return result
}

func (ds *alertingChannelDataSource) mapOffice365ChannelToState(channel *restapi.AlertingChannelDS) map[string]interface{} {
	return map[string]interface{}{
		AlertingChannelFieldName: channel.Name,
		AlertingChannelFieldChannelOffice365: []interface{}{
			map[string]interface{}{
				AlertingChannelWebhookBasedFieldWebhookURL: channel.WebhookURL,
			},
		},
	}
}

func (ds *alertingChannelDataSource) mapGoogleChatChannelToState(channel *restapi.AlertingChannelDS) map[string]interface{} {
	return map[string]interface{}{
		AlertingChannelFieldName: channel.Name,
		AlertingChannelFieldChannelGoogleChat: []interface{}{
			map[string]interface{}{
				AlertingChannelWebhookBasedFieldWebhookURL: channel.WebhookURL,
			},
		},
	}
}
