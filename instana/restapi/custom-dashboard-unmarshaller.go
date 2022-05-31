package restapi

import (
	"encoding/json"
)

//NewCustomDashboardUnmarshaller creates a new Unmarshaller instance for CustomDashboard
func NewCustomDashboardUnmarshaller() JSONUnmarshaller {
	return &customDashboardUnmarshaller{}
}

type customDashboardUnmarshaller struct{}

//Unmarshal Unmarshaller interface implementation
func (u *customDashboardUnmarshaller) Unmarshal(data []byte) (interface{}, error) {
	temp := &tempCustomDashboard{}
	if err := json.Unmarshal(data, temp); err != nil {
		return &CustomDashboard{}, err
	}

	var widgets string
	if temp.Widgets != nil {
		widgets = string(temp.Widgets)
	}
	return &CustomDashboard{
		ID:          temp.ID,
		Title:       temp.Title,
		AccessRules: temp.AccessRules,
		Widgets:     widgets,
	}, nil
}

type tempCustomDashboard struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	AccessRules []AccessRule    `json:"accessRules"`
	Widgets     json.RawMessage `json:"widgets"`
}
