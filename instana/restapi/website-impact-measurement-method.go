package restapi

// WebsiteImpactMeasurementMethod custom type for impact measurement method of website alert rules
type WebsiteImpactMeasurementMethod string

// WebsiteImpactMeasurementMethods custom type for a slice of WebsiteImpactMeasurementMethod
type WebsiteImpactMeasurementMethods []WebsiteImpactMeasurementMethod

// ToStringSlice Returns the corresponding string representations
func (methods WebsiteImpactMeasurementMethods) ToStringSlice() []string {
	result := make([]string, len(methods))
	for i, v := range methods {
		result[i] = string(v)
	}
	return result
}

const (
	//WebsiteImpactMeasurementMethodAggregated constant value for the website impact measurement method aggregated
	WebsiteImpactMeasurementMethodAggregated = WebsiteImpactMeasurementMethod("AGGREGATED")
	//WebsiteImpactMeasurementMethodPerWindow constant value for the website impact measurement method per_window
	WebsiteImpactMeasurementMethodPerWindow = WebsiteImpactMeasurementMethod("PER_WINDOW")
)

// SupportedWebsiteImpactMeasurementMethods list of all supported WebsiteImpactMeasurementMethod
var SupportedWebsiteImpactMeasurementMethods = WebsiteImpactMeasurementMethods{WebsiteImpactMeasurementMethodAggregated, WebsiteImpactMeasurementMethodPerWindow}
