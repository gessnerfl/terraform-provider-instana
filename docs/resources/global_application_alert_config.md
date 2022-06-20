# Global Application Alert Configuration Resource

Management of global application alert configurations (Global Application Smart Alerts).

API Documentation: <https://instana.github.io/openapi/#tag/Global-Application-Alert-Configuration>

The ID of the resource which is also used as unique identifier in Instana is auto generated!
The resource supports `default_name_prefix` and `default_name_suffix`. The configured strings
will be appended automatically to the name of the alert when activated.

## Example Usage

```hcl
resource "instana_global_application_alert_config" "example" {
	name              = "test-alert"
    description       = "test-alert-description"
    boundary_scope    = "ALL"
    severity          = "warning"
    triggering        = false
    include_internal  = false
    include_synthetic = false
    alert_channel_ids = [ instana_alerting_channel_email.example.id ]
    granularity       = 600000
	evaluation_type   = "PER_AP"

	tag_filter        = "call.type@na EQUALS 'HTTP'"
    
    application {
		application_id = instana_application_config.example.id
		inclusive 	   = true
	}

	rule {
		slowness {
			metric_name = "latency"
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

* `name` - Required - The name for the global application alert configuration
* `description` - Required - The description text of the global application alert config
* `severity` - Required - The severity of the alert when triggered (`critical` or `warning`)
* `boundary_scope` - Required - The boundary scope of the global application alert config. Allowed values: `INBOUND`, `ALL`, `DEFAULT`
* `triggering` - Optional - default `false` - Flag to indicate whether also an Incident is triggered or not. The default is false
* `include_internal` - Optional - default `false` - Flag to indicate whether also internal calls are included in the scope or not
* `include_synthetic` - Optional - default `false` - Flag to indicate whether also synthetic calls are included in the scope or not
* `alert_channel_ids` - Optional - List of IDs of alert channels defined in Instana.
* `granularity` - Optional - default `600000` - The evaluation granularity used for detection of violations of the defined threshold. In other words, it defines the size of the tumbling window used. Allowed values: `300000`, `600000`, `900000`, `1200000`, `800000`
* `evaluation_type` - Required - The evaluation type of the global application alert config. Allowed values: `PER_AP`, `PER_AP_SERVICE`, `PER_AP_ENDPOINT`
* `tag_filter` - Optional - The tag filter of the global application alert config. [Details](#tag-filter-argument-reference)
* `application` - Required - Selection/Set of applications in scope. [Details](#application-argument-reference)
* `rule` - Required - Indicates the type of rule this alert configuration is about. [Details](#rule-argument-reference)
* `custom_payload_filed` - Optional - An optional list of custom payload fields (static key/value pairs added to the event).  [Details](#custom-payload-field-argument-reference)
* `threshold` - Required - Indicates the type of threshold this alert rule is evaluated on.  [Details](#threshold-argument-reference)
* `time_threshold` - Required - Indicates the type of violation of the defined threshold.  [Details](#time-threshold-argument-reference)

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

### Application Argument Reference

* `application_id` - Required - ID of the included application
* `inclusive` - Required - Defines whether this node and his child nodes are included (true) or excluded (false)
* `service` - Optional - Selection of services in scope. [Details](#service-argument-reference)

#### Service Argument Reference

* `service_id` - Required - ID of the included service
* `inclusive` - Required - Defines whether this node and his child nodes are included (true) or excluded (false)
* `endpoint` - Optional - Selection of endpoints in scope. [Details](#endpoint-argument-reference)

##### Endpoint Argument Reference

* `endpoint_id` - Required - ID of the included endpoint
* `inclusive` - Required - Defines whether this node and his child nodes are included (true) or excluded (false)

### Rule Argument Reference

Exactly one of the elements below must be configured

* `error_rate` - Optional - Rule based on the error rate of the configured alert configuration target. [Details](#error-rate-rule-argument-reference)
* `logs` - Optional - Rule based on logs of the configured alert configuration target. [Details](#logs-rule-argument-reference)
* `slowness` - Optional - Rule based on the slowness of the configured alert configuration target. [Details](#slowness-rule-argument-reference)
* `status_code` - Optional - Rule based on the HTTP status code of the configured alert configuration target. [Details](#status-code-rule-argument-reference)
* `throughput` - Optional - Rule based on the throughput of the configured alert configuration target. [Details](#throughput-rule-argument-reference)

#### Error Rate Rule Argument Reference

* `metric_name` - Required - The metric name of the application alert rule
* `aggregation` - Required - The aggregation function of the application alert rule. Supported values `SUM`, `MEAN`, `MAX`, `MIN`, `P25`, `P50`, `P75`, `P90`, `P95`, `P98`, `P99`, `P99_9`, `P99_99`, `DISTRIBUTION`, `DISTINCT_COUNT`, `SUM_POSITIVE`
* `stable_hast` - Optional - The stable hash used for the application alert rule

#### Logs Rule Argument Reference

* `metric_name` - Required - The metric name of the application alert rule
* `aggregation` - Required - The aggregation function of the application alert rule. Supported values `SUM`, `MEAN`, `MAX`, `MIN`, `P25`, `P50`, `P75`, `P90`, `P95`, `P98`, `P99`, `P99_9`, `P99_99`, `DISTRIBUTION`, `DISTINCT_COUNT`, `SUM_POSITIVE`
* `stable_hast` - Optional - The stable hash used for the application alert rule
* `level` - Required - The log level for which this rule applies to. Supported values: `WARN`, `ERROR`, `ANY`
* `message` - Optional - The log message for which this rule applies to.
* `operator` - Required - The operator which will be applied to evaluate this rule. Supported values: `EQUALS`, `NOT_EQUAL`, `CONTAINS`, `NOT_CONTAIN`, `IS_EMPTY`, `NOT_EMPTY`, `IS_BLANK`, `IS_BLANK`, `NOT_BLANK`, `STARTS_WITH`, `ENDS_WITH`, `NOT_STARTS_WITH`, `NOT_ENDS_WITH`, `GREATER_OR_EQUAL_THAN`, `LESS_OR_EQUAL_THAN`, `GREATER_THAN`, `LESS_THAN`

#### Slowness Rule Argument Reference

* `metric_name` - Required - The metric name of the application alert rule
* `aggregation` - Required - The aggregation function of the application alert rule. Supported values `SUM`, `MEAN`, `MAX`, `MIN`, `P25`, `P50`, `P75`, `P90`, `P95`, `P98`, `P99`, `P99_9`, `P99_99`, `DISTRIBUTION`, `DISTINCT_COUNT`, `SUM_POSITIVE`
* `stable_hast` - Optional - The stable hash used for the application alert rule

#### Status Code Rule Argument Reference

* `metric_name` - Required - The metric name of the application alert rule
* `aggregation` - Required - The aggregation function of the application alert rule. Supported values `SUM`, `MEAN`, `MAX`, `MIN`, `P25`, `P50`, `P75`, `P90`, `P95`, `P98`, `P99`, `P99_9`, `P99_99`, `DISTRIBUTION`, `DISTINCT_COUNT`, `SUM_POSITIVE`
* `stable_hast` - Optional - The stable hash used for the application alert rule
* `status_code_start` - Optional - minimal HTTP status code applied for this rule
* `status_code_end` - Optional - maximum HTTP status code applied for this rule

#### Throughput Rule Argument Reference

* `metric_name` - Required - The metric name of the application alert rule
* `aggregation` - Required - The aggregation function of the application alert rule. Supported values `SUM`, `MEAN`, `MAX`, `MIN`, `P25`, `P50`, `P75`, `P90`, `P95`, `P98`, `P99`, `P99_9`, `P99_99`, `DISTRIBUTION`, `DISTINCT_COUNT`, `SUM_POSITIVE`
* `stable_hast` - Optional - The stable hash used for the application alert rule

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

* `request_impact` - Optional - Time threshold base on request impact. [Details](#request-impact-time-threshold-argument-reference)
* `violations_in_period` - Optional - Time threshold base on violations in period. [Details](#violations-in-period-time-threshold-argument-reference)
* `violations_in_sequence` - Optional - Time threshold base on violations in sequence. [Details](#violations-in-sequence-time-threshold-argument-reference)

#### Request Impact Time Threshold Argument Reference

* `time_window` - Optional - The time window if the time threshold
* `request` - Optional - The number of requests in the given window

#### Violations In Period Time Threshold Argument Reference

* `time_window` - Optional - The time window if the time threshold
* `violations` - Optional - The violations appeared in the period

#### Violations In Sequence Time Threshold Argument Reference

* `time_window` - Optional - The time window if the time threshold

## Import

Application Alert Configs can be imported using the `id`, e.g.:

```
$ terraform import instana_application_alert_config.example 60845e4e5e6b9cf8fc2868da
```