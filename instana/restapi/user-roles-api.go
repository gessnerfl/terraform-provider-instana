package restapi

import "errors"

//UserRolesResourcePath path to User Role resource of Instana RESTful API
const UserRolesResourcePath = SettingsBasePath + "/roles"

//UserRole is the representation of a user role in Instana
type UserRole struct {
	ID                                string `json:"id"`
	Name                              string `json:"name"`
	CanConfigureServiceMapping        bool   `json:"canConfigureServiceMapping"`
	CanConfigureEumApplications       bool   `json:"canConfigureEumApplications"`
	CanConfigureMobileAppMonitoring   bool   `json:"canConfigureMobileAppMonitoring"` //NEW
	CanConfigureUsers                 bool   `json:"canConfigureUsers"`
	CanInstallNewAgents               bool   `json:"canInstallNewAgents"`
	CanSeeUsageInformation            bool   `json:"canSeeUsageInformation"`
	CanConfigureIntegrations          bool   `json:"canConfigureIntegrations"`
	CanSeeOnPremiseLicenseInformation bool   `json:"canSeeOnPremLicenseInformation"`
	CanConfigureRoles                 bool   `json:"canConfigureRoles"`
	CanConfigureCustomAlerts          bool   `json:"canConfigureCustomAlerts"`
	CanConfigureAPITokens             bool   `json:"canConfigureApiTokens"`
	CanConfigureAgentRunMode          bool   `json:"canConfigureAgentRunMode"`
	CanViewAuditLog                   bool   `json:"canViewAuditLog"`
	CanConfigureObjectives            bool   `json:"canConfigureObjectives"`
	CanConfigureAgents                bool   `json:"canConfigureAgents"`
	CanConfigureAuthenticationMethods bool   `json:"canConfigureAuthenticationMethods"`
	CanConfigureApplications          bool   `json:"canConfigureApplications"`
	CanConfigureTeams                 bool   `json:"canConfigureTeams"`
	RestrictedAccess                  bool   `json:"restrictedAccess"`
	CanConfigureReleases              bool   `json:"canConfigureReleases"`
	CanConfigureLogManagement         bool   `json:"canConfigureLogManagement"`
	CanCreatePublicCustomDashboards   bool   `json:"canCreatePublicCustomDashboards"`
	CanViewLogs                       bool   `json:"canViewLogs"`
	CanViewTraceDetails               bool   `json:"canViewTraceDetails"`
}

//GetID implemention of the interface InstanaDataObject
func (r *UserRole) GetID() string {
	return r.ID
}

//Validate implementation of the interface InstanaDataObject to verify if data object is correct
func (r *UserRole) Validate() error {
	if len(r.ID) == 0 {
		return errors.New("ID is missing")
	}
	if len(r.Name) == 0 {
		return errors.New("Name is missing")
	}
	return nil
}
