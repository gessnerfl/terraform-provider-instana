package instana

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rs/xid"
)

// RandomID generates a random ID for a resource
func RandomID() string {
	id := xid.New()
	return id.String()
}

// ReadArrayParameterFromMap reads an array parameter from a resource
func ReadArrayParameterFromMap[T any](data map[string]interface{}, key string) []T {
	if attr, ok := data[key]; ok {
		var array []T
		items := attr.([]interface{})
		for _, x := range items {
			item := x.(T)
			array = append(array, item)
		}
		return array
	}
	return nil
}

// ReadStringSetParameterFromResource reads a string set parameter from a resource and returns it as a slice of strings
func ReadStringSetParameterFromResource(d *schema.ResourceData, key string) []string {
	if attr, ok := d.GetOk(key); ok {
		var array []string
		set := attr.(*schema.Set)
		for _, x := range set.List() {
			item := x.(string)
			array = append(array, item)
		}
		return array
	}
	return nil
}

// ReadSetParameterFromMap reads a set parameter from a map and returns it as a slice
func ReadSetParameterFromMap[T any](data map[string]interface{}, key string) []T {
	if attr, ok := data[key]; ok {
		var array []T
		set := attr.(*schema.Set)
		for _, x := range set.List() {
			item := x.(T)
			array = append(array, item)
		}
		return array
	}
	return nil
}

// ConvertSeverityFromInstanaAPIToTerraformRepresentation converts the integer representation of the Instana API to the string representation of the Terraform provider
func ConvertSeverityFromInstanaAPIToTerraformRepresentation(severity int) (string, error) {
	if severity == restapi.SeverityWarning.GetAPIRepresentation() {
		return restapi.SeverityWarning.GetTerraformRepresentation(), nil
	} else if severity == restapi.SeverityCritical.GetAPIRepresentation() {
		return restapi.SeverityCritical.GetTerraformRepresentation(), nil
	} else {
		return "INVALID", fmt.Errorf("%d is not a valid severity", severity)
	}
}

// ConvertSeverityFromTerraformToInstanaAPIRepresentation converts the string representation of the Terraform to the int representation of the Instana API provider
func ConvertSeverityFromTerraformToInstanaAPIRepresentation(severity string) (int, error) {
	if severity == restapi.SeverityWarning.GetTerraformRepresentation() {
		return restapi.SeverityWarning.GetAPIRepresentation(), nil
	} else if severity == restapi.SeverityCritical.GetTerraformRepresentation() {
		return restapi.SeverityCritical.GetAPIRepresentation(), nil
	} else {
		return -1, fmt.Errorf("%s is not a valid severity", severity)
	}
}

// GetIntPointerFromResourceData gets a int value from the resource data and either returns a pointer to the value or nil if the value is not defined
func GetIntPointerFromResourceData(d *schema.ResourceData, key string) *int {
	val, ok := d.GetOk(key)
	if ok {
		intValue := val.(int)
		return &intValue
	}
	return nil
}

// GetInt32PointerFromResourceData gets a int32 value from the resource data and either returns a pointer to the value or nil if the value is not defined
func GetInt32PointerFromResourceData(d *schema.ResourceData, key string) *int32 {
	val, ok := d.GetOk(key)
	if ok {
		intValue := int32(val.(int))
		return &intValue
	}
	return nil
}

// GetFloat64PointerFromResourceData gets a float64 value from the resource data and either returns a pointer to the value or nil if the value is not defined
func GetFloat64PointerFromResourceData(d *schema.ResourceData, key string) *float64 {
	val, ok := d.GetOk(key)
	if ok {
		floatValue := val.(float64)
		return &floatValue
	}
	return nil
}

// GetFloat32PointerFromResourceData gets a float32 value from the resource data and either returns a pointer to the value or nil if the value is not defined
func GetFloat32PointerFromResourceData(d *schema.ResourceData, key string) *float32 {
	val, ok := d.GetOk(key)
	if ok {
		floatValue := float32(val.(float64))
		return &floatValue
	}
	return nil
}

// GetStringPointerFromResourceData gets a string value from the resource data and either returns a pointer to the value or nil if the value is not defined
func GetStringPointerFromResourceData(d *schema.ResourceData, key string) *string {
	val, ok := d.GetOk(key)
	if ok {
		stringValue := val.(string)
		return &stringValue
	}
	return nil
}

// GetPointerFromMap gets a pointer of the given type from the map or nil if the value is not defined
func GetPointerFromMap[T any](d map[string]interface{}, key string) *T {
	var defaultValue T
	val, ok := d[key]
	if ok && val != nil && !cmp.Equal(defaultValue, val) {
		value := val.(T)
		return &value
	}
	return nil
}

// MergeSchemaMap merges the provided maps into a single map
func MergeSchemaMap(mapA map[string]*schema.Schema, mapB map[string]*schema.Schema) map[string]*schema.Schema {
	mergedMap := make(map[string]*schema.Schema)

	for k, v := range mapA {
		mergedMap[k] = v
	}
	for k, v := range mapB {
		mergedMap[k] = v
	}

	return mergedMap
}

// ConvertInterfaceSlice converts the given interface slice to the desired target slice
func ConvertInterfaceSlice[T any](input []interface{}) []T {
	result := make([]T, len(input))
	for i, v := range input {
		result[i] = v.(T)
	}
	return result
}

func NormalizeJSONString(jsonString string) string {
	var raw []map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &raw)
	if err != nil {
		return jsonString
	}
	bytes, err := json.Marshal(raw)
	if err != nil {
		return jsonString
	}
	return string(bytes)
}
