provider "instana" {
  api_key     = "my-api-key"
  endpoint    = "https://tisdev-tis.instana.io"
  //timeout     = 60
  //max_retries = 5
}

resource "instana_rule" "my-custom-rule" {
  id                = "my-rule-id"
  name              = "my custom rule"
  entityType        = "nomadScheduler"
  metricName        = "nomad.client.allocations.blocked"
  rollup            = 0
  window            = 300000
  aggregation       = "sum"
  conditionOperator = ">"
  conditionValue    = 0.0
}