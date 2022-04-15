package restapi

//CustomPayloadType custom type for the type of a custom payload
type CustomPayloadType string

//CustomPayloadTypes custom type for a slice of CustomPayloadType
type CustomPayloadTypes []CustomPayloadType

//IsSupported check if the provided CustomPayloadType is supported
func (types CustomPayloadTypes) IsSupported(cpt CustomPayloadType) bool {
	for _, t := range types {
		if t == cpt {
			return true
		}
	}
	return false
}

//ToStringSlice Returns the corresponding string representations
func (types CustomPayloadTypes) ToStringSlice() []string {
	result := make([]string, len(types))
	for i, v := range types {
		result[i] = string(v)
	}
	return result
}

const (
	//StaticCustomPayloadType constant value for the static CustomPayloadType
	StaticCustomPayloadType = CustomPayloadType("staticString")
	//DynamicCustomPayloadType constant value for the dynamic CustomPayloadType
	DynamicCustomPayloadType = CustomPayloadType("dynamic")
)

//SupportedCustomPayloadTypes list of all supported CustomPayloadType
var SupportedCustomPayloadTypes = CustomPayloadTypes{StaticCustomPayloadType, DynamicCustomPayloadType}

//StaticStringCustomPayloadFieldValue type for static string values of custom payload field
type StaticStringCustomPayloadFieldValue string

//DynamicCustomPayloadFieldValue type for dynamic values of custom payload field
type DynamicCustomPayloadFieldValue struct {
	Key     *string `json:"key"`
	TagName string  `json:"tagName"`
}

//CustomPayloadField custom type to represent static fields with a string value for custom payloads
type CustomPayloadField[T any] struct {
	Type  CustomPayloadType `json:"type"`
	Key   string            `json:"key"`
	Value T                 `json:"value"`
}
