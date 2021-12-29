package restapi

//ApplicationAlertEvaluationType custom type representing the application alert evaluation type from the instana API
type ApplicationAlertEvaluationType string

//ApplicationAlertEvaluationTypes custom type representing a slice of ApplicationAlertEvaluationType
type ApplicationAlertEvaluationTypes []ApplicationAlertEvaluationType

//IsSupported check if the provided evaluation type is supported
func (types ApplicationAlertEvaluationTypes) IsSupported(evalType ApplicationAlertEvaluationType) bool {
	for _, t := range types {
		if t == evalType {
			return true
		}
	}
	return false
}

//ToStringSlice Returns the corresponding string representations
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
)

//SupportedApplicationAlertEvaluationTypes list of all supported ApplicationAlertEvaluationTypes
var SupportedApplicationAlertEvaluationTypes = ApplicationAlertEvaluationTypes{EvaluationTypePerApplication, EvaluationTypePerApplicationService}
