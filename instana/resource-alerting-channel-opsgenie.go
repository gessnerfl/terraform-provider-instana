package instana

import (
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"strings"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
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

func createOpsGenieRegionSlice() []string {
	opsGenieRegions := make([]string, len(restapi.SupportedOpsGenieRegions))
	for i, r := range restapi.SupportedOpsGenieRegions {
		opsGenieRegions[i] = string(r)
	}
	return opsGenieRegions
}

var (
	opsGenieRegions                     = createOpsGenieRegionSlice()
	alertingChannelOpsGenieSchemaAPIKey = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The OpsGenie API Key of the OpsGenie alerting channel",
		Sensitive:   true,
	}
	alertingChannelOpsGenieSchemaTags = &schema.Schema{
		Type:     schema.TypeList,
		MinItems: 1,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Required:    true,
		Description: "The OpsGenie tags of the OpsGenie alerting channel",
	}
	alertingChannelOpsGenieSchemaRegion = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice(opsGenieRegions, false),
		Description:  fmt.Sprintf("The OpsGenie region (%s) of the OpsGenie alerting channel", strings.Join(opsGenieRegions, "/")),
	}
)

// NewAlertingChannelOpsGenieResourceHandle creates the resource handle for Alerting Channels of type OpsGenie
func NewAlertingChannelOpsGenieResourceHandle() ResourceHandle[*restapi.AlertingChannel] {
	return &alertingChannelOpsGenieResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaAlertingChannelOpsGenie,
			Schema: map[string]*schema.Schema{
				AlertingChannelFieldName:           alertingChannelNameSchemaField,
				AlertingChannelOpsGenieFieldAPIKey: alertingChannelOpsGenieSchemaAPIKey,
				AlertingChannelOpsGenieFieldTags:   alertingChannelOpsGenieSchemaTags,
				AlertingChannelOpsGenieFieldRegion: alertingChannelOpsGenieSchemaRegion,
			},
			SchemaVersion: 1,
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
	return []schema.StateUpgrader{
		{
			Type:    r.schemaV0().CoreConfigSchema().ImpliedType(),
			Upgrade: migrateFullNameToName,
			Version: 0,
		},
	}
}

func (r *alertingChannelOpsGenieResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.AlertingChannel] {
	return api.AlertingChannels()
}

func (r *alertingChannelOpsGenieResource) SetComputedFields(_ *schema.ResourceData) error {
	return nil
}

func (r *alertingChannelOpsGenieResource) UpdateState(d *schema.ResourceData, alertingChannel *restapi.AlertingChannel) error {
	tags := r.convertCommaSeparatedListToSlice(*alertingChannel.Tags)
	d.SetId(alertingChannel.ID)
	return tfutils.UpdateState(d, map[string]interface{}{
		AlertingChannelFieldName:           alertingChannel.Name,
		AlertingChannelOpsGenieFieldAPIKey: alertingChannel.APIKey,
		AlertingChannelOpsGenieFieldRegion: alertingChannel.Region,
		AlertingChannelOpsGenieFieldTags:   tags,
	})
}

func (r *alertingChannelOpsGenieResource) MapStateToDataObject(d *schema.ResourceData) (*restapi.AlertingChannel, error) {
	apiKey := d.Get(AlertingChannelOpsGenieFieldAPIKey).(string)
	region := restapi.OpsGenieRegionType(d.Get(AlertingChannelOpsGenieFieldRegion).(string))
	tags := strings.Join(ReadStringArrayParameterFromResource(d, AlertingChannelOpsGenieFieldTags), ",")

	return &restapi.AlertingChannel{
		ID:     d.Id(),
		Name:   d.Get(AlertingChannelFieldName).(string),
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

func (r *alertingChannelOpsGenieResource) schemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			AlertingChannelFieldName:           alertingChannelNameSchemaField,
			AlertingChannelFieldFullName:       alertingChannelFullNameSchemaField,
			AlertingChannelOpsGenieFieldAPIKey: alertingChannelOpsGenieSchemaAPIKey,
			AlertingChannelOpsGenieFieldTags:   alertingChannelOpsGenieSchemaTags,
			AlertingChannelOpsGenieFieldRegion: alertingChannelOpsGenieSchemaRegion,
		},
	}
}
