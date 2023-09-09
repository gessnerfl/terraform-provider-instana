# Alerting Channel Email Resource

**Deprecated** This feature will be removed in version 2.x and should be replaced with `instana_alerting_channel`.

Alerting channel configuration for notifications to a specified list of email addresses.

API Documentation: <https://instana.github.io/openapi/#operation/getAlertingChannels>

## Example Usage

```hcl
resource "instana_alerting_channel_email" "example" {
  name = "my-email-alerting-channel"
  emails = [ "email1@example.com", "email2@example.com" ]
}
```

## Argument Reference

* `name` - Required - the name of the alerting channel
* `emails` - Required - the list of target email addresses

## Import

Email alerting channels can be imported using the `id`, e.g.:

```
$ terraform import instana_alerting_channel_email.my_channel 60845e4e5e6b9cf8fc2868da
```