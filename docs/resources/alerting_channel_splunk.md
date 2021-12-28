# Alerting Channel Splunk Resource

Alerting channel configuration to integrate with Splunk.

API Documentation: <https://instana.github.io/openapi/#operation/getAlertingChannels>

The resource supports `default_name_prefix` and `default_name_suffix`. The string will be appended automatically
to the name of the alerting channel.

## Example Usage

```hcl
resource "instana_alerting_channel_splunk" "example" {
  name  = "my-splunk-alerting-channel"
  url   = "https://my.splunk.url.example.com"
  token = "my-splunk-token"
}
```

## Argument Reference

* `name` - Required - the name of the alerting channel
* `url` - Required - the target Splunk endpoint URL
* `token` - Required - the authentication token to login at the Splunk API

## Import

Splunk alerting channels can be imported using the `id`, e.g.:

```
$ terraform import instana_alerting_channel_splunk.my_channel 60845e4e5e6b9cf8fc2868da
```
