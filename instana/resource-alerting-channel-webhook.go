package instana

import (
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

const (
	//AlertingChannelWebhookFieldWebhookURLs const for the webhooks field of the Webhook alerting channel
	AlertingChannelWebhookFieldWebhookURLs = "webhook_urls"
	//AlertingChannelWebhookFieldHttpHeaders const for the http headers field of the Webhook alerting channel
	AlertingChannelWebhookFieldHttpHeaders = "http_headers"
)

//NewAlertingChannelWebhookResource creates the terraform resource for Alerting Channels of type Webhook
func NewAlertingChannelWebhookResource() TerraformResource {
	return NewTerraformResource(NewAlertingChannelWebhookResourceHandle())
}

//NewAlertingChannelWebhookResourceHandle creates the resource handle for Alerting Channels of type Webhook
func NewAlertingChannelWebhookResourceHandle() ResourceHandle {
	return &alertingChannelWebhookResourceHandle{}
}

type alertingChannelWebhookResourceHandle struct {
}

func (h *alertingChannelWebhookResourceHandle) GetResource(api restapi.InstanaAPI) restapi.RestResource {
	return api.AlertingChannels()
}

func (h *alertingChannelWebhookResourceHandle) GetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
		AlertingChannelWebhookFieldHttpHeaders: &schema.Schema{
			Type: schema.TypeMap,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional:    true,
			Description: "The optional map of HTTP headers of the Webhook alerting channel",
		},
	}
}

func (h *alertingChannelWebhookResourceHandle) GetResourceName() string {
	return "instana_alerting_channel_webhook"
}

func (h *alertingChannelWebhookResourceHandle) UpdateState(d *schema.ResourceData, obj restapi.InstanaDataObject) {
	alertingChannel := obj.(restapi.AlertingChannel)
	urls := alertingChannel.WebhookURLs
	headers := h.createHeaderMapFromList(alertingChannel.Headers)
	d.Set(AlertingChannelFieldFullName, alertingChannel.Name)
	d.Set(AlertingChannelWebhookFieldWebhookURLs, urls)
	d.Set(AlertingChannelWebhookFieldHttpHeaders, headers)
	d.SetId(alertingChannel.ID)
}

func (h *alertingChannelWebhookResourceHandle) createHeaderMapFromList(headers []string) map[string]interface{} {
	result := make(map[string]interface{})
	for _, header := range headers {
		keyValue := strings.Split(header, ":")
		if len(keyValue) == 2 {
			result[strings.TrimSpace(keyValue[0])] = strings.TrimSpace(keyValue[1])
		} else {
			result[header] = "<undefined>"
		}
	}
	return result
}

func (h *alertingChannelWebhookResourceHandle) ConvertStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) restapi.InstanaDataObject {
	name := computeFullAlertingChannelNameString(d, formatter)
	headers := h.createHeaderListFromMap(d)
	return restapi.AlertingChannel{
		ID:          d.Id(),
		Name:        name,
		Kind:        restapi.WebhookChannelType,
		WebhookURLs: ReadStringArrayParameterFromResource(d, AlertingChannelWebhookFieldWebhookURLs),
		Headers:     headers,
	}
}

func (h *alertingChannelWebhookResourceHandle) createHeaderListFromMap(d *schema.ResourceData) []string {
	if attr, ok := d.GetOk(AlertingChannelWebhookFieldHttpHeaders); ok {
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
	return nil
}
