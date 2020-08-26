# Alerting Channel Ops Genie Resource

Alerting channel configuration integration with OpsGenie.

API Documentation: <https://instana.github.io/openapi/#operation/getAlertingChannels>

The resource supports `default_name_prefix` and `default_name_suffix`. The string will be appended automatically
to the name of the alerting channel.

## Example Usage

```hcl
resource "instana_alerting_channel_ops_genie" "example" {
  name = "my-ops-genie-alerting-channel"
  api_key = "my-secure-api-key"
  tags = [ "tag1", "tag2" ]
  region = "EU"
}
```

## Argument Reference

* `name` - Required - the name of the alerting channel
* `api_key` - Required - the API Key for authentication at the Ops Genie API
* `tags` - Required - a list of tags (strings) for the alert in Ops Genie
* `region` - Required - the target Ops Genie region
