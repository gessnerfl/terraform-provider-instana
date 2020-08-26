# Terraform Provider Instana

[![Build Status](https://travis-ci.org/gessnerfl/terraform-provider-instana.svg?branch=master)](https://travis-ci.org/gessnerfl/terraform-provider-instana)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=de.gessnerfl.terraform-provider-instana&metric=alert_status)](https://sonarcloud.io/dashboard?id=de.gessnerfl.terraform-provider-instana)

Terraform provider implementation for the Instana REST API.

Changes Log: **[CHANGELOG.md](https://github.com/gessnerfl/terraform-provider-instana/blob/master/CHANGELOG.md)**

- [Terraform Provider Instana](#terraform-provider-instana)
  - [How to Use](#how-to-use)
    - [Provider Configuration](#provider-configuration)
    - [Resources](#resources)
      - [Application Settings](#application-settings)
        - [Application Configuration](#application-configuration)
      - [Event Settings](#event-settings)
        - [Custom Event Specification](#custom-event-specification)
          - [Custom Event Specification with System Rules](#custom-event-specification-with-system-rules)
          - [Custom Event Specification with Entity Verification Rules](#custom-event-specification-with-entity-verification-rules)
          - [Custom Event Specification with Threshold Rules](#custom-event-specification-with-threshold-rules)
        - [Alerting Channels](#alerting-channels)
          - [Email](#email)
          - [Google Chat](#google-chat)
          - [Office 365](#office-365)
          - [OpsGenie](#opsgenie)
          - [Pager Duty](#pager-duty)
          - [Slack](#slack)
          - [Splunk](#splunk)
          - [VictorOps](#victorops)
          - [Generic Webhook](#generic-webhook)
        - [Alerting Configuration](#alerting-configuration)
      - [Settings](#settings)
        - [User Roles](#user-roles)
  - [Implementation Details](#implementation-details)
    - [Testing](#testing)
    - [Release a new version](#release-a-new-version)

## How to Use

**NOTE:** Starting with version 0.6.0 Terraform version 0.12.x or later is required.

The implementation is based on the Instana REST API. The configuration reflects one by one the REST API of Instana.
Because of this the semantics of the configuration options is not described in this documentation. Instead of this
a link to the official API documentation will be provided to avoid that the documentation of this implementation
diverge from the documentation of the official API.

### Provider Configuration

with the provider configuration the basic requirements for the Instana REST API has to be provided. There are only
two configuration options needed to setup the Instana Terraform Provider:

- api_token: The API token which is created in the Settings area of Instana for remote access through the REST API. You have to make sure that you assign the proper permissions for this token to configure the desired resources with this provider. E.g. when User Roles should be provisioned by terraform using this provider implementation then the permission 'Access role configuration' must be activated
- endpoint: The endpoint of the instana backend. For SaaS the endpoint URL has the pattern _tenant_-_organization_.instana.io. For onPremise installation the endpoint URL depends on your local setup.
- default_name_prefix: optional - default "" - string which should be added in front the resource UI name or label by default (not supported by all resources). For existing resources the string will only be added when the name/label is changed.
- default_name_suffix: optional - default " (TF managed)" - string which should be appended to the resource UI name or label by default (not supported by all resources). For existing resources the string will only be appended when the name/label is changed.

```hcl
provider "instana" {
  api_token = "secure-api-token"  
  endpoint = "<tenant>-<org>.instana.io"
  default_name_prefix = ""
  default_name_suffix = "(TF managed)"
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
The resource supports `default_name_prefix` and `default_name_suffix` and will append the string automatically
to the application config label when active.

```hcl
resource "instana_application_config" "example" {
  label               = "label"
  scope               = "INCLUDE_ALL_DOWNSTREAM"  #Optional, default = INCLUDE_NO_DOWNSTREAM
  boundary_scope      = "INBOUND"  #Optional, default = INBOUND
  match_specification = "agent.tag.stage EQUALS 'test' OR aws.ec2.tag.stage EQUALS 'test' OR call.tag.stage EQUALS 'test'"
}
```

For **scope** the following three options are allowed:

- INCLUDE_ALL_DOWNSTREAM
- INCLUDE_NO_DOWNSTREAM
- INCLUDE_IMMEDIATE_DOWNSTREAM_DATABASE_AND_MESSAGING

For **boundary_scope** the following three options are allowed:

- INBOUND
- ALL
- DEFAULT

The **match_specification** defines which entities should be included into the application. It supports:

- logical AND and/or logical OR conjunctions whereas AND has higher precedence then OR
- comparisons EQUALS, NOT_EQUAL, CONTAINS, NOT_CONTAIN
- unary operators IS_EMPTY, NOT_EMPTY, IS_BLANK, NOT_BLANK.

The **match_specification** is defined by the following eBNF:

```plain
match_specification       := logical_or
binary_operation          := logical_and OR logical_or | logical_and
logical_and               := primary_expression AND logical_and | primary_expression
primary_expression        := comparison | unary_operator_expression
comparison                := key comparison_operator value
comparison_operator       := EQUALS | NOT_EQUAL | CONTAINS | NOT_CONTAIN | STARTS_WITH | ENDS_WITH | NOT_STARTS_WITH | NOT_ENDS_WITH | GREATER_OR_EQUAL_THAN | LESS_OR_EQUAL_THAN | LESS_THAN | GREATER_THAN
unary_operator_expression := key unary_operator
unary_operator            := IS_EMPTY | NOT_EMPTY | IS_BLANK | NOT_BLANK
key                       := [a-zA-Z][\.a-zA-Z0-9_\-]*
value                     := "'" <string> "'"

```

#### Event Settings

API Documentation: <https://instana.github.io/openapi/#tag/Event-Settings>

##### Custom Event Specification

Management of custom event specifications

API Documentation: <https://instana.github.io/openapi/#operation/getCustomEventSpecification>

Custom Event Specification support two different flavors:

- System Rules - defines an event triggered by a system rule
- Entity Verification Rules - defines an event which is triggered when an entity is not running on the selected systems.
- Threshold Rules - defines an event triggered by a rule for a certain metric comparing the value with a given value over a time window

Custom event resources supports `default_name_prefix` and `default_name_suffix`. The string will be appended automatically
to the name of the custom event.

###### Custom Event Specification with System Rules

```hcl
resource "instana_custom_event_spec_system_rule" "example" {
  name            = "name"
  query           = "query"        #Optional
  enabled         = true           #Optional, default = true
  triggering      = true           #Optional, default = false
  description     = "description"  #Optional
  expiration_time = 60000          #Optional, only when triggering is active

  rule_severity       = "warning"
  rule_system_rule_id = "system-rule-id"
}
```

###### Custom Event Specification with Entity Verification Rules

Entity verification rules is a specialized system rule to check for hosts which do not have matching entities running on them.

```hcl
resource "instana_custom_event_spec_entity_verification_rule" "example" {
  name            = "name"
  query           = "query"        #Optional
  enabled         = true           #Optional, default = true
  triggering      = true           #Optional, default = false
  description     = "description"  #Optional
  expiration_time = 60000          #Optional, only when triggering is active

  rule_severity              = "warning"
  rule_matching_entity_type  = "process"
  rule_matching_operator     = "is"             #allowed values: is, contains, startsWith, starts_with, endsWith, ends_with
  rule_matching_entity_label = "entity-label"
  rule_offline_duration      = 60000
}
```

###### Custom Event Specification with Threshold Rules

A threshold rule is verifies if a certain condition applies to a given metric. Therefore you can either use `rule_rollup` or `rule_window` or
both to define the data points which should be evaluated. Instana API always returns max. 600 data points for validation.

- `rule_window` = the time frame in seconds where the aggregation is applied to
- `rule_rollup` = the resolution of the data points which are considered for this event (See also <https://instana.github.io/openapi/#tag/Infrastructure-Metrics>)

Both are optional in the Instana API. Usually configurations define a **window** for calculating the event.

```hcl
resource "instana_custom_event_spec_threshold_rule" "example" {
  name            = "name"
  query           = "query"        #Optional
  enabled         = true           #Optional, default = true
  triggering      = true           #Optional, default = false
  description     = "description"  #Optional
  expiration_time = 60000          #Optional, only when triggering is active
  entity_type     = "entity_type"

  rule_severity           = "warning"
  rule_metric_name        = "metric_name"
  rule_window             = 60000          #Optional
  rule_rollup             = 500            #Optional
  rule_aggregation        = "sum"          #Optional depending on metric type, allowed values: sum, avg, min, max
  rule_condition_operator = "=="           #allowed values: ==, !=, <=, <, >, =>
  rule_condition_value    = 1.2

  #For built-in dynamic metrics
  rule_metric_pattern_prefix      = "prefix"        #Optional - required only for built in dynamic metrics
  rule_metric_pattern_postfix     = "postfix"       #Optional
  rule_metric_pattern_placeholder = "placeholder"   #Optional
  rule_metric_pattern_operator    = "is"            #Optional - required only for built in dynamic metrics; allowed values: is, contains, any, startsWith, endsWith
}
```

##### Alerting Channels

Management of Alerting channels in Instana. A dedicated terraform resource type is available for each altering channel type.

API Documentation: <https://instana.github.io/openapi/#operation/getAlertingChannels>

###### Email

Alerting channel configuration for notifications to a specified list of email addresses.

```hcl
resource "instana_alerting_channel_email" "example" {
  name = "my-email-alerting-channel"
  emails = [ "email1@example.com", "email2@example.com" ]
}
```

###### Google Chat

Alerting channel configuration for notifications to Google Chat.

```hcl
resource "instana_alerting_channel_google_chat" "example" {
  name        = "my-google-chat-alerting-channel"
  webhook_url = "https://my.google.chat.weebhook.exmaple.com/"
}
```

###### Office 365

Alerting channel configuration for notifications to Office356.

```hcl
resource "instana_alerting_channel_office_365" "example" {
  name        = "my-office365-alerting-channel"
  webhook_url = "https://my.office365.weebhook.exmaple.com/"
}
```

###### OpsGenie

Alerting channel configuration integration with OpsGenie.

```hcl
resource "instana_alerting_channel_ops_genie" "example" {
  name = "my-ops-genie-alerting-channel"
  api_key = "my-secure-api-key"
  tags = [ "tag1", "tag2" ]
  region = "EU"
}
```

###### Pager Duty

Alerting channel configuration integration with PagerDuty.

```hcl
resource "instana_alerting_channel_pager_duty" "example" {
  name = "my-pager-duty-alerting-channel"
  service_integration_key = "my-service-integration-key"
}
```

###### Slack

Alerting channel configuration notifications to Slack.

```hcl
resource "instana_alerting_channel_slack" "example" {
  name        = "my-slack-alerting-channel"
  webhook_url = "https://my.slack.weebhook.exmaple.com/"
  icon_url    = "https://my.slack.icon.exmaple.com/"   #Optional
  channel     = "my-channel"                           #Optional
}
```

###### Splunk

Alerting channel configuration to integrate with Splunk.

```hcl
resource "instana_alerting_channel_splunk" "example" {
  name  = "my-splunk-alerting-channel"
  url   = "https://my.splunk.url.example.com"
  token = "my-splunk-token"
}
```

###### VictorOps

Alerting channel configuration to integrate with VictorOps.

```hcl
resource "instana_alerting_channel_victor_ops" "example" {
  name        = "my-victor-ops-alerting-channel"
  api_key     = "my-victor-ops-api-key"
  routing_key = "my-victor-ops-routing-key"
}
```

###### Generic Webhook

Alerting channel configuration to integrate with WebHooks.

```hcl
resource "instana_alerting_channel_webhook" "example" {
  name         = "my-generic-webhook-alerting-channel"
  webhook_urls = [ "https://my.weebhook1.exmaple.com/", "https://my.weebhook2.exmaple.com/" ]
  http_headers = {      #Optional
    header1 = "headerValue1"
    header2 = "headerValue2"
  }
}
```

##### Alerting Configuration

Management of alert configurations. Alert configurations define how either event types or 
event (aka rules) are reported to integrated services (Alerting Channels).

API Documentation: <https://instana.github.io/openapi/#operation/putAlert>

The ID of the resource which is also used as unique identifier in Instana is auto generated!

Alerting configurations support `default_name_prefix` and `default_name_suffix`. The string will be appended automatically
to the alert_name.

Configuration for an alert configuration using dedicated event/rule ids

```hcl
resource "instana_alerting_config" "example" {
  alert_name            = "name"
  integration_ids       = [ "alerting-channel-id1", "alerting-channel-id2" ]  # Optional, you can also use references to existing alerting channel configurations
  event_filter_query    = "query"                                             # Optional
  event_filter_rule_ids = [ "rule-1", "rule-2" ]                              # You can also use references to existing custom events
}
``` 

Configuration for an alert configuration using event types

```hcl
resource "instana_alerting_config" "example" {
  alert_name               = "name"
  integration_ids          = [ "alerting-channel-id1", "alerting-channel-id2" ]  # Optional, you can also use references to existing alerting channel configurations
  event_filter_query       = "query"                                             # Optional
  event_filter_event_types = [ "incident", "critical" ]                          # Allowed values: incident, critical, warning, change, online, offline, agent_monitoring_issue, none
}
``` 

#### Settings

##### User Roles

Management of user roles.
API Documentation: <https://instana.github.io/openapi/#operation/getRole>

The ID of the resource which is also used as unique identifier in Instana is auto generated!
The resource does NOT support `default_name_prefix` and `default_name_suffix`.

```hcl
resource "instana_user_role" "example" {
  name                                   = "name"
  implicit_view_filter                   = "view filter" #Optional
  can_configure_service_mapping          = true          #Optional, default = false
  can_configure_eum_applications         = true          #Optional, default = false
  can_configure_users                    = true          #Optional, default = false
  can_install_new_agents                 = true          #Optional, default = false
  can_see_usage_information              = true          #Optional, default = false
  can_configure_integrations             = true          #Optional, default = false
  can_see_on_premise_license_information = true          #Optional, default = false
  can_configure_roles                    = true          #Optional, default = false
  can_configure_custom_alerts            = true          #Optional, default = false
  can_configure_api_tokens               = true          #Optional, default = false
  can_configure_agent_run_mode           = true          #Optional, default = false
  can_view_audit_log                     = true          #Optional, default = false
  can_configure_objectives               = true          #Optional, default = false
  can_configure_agents                   = true          #Optional, default = false
  can_configure_authentication_methods   = true          #Optional, default = false
  can_configure_applications             = true          #Optional, default = false
}
```

## Implementation Details

### Testing

 Mocking:
 Tests are co-located in the package next to the implementation. We use gomock (<https://github.com/golang/mock)> for mocking. To generate mocks you need to use the package options to create the mocks in the same package:

```bash
mockgen -source=<source_file> -destination=mocks/<source_package>/<source_file_name>_mocks.go package=<source_package>_mocks -self_package=github.com/gessnerfl/terraform-provider-instana/<source_package>
```

### Release a new version

1. Create a new tag follow semantic versioning approach
2. Update changelog before creating a new release by using [github-changelog-generator](https://github.com/github-changelog-generator/github-changelog-generator)
3. Push the tag to the remote to build the new release
