package restapi

//ApplicationAlertConfigsResourcePath the base path of the Instana REST API for application alert configs
const ApplicationAlertConfigsResourcePath = EventSettingsBasePath + "/application-alert-configs"

//ApplicationAlertRule is the representation of an application alert rule in Instana
type ApplicationAlertRule struct {
	AlertType   string      `json:"alertType"`
	MetricName  string      `json:"metricName"`
	Aggregation Aggregation `json:"aggregation"`

	StatusCodeStart *int32 `json:"statusCodeStart"`
	StatusCodeEnd   *int32 `json:"statusCodeEnd"`

	Level    *LogLevel           `json:"level"`
	Message  *string             `json:"message"`
	Operator *ExpressionOperator `json:"operator"`
}

//IncludedEndpoint custom type to include of a specific endpoint in an alert config
type IncludedEndpoint struct {
	EndpointID string `json:"endpointId"`
	Inclusive  bool   `json:"inclusive"`
}

//IncludedService custom type to include of a specific service in an alert config
type IncludedService struct {
	ServiceID string `json:"serviceId"`
	Inclusive bool   `json:"inclusive"`

	Endpoints map[string]IncludedEndpoint `json:"endpoints"`
}

//IncludedApplication custom type to include specific applications in an alert config
type IncludedApplication struct {
	ApplicationID string `json:"applicationId"`
	Inclusive     bool   `json:"inclusive"`

	Services map[string]IncludedService `json:"services"`
}

//StaticStringField custom type to represent static fields with a string value for custom payloads
type StaticStringField struct {
	Type  string `json:"type"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

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
	CustomerPayloadFields []StaticStringField            `json:"customerPayloadFields"`
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
	//TODO add validation when check cannot be covered in TF schema
	return nil
}
