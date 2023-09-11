# Alerting Channel Data Source

Data source to retrieve details about existing alerting channels

API Documentation: <https://instana.github.io/openapi/#operation/getAlertingChannels>

## Example Usage

```hcl
data "instana_alerting_channel" "example" {
  name = "my-alerting-channel"
}
```

## Argument Reference

* `name` - Required - the name of the alerting channel

## Attribute Reference

Exactly one of the following items is provided depending on the type of the alerting channel:

* `email` - configuration of a email alerting channel - [Details](#email)
* `google_chat` - configuration of a Google Chat alerting channel - [Details](#google-chat)
* `office_365` - configuration of a Office 365 alerting channel - [Details](#office-365)
* `ops_genie` - configuration of a OpsGenie alerting channel - [Details](#opsgenie)
* `pager_duty` - configuration of a PagerDuty alerting channel - [Details](#pagerduty)
* `slack` - configuration of a Slack alerting channel - [Details](#slack)
* `splunk` - configuration of a Splunk alerting channel - [Details](#splunk)
* `victor_ops` - configuration of a VictorOps alerting channel - [Details](#victorops)
* `webhook` - configuration of a webhook alerting channel - [Details](#webhook)

### Email

* `emails` - the list of target email addresses

### Google Chat

* `webhook_url` - the URL of the Google Chat Webhook where the alert will be sent to

### Office 365

* `webhook_url` - the URL of the Google Chat Webhook where the alert will be sent to

### OpsGenie

* `api_key` - the API Key for authentication at the Ops Genie API
* `tags` - a list of tags (strings) for the alert in Ops Genie
* `region` - the target Ops Genie region

### PagerDuty

* `service_integration_key` - the key for the service integration in pager duty

### Slack

* `webhook_url` - the URL of the Slack webhook to send alerts to
* `icon_url` - the URL to the icon which should be rendered in the slack message
* `channel` - the target Slack channel where the alert should be posted

### Splunk

* `url` - the target Splunk endpoint URL
* `token` - the authentication token to login at the Splunk API

### VictorOps

* `api_key` - the api key to authenticate at the VictorOps API
* `routing_key` - the routing key used by VictoryOps to route the alert to the desired targe

### Webhook

* `webhook_urls` - the list of webhook URLs where the alert will be sent to
* `http_headers` - key/value map of additional http headers which will be sent to the webhook
