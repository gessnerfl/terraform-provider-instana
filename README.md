# Terraform Provider Instana

[![Build Status](https://travis-ci.org/gessnerfl/terraform-provider-instana.svg?branch=master)](https://travis-ci.org/gessnerfl/terraform-provider-instana)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=de.gessnerfl.terraform-provider-instana&metric=alert_status)](https://sonarcloud.io/dashboard?id=de.gessnerfl.terraform-provider-instana)

Terraform provider implementation for the Instana REST API.

- [Terraform Provider Instana](#terraform-provider-instana)
  - [How to Use](#how-to-use)
    - [Provider Configuration](#provider-configuration)
    - [Resources](#resources)
      - [Application Settings](#application-settings)
        - [Application Configuration](#application-configuration)
      - [Event Settings](#event-settings)
        - [Rules](#rules)
        - [Rule Bindings](#rule-bindings)
      - [Settings](#settings)
        - [User Roles](#user-roles)
  - [Implementation Details](#implementation-details)

## How to Use

The implementation is based on the Instana REST API. The configuration reflects one by one the REST API of Instana.
Because of this the semantics of the configuration options is not described in this documentation. Instead of this
a link to the official API documentation will be provided to avoid that the documentation of this implementation
diverge from the documentation of the official API.

### Provider Configuration

with the provider configuration the basic requirements for the Instana REST API has to be provided. There are only
two configuration options needed to setup the Instana Terraform Provider:

- api_token: The API token which is created in the Settings area of Instana for remote access through the REST API. You have to make sure that you assign the proper permissions for this token to configure the desired resources with this provider. E.g. when User Roles should be provisioned by terraform using this provider implementation then the permission 'Access role configuration' must be activated
- endpoint: The endpoint of the instana backend. For SaaS the endpoint URL has the pattern _tenant_-_organization_.instana.io. For onPremise installation the endpoint URL depends on your local setup.

```hcl
provider "instana" {
  api_token = "secure-api-token"  
  endpoint = "<tenant>-<org>.instana.io"
}
```

### Resources

In this section we will list all provided endpoints with the full list of available configuration options. Not all
resources of the Instana API are implemented by the terraform-provider-instana. Please open a ticket of provide a
Pull Request when a resource or a configuration option is missing.

#### Application Settings

API Documentation: <https://instana.github.io/openapi/#tag/Application-Settings>

##### Application Configuration

Management of application configurations (definition of application perspectives).
API Documentation: <https://instana.github.io/openapi/#operation/putApplicationConfig>

The ID of the resource which is also used as unique identifier in Instana is auto generated!

```hcl
resource "instana_rule" "example" {
  label = "label"
  scope = "INCLUDE_ALL_DOWNSTREAM"
  match_specification = "agent.tag.stage = 'test' OR aws.ec2.tag.stage = 'test' OR call.tag.stage = 'test'"
}
```

For **scope** the following three options are allowed:

* INCLUDE_ALL_DOWNSTREAM
* INCLUDE_NO_DOWNSTREAM
* INCLUDE_IMMEDIATE_DOWNSTREAM_DATABASE_AND_MESSAGING

The **match_specification** defines which entities should be included into the application. It supports:

* logical AND and/or logical OR conjunctions whereas AND has higher precedence then OR
* comparisons EQUALS, NOT_EQUAL, CONTAINS, NOT_CONTAIN
* unary operators IS_EMPTY, NOT_EMPTY, IS_BLANK, NOT_BLANK.

The **match_specification** is defined by the following eBNF:

```plain
match_specification       := logical_or
binary_operation          := logical_and OR logical_or | logical_and
logical_and               := primary_expression AND logical_and | primary_expression
primary_expression        := comparison | unary_operator_expression
comparison                := key comparison_operator value
comparison_operator       := EQUALS | NOT_EQUAL | CONTAINS | NOT_CONTAIN
unary_operator_expression := key unary_operator
unary_operator            := IS_EMPTY | NOT_EMPTY | IS_BLANK | NOT_BLANK
key                       := [a-zA-Z][\.a-zA-Z0-9]*
value                     := "'" <string> "'"

```


#### Event Settings

API Documentation: <https://instana.github.io/openapi/#tag/Event-Settings>

##### Rules

Management of custom rules.
API Documentation: <https://instana.github.io/openapi/#operation/getRule>

The ID of the resource which is also used as unique identifier in Instana is auto generated!

```hcl
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

##### Rule Bindings

Management of Rule Bindings. Rule bindings represent incident configurations.
API Documentation: <https://instana.github.io/openapi/#operation/getRuleBinding>

The ID of the resource which is also used as unique identifier in Instana is auto generated!

```hcl
resource "instana_rule_binding" "example" {
  enabled = true
  triggering = true
  severity = warning
  text = "text"
  description = "description"
  expiration_time = 60000
  query = "query"
  rule_ids = [ "rule-id-1", "rule-id-2" ]
}
```

#### Settings

##### User Roles

Management of user roles.
API Documentation: <https://instana.github.io/openapi/#operation/getRole>

The ID of the resource which is also used as unique identifier in Instana is auto generated!

```hcl
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

## Implementation Details

 Mocking:
 Tests are co-located in the package next to the implementation. We use gomock (<https://github.com/golang/mock)> for mocking. To generate mocks you need to use the package options to create the mocks in the same package:

```hcl
mockgen -source=<source_file> -destination=mocks/<source_package>/<source_file_name>_mocks.go package=<source_package>_mocks -self_package=github.com/gessnerfl/terraform-provider-instana/<source_package>
```