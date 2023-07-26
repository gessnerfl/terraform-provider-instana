package restapi

import (
	"errors"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/utils"
)

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

// GetIDForResourcePath implemention of the interface InstanaDataObject
func (r *AlertingChannel) GetIDForResourcePath() string {
	return r.ID
}

// Validate implementation of the interface InstanaDataObject to verify if data object is correct
func (r *AlertingChannel) Validate() error {
	if utils.IsBlank(r.ID) {
		return errors.New("ID is missing")
	}
	if utils.IsBlank(r.Name) {
		return errors.New("Name is missing")
	}
	if len(r.Kind) == 0 {
		return errors.New("Kind is missing")
	}

	switch r.Kind {
	case EmailChannelType:
		return r.validateEmailIntegration()
	case GoogleChatChannelType, Office365ChannelType, SlackChannelType:
		return r.validateWebHookBasedIntegrations()
	case OpsGenieChannelType:
		return r.validateOpsGenieIntegration()
	case PagerDutyChannelType:
		return r.validatePagerDutyIntegration()
	case SplunkChannelType:
		return r.validateSplunkIntegration()
	case VictorOpsChannelType:
		return r.validateVictorOpsIntegration()
	case WebhookChannelType:
		return r.validateGenericWebHookIntegration()
	default:
		return fmt.Errorf("unsupported alerting channel type %s", r.Kind)
	}
}

func (r *AlertingChannel) validateEmailIntegration() error {
	return acValidateList(r.Emails, "Email addresses are missing")
}

func (r *AlertingChannel) validateWebHookBasedIntegrations() error {
	return acValidateOpt(r.WebhookURL, "Webhook URL is missing")
}

func (r *AlertingChannel) validateOpsGenieIntegration() error {
	return acValidateOpsGenieIntegration(r.APIKey, r.Tags, r.Region)
}

func (r *AlertingChannel) validatePagerDutyIntegration() error {
	return acValidateOpt(r.ServiceIntegrationKey, "Service integration key is missing")
}

func (r *AlertingChannel) validateSplunkIntegration() error {
	m := make(map[string]*string)
	m["URL is missing"] = r.URL
	m["Token is missing"] = r.Token
	return acValidateOpts(m)
}

func (r *AlertingChannel) validateVictorOpsIntegration() error {
	m := make(map[string]*string)
	m["API Key is missing"] = r.APIKey
	m["Routing Key is missing"] = r.RoutingKey
	return acValidateOpts(m)
}

func (r *AlertingChannel) validateGenericWebHookIntegration() error {
	return acValidateList(r.WebhookURLs, "Webhook URLs are missing")
}
