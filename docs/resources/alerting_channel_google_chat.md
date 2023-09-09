# Alerting Channel Google Chat Resource

**Deprecated** This feature will be removed in version 2.x and should be replaced with `instana_alerting_channel`.

Alerting channel configuration for notifications to Google Chat.

API Documentation: <https://instana.github.io/openapi/#operation/getAlertingChannels>

## Example Usage

```hcl
resource "instana_alerting_channel_google_chat" "example" {
  name        = "my-google-chat-alerting-channel"
  webhook_url = "https://my.google.chat.weebhook.exmaple.com/"
}
```

## Argument Reference

* `name` - Required - the name of the alerting channel
* `webhook_url` - Required - the URL of the Google Chat Webhook where the alert will be sent to

## Import

Google Chat alerting channels can be imported using the `id`, e.g.:

```
$ terraform import instana_alerting_channel_google_chat.my_channel 60845e4e5e6b9cf8fc2868da
```
