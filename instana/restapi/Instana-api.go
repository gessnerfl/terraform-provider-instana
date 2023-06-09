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
	//RBACSettingsBasePath path to Role Based Access Control Settings resources of Instana RESTful API
	RBACSettingsBasePath = SettingsBasePath + "/rbac"
	//WebsiteMonitoringResourcePath path to website monitoring
	WebsiteMonitoringResourcePath = InstanaAPIBasePath + "/website-monitoring"
	//SyntheticSettingsBasePath path to synthetic monitoring
	SyntheticSettingsBasePath = InstanaAPIBasePath + "/synthetics" + settingsPathElement
	//SyntheticMonitoringTestPath path to synthetic monitoring tests
	SyntheticMonitorResourcePath = SyntheticSettingsBasePath + "/tests"
)

// InstanaAPI is the interface to all resources of the Instana Rest API
type InstanaAPI interface {
	CustomEventSpecifications() RestResource
	BuiltinEventSpecifications() ReadOnlyRestResource
	APITokens() RestResource
	ApplicationConfigs() RestResource
	ApplicationAlertConfigs() RestResource
	GlobalApplicationAlertConfigs() RestResource
	AlertingChannels() RestResource
	AlertingConfigurations() RestResource
	SliConfigs() RestResource
	WebsiteMonitoringConfig() RestResource
	WebsiteAlertConfig() RestResource
	Groups() RestResource
	CustomDashboards() RestResource
	SyntheticMonitorConfig() RestResource
}

// NewInstanaAPI creates a new instance of the instana API
func NewInstanaAPI(apiToken string, endpoint string, skipTlsVerification bool) InstanaAPI {
	client := NewClient(apiToken, endpoint, skipTlsVerification)
	return &baseInstanaAPI{client: client}
}

type baseInstanaAPI struct {
	client RestClient
}

// CustomEventSpecifications implementation of InstanaAPI interface
func (api *baseInstanaAPI) CustomEventSpecifications() RestResource {
	return NewCreatePUTUpdatePUTRestResource(CustomEventSpecificationResourcePath, NewDefaultJSONUnmarshaller(&CustomEventSpecification{}), api.client)
}

// BuiltinEventSpecifications implementation of InstanaAPI interface
func (api *baseInstanaAPI) BuiltinEventSpecifications() ReadOnlyRestResource {
	return NewReadOnlyRestResource(BuiltinEventSpecificationResourcePath, NewDefaultJSONUnmarshaller(&BuiltinEventSpecification{}), NewDefaultJSONUnmarshaller(&[]BuiltinEventSpecification{}), api.client)
}

// APITokens implementation of InstanaAPI interface
func (api *baseInstanaAPI) APITokens() RestResource {
	return NewCreatePOSTUpdatePUTRestResource(APITokensResourcePath, NewDefaultJSONUnmarshaller(&APIToken{}), api.client)
}

// ApplicationConfigs implementation of InstanaAPI interface
func (api *baseInstanaAPI) ApplicationConfigs() RestResource {
	return NewCreatePUTUpdatePUTRestResource(ApplicationConfigsResourcePath, NewApplicationConfigUnmarshaller(), api.client)
}

// ApplicationAlertConfigs implementation of InstanaAPI interface
func (api *baseInstanaAPI) ApplicationAlertConfigs() RestResource {
	return NewCreatePOSTUpdatePOSTRestResource(ApplicationAlertConfigsResourcePath, NewApplicationAlertConfigUnmarshaller(), api.client)
}

// GlobalApplicationAlertConfigs implementation of InstanaAPI interface
func (api *baseInstanaAPI) GlobalApplicationAlertConfigs() RestResource {
	return NewCreatePOSTUpdatePOSTRestResource(GlobalApplicationAlertConfigsResourcePath, NewApplicationAlertConfigUnmarshaller(), api.client)
}

// AlertingChannels implementation of InstanaAPI interface
func (api *baseInstanaAPI) AlertingChannels() RestResource {
	return NewCreatePUTUpdatePUTRestResource(AlertingChannelsResourcePath, NewDefaultJSONUnmarshaller(&AlertingChannel{}), api.client)
}

// AlertingConfigurations implementation of InstanaAPI interface
func (api *baseInstanaAPI) AlertingConfigurations() RestResource {
	return NewCreatePUTUpdatePUTRestResource(AlertsResourcePath, NewDefaultJSONUnmarshaller(&AlertingConfiguration{}), api.client)
}

func (api *baseInstanaAPI) SliConfigs() RestResource {
	return NewCreatePUTUpdatePUTRestResource(SliConfigResourcePath, NewDefaultJSONUnmarshaller(&SliConfig{}), api.client)
}

func (api *baseInstanaAPI) WebsiteMonitoringConfig() RestResource {
	return NewWebsiteMonitoringConfigRestResource(NewDefaultJSONUnmarshaller(&WebsiteMonitoringConfig{}), api.client)
}

func (api *baseInstanaAPI) WebsiteAlertConfig() RestResource {
	return NewCreatePOSTUpdatePOSTRestResource(WebsiteAlertConfigResourcePath, NewWebsiteAlertConfigUnmarshaller(), api.client)
}

func (api *baseInstanaAPI) Groups() RestResource {
	return NewCreatePOSTUpdatePUTRestResource(GroupsResourcePath, NewDefaultJSONUnmarshaller(&Group{}), api.client)
}

func (api *baseInstanaAPI) CustomDashboards() RestResource {
	return NewCreatePOSTUpdatePUTRestResource(CustomDashboardsResourcePath, NewDefaultJSONUnmarshaller(&CustomDashboard{}), api.client)
}

func (api *baseInstanaAPI) SyntheticMonitorConfig() RestResource {
	return NewSyntheticMonitorRestResource(NewDefaultJSONUnmarshaller(&SyntheticMonitor{}), api.client)
}
