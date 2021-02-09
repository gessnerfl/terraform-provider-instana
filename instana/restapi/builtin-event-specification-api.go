package restapi

//BuiltinEventSpecificationResourcePath path to Builtin Event Specification settings resource of Instana RESTful API
const BuiltinEventSpecificationResourcePath = EventSpecificationBasePath + "/built-in"

//BuiltinEventSpecification is the representation of a builtin event specification in Instana
type BuiltinEventSpecification struct {
	ID            string  `json:"id"`
	ShortPluginID string  `json:"shortPluginId"`
	Name          string  `json:"name"`
	Description   *string `json:"description"`
	Severity      int     `json:"severity"`
	Triggering    bool    `json:"triggering"`
	Enabled       bool    `json:"enabled"`
}

//GetID implemention of the interface InstanaDataObject
func (spec BuiltinEventSpecification) GetID() string {
	return spec.ID
}

//Validate implementation of the interface InstanaDataObject to verify if data object is correct. As this is read only datasource no validation is applied
func (spec BuiltinEventSpecification) Validate() error {
	return nil
}
