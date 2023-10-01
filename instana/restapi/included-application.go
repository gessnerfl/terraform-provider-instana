package restapi

// IncludedEndpoint custom type to include of a specific endpoint in an alert config
type IncludedEndpoint struct {
	EndpointID string `json:"endpointId"`
	Inclusive  bool   `json:"inclusive"`
}

// IncludedService custom type to include of a specific service in an alert config
type IncludedService struct {
	ServiceID string `json:"serviceId"`
	Inclusive bool   `json:"inclusive"`

	Endpoints map[string]IncludedEndpoint `json:"endpoints"`
}

// IncludedApplication custom type to include specific applications in an alert config
type IncludedApplication struct {
	ApplicationID string `json:"applicationId"`
	Inclusive     bool   `json:"inclusive"`

	Services map[string]IncludedService `json:"services"`
}
