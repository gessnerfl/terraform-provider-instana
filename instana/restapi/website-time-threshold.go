package restapi

// WebsiteTimeThreshold struct representing the API model of a website time threshold
type WebsiteTimeThreshold struct {
	Type                    string                          `json:"type"`
	TimeWindow              *int64                          `json:"timeWindow"`
	Violations              *int32                          `json:"violations"`
	ImpactMeasurementMethod *WebsiteImpactMeasurementMethod `json:"impactMeasurementMethod"`
	UserPercentage          *float64                        `json:"userPercentage"`
	Users                   *int32                          `json:"users"`
}
