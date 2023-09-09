# Custom Event Specification with Entity Verification Rule Resource

Configuration of a custom event specification based on an entity verification rule. This rule type is used
to check for hosts which do not have matching entities running on them.

API Documentation: <https://instana.github.io/openapi/#operation/putCustomEventSpecification>

`default_name_prefix` and `default_name_prefix` is **NOT** supported for this resource as this feature will
be removed in version 2.x.

## Example Usage

### Entity Verification Rule

```hcl
resource "instana_custom_event_specification" "example" {
  name            = "name"
  description     = "description"
  query           = "query"
  enabled         = true
  triggering      = true
  expiration_time = 60000
  entity_type     = "host"

  rules {  
    entity_verification {
      severity              = "warning"
      matching_entity_type  = "process"
      matching_operator     = "is"
      matching_entity_label = "entity-label"
      offline_duration      = 60000
    }
  } 
}
```

### System Rule

```hcl
resource "instana_custom_event_specification" "example" {
  name            = "name"
  description     = "description"
  query           = "query"
  enabled         = true
  triggering      = true
  expiration_time = 60000
  entity_type     = "any"

  rules { 
    system {
      severity       = "warning"
      system_rule_id = "system-rule-id"
    }
  } 
}
```

### Threshold Rule

```hcl
resource "instana_custom_event_specification" "example" {
  name            = "name"
  description     = "description"
  query           = "query"
  enabled         = true
  triggering      = true
  expiration_time = 60000
  entity_type     = "nomadScheduler"

  rules { 
    threshold {
      severity           = "critical"
      metric_name        = "nomad.client.allocations.pending"
      window             = 60000
      aggregation        = "avg"
      condition_operator = ">"
      condition_value    = 0
    }
  } 
}
```

## Argument Reference

* `name` - Required - The name of the custom event specification
* `description` - Required - The description text of the custom event specification
* `entity_type` - Required - The entity type/plugin for which the verification rule will be defined. Must be set to
  `any` for [System Rules](#system-rule) and `host` for [Entity Verification Rules](#entity-verification-rule). For
  threshold rules the supported entity types (plugins) can be retrieved from the Instana REST API using the path
  `/api/infrastructure-monitoring/catalog/plugins`.
* `query` - Optional - The dynamic filter query for which the rule should be applied to
* `enabled` - Optional - Boolean flag if the rule should be enabled - default = true
* `triggering` - Optional - Boolean flag if the rule should trigger an incident - default = false
* `expiration_time` - Optional - The grace period in milliseconds until the issue is closed
* `rules` - Required - The configuration of the specific rule of the custom event [Details](#rules)

### Rules

Exactly one of the elements below must be configured:

* `entity_verifiation` - Optional - configuration of entity verification rules [Details](#entity-verification-rule)
* `system` - Optional - configuration of system rules [Details](#system-rule)
* `threshold` - Optional - configuration of threshold rule [Details](#threshold-rule)

#### Entity Verification Rule

* `severity` - Required - The severity of the rule - allowed values: `warning`, `critical`
* `matching_entity_type` - Required - The entity type used to check for matching entities on the selected hosts.
  Supported entity types (plugins) can be retrieved from the Instana REST API using the path
  `/api/infrastructure-monitoring/catalog/plugins`.
* `matching_operator` - Required - The comparison operator used to check for matching entities on the selected hosts.
  Allowed values: `is`, `contains`, `startsWith`, `endsWith`
* `matching_entity_label` - Required - The label/string to check for matching entities on the selected hosts
* `offline_duration` - Required - The duration in milliseconds to wait until the entity is considered as offline

#### System Rule

* `severity` - Required - The severity of the rule - allowed values: `warning`, `critical`
* `system_rule_id` - Required - The id of the instana system rule of the given even

#### Threshold Rule

* `severity` - Required - The severity of the rule - allowed values: `warning`, `critical`
* `metric_name` - Optional (required for Built-In and Custom Metrics only) The name of the built in or custom metric
  name (supported
  built in metrics can be retrieved from the REST API using the
  endpoint `/api/infrastructure-monitoring/catalog/metrics/{plugin}`) - conflicts with `metric_pattern`, exactly one of
  them must be provided depending on the use case.
* `metric_pattern` - Optional  (required for Dynamic Built-In Metrics only) the metric pattern of the dynamic built in
  metric - [Details](#metric-pattern) - conflicts with `metric_name`, exactly one of them must be provided depending on the use case.
* `window` - Required - The time window in milliseconds within the rule condition is applied to
* `rollup` - Optional - The resolution of the monitored metrics
* `aggregation` - Optional (depending on metric type) - the aggregation used to calculate the metric value for the given
  time window and/or rollup. Supported value: `sum`, `avg`, `min`, `max`
* `condition_operator` - Required - The condition operator used to check against the calculated metric value for the
  given time window and/or rollup. Supported values: `=`, `!=`, `<=`,`<`, `>`, `=>`
* `condition_value` - Required - The numeric condition value used to check against the calculated metric value for the
  given time window and/or rollup.

##### Metric Pattern
* `prefix` - Required - the prefix of the built-in dynamic metric
* `postfix` - Optional - the postfix of the built-in dynamic metric
* `placeholder` - Required - the placeholder string of the dynamic metric
* `operator` - Required - the operation used to check for matching
  placeholder string. Allowed values:  `is`, `contains`, `any`, `startsWith`, `endsWith`

## Import

Custom event specifications with entity verification rule can be imported using the `id`, e.g.:

```
$ terraform import instana_custom_event_spec_entity_verification_rule.my_event_spec 60845e4e5e6b9cf8fc2868da
```
