# Custom Event Specification with Threshold Rule Resource

Configuration of a custom event specification based on a threshold rule. A threshold rule is verifies if a certain 
condition applies to a given metric. Therefore you can either use `rule_rollup` or `rule_window` or both to define 
the data points which should be evaluated. Instana API always returns max. 600 data points for validation.

- `rule_window` = the time frame in seconds where the aggregation is applied to
- `rule_rollup` = the resolution of the data points which are considered for this event (See also <https://instana.github.io/openapi/#tag/Infrastructure-Metrics>)

Both are optional in the Instana API. Usually configurations define a **window** for calculating the event.

API Documentation: <https://instana.github.io/openapi/#operation/putCustomEventSpecification>

Custom event resources support `default_name_prefix` and `default_name_suffix`. The string will be appended automatically
to the name of the custom event.

## Example Usage

### Built in metric

```hcl
resource "instana_custom_event_spec_threshold_rule" "nomad_pending_allocations" {
  name            = "Nomad pending allocations"
  description     = "Pending allocations in nomad. Ensure the EC2 auto-scaling group is configured correct and sufficient instances are running"
  query           = "entity.tag:\"stage=prod\""
  enabled         = true
  triggering      = true
  expiration_time = 60000
  entity_type     = "nomadScheduler"

  rule_severity           = "critical"
  rule_metric_name        = "nomad.client.allocations.pending"
  rule_window             = 60000
  rule_aggregation        = "avg"
  rule_condition_operator = ">"
  rule_condition_value    = 0
}
```

### Built in dynamic metric

```hcl
resource "instana_custom_event_spec_threshold_rule" "mysql_number_of_writes" {
  name            = "mysql-number-of-writes-exceeded"
  description     = "Insert count exceeds supported write operations"
  query           = "entity.tag:\"stage=prod\""
  enabled         = true
  triggering      = true
  expiration_time = 60000
  entity_type     = "mySqlDatabase"

  rule_severity           = "warning"
  rule_window             = 60000
  rule_aggregation        = "avg"
  rule_condition_operator = ">"
  rule_condition_value    = 1000000

  rule_metric_pattern_prefix      = "databases"
  rule_metric_pattern_postfix     = "insert_count"
  rule_metric_pattern_placeholder = "myschema"
  rule_metric_pattern_operator    = "is"
}
```

### Custom metric

```hcl
resource "instana_custom_event_spec_threshold_rule" "custom_health_metrics_unhealthy" {
  name            = "MyApp is unhealthy"
  description     = "The custom health indicator of the application switched to unhealthy"
  query           = "entity.tag:\"stage=prod\" AND entity.service.name:myapp"
  enabled         = true
  triggering      = true
  expiration_time = 60000
  entity_type     = "jvmRuntimePlatform"

  rule_severity           = "warning"
  rule_metric_name        = "metrics.gauges.myapp.healthIndicator"
  rule_window             = 10000
  rule_aggregation        = "avg"
  rule_condition_operator = "="
  rule_condition_value    = 1
}
```

## Argument Reference

* `name` - Required - The name of the custom event specification
* `description` - Required - The description text of the custom event specification
* `query` - Optional - The dynamic filter query for which the rule should be applied to
* `enabled` - Optional - Boolean flag if the rule should be enabled - default = true
* `triggering` - Optional - Boolean flag if the rule should trigger an incident - default = false
* `expiration_time` - Optional - The grace period in milliseconds until the issue is closed
* `entity_type` - Required - The entity type/plugin for which the verification rule will be defined
Supported entity types (plugins) can be retrieved from the Instana REST API using the path
`/api/infrastructure-monitoring/catalog/plugins`.
* `rule_severity` - Required - The severity of the rule - allowed values: `warning`, `critical`
  
* `rule_metric_name` - Required (Built-In and Custom Metrics only) The name of the built in or custom metric name (supported
built in metrics can be retrieved from the REST API using the endpoint `/api/infrastructure-monitoring/catalog/metrics/{plugin}`)

* `rule_metric_pattern_prefix` - Required (Dynamic Built-In Metrics only) The prefix of the built in dynamic metric
* `rule_metric_pattern_postfix` - Optional (Dynamic Built-In Metrics only) The postfix of the built in dynamic metric
* `rule_metric_pattern_placeholder` - Required (Dynamic Built-In Metrics only) The placeholder string of the dynamic metric
* `rule_metric_pattern_operator` - Required (Dynamic Built-In Metrics only) The operation used to check for matching
placeholder string. Allowed values:  `is`, `contains`, `any`, `startsWith`, `endsWith`

* `rule_window` - Optional - The time window in milliseconds within the rule condition is applied to
* `rule_rollup` - Optional - The resolution of the monitored metrics
* `rule_aggregation` - Optional (depending on metric type) - the aggregation used to calculate the metric value for the given
time window and/or rollup. Supported value: `sum`, `avg`, `min`, `max`
* `rule_condition_operator` - Required - The condition operator used to check against the calculated metric value for the given
time window and/or rollup. Supported values: `=` (`==` also supported as an alternative representation for equals), `!=`, `<=`, 
`<`, `>`, `=>`
* `rule_condition_value` - Required - The numeric condition value used to check against the calculated metric value for the given
time window and/or rollup.