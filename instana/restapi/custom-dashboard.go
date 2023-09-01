package restapi

import "encoding/json"

// CustomDashboardsResourcePath the API resource path for Custom Dashboards
const CustomDashboardsResourcePath = InstanaAPIBasePath + "/custom-dashboard"

type CustomDashboard struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	AccessRules []AccessRule    `json:"accessRules"`
	Widgets     json.RawMessage `json:"widgets"`
}

// GetIDForResourcePath implementation of the interface InstanaDataObject for CustomDashboard
func (a *CustomDashboard) GetIDForResourcePath() string {
	return a.ID
}
