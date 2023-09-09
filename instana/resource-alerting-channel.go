package instana

import (
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strings"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	//ResourceInstanaAlertingChannel the name of the terraform-provider-instana resource to manage alerting channels
	ResourceInstanaAlertingChannel = "instana_alerting_channel"

	//AlertingChannelFieldChannelEmail const for schema field of the email channel
	AlertingChannelFieldChannelEmail = "email"
	//AlertingChannelFieldChannelOpsGenie const for schema field of the OpsGenie channel
	AlertingChannelFieldChannelOpsGenie = "ops_genie"
	//AlertingChannelFieldChannelPageDuty const for schema field of the PagerDuty channel
	AlertingChannelFieldChannelPageDuty = "pager_duty"
	//AlertingChannelFieldChannelSlack const for schema field of the Slack channel
	AlertingChannelFieldChannelSlack = "slack"
	//AlertingChannelFieldChannelSplunk const for schema field of the Splunk channel
	AlertingChannelFieldChannelSplunk = "splunk"
	//AlertingChannelFieldChannelVictorOps const for schema field of the Victor Ops channel
	AlertingChannelFieldChannelVictorOps = "victor_ops"
	//AlertingChannelFieldChannelWebhook const for schema field of the Webhook channel
	AlertingChannelFieldChannelWebhook = "webhook"
	//AlertingChannelFieldChannelOffice365 const for schema field of the Office 365 channel
	AlertingChannelFieldChannelOffice365 = "office_365"
	//AlertingChannelFieldChannelGoogleChat const for schema field of the Google Chat channel
	AlertingChannelFieldChannelGoogleChat = "google_chat"
)

var AlertingChannelTypeFields = []string{
	AlertingChannelFieldChannelEmail,
	AlertingChannelFieldChannelOpsGenie,
	AlertingChannelFieldChannelPageDuty,
	AlertingChannelFieldChannelSlack,
	AlertingChannelFieldChannelSplunk,
	AlertingChannelFieldChannelVictorOps,
	AlertingChannelFieldChannelWebhook,
	AlertingChannelFieldChannelOffice365,
	AlertingChannelFieldChannelGoogleChat,
}

// NewAlertingChannelResourceHandle creates the resource handle for Alerting Channels
func NewAlertingChannelResourceHandle() ResourceHandle[*restapi.AlertingChannel] {
	supportedOpsGenieRegions := []string{"EU", "US"}
	return &alertingChannelResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaAlertingChannel,
			Schema: map[string]*schema.Schema{
				AlertingChannelFieldName: alertingChannelNameSchemaField,
				AlertingChannelFieldChannelEmail: {
					Type:         schema.TypeList,
					Optional:     true,
					MinItems:     1,
					MaxItems:     1,
					Description:  "The configuration of the Email channel",
					ExactlyOneOf: AlertingChannelTypeFields,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							AlertingChannelEmailFieldEmails: {
								Type:     schema.TypeSet,
								MinItems: 1,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
								Required:    true,
								Description: "The list of emails of the Email alerting channel",
							},
						},
					},
				},
				AlertingChannelFieldChannelOpsGenie: {
					Type:         schema.TypeList,
					Optional:     true,
					MinItems:     1,
					MaxItems:     1,
					Description:  "The configuration of the Ops Genie channel",
					ExactlyOneOf: AlertingChannelTypeFields,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							AlertingChannelOpsGenieFieldAPIKey: {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The OpsGenie API Key of the OpsGenie alerting channel",
							},
							AlertingChannelOpsGenieFieldTags: {
								Type:     schema.TypeList,
								MinItems: 1,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
								Required:    true,
								Description: "The OpsGenie tags of the OpsGenie alerting channel",
							},
							AlertingChannelOpsGenieFieldRegion: {
								Type:         schema.TypeString,
								Required:     true,
								ValidateFunc: validation.StringInSlice(supportedOpsGenieRegions, false),
								Description:  fmt.Sprintf("The OpsGenie region (%s) of the OpsGenie alerting channel", strings.Join(supportedOpsGenieRegions, ", ")),
							},
						},
					},
				},
				AlertingChannelFieldChannelPageDuty: {
					Type:         schema.TypeList,
					Optional:     true,
					MinItems:     1,
					MaxItems:     1,
					Description:  "The configuration of the Pager Duty channel",
					ExactlyOneOf: AlertingChannelTypeFields,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							AlertingChannelPagerDutyFieldServiceIntegrationKey: {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The Service Integration Key of the PagerDuty alerting channel",
							},
						},
					},
				},
				AlertingChannelFieldChannelSlack: {
					Type:         schema.TypeList,
					Optional:     true,
					MinItems:     1,
					MaxItems:     1,
					Description:  "The configuration of the Slack channel",
					ExactlyOneOf: AlertingChannelTypeFields,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							AlertingChannelSlackFieldWebhookURL: {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The webhook URL of the Slack alerting channel",
							},
							AlertingChannelSlackFieldIconURL: {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "The icon URL of the Slack alerting channel",
							},
							AlertingChannelSlackFieldChannel: {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "The Slack channel of the Slack alerting channel",
							},
						},
					},
				},
				AlertingChannelFieldChannelSplunk: {
					Type:         schema.TypeList,
					Optional:     true,
					MinItems:     1,
					MaxItems:     1,
					Description:  "The configuration of the Splunk channel",
					ExactlyOneOf: AlertingChannelTypeFields,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
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
					},
				},
				AlertingChannelFieldChannelVictorOps: {
					Type:         schema.TypeList,
					Optional:     true,
					MinItems:     1,
					MaxItems:     1,
					Description:  "The configuration of the ViktorOps channel",
					ExactlyOneOf: AlertingChannelTypeFields,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							AlertingChannelVictorOpsFieldAPIKey: {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The API Key of the VictorOps alerting channel",
							},
							AlertingChannelVictorOpsFieldRoutingKey: {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The Routing Key of the VictorOps alerting channel",
							},
						},
					},
				},
				AlertingChannelFieldChannelWebhook: {
					Type:         schema.TypeList,
					Optional:     true,
					MinItems:     1,
					MaxItems:     1,
					Description:  "The configuration of the Webhook channel",
					ExactlyOneOf: AlertingChannelTypeFields,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							AlertingChannelWebhookFieldWebhookURLs: {
								Type:     schema.TypeSet,
								MinItems: 1,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
								Required:    true,
								Description: "The list of webhook urls of the Webhook alerting channel",
							},
							AlertingChannelWebhookFieldHTTPHeaders: {
								Type: schema.TypeMap,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
								Optional:    true,
								Description: "The optional map of HTTP headers of the Webhook alerting channel",
							},
						},
					},
				},
				AlertingChannelFieldChannelOffice365: {
					Type:         schema.TypeList,
					Optional:     true,
					MinItems:     1,
					MaxItems:     1,
					Description:  "The configuration of the Office 365 channel",
					ExactlyOneOf: AlertingChannelTypeFields,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							AlertingChannelWebhookBasedFieldWebhookURL: {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The webhook URL of the Office 365 alerting channel",
							},
						},
					},
				},
				AlertingChannelFieldChannelGoogleChat: {
					Type:         schema.TypeList,
					Optional:     true,
					MinItems:     1,
					MaxItems:     1,
					Description:  "The configuration of the Google Chat channel",
					ExactlyOneOf: AlertingChannelTypeFields,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							AlertingChannelWebhookBasedFieldWebhookURL: {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The webhook URL of the Google Chat alerting channel",
							},
						},
					},
				},
			},
			SchemaVersion: 0,
		},
	}
}

type alertingChannelResource struct {
	metaData ResourceMetaData
}

func (r *alertingChannelResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *alertingChannelResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{}
}

func (r *alertingChannelResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.AlertingChannel] {
	return api.AlertingChannels()
}

func (r *alertingChannelResource) SetComputedFields(_ *schema.ResourceData) error {
	return nil
}

func (r *alertingChannelResource) UpdateState(d *schema.ResourceData, alertingChannel *restapi.AlertingChannel) error {
	data, err := r.mapChannelToState(alertingChannel)
	if err != nil {
		return err
	}

	d.SetId(alertingChannel.ID)
	return tfutils.UpdateState(d, data)
}

func (r *alertingChannelResource) mapChannelToState(channel *restapi.AlertingChannel) (map[string]interface{}, error) {
	if channel.Kind == restapi.EmailChannelType {
		return r.mapEmailChannelToState(channel), nil
	}
	if channel.Kind == restapi.OpsGenieChannelType {
		return r.mapOpsGenieChannelToState(channel), nil
	}
	if channel.Kind == restapi.PagerDutyChannelType {
		return r.mapPagerDutyChannelToState(channel), nil
	}
	if channel.Kind == restapi.SlackChannelType {
		return r.mapSlackChannelToState(channel), nil
	}
	if channel.Kind == restapi.SplunkChannelType {
		return r.mapSplunkChannelToState(channel), nil
	}
	if channel.Kind == restapi.VictorOpsChannelType {
		return r.mapVictorOpsChannelToState(channel), nil
	}
	if channel.Kind == restapi.WebhookChannelType {
		return r.mapWebhookChannelToState(channel), nil
	}
	if channel.Kind == restapi.Office365ChannelType {
		return r.mapOffice365ChannelToState(channel), nil
	}
	if channel.Kind == restapi.GoogleChatChannelType {
		return r.mapGoogleChatChannelToState(channel), nil
	}
	return nil, fmt.Errorf("received unsupported alerting channel of type %s", channel.Kind)
}

func (r *alertingChannelResource) mapEmailChannelToState(channel *restapi.AlertingChannel) map[string]interface{} {
	return map[string]interface{}{
		AlertingChannelFieldName: channel.Name,
		AlertingChannelFieldChannelEmail: []interface{}{
			map[string]interface{}{
				AlertingChannelEmailFieldEmails: channel.Emails,
			},
		},
	}
}

func (r *alertingChannelResource) mapOpsGenieChannelToState(channel *restapi.AlertingChannel) map[string]interface{} {
	tags := r.convertCommaSeparatedListToSlice(*channel.Tags)
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

func (r *alertingChannelResource) convertCommaSeparatedListToSlice(csv string) []string {
	entries := strings.Split(csv, ",")
	result := make([]string, len(entries))
	for i, e := range entries {
		result[i] = strings.TrimSpace(e)
	}
	return result
}

func (r *alertingChannelResource) mapPagerDutyChannelToState(channel *restapi.AlertingChannel) map[string]interface{} {
	return map[string]interface{}{
		AlertingChannelFieldName: channel.Name,
		AlertingChannelFieldChannelPageDuty: []interface{}{
			map[string]interface{}{
				AlertingChannelPagerDutyFieldServiceIntegrationKey: channel.ServiceIntegrationKey,
			},
		},
	}
}

func (r *alertingChannelResource) mapSlackChannelToState(channel *restapi.AlertingChannel) map[string]interface{} {
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

func (r *alertingChannelResource) mapSplunkChannelToState(channel *restapi.AlertingChannel) map[string]interface{} {
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

func (r *alertingChannelResource) mapVictorOpsChannelToState(channel *restapi.AlertingChannel) map[string]interface{} {
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

func (r *alertingChannelResource) mapWebhookChannelToState(channel *restapi.AlertingChannel) map[string]interface{} {
	headers := r.createHTTPHeaderMapFromList(channel.Headers)
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

func (r *alertingChannelResource) createHTTPHeaderMapFromList(headers []string) map[string]interface{} {
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

func (r *alertingChannelResource) mapOffice365ChannelToState(channel *restapi.AlertingChannel) map[string]interface{} {
	return map[string]interface{}{
		AlertingChannelFieldName: channel.Name,
		AlertingChannelFieldChannelOffice365: []interface{}{
			map[string]interface{}{
				AlertingChannelWebhookBasedFieldWebhookURL: channel.WebhookURL,
			},
		},
	}
}

func (r *alertingChannelResource) mapGoogleChatChannelToState(channel *restapi.AlertingChannel) map[string]interface{} {
	return map[string]interface{}{
		AlertingChannelFieldName: channel.Name,
		AlertingChannelFieldChannelGoogleChat: []interface{}{
			map[string]interface{}{
				AlertingChannelWebhookBasedFieldWebhookURL: channel.WebhookURL,
			},
		},
	}
}

func (r *alertingChannelResource) MapStateToDataObject(d *schema.ResourceData) (*restapi.AlertingChannel, error) {
	if channel, ok := d.GetOk(AlertingChannelFieldChannelEmail); ok && len(channel.([]interface{})) == 1 {
		return r.mapStateToEmailObject(d, channel.([]interface{})[0].(map[string]interface{})), nil
	}
	if channel, ok := d.GetOk(AlertingChannelFieldChannelOpsGenie); ok && len(channel.([]interface{})) == 1 {
		return r.mapStateToOpsGenieObject(d, channel.([]interface{})[0].(map[string]interface{})), nil
	}
	if channel, ok := d.GetOk(AlertingChannelFieldChannelPageDuty); ok && len(channel.([]interface{})) == 1 {
		return r.mapStateToPagerDutyObject(d, channel.([]interface{})[0].(map[string]interface{})), nil
	}
	if channel, ok := d.GetOk(AlertingChannelFieldChannelSlack); ok && len(channel.([]interface{})) == 1 {
		return r.mapStateToSlackObject(d, channel.([]interface{})[0].(map[string]interface{})), nil
	}
	if channel, ok := d.GetOk(AlertingChannelFieldChannelSplunk); ok && len(channel.([]interface{})) == 1 {
		return r.mapStateToSplunkObject(d, channel.([]interface{})[0].(map[string]interface{})), nil
	}
	if channel, ok := d.GetOk(AlertingChannelFieldChannelVictorOps); ok && len(channel.([]interface{})) == 1 {
		return r.mapStateToVictorOpsObject(d, channel.([]interface{})[0].(map[string]interface{})), nil
	}
	if channel, ok := d.GetOk(AlertingChannelFieldChannelWebhook); ok && len(channel.([]interface{})) == 1 {
		return r.mapStateToWebhookObject(d, channel.([]interface{})[0].(map[string]interface{})), nil
	}
	if channel, ok := d.GetOk(AlertingChannelFieldChannelOffice365); ok && len(channel.([]interface{})) == 1 {
		return r.mapStateToWebhookBasedObject(restapi.Office365ChannelType, d, channel.([]interface{})[0].(map[string]interface{})), nil
	}
	if channel, ok := d.GetOk(AlertingChannelFieldChannelGoogleChat); ok && len(channel.([]interface{})) == 1 {
		return r.mapStateToWebhookBasedObject(restapi.GoogleChatChannelType, d, channel.([]interface{})[0].(map[string]interface{})), nil
	}
	return nil, fmt.Errorf("no supported alerting channel defined")
}

func (r *alertingChannelResource) mapStateToEmailObject(d *schema.ResourceData, channelState map[string]interface{}) *restapi.AlertingChannel {
	return &restapi.AlertingChannel{
		ID:     d.Id(),
		Name:   d.Get(AlertingChannelFieldName).(string),
		Kind:   restapi.EmailChannelType,
		Emails: ReadSetParameterFromMap[string](channelState, AlertingChannelEmailFieldEmails),
	}
}

func (r *alertingChannelResource) mapStateToOpsGenieObject(d *schema.ResourceData, channelState map[string]interface{}) *restapi.AlertingChannel {
	apiKey := channelState[AlertingChannelOpsGenieFieldAPIKey].(string)
	region := restapi.OpsGenieRegionType(channelState[AlertingChannelOpsGenieFieldRegion].(string))
	tags := strings.Join(ReadArrayParameterFromMap[string](channelState, AlertingChannelOpsGenieFieldTags), ",")

	return &restapi.AlertingChannel{
		ID:     d.Id(),
		Name:   d.Get(AlertingChannelFieldName).(string),
		Kind:   restapi.OpsGenieChannelType,
		APIKey: &apiKey,
		Region: &region,
		Tags:   &tags,
	}
}

func (r *alertingChannelResource) mapStateToPagerDutyObject(d *schema.ResourceData, channelState map[string]interface{}) *restapi.AlertingChannel {
	integrationKey := channelState[AlertingChannelPagerDutyFieldServiceIntegrationKey].(string)
	return &restapi.AlertingChannel{
		ID:                    d.Id(),
		Name:                  d.Get(AlertingChannelFieldName).(string),
		Kind:                  restapi.PagerDutyChannelType,
		ServiceIntegrationKey: &integrationKey,
	}
}

func (r *alertingChannelResource) mapStateToSlackObject(d *schema.ResourceData, channelState map[string]interface{}) *restapi.AlertingChannel {
	webhookURL := channelState[AlertingChannelSlackFieldWebhookURL].(string)
	iconURL := channelState[AlertingChannelSlackFieldIconURL].(string)
	channel := channelState[AlertingChannelSlackFieldChannel].(string)
	return &restapi.AlertingChannel{
		ID:         d.Id(),
		Name:       d.Get(AlertingChannelFieldName).(string),
		Kind:       restapi.SlackChannelType,
		WebhookURL: &webhookURL,
		IconURL:    &iconURL,
		Channel:    &channel,
	}
}

func (r *alertingChannelResource) mapStateToSplunkObject(d *schema.ResourceData, channelState map[string]interface{}) *restapi.AlertingChannel {
	url := channelState[AlertingChannelSplunkFieldURL].(string)
	token := channelState[AlertingChannelSplunkFieldToken].(string)
	return &restapi.AlertingChannel{
		ID:    d.Id(),
		Name:  d.Get(AlertingChannelFieldName).(string),
		Kind:  restapi.SplunkChannelType,
		URL:   &url,
		Token: &token,
	}
}

func (r *alertingChannelResource) mapStateToVictorOpsObject(d *schema.ResourceData, channelState map[string]interface{}) *restapi.AlertingChannel {
	apiKey := channelState[AlertingChannelVictorOpsFieldAPIKey].(string)
	routingKey := channelState[AlertingChannelVictorOpsFieldRoutingKey].(string)
	return &restapi.AlertingChannel{
		ID:         d.Id(),
		Name:       d.Get(AlertingChannelFieldName).(string),
		Kind:       restapi.VictorOpsChannelType,
		APIKey:     &apiKey,
		RoutingKey: &routingKey,
	}
}

func (r *alertingChannelResource) mapStateToWebhookObject(d *schema.ResourceData, channelState map[string]interface{}) *restapi.AlertingChannel {
	headers := r.createHTTPHeaderListFromMap(channelState)
	return &restapi.AlertingChannel{
		ID:          d.Id(),
		Name:        d.Get(AlertingChannelFieldName).(string),
		Kind:        restapi.WebhookChannelType,
		WebhookURLs: ReadSetParameterFromMap[string](channelState, AlertingChannelWebhookFieldWebhookURLs),
		Headers:     headers,
	}
}

func (r *alertingChannelResource) createHTTPHeaderListFromMap(channelState map[string]interface{}) []string {
	if attr, ok := channelState[AlertingChannelWebhookFieldHTTPHeaders]; ok {
		headerMap := attr.(map[string]interface{})
		result := make([]string, len(headerMap))
		i := 0
		for key, value := range headerMap {
			header := fmt.Sprintf("%s: %s", key, value)
			result[i] = header
			i++
		}

		return result
	}
	return []string{}
}

func (r *alertingChannelResource) mapStateToWebhookBasedObject(channelType restapi.AlertingChannelType, d *schema.ResourceData, channelState map[string]interface{}) *restapi.AlertingChannel {
	webhookURL := channelState[AlertingChannelWebhookBasedFieldWebhookURL].(string)
	return &restapi.AlertingChannel{
		ID:         d.Id(),
		Name:       d.Get(AlertingChannelFieldName).(string),
		Kind:       channelType,
		WebhookURL: &webhookURL,
	}
}
