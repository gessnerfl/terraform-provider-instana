# Alerting Channel Resource

Alerting channel configuration for notifications to a specified target channel.

API Documentation: <https://instana.github.io/openapi/#operation/getAlertingChannels>

## Example Usage

### Email Alerting Channel

```hcl
resource "instana_alerting_channel" "example" {
  name = "my-email-alerting-channel"
  
  email {
    emails = [ "email1@example.com", "email2@example.com" ]
  }
}
```

### Google Chat Alerting Channel

```hcl
resource "instana_alerting_channel" "example" {
  name        = "my-google-chat-alerting-channel"
  
  google_chat {
    webhook_url = "https://my.google.chat.weebhook.exmaple.com/"
  }
}
```

### Office 365 Alerting Channel

```hcl
resource "instana_alerting_channel" "example" {
  name        = "my-google-chat-alerting-channel"
  
  office_365 {
    webhook_url = "https://my.google.chat.weebhook.exmaple.com/"
  }
}
```

### OpsGenie Alerting Channel

```hcl
resource "instana_alerting_channel" "example" {
  name = "my-ops-genie-alerting-channel"
  
  ops_genie {
    api_key = "my-secure-api-key"
    tags = [ "tag1", "tag2" ]
    region = "EU"
  }
}
```

### PagerDuty Alerting Channel

```hcl
resource "instana_alerting_channel" "example" {
  name = "my-pager-duty-alerting-channel"
  
  pager_duty {
    service_integration_key = "my-service-integration-key"
  }
}
```

### Slack Alerting Channel

```hcl
resource "instana_alerting_channel" "example" {
  name        = "my-slack-alerting-channel"
  
  slack {
    webhook_url = "https://my.slack.weebhook.exmaple.com/"
    icon_url    = "https://my.slack.icon.exmaple.com/"
    channel     = "my-channel"
  }
}
```

### Splunk Alerting Channel

```hcl
resource "instana_alerting_channel" "example" {
  name  = "my-splunk-alerting-channel"
  
  splunk {
    url   = "https://my.splunk.url.example.com"
    token = "my-splunk-token"
  }
}
```

### VictorOps Alerting Channel

```hcl
resource "instana_alerting_channel" "example" {
  name        = "my-victor-ops-alerting-channel"
  
  victor_ops {
    api_key     = "my-victor-ops-api-key"
    routing_key = "my-victor-ops-routing-key"
  }
}
```

### Webhook Alerting Channel

```hcl
resource "instana_alerting_channel" "example" {
  name         = "my-generic-webhook-alerting-channel"
  
  webhook {
    webhook_urls = [ 
      "https://my.weebhook1.exmaple.com/", 
      "https://my.weebhook2.exmaple.com/" 
    ]
    
    http_headers = {
      header1 = "headerValue1"
      header2 = "headerValue2"
    }
  }
}
```

## Argument Reference

* `name` - Required - the name of the alerting channel

Exactly one of the following channel types must be configured:

* `email` - Optional - configuration of a email alerting channel - [Details](#email)
* `google_chat` - Optional - configuration of a Google Chat alerting channel - [Details](#google-chat)
* `office_365` - Optional - configuration of a Office 365 alerting channel - [Details](#office-365)
* `ops_genie` - Optional - configuration of a OpsGenie alerting channel - [Details](#opsgenie)
* `pager_duty` - Optional - configuration of a PagerDuty alerting channel - [Details](#pagerduty)
* `slack` - Optional - configuration of a Slack alerting channel - [Details](#slack)
* `splunk` - Optional - configuration of a Splunk alerting channel - [Details](#splunk)
* `victor_ops` - Optional - configuration of a VictorOps alerting channel - [Details](#victorops)
* `webhook` - Optional - configuration of a webhook alerting channel - [Details](#webhook)

### Email

* `emails` - Required - the list of target email addresses

### Google Chat

* `webhook_url` - Required - the URL of the Google Chat Webhook where the alert will be sent to

### Office 365

* `webhook_url` - Required - the URL of the Google Chat Webhook where the alert will be sent to

### OpsGenie

* `api_key` - Required - the API Key for authentication at the Ops Genie API
* `tags` - Required - a list of tags (strings) for the alert in Ops Genie
* `region` - Required - the target Ops Genie region

### PagerDuty

* `service_integration_key` - Required - the key for the service integration in pager duty

### Slack

* `webhook_url` - Required - the URL of the Slack webhook to send alerts to
* `icon_url` - Optional - the URL to the icon which should be rendered in the slack message
* `channel` - Optional - the target Slack channel where the alert should be posted

### Splunk

* `url` - Required - the target Splunk endpoint URL
* `token` - Required - the authentication token to login at the Splunk API

### VictorOps

* `api_key` - Required - the api key to authenticate at the VictorOps API
* `routing_key` - Required - the routing key used by VictoryOps to route the alert to the desired targe

### Webhook

* `webhook_urls` - Required - the list of webhook URLs where the alert will be sent to
* `http_headers` - Optional - key/value map of additional http headers which will be sent to the webhook

## Import

Email alerting channels can be imported using the `id`, e.g.:

```
$ terraform import instana_alerting_channel_email.my_channel 60845e4e5e6b9cf8fc2868da
```