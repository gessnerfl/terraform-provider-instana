provider "instana" {
  api_token = ""
  endpoint  = ""
}

resource "instana_application_config" "terraform-test" {
  label               = "Terraform Test"
  match_specification = "aws.ec2.tag.stage EQUALS 'live-test' AND entity.type EQUALS 'mysql' OR entity.type EQUALS 'elasticsearch'"
}

resource "instana_custom_event_spec_system_rule" "example" {
  name                = "test-instana-system-rule"
  query               = "entity.service.name:\"btm-payment-export\" AND entity.tag:stage=live-test"
  enabled             = true
  triggering          = true
  description         = "Terraform test of system rule"
  expiration_time     = 60000
  rule_severity       = "warning"
  rule_system_rule_id = "entity.offline"
}

resource "instana_custom_event_spec_threshold_rule" "example" {
  name                    = "test-instana-threshold-rule"
  entity_type             = "nomadScheduler"
  query                   = "entity.tag:\"stage=dev\" AND entity.tag:\"region=aws-us-east-1\""
  enabled                 = true
  triggering              = true
  description             = "Terraform test of threshold rule"
  expiration_time         = 60000
  rule_severity           = "warning"
  rule_metric_name        = "nomad.client.allocations.pending"
  rule_window             = 60000
  rule_aggregation        = "sum"
  rule_condition_operator = ">"
  rule_condition_value    = 0.0
}

resource "instana_custom_event_spec_entity_verification_rule" "example" {
  name                       = "name"
  query                      = "query"
  enabled                    = true
  triggering                 = true
  description                = "description"
  expiration_time            = 60000
  rule_severity              = "warning"
  rule_matching_entity_type  = "process"
  rule_matching_operator     = "is"
  rule_matching_entity_label = "entity-label"
  rule_offline_duration      = 60000
}

resource "instana_alerting_channel_email" "example" {
  name   = "my-email-alerting-channel"
  emails = ["email1@example.com", "email2@example.com"]
}

resource "instana_alerting_channel_google_chat" "example" {
  name        = "my-google-chat-alerting-channel"
  webhook_url = "https://my.google.chat.weebhook.exmaple.com/"
}

resource "instana_alerting_channel_office_365" "example" {
  name        = "my-office365-alerting-channel"
  webhook_url = "https://my.office365.weebhook.exmaple.com/"
}

resource "instana_alerting_channel_ops_genie" "example" {
  name    = "my-ops-genie-alerting-channel"
  api_key = "my-secure-api-key"
  tags    = ["tag1", "tag2"]
  region  = "EU"
}

resource "instana_alerting_channel_pager_duty" "example" {
  name                    = "my-pager-duty-alerting-channel"
  service_integration_key = "my-service-integration-key"
}

resource "instana_alerting_channel_slack" "example" {
  name        = "my-slack-alerting-channel"
  webhook_url = "https://my.slack.weebhook.exmaple.com/"
  icon_url    = "https://my.slack.icon.exmaple.com/" #Optional
  channel     = "my-channel"                         #Optional
}

resource "instana_alerting_channel_splunk" "example" {
  name  = "my-splunk-alerting-channel"
  url   = "https://my.splunk.url.example.com"
  token = "my-splunk-token"
}

resource "instana_alerting_channel_victor_ops" "example" {
  name        = "my-victor-ops-alerting-channel"
  api_key     = "my-victor-ops-api-key"
  routing_key = "my-victor-ops-routing-key"
}

resource "instana_alerting_channel_webhook" "example" {
  name         = "my-generic-webhook-alerting-channel"
  webhook_urls = ["https://my.weebhook1.exmaple.com/", "https://my.weebhook2.exmaple.com/"]
  http_headers = { #Optional
    header1 = "headerValue1"
    header2 = "headerValue2"
  }
}

resource "instana_alerting_config" "alerting_for_rules" {
  alert_name = "name"
  integration_ids = [
    "${instana_alerting_channel_email.example.id}",
    "${instana_alerting_channel_google_chat.example.id}"
  ]
  event_filter_query = "entity.tag:stage=live-test"
  event_filter_rule_ids = [
    "${instana_custom_event_spec_system_rule.example.id}",
    "${instana_custom_event_spec_threshold_rule.example.id}",
    "${instana_custom_event_spec_entity_verification_rule.example.id}"
  ]
}

resource "instana_alerting_config" "alerting_for_event_types" {
  alert_name = "name"
  integration_ids = [
    "${instana_alerting_channel_pager_duty.example.id}",
    "${instana_alerting_channel_email.example.id}"
  ]
  event_filter_query       = "entity.tag:stage=live-test"
  event_filter_event_types = ["incident", "critical"]
}

data "instana_synthetic_location" "monitor" {
  label         = "label-name"
  location_type = "Private"
}

resource "instana_synthetic_test" "uptime_check" {
  label          = "test"
  active         = true
  locations      = [data.instana_synthetic_location.monitor.id]
  test_frequency = 10
  playback_mode  = "Staggered"
  configuration {
    mark_synthetic_call = true
    retries             = 0
    retry_interval      = 1
    synthetic_type      = "HTTPAction"
    timeout             = "3m"
    url                 = "https://example.com"
    operation           = "GET"
  }
  custom_properties = {
    "foo" = "bar"
  }
}

resource "instana_synthetic_test" "http_action" {
  label     = "test"
  active    = true
  locations = [data.instana_synthetic_location.monitor.id]
  configuration {
    synthetic_type = "HTTPScript"
    script         = <<EOF
      const assert = require('assert');

      $http.get('https://terraform.io',
      function(err, response, body) {
          if (err) {
              console.error(err);
              return;
          }
          assert.equal(response.statusCode, 200, 'Expected a 200 OK response');
      });
    EOF
  }
}
