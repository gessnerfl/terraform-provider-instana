package instana

import (
	"fmt"
	"strings"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	//AlertingChannelWebhookFieldWebhookURLs const for the webhooks field of the Webhook alerting channel
	AlertingChannelWebhookFieldWebhookURLs = "webhook_urls"
	//AlertingChannelWebhookFieldHTTPHeaders const for the http headers field of the Webhook alerting channel
	AlertingChannelWebhookFieldHTTPHeaders = "http_headers"
	//ResourceInstanaAlertingChannelWebhook the name of the terraform-provider-instana resource to manage alerting channels of type webhook
	ResourceInstanaAlertingChannelWebhook = "instana_alerting_channel_webhook"
)

//AlertingChannelWebhookWebhookURLsSchemaField schema field definition of instana_alerting_channel_webhook field webhook_urls
var AlertingChannelWebhookWebhookURLsSchemaField = &schema.Schema{
	Type:     schema.TypeSet,
	MinItems: 1,
	Elem: &schema.Schema{
		Type: schema.TypeString,
	},
	Required:    true,
	Description: "The list of webhook urls of the Webhook alerting channel",
}

//AlertingChannelWebhookHTTPHeadersSchemaField schema field definition of instana_alerting_channel_webhook field http_headers
var AlertingChannelWebhookHTTPHeadersSchemaField = &schema.Schema{
	Type: schema.TypeMap,
	Elem: &schema.Schema{
		Type: schema.TypeString,
	},
	Optional:    true,
	Description: "The optional map of HTTP headers of the Webhook alerting channel",
}

//NewAlertingChannelWebhookResourceHandle creates the resource handle for Alerting Channels of type Webhook
func NewAlertingChannelWebhookResourceHandle() *ResourceHandle {
	return &ResourceHandle{
		ResourceName: ResourceInstanaAlertingChannelWebhook,
		Schema: map[string]*schema.Schema{
			AlertingChannelFieldName:               alertingChannelNameSchemaField,
			AlertingChannelFieldFullName:           alertingChannelFullNameSchemaField,
			AlertingChannelWebhookFieldWebhookURLs: AlertingChannelWebhookWebhookURLsSchemaField,
			AlertingChannelWebhookFieldHTTPHeaders: AlertingChannelWebhookHTTPHeadersSchemaField,
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type: alertingChannelWebhookSchemaV0().CoreConfigSchema().ImpliedType(),
				Upgrade: func(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
					return rawState, nil
				},
				Version: 0,
			},
		},
		RestResourceFactory:  func(api restapi.InstanaAPI) restapi.RestResource { return api.AlertingChannels() },
		UpdateState:          updateStateForAlertingChannelWebhook,
		MapStateToDataObject: mapStateToDataObjectForAlertingChannelWebhook,
	}
}

func updateStateForAlertingChannelWebhook(d *schema.ResourceData, obj restapi.InstanaDataObject) error {
	alertingChannel := obj.(restapi.AlertingChannel)
	urls := alertingChannel.WebhookURLs
	headers := createHTTPHeaderMapFromList(alertingChannel.Headers)
	d.Set(AlertingChannelFieldFullName, alertingChannel.Name)
	d.Set(AlertingChannelWebhookFieldWebhookURLs, urls)
	d.Set(AlertingChannelWebhookFieldHTTPHeaders, headers)
	d.SetId(alertingChannel.ID)
	return nil
}

func createHTTPHeaderMapFromList(headers []string) map[string]interface{} {
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

func mapStateToDataObjectForAlertingChannelWebhook(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	name := computeFullAlertingChannelNameString(d, formatter)
	headers := createHTTPHeaderListFromMap(d)
	return restapi.AlertingChannel{
		ID:          d.Id(),
		Name:        name,
		Kind:        restapi.WebhookChannelType,
		WebhookURLs: ReadStringSetParameterFromResource(d, AlertingChannelWebhookFieldWebhookURLs),
		Headers:     headers,
	}, nil
}

func createHTTPHeaderListFromMap(d *schema.ResourceData) []string {
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

func alertingChannelWebhookSchemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			AlertingChannelFieldName:     alertingChannelNameSchemaField,
			AlertingChannelFieldFullName: alertingChannelFullNameSchemaField,
			AlertingChannelWebhookFieldWebhookURLs: &schema.Schema{
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
