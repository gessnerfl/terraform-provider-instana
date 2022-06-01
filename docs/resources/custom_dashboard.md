# Custom Dashboard

Management of custom dashboards.

API Documentation: <https://instana.github.io/openapi/#tag/Custom-Dashboards>

Requires permission `canCreatePublicCustomDashboards` (Creation of public custom dashboards) and 
`canEditAllAccessibleCustomDashboards` (Management of all accessible custom dashboards) to be activated

The ID of the resource which is also used as unique identifier in Instana is auto generated!

Custom Dashboards support `default_name_prefix` and `default_name_suffix`. The string will be appended automatically to the title.

## Example Usage

```hcl
resource "instana_custom_dashboard" "example" {
  title = "Example Dashboard"

  access_rule { 
    access_type = "READ_WRITE"
    relation_type = "USER"
    related_id = "user-id-1"
  }

  access_rule {
    access_type = "READ_WRITE"
    related_id = "user-id-2"
    relation_type = "USER"
  }
  
  access_rule { 
    access_type = "READ"
    relation_type = "GLOBAL"
  }

  widgets = file("${path.module}/widgets.json")
}
``` 

## Argument Reference

* `title` - Required - the name of the custom dashboard
* `access_rule` - Required - configuration of access rules (sharing/permissions) of the custom dashboard
    * `access_type` - Required - type of granted access. Supported values are `READ` and `READ_WRITE`
    * `relation_type` - Required - type of the entity for which the access is granted. Supported values are: 
       `USER`, `API_TOKEN`, `ROLE`, `TEAM`, `GLOBAL` 
    * `related_id` - Optional - the id of the related entity for which access is granted. Required for all 
      `relation_type` except `GLOBAL`
* `widgets` - Required - JSON array of widget configurations. It is recommended to get this configuration via the 
  `Edit as Json` feature of custom dashboards in Instana UI and to adopt the configuration afterwards. It is also 
  recommended to store the configuration in dedicated json files. This allows the use of the built-in terraform functions
  `file` (<https://www.terraform.io/language/functions/file>) or `templatefile` (https://www.terraform.io/language/functions/templatefile)

## Import

Custom Dashboards can be imported using the `id` of the custom dashboard, e.g.:

```
$ terraform import instana_custom_dashboard.example 60845e4e5e6b9cf8fc2868da
```
