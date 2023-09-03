# Custom Event Specification with Entity Verification Rule Resource

Configuration of a custom event specification based on an entity verification rule. This rule type is used
to check for hosts which do not have matching entities running on them.

API Documentation: <https://instana.github.io/openapi/#operation/putCustomEventSpecification>

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

  rule {  
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

  rule { 
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

  rule { 
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
* `query` - Optional - The dynamic filter query for which the rule should be applied to
* `enabled` - Optional - Boolean flag if the rule should be enabled - default = true
* `triggering` - Optional - Boolean flag if the rule should trigger an incident - default = false
* `expiration_time` - Optional - The grace period in milliseconds until the issue is closed
* `rule_severity` - Required - The severity of the rule - allowed values: `warning`, `critical`
* `rule_matching_entity_type` - Required - The entity type used to check for matching entities on the selected hosts. 
Supported entity types (plugins) can be retrieved from the Instana REST API using the path
`/api/infrastructure-monitoring/catalog/plugins`.
* `rule_matching_operator` - Required - The comparison operator used to check for matching entities on the selected hosts. 
Allowed values: `is`, `contains`, `startsWith`, `starts_with`, `endsWith`, `ends_with`
* `rule_matching_entity_label` - Required - The label/string to check for matching entities on the selected hosts
* `rule_offline_duration` - Required - The duration in milliseconds to wait until the entity is considered as offline

## Import

Custom event specifications with entity verification rule can be imported using the `id`, e.g.:

```
$ terraform import instana_custom_event_spec_entity_verification_rule.my_event_spec 60845e4e5e6b9cf8fc2868da
```
