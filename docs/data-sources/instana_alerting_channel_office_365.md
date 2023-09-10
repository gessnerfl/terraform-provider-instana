# Data Source: instana_alerting_channel_office_365

**Deprecated:** This feature will be removed in version 2.x and should be replaced with `instana_alerting_channel` data
source

Retrieve ID for specific alerting channel of type office 365.

API Documentation: <https://instana.github.io/openapi/#operation/getAlertingChannels>

## Example Usage

```hcl
data "instana_alerting_channel_office_365" "example" {
  name = "My custom alerting channel name"
}
```

## Argument Reference

* `name` - Required - the name of the alerting channel. This is an exact string match and not a regex. Therefore it is
  case-sensitive.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `id` - Unique identifier of this resource