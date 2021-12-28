# RBAC Group

Management of Groups for role based access control.

API Documentation: <https://instana.github.io/openapi/#tag/Groups>

The ID of the resource which is also used as unique identifier in Instana is auto generated!

Groups support `default_name_prefix` and `default_name_suffix`. The string will be appended automatically to the name.

## Example Usage

```hcl
resource "instana_rbac_group" "example" {
  name = "test"

  permission_set {
    application_ids             = ["app_id1", "app_id2"]
    kubernetes_cluster_uuids    = ["k8s_cluster_id1", "k8s_cluster_id2"]
    kubernetes_namespaces_uuids = ["k8s_namespace_id1", "k8s_namespace_id2"]
    mobile_app_ids              = ["mobile_app_id1", "mobile_app_id2"]
    website_ids                 = ["website_id1", "website_id2"]
    infra_dfq_filter            = "infra_dfq"
    permissions                 = ["CAN_CONFIGURE_APPLICATIONS", "CAN_CONFIGURE_AGENTS"]
  }
}
``` 

## Argument Reference

* `name` - Required - the name of the RBAC group
* `permission_set` - Optional - resource block to describe the assigned permissions
    * `application_ids` - Optional - list of application ids which are permitted to the given group
    * `kubernetes_cluster_uuids` - Optional - list of Kubernetes Cluster UUIDs which are permitted to the given group
    * `kubernetes_namespaces_uuids` - Optional - list of Kubernetes Namespaces UUIDs which are permitted to the given
      group
    * `mobile_app_ids` - Optional - list of mobile app ids which are permitted to the given group
    * `website_ids` - Optional -list of website ids which are permitted to the given group
    * `infra_dfq_filter` - Optional - a dynamic focus query to restrict access to a limited set of infrastructure
      resources
    * `permissions` - Optional - the list of permissions granted to the given group. Allowed values
      are: 
      * `CAN_CONFIGURE_APPLICATIONS`
      * `CAN_SEE_ON_PREM_LICENE_INFORMATION`
      * `CAN_CONFIGURE_EUM_APPLICATIONS`
      * `CAN_CONFIGURE_AGENTS`
      * `CAN_VIEW_TRACE_DETAILS`
      * `CAN_VIEW_LOGS`
      * `CAN_CONFIGURE_SESSION_SETTINGS`
      * `CAN_CONFIGURE_INTEGRATIONS`
      * `CAN_CONFIGURE_GLOBAL_ALERT_CONFIGS`
      * `CAN_CONFIGURE_GLOBAL_ALERT_PAYLOAD`
      * `CAN_CONFIGURE_MOBILE_APP_MONITORING`
      * `CAN_CONFIGURE_API_TOKENS`
      * `CAN_CONFIGURE_SERVICE_LEVEL_INDICATORS`
      * `CAN_CONFIGURE_AUTHENTICATION_METHODS`
      * `CAN_CONFIGURE_RELEASES`
      * `CAN_VIEW_AUDIT_LOG`
      * `CAN_CONFIGURE_CUSTOM_ALERTS`
      * `CAN_CONFIGURE_AGENT_RUN_MODE`
      * `CAN_CONFIGURE_SERVICE_MAPPING`
      * `CAN_SEE_USAGE_INFORMATION`
      * `CAN_EDIT_ALL_ACCESSIBLE_CUSTOM_DASHBOARDS`
      * `CAN_CONFIGURE_USERS`
      * `CAN_INSTALL_NEW_AGENTS`
      * `CAN_CONFIGURE_TEAMS`
      * `CAN_CREATE_PUBLIC_CUSTOM_DASHBOARDS`
      * `CAN_CONFIGURE_LOG_MANAGEMENT`
      * `CAN_VIEW_ACCOUNT_AND_BILLING_INFORMATION`

## Import

RBAC Groups can be imported using the `id` of the group, e.g.:

```
$ terraform import instana_rbac_group.my_group 60845e4e5e6b9cf8fc2868da
```
