package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//ResourceInstanaUserRole the name of the terraform-provider-instana resource to manage user roles
const ResourceInstanaUserRole = "instana_user_role"

const (
	//UserRoleFieldName constant value for the schema field name
	UserRoleFieldName = "name"
	//UserRoleFieldCanConfigureServiceMapping constant value for the schema field can_configure_service_mapping
	UserRoleFieldCanConfigureServiceMapping = "can_configure_service_mapping"
	//UserRoleFieldCanConfigureEumApplications constant value for the schema field can_configure_eum_applications
	UserRoleFieldCanConfigureEumApplications = "can_configure_eum_applications"
	//UserRoleFieldCanConfigureMobileAppMonitoring constant value for the schema field can_configure_mobile_app_monitoring
	UserRoleFieldCanConfigureMobileAppMonitoring = "can_configure_mobile_app_monitoring"
	//UserRoleFieldCanConfigureUsers constant value for the schema field can_configure_users
	UserRoleFieldCanConfigureUsers = "can_configure_users"
	//UserRoleFieldCanInstallNewAgents constant value for the schema field can_install_new_agents
	UserRoleFieldCanInstallNewAgents = "can_install_new_agents"
	//UserRoleFieldCanSeeUsageInformation constant value for the schema field can_see_usage_information
	UserRoleFieldCanSeeUsageInformation = "can_see_usage_information"
	//UserRoleFieldCanConfigureIntegrations constant value for the schema field can_configure_integrations
	UserRoleFieldCanConfigureIntegrations = "can_configure_integrations"
	//UserRoleFieldCanSeeOnPremiseLicenseInformation constant value for the schema field can_see_on_premise_license_information
	UserRoleFieldCanSeeOnPremiseLicenseInformation = "can_see_on_premise_license_information"
	//UserRoleFieldCanConfigureRoles constant value for the schema field can_configure_roles
	UserRoleFieldCanConfigureRoles = "can_configure_roles"
	//UserRoleFieldCanConfigureCustomAlerts constant value for the schema field can_configure_custom_alerts
	UserRoleFieldCanConfigureCustomAlerts = "can_configure_custom_alerts"
	//UserRoleFieldCanConfigureAPITokens constant value for the schema field can_configure_api_tokens
	UserRoleFieldCanConfigureAPITokens = "can_configure_api_tokens"
	//UserRoleFieldCanConfigureAgentRunMode constant value for the schema field can_configure_agent_run_mode
	UserRoleFieldCanConfigureAgentRunMode = "can_configure_agent_run_mode"
	//UserRoleFieldCanViewAuditLog constant value for the schema field can_view_audit_log
	UserRoleFieldCanViewAuditLog = "can_view_audit_log"
	//UserRoleFieldCanConfigureObjectives constant value for the schema field can_configure_objectives
	UserRoleFieldCanConfigureObjectives = "can_configure_objectives"
	//UserRoleFieldCanConfigureAgents constant value for the schema field can_configure_agents
	UserRoleFieldCanConfigureAgents = "can_configure_agents"
	//UserRoleFieldCanConfigureAuthenticationMethods constant value for the schema field can_configure_authentication_methods
	UserRoleFieldCanConfigureAuthenticationMethods = "can_configure_authentication_methods"
	//UserRoleFieldCanConfigureApplications constant value for the schema field can_configure_applications
	UserRoleFieldCanConfigureApplications = "can_configure_applications"
	//UserRoleFieldCanConfigureTeams constant value for the schema field can_configure_teams
	UserRoleFieldCanConfigureTeams = "can_configure_teams"
	//UserRoleFieldRestrictedAccess constant value for the schema field restricted_access
	UserRoleFieldRestrictedAccess = "restricted_access"
	//UserRoleFieldCanConfigureReleases constant value for the schema field can_configure_releases
	UserRoleFieldCanConfigureReleases = "can_configure_releases"
	//UserRoleFieldCanConfigureLogManagement constant value for the schema field can_configure_log_management
	UserRoleFieldCanConfigureLogManagement = "can_configure_log_management"
	//UserRoleFieldCanCreatePublicCustomDashboards constant value for the schema field can_create_public_custom_dashboards
	UserRoleFieldCanCreatePublicCustomDashboards = "can_create_public_custom_dashboards"
	//UserRoleFieldCanViewLogs constant value for the schema field can_view_logs
	UserRoleFieldCanViewLogs = "can_view_logs"
	//UserRoleFieldCanViewTraceDetails constant value for the schema field can_view_trace_details
	UserRoleFieldCanViewTraceDetails = "can_view_trace_details"
)

var (
	userRoleSchemaName = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the user role",
	}
	userRoleSchemaCanConfigureServiceMapping = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if users of the role are allowed to configure service mappings",
	}
	userRoleSchemaCanConfigureEumApplications = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if users of the role are allowed to configure End User Monitoring applications",
	}
	userRoleSchemaCanConfigureMobileAppMonitoring = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if users of the role are allowed to configure Mobile App Monitoring",
	}
	userRoleSchemaCanConfigureUsers = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if users of the role are allowed to configure users",
	}
	userRoleSchemaCanInstallNewAgents = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if users of the role are allowed to install new agents",
	}
	userRoleSchemaCanSeeUsageInformation = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if users of the role are allowed to see usage information",
	}
	userRoleSchemaCanConfigureIntegrations = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if users of the role are allowed to configure integrations",
	}
	userRoleSchemaCanSeeOnPremiseLicenseInformation = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if users of the role are allowed to see onPremise license information",
	}
	userRoleSchemaCanConfigureRoles = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if users of the role are allowed to configure user roles",
	}
	userRoleSchemaCanConfigureCustomAlerts = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if users of the role are allowed to configure custom alerts",
	}
	userRoleSchemaCanConfigureAPITokens = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if users of the role are allowed to configure API tokens",
	}
	userRoleSchemaCanConfigureAgentRunMode = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if users of the role are allowed to configure agent run mode",
	}
	userRoleSchemaCanViewAuditLog = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if users of the role are allowed to view the audit log",
	}
	userRoleSchemaCanConfigureObjectives = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if users of the role are allowed to configure objectives",
	}
	userRoleSchemaCanConfigureAgents = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if users of the role are allowed to configure agents",
	}
	userRoleSchemaCanConfigureAuthenticationMethods = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if users of the role are allowed to configure authentication methods",
	}
	userRoleSchemaCanConfigureApplications = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if users of the role are allowed to configure applications",
	}

	userRoleFieldCanConfigureTeams = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if users of the role are allowed to configure teams (Groups)",
	}
	userRoleFieldRestrictedAccess = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if users of the role has limited access by group access scopes",
	}
	userRoleFieldCanConfigureReleases = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if users of the role are allowed to configure releases",
	}
	userRoleFieldCanConfigureLogManagement = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if users of the role are allowed to configure log management",
	}
	userRoleFieldCanCreatePublicCustomDashboards = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if users of the role are allowed to create public custom dashboards",
	}
	userRoleFieldCanViewLogs = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if users of the role are allowed to view logs",
	}
	userRoleFieldCanViewTraceDetails = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if users of the role are allowed to view trace details",
	}
)

//NewUserRoleResourceHandle creates a ResourceHandle instance for the terraform resource user role
func NewUserRoleResourceHandle() ResourceHandle {
	return &useerRoleResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaUserRole,
			Schema: map[string]*schema.Schema{
				UserRoleFieldName:                              userRoleSchemaName,
				UserRoleFieldCanConfigureServiceMapping:        userRoleSchemaCanConfigureServiceMapping,
				UserRoleFieldCanConfigureEumApplications:       userRoleSchemaCanConfigureEumApplications,
				UserRoleFieldCanConfigureMobileAppMonitoring:   userRoleSchemaCanConfigureMobileAppMonitoring,
				UserRoleFieldCanConfigureUsers:                 userRoleSchemaCanConfigureUsers,
				UserRoleFieldCanInstallNewAgents:               userRoleSchemaCanInstallNewAgents,
				UserRoleFieldCanSeeUsageInformation:            userRoleSchemaCanSeeUsageInformation,
				UserRoleFieldCanConfigureIntegrations:          userRoleSchemaCanConfigureIntegrations,
				UserRoleFieldCanSeeOnPremiseLicenseInformation: userRoleSchemaCanSeeOnPremiseLicenseInformation,
				UserRoleFieldCanConfigureRoles:                 userRoleSchemaCanConfigureRoles,
				UserRoleFieldCanConfigureCustomAlerts:          userRoleSchemaCanConfigureCustomAlerts,
				UserRoleFieldCanConfigureAPITokens:             userRoleSchemaCanConfigureAPITokens,
				UserRoleFieldCanConfigureAgentRunMode:          userRoleSchemaCanConfigureAgentRunMode,
				UserRoleFieldCanViewAuditLog:                   userRoleSchemaCanViewAuditLog,
				UserRoleFieldCanConfigureObjectives:            userRoleSchemaCanConfigureObjectives,
				UserRoleFieldCanConfigureAgents:                userRoleSchemaCanConfigureAgents,
				UserRoleFieldCanConfigureAuthenticationMethods: userRoleSchemaCanConfigureAuthenticationMethods,
				UserRoleFieldCanConfigureApplications:          userRoleSchemaCanConfigureApplications,
				UserRoleFieldCanConfigureTeams:                 userRoleFieldCanConfigureTeams,
				UserRoleFieldRestrictedAccess:                  userRoleFieldRestrictedAccess,
				UserRoleFieldCanConfigureReleases:              userRoleFieldCanConfigureReleases,
				UserRoleFieldCanConfigureLogManagement:         userRoleFieldCanConfigureLogManagement,
				UserRoleFieldCanCreatePublicCustomDashboards:   userRoleFieldCanCreatePublicCustomDashboards,
				UserRoleFieldCanViewLogs:                       userRoleFieldCanViewLogs,
				UserRoleFieldCanViewTraceDetails:               userRoleFieldCanViewTraceDetails,
			},
			SchemaVersion: 1,
		},
	}
}

type useerRoleResource struct {
	metaData ResourceMetaData
}

func (r *useerRoleResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *useerRoleResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{
		{
			Type:    r.schemaV0().CoreConfigSchema().ImpliedType(),
			Upgrade: r.migrateVersion0ToVersion1,
			Version: 0,
		},
	}
}

func (r *useerRoleResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource {
	return api.UserRoles()
}

func (r *useerRoleResource) SetComputedFields(d *schema.ResourceData) {
	//No computed fields defined
}

func (r *useerRoleResource) UpdateState(d *schema.ResourceData, obj restapi.InstanaDataObject) error {
	userRole := obj.(*restapi.UserRole)
	d.Set(UserRoleFieldName, userRole.Name)
	d.Set(UserRoleFieldCanConfigureServiceMapping, userRole.CanConfigureServiceMapping)
	d.Set(UserRoleFieldCanConfigureEumApplications, userRole.CanConfigureEumApplications)
	d.Set(UserRoleFieldCanConfigureMobileAppMonitoring, userRole.CanConfigureMobileAppMonitoring)
	d.Set(UserRoleFieldCanConfigureUsers, userRole.CanConfigureUsers)
	d.Set(UserRoleFieldCanInstallNewAgents, userRole.CanInstallNewAgents)
	d.Set(UserRoleFieldCanSeeUsageInformation, userRole.CanSeeUsageInformation)
	d.Set(UserRoleFieldCanConfigureIntegrations, userRole.CanConfigureIntegrations)
	d.Set(UserRoleFieldCanSeeOnPremiseLicenseInformation, userRole.CanSeeOnPremiseLicenseInformation)
	d.Set(UserRoleFieldCanConfigureRoles, userRole.CanConfigureRoles)
	d.Set(UserRoleFieldCanConfigureCustomAlerts, userRole.CanConfigureCustomAlerts)
	d.Set(UserRoleFieldCanConfigureAPITokens, userRole.CanConfigureAPITokens)
	d.Set(UserRoleFieldCanConfigureAgentRunMode, userRole.CanConfigureAgentRunMode)
	d.Set(UserRoleFieldCanViewAuditLog, userRole.CanViewAuditLog)
	d.Set(UserRoleFieldCanConfigureObjectives, userRole.CanConfigureObjectives)
	d.Set(UserRoleFieldCanConfigureAgents, userRole.CanConfigureAgents)
	d.Set(UserRoleFieldCanConfigureAuthenticationMethods, userRole.CanConfigureAuthenticationMethods)
	d.Set(UserRoleFieldCanConfigureApplications, userRole.CanConfigureApplications)
	d.Set(UserRoleFieldCanConfigureTeams, userRole.CanConfigureTeams)
	d.Set(UserRoleFieldRestrictedAccess, userRole.RestrictedAccess)
	d.Set(UserRoleFieldCanConfigureReleases, userRole.CanConfigureReleases)
	d.Set(UserRoleFieldCanConfigureLogManagement, userRole.CanConfigureLogManagement)
	d.Set(UserRoleFieldCanCreatePublicCustomDashboards, userRole.CanCreatePublicCustomDashboards)
	d.Set(UserRoleFieldCanViewLogs, userRole.CanViewLogs)
	d.Set(UserRoleFieldCanViewTraceDetails, userRole.CanViewTraceDetails)

	d.SetId(userRole.ID)
	return nil
}

func (r *useerRoleResource) MapStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	return &restapi.UserRole{
		ID:                                d.Id(),
		Name:                              d.Get(UserRoleFieldName).(string),
		CanConfigureServiceMapping:        d.Get(UserRoleFieldCanConfigureServiceMapping).(bool),
		CanConfigureEumApplications:       d.Get(UserRoleFieldCanConfigureEumApplications).(bool),
		CanConfigureMobileAppMonitoring:   d.Get(UserRoleFieldCanConfigureMobileAppMonitoring).(bool),
		CanConfigureUsers:                 d.Get(UserRoleFieldCanConfigureUsers).(bool),
		CanInstallNewAgents:               d.Get(UserRoleFieldCanInstallNewAgents).(bool),
		CanSeeUsageInformation:            d.Get(UserRoleFieldCanSeeUsageInformation).(bool),
		CanConfigureIntegrations:          d.Get(UserRoleFieldCanConfigureIntegrations).(bool),
		CanSeeOnPremiseLicenseInformation: d.Get(UserRoleFieldCanSeeOnPremiseLicenseInformation).(bool),
		CanConfigureRoles:                 d.Get(UserRoleFieldCanConfigureRoles).(bool),
		CanConfigureCustomAlerts:          d.Get(UserRoleFieldCanConfigureCustomAlerts).(bool),
		CanConfigureAPITokens:             d.Get(UserRoleFieldCanConfigureAPITokens).(bool),
		CanConfigureAgentRunMode:          d.Get(UserRoleFieldCanConfigureAgentRunMode).(bool),
		CanViewAuditLog:                   d.Get(UserRoleFieldCanViewAuditLog).(bool),
		CanConfigureObjectives:            d.Get(UserRoleFieldCanConfigureObjectives).(bool),
		CanConfigureAgents:                d.Get(UserRoleFieldCanConfigureAgents).(bool),
		CanConfigureAuthenticationMethods: d.Get(UserRoleFieldCanConfigureAuthenticationMethods).(bool),
		CanConfigureApplications:          d.Get(UserRoleFieldCanConfigureApplications).(bool),
		CanConfigureTeams:                 d.Get(UserRoleFieldCanConfigureTeams).(bool),
		RestrictedAccess:                  d.Get(UserRoleFieldRestrictedAccess).(bool),
		CanConfigureReleases:              d.Get(UserRoleFieldCanConfigureReleases).(bool),
		CanConfigureLogManagement:         d.Get(UserRoleFieldCanConfigureLogManagement).(bool),
		CanCreatePublicCustomDashboards:   d.Get(UserRoleFieldCanCreatePublicCustomDashboards).(bool),
		CanViewLogs:                       d.Get(UserRoleFieldCanViewLogs).(bool),
		CanViewTraceDetails:               d.Get(UserRoleFieldCanViewTraceDetails).(bool),
	}, nil
}

func (r *useerRoleResource) schemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			UserRoleFieldName: userRoleSchemaName,
			"implicit_view_filter": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "The an implicit view filter which is applied for users of the given role",
			},
			UserRoleFieldCanConfigureServiceMapping:        userRoleSchemaCanConfigureServiceMapping,
			UserRoleFieldCanConfigureEumApplications:       userRoleSchemaCanConfigureEumApplications,
			UserRoleFieldCanConfigureUsers:                 userRoleSchemaCanConfigureUsers,
			UserRoleFieldCanInstallNewAgents:               userRoleSchemaCanInstallNewAgents,
			UserRoleFieldCanSeeUsageInformation:            userRoleSchemaCanSeeUsageInformation,
			UserRoleFieldCanConfigureIntegrations:          userRoleSchemaCanConfigureIntegrations,
			UserRoleFieldCanSeeOnPremiseLicenseInformation: userRoleSchemaCanSeeOnPremiseLicenseInformation,
			UserRoleFieldCanConfigureRoles:                 userRoleSchemaCanConfigureRoles,
			UserRoleFieldCanConfigureCustomAlerts:          userRoleSchemaCanConfigureCustomAlerts,
			UserRoleFieldCanConfigureAPITokens:             userRoleSchemaCanConfigureAPITokens,
			UserRoleFieldCanConfigureAgentRunMode:          userRoleSchemaCanConfigureAgentRunMode,
			UserRoleFieldCanViewAuditLog:                   userRoleSchemaCanViewAuditLog,
			UserRoleFieldCanConfigureObjectives:            userRoleSchemaCanConfigureObjectives,
			UserRoleFieldCanConfigureAgents:                userRoleSchemaCanConfigureAgents,
			UserRoleFieldCanConfigureAuthenticationMethods: userRoleSchemaCanConfigureAuthenticationMethods,
			UserRoleFieldCanConfigureApplications:          userRoleSchemaCanConfigureApplications,
		},
	}
}

func (r *useerRoleResource) migrateVersion0ToVersion1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	delete(rawState, "implicit_view_filter")
	return rawState, nil
}
