# Alerting Channel Ops Genie Resource

**Deprecated** This feature will be removed in version 2.x and should be replaced with `instana_alerting_channel`.

Alerting channel configuration integration with OpsGenie.

API Documentation: <https://instana.github.io/openapi/#operation/getAlertingChannels>

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

## Import

OpsGenie alerting channels can be imported using the `id`, e.g.:

```
$ terraform import instana_alerting_channel_ops_genie.my_channel 60845e4e5e6b9cf8fc2868da
```
