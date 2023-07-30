package restapi_test

import (
	"fmt"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnIDOfAlteringChannelDS(t *testing.T) {
	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:     idFieldValue,
			Name:   nameFieldValue,
			Kind:   EmailChannelType,
			Emails: []string{email1FieldValue, email2FieldValue},
		},
	}

	assert.Equal(t, idFieldValue, alertingChannel.GetIDForResourcePath())
}

func TestShouldSuccussullyValidateConsistentEmailAlteringChannelDS(t *testing.T) {
	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:     idFieldValue,
			Name:   nameFieldValue,
			Kind:   EmailChannelType,
			Emails: []string{email1FieldValue, email2FieldValue},
		},
	}

	err := alertingChannel.Validate()
	assert.Nil(t, err)
}

func TestShouldFailToValidateAlteringChannelDSWhenIdIsMissing(t *testing.T) {
	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			Name:   nameFieldValue,
			Kind:   EmailChannelType,
			Emails: []string{email1FieldValue, email2FieldValue},
		},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "ID")
}

func TestShouldFailToValidateAlteringChannelDSWhenIdIsBlank(t *testing.T) {
	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:     " ",
			Name:   nameFieldValue,
			Kind:   EmailChannelType,
			Emails: []string{email1FieldValue, email2FieldValue},
		},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "ID")
}

func TestShouldFailToValidateAlteringChannelDSWhenNameIsMissing(t *testing.T) {
	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:     idFieldValue,
			Kind:   EmailChannelType,
			Emails: []string{email1FieldValue, email2FieldValue},
		},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Name")
}

func TestShouldFailToValidateAlteringChannelDSWhenNameIsBlank(t *testing.T) {
	name := " "
	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:     idFieldValue,
			Kind:   EmailChannelType,
			Name:   name,
			Emails: []string{email1FieldValue, email2FieldValue},
		},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Name")
}

func TestShouldFailToValidateAlteringChannelDSWhenKindIsMissing(t *testing.T) {
	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:     idFieldValue,
			Name:   nameFieldValue,
			Emails: []string{email1FieldValue, email2FieldValue},
		},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Kind")
}

func TestShouldFailToValidateAlteringChannelDSWhenKindIsNotValid(t *testing.T) {
	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:     idFieldValue,
			Name:   nameFieldValue,
			Kind:   AlertingChannelType("invalid"),
			Emails: []string{email1FieldValue, email2FieldValue},
		},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unsupported alerting channel type")
}
func TestShouldFailToValidateEmailAlteringChannelDSWhenNoEmailIsProvided(t *testing.T) {
	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:     idFieldValue,
			Name:   nameFieldValue,
			Kind:   EmailChannelType,
			Emails: []string{},
		},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Email addresses")
}

func TestShouldSuccussullyValidateConsistentWebhhokBasedAlteringChannelDS(t *testing.T) {
	for _, channelType := range []AlertingChannelType{GoogleChatChannelType, Office365ChannelType, SlackChannelType} {
		t.Run(fmt.Sprintf("TestShouldSuccussullyValidateConsistentWebhhokBasedAlteringChannelDS%s", channelType), func(t *testing.T) {
			webhookURL := "https://my-webhook.example.com"
			alertingChannel := AlertingChannelDS{
				AlertingChannel{
					ID:         idFieldValue,
					Name:       nameFieldValue,
					Kind:       channelType,
					WebhookURL: &webhookURL,
				},
			}

			err := alertingChannel.Validate()

			assert.Nil(t, err)
		})
	}
}

func TestShouldFailToValidateWebhhokBasedAlteringChannelDSWhenWebhookUrlIsMissing(t *testing.T) {
	for _, channelType := range []AlertingChannelType{GoogleChatChannelType, Office365ChannelType, SlackChannelType} {
		t.Run(fmt.Sprintf("TestShouldFailToValidateWebhhokBasedAlteringChannelDS%sWhenWebhookUrlIsMissing", channelType), func(t *testing.T) {
			alertingChannel := AlertingChannelDS{
				AlertingChannel{
					ID:   idFieldValue,
					Name: nameFieldValue,
					Kind: channelType,
				},
			}

			err := alertingChannel.Validate()

			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "Webhook URL")
		})
	}
}

func TestShouldFailToValidateWebhhokBasedAlteringChannelDSWhenWebhookUrlIsBlank(t *testing.T) {
	for _, channelType := range []AlertingChannelType{GoogleChatChannelType, Office365ChannelType, SlackChannelType} {
		t.Run(fmt.Sprintf("TestShouldFailToValidateWebhhokBasedAlteringChannelDS%sWhenWebhookUrlIsBlank", channelType), func(t *testing.T) {
			webhookURL := " "
			alertingChannel := AlertingChannelDS{
				AlertingChannel{
					ID:         idFieldValue,
					Name:       nameFieldValue,
					Kind:       channelType,
					WebhookURL: &webhookURL,
				},
			}

			err := alertingChannel.Validate()

			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "Webhook URL")
		})
	}
}

func TestShouldSuccussullyValidateConsistentOpsGenieAlteringChannelDS(t *testing.T) {
	apiKey := apiKeyFieldValue
	region := EuOpsGenieRegion
	tags := tagsFieldValue

	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:     idFieldValue,
			Name:   nameFieldValue,
			Kind:   OpsGenieChannelType,
			APIKey: &apiKey,
			Region: &region,
			Tags:   &tags,
		},
	}

	err := alertingChannel.Validate()

	assert.Nil(t, err)
}

func TestShouldFailToValidateOpsGenieAlteringChannelDSWhenApiKeyIsMissing(t *testing.T) {
	region := EuOpsGenieRegion
	tags := tagsFieldValue

	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:     idFieldValue,
			Name:   nameFieldValue,
			Kind:   OpsGenieChannelType,
			Region: &region,
			Tags:   &tags,
		},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "API key")
}

func TestShouldFailToValidateOpsGenieAlteringChannelDSWhenApiKeyIsBlank(t *testing.T) {
	region := EuOpsGenieRegion
	tags := tagsFieldValue
	apiKey := " "

	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:     idFieldValue,
			Name:   nameFieldValue,
			Kind:   OpsGenieChannelType,
			APIKey: &apiKey,
			Region: &region,
			Tags:   &tags,
		},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "API key")
}

func TestShouldFailToValidateOpsGenieAlteringChannelDSWhenRegionIsMissing(t *testing.T) {
	apiKey := apiKeyFieldValue
	tags := tagsFieldValue

	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:     idFieldValue,
			Name:   nameFieldValue,
			Kind:   OpsGenieChannelType,
			APIKey: &apiKey,
			Tags:   &tags,
		},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "region")
}

func TestShouldFailToValidateOpsGenieAlteringChannelDSWhenRegionIsNotValid(t *testing.T) {
	apiKey := apiKeyFieldValue
	region := OpsGenieRegionType("Invalid")
	tags := tagsFieldValue

	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:     idFieldValue,
			Name:   nameFieldValue,
			Kind:   OpsGenieChannelType,
			APIKey: &apiKey,
			Region: &region,
			Tags:   &tags,
		},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "region")
}

func TestShouldFailToValidateOpsGenieAlteringChannelDSWhenTagsAreMissing(t *testing.T) {
	apiKey := apiKeyFieldValue
	region := EuOpsGenieRegion

	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:     idFieldValue,
			Name:   nameFieldValue,
			Kind:   OpsGenieChannelType,
			APIKey: &apiKey,
			Region: &region,
		},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "tags")
}

func TestShouldFailToValidateOpsGenieAlteringChannelDSWhenTagsAreBlank(t *testing.T) {
	apiKey := apiKeyFieldValue
	region := EuOpsGenieRegion
	tags := " "

	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:     idFieldValue,
			Name:   nameFieldValue,
			Kind:   OpsGenieChannelType,
			APIKey: &apiKey,
			Region: &region,
			Tags:   &tags,
		},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "tags")
}

func TestShouldSuccussullyValidateConsistentPagerDutyAlteringChannelDS(t *testing.T) {
	integrationId := serviceIntegrationKeyFieldValue

	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:                    idFieldValue,
			Name:                  nameFieldValue,
			Kind:                  PagerDutyChannelType,
			ServiceIntegrationKey: &integrationId,
		},
	}

	err := alertingChannel.Validate()

	assert.Nil(t, err)
}

func TestShouldFailToValidatePagerDutyAlteringChannelDSWhenServiceIntegrationKeyIsMissing(t *testing.T) {
	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:   idFieldValue,
			Name: nameFieldValue,
			Kind: PagerDutyChannelType,
		},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Service integration key")
}

func TestShouldFailToValidatePagerdutyAlteringChannelDSWhenServiceIntegrationKeyIsBlank(t *testing.T) {
	integrationId := "  "

	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:                    idFieldValue,
			Name:                  nameFieldValue,
			Kind:                  PagerDutyChannelType,
			ServiceIntegrationKey: &integrationId,
		},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Service integration key")
}

func TestShouldSuccussullyValidateConsistentSplunkAlteringChannelDS(t *testing.T) {
	url := urlFieldValue
	token := tokenFieldValue

	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:    idFieldValue,
			Name:  nameFieldValue,
			Kind:  SplunkChannelType,
			URL:   &url,
			Token: &token,
		},
	}

	err := alertingChannel.Validate()

	assert.Nil(t, err)
}

func TestShouldFailToValidateSplunkAlteringChannelDSWhenUrlIsMissing(t *testing.T) {
	token := tokenFieldValue

	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:    idFieldValue,
			Name:  nameFieldValue,
			Kind:  SplunkChannelType,
			Token: &token,
		},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "URL")
}

func TestShouldFailToValidateSplunkAlteringChannelDSWhenUrlIsBlank(t *testing.T) {
	url := " "
	token := tokenFieldValue

	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:    idFieldValue,
			Name:  nameFieldValue,
			Kind:  SplunkChannelType,
			URL:   &url,
			Token: &token,
		},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "URL")
}

func TestShouldFailToValidateSplunkAlteringChannelDSWhenTokenIsMissing(t *testing.T) {
	url := urlFieldValue

	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:   idFieldValue,
			Name: nameFieldValue,
			Kind: SplunkChannelType,
			URL:  &url,
		},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Token")
}

func TestShouldFailToValidateSplunkAlteringChannelDSWhenTokenIsBlank(t *testing.T) {
	url := urlFieldValue
	token := " "

	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:    idFieldValue,
			Name:  nameFieldValue,
			Kind:  SplunkChannelType,
			URL:   &url,
			Token: &token,
		},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Token")
}

func TestShouldSuccussullyValidateConsistentVictorOpsAlteringChannelDS(t *testing.T) {
	apiKey := apiKeyFieldValue
	routingKey := routingKeyFieldValue

	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:         idFieldValue,
			Name:       nameFieldValue,
			Kind:       VictorOpsChannelType,
			APIKey:     &apiKey,
			RoutingKey: &routingKey,
		},
	}

	err := alertingChannel.Validate()

	assert.Nil(t, err)
}

func TestShouldFailToValidateVictorOpsAlteringChannelDSWhenApiKeyIsMissing(t *testing.T) {
	routingKey := routingKeyFieldValue

	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:         idFieldValue,
			Name:       nameFieldValue,
			Kind:       VictorOpsChannelType,
			RoutingKey: &routingKey,
		},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "API Key")
}

func TestShouldFailToValidateVictorOpsAlteringChannelDSWhenApiKeyIsBlank(t *testing.T) {
	apiKey := " "
	routingKey := routingKeyFieldValue

	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:         idFieldValue,
			Name:       nameFieldValue,
			Kind:       VictorOpsChannelType,
			APIKey:     &apiKey,
			RoutingKey: &routingKey,
		},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "API Key")
}

func TestShouldFailToValidateVictorOpsAlteringChannelDSWhenRoutingKeyIsMissing(t *testing.T) {
	apiKey := apiKeyFieldValue

	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:     idFieldValue,
			Name:   nameFieldValue,
			Kind:   VictorOpsChannelType,
			APIKey: &apiKey,
		},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Routing Key")
}

func TestShouldFailToValidateVictorOpsAlteringChannelDSWhenRoutingKeyIsBlank(t *testing.T) {
	apiKey := apiKeyFieldValue
	routingKey := " "

	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:         idFieldValue,
			Name:       nameFieldValue,
			Kind:       VictorOpsChannelType,
			APIKey:     &apiKey,
			RoutingKey: &routingKey,
		},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Routing Key")
}

func TestShouldSuccussullyValidateConsistentMinimalGenericWebhookAlteringChannelDS(t *testing.T) {
	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:          idFieldValue,
			Name:        nameFieldValue,
			Kind:        WebhookChannelType,
			WebhookURLs: []string{"url"},
		},
	}

	err := alertingChannel.Validate()

	assert.Nil(t, err)
}

func TestShouldSuccussullyValidateConsistentFullGenericWebhookAlteringChannelDS(t *testing.T) {
	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:          idFieldValue,
			Name:        nameFieldValue,
			Kind:        WebhookChannelType,
			WebhookURLs: []string{"url1", "url2"},
			Headers:     []string{"key1: value1", "key2: value2"},
		},
	}

	err := alertingChannel.Validate()

	assert.Nil(t, err)
}

func TestShouldFailToValidateGenericWebhookAlteringChannelDSWhenNoWebhookUrlIsProvided(t *testing.T) {
	alertingChannel := AlertingChannelDS{
		AlertingChannel{
			ID:   idFieldValue,
			Name: nameFieldValue,
			Kind: WebhookChannelType,
		},
	}

	err := alertingChannel.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Webhook URLs")
}
