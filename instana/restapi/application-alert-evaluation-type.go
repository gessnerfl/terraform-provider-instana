package restapi

// ApplicationAlertEvaluationType custom type representing the application alert evaluation type from the instana API
type ApplicationAlertEvaluationType string

// ApplicationAlertEvaluationTypes custom type representing a slice of ApplicationAlertEvaluationType
type ApplicationAlertEvaluationTypes []ApplicationAlertEvaluationType

// ToStringSlice Returns the corresponding string representations
func (types ApplicationAlertEvaluationTypes) ToStringSlice() []string {
	result := make([]string, len(types))
	for i, v := range types {
		result[i] = string(v)
	}
	return result
}

const (
	//EvaluationTypePerApplication constant value for ApplicationAlertEvaluationType PER_AP
	EvaluationTypePerApplication = ApplicationAlertEvaluationType("PER_AP")
	//EvaluationTypePerApplicationService constant value for ApplicationAlertEvaluationType PER_AP_SERVICE
	EvaluationTypePerApplicationService = ApplicationAlertEvaluationType("PER_AP_SERVICE")
	//EvaluationTypePerApplicationEndpoint constant value for ApplicationAlertEvaluationType PER_AP_ENDPOINT
	EvaluationTypePerApplicationEndpoint = ApplicationAlertEvaluationType("PER_AP_ENDPOINT")
)

// SupportedApplicationAlertEvaluationTypes list of all supported ApplicationAlertEvaluationTypes
var SupportedApplicationAlertEvaluationTypes = ApplicationAlertEvaluationTypes{EvaluationTypePerApplication, EvaluationTypePerApplicationService, EvaluationTypePerApplicationEndpoint}
