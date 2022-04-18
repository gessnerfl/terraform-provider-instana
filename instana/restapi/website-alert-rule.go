package restapi

type WebsiteAlertRule struct {
	AlertType   string              `json:"alertType"`
	MetricName  string              `json:"metricName"`
	Aggregation Aggregation         `json:"aggregation"`
	Operator    *ExpressionOperator `json:"operator"`
	Value       *string             `json:"value"`
}
