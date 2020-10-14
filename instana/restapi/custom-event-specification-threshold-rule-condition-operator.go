package restapi

import "fmt"

//ConditionOperator representation of a ConditionOperator of a threshold rule of a custom event specification  of the Instana Web REST API
type ConditionOperator interface {
	InstanaAPIValue() string
	TerraformSupportedValues() []string
}

func newBasicConditionOperator(instanaAPIValue string, additionalSupportedTerraformValues ...string) ConditionOperator {
	return &baseConditionOperator{instanaAPIValue: instanaAPIValue, terraformSupportedValues: append(additionalSupportedTerraformValues, instanaAPIValue)}
}

//ConditionOperatorType custom type representing a matching operator of a custom event specification rule
type baseConditionOperator struct {
	instanaAPIValue          string
	terraformSupportedValues []string
}

//InstanaAPIValue implementation of ConditionOperator interace
func (b *baseConditionOperator) InstanaAPIValue() string {
	return b.instanaAPIValue
}

//TerraformSupportedValues implementation of ConditionOperator interace
func (b *baseConditionOperator) TerraformSupportedValues() []string {
	return b.terraformSupportedValues
}

//ConditionOperators custom type representing a slice of ConditionOperatorType
type ConditionOperators []ConditionOperator

//TerrafromSupportedValues Returns the terraform string representations fo the matching operators
func (types ConditionOperators) TerrafromSupportedValues() []string {
	result := make([]string, 0)
	for _, v := range types {
		result = append(result, v.TerraformSupportedValues()...)
	}
	return result
}

//InstanaAPISupportedValues Returns the terraform string representations fo the matching operators
func (types ConditionOperators) InstanaAPISupportedValues() []string {
	result := make([]string, len(types))
	for i, v := range types {
		result[i] = v.InstanaAPIValue()
	}
	return result
}

//IsSupportedInstanaAPIConditionOperator check if the provided matching operator type is a supported instana api value
func (types ConditionOperators) IsSupportedInstanaAPIConditionOperator(operator string) bool {
	for _, t := range types {
		if t.InstanaAPIValue() == operator {
			return true
		}
	}
	return false
}

//FromInstanaAPIValue returns the ConditionOperator for the given instana apistring value or an error if the operator type does not exist
func (types ConditionOperators) FromInstanaAPIValue(instanaAPIvalue string) (ConditionOperator, error) {
	for _, t := range types {
		if t.InstanaAPIValue() == instanaAPIvalue {
			return t, nil
		}
	}
	return ConditionOperatorEquals, fmt.Errorf("%s is not a supported condition operator of the Instana Web REST API", instanaAPIvalue)
}

//FromTerraformValue returns the ConditionOperator for the given terraform string value or an error if the operator type does not exist
func (types ConditionOperators) FromTerraformValue(terraformRepresentation string) (ConditionOperator, error) {
	for _, t := range types {
		for _, v := range t.TerraformSupportedValues() {
			if v == terraformRepresentation {
				return t, nil
			}
		}
	}
	return ConditionOperatorEquals, fmt.Errorf("%s is not a supported condition operator of the Instana Terraform provider", terraformRepresentation)
}

var (
	//ConditionOperatorEquals const for a equals (==) condition operator
	ConditionOperatorEquals = newBasicConditionOperator("=", "==")
	//ConditionOperatorNotEqual const for a not equal (!=) condition operator
	ConditionOperatorNotEqual = newBasicConditionOperator("!=")
	//ConditionOperatorLessThan const for a less than (<) condition operator
	ConditionOperatorLessThan = newBasicConditionOperator("<")
	//ConditionOperatorLessThanOrEqual const for a less than or equal (<=) condition operator
	ConditionOperatorLessThanOrEqual = newBasicConditionOperator("<=")
	//ConditionOperatorGreaterThan const for a greater than (>) condition operator
	ConditionOperatorGreaterThan = newBasicConditionOperator(">")
	//ConditionOperatorGreaterThanOrEqual const for a greater than or equal (<=) condition operator
	ConditionOperatorGreaterThanOrEqual = newBasicConditionOperator(">=")
)

//SupportedConditionOperators slice of supported matching operatorTypes types
var SupportedConditionOperators = ConditionOperators{ConditionOperatorEquals, ConditionOperatorNotEqual, ConditionOperatorLessThan, ConditionOperatorLessThanOrEqual, ConditionOperatorGreaterThan, ConditionOperatorGreaterThanOrEqual}
