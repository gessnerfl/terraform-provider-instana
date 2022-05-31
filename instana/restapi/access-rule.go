package restapi

type AccessRule struct {
	AccessType   AccessType   `json:"accessType"`
	RelatedID    *string      `json:"relatedId"`
	RelationType RelationType `json:"relationType"`
}
