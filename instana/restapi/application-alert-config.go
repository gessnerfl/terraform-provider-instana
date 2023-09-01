package restapi

// ApplicationAlertConfigsResourcePath the base path of the Instana REST API for application alert configs
const ApplicationAlertConfigsResourcePath = EventSettingsBasePath + "/application-alert-configs"

// GlobalApplicationAlertConfigsResourcePath the base path of the Instana REST API for global application alert configs
const GlobalApplicationAlertConfigsResourcePath = EventSettingsBasePath + "/global-alert-configs/applications"

// ApplicationAlertConfig is the representation of an application alert configuration in Instana
type ApplicationAlertConfig struct {
	ID                    string                         `json:"id"`
	Name                  string                         `json:"name"`
	Description           string                         `json:"description"`
	Severity              int                            `json:"severity"`
	Triggering            bool                           `json:"triggering"`
	Applications          map[string]IncludedApplication `json:"applications"`
	BoundaryScope         BoundaryScope                  `json:"boundaryScope"`
	TagFilterExpression   interface{}                    `json:"tagFilterExpression"`
	IncludeInternal       bool                           `json:"includeInternal"`
	IncludeSynthetic      bool                           `json:"includeSynthetic"`
	EvaluationType        ApplicationAlertEvaluationType `json:"evaluationType"`
	AlertChannelIDs       []string                       `json:"alertChannelIds"`
	Granularity           Granularity                    `json:"granularity"`
	CustomerPayloadFields []CustomPayloadField[any]      `json:"customPayloadFields"`
	Rule                  ApplicationAlertRule           `json:"rule"`
	Threshold             Threshold                      `json:"threshold"`
	TimeThreshold         TimeThreshold                  `json:"timeThreshold"`
}

// GetIDForResourcePath implementation of the interface InstanaDataObject
func (a *ApplicationAlertConfig) GetIDForResourcePath() string {
	return a.ID
}
