package restapi_test

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
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

	if idFieldValue != alertingChannel.GetID() {
		t.Fatal("GetID should return id value of alerting channel")
	}
}

func TestShouldSuccussullyValidateConsistentEmailAlteringChannel(t *testing.T) {
	alertingChannel := AlertingChannel{
		ID:     idFieldValue,
		Name:   nameFieldValue,
		Kind:   EmailChannelType,
		Emails: []string{email1FieldValue, email2FieldValue},
	}

	if err := alertingChannel.Validate(); err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
}

func TestShouldFailToValidateAlteringChannelWhenIdIsMissing(t *testing.T) {
	alertingChannel := AlertingChannel{
		Name:   nameFieldValue,
		Kind:   EmailChannelType,
		Emails: []string{email1FieldValue, email2FieldValue},
	}

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "ID") {
		t.Fatal("Expected validate to fail as ID is not provided")
	}
}

func TestShouldFailToValidateAlteringChannelWhenIdIsBlank(t *testing.T) {
	alertingChannel := AlertingChannel{
		ID:     " ",
		Name:   nameFieldValue,
		Kind:   EmailChannelType,
		Emails: []string{email1FieldValue, email2FieldValue},
	}

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "ID") {
		t.Fatal("Expected validate to fail as ID is not provided")
	}
}

func TestShouldFailToValidateAlteringChannelWhenNameIsMissing(t *testing.T) {
	alertingChannel := AlertingChannel{
		ID:     idFieldValue,
		Kind:   EmailChannelType,
		Emails: []string{email1FieldValue, email2FieldValue},
	}

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "Name") {
		t.Fatal("Expected validate to fail as name is not provided")
	}
}

func TestShouldFailToValidateAlteringChannelWhenNameIsBlank(t *testing.T) {
	name := " "
	alertingChannel := AlertingChannel{
		ID:     idFieldValue,
		Kind:   EmailChannelType,
		Name:   name,
		Emails: []string{email1FieldValue, email2FieldValue},
	}

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "Name") {
		t.Fatal("Expected validate to fail as name is not provided")
	}
}

func TestShouldFailToValidateAlteringChannelWhenKindIsMissing(t *testing.T) {
	alertingChannel := AlertingChannel{
		ID:     idFieldValue,
		Name:   nameFieldValue,
		Emails: []string{email1FieldValue, email2FieldValue},
	}

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "Kind") {
		t.Fatal("Expected validate to fail as kind is not provided")
	}
}

func TestShouldFailToValidateAlteringChannelWhenKindIsNotValid(t *testing.T) {
	alertingChannel := AlertingChannel{
		ID:     idFieldValue,
		Name:   nameFieldValue,
		Kind:   AlertingChannelType("invalid"),
		Emails: []string{email1FieldValue, email2FieldValue},
	}

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "unsupported alerting channel type") {
		t.Fatal("Expected validate to fail as kind is not valid")
	}
}
func TestShouldFailToValidateEmailAlteringChannelWhenNoEmailIsProvided(t *testing.T) {
	alertingChannel := AlertingChannel{
		ID:     idFieldValue,
		Name:   nameFieldValue,
		Kind:   EmailChannelType,
		Emails: []string{},
	}

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "Email addresses") {
		t.Fatal("Expected validate to fail as at least one email is missing")
	}
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

			if err := alertingChannel.Validate(); err != nil {
				t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
			}
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

			if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "Webhook URL") {
				t.Fatal("Expected validate to fail as webhook URL is missing")
			}
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

			if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "Webhook URL") {
				t.Fatal("Expected validate to fail as webhook URL is missing")
			}
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

	if err := alertingChannel.Validate(); err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
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

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "API key") {
		t.Fatal("Expected validate to fail as API key is missing")
	}
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

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "API key") {
		t.Fatal("Expected validate to fail as API key is missing")
	}
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

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "Region") {
		t.Fatal("Expected validate to fail as region is missing")
	}
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

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "Region") {
		t.Fatal("Expected validate to fail as region is missing")
	}
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

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "Tags") {
		t.Fatal("Expected validate to fail as tags are missing")
	}
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

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "Tags") {
		t.Fatal("Expected validate to fail as tags are missing")
	}
}

func TestShouldSuccussullyValidateConsistentPagerDutyAlteringChannel(t *testing.T) {
	integrationId := serviceIntegrationKeyFieldValue

	alertingChannel := AlertingChannel{
		ID:                    idFieldValue,
		Name:                  nameFieldValue,
		Kind:                  PagerDutyChannelType,
		ServiceIntegrationKey: &integrationId,
	}

	if err := alertingChannel.Validate(); err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
}

func TestShouldFailToValidatePagerDutyAlteringChannelWhenServiceIntegrationKeyIsMissing(t *testing.T) {
	alertingChannel := AlertingChannel{
		ID:   idFieldValue,
		Name: nameFieldValue,
		Kind: PagerDutyChannelType,
	}

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "Service integration key") {
		t.Fatal("Expected validate to fail as service integration key is missing")
	}
}

func TestShouldFailToValidatePagerdutyAlteringChannelWhenServiceIntegrationKeyIsBlank(t *testing.T) {
	integrationId := "  "

	alertingChannel := AlertingChannel{
		ID:                    idFieldValue,
		Name:                  nameFieldValue,
		Kind:                  PagerDutyChannelType,
		ServiceIntegrationKey: &integrationId,
	}

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "Service integration key") {
		t.Fatal("Expected validate to fail as service integration key is missing")
	}
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

	if err := alertingChannel.Validate(); err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
}

func TestShouldFailToValidateSplunkAlteringChannelWhenUrlIsMissing(t *testing.T) {
	token := tokenFieldValue

	alertingChannel := AlertingChannel{
		ID:    idFieldValue,
		Name:  nameFieldValue,
		Kind:  SplunkChannelType,
		Token: &token,
	}

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "URL") {
		t.Fatal("Expected validate to fail as URL is missing")
	}
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

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "URL") {
		t.Fatal("Expected validate to fail as URL is missing")
	}
}

func TestShouldFailToValidateSplunkAlteringChannelWhenTokenIsMissing(t *testing.T) {
	url := urlFieldValue

	alertingChannel := AlertingChannel{
		ID:   idFieldValue,
		Name: nameFieldValue,
		Kind: SplunkChannelType,
		URL:  &url,
	}

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "Token") {
		t.Fatal("Expected validate to fail as token is missing")
	}
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

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "Token") {
		t.Fatal("Expected validate to fail as token is missing")
	}
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

	if err := alertingChannel.Validate(); err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
}

func TestShouldFailToValidateVictorOpsAlteringChannelWhenApiKeyIsMissing(t *testing.T) {
	routingKey := routingKeyFieldValue

	alertingChannel := AlertingChannel{
		ID:         idFieldValue,
		Name:       nameFieldValue,
		Kind:       VictorOpsChannelType,
		RoutingKey: &routingKey,
	}

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "API Key") {
		t.Fatal("Expected validate to fail as API Key is missing")
	}
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

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "API Key") {
		t.Fatal("Expected validate to fail as API Key is missing")
	}
}

func TestShouldFailToValidateVictorOpsAlteringChannelWhenRoutingKeyIsMissing(t *testing.T) {
	apiKey := apiKeyFieldValue

	alertingChannel := AlertingChannel{
		ID:     idFieldValue,
		Name:   nameFieldValue,
		Kind:   VictorOpsChannelType,
		APIKey: &apiKey,
	}

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "Routing Key") {
		t.Fatal("Expected validate to fail as Routing Key is missing")
	}
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

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "Routing Key") {
		t.Fatal("Expected validate to fail as Routing Key is missing")
	}
}

func TestShouldSuccussullyValidateConsistentMinimalGenericWebhookAlteringChannel(t *testing.T) {
	alertingChannel := AlertingChannel{
		ID:          idFieldValue,
		Name:        nameFieldValue,
		Kind:        WebhookChannelType,
		WebhookURLs: []string{"url"},
	}

	if err := alertingChannel.Validate(); err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
}

func TestShouldSuccussullyValidateConsistentFullGenericWebhookAlteringChannel(t *testing.T) {
	alertingChannel := AlertingChannel{
		ID:          idFieldValue,
		Name:        nameFieldValue,
		Kind:        WebhookChannelType,
		WebhookURLs: []string{"url1", "url2"},
		Headers:     []string{"key1: value1", "key2: value2"},
	}

	if err := alertingChannel.Validate(); err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
}

func TestShouldFailToValidateGenericWebhookAlteringChannelWhenNoWebhookUrlIsProvided(t *testing.T) {
	alertingChannel := AlertingChannel{
		ID:   idFieldValue,
		Name: nameFieldValue,
		Kind: WebhookChannelType,
	}

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "Webhook URLs") {
		t.Fatal("Expected validate to fail as Webhook URLs are missing")
	}
}
