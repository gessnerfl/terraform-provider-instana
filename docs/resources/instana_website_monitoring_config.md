# Website Monitoring Config Resource

Resource to configure websites in Instana

API Documentation: <https://instana.github.io/openapi/#tag/Website-Configuration>

The resource supports `default_name_prefix` and `default_name_suffix`. The string will be appended automatically
to the name of the website monitoring config.

## Example Usage

```hcl
resource "instana_website_monitoring_config" "example" {
  name = "my-website-monitoring-config"
}
```

## Argument Reference

* `name` - Required - the name of the website monitoring config
