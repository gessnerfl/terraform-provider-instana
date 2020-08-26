# Alerting Channel Pager Duty Resource

Alerting channel configuration integration with PagerDuty.

API Documentation: <https://instana.github.io/openapi/#operation/getAlertingChannels>

The resource supports `default_name_prefix` and `default_name_suffix`. The string will be appended automatically
to the name of the alerting channel.

## Example Usage

```hcl
resource "instana_alerting_channel_pager_duty" "example" {
  name = "my-pager-duty-alerting-channel"
  service_integration_key = "my-service-integration-key"
}
```

## Argument Reference

* `name` - Required - the name of the alerting channel
* `service_integration_key` - Required - the key for the service integration in pager duty
