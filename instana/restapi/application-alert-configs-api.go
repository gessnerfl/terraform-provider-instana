package restapi

import (
	"errors"
	"github.com/gessnerfl/terraform-provider-instana/utils"
)

/// api/events/settings/application-alert-configs/
const ApplicationAlertConfigsResourcePath = EventSettingsBasePath + "/application-alert-configs"

type ApplicationAlertConfigsRule struct {
	AlertType   string `json:"alertType"`
	MetricName  string `json:"metricName"`
	Aggregation string `json:"aggregation"`

	Operator string `json:"operator"`
	Message  string `json:"message"`
	Level    string `json:"level"`
}

type Threshold struct {
	Type        string  `json:"type"`
	Operator    string  `json:"operator"`
	LastUpdated int     `json:"lastUpdate"`
	Value       float64 `json:"value"`
}

type ApplicationAlertConfigsTagFilter struct {
	Type         string  `json:"type"`
	Name         string  `json:"name"`
	StringValue  string  `json:"stringValue"`
	NumberValue  float64 `json:"numberValue"`
	BooleanValue bool    `json:"booleanValue"`
	Operator     string  `json:"operator"`
	Entity       string  `json:"entity"`
}

type ApplicationAlertConfigs struct {
	ID            string                      `json:"id"`
	AlertName     string                      `json:"name"`
	ApplicationId string                      `json:"applicationId"`
	Rule          ApplicationAlertConfigsRule `json:"rule"`
	// Complex TypeThreshold		            string						`json:"threshold"`
	Description string    `json:"description"`
	Severity    int       `json:"severity"`
	Threshold   Threshold `json:"threshold"`

	AlertChannelIds []string `json:"alertChannelIds"`
	//	IntegrationIDs              []string                    `json:"integrationIds"`
	//	EventFilteringConfiguration EventFilteringConfiguration `json:"eventFilteringConfiguration"`
	TagFilters []ApplicationAlertConfigsTagFilter `json:"tagFilters"`
}

//
//GetID implemention of the interface InstanaDataObject
func (c ApplicationAlertConfigs) GetID() string {
	return c.ID
}

//

func (c ApplicationAlertConfigs) Validate() error {
	if utils.IsBlank(c.ID) {
		return errors.New("ID is missing")
	}
	if utils.IsBlank(c.AlertName) {
		return errors.New("AlertName is missing")
	}
	if len(c.AlertName) > 256 {
		return errors.New("AlertName not valid; Maximum length of AlertName is 256 characters")
	}
	return nil
}
