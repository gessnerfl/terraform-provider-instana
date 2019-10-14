package instana

import (
	"errors"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	//UserRoleFieldName constant value for the schema field name
	UserRoleFieldName = "name"
	//UserRoleFieldImplicitViewFilter constant value for the schema field implicit_view_filter
	UserRoleFieldImplicitViewFilter = "implicit_view_filter"
	//UserRoleFieldCanConfigureServiceMapping constant value for the schema field can_configure_service_mapping
	UserRoleFieldCanConfigureServiceMapping = "can_configure_service_mapping"
	//UserRoleFieldCanConfigureEumApplications constant value for the schema field can_configure_eum_applications
	UserRoleFieldCanConfigureEumApplications = "can_configure_eum_applications"
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
)

//CreateResourceUserRole creates the resource definition for the resource instana_user_role
func CreateResourceUserRole() *schema.Resource {
	return &schema.Resource{
		Create: CreateUserRole,
		Read:   ReadUserRole,
		Update: UpdateUserRole,
		Delete: DeleteUserRole,

		Schema: map[string]*schema.Schema{
			UserRoleFieldName: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the user role",
			},
			UserRoleFieldImplicitViewFilter: &schema.Schema{
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "The an implicit view filter which is applied for users of the given role",
			},
			UserRoleFieldCanConfigureServiceMapping: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Configures if users of the role are allowed to configure service mappings",
			},
			UserRoleFieldCanConfigureEumApplications: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Configures if users of the role are allowed to configure End User Monitoring applications",
			},
			UserRoleFieldCanConfigureUsers: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Configures if users of the role are allowed to configure users",
			},
			UserRoleFieldCanInstallNewAgents: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Configures if users of the role are allowed to install new agents",
			},
			UserRoleFieldCanSeeUsageInformation: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Configures if users of the role are allowed to see usage information",
			},
			UserRoleFieldCanConfigureIntegrations: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Configures if users of the role are allowed to configure integrations",
			},
			UserRoleFieldCanSeeOnPremiseLicenseInformation: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Configures if users of the role are allowed to see onPremise license information",
			},
			UserRoleFieldCanConfigureRoles: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Configures if users of the role are allowed to configure user roles",
			},
			UserRoleFieldCanConfigureCustomAlerts: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Configures if users of the role are allowed to configure custom alerts",
			},
			UserRoleFieldCanConfigureAPITokens: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Configures if users of the role are allowed to configure API tokens",
			},
			UserRoleFieldCanConfigureAgentRunMode: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Configures if users of the role are allowed to configure agent run mode",
			},
			UserRoleFieldCanViewAuditLog: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Configures if users of the role are allowed to view the audit log",
			},
			UserRoleFieldCanConfigureObjectives: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Configures if users of the role are allowed to configure objectives",
			},
			UserRoleFieldCanConfigureAgents: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Configures if users of the role are allowed to configure agents",
			},
			UserRoleFieldCanConfigureAuthenticationMethods: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Configures if users of the role are allowed to configure authentication methods",
			},
			UserRoleFieldCanConfigureApplications: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Configures if users of the role are allowed to configure applications",
			},
		},
	}
}

//CreateUserRole defines the create operation for the resource instana_user_role
func CreateUserRole(d *schema.ResourceData, meta interface{}) error {
	d.SetId(RandomID())
	return UpdateUserRole(d, meta)
}

//ReadUserRole defines the read operation for the resource instana_user_role
func ReadUserRole(d *schema.ResourceData, meta interface{}) error {
	providerMeta := meta.(*ProviderMeta)
	instanaAPI := providerMeta.InstanaAPI
	ruleID := d.Id()
	if len(ruleID) == 0 {
		return errors.New("ID of user role is missing")
	}
	rule, err := instanaAPI.UserRoles().GetOne(ruleID)
	if err != nil {
		if err == restapi.ErrEntityNotFound {
			d.SetId("")
			return nil
		}
		return err
	}
	updateUserRoleState(d, rule)
	return nil
}

//UpdateUserRole defines the update operation for the resource instana_user_role
func UpdateUserRole(d *schema.ResourceData, meta interface{}) error {
	providerMeta := meta.(*ProviderMeta)
	instanaAPI := providerMeta.InstanaAPI
	rule := createUserRoleFromResourceData(d)
	updatedUserRole, err := instanaAPI.UserRoles().Upsert(rule)
	if err != nil {
		return err
	}
	updateUserRoleState(d, updatedUserRole)
	return nil
}

//DeleteUserRole defines the delete operation for the resource instana_user_role
func DeleteUserRole(d *schema.ResourceData, meta interface{}) error {
	providerMeta := meta.(*ProviderMeta)
	instanaAPI := providerMeta.InstanaAPI
	userRole := createUserRoleFromResourceData(d)
	err := instanaAPI.UserRoles().DeleteByID(userRole.ID)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func createUserRoleFromResourceData(d *schema.ResourceData) restapi.UserRole {
	return restapi.UserRole{
		ID:                                d.Id(),
		Name:                              d.Get(UserRoleFieldName).(string),
		ImplicitViewFilter:                d.Get(UserRoleFieldImplicitViewFilter).(string),
		CanConfigureServiceMapping:        d.Get(UserRoleFieldCanConfigureServiceMapping).(bool),
		CanConfigureEumApplications:       d.Get(UserRoleFieldCanConfigureEumApplications).(bool),
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
	}
}

func updateUserRoleState(d *schema.ResourceData, userRole restapi.UserRole) {
	d.Set(UserRoleFieldName, userRole.Name)
	d.Set(UserRoleFieldImplicitViewFilter, userRole.ImplicitViewFilter)
	d.Set(UserRoleFieldCanConfigureServiceMapping, userRole.CanConfigureServiceMapping)
	d.Set(UserRoleFieldCanConfigureEumApplications, userRole.CanConfigureEumApplications)
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

	d.SetId(userRole.ID)
}
