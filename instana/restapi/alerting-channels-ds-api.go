package restapi

import (
	"errors"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/utils"
)

// AlertingChannelDS is embedding AlertingChannel so we could satisfy the restapi.InstanaDataObject interface
// for datasource read operations.
type AlertingChannelDS struct {
	AlertingChannel
}

// GetIDForResourcePath implementation of the interface InstanaDataObject
func (r *AlertingChannelDS) GetIDForResourcePath() string {
	return r.ID
}

// Validate implementation of the interface InstanaDataObject to verify if data object is correct
func (r *AlertingChannelDS) Validate() error {
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

func (r *AlertingChannelDS) validateEmailIntegration() error {
	return acValidateList(r.Emails, "Email addresses are missing")
}

func (r *AlertingChannelDS) validateWebHookBasedIntegrations() error {
	return acValidateOpt(r.WebhookURL, "Webhook URL is missing")
}

func (r *AlertingChannelDS) validateOpsGenieIntegration() error {
	return acValidateOpsGenieIntegration(r.APIKey, r.Tags, r.Region)
}

func (r *AlertingChannelDS) validatePagerDutyIntegration() error {
	return acValidateOpt(r.ServiceIntegrationKey, "Service integration key is missing")
}

func (r *AlertingChannelDS) validateSplunkIntegration() error {
	m := make(map[string]*string)
	m["URL is missing"] = r.URL
	m["Token is missing"] = r.Token
	return acValidateOpts(m)
}

func (r *AlertingChannelDS) validateVictorOpsIntegration() error {
	m := make(map[string]*string)
	m["API Key is missing"] = r.APIKey
	m["Routing Key is missing"] = r.RoutingKey
	return acValidateOpts(m)
}

func (r *AlertingChannelDS) validateGenericWebHookIntegration() error {
	return acValidateList(r.WebhookURLs, "Webhook URLs are missing")
}
