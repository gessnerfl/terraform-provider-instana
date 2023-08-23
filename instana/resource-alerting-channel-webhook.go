package instana

import (
	"context"
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"strings"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	//AlertingChannelWebhookFieldWebhookURLs const for the webhooks field of the Webhook alerting channel
	AlertingChannelWebhookFieldWebhookURLs = "webhook_urls"
	//AlertingChannelWebhookFieldHTTPHeaders const for the http headers field of the Webhook alerting channel
	AlertingChannelWebhookFieldHTTPHeaders = "http_headers"
	//ResourceInstanaAlertingChannelWebhook the name of the terraform-provider-instana resource to manage alerting channels of type webhook
	ResourceInstanaAlertingChannelWebhook = "instana_alerting_channel_webhook"
)

// AlertingChannelWebhookWebhookURLsSchemaField schema field definition of instana_alerting_channel_webhook field webhook_urls
var AlertingChannelWebhookWebhookURLsSchemaField = &schema.Schema{
	Type:     schema.TypeSet,
	MinItems: 1,
	Elem: &schema.Schema{
		Type: schema.TypeString,
	},
	Required:    true,
	Description: "The list of webhook urls of the Webhook alerting channel",
}

// AlertingChannelWebhookHTTPHeadersSchemaField schema field definition of instana_alerting_channel_webhook field http_headers
var AlertingChannelWebhookHTTPHeadersSchemaField = &schema.Schema{
	Type: schema.TypeMap,
	Elem: &schema.Schema{
		Type: schema.TypeString,
	},
	Optional:    true,
	Description: "The optional map of HTTP headers of the Webhook alerting channel",
}

// NewAlertingChannelWebhookResourceHandle creates the resource handle for Alerting Channels of type Webhook
func NewAlertingChannelWebhookResourceHandle() ResourceHandle[*restapi.AlertingChannel] {
	return &alertingChannelWebhookResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaAlertingChannelWebhook,
			Schema: map[string]*schema.Schema{
				AlertingChannelFieldName:               alertingChannelNameSchemaField,
				AlertingChannelWebhookFieldWebhookURLs: AlertingChannelWebhookWebhookURLsSchemaField,
				AlertingChannelWebhookFieldHTTPHeaders: AlertingChannelWebhookHTTPHeadersSchemaField,
			},
			SchemaVersion: 2,
		},
	}
}

type alertingChannelWebhookResource struct {
	metaData ResourceMetaData
}

func (r *alertingChannelWebhookResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *alertingChannelWebhookResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{
		{
			Type: r.schemaV0().CoreConfigSchema().ImpliedType(),
			Upgrade: func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
				return rawState, nil
			},
			Version: 0,
		},
		{
			Type:    r.schemaV1().CoreConfigSchema().ImpliedType(),
			Upgrade: migrateFullNameToName,
			Version: 1,
		},
	}
}

func (r *alertingChannelWebhookResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.AlertingChannel] {
	return api.AlertingChannels()
}

func (r *alertingChannelWebhookResource) SetComputedFields(_ *schema.ResourceData) error {
	return nil
}

func (r *alertingChannelWebhookResource) UpdateState(d *schema.ResourceData, alertingChannel *restapi.AlertingChannel) error {
	urls := alertingChannel.WebhookURLs
	headers := r.createHTTPHeaderMapFromList(alertingChannel.Headers)
	d.SetId(alertingChannel.ID)
	return tfutils.UpdateState(d, map[string]interface{}{
		AlertingChannelFieldName:               alertingChannel.Name,
		AlertingChannelWebhookFieldWebhookURLs: urls,
		AlertingChannelWebhookFieldHTTPHeaders: headers,
	})
}

func (r *alertingChannelWebhookResource) createHTTPHeaderMapFromList(headers []string) map[string]interface{} {
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

func (r *alertingChannelWebhookResource) MapStateToDataObject(d *schema.ResourceData) (*restapi.AlertingChannel, error) {
	headers := r.createHTTPHeaderListFromMap(d)
	return &restapi.AlertingChannel{
		ID:          d.Id(),
		Name:        d.Get(AlertingChannelFieldName).(string),
		Kind:        restapi.WebhookChannelType,
		WebhookURLs: ReadStringSetParameterFromResource(d, AlertingChannelWebhookFieldWebhookURLs),
		Headers:     headers,
	}, nil
}

func (r *alertingChannelWebhookResource) createHTTPHeaderListFromMap(d *schema.ResourceData) []string {
	if attr, ok := d.GetOk(AlertingChannelWebhookFieldHTTPHeaders); ok {
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

func (r *alertingChannelWebhookResource) schemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			AlertingChannelFieldName:     alertingChannelNameSchemaField,
			AlertingChannelFieldFullName: alertingChannelFullNameSchemaField,
			AlertingChannelWebhookFieldWebhookURLs: {
				Type:     schema.TypeList,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "The list of webhook urls of the Webhook alerting channel",
			},
			AlertingChannelWebhookFieldHTTPHeaders: AlertingChannelWebhookHTTPHeadersSchemaField,
		},
	}
}

func (r *alertingChannelWebhookResource) schemaV1() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			AlertingChannelFieldName:               alertingChannelNameSchemaField,
			AlertingChannelFieldFullName:           alertingChannelFullNameSchemaField,
			AlertingChannelWebhookFieldWebhookURLs: AlertingChannelWebhookWebhookURLsSchemaField,
			AlertingChannelWebhookFieldHTTPHeaders: AlertingChannelWebhookHTTPHeadersSchemaField,
		},
	}
}
