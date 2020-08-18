package restapi

import "fmt"

//MatchingOperator representation of a MatchingOperator of a threshold rule of a custom event specification  of the Instana Web REST API
type MatchingOperator interface {
	InstanaAPIValue() string
	TerraformSupportedValues() []string
}

func newBasicMatchingOperator(instanaAPIValue string, additionalSupportedTerraformValues ...string) MatchingOperator {
	return &baseMatchingOperator{instanaAPIValue: instanaAPIValue, terraformSupportedValues: append(additionalSupportedTerraformValues, instanaAPIValue)}
}

//MatchingOperatorType custom type representing a matching operator of a custom event specification rule
type baseMatchingOperator struct {
	instanaAPIValue          string
	terraformSupportedValues []string
}

//InstanaAPIValue implementation of MatchingOperator interace
func (b *baseMatchingOperator) InstanaAPIValue() string {
	return b.instanaAPIValue
}

//TerraformSupportedValues implementation of MatchingOperator interace
func (b *baseMatchingOperator) TerraformSupportedValues() []string {
	return b.terraformSupportedValues
}

//MatchingOperators custom type representing a slice of MatchingOperatorType
type MatchingOperators []MatchingOperator

//TerrafromSupportedValues Returns the terraform string representations fo the matching operators
func (types MatchingOperators) TerrafromSupportedValues() []string {
	result := make([]string, 0)
	for _, v := range types {
		result = append(result, v.TerraformSupportedValues()...)
	}
	return result
}

//InstanaAPISupportedValues Returns the terraform string representations fo the matching operators
func (types MatchingOperators) InstanaAPISupportedValues() []string {
	result := make([]string, len(types))
	for i, v := range types {
		result[i] = v.InstanaAPIValue()
	}
	return result
}

//IsSupportedInstanaAPIMatchingOperator check if the provided matching operator type is a supported instana api value
func (types MatchingOperators) IsSupportedInstanaAPIMatchingOperator(operator string) bool {
	for _, t := range types {
		if t.InstanaAPIValue() == operator {
			return true
		}
	}
	return false
}

//FromInstanaAPIValue returns the MatchingOperator for the given instana apistring value or an error if the operator type does not exist
func (types MatchingOperators) FromInstanaAPIValue(instanaAPIvalue string) (MatchingOperator, error) {
	for _, t := range types {
		if t.InstanaAPIValue() == instanaAPIvalue {
			return t, nil
		}
	}
	return MatchingOperatorIs, fmt.Errorf("%s is not a supported matching operator of the Instana Web REST API", instanaAPIvalue)
}

//FromTerraformValue returns the MatchingOperator for the given terraform string value or an error if the operator type does not exist
func (types MatchingOperators) FromTerraformValue(terraformRepresentation string) (MatchingOperator, error) {
	for _, t := range types {
		for _, v := range t.TerraformSupportedValues() {
			if v == terraformRepresentation {
				return t, nil
			}
		}
	}
	return MatchingOperatorIs, fmt.Errorf("%s is not a supported matching operator of the Instana Terraform provider", terraformRepresentation)
}

var (
	//MatchingOperatorIs const for IS condition operator
	MatchingOperatorIs = newBasicMatchingOperator("is")
	//MatchingOperatorContains const for CONTAINS condition operator
	MatchingOperatorContains = newBasicMatchingOperator("contains")
	//MatchingOperatorStartsWith const for STARTS_WITH condition operator
	MatchingOperatorStartsWith = newBasicMatchingOperator("startsWith", "starts_with")
	//MatchingOperatorEndsWith const for ENDS_WITH condition operator
	MatchingOperatorEndsWith = newBasicMatchingOperator("endsWith", "ends_with")
)

//SupportedMatchingOperators slice of supported matching operatorTypes types
var SupportedMatchingOperators = MatchingOperators{MatchingOperatorIs, MatchingOperatorContains, MatchingOperatorStartsWith, MatchingOperatorEndsWith}
