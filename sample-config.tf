provider "instana" {
  api_key     = "my-api-key"
  endpoint    = "https://mytenant-mycustomername.instana.io"
  //timeout     = 60
  //max_retries = 5
}

resource "instana_rule" "my-custom-rule" {
  identifier         = "my-rule-id"
  name               = "my custom rule"
  entity_type        = "nomadScheduler"
  metric_name        = "nomad.client.allocations.blocked"
  rollup             = 0
  window             = 300000
  aggregation        = "sum"
  condition_operator = ">"
  condition_value    = 0.0
}