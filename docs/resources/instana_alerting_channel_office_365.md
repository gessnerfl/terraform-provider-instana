# Alerting Channel Office 365 Resource

Alerting channel configuration for notifications to Office356.

API Documentation: <https://instana.github.io/openapi/#operation/getAlertingChannels>

The resource supports `default_name_prefix` and `default_name_suffix`. The string will be appended automatically
to the name of the alerting channel.

## Example Usage

```hcl
resource "instana_alerting_channel_office_365" "example" {
  name        = "my-office365-alerting-channel"
  webhook_url = "https://my.office365.weebhook.exmaple.com/"
}
```

## Argument Reference

* `name` - Required - the name of the alerting channel
* `webhook_url` - Required - the URL of the Office 365 Webhook where the alert will be sent to
