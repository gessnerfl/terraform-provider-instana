package instana

import (
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"strings"
)

const (
	//AlertingChannelOpsGenieFieldAPIKey const for the api key field of the alerting channel OpsGenie
	AlertingChannelOpsGenieFieldAPIKey = "api_key"
	//AlertingChannelOpsGenieFieldTags const for the tags field of the alerting channel OpsGenie
	AlertingChannelOpsGenieFieldTags = "tags"
	//AlertingChannelOpsGenieFieldRegion const for the region field of the alerting channel OpsGenie
	AlertingChannelOpsGenieFieldRegion = "region"
	//ResourceInstanaAlertingChannelOpsGenie the name of the terraform-provider-instana resource to manage alerting channels of type OpsGenie
	ResourceInstanaAlertingChannelOpsGenie = "instana_alerting_channel_ops_genie"
)

var opsGenieRegions = convertOpsGenieRegionsToStringSlice()

//NewAlertingChannelOpsGenieResourceHandle creates the resource handle for Alerting Channels of type OpsGenie
func NewAlertingChannelOpsGenieResourceHandle() ResourceHandle {
	return &alertingChannelOpsGenieResourceHandle{}
}

type alertingChannelOpsGenieResourceHandle struct {
}

func (h *alertingChannelOpsGenieResourceHandle) GetResource(api restapi.InstanaAPI) restapi.RestResource {
	return api.AlertingChannels()
}

func (h *alertingChannelOpsGenieResourceHandle) GetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		AlertingChannelFieldName:     alertingChannelNameSchemaField,
		AlertingChannelFieldFullName: alertingChannelFullNameSchemaField,
		AlertingChannelOpsGenieFieldAPIKey: &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "The OpsGenie API Key of the OpsGenie alerting channel",
		},
		AlertingChannelOpsGenieFieldTags: &schema.Schema{
			Type:     schema.TypeList,
			MinItems: 1,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Required:    true,
			Description: "The OpsGenie tags of the OpsGenie alerting channel",
		},
		AlertingChannelOpsGenieFieldRegion: &schema.Schema{
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(opsGenieRegions, false),
			Description:  fmt.Sprintf("The OpsGenie region (%s) of the OpsGenie alerting channel", strings.Join(opsGenieRegions, "/")),
		},
	}
}

func (h *alertingChannelOpsGenieResourceHandle) GetResourceName() string {
	return ResourceInstanaAlertingChannelOpsGenie
}

func (h *alertingChannelOpsGenieResourceHandle) UpdateState(d *schema.ResourceData, obj restapi.InstanaDataObject) {
	alertingChannel := obj.(restapi.AlertingChannel)
	tags := convertCommaSeparatedListToSlice(*alertingChannel.Tags)
	d.Set(AlertingChannelFieldFullName, alertingChannel.Name)
	d.Set(AlertingChannelOpsGenieFieldAPIKey, alertingChannel.APIKey)
	d.Set(AlertingChannelOpsGenieFieldRegion, alertingChannel.Region)
	d.Set(AlertingChannelOpsGenieFieldTags, tags)
	d.SetId(alertingChannel.ID)
}

func (h *alertingChannelOpsGenieResourceHandle) ConvertStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) restapi.InstanaDataObject {
	name := computeFullAlertingChannelNameString(d, formatter)
	apiKey := d.Get(AlertingChannelOpsGenieFieldAPIKey).(string)
	region := restapi.OpsGenieRegionType(d.Get(AlertingChannelOpsGenieFieldRegion).(string))
	tags := strings.Join(ReadStringArrayParameterFromResource(d, AlertingChannelOpsGenieFieldTags), ",")

	return restapi.AlertingChannel{
		ID:     d.Id(),
		Name:   name,
		Kind:   restapi.OpsGenieChannelType,
		APIKey: &apiKey,
		Region: &region,
		Tags:   &tags,
	}
}

func convertCommaSeparatedListToSlice(csv string) []string {
	entries := strings.Split(csv, ",")
	result := make([]string, len(entries))
	for i, e := range entries {
		result[i] = strings.TrimSpace(e)
	}
	return result
}

func convertOpsGenieRegionsToStringSlice() []string {
	result := make([]string, len(restapi.SupportedOpsGenieRegions))
	for i, r := range restapi.SupportedOpsGenieRegions {
		result[i] = string(r)
	}
	return result
}
