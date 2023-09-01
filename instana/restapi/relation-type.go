package restapi

// RelationType custom type for relation type
type RelationType string

// RelationTypes custom type for a slice of RelationType
type RelationTypes []RelationType

// ToStringSlice Returns the corresponding string representations
func (types RelationTypes) ToStringSlice() []string {
	result := make([]string, len(types))
	for i, v := range types {
		result[i] = string(v)
	}
	return result
}

const (
	//RelationTypeUser constant value for the USER RelationType
	RelationTypeUser = RelationType("USER")
	//RelationTypeApiToken constant value for the API_TOKEN RelationType
	RelationTypeApiToken = RelationType("API_TOKEN")
	//RelationTypeRole constant value for the ROLE RelationType
	RelationTypeRole = RelationType("ROLE")
	//RelationTypeTeam constant value for the TEAM RelationType
	RelationTypeTeam = RelationType("TEAM")
	//RelationTypeGlobal constant value for the GLOBAL RelationType
	RelationTypeGlobal = RelationType("GLOBAL")
)

// SupportedRelationTypes list of all supported RelationType
var SupportedRelationTypes = RelationTypes{RelationTypeUser, RelationTypeApiToken, RelationTypeRole, RelationTypeTeam, RelationTypeGlobal}
