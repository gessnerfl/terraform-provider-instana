package restapi

const (
	//InstanaAPIBasePath path to Instana RESTful API
	InstanaAPIBasePath = "/api"
	//EventsBasePath path to Events resource of Instana RESTful API
	EventsBasePath = InstanaAPIBasePath + "/events"
	//settingsPathElement path element to settings
	settingsPathElement = "/settings"
	//EventSettingsBasePath path to Event Settings resource of Instana RESTful API
	EventSettingsBasePath = EventsBasePath + settingsPathElement
	//SettingsBasePath path to Event Settings resource of Instana RESTful API
	SettingsBasePath = InstanaAPIBasePath + settingsPathElement
	//WebsiteMonitoringResourcePath path to website monitoring
	WebsiteMonitoringResourcePath = InstanaAPIBasePath + "/website-monitoring"
)

//InstanaAPI is the interface to all resources of the Instana Rest API
type InstanaAPI interface {
	CustomEventSpecifications() RestResource
	BuiltinEventSpecifications() ReadOnlyRestResource
	UserRoles() RestResource
	ApplicationConfigs() RestResource
	AlertingChannels() RestResource
	AlertingConfigurations() RestResource
	SliConfigs() RestResource
	WebsiteMonitoringConfig() RestResource
}

//NewInstanaAPI creates a new instance of the instana API
func NewInstanaAPI(apiToken string, endpoint string) InstanaAPI {
	client := NewClient(apiToken, endpoint)
	return &baseInstanaAPI{client: client}
}

type baseInstanaAPI struct {
	client RestClient
}

//CustomEventSpecifications implementation of InstanaAPI interface
func (api *baseInstanaAPI) CustomEventSpecifications() RestResource {
	return NewCreatePUTUpdatePUTRestResource(CustomEventSpecificationResourcePath, NewDefaultJSONUnmarshaller(&CustomEventSpecification{}), api.client)
}

//CustomEventSpecifications implementation of InstanaAPI interface
func (api *baseInstanaAPI) BuiltinEventSpecifications() ReadOnlyRestResource {
	return NewReadOnlyRestResource(BuiltinEventSpecificationResourcePath, NewDefaultJSONUnmarshaller(&BuiltinEventSpecification{}), NewDefaultJSONUnmarshaller(&[]BuiltinEventSpecification{}), api.client)
}

//UserRoles implementation of InstanaAPI interface
func (api *baseInstanaAPI) UserRoles() RestResource {
	return NewCreatePUTUpdatePUTRestResource(UserRolesResourcePath, NewDefaultJSONUnmarshaller(&UserRole{}), api.client)
}

//ApplicationConfigs implementation of InstanaAPI interface
func (api *baseInstanaAPI) ApplicationConfigs() RestResource {
	return NewCreatePUTUpdatePUTRestResource(ApplicationConfigsResourcePath, NewApplicationConfigUnmarshaller(), api.client)
}

//AlertingChannels implementation of InstanaAPI interface
func (api *baseInstanaAPI) AlertingChannels() RestResource {
	return NewCreatePUTUpdatePUTRestResource(AlertingChannelsResourcePath, NewDefaultJSONUnmarshaller(&AlertingChannel{}), api.client)
}

//AlertingConfigurations implementation of InstanaAPI interface
func (api *baseInstanaAPI) AlertingConfigurations() RestResource {
	return NewCreatePUTUpdatePUTRestResource(AlertsResourcePath, NewDefaultJSONUnmarshaller(&AlertingConfiguration{}), api.client)
}

func (api *baseInstanaAPI) SliConfigs() RestResource {
	return NewCreatePUTUpdatePUTRestResource(SliConfigResourcePath, NewDefaultJSONUnmarshaller(&SliConfig{}), api.client)
}

func (api *baseInstanaAPI) WebsiteMonitoringConfig() RestResource {
	return NewWebsiteMonitoringConfigRestResource(NewDefaultJSONUnmarshaller(&WebsiteMonitoringConfig{}), api.client)
}
