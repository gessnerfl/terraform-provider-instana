# Synthetic Location Data Source

Data source to get the synthetic locations from Instana API. This allows you to retrieve the specification
by label and location type and reference it in other resources such as Synthetic tests.

API Documentation: <https://instana.github.io/openapi/#operation/getSyntheticLocation>

## Example Usage

```hcl
data "instana_synthetic_location" "locations" {}
```

## Argument Reference

* `label` - Required - the label of the synthetic location
* `location_type` - Required - indicates if the location is public or private
