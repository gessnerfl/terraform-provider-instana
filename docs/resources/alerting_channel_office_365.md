# Alerting Channel Office 365 Resource

**Deprecated** This feature will be removed in version 2.x and should be replaced with `instana_alerting_channel`.

Alerting channel configuration for notifications to Office356.

API Documentation: <https://instana.github.io/openapi/#operation/getAlertingChannels>

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

## Import

Office 365 alerting channels can be imported using the `id`, e.g.:

```
$ terraform import instana_alerting_channel_office_365.my_channel 60845e4e5e6b9cf8fc2868da
```
