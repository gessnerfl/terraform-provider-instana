package restapi

type SyntheticLocation struct {
	ID           string `json:"id"`
	Label        string `json:"label"`
	Description  string `json:"description"`
	LocationType string `json:"locationType"`
}

// GetIDForResourcePath implementation of the interface InstanaDataObject for SyntheticLocation
func (s *SyntheticLocation) GetIDForResourcePath() string {
	return s.ID
}
