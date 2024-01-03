package restapi

// ApplicationAlertRule is the representation of an application alert rule in Instana
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
