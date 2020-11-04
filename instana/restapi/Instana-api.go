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
)

//InstanaAPI is the interface to all resources of the Instana Rest API
type InstanaAPI interface {
	CustomEventSpecifications() RestResource
	UserRoles() RestResource
	ApplicationConfigs() RestResource
	ApplicationAlertConfigs() RestResource
	AlertingChannels() RestResource
	AlertingConfigurations() RestResource
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
	return NewRestResource(CustomEventSpecificationResourcePath, NewCustomEventSpecificationUnmarshaller(), api.client)
}

//UserRoles implementation of InstanaAPI interface
func (api *baseInstanaAPI) UserRoles() RestResource {
	return NewRestResource(UserRolesResourcePath, NewUserRoleUnmarshaller(), api.client)
}

//ApplicationConfigs implementation of InstanaAPI interface
func (api *baseInstanaAPI) ApplicationConfigs() RestResource {
	return NewRestResource(ApplicationConfigsResourcePath, NewApplicationConfigUnmarshaller(), api.client)
}

//ApplicationConfigs implementation of InstanaAPI interface
func (api *baseInstanaAPI) ApplicationAlertConfigs() RestResource {
	return NewPostingRestResource(ApplicationAlertConfigsResourcePath, NewApplicationAlertConfigsUnmarshaller(), api.client)
}


//AlertingChannels implementation of InstanaAPI interface
func (api *baseInstanaAPI) AlertingChannels() RestResource {
	return NewRestResource(AlertingChannelsResourcePath, NewAlertingChannelUnmarshaller(), api.client)
}

//AlertingConfigurations implementation of InstanaAPI interface
func (api *baseInstanaAPI) AlertingConfigurations() RestResource {
	return NewRestResource(AlertsResourcePath, NewAlertingConfigurationUnmarshaller(), api.client)
}
