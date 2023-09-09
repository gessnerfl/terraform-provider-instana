# Alerting Channel Webhook Resource

**Deprecated** This feature will be removed in version 2.x and should be replaced with `instana_alerting_channel`.

Alerting channel configuration to integrate with WebHooks.

API Documentation: <https://instana.github.io/openapi/#operation/getAlertingChannels>

## Example Usage

```hcl
resource "instana_alerting_channel_webhook" "example" {
  name         = "my-generic-webhook-alerting-channel"
  webhook_urls = [ "https://my.weebhook1.exmaple.com/", "https://my.weebhook2.exmaple.com/" ]
  http_headers = {      #Optional
    header1 = "headerValue1"
    header2 = "headerValue2"
  }
}
```

## Argument Reference

* `name` - Required - the name of the alerting channel
* `webhook_urls` - Required - the list of webhook URLs where the alert will be sent to
* `http_headers` - Optional - key/value map of additional http headers which will be sent to the webhook

## Import

Webhook alerting channels can be imported using the `id`, e.g.:

```
$ terraform import instana_alerting_channel_webhook.my_channel 60845e4e5e6b9cf8fc2868da
```
