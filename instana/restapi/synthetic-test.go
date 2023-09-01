package restapi

type SyntheticTestConfig struct {
	MarkSyntheticCall bool   `json:"markSyntheticCall"`
	Retries           int32  `json:"retries"`
	RetryInterval     int32  `json:"retryInterval"`
	SyntheticType     string `json:"syntheticType"`
	Timeout           string `json:"timeout"`
	// HttpAction
	URL              string                 `json:"url"`
	Operation        string                 `json:"operation"`
	Headers          map[string]interface{} `json:"headers"`
	Body             string                 `json:"body"`
	ValidationString string                 `json:"validationString"`
	FollowRedirect   bool                   `json:"followRedirect"`
	AllowInsecure    bool                   `json:"allowInsecure"`
	ExpectStatus     int32                  `json:"expectStatus"`
	ExpectMatch      string                 `json:"expectMatch"`
	// HttpScript
	Script string `json:"script"`
}

type SyntheticTest struct {
	ID               string                 `json:"id"`
	Label            string                 `json:"label"`
	Description      string                 `json:"description"`
	Active           bool                   `json:"active"`
	ApplicationID    *string                `json:"applicationId"`
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
