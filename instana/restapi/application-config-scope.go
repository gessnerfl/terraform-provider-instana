package restapi

// ApplicationConfigScope type definition of the application config scope of the Instana Web REST API
type ApplicationConfigScope string

// ApplicationConfigScopes type definition of slice of ApplicationConfigScope
type ApplicationConfigScopes []ApplicationConfigScope

// ToStringSlice returns a slice containing the string representations of the given scopes
func (scopes ApplicationConfigScopes) ToStringSlice() []string {
	result := make([]string, len(scopes))
	for i, s := range scopes {
		result[i] = string(s)
	}
	return result
}

const (
	//ApplicationConfigScopeIncludeNoDownstream constant for the scope INCLUDE_NO_DOWNSTREAM
	ApplicationConfigScopeIncludeNoDownstream = ApplicationConfigScope("INCLUDE_NO_DOWNSTREAM")
	//ApplicationConfigScopeIncludeImmediateDownstreamDatabaseAndMessaging constant for the scope INCLUDE_IMMEDIATE_DOWNSTREAM_DATABASE_AND_MESSAGING
	ApplicationConfigScopeIncludeImmediateDownstreamDatabaseAndMessaging = ApplicationConfigScope("INCLUDE_IMMEDIATE_DOWNSTREAM_DATABASE_AND_MESSAGING")
	//ApplicationConfigScopeIncludeAllDownstream constant for the scope INCLUDE_ALL_DOWNSTREAM
	ApplicationConfigScopeIncludeAllDownstream = ApplicationConfigScope("INCLUDE_ALL_DOWNSTREAM")
)

// SupportedApplicationConfigScopes supported ApplicationConfigScopes of the Instana Web REST API
var SupportedApplicationConfigScopes = ApplicationConfigScopes{
	ApplicationConfigScopeIncludeNoDownstream,
	ApplicationConfigScopeIncludeImmediateDownstreamDatabaseAndMessaging,
	ApplicationConfigScopeIncludeAllDownstream,
}
