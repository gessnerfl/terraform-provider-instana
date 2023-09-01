package restapi

// BuiltinEventSpecificationResourcePath path to Builtin Event Specification settings resource of Instana RESTful API
const BuiltinEventSpecificationResourcePath = EventSpecificationBasePath + "/built-in"

// BuiltinEventSpecification is the representation of a builtin event specification in Instana
type BuiltinEventSpecification struct {
	ID            string  `json:"id"`
	ShortPluginID string  `json:"shortPluginId"`
	Name          string  `json:"name"`
	Description   *string `json:"description"`
	Severity      int     `json:"severity"`
	Triggering    bool    `json:"triggering"`
	Enabled       bool    `json:"enabled"`
}

// GetIDForResourcePath implemention of the interface InstanaDataObject
func (spec *BuiltinEventSpecification) GetIDForResourcePath() string {
	return spec.ID
}
