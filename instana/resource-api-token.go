package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//ResourceInstanaAPIToken the name of the terraform-provider-instana resource to manage API tokens
const ResourceInstanaAPIToken = "instana_api_token"

const (
	//APITokenFieldAccessGrantingToken constant value for the schema field access_granting_token
	APITokenFieldAccessGrantingToken = "access_granting_token"
	//APITokenFieldInternalID constant value for the schema field internal_id
	APITokenFieldInternalID = "internal_id"
	//APITokenFieldName constant value for the schema field name
	APITokenFieldName = "name"
	//APITokenFieldFullName constant value for the schema field full_name
	APITokenFieldFullName = "full_name"
	//APITokenFieldCanConfigureServiceMapping constant value for the schema field can_configure_service_mapping
	APITokenFieldCanConfigureServiceMapping = "can_configure_service_mapping"
	//APITokenFieldCanConfigureEumApplications constant value for the schema field can_configure_eum_applications
	APITokenFieldCanConfigureEumApplications = "can_configure_eum_applications"
	//APITokenFieldCanConfigureMobileAppMonitoring constant value for the schema field can_configure_mobile_app_monitoring
	APITokenFieldCanConfigureMobileAppMonitoring = "can_configure_mobile_app_monitoring"
	//APITokenFieldCanConfigureUsers constant value for the schema field can_configure_users
	APITokenFieldCanConfigureUsers = "can_configure_users"
	//APITokenFieldCanInstallNewAgents constant value for the schema field can_install_new_agents
	APITokenFieldCanInstallNewAgents = "can_install_new_agents"
	//APITokenFieldCanSeeUsageInformation constant value for the schema field can_see_usage_information
	APITokenFieldCanSeeUsageInformation = "can_see_usage_information"
	//APITokenFieldCanConfigureIntegrations constant value for the schema field can_configure_integrations
	APITokenFieldCanConfigureIntegrations = "can_configure_integrations"
	//APITokenFieldCanSeeOnPremiseLicenseInformation constant value for the schema field can_see_on_premise_license_information
	APITokenFieldCanSeeOnPremiseLicenseInformation = "can_see_on_premise_license_information"
	//APITokenFieldCanConfigureCustomAlerts constant value for the schema field can_configure_custom_alerts
	APITokenFieldCanConfigureCustomAlerts = "can_configure_custom_alerts"
	//APITokenFieldCanConfigureAPITokens constant value for the schema field can_configure_api_tokens
	APITokenFieldCanConfigureAPITokens = "can_configure_api_tokens"
	//APITokenFieldCanConfigureAgentRunMode constant value for the schema field can_configure_agent_run_mode
	APITokenFieldCanConfigureAgentRunMode = "can_configure_agent_run_mode"
	//APITokenFieldCanViewAuditLog constant value for the schema field can_view_audit_log
	APITokenFieldCanViewAuditLog = "can_view_audit_log"
	//APITokenFieldCanConfigureAgents constant value for the schema field can_configure_agents
	APITokenFieldCanConfigureAgents = "can_configure_agents"
	//APITokenFieldCanConfigureAuthenticationMethods constant value for the schema field can_configure_authentication_methods
	APITokenFieldCanConfigureAuthenticationMethods = "can_configure_authentication_methods"
	//APITokenFieldCanConfigureApplications constant value for the schema field can_configure_applications
	APITokenFieldCanConfigureApplications = "can_configure_applications"
	//APITokenFieldCanConfigureTeams constant value for the schema field can_configure_teams
	APITokenFieldCanConfigureTeams = "can_configure_teams"
	//APITokenFieldCanConfigureReleases constant value for the schema field can_configure_releases
	APITokenFieldCanConfigureReleases = "can_configure_releases"
	//APITokenFieldCanConfigureLogManagement constant value for the schema field can_configure_log_management
	APITokenFieldCanConfigureLogManagement = "can_configure_log_management"
	//APITokenFieldCanCreatePublicCustomDashboards constant value for the schema field can_create_public_custom_dashboards
	APITokenFieldCanCreatePublicCustomDashboards = "can_create_public_custom_dashboards"
	//APITokenFieldCanViewLogs constant value for the schema field can_view_logs
	APITokenFieldCanViewLogs = "can_view_logs"
	//APITokenFieldCanViewTraceDetails constant value for the schema field can_view_trace_details
	APITokenFieldCanViewTraceDetails = "can_view_trace_details"
	//APITokenFieldCanConfigureSessionSettings constant value for the schema field can_configure_session_settings
	APITokenFieldCanConfigureSessionSettings = "can_configure_session_settings"
	//APITokenFieldCanConfigureServiceLevelIndicators constant value for the schema field can_configure_service_level_indicators
	APITokenFieldCanConfigureServiceLevelIndicators = "can_configure_service_level_indicators"
	//APITokenFieldCanConfigureGlobalAlertPayload constant value for the schema field can_configure_global_alert_payload
	APITokenFieldCanConfigureGlobalAlertPayload = "can_configure_global_alert_payload"
	//APITokenFieldCanConfigureGlobalAlertConfigs constant value for the schema field can_configure_global_alert_configs
	APITokenFieldCanConfigureGlobalAlertConfigs = "can_configure_global_alert_configs"
	//APITokenFieldCanViewAccountAndBillingInformation constant value for the schema field can_view_account_and_billing_information
	APITokenFieldCanViewAccountAndBillingInformation = "can_view_account_and_billing_information"
	//APITokenFieldCanEditAllAccessibleCustomDashboards constant value for the schema field can_edit_all_accessible_custom_dashboards
	APITokenFieldCanEditAllAccessibleCustomDashboards = "can_edit_all_accessible_custom_dashboards"
)

var (
	apiTokenSchemaAccessGrantingToken = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The token used for the api Client used in the Authorization header to authenticate the client",
	}
	apiTokenSchemaInternalID = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The internal ID of the access token from the Instana platform",
	}
	apiTokenSchemaName = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the API token",
	}
	apiTokenSchemaFullName = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The full name of the API token including prefix in suffix",
	}
	apiTokenSchemaCanConfigureServiceMapping = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure service mappings",
	}
	apiTokenSchemaCanConfigureEumApplications = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure End User Monitoring applications",
	}
	apiTokenSchemaCanConfigureMobileAppMonitoring = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure Mobile App Monitoring",
	}
	apiTokenSchemaCanConfigureUsers = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure users",
	}
	apiTokenSchemaCanInstallNewAgents = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to install new agents",
	}
	apiTokenSchemaCanSeeUsageInformation = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to see usage information",
	}
	apiTokenSchemaCanConfigureIntegrations = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure integrations",
	}
	apiTokenSchemaCanSeeOnPremiseLicenseInformation = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to see onPremise license information",
	}
	apiTokenSchemaCanConfigureCustomAlerts = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure custom alerts",
	}
	apiTokenSchemaCanConfigureAPITokens = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure API tokens",
	}
	apiTokenSchemaCanConfigureAgentRunMode = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure agent run mode",
	}
	apiTokenSchemaCanViewAuditLog = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to view the audit log",
	}
	apiTokenSchemaCanConfigureAgents = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure agents",
	}
	apiTokenSchemaCanConfigureAuthenticationMethods = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure authentication methods",
	}
	apiTokenSchemaCanConfigureApplications = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure applications",
	}
	apiTokenSchemaCanConfigureTeams = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure teams (Groups)",
	}
	apiTokenSchemaCanConfigureReleases = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure releases",
	}
	apiTokenSchemaCanConfigureLogManagement = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure log management",
	}
	apiTokenSchemaCanCreatePublicCustomDashboards = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to create public custom dashboards",
	}
	apiTokenSchemaCanViewLogs = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to view logs",
	}
	apiTokenSchemaCanViewTraceDetails = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to view trace details",
	}
	apiTokenSchemaCanConfigureSessionSettings = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure session settings",
	}
	apiTokenSchemaCanConfigureServiceLevelIndicators = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure service level indicators",
	}
	apiTokenSchemaCanConfigureGlobalAlertPayload = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure global alert payload",
	}
	apiTokenSchemaCanConfigureGlobalAlertConfigs = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to configure global alert configs",
	}
	apiTokenSchemaCanViewAccountAndBillingInformation = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to view account and billing information",
	}
	apiTokenSchemaCanEditAllAccessibleCustomDashboards = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Configures if the API token is allowed to edit all accessible custom dashboards",
	}
)

//NewAPITokenResourceHandle creates a ResourceHandle instance for the terraform resource API token
func NewAPITokenResourceHandle() ResourceHandle {
	internalIDFieldName := APITokenFieldInternalID
	return &apiTokenResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaAPIToken,
			Schema: map[string]*schema.Schema{
				APITokenFieldAccessGrantingToken:                  apiTokenSchemaAccessGrantingToken,
				APITokenFieldInternalID:                           apiTokenSchemaInternalID,
				APITokenFieldName:                                 apiTokenSchemaName,
				APITokenFieldFullName:                             apiTokenSchemaFullName,
				APITokenFieldCanConfigureServiceMapping:           apiTokenSchemaCanConfigureServiceMapping,
				APITokenFieldCanConfigureEumApplications:          apiTokenSchemaCanConfigureEumApplications,
				APITokenFieldCanConfigureMobileAppMonitoring:      apiTokenSchemaCanConfigureMobileAppMonitoring,
				APITokenFieldCanConfigureUsers:                    apiTokenSchemaCanConfigureUsers,
				APITokenFieldCanInstallNewAgents:                  apiTokenSchemaCanInstallNewAgents,
				APITokenFieldCanSeeUsageInformation:               apiTokenSchemaCanSeeUsageInformation,
				APITokenFieldCanConfigureIntegrations:             apiTokenSchemaCanConfigureIntegrations,
				APITokenFieldCanSeeOnPremiseLicenseInformation:    apiTokenSchemaCanSeeOnPremiseLicenseInformation,
				APITokenFieldCanConfigureCustomAlerts:             apiTokenSchemaCanConfigureCustomAlerts,
				APITokenFieldCanConfigureAPITokens:                apiTokenSchemaCanConfigureAPITokens,
				APITokenFieldCanConfigureAgentRunMode:             apiTokenSchemaCanConfigureAgentRunMode,
				APITokenFieldCanViewAuditLog:                      apiTokenSchemaCanViewAuditLog,
				APITokenFieldCanConfigureAgents:                   apiTokenSchemaCanConfigureAgents,
				APITokenFieldCanConfigureAuthenticationMethods:    apiTokenSchemaCanConfigureAuthenticationMethods,
				APITokenFieldCanConfigureApplications:             apiTokenSchemaCanConfigureApplications,
				APITokenFieldCanConfigureTeams:                    apiTokenSchemaCanConfigureTeams,
				APITokenFieldCanConfigureReleases:                 apiTokenSchemaCanConfigureReleases,
				APITokenFieldCanConfigureLogManagement:            apiTokenSchemaCanConfigureLogManagement,
				APITokenFieldCanCreatePublicCustomDashboards:      apiTokenSchemaCanCreatePublicCustomDashboards,
				APITokenFieldCanViewLogs:                          apiTokenSchemaCanViewLogs,
				APITokenFieldCanViewTraceDetails:                  apiTokenSchemaCanViewTraceDetails,
				APITokenFieldCanConfigureSessionSettings:          apiTokenSchemaCanConfigureSessionSettings,
				APITokenFieldCanConfigureServiceLevelIndicators:   apiTokenSchemaCanConfigureServiceLevelIndicators,
				APITokenFieldCanConfigureGlobalAlertPayload:       apiTokenSchemaCanConfigureGlobalAlertPayload,
				APITokenFieldCanConfigureGlobalAlertConfigs:       apiTokenSchemaCanConfigureGlobalAlertConfigs,
				APITokenFieldCanViewAccountAndBillingInformation:  apiTokenSchemaCanViewAccountAndBillingInformation,
				APITokenFieldCanEditAllAccessibleCustomDashboards: apiTokenSchemaCanEditAllAccessibleCustomDashboards,
			},
			SchemaVersion:    0,
			SkipIDGeneration: true,
			ResourceIDField:  &internalIDFieldName,
		},
	}
}

type apiTokenResource struct {
	metaData ResourceMetaData
}

func (r *apiTokenResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *apiTokenResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{}
}

func (r *apiTokenResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource {
	return api.APITokens()
}

func (r *apiTokenResource) SetComputedFields(d *schema.ResourceData) {
	d.Set(APITokenFieldInternalID, RandomID())
	d.Set(APITokenFieldAccessGrantingToken, RandomID())
}

func (r *apiTokenResource) UpdateState(d *schema.ResourceData, obj restapi.InstanaDataObject) error {
	apiToken := obj.(*restapi.APIToken)
	d.Set(APITokenFieldAccessGrantingToken, apiToken.AccessGrantingToken)
	d.Set(APITokenFieldInternalID, apiToken.InternalID)
	d.Set(APITokenFieldFullName, apiToken.Name)
	d.Set(APITokenFieldCanConfigureServiceMapping, apiToken.CanConfigureServiceMapping)
	d.Set(APITokenFieldCanConfigureEumApplications, apiToken.CanConfigureEumApplications)
	d.Set(APITokenFieldCanConfigureMobileAppMonitoring, apiToken.CanConfigureMobileAppMonitoring)
	d.Set(APITokenFieldCanConfigureUsers, apiToken.CanConfigureUsers)
	d.Set(APITokenFieldCanInstallNewAgents, apiToken.CanInstallNewAgents)
	d.Set(APITokenFieldCanSeeUsageInformation, apiToken.CanSeeUsageInformation)
	d.Set(APITokenFieldCanConfigureIntegrations, apiToken.CanConfigureIntegrations)
	d.Set(APITokenFieldCanSeeOnPremiseLicenseInformation, apiToken.CanSeeOnPremiseLicenseInformation)
	d.Set(APITokenFieldCanConfigureCustomAlerts, apiToken.CanConfigureCustomAlerts)
	d.Set(APITokenFieldCanConfigureAPITokens, apiToken.CanConfigureAPITokens)
	d.Set(APITokenFieldCanConfigureAgentRunMode, apiToken.CanConfigureAgentRunMode)
	d.Set(APITokenFieldCanViewAuditLog, apiToken.CanViewAuditLog)
	d.Set(APITokenFieldCanConfigureAgents, apiToken.CanConfigureAgents)
	d.Set(APITokenFieldCanConfigureAuthenticationMethods, apiToken.CanConfigureAuthenticationMethods)
	d.Set(APITokenFieldCanConfigureApplications, apiToken.CanConfigureApplications)
	d.Set(APITokenFieldCanConfigureTeams, apiToken.CanConfigureTeams)
	d.Set(APITokenFieldCanConfigureReleases, apiToken.CanConfigureReleases)
	d.Set(APITokenFieldCanConfigureLogManagement, apiToken.CanConfigureLogManagement)
	d.Set(APITokenFieldCanCreatePublicCustomDashboards, apiToken.CanCreatePublicCustomDashboards)
	d.Set(APITokenFieldCanViewLogs, apiToken.CanViewLogs)
	d.Set(APITokenFieldCanViewTraceDetails, apiToken.CanViewTraceDetails)
	d.Set(APITokenFieldCanConfigureSessionSettings, apiToken.CanConfigureSessionSettings)
	d.Set(APITokenFieldCanConfigureServiceLevelIndicators, apiToken.CanConfigureServiceLevelIndicators)
	d.Set(APITokenFieldCanConfigureGlobalAlertPayload, apiToken.CanConfigureGlobalAlertPayload)
	d.Set(APITokenFieldCanConfigureGlobalAlertConfigs, apiToken.CanConfigureGlobalAlertConfigs)
	d.Set(APITokenFieldCanViewAccountAndBillingInformation, apiToken.CanViewAccountAndBillingInformation)
	d.Set(APITokenFieldCanEditAllAccessibleCustomDashboards, apiToken.CanEditAllAccessibleCustomDashboards)
	d.SetId(apiToken.ID)
	return nil
}

func (r *apiTokenResource) MapStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	name := r.computeFullNameString(d, formatter)
	return &restapi.APIToken{
		ID:                                   d.Id(),
		AccessGrantingToken:                  d.Get(APITokenFieldAccessGrantingToken).(string),
		InternalID:                           d.Get(APITokenFieldInternalID).(string),
		Name:                                 name,
		CanConfigureServiceMapping:           d.Get(APITokenFieldCanConfigureServiceMapping).(bool),
		CanConfigureEumApplications:          d.Get(APITokenFieldCanConfigureEumApplications).(bool),
		CanConfigureMobileAppMonitoring:      d.Get(APITokenFieldCanConfigureMobileAppMonitoring).(bool),
		CanConfigureUsers:                    d.Get(APITokenFieldCanConfigureUsers).(bool),
		CanInstallNewAgents:                  d.Get(APITokenFieldCanInstallNewAgents).(bool),
		CanSeeUsageInformation:               d.Get(APITokenFieldCanSeeUsageInformation).(bool),
		CanConfigureIntegrations:             d.Get(APITokenFieldCanConfigureIntegrations).(bool),
		CanSeeOnPremiseLicenseInformation:    d.Get(APITokenFieldCanSeeOnPremiseLicenseInformation).(bool),
		CanConfigureCustomAlerts:             d.Get(APITokenFieldCanConfigureCustomAlerts).(bool),
		CanConfigureAPITokens:                d.Get(APITokenFieldCanConfigureAPITokens).(bool),
		CanConfigureAgentRunMode:             d.Get(APITokenFieldCanConfigureAgentRunMode).(bool),
		CanViewAuditLog:                      d.Get(APITokenFieldCanViewAuditLog).(bool),
		CanConfigureAgents:                   d.Get(APITokenFieldCanConfigureAgents).(bool),
		CanConfigureAuthenticationMethods:    d.Get(APITokenFieldCanConfigureAuthenticationMethods).(bool),
		CanConfigureApplications:             d.Get(APITokenFieldCanConfigureApplications).(bool),
		CanConfigureTeams:                    d.Get(APITokenFieldCanConfigureTeams).(bool),
		CanConfigureReleases:                 d.Get(APITokenFieldCanConfigureReleases).(bool),
		CanConfigureLogManagement:            d.Get(APITokenFieldCanConfigureLogManagement).(bool),
		CanCreatePublicCustomDashboards:      d.Get(APITokenFieldCanCreatePublicCustomDashboards).(bool),
		CanViewLogs:                          d.Get(APITokenFieldCanViewLogs).(bool),
		CanViewTraceDetails:                  d.Get(APITokenFieldCanViewTraceDetails).(bool),
		CanConfigureSessionSettings:          d.Get(APITokenFieldCanConfigureSessionSettings).(bool),
		CanConfigureServiceLevelIndicators:   d.Get(APITokenFieldCanConfigureServiceLevelIndicators).(bool),
		CanConfigureGlobalAlertPayload:       d.Get(APITokenFieldCanConfigureGlobalAlertPayload).(bool),
		CanConfigureGlobalAlertConfigs:       d.Get(APITokenFieldCanConfigureGlobalAlertConfigs).(bool),
		CanViewAccountAndBillingInformation:  d.Get(APITokenFieldCanViewAccountAndBillingInformation).(bool),
		CanEditAllAccessibleCustomDashboards: d.Get(APITokenFieldCanEditAllAccessibleCustomDashboards).(bool),
	}, nil
}

func (r *apiTokenResource) computeFullNameString(d *schema.ResourceData, formatter utils.ResourceNameFormatter) string {
	if d.HasChange(APITokenFieldName) {
		return formatter.Format(d.Get(APITokenFieldName).(string))
	}
	return d.Get(APITokenFieldFullName).(string)
}
