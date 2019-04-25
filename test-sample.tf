provider "instana" {
  api_token   = ""
  endpoint    = ""
}

resource "instana_application_config" "terraform-test" {
  label              = "Terraform Test"
  match_specification = "aws.ec2.tag.stage EQUALS 'live-test' AND entity.type EQUALS 'mysql' OR entity.type EQUALS 'elasticsearch'"
}

resource "instana_custom_event_spec_system_rule" "example" {
  name = "test-instana-system-rule"
  entity_type = "any"
  query = "entity.service.name:\"btm-payment-export\" AND entity.tag:stage=live-test"
  enabled = true
  triggering = true
  description = "Terraform test of system rule"
  expiration_time = "60000"
	rule_severity = "warning"
	rule_system_rule_id = "entity.offline"
}

resource "instana_custom_event_spec_threshold_rule" "example" {
  name = "test-instana-threshold-rule"
  entity_type = "nomadScheduler"
  query = "entity.tag:\"stage=dev\" AND entity.tag:\"region=aws-us-east-1\""
  enabled = true
  triggering = true
  description = "Terraform test of threshold rule"
  expiration_time = "60000"
  rule_severity = "warning"
  rule_metric_name = "nomad.client.allocations.pending"
  rule_window = "60000"
  rule_aggregation = "sum"
  rule_condition_operator = ">"
  rule_condition_value = "0.0"
}