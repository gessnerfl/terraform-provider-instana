# Alerting Channel Email Resource

Alerting channel configuration for notifications to a specified list of email addresses.

API Documentation: <https://instana.github.io/openapi/#operation/getAlertingChannels>

The resource supports `default_name_prefix` and `default_name_suffix`. The string will be appended automatically
to the name of the alerting channel.

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
