# terraform-provider-instana

[![Build Status](https://travis-ci.org/gessnerfl/terraform-provider-instana.svg?branch=master)](https://travis-ci.org/gessnerfl/terraform-provider-instana)
[![Sonarcloud Status](https://sonarcloud.io/api/project_badges/measure?project=de.gessnerfl.terraform-provider-instana&metric=alert_status)](https://sonarcloud.io/dashboard/index/terraform-provider-instana)

Terraform provider implementation for Instana REST API

#How to Use

## Provider Configuration

```
provider "instana" {
  api_token = "secure-api-token"  
  endpoint = "<mytenant>-<myorg>.instana.io"
}
```

## Resources

### Event Settings

#### Rules

```
resource "instana_rule" "example" {
  name = "name"
  entity_type = "entity_type"
  metric_name = "metric_name"
  rollup = 100
  window = 20000
  aggregation = "sum"
  condition_operator = ">"
  condition_value = 1.1
}
```

#### Rule Bindings

```
resource "instana_rule_binding" "example" {
  enabled = true
  triggering = true
  severity = 5
  text = "text"
  description = "description"
  expiration_time = 60000
  query = "query"
  rule_ids = [ "rule-id-1", "rule-id-2" ]
}
```

### Settings

#### User Roles

```
resource "instana_user_role" "example" {
  name = "name"
  implicit_view_filter = "view filter"
  can_configure_service_mapping = true
  can_configure_eum_applications = true
  can_configure_users = true
  can_install_new_agents = true
  can_see_usage_information = true
  can_configure_integrations = true
  can_see_on_premise_license_information = true
  can_configure_roles = true
  can_configure_custom_alerts = true
  can_configure_api_tokens = true
  can_configure_agent_run_mode = true
  can_view_audit_log = true
  can_configure_objectives = true
  can_configure_agents = true
  can_configure_authentication_methods = true
  can_configure_applications = true
}
```

# Implementation
 Mocking:
 Tests are colocated in the package next to the implementation. We use gomock (https://github.com/golang/mock) for mocking. To generate mocks you need to use the package options to create the mocks in the same package:

```
mockgen -source=<source_file> -destination=mocks/<source_package>/<source_file_name>_mocks.go package=<source_package>_mocks -self_package=github.com/gessnerfl/terraform-provider-instana/<source_package>
```