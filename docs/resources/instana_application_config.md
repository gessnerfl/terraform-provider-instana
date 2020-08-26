# Application Configuration Resource

Management of application configurations (definition of application perspectives).

API Documentation: <https://instana.github.io/openapi/#operation/putApplicationConfig>

The ID of the resource which is also used as unique identifier in Instana is auto generated!
The resource supports `default_name_prefix` and `default_name_suffix` and will append the string automatically
to the application config label when active.

## Example Usage

```hcl
resource "instana_application_config" "example" {
  label               = "label"
  scope               = "INCLUDE_ALL_DOWNSTREAM"  #Optional, default = INCLUDE_NO_DOWNSTREAM
  boundary_scope      = "INBOUND"  #Optional, default = INBOUND
  match_specification = "agent.tag.stage EQUALS 'test' OR aws.ec2.tag.stage EQUALS 'test' OR call.tag.stage EQUALS 'test'"
}
```

## Argument Reference

* `label` - Required - The name/label of the application perspective
* `scope` - Optional - The scope of the application perspective. Default value: `INCLUDE_NO_DOWNSTREAM`. Allowed valued: `INCLUDE_ALL_DOWNSTREAM`, `INCLUDE_NO_DOWNSTREAM`, `INCLUDE_IMMEDIATE_DOWNSTREAM_DATABASE_AND_MESSAGING`
* `boundary_scope` - Optional - The boundary scope of the application perspective. Default value `DEFAULT`. Allowed values: `INBOUND`, `ALL`, `DEFAULT`
* `match_specification` - Required - specifies which entities should be included in the application

### Match Specification
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
comparison_operator       := EQUALS | NOT_EQUAL | CONTAINS | NOT_CONTAIN | STARTS_WITH | ENDS_WITH | NOT_STARTS_WITH | NOT_ENDS_WITH | GREATER_OR_EQUAL_THAN | LESS_OR_EQUAL_THAN | LESS_THAN | GREATER_THAN
unary_operator_expression := key unary_operator
unary_operator            := IS_EMPTY | NOT_EMPTY | IS_BLANK | NOT_BLANK
key                       := [a-zA-Z][\.a-zA-Z0-9_\-]*
value                     := "'" <string> "'"

```