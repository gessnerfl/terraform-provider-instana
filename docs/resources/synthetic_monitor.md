# Synthetic Monitor Resource

Synthetic monitor configuration used to manage synthetic monitors in Instana API.

API Documentation: <https://instana.github.io/openapi/#operation/getSyntheticTests>

## Example Usage


### Create a HTTPAction monitor
```hcl
resource "instana_synthetic_monitor" "uptime_check" {
  label          = "test"
  active         = true
  locations      = [data.instana_synthetic_location.monitor.id]
  test_frequency = 10
  playback_mode  = "Staggered"
  configuration {
    mark_synthetic_call = true
    retries             = 0
    retry_interval      = 1
    synthetic_type      = "HTTPAction"
    timeout             = "3m"
    url                 = "https://example.com"
    operation           = "GET"
  }
  custom_properties = {
    "foo" = "bar"
  }
}
```

### Create a HTTPScript monitor
```hcl
resource "instana_synthetic_monitor" "http_action" {
  label     = "test"
  active    = true
  locations = [data.instana_synthetic_location.monitor.id]
  configuration {
    synthetic_type = "HTTPScript"
    script         = <<EOF
      const assert = require('assert');

      $http.get('https://terraform.io',
      function(err, response, body) {
          if (err) {
              console.error(err);
              return;
          }
          assert.equal(response.statusCode, 200, 'Expected a 200 OK response');
      });
    EOF
  }
}
```

## Argument Reference

* `label` - Required - The name of the synthetic monitor
* `description` - Optional - The name of the synthetic monitor
* `active` - Optional - Enables/disables the synthetic monitor (defaults to true)
* `configuration` - Required - Configuration block
* `custom_properties` - Optional - A map of key/values which are used as tags
* `locations` - Required - A list of strings with location IDs 
* `playback_mode` - Optional - Defines how the Synthetic test should be executed across multiple PoPs (defaults to Simultaneous)
* `test_frequency` - Optional - how often the playback for a synthetic monitor is scheduled (defaults to 15 seconds)

### Configuration

* `mark_synthetic_call` - Optional - flag used to control if HTTP calls will be marked as synthetic calls
* `retries` - Optional - Indicates how many attempts will be allowed to get a successful connection (defaults to 0)
* `retry_interval` - Optional - The time interval between retries in seconds (defaults to 1)
* `synthetic_type` - Required - The type of the Synthetic test (currently supports HTTPAction or HTTPScript)
* `timeout` - Optional - The timeout to be used by the PoP playback engines running the test 
* `url` - Required when synthetic_type is set to HTTPAction - The URL which is being tested
* `operation` - Optional - The HTTP operation (defaults to GET)
* `script` - Required  when synthetic_type is set to HTTPScript - The Javascript content in plain text

## Import

Synthetic monitors can be imported using the `id`, e.g.:

```
$ terraform import instana_synthetic_monitor.http_action cl1g4qrmo26x930s17i2
```