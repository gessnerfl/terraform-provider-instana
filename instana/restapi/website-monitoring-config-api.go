package restapi

// WebsiteMonitoringConfigResourcePath path to website monitoring config resource of Instana RESTful API
const WebsiteMonitoringConfigResourcePath = WebsiteMonitoringResourcePath + "/config"

// WebsiteMonitoringConfig data structure of a Website Monitoring Configuration of the Instana API
type WebsiteMonitoringConfig struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	AppName string `json:"appName"`
}

// GetIDForResourcePath implemention of the interface InstanaDataObject
func (r *WebsiteMonitoringConfig) GetIDForResourcePath() string {
	return r.ID
}
