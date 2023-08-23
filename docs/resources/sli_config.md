# SLI Configuration

Management of SLI configurations. A service level indicator (SLI) is the defined quantitative measure of one characteristic of the level of service that is provided to a customer. Common examples of such indicators are error rate or response latency of a service.

API Documentation: <https://instana.github.io/openapi/#operation/createSli>

The ID of the resource which is also used as unique identifier in Instana is auto generated!

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
* `metric_configuration` - Required - resource block to describe the metric the SLI config is based on
  * `metric_name` - Required - name of the metric
  * `aggregation` - Required - the aggregation type for the metric configuration
  Allowed values: `SUM`, `MEAN`, `MAX`, `MIN`, `P25`, `P50`, `P75`, `P90`, `P95`, `P98`, `P99`, `DISTINCT_COUNT`
  * `threshold` - Required - threshold for the metric configuration
* `sli_entity` - Required - resource block to describe the entity the SLI config is based on
  * `type` - Required - the entity type
  Allowed values: `application`, `custom`, `availability`
  * `application_id` - Optional - the application ID of the entity
  * `service_id` - Optional - the service ID of the entity
  * `endpoint_id` - Optional - the endpoint ID of the entity
  * `boundary_scope` - Required - the boundary scope of the entity
  Allowed values: `ALL`, `INBOUND`

## Import

SLI Configs can be imported using the `id`, e.g.:

```
$ terraform import instana_sli_config.my_sli 60845e4e5e6b9cf8fc2868da
```
