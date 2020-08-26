# User Role Resource

Management of user roles.

API Documentation: <https://instana.github.io/openapi/#operation/getRole>

The ID of the resource which is also used as unique identifier in Instana is auto generated!
The resource does NOT support `default_name_prefix` and `default_name_suffix`.

## Example Usage

```hcl
resource "instana_user_role" "example" {
  name                                   = "name"
  implicit_view_filter                   = "view filter"
  can_configure_service_mapping          = true
  can_configure_eum_applications         = true
  can_configure_users                    = true
  can_install_new_agents                 = true
  can_see_usage_information              = true
  can_configure_integrations             = true
  can_see_on_premise_license_information = true
  can_configure_roles                    = true
  can_configure_custom_alerts            = true
  can_configure_api_tokens               = true
  can_configure_agent_run_mode           = true
  can_view_audit_log                     = true
  can_configure_objectives               = true
  can_configure_agents                   = true
  can_configure_authentication_methods   = true
  can_configure_applications             = true
}
```
```

## Argument Reference

* `name` - Required - the name of the alerting channel
* `implicit_view_filter` - Optional - NOT Supported anymore!!
* `can_configure_service_mapping` - Optional - default false - enables permission to configure service mappings
* `can_configure_eum_applications` - Optional - default false - enables permission to configure EUM applications
* `can_configure_users` - Optional - default false - enables permission to configure users
* `can_install_new_agents` - Optional - default false - enables permission to install new agents
* `can_see_usage_information` - Optional - default false - enables permission to see usage information
* `can_configure_integrations` - Optional - default false - enables permission to configure integrations
* `can_see_on_premise_license_information` - Optional - default false - enables permission to see on premise license information
* `can_configure_roles` - Optional - default false - enables permission to configure roles
* `can_configure_custom_alerts` - Optional - default false - enables permission to configure custom alerts
* `can_configure_api_tokens` - Optional - default false - enables permission to configure api tokes
* `can_configure_agent_run_mode` - Optional - default false - enables permission to configure agent run mode
* `can_view_audit_log` - Optional - default false - enables permission to view audit logs
* `can_configure_objectives` - Optional - default false - enables permission to configure objectives
* `can_configure_agents` - Optional - default false - enables permission to configure agents
* `can_configure_authentication_methods` - Optional - default false - enables permission to configure authentication methods
* `can_configure_applications` - Optional - default false - enables permission to configure applications
