package restapi

// WebsiteAlertRule struct representing the API model of a website alert rule
type WebsiteAlertRule struct {
	AlertType   string              `json:"alertType"`
	MetricName  string              `json:"metricName"`
	Aggregation *Aggregation        `json:"aggregation"`
	Operator    *ExpressionOperator `json:"operator"`
	Value       *string             `json:"value"`
}
