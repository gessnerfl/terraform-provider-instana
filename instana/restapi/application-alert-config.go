package restapi

//ApplicationAlertConfigsResourcePath the base path of the Instana REST API for application alert configs
const ApplicationAlertConfigsResourcePath = EventSettingsBasePath + "/application-alert-configs"

//ApplicationAlertConfig is the representation of an application alert configuration in Instana
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

//GetIDForResourcePath implementation of the interface InstanaDataObject
func (a *ApplicationAlertConfig) GetIDForResourcePath() string {
	return a.ID
}

//Validate implementation of the interface InstanaDataObject for ApplicationConfig
func (a *ApplicationAlertConfig) Validate() error {
	//No validation required validation part of terraform schema
	return nil
}
