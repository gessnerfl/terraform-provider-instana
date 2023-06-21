package restapi

import (
	"errors"

	"github.com/gessnerfl/terraform-provider-instana/utils"
)

type SyntheticTestConfig struct {
	MarkSyntheticCall bool   `json:"markSyntheticCall"`
	Retries           int32  `json:"retries"`
	RetryInterval     int32  `json:"retryInterval"`
	SyntheticType     string `json:"syntheticType"`
	Timeout           string `json:"timeout"`
	URL               string `json:"url"`
	Operation         string `json:"operation"`
	Script            string `json:"script"`
}

type SyntheticTest struct {
	ID               string                 `json:"id"`
	Label            string                 `json:"label"`
	Description      string                 `json:"description"`
	Active           bool                   `json:"active"`
	Configuration    SyntheticTestConfig    `json:"configuration"`
	CustomProperties map[string]interface{} `json:"customProperties"`
	Locations        []string               `json:"locations"`
	PlaybackMode     string                 `json:"playbackMode"`
	TestFrequency    int32                  `json:"testFrequency"`
}

// GetIDForResourcePath implementation of the interface InstanaDataObject for SyntheticTest
func (s *SyntheticTest) GetIDForResourcePath() string {
	return s.ID
}

// Validate implementation of the interface InstanaDataObject for SyntheticTest
func (s *SyntheticTest) Validate() error {
	if utils.IsBlank(s.ID) {
		return errors.New("id is missing")
	}
	if utils.IsBlank(s.Label) {
		return errors.New("label is missing")
	}

	if len(s.Locations) == 0 {
		return errors.New("locations cannot be empty")
	}

	switch s.Configuration.SyntheticType {
	case "HTTPAction":
		if utils.IsBlank(s.Configuration.URL) {
			return errors.New("configuration.url is missing")
		}
		if utils.IsBlank(s.Configuration.Operation) {
			return errors.New("configuration.operation is missing")
		}
	case "HTTPScript":
		if utils.IsBlank(s.Configuration.Script) {
			return errors.New("configuration.script is missing")
		}
	default:
		return errors.New("incorrect synthetic_type")
	}
	return nil
}
