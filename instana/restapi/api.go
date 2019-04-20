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

//RulesResourcePath path to Rule resource of Instana RESTful API
const RulesResourcePath = EventSettingsBasePath + "/rules"

//RuleBindingsResourcePath path to Rule Binding resource of Instana RESTful API
const RuleBindingsResourcePath = EventSettingsBasePath + "/rule-bindings"

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

//InstanaAPI is the interface to all resources of the Instana Rest API
type InstanaAPI interface {
	Rules() RuleResource
	RuleBindings() RuleBindingResource
	UserRoles() UserRoleResource
	ApplicationConfigs() ApplicationConfigResource
}

//ErrEntityNotFound error message which is returned when the entity cannot be found at the server
var ErrEntityNotFound = errors.New("Failed to get resource from Instana API. 404 - Resource not found")
