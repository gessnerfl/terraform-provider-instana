package instana

import (
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	//DefaultCustomPayloadFieldsName constant value for field custom_payload_field
	DefaultCustomPayloadFieldsName = "custom_payload_field"
	//CustomPayloadFieldsFieldKey constant value for field key
	CustomPayloadFieldsFieldKey = "key"
	//CustomPayloadFieldsFieldStaticStringValue constant value for field value
	CustomPayloadFieldsFieldStaticStringValue = "value"
	//CustomPayloadFieldsFieldDynamicValue constant value for field dynamic_value
	CustomPayloadFieldsFieldDynamicValue = "dynamic_value"
	//CustomPayloadFieldsFieldDynamicKey constant value for field key of dynamic custom payload fields
	CustomPayloadFieldsFieldDynamicKey = "key"
	//CustomPayloadFieldsFieldDynamicTagName constant value for field tag_name of dynamic custom payload fields
	CustomPayloadFieldsFieldDynamicTagName = "tag_name"
)

func buildCustomPayloadFields() *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeSet,
		Set: func(i interface{}) int {
			return schema.HashString(i.(map[string]interface{})[CustomPayloadFieldsFieldKey])
		},
		Optional: true,
		MinItems: 0,
		MaxItems: 20,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				CustomPayloadFieldsFieldKey: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The key of the custom payload field",
				},
				CustomPayloadFieldsFieldStaticStringValue: {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The value of a static string custom payload field",
				},
				CustomPayloadFieldsFieldDynamicValue: {
					Type:        schema.TypeList,
					MinItems:    0,
					MaxItems:    1,
					Optional:    true,
					Description: "The value of a dynamic custom payload field",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							CustomPayloadFieldsFieldDynamicKey: {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "The key of the dynamic custom payload field",
							},
							CustomPayloadFieldsFieldDynamicTagName: {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The name of the tag of the dynamic custom payload field",
							},
						},
					},
				},
			},
		},
		Description: "A list of custom payload fields",
	}
}

func mapCustomPayloadFieldsToSchema(input restapi.CustomPayloadFieldsAware) []map[string]interface{} {
	customPayloadField := input.GetCustomerPayloadFields()
	result := make([]map[string]interface{}, len(customPayloadField))
	for i, v := range customPayloadField {
		field := make(map[string]interface{})
		field[CustomPayloadFieldsFieldKey] = v.Key
		if v.Type == restapi.DynamicCustomPayloadType {
			value := v.Value.(restapi.DynamicCustomPayloadFieldValue)
			field[CustomPayloadFieldsFieldDynamicValue] = []interface{}{map[string]interface{}{
				CustomPayloadFieldsFieldDynamicKey:     value.Key,
				CustomPayloadFieldsFieldDynamicTagName: value.TagName,
			}}
		} else {
			value := v.Value.(restapi.StaticStringCustomPayloadFieldValue)
			field[CustomPayloadFieldsFieldStaticStringValue] = string(value)
		}
		result[i] = field
	}
	return result
}

func mapDefaultCustomPayloadFieldsFromSchema(d *schema.ResourceData) ([]restapi.CustomPayloadField[any], error) {
	return mapCustomPayloadFieldsFromSchema(d, DefaultCustomPayloadFieldsName)
}

func mapCustomPayloadFieldsFromSchema(d *schema.ResourceData, key string) ([]restapi.CustomPayloadField[any], error) {
	val, ok := d.GetOk(key)
	if ok && val != nil {
		fields := val.(*schema.Set).List()
		result := make([]restapi.CustomPayloadField[any], len(fields))
		for i, v := range fields {
			field := v.(map[string]interface{})
			key := field[CustomPayloadFieldsFieldKey].(string)
			dynamicValue, dynamicValueExists := extractDynamicCustomPayloadFieldValue(field)
			staticStringValue, staticStringValueExists := extractStaticStringCustomPayloadFieldValue(field)

			if (dynamicValueExists && staticStringValueExists) || (!dynamicValueExists && !staticStringValueExists) {
				return []restapi.CustomPayloadField[any]{}, fmt.Errorf("either a static string value or a dynamic value must be defined for custom field with key %s", key)
			}

			if dynamicValueExists {
				var dynamicValueKey *string
				if val, ok := dynamicValue[CustomPayloadFieldsFieldDynamicKey]; ok {
					dynamicValueKeyString := val.(string)
					dynamicValueKey = &dynamicValueKeyString
				}
				result[i] = restapi.CustomPayloadField[any]{
					Type: restapi.DynamicCustomPayloadType,
					Key:  key,
					Value: restapi.DynamicCustomPayloadFieldValue{
						TagName: dynamicValue[CustomPayloadFieldsFieldDynamicTagName].(string),
						Key:     dynamicValueKey,
					},
				}
			} else if staticStringValueExists {
				result[i] = restapi.CustomPayloadField[any]{
					Type:  restapi.StaticStringCustomPayloadType,
					Key:   key,
					Value: restapi.StaticStringCustomPayloadFieldValue(staticStringValue),
				}
			}
		}
		return result, nil
	}
	return []restapi.CustomPayloadField[any]{}, nil
}

func extractDynamicCustomPayloadFieldValue(field map[string]interface{}) (map[string]interface{}, bool) {
	val, ok := field[CustomPayloadFieldsFieldDynamicValue]
	if ok && len(val.([]interface{})) > 0 {
		return val.([]interface{})[0].(map[string]interface{}), true
	}
	return nil, false
}

func extractStaticStringCustomPayloadFieldValue(field map[string]interface{}) (string, bool) {
	val, ok := field[CustomPayloadFieldsFieldStaticStringValue]
	if ok && len(val.(string)) > 0 {
		return val.(string), true
	}
	return "", false
}
