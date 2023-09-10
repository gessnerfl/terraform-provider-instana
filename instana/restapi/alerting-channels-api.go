package restapi

// AlertingChannelsResourcePath path to Alerting channels resource of Instana RESTful API
const AlertingChannelsResourcePath = EventSettingsBasePath + "/alertingChannels"

// AlertingChannel is the representation of an alerting channel in Instana
type AlertingChannel struct {
	ID                    string              `json:"id"`
	Name                  string              `json:"name"`
	Kind                  AlertingChannelType `json:"kind"`
	Emails                []string            `json:"emails"`
	WebhookURL            *string             `json:"webhookUrl"`
	APIKey                *string             `json:"apiKey"`
	Tags                  *string             `json:"tags"`
	Region                *string             `json:"region"`
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
