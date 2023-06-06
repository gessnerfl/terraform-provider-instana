package restapi

import "encoding/json"

type MonitorConfig struct{}

type SyntheticTest struct {
	ID               string            `json:"id"`
	Label            string            `json:"label"`
	Description      string            `json:"description"`
	Active           bool              `json:"active"`
	ApplicationID    string            `json:"applicationId"`
	Configuration    json.RawMessage   `json:"configuration"`
	CustomProperties map[string]string `json:"customProperties"`
	Locations        []string          `json:"locations"`
	PlaybackMode     string            `json:"playbackMode"`
	TestFrequency    int32             `json:"testFrequency"`
}

// GetIDForResourcePath implementation of the interface InstanaDataObject for SyntheticTest
func (a *SyntheticTest) GetIDForResourcePath() string {
	return a.ID
}

// Validate implementation of the interface InstanaDataObject for SyntheticTest
func (a *SyntheticTest) Validate() error {
	//No validation required validation part of terraform schema
	return nil
}
