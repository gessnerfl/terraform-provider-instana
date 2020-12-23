package restapi

import "errors"

//WebsiteMonitoringConfigResourcePath path to website monitoring config resource of Instana RESTful API
const WebsiteMonitoringConfigResourcePath = WebsiteMonitoringResourcePath + "/config"

//WebsiteMonitoringConfig data structure of a Website Monitoring Configuration of the Instana API
type WebsiteMonitoringConfig struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	AppName string `json:"appName"`
}

//GetID implemention of the interface InstanaDataObject
func (r WebsiteMonitoringConfig) GetID() string {
	return r.ID
}

//Validate implementation of the interface InstanaDataObject to verify if data object is correct
func (r WebsiteMonitoringConfig) Validate() error {
	if len(r.Name) == 0 {
		return errors.New("Name is missing")
	}
	return nil
}
