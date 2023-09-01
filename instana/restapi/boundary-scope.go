package restapi

// BoundaryScope type definition of the application config boundary scope of the Instana Web REST API
type BoundaryScope string

// BoundaryScopes type definition of slice of BoundaryScopes
type BoundaryScopes []BoundaryScope

// ToStringSlice returns a slice containing the string representations of the given boundary scopes
func (scopes BoundaryScopes) ToStringSlice() []string {
	result := make([]string, len(scopes))
	for i, s := range scopes {
		result[i] = string(s)
	}
	return result
}

const (
	//BoundaryScopeAll constant value for the boundary scope ALL of an application config of the Instana Web REST API
	BoundaryScopeAll = BoundaryScope("ALL")
	//BoundaryScopeInbound constant value for the boundary scope INBOUND of an application config of the Instana Web REST API
	BoundaryScopeInbound = BoundaryScope("INBOUND")
	//BoundaryScopeDefault constant value for the boundary scope DEFAULT of an application config of the Instana Web REST API
	BoundaryScopeDefault = BoundaryScope("DEFAULT")
)

// SupportedApplicationConfigBoundaryScopes supported BoundaryScopes of the Instana Web REST API
var SupportedApplicationConfigBoundaryScopes = BoundaryScopes{BoundaryScopeAll, BoundaryScopeInbound, BoundaryScopeDefault}

// SupportedApplicationAlertConfigBoundaryScopes supported BoundaryScopes of the Instana Web REST API
var SupportedApplicationAlertConfigBoundaryScopes = BoundaryScopes{BoundaryScopeAll, BoundaryScopeInbound}
