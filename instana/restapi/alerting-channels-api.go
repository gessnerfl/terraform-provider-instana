package restapi

// AlertingChannelsResourcePath path to Alerting channels resource of Instana RESTful API
const AlertingChannelsResourcePath = EventSettingsBasePath + "/alertingChannels"

// AlertingChannelType type of the alerting channel
type AlertingChannelType string

const (
	//EmailChannelType constant value for alerting channel type EMAIL
	EmailChannelType = AlertingChannelType("EMAIL")
	//GoogleChatChannelType constant value for alerting channel type GOOGLE_CHAT
	GoogleChatChannelType = AlertingChannelType("GOOGLE_CHAT")
	//Office365ChannelType constant value for alerting channel type OFFICE_365
	Office365ChannelType = AlertingChannelType("OFFICE_365")
	//OpsGenieChannelType constant value for alerting channel type OPS_GENIE
	OpsGenieChannelType = AlertingChannelType("OPS_GENIE")
	//PagerDutyChannelType constant value for alerting channel type PAGER_DUTY
	PagerDutyChannelType = AlertingChannelType("PAGER_DUTY")
	//SlackChannelType constant value for alerting channel type SLACK
	SlackChannelType = AlertingChannelType("SLACK")
	//SplunkChannelType constant value for alerting channel type SPLUNK
	SplunkChannelType = AlertingChannelType("SPLUNK")
	//VictorOpsChannelType constant value for alerting channel type VICTOR_OPS
	VictorOpsChannelType = AlertingChannelType("VICTOR_OPS")
	//WebhookChannelType constant value for alerting channel type WEB_HOOK
	WebhookChannelType = AlertingChannelType("WEB_HOOK")
)

// SupportedAlertingChannels list of supported calerting channels of Instana API
var SupportedAlertingChannels = []AlertingChannelType{
	EmailChannelType,
	GoogleChatChannelType,
	Office365ChannelType,
	OpsGenieChannelType,
	PagerDutyChannelType,
	SlackChannelType,
	SplunkChannelType,
	VictorOpsChannelType,
	WebhookChannelType,
}

// OpsGenieRegionType type of the OpsGenie region
type OpsGenieRegionType string

const (
	//EuOpsGenieRegion constatnt field for OpsGenie region type EU
	EuOpsGenieRegion = OpsGenieRegionType("EU")
	//UsOpsGenieRegion constatnt field for OpsGenie region type US
	UsOpsGenieRegion = OpsGenieRegionType("US")
)

// SupportedOpsGenieRegions list of supported OpsGenie regions of Instana API
var SupportedOpsGenieRegions = []OpsGenieRegionType{EuOpsGenieRegion, UsOpsGenieRegion}

// IsSupportedOpsGenieRegionType checks if the given OpsGenie region is supported by Instana
func IsSupportedOpsGenieRegionType(regionType OpsGenieRegionType) bool {
	return isInOpsGenieRegionTypeSlice(SupportedOpsGenieRegions, regionType)
}

func isInOpsGenieRegionTypeSlice(allRegionTypes []OpsGenieRegionType, regionType OpsGenieRegionType) bool {
	for _, v := range allRegionTypes {
		if v == regionType {
			return true
		}
	}
	return false
}

// AlertingChannel is the representation of an alerting channel in Instana
type AlertingChannel struct {
	ID                    string              `json:"id"`
	Name                  string              `json:"name"`
	Kind                  AlertingChannelType `json:"kind"`
	Emails                []string            `json:"emails"`
	WebhookURL            *string             `json:"webhookUrl"`
	APIKey                *string             `json:"apiKey"`
	Tags                  *string             `json:"tags"`
	Region                *OpsGenieRegionType `json:"region"`
	RoutingKey            *string             `json:"routingKey"`
	ServiceIntegrationKey *string             `json:"serviceIntegrationKey"`
	IconURL               *string             `json:"iconUrl"`
	Channel               *string             `json:"channel"`
	URL                   *string             `json:"url"`
	Token                 *string             `json:"token"`
	WebhookURLs           []string            `json:"webhookUrls"`
	Headers               []string            `json:"headers"`
}

// GetIDForResourcePath implementation of the interface InstanaDataObject
func (r *AlertingChannel) GetIDForResourcePath() string {
	return r.ID
}
