package restapi

const (
	//SliConfigResourcePath path to sli config resource of Instana RESTful API
	SliConfigResourcePath = SettingsBasePath + "/v2/sli"
)

// MetricConfiguration represents the nested object metric configuration of the sli config REST resource at Instana
type MetricConfiguration struct {
	Name        string  `json:"metricName"`
	Aggregation string  `json:"metricAggregation"`
	Threshold   float64 `json:"threshold"`
}

// SliEntity represents the nested object sli entity of the sli config REST resource at Instana
type SliEntity struct {
	Type                      string     `json:"sliType"`
	ApplicationID             *string    `json:"applicationId"`
	ServiceID                 *string    `json:"serviceId"`
	EndpointID                *string    `json:"endpointId"`
	BoundaryScope             *string    `json:"boundaryScope"`
	IncludeSynthetic          *bool      `json:"includeSynthetic"`
	IncludeInternal           *bool      `json:"includeInternal"`
	WebsiteId                 *string    `json:"websiteId"`
	BeaconType                *string    `json:"beaconType"`
	GoodEventFilterExpression *TagFilter `json:"goodEventFilterExpression"`
	BadEventFilterExpression  *TagFilter `json:"badEventFilterExpression"`
	FilterExpression          *TagFilter `json:"filterExpression"`
}

// SliConfig represents the REST resource of sli configuration at Instana
type SliConfig struct {
	ID                         string               `json:"id"`
	Name                       string               `json:"sliName"`
	InitialEvaluationTimestamp int                  `json:"initialEvaluationTimestamp"`
	MetricConfiguration        *MetricConfiguration `json:"metricConfiguration"`
	SliEntity                  SliEntity            `json:"sliEntity"`
}

// GetIDForResourcePath implementation of the interface InstanaDataObject
func (s *SliConfig) GetIDForResourcePath() string {
	return s.ID
}
