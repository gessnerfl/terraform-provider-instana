package restapi

import (
	"errors"
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/utils"
)

//InstanaPermission data type representing an Instana permission string
type InstanaPermission string

const (
	//PermissionCanConfigureApplications const for Instana permission CAN_CONFIGURE_APPLICATIONS
	PermissionCanConfigureApplications = InstanaPermission("CAN_CONFIGURE_APPLICATIONS")
	//PermissionCanSeeOnPremLiceneInformation const for Instana permission CAN_SEE_ON_PREM_LICENE_INFORMATION
	PermissionCanSeeOnPremLiceneInformation = InstanaPermission("CAN_SEE_ON_PREM_LICENE_INFORMATION")
	//PermissionCanConfigureEumApplications const for Instana permission CAN_CONFIGURE_EUM_APPLICATIONS
	PermissionCanConfigureEumApplications = InstanaPermission("CAN_CONFIGURE_EUM_APPLICATIONS")
	//PermissionCanConfigureAgents const for Instana permission CAN_CONFIGURE_AGENTS
	PermissionCanConfigureAgents = InstanaPermission("CAN_CONFIGURE_AGENTS")
	//PermissionCanViewTraceDetails const for Instana permission CAN_VIEW_TRACE_DETAILS
	PermissionCanViewTraceDetails = InstanaPermission("CAN_VIEW_TRACE_DETAILS")
	//PermissionCanViewLogs const for Instana permission CAN_VIEW_LOGS
	PermissionCanViewLogs = InstanaPermission("CAN_VIEW_LOGS")
	//PermissionCanConfigureSessionSettings const for Instana permission CAN_CONFIGURE_SESSION_SETTINGS
	PermissionCanConfigureSessionSettings = InstanaPermission("CAN_CONFIGURE_SESSION_SETTINGS")
	//PermissionCanConfigureIntegrations const for Instana permission CAN_CONFIGURE_INTEGRATIONS
	PermissionCanConfigureIntegrations = InstanaPermission("CAN_CONFIGURE_INTEGRATIONS")
	//PermissionCanConfigureGlobalAlertConfigs const for Instana permission CAN_CONFIGURE_GLOBAL_ALERT_CONFIGS
	PermissionCanConfigureGlobalAlertConfigs = InstanaPermission("CAN_CONFIGURE_GLOBAL_ALERT_CONFIGS")
	//PermissionCanConfigureGlobalAlertPayload const for Instana permission CAN_CONFIGURE_GLOBAL_ALERT_PAYLOAD
	PermissionCanConfigureGlobalAlertPayload = InstanaPermission("CAN_CONFIGURE_GLOBAL_ALERT_PAYLOAD")
	//PermissionCanConfigureMobileAppMonitoring const for Instana permission CAN_CONFIGURE_MOBILE_APP_MONITORING
	PermissionCanConfigureMobileAppMonitoring = InstanaPermission("CAN_CONFIGURE_MOBILE_APP_MONITORING")
	//PermissionCanConfigureAPITokens const for Instana permission CAN_CONFIGURE_API_TOKENS
	PermissionCanConfigureAPITokens = InstanaPermission("CAN_CONFIGURE_API_TOKENS")
	//PermissionCanConfigureServiceLevelIndicators const for Instana permission CAN_CONFIGURE_SERVICE_LEVEL_INDICATORS
	PermissionCanConfigureServiceLevelIndicators = InstanaPermission("CAN_CONFIGURE_SERVICE_LEVEL_INDICATORS")
	//PermissionCanConfigureAuthenticationMethods const for Instana permission CAN_CONFIGURE_AUTHENTICATION_METHODS
	PermissionCanConfigureAuthenticationMethods = InstanaPermission("CAN_CONFIGURE_AUTHENTICATION_METHODS")
	//PermissionCanConfigureReleases const for Instana permission CAN_CONFIGURE_RELEASES
	PermissionCanConfigureReleases = InstanaPermission("CAN_CONFIGURE_RELEASES")
	//PermissionCanViewAuditLog const for Instana permission CAN_VIEW_AUDIT_LOG
	PermissionCanViewAuditLog = InstanaPermission("CAN_VIEW_AUDIT_LOG")
	//PermissionCanConfigureCustomAlerts const for Instana permission CAN_CONFIGURE_CUSTOM_ALERTS
	PermissionCanConfigureCustomAlerts = InstanaPermission("CAN_CONFIGURE_CUSTOM_ALERTS")
	//PermissionCanConfigureAgentRunMode const for Instana permission CAN_CONFIGURE_AGENT_RUN_MODE
	PermissionCanConfigureAgentRunMode = InstanaPermission("CAN_CONFIGURE_AGENT_RUN_MODE")
	//PermissionCanConfigureServiceMapping const for Instana permission CAN_CONFIGURE_SERVICE_MAPPING
	PermissionCanConfigureServiceMapping = InstanaPermission("CAN_CONFIGURE_SERVICE_MAPPING")
	//PermissionCanSeeUsageInformation const for Instana permission CAN_SEE_USAGE_INFORMATION
	PermissionCanSeeUsageInformation = InstanaPermission("CAN_SEE_USAGE_INFORMATION")
	//PermissionCanEditAllAccessibleCustomDashboards const for Instana permission CAN_EDIT_ALL_ACCESSIBLE_CUSTOM_DASHBOARDS
	PermissionCanEditAllAccessibleCustomDashboards = InstanaPermission("CAN_EDIT_ALL_ACCESSIBLE_CUSTOM_DASHBOARDS")
	//PermissionCanConfigureUsers const for Instana permission CAN_CONFIGURE_USERS
	PermissionCanConfigureUsers = InstanaPermission("CAN_CONFIGURE_USERS")
	//PermissionCanInstallNewAgents const for Instana permission CAN_INSTALL_NEW_AGENTS
	PermissionCanInstallNewAgents = InstanaPermission("CAN_INSTALL_NEW_AGENTS")
	//PermissionCanConfigureTeams const for Instana permission CAN_CONFIGURE_TEAMS
	PermissionCanConfigureTeams = InstanaPermission("CAN_CONFIGURE_TEAMS")
	//PermissionCanCreatePublicCustomDashboards const for Instana permission CAN_CREATE_PUBLIC_CUSTOM_DASHBOARDS
	PermissionCanCreatePublicCustomDashboards = InstanaPermission("CAN_CREATE_PUBLIC_CUSTOM_DASHBOARDS")
	//PermissionCanConfigureLogManagement const for Instana permission CAN_CONFIGURE_LOG_MANAGEMENT
	PermissionCanConfigureLogManagement = InstanaPermission("CAN_CONFIGURE_LOG_MANAGEMENT")
	//PermissionCanViewAccountAndBillingInformation const for Instana permission CAN_VIEW_ACCOUNT_AND_BILLING_INFORMATION
	PermissionCanViewAccountAndBillingInformation = InstanaPermission("CAN_VIEW_ACCOUNT_AND_BILLING_INFORMATION")
)

//InstanaPermissions data type representing a slice of Instana permissions
type InstanaPermissions []InstanaPermission

//ToStringSlice converts the slice of InstanaPermissions to its string representation
func (permissions InstanaPermissions) ToStringSlice() []string {
	result := make([]string, len(permissions))
	for i, v := range permissions {
		result[i] = string(v)
	}
	return result
}

//IsSupported checks if the provided InstanaPermission is a supported Instana permission
func (permissions InstanaPermissions) IsSupported(toBeChecked InstanaPermission) bool {
	for _, p := range permissions {
		if p == toBeChecked {
			return true
		}
	}
	return false
}

//SupportedInstanaPermissions slice of all supported Permissions of the Instana API
var SupportedInstanaPermissions = InstanaPermissions{
	PermissionCanConfigureApplications,
	PermissionCanSeeOnPremLiceneInformation,
	PermissionCanConfigureEumApplications,
	PermissionCanConfigureAgents,
	PermissionCanViewTraceDetails,
	PermissionCanViewLogs,
	PermissionCanConfigureSessionSettings,
	PermissionCanConfigureIntegrations,
	PermissionCanConfigureGlobalAlertConfigs,
	PermissionCanConfigureGlobalAlertPayload,
	PermissionCanConfigureMobileAppMonitoring,
	PermissionCanConfigureAPITokens,
	PermissionCanConfigureServiceLevelIndicators,
	PermissionCanConfigureAuthenticationMethods,
	PermissionCanConfigureReleases,
	PermissionCanViewAuditLog,
	PermissionCanConfigureCustomAlerts,
	PermissionCanConfigureAgentRunMode,
	PermissionCanConfigureServiceMapping,
	PermissionCanSeeUsageInformation,
	PermissionCanEditAllAccessibleCustomDashboards,
	PermissionCanConfigureUsers,
	PermissionCanInstallNewAgents,
	PermissionCanConfigureTeams,
	PermissionCanCreatePublicCustomDashboards,
	PermissionCanConfigureLogManagement,
	PermissionCanViewAccountAndBillingInformation,
}

//GroupsResourcePath path to Group resource of Instana RESTful API
const GroupsResourcePath = RBACSettingsBasePath + "/alerts"

//ScopeBinding data structure for the Instana API model for scope bindings
type ScopeBinding struct {
	ScopeID     string  `json:"scopeId"`
	ScopeRoleID *string `json:"scopeRoleId"`
}

func (b *ScopeBinding) validate() error {
	if utils.IsBlank(b.ScopeID) {
		return errors.New("scopeId of scope binding is missing")
	}
	return nil
}

//APIPermissionSetWithRoles data structure for the Instana API model for permissions with roles
type APIPermissionSetWithRoles struct {
	ApplicationIDs          []ScopeBinding      `json:"applicationIds"`
	InfraDFQFilter          *ScopeBinding       `json:"infraDfqFilter"`
	KubernetesClusterUUIDs  []ScopeBinding      `json:"kubernetesClusterUUIDs"`
	KubernetesNamespaceUIDs []ScopeBinding      `json:"kubernetesNamespaceUIDs"`
	MobileAppIDs            []ScopeBinding      `json:"mobileAppIds"`
	WebsiteIDs              []ScopeBinding      `json:"websiteIds"`
	Permissions             []InstanaPermission `json:"permissions"`
}

func (m *APIPermissionSetWithRoles) validate() error {
	for _, a := range m.ApplicationIDs {
		if err := a.validate(); err != nil {
			return err
		}
	}
	if m.InfraDFQFilter != nil {
		if err := m.InfraDFQFilter.validate(); err != nil {
			return err
		}
	}
	for _, u := range m.KubernetesClusterUUIDs {
		if err := u.validate(); err != nil {
			return err
		}
	}
	for _, u := range m.KubernetesNamespaceUIDs {
		if err := u.validate(); err != nil {
			return err
		}
	}
	for _, m := range m.MobileAppIDs {
		if err := m.validate(); err != nil {
			return err
		}
	}
	for _, w := range m.WebsiteIDs {
		if err := w.validate(); err != nil {
			return err
		}
	}
	for _, p := range m.Permissions {
		if !SupportedInstanaPermissions.IsSupported(p) {
			return fmt.Errorf("%s is not a supported Instana permission", p)
		}
	}
	return nil
}

//APIMember data structure for the Instana API model for group members
type APIMember struct {
	UserID string  `json:"userId"`
	Email  *string `json:"email"`
}

func (m *APIMember) validate() error {
	if utils.IsBlank(m.UserID) {
		return errors.New("userId of group member is missing")
	}
	return nil
}

//Group data structure for the Instana API model for groups
type Group struct {
	ID            string                    `json:"id"`
	Name          string                    `json:"name"`
	Members       []APIMember               `json:"members"`
	PermissionSet APIPermissionSetWithRoles `json:"permissionSet"`
}

//GetIDForResourcePath implementation of the interface InstanaDataObject
func (c *Group) GetIDForResourcePath() string {
	return c.ID
}

//Validate implementation of the interface InstanaDataObject to verify if data object is correct
func (c *Group) Validate() error {
	if utils.IsBlank(c.ID) {
		return errors.New("id is missing")
	}
	if utils.IsBlank(c.Name) {
		return errors.New("name is missing")
	}
	for _, m := range c.Members {
		err := m.validate()
		if err != nil {
			return err
		}
	}
	return c.PermissionSet.validate()
}
