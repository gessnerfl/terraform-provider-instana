package restapi

// Severity representation of the severity in both worlds Instana API and Terraform Provider
type Severity struct {
	apiRepresentation       int
	terraformRepresentation string
}

// GetAPIRepresentation returns the integer representation of the Instana API
func (s Severity) GetAPIRepresentation() int { return s.apiRepresentation }

// GetTerraformRepresentation returns the string representation of the Terraform Provider
func (s Severity) GetTerraformRepresentation() string { return s.terraformRepresentation }

func (s Severity) equals(other Severity) bool {
	return s.apiRepresentation == other.apiRepresentation && s.terraformRepresentation == other.terraformRepresentation
}

// SeverityCritical representation of the critical severity
var SeverityCritical = Severity{apiRepresentation: 10, terraformRepresentation: "critical"}

// SeverityWarning representation of the warning severity
var SeverityWarning = Severity{apiRepresentation: 5, terraformRepresentation: "warning"}

// Severities custom type representing a slice of Severity
type Severities []Severity

// TerraformRepresentations returns the corresponding Terraform representations as string slice
func (severities Severities) TerraformRepresentations() []string {
	result := make([]string, len(severities))
	for i, v := range severities {
		result[i] = v.terraformRepresentation
	}
	return result
}

// APIRepresentations returns the corresponding Instana API representations as int slice
func (severities Severities) APIRepresentations() []int {
	result := make([]int, len(severities))
	for i, v := range severities {
		result[i] = v.apiRepresentation
	}
	return result
}

// SupportedSeverities slice of all supported severities of the Instana REST API
var SupportedSeverities = Severities{SeverityWarning, SeverityCritical}
