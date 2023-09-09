# Alerting Channel Slack Resource

**Deprecated** This feature will be removed in version 2.x and should be replaced with `instana_alerting_channel`.

Alerting channel configuration notifications to Slack.

API Documentation: <https://instana.github.io/openapi/#operation/getAlertingChannels>

The resource supports `default_name_prefix` and `default_name_suffix`. The string will be appended automatically
to the name of the alerting channel.

## Example Usage

```hcl
resource "instana_alerting_channel_slack" "example" {
  name        = "my-slack-alerting-channel"
  webhook_url = "https://my.slack.weebhook.exmaple.com/"
  icon_url    = "https://my.slack.icon.exmaple.com/"
  channel     = "my-channel"
}
```

## Argument Reference

* `name` - Required - the name of the alerting channel
* `webhook_url` - Required - the URL of the Slack webhook to send alerts to
* `icon_url` - Optional - the URL to the icon which should be rendered in the slack message
* `channel` - Optional - the target Slack channel where the alert should be posted 

## Import

Slack alerting channels can be imported using the `id`, e.g.:

```
$ terraform import instana_alerting_channel_slack.my_channel 60845e4e5e6b9cf8fc2868da
```
