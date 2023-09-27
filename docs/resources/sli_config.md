# SLI Configuration

Management of SLI configurations. A service level indicator (SLI) is the defined quantitative measure of one
characteristic of the level of service that is provided to a customer. Common examples of such indicators are error rate
or response latency of a service.

API Documentation: <https://instana.github.io/openapi/#operation/createSli>

The ID of the resource which is also used as unique identifier in Instana is auto generated!

**Note:** SLI Configurations cannot be changed. An update of the resource will result in an error. To update an SLI you
need to create a new SLI and delete the old one.

## Example Usage

```hcl
resource "instana_sli_config" "example" {
    name                         = "sli_name_example"
    initial_evaluation_timestamp = 0
    metric_configuration {
	    metric_name = "metric_name_example"
	    aggregation = "SUM"
	    threshold   = 1
    }
    sli_entity {
        type           = "application"
        application_id = "application_id_example"
        service_id     = "service_id_example"
        endpoint_id    = "endpoint_id_example"
        boundary_scope = "ALL"
    }
}
``` 

## Argument Reference

* `name` - Required - the name of the SLI configuration
* `initial_evaluation_timestamp` - Optional - the initial evaluation timestamp for the SLI config
* `metric_configuration` - Optional - resource block to describe the metric the SLI config is based
  on [Details](#metric-configuration-reference), Required
  for [application_time_based](#application-time-based-sli-entity-reference)
  and [website_time_based](#website-time-based-sli-entity-reference) sli entities
* `sli_entity` - Required - resource block to describe the entity the SLI config is based
  on. [Details](#sli-entity-reference)

### Metric Configuration Reference

* `metric_name` - Required - name of the metric
* `aggregation` - Required - the aggregation type for the metric configuration. Allowed
  values: `SUM`, `MEAN`, `MAX`, `MIN`, `P25`, `P50`, `P75`, `P90`, `P95`, `P98`, `P99`, `P99_9`, `P99_99`, `DISTRIBUTION`, `DISTINCT_COUNT`, `SUM_POSITIVE`, `PER_SECOND`
* `threshold` - Required - threshold for the metric configuration

### Sli Entity Reference

Exactly one of the elements below must be configured:

* `application_event_based` - Optional - event-base sli entity configuration for
  applications [Details](#application-event-based-sli-entity-reference)
* `application_time_based` - Optional - time-base sli entity configuration for
  applications [Details](#application-time-based-sli-entity-reference)
* `website_event_based` - Optional - event-base sli entity configuration for
  websites [Details](#website-event-based-sli-entity-reference)
* `website_time_based` - Optional - time-base sli entity configuration for
  websites [Details](#website-time-based-sli-entity-reference)

#### Application event-based Sli Entity Reference

* `application_id` - Required - the application ID of the entity
* `boundary_scope` - Required - the boundary scope of the entity. Allowed values: `ALL`, `INBOUND`
* `includes_internal` - Optional - flag to indicate whether also internal calls are included in the scope or not. The
  default is `false`
* `includes_synthetic` - Optional - flag to indicate whether also synthetic calls are included in the scope or not. The
  default is `false`
* `good_event_filter_expression` - Required - tag filter expression to match good events /
  calls [Details](#tag-filter-expression-reference)
* `bad_event_filter_expression` - Required - tag filter expression to match bad events /
  calls [Details](#tag-filter-expression-reference)

#### Application time-based Sli Entity Reference

* `application_id` - Required - the application ID of the entity
* `service_id` - Optional - the service ID of the entity
* `endpoint_id` - Optional - the endpoint ID of the entity
* `boundary_scope` - Required - the boundary scope of the entity. Allowed values: `ALL`, `INBOUND`

#### Website event-based Sli Entity Reference

* `website_id` - Required - the website ID of the entity
* `beacon_type` - Required - the beacon type of the entity. Allowed
  values: `pageLoad`, `resourceLoad`, `httpRequest`, `error`, `custom`, `pageChange`
* `good_event_filter_expression` - Required - tag filter expression to match good events /
  calls [Details](#tag-filter-expression-reference)
* `bad_event_filter_expression` - Required - tag filter expression to match bad events /
  calls [Details](#tag-filter-expression-reference)

#### Website time-based Sli Entity Reference

* `website_id` - Required - the website ID of the entity
* `beacon_type` - Required - the beacon type of the entity. Allowed
  values: `pageLoad`, `resourceLoad`, `httpRequest`, `error`, `custom`, `pageChange`
* `filter_expression` - Optional - tag filter expression to match events /
  calls [Details](#tag-filter-expression-reference)

#### Tag Filter Expression Reference

The **tag_filter** defines which calls/events should be included. It supports:

* logical AND and/or logical OR conjunctions whereas AND has higher precedence then OR
* comparison operators EQUALS, NOT_EQUAL, CONTAINS | NOT_CONTAIN, STARTS_WITH, ENDS_WITH, NOT_STARTS_WITH,
  NOT_ENDS_WITH, GREATER_OR_EQUAL_THAN, LESS_OR_EQUAL_THAN, LESS_THAN, GREATER_THAN
* unary operators IS_EMPTY, NOT_EMPTY, IS_BLANK, NOT_BLANK.

The **tag_filter** is defined by the following eBNF:

```plain
tag_filter                := logical_or
logical_or                := logical_and OR logical_or | logical_and
logical_and               := primary_expression AND logical_and | primary_expression
primary_expression        := comparison | unary_operator_expression
comparison                := identifier comparison_operator value | identifier@entity_origin comparison_operator value | identifier:tag_key comparison_operator value | identifier:tag_key@entity_origin comparison_operator value
comparison_operator       := EQUALS | NOT_EQUAL | CONTAINS | NOT_CONTAIN | STARTS_WITH | ENDS_WITH | NOT_STARTS_WITH | NOT_ENDS_WITH | GREATER_OR_EQUAL_THAN | LESS_OR_EQUAL_THAN | LESS_THAN | GREATER_THAN
unary_operator_expression := identifier unary_operator | identifier@entity_origin unary_operator
unary_operator            := IS_EMPTY | NOT_EMPTY | IS_BLANK | NOT_BLANK
tag_key                   := identifier | string_value
entity_origin             := src | dest | na
value                     := string_value | number_value | boolean_value
string_value              := "'" <string> "'"
number_value              := (+-)?[0-9]+
boolean_value             := TRUE | FALSE
identifier                := [a-zA-Z_][\.a-zA-Z0-9_\-/]*
```

## Import

SLI Configs can be imported using the `id`, e.g.:

```
$ terraform import instana_sli_config.my_sli 60845e4e5e6b9cf8fc2868da
```
