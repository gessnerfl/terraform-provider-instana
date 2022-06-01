package testutils

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//NewTerraformSchemaAssert creates a new instance of TerraformSchemaAssert
func NewTerraformSchemaAssert(schemaMap map[string]*schema.Schema, t *testing.T) TerraformSchemaAssert {
	return &terraformSchemaAssertImpl{schemaMap: schemaMap, t: t}
}

//TerraformSchemaAssert a test util to verify terraform schema fields
type TerraformSchemaAssert interface {
	//AssertSchemaIsRequiredAndOfTypeString checks if the given schema field is required and of type string
	AssertSchemaIsRequiredAndOfTypeString(fieldName string)
	//AssertSchemaIsRequiredAndOfTypeInt checks if the given schema field is required and of type int
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
	//AssertSchemaIsOptionalAndOfTypeListOfStrings checks if the given schema field is optional and of type list of string
	AssertSchemaIsOptionalAndOfTypeListOfStrings(fieldName string)
	//AssertSchemaIsRequiredAndOfTypeSetOfStrings checks if the given schema field is required and of type set of string
	AssertSchemaIsRequiredAndOfTypeSetOfStrings(fieldName string)
	//AssertSchemaIsOptionalAndOfTypeSetOfStrings checks if the given schema field is required and of type set of string
	AssertSchemaIsOptionalAndOfTypeSetOfStrings(fieldName string)
	//AssertSchemaIsComputedAndOfTypeString checks if the given schema field is computed and of type string
	AssertSchemaIsComputedAndOfTypeString(fieldName string)
	//AssertSchemaIsComputedAndOfTypeInt checks if the given schema field is computed and of type int
	AssertSchemaIsComputedAndOfTypeInt(fieldName string)
	//AssertSchemaIsComputedAndOfTypeBool checks if the given schema field is computed and of type bool
	AssertSchemaIsComputedAndOfTypeBool(fieldName string)
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
	require.NotNil(inst.t, s)

	inst.assertSchemaIsOfType(s, dataType)
	require.True(inst.t, s.Required)
}

func (inst *terraformSchemaAssertImpl) AssertSchemaIsOptionalAndOfTypeString(schemaField string) {
	inst.assertSchemaIsOptionalAndOfType(schemaField, schema.TypeString)
}

func (inst *terraformSchemaAssertImpl) AssertSchemaIsOptionalAndOfTypeStringWithDefault(schemaField string, defaultValue string) {
	inst.assertSchemaIsOptionalAndOfType(schemaField, schema.TypeString)
	s := inst.schemaMap[schemaField]

	require.Equal(inst.t, defaultValue, s.Default)
}

func (inst *terraformSchemaAssertImpl) AssertSchemaIsOptionalAndOfTypeInt(schemaField string) {
	inst.assertSchemaIsOptionalAndOfType(schemaField, schema.TypeInt)
}

func (inst *terraformSchemaAssertImpl) AssertSchemaIsOptionalAndOfTypeFloat(schemaField string) {
	inst.assertSchemaIsOptionalAndOfType(schemaField, schema.TypeFloat)
}

func (inst *terraformSchemaAssertImpl) assertSchemaIsOptionalAndOfType(schemaField string, dataType schema.ValueType) {
	s := inst.schemaMap[schemaField]
	require.NotNil(inst.t, s)

	inst.assertSchemaIsOfType(s, dataType)
	require.True(inst.t, s.Optional)
}

func (inst *terraformSchemaAssertImpl) AssertSchemaIsOfTypeBooleanWithDefault(schemaField string, defaultValue bool) {
	s := inst.schemaMap[schemaField]
	require.NotNil(inst.t, s)

	inst.assertSchemaIsOfType(s, schema.TypeBool)

	require.False(inst.t, s.Required)
	require.Equal(inst.t, defaultValue, s.Default)
}

func (inst *terraformSchemaAssertImpl) assertSchemaIsOfType(s *schema.Schema, dataType schema.ValueType) {
	require.Equal(inst.t, dataType, s.Type)
	require.Greater(inst.t, len(s.Description), 0)
}

func (inst *terraformSchemaAssertImpl) AssertSchemaIsOptionalAndOfTypeListOfStrings(schemaField string) {
	s := inst.schemaMap[schemaField]

	require.NotNil(inst.t, s)
	require.False(inst.t, s.Required)
	require.True(inst.t, s.Optional)
	inst.assertSchemaIsOfTypeListOfStrings(s)
}

func (inst *terraformSchemaAssertImpl) AssertSchemaIsRequiredAndOfTypeListOfStrings(schemaField string) {
	s := inst.schemaMap[schemaField]

	require.NotNil(inst.t, s)
	require.True(inst.t, s.Required)
	inst.assertSchemaIsOfTypeListOfStrings(s)
}

func (inst *terraformSchemaAssertImpl) assertSchemaIsOfTypeListOfStrings(s *schema.Schema) {
	require.Equal(inst.t, schema.TypeList, s.Type)
	require.NotNil(inst.t, s.Elem)
	require.Equal(inst.t, schema.TypeString, s.Elem.(*schema.Schema).Type)
	require.Greater(inst.t, len(s.Description), 0)
}

func (inst *terraformSchemaAssertImpl) AssertSchemaIsOptionalAndOfTypeSetOfStrings(schemaField string) {
	s := inst.schemaMap[schemaField]

	require.NotNil(inst.t, s)
	require.False(inst.t, s.Required)
	require.True(inst.t, s.Optional)
	inst.assertSchemaIsOfTypeSetOfStrings(s)
}

func (inst *terraformSchemaAssertImpl) AssertSchemaIsRequiredAndOfTypeSetOfStrings(schemaField string) {
	s := inst.schemaMap[schemaField]

	require.NotNil(inst.t, s)
	require.True(inst.t, s.Required)
	inst.assertSchemaIsOfTypeSetOfStrings(s)
}

func (inst *terraformSchemaAssertImpl) assertSchemaIsOfTypeSetOfStrings(s *schema.Schema) {
	require.Equal(inst.t, schema.TypeSet, s.Type)
	require.NotNil(inst.t, s.Elem)
	require.Equal(inst.t, schema.TypeString, s.Elem.(*schema.Schema).Type)
	require.Greater(inst.t, len(s.Description), 0)
}

func (inst *terraformSchemaAssertImpl) AssertSchemaIsComputedAndOfTypeString(schemaField string) {
	s := inst.schemaMap[schemaField]

	require.NotNil(inst.t, s)
	inst.assertSchemaIsOfType(s, schema.TypeString)
	require.True(inst.t, s.Computed)
}

func (inst *terraformSchemaAssertImpl) AssertSchemaIsComputedAndOfTypeInt(schemaField string) {
	s := inst.schemaMap[schemaField]

	require.NotNil(inst.t, s)
	inst.assertSchemaIsOfType(s, schema.TypeInt)
	require.True(inst.t, s.Computed)
}

func (inst *terraformSchemaAssertImpl) AssertSchemaIsComputedAndOfTypeBool(schemaField string) {
	s := inst.schemaMap[schemaField]

	require.NotNil(inst.t, s)
	inst.assertSchemaIsOfType(s, schema.TypeBool)
	require.True(inst.t, s.Computed)
}
