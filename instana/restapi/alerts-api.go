package restapi

// AlertsResourcePath path to Alerts resource of Instana RESTful API
const AlertsResourcePath = EventSettingsBasePath + "/alerts"

// EventFilteringConfiguration type definiton of an EventFilteringConfiguration of a AlertingConfiguration of the Instana ReST AOI
type EventFilteringConfiguration struct {
	Query      *string          `json:"query"`
	RuleIDs    []string         `json:"ruleIds"`
	EventTypes []AlertEventType `json:"eventTypes"`
}

// AlertingConfiguration type definition of an Alertinng Configruation in Instana REST API
type AlertingConfiguration struct {
	ID                          string                      `json:"id"`
	AlertName                   string                      `json:"alertName"`
	IntegrationIDs              []string                    `json:"integrationIds"`
	EventFilteringConfiguration EventFilteringConfiguration `json:"eventFilteringConfiguration"`
}

// GetIDForResourcePath implementation of the interface InstanaDataObject
func (c *AlertingConfiguration) GetIDForResourcePath() string {
	return c.ID
}
