package instana_test

import (
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/hashicorp/terraform/helper/schema"
)

func validateRequiredSchemaOfTypeString(schemaField string, schemaMap map[string]*schema.Schema, t *testing.T) {
	validateRequiredSchemaOfType(schemaField, schemaMap, schema.TypeString, t)
}

func validateRequiredSchemaOfTypeInt(schemaField string, schemaMap map[string]*schema.Schema, t *testing.T) {
	validateRequiredSchemaOfType(schemaField, schemaMap, schema.TypeInt, t)
}

func validateRequiredSchemaOfTypeFloat(schemaField string, schemaMap map[string]*schema.Schema, t *testing.T) {
	validateRequiredSchemaOfType(schemaField, schemaMap, schema.TypeFloat, t)
}

func validateRequiredSchemaOfType(schemaField string, schemaMap map[string]*schema.Schema, dataType schema.ValueType, t *testing.T) {
	s := schemaMap[schemaField]
	if s == nil {
		t.Fatalf("Expected configuration for %s", schemaField)
	}
	validateSchemaOfType(s, dataType, t)
	if !s.Required {
		t.Fatalf("Expected %s to be required", schemaField)
	}
}

func validateOptionalSchemaOfTypeString(schemaField string, schemaMap map[string]*schema.Schema, t *testing.T) {
	validateOptionalSchemaOfType(schemaField, schemaMap, schema.TypeString, t)
}

func validateOptionalSchemaOfTypeInt(schemaField string, schemaMap map[string]*schema.Schema, t *testing.T) {
	validateOptionalSchemaOfType(schemaField, schemaMap, schema.TypeInt, t)
}

func validateOptionalSchemaOfType(schemaField string, schemaMap map[string]*schema.Schema, dataType schema.ValueType, t *testing.T) {
	s := schemaMap[schemaField]
	if s == nil {
		t.Fatalf("Expected configuration for %s", schemaField)
	}
	validateSchemaOfType(s, dataType, t)
	if !s.Optional {
		t.Fatalf("Expected %s to be optional", schemaField)
	}
}

func validateSchemaOfTypeBoolWithDefault(schemaField string, defaultValue bool, schemaMap map[string]*schema.Schema, t *testing.T) {
	s := schemaMap[schemaField]
	if s == nil {
		t.Fatalf("Expected configuration for %s", schemaField)
	}
	validateSchemaOfType(s, schema.TypeBool, t)
	if s.Required {
		t.Fatalf("Expected %s to be optional", schemaField)
	}
	if s.Default != defaultValue {
		t.Fatalf("Expected default value %t", defaultValue)
	}
}

func validateSchemaOfType(s *schema.Schema, dataType schema.ValueType, t *testing.T) {
	if s.Type != dataType {
		t.Fatalf("Expected field to be of type %d", dataType)
	}
	if len(s.Description) == 0 {
		t.Fatal("Expected description for schema")
	}
}

func validateRequiredSchemaOfTypeListOfString(schemaField string, schemaMap map[string]*schema.Schema, t *testing.T) {
	s := schemaMap[schemaField]
	if s == nil {
		t.Fatalf("Expected configuration for %s", schemaField)
	}
	if s.Type != schema.TypeList {
		t.Fatal("Expected field to be of type list")
	}
	if s.Elem == nil {
		t.Fatal("Expected list element to be defined")
	}
	if s.Elem.(schema.Schema).Type != schema.TypeString {
		t.Fatal("Expected list element to be of type string")
	}
	if len(s.Description) == 0 {
		t.Fatal("Expected description for schema")
	}
	if s.Required {
		t.Fatalf("Expected %s to be optional", schemaField)
	}
}

func createEmptyRuleResourceData(t *testing.T) *schema.ResourceData {
	data := make(map[string]interface{})
	return createRuleResourceData(t, data)
}

func createRuleResourceData(t *testing.T, data map[string]interface{}) *schema.ResourceData {
	schemaMap := instana.CreateResourceRule().Schema
	return schema.TestResourceDataRaw(t, schemaMap, data)
}

func createEmptyRuleBindingResourceData(t *testing.T) *schema.ResourceData {
	data := make(map[string]interface{})
	return createRuleBindingResourceData(t, data)
}

func createRuleBindingResourceData(t *testing.T, data map[string]interface{}) *schema.ResourceData {
	schemaMap := instana.CreateResourceRuleBinding().Schema
	return schema.TestResourceDataRaw(t, schemaMap, data)
}

func createEmptyUserRoleResourceData(t *testing.T) *schema.ResourceData {
	data := make(map[string]interface{})
	return createUserRoleResourceData(t, data)
}

func createUserRoleResourceData(t *testing.T, data map[string]interface{}) *schema.ResourceData {
	schemaMap := instana.CreateResourceUserRole().Schema
	return schema.TestResourceDataRaw(t, schemaMap, data)
}
