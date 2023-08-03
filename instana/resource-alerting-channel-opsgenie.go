package instana

import (
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"strings"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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

// NewAlertingChannelOpsGenieResourceHandle creates the resource handle for Alerting Channels of type OpsGenie
func NewAlertingChannelOpsGenieResourceHandle() ResourceHandle[*restapi.AlertingChannel] {
	opsGenieRegions := make([]string, len(restapi.SupportedOpsGenieRegions))
	for i, r := range restapi.SupportedOpsGenieRegions {
		opsGenieRegions[i] = string(r)
	}

	return &alertingChannelOpsGenieResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaAlertingChannelOpsGenie,
			Schema: map[string]*schema.Schema{
				AlertingChannelFieldName:     alertingChannelNameSchemaField,
				AlertingChannelFieldFullName: alertingChannelFullNameSchemaField,
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
					ValidateFunc: validation.StringInSlice(opsGenieRegions, false),
					Description:  fmt.Sprintf("The OpsGenie region (%s) of the OpsGenie alerting channel", strings.Join(opsGenieRegions, "/")),
				},
			},
		},
	}
}

type alertingChannelOpsGenieResource struct {
	metaData ResourceMetaData
}

func (r *alertingChannelOpsGenieResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *alertingChannelOpsGenieResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{}
}

func (r *alertingChannelOpsGenieResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.AlertingChannel] {
	return api.AlertingChannels()
}

func (r *alertingChannelOpsGenieResource) SetComputedFields(_ *schema.ResourceData) error {
	return nil
}

func (r *alertingChannelOpsGenieResource) UpdateState(d *schema.ResourceData, alertingChannel *restapi.AlertingChannel, formatter utils.ResourceNameFormatter) error {
	tags := r.convertCommaSeparatedListToSlice(*alertingChannel.Tags)
	d.SetId(alertingChannel.ID)
	return tfutils.UpdateState(d, map[string]interface{}{
		AlertingChannelFieldName:           formatter.UndoFormat(alertingChannel.Name),
		AlertingChannelFieldFullName:       alertingChannel.Name,
		AlertingChannelOpsGenieFieldAPIKey: alertingChannel.APIKey,
		AlertingChannelOpsGenieFieldRegion: alertingChannel.Region,
		AlertingChannelOpsGenieFieldTags:   tags,
	})
}

func (r *alertingChannelOpsGenieResource) MapStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (*restapi.AlertingChannel, error) {
	name := computeFullAlertingChannelNameString(d, formatter)
	apiKey := d.Get(AlertingChannelOpsGenieFieldAPIKey).(string)
	region := restapi.OpsGenieRegionType(d.Get(AlertingChannelOpsGenieFieldRegion).(string))
	tags := strings.Join(ReadStringArrayParameterFromResource(d, AlertingChannelOpsGenieFieldTags), ",")

	return &restapi.AlertingChannel{
		ID:     d.Id(),
		Name:   name,
		Kind:   restapi.OpsGenieChannelType,
		APIKey: &apiKey,
		Region: &region,
		Tags:   &tags,
	}, nil
}

func (r *alertingChannelOpsGenieResource) convertCommaSeparatedListToSlice(csv string) []string {
	entries := strings.Split(csv, ",")
	result := make([]string, len(entries))
	for i, e := range entries {
		result[i] = strings.TrimSpace(e)
	}
	return result
}
