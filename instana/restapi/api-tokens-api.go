package restapi

import "errors"

//APITokensResourcePath path to API Tokens resource of Instana RESTful API
const APITokensResourcePath = SettingsBasePath + "/api-tokens"

//APIToken is the representation of a API Token in Instana
type APIToken struct {
	ID                                   string `json:"id"`
	AccessGrantingToken                  string `json:"accessGrantingToken"`
	InternalID                           string `json:"internalId"`
	Name                                 string `json:"name"`
	CanConfigureServiceMapping           bool   `json:"canConfigureServiceMapping"`
	CanConfigureEumApplications          bool   `json:"canConfigureEumApplications"`
	CanConfigureMobileAppMonitoring      bool   `json:"canConfigureMobileAppMonitoring"` //NEW
	CanConfigureUsers                    bool   `json:"canConfigureUsers"`
	CanInstallNewAgents                  bool   `json:"canInstallNewAgents"`
	CanSeeUsageInformation               bool   `json:"canSeeUsageInformation"`
	CanConfigureIntegrations             bool   `json:"canConfigureIntegrations"`
	CanSeeOnPremiseLicenseInformation    bool   `json:"canSeeOnPremLicenseInformation"`
	CanConfigureCustomAlerts             bool   `json:"canConfigureCustomAlerts"`
	CanConfigureAPITokens                bool   `json:"canConfigureApiTokens"`
	CanConfigureAgentRunMode             bool   `json:"canConfigureAgentRunMode"`
	CanViewAuditLog                      bool   `json:"canViewAuditLog"`
	CanConfigureAgents                   bool   `json:"canConfigureAgents"`
	CanConfigureAuthenticationMethods    bool   `json:"canConfigureAuthenticationMethods"`
	CanConfigureApplications             bool   `json:"canConfigureApplications"`
	CanConfigureTeams                    bool   `json:"canConfigureTeams"`
	CanConfigureReleases                 bool   `json:"canConfigureReleases"`
	CanConfigureLogManagement            bool   `json:"canConfigureLogManagement"`
	CanCreatePublicCustomDashboards      bool   `json:"canCreatePublicCustomDashboards"`
	CanViewLogs                          bool   `json:"canViewLogs"`
	CanViewTraceDetails                  bool   `json:"canViewTraceDetails"`
	CanConfigureSessionSettings          bool   `json:"canConfigureSessionSettings"`
	CanConfigureServiceLevelIndicators   bool   `json:"canConfigureServiceLevelIndicators"`
	CanConfigureGlobalAlertPayload       bool   `json:"canConfigureGlobalAlertPayload"`
	CanConfigureGlobalAlertConfigs       bool   `json:"canConfigureGlobalAlertConfigs"`
	CanViewAccountAndBillingInformation  bool   `json:"canViewAccountAndBillingInformation"`
	CanEditAllAccessibleCustomDashboards bool   `json:"canEditAllAccessibleCustomDashboards"`
}

//GetID implemention of the interface InstanaDataObject
func (r *APIToken) GetID() string {
	return r.ID
}

//Validate implementation of the interface InstanaDataObject to verify if data object is correct
func (r *APIToken) Validate() error {
	if len(r.AccessGrantingToken) == 0 {
		return errors.New("Access Granting Token is missing")
	}
	if len(r.Name) == 0 {
		return errors.New("Name is missing")
	}
	return nil
}
