# Alerting Channel Office 365 Resource

Alerting channel configuration to integrate with WebHooks.

API Documentation: <https://instana.github.io/openapi/#operation/getAlertingChannels>

The resource supports `default_name_prefix` and `default_name_suffix`. The string will be appended automatically
to the name of the alerting channel.

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