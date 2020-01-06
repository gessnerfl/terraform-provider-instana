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
)

func TestShouldSuccussullyValididateConsistentEmailAlteringChannel(t *testing.T) {
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

func TestShouldFailToValididateAlteringChannelWhenIdIsMissing(t *testing.T) {
	alertingChannel := AlertingChannel{
		Name:   nameFieldValue,
		Kind:   EmailChannelType,
		Emails: []string{email1FieldValue, email2FieldValue},
	}

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "ID") {
		t.Fatal("Expected validate to fail as ID is not provided")
	}
}

func TestShouldFailToValididateAlteringChannelWhenIdIsBlank(t *testing.T) {
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

func TestShouldFailToValididateAlteringChannelWhenNameIsMissing(t *testing.T) {
	alertingChannel := AlertingChannel{
		ID:     idFieldValue,
		Kind:   EmailChannelType,
		Emails: []string{email1FieldValue, email2FieldValue},
	}

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "Name") {
		t.Fatal("Expected validate to fail as name is not provided")
	}
}

func TestShouldFailToValididateAlteringChannelWhenNameIsBlank(t *testing.T) {
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

func TestShouldFailToValididateAlteringChannelWhenKindIsMissing(t *testing.T) {
	alertingChannel := AlertingChannel{
		ID:     idFieldValue,
		Name:   nameFieldValue,
		Emails: []string{email1FieldValue, email2FieldValue},
	}

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "Kind") {
		t.Fatal("Expected validate to fail as kind is not provided")
	}
}

func TestShouldFailToValididateAlteringChannelWhenKindIsNotValid(t *testing.T) {
	alertingChannel := AlertingChannel{
		ID:     idFieldValue,
		Name:   nameFieldValue,
		Kind:   AlertingChannelType("invalid"),
		Emails: []string{email1FieldValue, email2FieldValue},
	}

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "Kind") {
		t.Fatal("Expected validate to fail as kind is not valid")
	}
}
func TestShouldFailToValididateEmailAlteringChannelWhenNoEmailIsProvided(t *testing.T) {
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

func TestShouldSuccussullyValididateConsistentWebhhokBasedAlteringChannel(t *testing.T) {
	for _, channelType := range []AlertingChannelType{GoogleChatChannelType, Office365ChannelType, SlackChannelType} {
		t.Run(fmt.Sprintf("TestShouldSuccussullyValididateConsistentWebhhokBasedAlteringChannel%s", channelType), func(t *testing.T) {
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

func TestShouldFailToValididateWebhhokBasedAlteringChannelWhenWebhookUrlIsMissing(t *testing.T) {
	for _, channelType := range []AlertingChannelType{GoogleChatChannelType, Office365ChannelType, SlackChannelType} {
		t.Run(fmt.Sprintf("TestShouldFailToValididateWebhhokBasedAlteringChannel%sWhenWebhookUrlIsMissing", channelType), func(t *testing.T) {
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

func TestShouldFailToValididateWebhhokBasedAlteringChannelWhenWebhookUrlIsBlank(t *testing.T) {
	for _, channelType := range []AlertingChannelType{GoogleChatChannelType, Office365ChannelType, SlackChannelType} {
		t.Run(fmt.Sprintf("TestShouldFailToValididateWebhhokBasedAlteringChannel%sWhenWebhookUrlIsBlank", channelType), func(t *testing.T) {
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

func TestShouldSuccussullyValididateConsistentOpsGenieAlteringChannel(t *testing.T) {
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

func TestShouldFailToValididateOpsGenieAlteringChannelWhenApiKeyIsMissing(t *testing.T) {
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

func TestShouldFailToValididateOpsGenieAlteringChannelWhenApiKeyIsBlank(t *testing.T) {
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

func TestShouldFailToValididateOpsGenieAlteringChannelWhenRegionIsMissing(t *testing.T) {
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

func TestShouldFailToValididateOpsGenieAlteringChannelWhenRegionIsNotValid(t *testing.T) {
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

func TestShouldFailToValididateOpsGenieAlteringChannelWhenTagsAreMissing(t *testing.T) {
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

func TestShouldFailToValididateOpsGenieAlteringChannelWhenTagsAreBlank(t *testing.T) {
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

func TestShouldSuccussullyValididateConsistentPagerDutyAlteringChannel(t *testing.T) {
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

func TestShouldFailToValididatePagerDutyAlteringChannelWhenServiceIntegrationKeyIsMissing(t *testing.T) {
	alertingChannel := AlertingChannel{
		ID:   idFieldValue,
		Name: nameFieldValue,
		Kind: PagerDutyChannelType,
	}

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "Service integration key") {
		t.Fatal("Expected validate to fail as API key is missing")
	}
}

func TestShouldFailToValididatePagerdutyAlteringChannelWhenServiceIntegrationKeyIsBlank(t *testing.T) {
	integrationId := "  "

	alertingChannel := AlertingChannel{
		ID:                    idFieldValue,
		Name:                  nameFieldValue,
		Kind:                  PagerDutyChannelType,
		ServiceIntegrationKey: &integrationId,
	}

	if err := alertingChannel.Validate(); err == nil || !strings.Contains(err.Error(), "Service integration key") {
		t.Fatal("Expected validate to fail as API key is missing")
	}
}
