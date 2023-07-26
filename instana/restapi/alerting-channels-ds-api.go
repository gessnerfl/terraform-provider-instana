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

// GetIDForResourcePath implemention of the interface InstanaDataObject
func (r AlertingChannelDS) GetIDForResourcePath() string {
	return r.ID
}

// Validate implementation of the interface InstanaDataObject to verify if data object is correct
func (r AlertingChannelDS) Validate() error {
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

func (r AlertingChannelDS) validateEmailIntegration() error {
	if len(r.Emails) == 0 {
		return errors.New("Email addresses are missing")
	}
	return nil
}

func (r AlertingChannelDS) validateWebHookBasedIntegrations() error {
	if r.WebhookURL == nil || utils.IsBlank(*r.WebhookURL) {
		return errors.New("Webhook URL is missing")
	}
	return nil
}

func (r AlertingChannelDS) validateOpsGenieIntegration() error {
	if r.APIKey == nil || utils.IsBlank(*r.APIKey) {
		return errors.New("API key is missing")
	}
	if r.Tags == nil || utils.IsBlank(*r.Tags) {
		return errors.New("Tags are missing")
	}
	if r.Region == nil {
		return errors.New("Region is missing")
	}
	if !IsSupportedOpsGenieRegionType(*r.Region) {
		return fmt.Errorf("Region %s is not valid", *r.Region)
	}
	return nil
}

func (r AlertingChannelDS) validatePagerDutyIntegration() error {
	if r.ServiceIntegrationKey == nil || utils.IsBlank(*r.ServiceIntegrationKey) {
		return errors.New("Service integration key is missing")
	}
	return nil
}

func (r AlertingChannelDS) validateSplunkIntegration() error {
	if r.URL == nil || utils.IsBlank(*r.URL) {
		return errors.New("URL is missing")
	}
	if r.Token == nil || utils.IsBlank(*r.Token) {
		return errors.New("Token is missing")
	}
	return nil
}

func (r AlertingChannelDS) validateVictorOpsIntegration() error {
	if r.APIKey == nil || utils.IsBlank(*r.APIKey) {
		return errors.New("API Key is missing")
	}
	if r.RoutingKey == nil || utils.IsBlank(*r.RoutingKey) {
		return errors.New("Routing Key is missing")
	}
	return nil
}

func (r AlertingChannelDS) validateGenericWebHookIntegration() error {
	if len(r.WebhookURLs) == 0 {
		return errors.New("Webhook URLs are missing")
	}
	return nil
}
