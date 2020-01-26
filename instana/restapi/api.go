package restapi

import "errors"

//InstanaAPIBasePath path to Instana RESTful API
const InstanaAPIBasePath = "/api"

//EventsBasePath path to Events resource of Instana RESTful API
const EventsBasePath = InstanaAPIBasePath + "/events"

//settingsPathElement path element to settings
const settingsPathElement = "/settings"

//EventSettingsBasePath path to Event Settings resource of Instana RESTful API
const EventSettingsBasePath = EventsBasePath + settingsPathElement

//EventSpecificationBasePath path to Event Specification settings of Instana RESTful API
const EventSpecificationBasePath = EventSettingsBasePath + "/event-specifications"

//AlertingChannelsResourcePath path to Alerting channels resource of Instana RESTful API
const AlertingChannelsResourcePath = EventSettingsBasePath + "/alertingChannels"

//CustomEventSpecificationResourcePath path to Custom Event Specification settings resource of Instana RESTful API
const CustomEventSpecificationResourcePath = EventSpecificationBasePath + "/custom"

//SettingsBasePath path to Event Settings resource of Instana RESTful API
const SettingsBasePath = InstanaAPIBasePath + settingsPathElement

//UserRolesResourcePath path to User Role resource of Instana RESTful API
const UserRolesResourcePath = SettingsBasePath + "/roles"

//ApplicationMonitoringBasePath path to application monitoring resource of Instana RESTful API
const ApplicationMonitoringBasePath = InstanaAPIBasePath + "/application-monitoring"

//ApplicationMonitoringSettingsBasePath path to application monitoring settings resource of Instana RESTful API
const ApplicationMonitoringSettingsBasePath = ApplicationMonitoringBasePath + settingsPathElement

//ApplicationConfigsResourcePath path to application config resource of Instana RESTful API
const ApplicationConfigsResourcePath = ApplicationMonitoringSettingsBasePath + "/application"

//Severity representation of the severity in both worlds Instana API and Terraform Provider
type Severity struct {
	apiRepresentation       int
	terraformRepresentation string
}

//GetAPIRepresentation returns the integer representation of the Instana API
func (s Severity) GetAPIRepresentation() int { return s.apiRepresentation }

//GetTerraformRepresentation returns the string representation of the Terraform Provider
func (s Severity) GetTerraformRepresentation() string { return s.terraformRepresentation }

//SeverityCritical representation of the critical severity
var SeverityCritical = Severity{apiRepresentation: 10, terraformRepresentation: "critical"}

//SeverityWarning representation of the warning severity
var SeverityWarning = Severity{apiRepresentation: 5, terraformRepresentation: "warning"}

//InstanaDataObject is a marker interface for any data object provided by any resource of the Instana REST API
type InstanaDataObject interface {
	GetID() string
	Validate() error
}

//RestClient interface to access REST resources of the Instana API
type RestClient interface {
	GetOne(id string, resourcePath string) ([]byte, error)
	Put(data InstanaDataObject, resourcePath string) ([]byte, error)
	Delete(resourceID string, resourceBasePath string) error
}

//RestResource interface definition of a instana REST resource.
type RestResource interface {
	GetOne(id string) (InstanaDataObject, error)
	Upsert(data InstanaDataObject) (InstanaDataObject, error)
	Delete(data InstanaDataObject) error
	DeleteByID(id string) error
}

//InstanaAPI is the interface to all resources of the Instana Rest API
type InstanaAPI interface {
	CustomEventSpecifications() CustomEventSpecificationResource
	UserRoles() UserRoleResource
	ApplicationConfigs() ApplicationConfigResource
	AlertingChannels() RestResource
}

//ErrEntityNotFound error message which is returned when the entity cannot be found at the server
var ErrEntityNotFound = errors.New("Failed to get resource from Instana API. 404 - Resource not found")
