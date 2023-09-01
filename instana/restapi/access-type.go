package restapi

// AccessType custom type for access type
type AccessType string

// AccessTypes custom type for a slice of AccessType
type AccessTypes []AccessType

// ToStringSlice Returns the corresponding string representations
func (types AccessTypes) ToStringSlice() []string {
	result := make([]string, len(types))
	for i, v := range types {
		result[i] = string(v)
	}
	return result
}

const (
	//AccessTypeRead constant value for the READ AccessType
	AccessTypeRead = AccessType("READ")
	//AccessTypeReadWrite constant value for the READ_WRITE AccessType
	AccessTypeReadWrite = AccessType("READ_WRITE")
)

// SupportedAccessTypes list of all supported AccessType
var SupportedAccessTypes = AccessTypes{AccessTypeRead, AccessTypeReadWrite}
