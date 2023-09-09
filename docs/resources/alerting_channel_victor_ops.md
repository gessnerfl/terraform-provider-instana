# Alerting Channel VictorOps Resource

**Deprecated** This feature will be removed in version 2.x and should be replaced with `instana_alerting_channel`.

Alerting channel configuration to integrate with VictorOps.

API Documentation: <https://instana.github.io/openapi/#operation/getAlertingChannels>

The resource supports `default_name_prefix` and `default_name_suffix`. The string will be appended automatically
to the name of the alerting channel.

## Example Usage

```hcl
resource "instana_alerting_channel_victor_ops" "example" {
  name        = "my-victor-ops-alerting-channel"
  api_key     = "my-victor-ops-api-key"
  routing_key = "my-victor-ops-routing-key"
}
```

## Argument Reference

* `name` - Required - the name of the alerting channel
* `api_key` - Required - the api key to authenticate at the VictorOps API
* `routing_key` - Required - the routing key used by VictoryOps to route the alert to the desired targe

## Import

VictorOps alerting channels can be imported using the `id`, e.g.:

```
$ terraform import instana_alerting_channel_victor_ops.my_channel 60845e4e5e6b9cf8fc2868da
```
