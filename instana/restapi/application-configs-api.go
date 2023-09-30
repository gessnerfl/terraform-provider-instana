package restapi

const (
	//ApplicationMonitoringBasePath path to application monitoring resource of Instana RESTful API
	ApplicationMonitoringBasePath = InstanaAPIBasePath + "/application-monitoring"
	//ApplicationMonitoringSettingsBasePath path to application monitoring settings resource of Instana RESTful API
	ApplicationMonitoringSettingsBasePath = ApplicationMonitoringBasePath + settingsPathElement
	//ApplicationConfigsResourcePath path to application config resource of Instana RESTful API
	ApplicationConfigsResourcePath = ApplicationMonitoringSettingsBasePath + "/application"
)

// ApplicationConfigResource represents the REST resource of application perspective configuration at Instana
type ApplicationConfigResource interface {
	GetOne(id string) (ApplicationConfig, error)
	Upsert(rule ApplicationConfig) (ApplicationConfig, error)
	Delete(rule ApplicationConfig) error
	DeleteByID(applicationID string) error
}

// ApplicationConfig is the representation of a application perspective configuration in Instana
type ApplicationConfig struct {
	ID                  string                 `json:"id"`
	Label               string                 `json:"label"`
	TagFilterExpression *TagFilter             `json:"tagFilterExpression"`
	Scope               ApplicationConfigScope `json:"scope"`
	BoundaryScope       BoundaryScope          `json:"boundaryScope"`
}

// GetIDForResourcePath implementation of the interface InstanaDataObject
func (a *ApplicationConfig) GetIDForResourcePath() string {
	return a.ID
}
