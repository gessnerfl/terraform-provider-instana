package restapi_test

import (
	"fmt"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/assert"
)

const (
	nameFieldValue                  = "name"
	email1FieldValue                = "email1"
	email2FieldValue                = "email2"
	apiKeyFieldValue                = "apiKey"
	tagsFieldValue                  = "tag1, tag2"
	serviceIntegrationKeyFieldValue = "serviceIntegrationKey"
	urlFieldValue                   = "urlFieldValue"
	tokenFieldValue                 = "tokenFieldValue"
	routingKeyFieldValue            = "routingKeyFieldValue"
)

func TestShouldReturnIDOfAlteringChannel(t *testing.T) {
	alertingChannel := AlertingChannel{
		ID:     idFieldValue,
		Name:   nameFieldValue,
		Kind:   EmailChannelType,
		Emails: []string{email1FieldValue, email2FieldValue},
	}

	assert.Equal(t, idFieldValue, alertingChannel.GetIDForResourcePath())
}

func TestShouldSuccussullyValidateConsistentEmailAlteringChannel(t *testing.T) {
	alertingChannel := AlertingChannel{
		ID:     idFieldValue,
		Name:   nameFieldValue,
		Kind:   EmailChannelType,
		Emails: []string{email1FieldValue, email2FieldValue},
	}

	err := alertingChannel.Validate()
	assert.Nil(t, err)
}

func TestShouldFailToValidateAlteringChannelWhenIdIsMissing(t *testing.T) {
	alertingChannel := AlertingChannel{
		Name:   nameFieldValue,
		Kind:   EmailChannelType,
		Emails: []string{email1FieldValue, email2FieldValue},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "ID")
}

func TestShouldFailToValidateAlteringChannelWhenIdIsBlank(t *testing.T) {
	alertingChannel := AlertingChannel{
		ID:     " ",
		Name:   nameFieldValue,
		Kind:   EmailChannelType,
		Emails: []string{email1FieldValue, email2FieldValue},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "ID")
}

func TestShouldFailToValidateAlteringChannelWhenNameIsMissing(t *testing.T) {
	alertingChannel := AlertingChannel{
		ID:     idFieldValue,
		Kind:   EmailChannelType,
		Emails: []string{email1FieldValue, email2FieldValue},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Name")
}

func TestShouldFailToValidateAlteringChannelWhenNameIsBlank(t *testing.T) {
	name := " "
	alertingChannel := AlertingChannel{
		ID:     idFieldValue,
		Kind:   EmailChannelType,
		Name:   name,
		Emails: []string{email1FieldValue, email2FieldValue},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Name")
}

func TestShouldFailToValidateAlteringChannelWhenKindIsMissing(t *testing.T) {
	alertingChannel := AlertingChannel{
		ID:     idFieldValue,
		Name:   nameFieldValue,
		Emails: []string{email1FieldValue, email2FieldValue},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Kind")
}

func TestShouldFailToValidateAlteringChannelWhenKindIsNotValid(t *testing.T) {
	alertingChannel := AlertingChannel{
		ID:     idFieldValue,
		Name:   nameFieldValue,
		Kind:   AlertingChannelType("invalid"),
		Emails: []string{email1FieldValue, email2FieldValue},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unsupported alerting channel type")
}
func TestShouldFailToValidateEmailAlteringChannelWhenNoEmailIsProvided(t *testing.T) {
	alertingChannel := AlertingChannel{
		ID:     idFieldValue,
		Name:   nameFieldValue,
		Kind:   EmailChannelType,
		Emails: []string{},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Email addresses")
}

func TestShouldSuccussullyValidateConsistentWebhhokBasedAlteringChannel(t *testing.T) {
	for _, channelType := range []AlertingChannelType{GoogleChatChannelType, Office365ChannelType, SlackChannelType} {
		t.Run(fmt.Sprintf("TestShouldSuccussullyValidateConsistentWebhhokBasedAlteringChannel%s", channelType), func(t *testing.T) {
			webhookURL := "https://my-webhook.example.com"
			alertingChannel := AlertingChannel{
				ID:         idFieldValue,
				Name:       nameFieldValue,
				Kind:       channelType,
				WebhookURL: &webhookURL,
			}

			err := alertingChannel.Validate()

			assert.Nil(t, err)
		})
	}
}

func TestShouldFailToValidateWebhhokBasedAlteringChannelWhenWebhookUrlIsMissing(t *testing.T) {
	for _, channelType := range []AlertingChannelType{GoogleChatChannelType, Office365ChannelType, SlackChannelType} {
		t.Run(fmt.Sprintf("TestShouldFailToValidateWebhhokBasedAlteringChannel%sWhenWebhookUrlIsMissing", channelType), func(t *testing.T) {
			alertingChannel := AlertingChannel{
				ID:   idFieldValue,
				Name: nameFieldValue,
				Kind: channelType,
			}

			err := alertingChannel.Validate()

			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "Webhook URL")
		})
	}
}

func TestShouldFailToValidateWebhhokBasedAlteringChannelWhenWebhookUrlIsBlank(t *testing.T) {
	for _, channelType := range []AlertingChannelType{GoogleChatChannelType, Office365ChannelType, SlackChannelType} {
		t.Run(fmt.Sprintf("TestShouldFailToValidateWebhhokBasedAlteringChannel%sWhenWebhookUrlIsBlank", channelType), func(t *testing.T) {
			webhookURL := " "
			alertingChannel := AlertingChannel{
				ID:         idFieldValue,
				Name:       nameFieldValue,
				Kind:       channelType,
				WebhookURL: &webhookURL,
			}

			err := alertingChannel.Validate()

			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "Webhook URL")
		})
	}
}

func TestShouldSuccussullyValidateConsistentOpsGenieAlteringChannel(t *testing.T) {
	apiKey := apiKeyFieldValue
	region := EuOpsGenieRegion
	tags := tagsFieldValue

	alertingChannel := AlertingChannel{
		ID:     idFieldValue,
		Name:   nameFieldValue,
		Kind:   OpsGenieChannelType,
		APIKey: &apiKey,
		Region: &region,
		Tags:   &tags,
	}

	err := alertingChannel.Validate()

	assert.Nil(t, err)
}

func TestShouldFailToValidateOpsGenieAlteringChannelWhenApiKeyIsMissing(t *testing.T) {
	region := EuOpsGenieRegion
	tags := tagsFieldValue

	alertingChannel := AlertingChannel{
		ID:     idFieldValue,
		Name:   nameFieldValue,
		Kind:   OpsGenieChannelType,
		Region: &region,
		Tags:   &tags,
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "API key")
}

func TestShouldFailToValidateOpsGenieAlteringChannelWhenApiKeyIsBlank(t *testing.T) {
	region := EuOpsGenieRegion
	tags := tagsFieldValue
	apiKey := " "

	alertingChannel := AlertingChannel{
		ID:     idFieldValue,
		Name:   nameFieldValue,
		Kind:   OpsGenieChannelType,
		APIKey: &apiKey,
		Region: &region,
		Tags:   &tags,
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "API key")
}

func TestShouldFailToValidateOpsGenieAlteringChannelWhenRegionIsMissing(t *testing.T) {
	apiKey := apiKeyFieldValue
	tags := tagsFieldValue

	alertingChannel := AlertingChannel{
		ID:     idFieldValue,
		Name:   nameFieldValue,
		Kind:   OpsGenieChannelType,
		APIKey: &apiKey,
		Tags:   &tags,
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "region")
}

func TestShouldFailToValidateOpsGenieAlteringChannelWhenRegionIsNotValid(t *testing.T) {
	apiKey := apiKeyFieldValue
	region := OpsGenieRegionType("Invalid")
	tags := tagsFieldValue

	alertingChannel := AlertingChannel{
		ID:     idFieldValue,
		Name:   nameFieldValue,
		Kind:   OpsGenieChannelType,
		APIKey: &apiKey,
		Region: &region,
		Tags:   &tags,
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "region")
}

func TestShouldFailToValidateOpsGenieAlteringChannelWhenTagsAreMissing(t *testing.T) {
	apiKey := apiKeyFieldValue
	region := EuOpsGenieRegion

	alertingChannel := AlertingChannel{
		ID:     idFieldValue,
		Name:   nameFieldValue,
		Kind:   OpsGenieChannelType,
		APIKey: &apiKey,
		Region: &region,
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "tags")
}

func TestShouldFailToValidateOpsGenieAlteringChannelWhenTagsAreBlank(t *testing.T) {
	apiKey := apiKeyFieldValue
	region := EuOpsGenieRegion
	tags := " "

	alertingChannel := AlertingChannel{
		ID:     idFieldValue,
		Name:   nameFieldValue,
		Kind:   OpsGenieChannelType,
		APIKey: &apiKey,
		Region: &region,
		Tags:   &tags,
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "tags")
}

func TestShouldSuccussullyValidateConsistentPagerDutyAlteringChannel(t *testing.T) {
	integrationId := serviceIntegrationKeyFieldValue

	alertingChannel := AlertingChannel{
		ID:                    idFieldValue,
		Name:                  nameFieldValue,
		Kind:                  PagerDutyChannelType,
		ServiceIntegrationKey: &integrationId,
	}

	err := alertingChannel.Validate()

	assert.Nil(t, err)
}

func TestShouldFailToValidatePagerDutyAlteringChannelWhenServiceIntegrationKeyIsMissing(t *testing.T) {
	alertingChannel := AlertingChannel{
		ID:   idFieldValue,
		Name: nameFieldValue,
		Kind: PagerDutyChannelType,
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Service integration key")
}

func TestShouldFailToValidatePagerdutyAlteringChannelWhenServiceIntegrationKeyIsBlank(t *testing.T) {
	integrationId := "  "

	alertingChannel := AlertingChannel{
		ID:                    idFieldValue,
		Name:                  nameFieldValue,
		Kind:                  PagerDutyChannelType,
		ServiceIntegrationKey: &integrationId,
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Service integration key")
}

func TestShouldSuccussullyValidateConsistentSplunkAlteringChannel(t *testing.T) {
	url := urlFieldValue
	token := tokenFieldValue

	alertingChannel := AlertingChannel{
		ID:    idFieldValue,
		Name:  nameFieldValue,
		Kind:  SplunkChannelType,
		URL:   &url,
		Token: &token,
	}

	err := alertingChannel.Validate()

	assert.Nil(t, err)
}

func TestShouldFailToValidateSplunkAlteringChannelWhenUrlIsMissing(t *testing.T) {
	token := tokenFieldValue

	alertingChannel := AlertingChannel{
		ID:    idFieldValue,
		Name:  nameFieldValue,
		Kind:  SplunkChannelType,
		Token: &token,
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "URL")
}

func TestShouldFailToValidateSplunkAlteringChannelWhenUrlIsBlank(t *testing.T) {
	url := " "
	token := tokenFieldValue

	alertingChannel := AlertingChannel{
		ID:    idFieldValue,
		Name:  nameFieldValue,
		Kind:  SplunkChannelType,
		URL:   &url,
		Token: &token,
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "URL")
}

func TestShouldFailToValidateSplunkAlteringChannelWhenTokenIsMissing(t *testing.T) {
	url := urlFieldValue

	alertingChannel := AlertingChannel{
		ID:   idFieldValue,
		Name: nameFieldValue,
		Kind: SplunkChannelType,
		URL:  &url,
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Token")
}

func TestShouldFailToValidateSplunkAlteringChannelWhenTokenIsBlank(t *testing.T) {
	url := urlFieldValue
	token := " "

	alertingChannel := AlertingChannel{
		ID:    idFieldValue,
		Name:  nameFieldValue,
		Kind:  SplunkChannelType,
		URL:   &url,
		Token: &token,
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Token")
}

func TestShouldSuccussullyValidateConsistentVictorOpsAlteringChannel(t *testing.T) {
	apiKey := apiKeyFieldValue
	routingKey := routingKeyFieldValue

	alertingChannel := AlertingChannel{
		ID:         idFieldValue,
		Name:       nameFieldValue,
		Kind:       VictorOpsChannelType,
		APIKey:     &apiKey,
		RoutingKey: &routingKey,
	}

	err := alertingChannel.Validate()

	assert.Nil(t, err)
}

func TestShouldFailToValidateVictorOpsAlteringChannelWhenApiKeyIsMissing(t *testing.T) {
	routingKey := routingKeyFieldValue

	alertingChannel := AlertingChannel{
		ID:         idFieldValue,
		Name:       nameFieldValue,
		Kind:       VictorOpsChannelType,
		RoutingKey: &routingKey,
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "API Key")
}

func TestShouldFailToValidateVictorOpsAlteringChannelWhenApiKeyIsBlank(t *testing.T) {
	apiKey := " "
	routingKey := routingKeyFieldValue

	alertingChannel := AlertingChannel{
		ID:         idFieldValue,
		Name:       nameFieldValue,
		Kind:       VictorOpsChannelType,
		APIKey:     &apiKey,
		RoutingKey: &routingKey,
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "API Key")
}

func TestShouldFailToValidateVictorOpsAlteringChannelWhenRoutingKeyIsMissing(t *testing.T) {
	apiKey := apiKeyFieldValue

	alertingChannel := AlertingChannel{
		ID:     idFieldValue,
		Name:   nameFieldValue,
		Kind:   VictorOpsChannelType,
		APIKey: &apiKey,
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Routing Key")
}

func TestShouldFailToValidateVictorOpsAlteringChannelWhenRoutingKeyIsBlank(t *testing.T) {
	apiKey := apiKeyFieldValue
	routingKey := " "

	alertingChannel := AlertingChannel{
		ID:         idFieldValue,
		Name:       nameFieldValue,
		Kind:       VictorOpsChannelType,
		APIKey:     &apiKey,
		RoutingKey: &routingKey,
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Routing Key")
}

func TestShouldSuccussullyValidateConsistentMinimalGenericWebhookAlteringChannel(t *testing.T) {
	alertingChannel := AlertingChannel{
		ID:          idFieldValue,
		Name:        nameFieldValue,
		Kind:        WebhookChannelType,
		WebhookURLs: []string{"url"},
	}

	err := alertingChannel.Validate()

	assert.Nil(t, err)
}

func TestShouldSuccussullyValidateConsistentFullGenericWebhookAlteringChannel(t *testing.T) {
	alertingChannel := AlertingChannel{
		ID:          idFieldValue,
		Name:        nameFieldValue,
		Kind:        WebhookChannelType,
		WebhookURLs: []string{"url1", "url2"},
		Headers:     []string{"key1: value1", "key2: value2"},
	}

	err := alertingChannel.Validate()

	assert.Nil(t, err)
}

func TestShouldFailToValidateGenericWebhookAlteringChannelWhenNoWebhookUrlIsProvided(t *testing.T) {
	alertingChannel := AlertingChannel{
		ID:   idFieldValue,
		Name: nameFieldValue,
		Kind: WebhookChannelType,
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Webhook URLs")
}
