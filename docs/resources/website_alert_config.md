# Website Alert Configuration Resource

Management of website alert configurations (Website Smart Alerts).

API Documentation: <https://instana.github.io/openapi/#operation/findActiveWebsiteAlertConfigs>

The ID of the resource which is also used as unique identifier in Instana is auto generated!
The resource supports `default_name_prefix` and `default_name_suffix`. The configured strings
will be appended automatically to the name of the alert when activated.

## Example Usage

```hcl
resource "instana_website_alert_config" "example" {
  name              = "test-alert"
  description       = "test-alert-description"
  severity          = "warning"
  triggering        = false
  alert_channel_ids = [instana_alerting_channel_email.example.id]
  granularity       = 600000
  tag_filter        = "beacon.user.id@na EQUALS '5678'"
  website_id        = instana_website_monitoring_config.example.id

  rule {
    slowness {
      metric_name = "onLoadTime"
      aggregation = "P90"
    }
  }
  threshold {
    static {
      operator = ">="
      value    = 5.0
    }
  }
  time_threshold {
    violations_in_sequence {
      time_window = 600000
    }
  }
  custom_payload_field {
    key   = "test"
    value = "test123"
  }
}
```

## Argument Reference

* `name` - Required - The name for the application alert configuration
* `description` - Required - The description text of the application alert config
* `severity` - Required - The severity of the alert when triggered (`critical` or `warning`)
* `triggering` - Optional - default `false` - Flag to indicate whether also an Incident is triggered or not. The default is false
* `alert_channel_ids` - Optional - List of IDs of alert channels defined in Instana.
* `granularity` - Optional - default `600000` - The evaluation granularity used for detection of violations of the defined threshold. In other words, it defines the size of the tumbling window used. Allowed values: `300000`, `600000`, `900000`, `1200000`, `800000`
* `tag_filter` - Optional - The tag filter of the application alert config. [Details](#tag-filter-argument-reference)
* `rule` - Required - Indicates the type of rule this alert configuration is about. [Details](#rule-argument-reference)
* `custom_payload_filed` - Optional - An optional list of custom payload fields (static key/value pairs added to the event).  [Details](#custom-payload-field-argument-reference)
* `threshold` - Required - Indicates the type of threshold this alert rule is evaluated on.  [Details](#threshold-argument-reference)
* `time_threshold` - Required - Indicates the type of violation of the defined threshold.  [Details](#time-threshold-argument-reference)
* `website_id` - Required - Unique ID of the website

### Tag Filter Argument Reference
The **tag_filter** defines which entities should be included into the application. It supports:

* logical AND and/or logical OR conjunctions whereas AND has higher precedence then OR
* comparison operators EQUALS, NOT_EQUAL, CONTAINS | NOT_CONTAIN, STARTS_WITH, ENDS_WITH, NOT_STARTS_WITH, NOT_ENDS_WITH, GREATER_OR_EQUAL_THAN, LESS_OR_EQUAL_THAN, LESS_THAN, GREATER_THAN
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
tag_key                   := identifier
entity_origin             := src | dest | na
value                     := string_value | number_value | boolean_value
string_value              := "'" <string> "'"
number_value              := (+-)?[0-9]+
boolean_value             := TRUE | FALSE
identifier                := [a-zA-Z_][\.a-zA-Z0-9_\-/]*
```

### Rule Argument Reference

Exactly one of the elements below must be configured

* `specific_js_error` - Optional - Rule based on a specific javascript error of the configured alert configuration target. [Details](#specific-js-error-rule-argument-reference)
* `slowness` - Optional - Rule based on the slowness of the configured alert configuration target. [Details](#slowness-rule-argument-reference)
* `status_code` - Optional - Rule based on the HTTP status code of the configured alert configuration target. [Details](#status-code-rule-argument-reference)
* `throughput` - Optional - Rule based on the throughput of the configured alert configuration target. [Details](#throughput-rule-argument-reference)

#### Specific JS Error Rule Argument Reference

* `metric_name` - Required - The metric name of the application alert rule
* `aggregation` - Optional - The aggregation function of the application alert rule. Supported values `SUM`, `MEAN`, `MAX`, `MIN`, `P25`, `P50`, `P75`, `P90`, `P95`, `P98`, `P99`, `P99_9`, `P99_99`, `DISTRIBUTION`, `DISTINCT_COUNT`, `SUM_POSITIVE`
* `operator`    - Required - The operator which will be applied to evaluate this rule. Supported values: `EQUALS`, `NOT_EQUAL`, `CONTAINS`, `NOT_CONTAIN`, `IS_EMPTY`, `NOT_EMPTY`, `IS_BLANK`, `IS_BLANK`, `NOT_BLANK`, `STARTS_WITH`, `ENDS_WITH`, `NOT_STARTS_WITH`, `NOT_ENDS_WITH`, `GREATER_OR_EQUAL_THAN`, `LESS_OR_EQUAL_THAN`, `GREATER_THAN`, `LESS_THAN`
* `value`       - Required - The value identify the specific javascript error.

#### Slowness Rule Argument Reference

* `metric_name` - Required - The metric name of the application alert rule
* `aggregation` - Required - The aggregation function of the application alert rule. Supported values `SUM`, `MEAN`, `MAX`, `MIN`, `P25`, `P50`, `P75`, `P90`, `P95`, `P98`, `P99`, `P99_9`, `P99_99`, `DISTRIBUTION`, `DISTINCT_COUNT`, `SUM_POSITIVE`

#### Status Code Rule Argument Reference

* `metric_name` - Required - The metric name of the application alert rule
* `aggregation` - Required - The aggregation function of the application alert rule. Supported values `SUM`, `MEAN`, `MAX`, `MIN`, `P25`, `P50`, `P75`, `P90`, `P95`, `P98`, `P99`, `P99_9`, `P99_99`, `DISTRIBUTION`, `DISTINCT_COUNT`, `SUM_POSITIVE`
* `operator`    - Required - The operator which will be applied to evaluate this rule. Supported values: `EQUALS`, `NOT_EQUAL`, `CONTAINS`, `NOT_CONTAIN`, `IS_EMPTY`, `NOT_EMPTY`, `IS_BLANK`, `IS_BLANK`, `NOT_BLANK`, `STARTS_WITH`, `ENDS_WITH`, `NOT_STARTS_WITH`, `NOT_ENDS_WITH`, `GREATER_OR_EQUAL_THAN`, `LESS_OR_EQUAL_THAN`, `GREATER_THAN`, `LESS_THAN`
* `value`       - Required - The value identify the specific http status code.

#### Throughput Rule Argument Reference

* `metric_name` - Required - The metric name of the application alert rule
* `aggregation` - Optional - The aggregation function of the application alert rule. Supported values `SUM`, `MEAN`, `MAX`, `MIN`, `P25`, `P50`, `P75`, `P90`, `P95`, `P98`, `P99`, `P99_9`, `P99_99`, `DISTRIBUTION`, `DISTINCT_COUNT`, `SUM_POSITIVE`

### Custom Payload Field Argument Reference

* `key` - Required - The key of the custom payload field
* `value` - Required - The value of the custom payload field

### Threshold Argument Reference

Exactly one of the elements below must be configured

* `historic_baseline` - Optional - Threshold based on a historic baseline. [Details](#historic-baseline-threshold-argument-reference)
* `static` - Optional - Static threshold definition. [Details](#static-threshold-argument-reference)

#### Historic Baseline Threshold Argument Reference

* `operator` - Required - The operator which will be applied to evaluate the threshold. Supported values: `>`, `>=`, `<`, `<=`
* `last_updated` - Optional - The last updated value of the threshold
* `baseline` - Optional - The baseline of the historic baseline threshold
* `deviation_factor` - Optional - The baseline of the historic baseline threshold
* `seasonality` - Required - The seasonality of the historic baseline threshold. Supported values: `WEEKLY`, `DAILY`

#### Static Threshold Argument Reference

* `operator` - Required - The operator which will be applied to evaluate the threshold. Supported values: `>`, `>=`, `<`, `<=`
* `last_updated` - Optional - The last updated value of the threshold
* `value` - Optional - The value of the static threshold

### Time Threshold Argument Reference

Exactly one of the elements below must be configured

* `user_impact_of_violations_in_sequence` - Optional - Time threshold base on user impact of violations in sequence. [Details](#user-impact-of-violations-in-sequence-time-threshold-argument-reference)
* `violations_in_period` - Optional - Time threshold base on violations in period. [Details](#violations-in-period-time-threshold-argument-reference)
* `violations_in_sequence` - Optional - Time threshold base on violations in sequence. [Details](#violations-in-sequence-time-threshold-argument-reference)

#### User Impact Of Violations in Sequence Time Threshold Argument Reference

* `time_window` - Optional - The time window if the time threshold
* `impact_measurement_method` - Required - The impact method of the time threshold based on user impact of violations in sequence. Supported valued: `AGGREGATED`, `PER_WINDOW`
* `user_percentage` - Optional - The percentage (expressed as floating point number from 0.0 to 1.0) of impacted users of the time threshold based on user impact of violations in sequence
* `users` - Optional - The number of impacted users (> 0) of the time threshold based on user impact of violations in sequence

#### Violations In Period Time Threshold Argument Reference

* `time_window` - Optional - The time window if the time threshold
* `violations` - Optional - The violations appeared in the period

#### Violations In Sequence Time Threshold Argument Reference

* `time_window` - Optional - The time window if the time threshold

## Import

Application Alert Configs can be imported using the `id`, e.g.:

```
$ terraform import instana_website_alert_config.example 60845e4e5e6b9cf8fc2868da
```