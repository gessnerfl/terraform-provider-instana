# Data Source: instana_alerting_channel_office_365

Retrieve ID for specific alerting channel of type office 365.

API Documentation: <https://instana.github.io/openapi/#operation/getAlertingChannels>

## Example Usage

```hcl
data "instana_alerting_channel_office_365" "example" {
  name = "My custom alerting channel name"
}
```

## Argument Reference

* `name` - Required - the name of the alerting channel

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `id` - Unique identifier of this resource