package testutils

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
)

const (
	messageExpectedToBeNotRequired = "Expected %s to be not required"
	messageExpectedToBeRequired    = "Expected %s to be required"
	messageExpectedToBeOptional    = "Expected %s to be optional"
	messageExpectedDefaultValue    = "Expected default value %t"
)

//NewTerraformSchemaAssert creates a new instance of TerraformSchemaAssert
func NewTerraformSchemaAssert(schemaMap map[string]*schema.Schema, t *testing.T) TerraformSchemaAssert {
	return &terraformSchemaAssertImpl{schemaMap: schemaMap, t: t}
}

//TerraformSchemaAssert a test util to verify terraform schema fields
type TerraformSchemaAssert interface {
	//AssertSchemaIsRequiredAndOfTypeString checks if the given schema field is required and of type string
	AssertSchemaIsRequiredAndOfTypeString(fieldName string)
	//AssertSchemaIsRequiredAndTypeInt checks if the given schema field is required and of type int
	AssertSchemaIsRequiredAndOfTypeInt(fieldName string)
	//AssertSchemaIsRequiredAndOfTypeFloat checks if the given schema field is required and of type float
	AssertSchemaIsRequiredAndOfTypeFloat(fieldName string)
	//AssertSchemaIsOptionalAndOfTypeString checks if the given schema field is optional and of type string
	AssertSchemaIsOptionalAndOfTypeString(fieldName string)
	//AssertSchemaIsOptionalAndOfTypeStringWithDefault checks if the given schema field is optional and of type string and has the given default value
	AssertSchemaIsOptionalAndOfTypeStringWithDefault(fieldName string, defaultValue string)
	//AssertSchemaIsOptionalAndOfTypeInt checks if the given schema field is optional and of type int
	AssertSchemaIsOptionalAndOfTypeInt(fieldName string)
	//AssertSchemaIsOptionalAndOfTypeFloat checks if the given schema field is required and of type float
	AssertSchemaIsOptionalAndOfTypeFloat(fieldName string)
	//AssertSchemaIsOfTypeBooleanWithDefault checks if the given schema field is an optional boolean field with an expected default value
	AssertSchemaIsOfTypeBooleanWithDefault(fieldName string, defaultValue bool)
	//AssertSchemaIsRequiredAndOfTypeListOfStrings checks if the given schema field is required and of type list of string
	AssertSchemaIsRequiredAndOfTypeListOfStrings(fieldName string)
	//AssertSchemaIsRequiredAndOfTypeListOfStrings checks if the given schema field is required and of type list of string
	AssertSChemaIsOptionalAndOfTypeListOfStrings(fieldName string)
}

type terraformSchemaAssertImpl struct {
	schemaMap map[string]*schema.Schema
	t         *testing.T
}

func (inst *terraformSchemaAssertImpl) AssertSchemaIsRequiredAndOfTypeString(schemaField string) {
	inst.AssertSchemaIsRequiredAndType(schemaField, schema.TypeString)
}

func (inst *terraformSchemaAssertImpl) AssertSchemaIsRequiredAndOfTypeInt(schemaField string) {
	inst.AssertSchemaIsRequiredAndType(schemaField, schema.TypeInt)
}

func (inst *terraformSchemaAssertImpl) AssertSchemaIsRequiredAndOfTypeFloat(schemaField string) {
	inst.AssertSchemaIsRequiredAndType(schemaField, schema.TypeFloat)
}

func (inst *terraformSchemaAssertImpl) AssertSchemaIsRequiredAndType(schemaField string, dataType schema.ValueType) {
	s := inst.schemaMap[schemaField]
	if s == nil {
		inst.t.Fatalf(ExpectedNoErrorButGotMessage, schemaField)
	}
	inst.assertSchemaIsOfType(s, dataType)
	if !s.Required {
		inst.t.Fatalf(messageExpectedToBeRequired, schemaField)
	}
}

func (inst *terraformSchemaAssertImpl) AssertSchemaIsOptionalAndOfTypeString(schemaField string) {
	inst.assertSchemaIsOptionalAndOfType(schemaField, schema.TypeString)
}

func (inst *terraformSchemaAssertImpl) AssertSchemaIsOptionalAndOfTypeStringWithDefault(schemaField string, defaultValue string) {
	inst.assertSchemaIsOptionalAndOfType(schemaField, schema.TypeString)
	s := inst.schemaMap[schemaField]
	if s.Default != defaultValue {
		inst.t.Fatalf("Expected default value %s", defaultValue)
	}
}

func (inst *terraformSchemaAssertImpl) AssertSchemaIsOptionalAndOfTypeInt(schemaField string) {
	inst.assertSchemaIsOptionalAndOfType(schemaField, schema.TypeInt)
}

func (inst *terraformSchemaAssertImpl) AssertSchemaIsOptionalAndOfTypeFloat(schemaField string) {
	inst.assertSchemaIsOptionalAndOfType(schemaField, schema.TypeFloat)
}

func (inst *terraformSchemaAssertImpl) assertSchemaIsOptionalAndOfType(schemaField string, dataType schema.ValueType) {
	s := inst.schemaMap[schemaField]
	if s == nil {
		inst.t.Fatalf(ExpectedNoErrorButGotMessage, schemaField)
	}
	inst.assertSchemaIsOfType(s, dataType)
	if !s.Optional {
		inst.t.Fatalf(messageExpectedToBeOptional, schemaField)
	}
}

func (inst *terraformSchemaAssertImpl) AssertSchemaIsOfTypeBooleanWithDefault(schemaField string, defaultValue bool) {
	s := inst.schemaMap[schemaField]
	if s == nil {
		inst.t.Fatalf(ExpectedNoErrorButGotMessage, schemaField)
	}
	inst.assertSchemaIsOfType(s, schema.TypeBool)
	if s.Required {
		inst.t.Fatalf(messageExpectedToBeOptional, schemaField)
	}
	if s.Default != defaultValue {
		inst.t.Fatalf(messageExpectedDefaultValue, defaultValue)
	}
}

func (inst *terraformSchemaAssertImpl) assertSchemaIsOfType(s *schema.Schema, dataType schema.ValueType) {
	if s.Type != dataType {
		inst.t.Fatalf("Expected field to be of type %d", dataType)
	}
	if len(s.Description) == 0 {
		inst.t.Fatal("Expected description for schema")
	}
}

func (inst *terraformSchemaAssertImpl) AssertSChemaIsOptionalAndOfTypeListOfStrings(schemaField string) {
	s := inst.schemaMap[schemaField]
	if s == nil {
		inst.t.Fatalf(ExpectedNoErrorButGotMessage, schemaField)
	}
	if s.Required {
		inst.t.Fatalf(messageExpectedToBeNotRequired, schemaField)
	}
	if !s.Optional {
		inst.t.Fatalf(messageExpectedToBeOptional, schemaField)
	}
	inst.assertSchemaIsOfTypeListOfStrings(s)
}

func (inst *terraformSchemaAssertImpl) AssertSchemaIsRequiredAndOfTypeListOfStrings(schemaField string) {
	s := inst.schemaMap[schemaField]
	if s == nil {
		inst.t.Fatalf(ExpectedNoErrorButGotMessage, schemaField)
	}
	if !s.Required {
		inst.t.Fatalf(messageExpectedToBeRequired, schemaField)
	}
	inst.assertSchemaIsOfTypeListOfStrings(s)
}

func (inst *terraformSchemaAssertImpl) assertSchemaIsOfTypeListOfStrings(s *schema.Schema) {
	if s.Type != schema.TypeList {
		inst.t.Fatal("Expected field to be of type list")
	}
	if s.Elem == nil {
		inst.t.Fatal("Expected list element to be defined")
	}
	if s.Elem.(*schema.Schema).Type != schema.TypeString {
		inst.t.Fatal("Expected list element to be of type string")
	}
	if len(s.Description) == 0 {
		inst.t.Fatal("Expected description for schema")
	}
}
