package restapi

// AlertsResourcePath path to Alerts resource of Instana RESTful API
const AlertsResourcePath = EventSettingsBasePath + "/alerts"

// EventFilteringConfiguration type definiton of an EventFilteringConfiguration of a AlertingConfiguration of the Instana ReST AOI
type EventFilteringConfiguration struct {
	Query      *string          `json:"query"`
	RuleIDs    []string         `json:"ruleIds"`
	EventTypes []AlertEventType `json:"eventTypes"`
}

// AlertingConfiguration type definition of an Alerting Configuration in Instana REST API
type AlertingConfiguration struct {
	ID                          string                      `json:"id"`
	AlertName                   string                      `json:"alertName"`
	IntegrationIDs              []string                    `json:"integrationIds"`
	EventFilteringConfiguration EventFilteringConfiguration `json:"eventFilteringConfiguration"`
	CustomerPayloadFields       []CustomPayloadField[any]   `json:"customPayloadFields"`
}

// GetIDForResourcePath implementation of the interface InstanaDataObject
func (c *AlertingConfiguration) GetIDForResourcePath() string {
	return c.ID
}

// GetCustomerPayloadFields implementation of the interface customPayloadFieldsAwareInstanaDataObject
func (a *AlertingConfiguration) GetCustomerPayloadFields() []CustomPayloadField[any] {
	return a.CustomerPayloadFields
}

// SetCustomerPayloadFields implementation of the interface customPayloadFieldsAwareInstanaDataObject
func (a *AlertingConfiguration) SetCustomerPayloadFields(fields []CustomPayloadField[any]) {
	a.CustomerPayloadFields = fields
}
