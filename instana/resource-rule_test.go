package instana_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	testutils "github.com/gessnerfl/terraform-provider-instana/test-utils"
)

var testRuleProviders = map[string]terraform.ResourceProvider{
	"instana": Provider(),
}

const resourceRuleDefinition = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:8080"
}

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
`

func TestCRUDOfRuleResourceWithMockServer(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, "/api/rules/{id}", testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, "/api/rules/{id}", testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, "/api/rules/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		json := strings.ReplaceAll(`
		{
			"id" : "{{id}}",
			"name" : "name",
			"entityType" : "entity_type",
			"metricName" : "metric_name",
			"rollup" : 100,
			"window" : 20000,
			"aggregation" : "sum",
			"conditionOperator" : ">",
			"conditionValue" : 1.1
		}
		`, "{{id}}", vars["id"])
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resource.UnitTest(t, resource.TestCase{
		Providers: testRuleProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: resourceRuleDefinition,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("instana_rule.example", "id"),
					resource.TestCheckResourceAttr("instana_rule.example", RuleFieldName, "name"),
					resource.TestCheckResourceAttr("instana_rule.example", RuleFieldEntityType, "entity_type"),
					resource.TestCheckResourceAttr("instana_rule.example", RuleFieldMetricName, "metric_name"),
					resource.TestCheckResourceAttr("instana_rule.example", RuleFieldRollup, "100"),
					resource.TestCheckResourceAttr("instana_rule.example", RuleFieldWindow, "20000"),
					resource.TestCheckResourceAttr("instana_rule.example", RuleFieldAggregation, "sum"),
					resource.TestCheckResourceAttr("instana_rule.example", RuleFieldConditionOperator, ">"),
					resource.TestCheckResourceAttr("instana_rule.example", RuleFieldConditionValue, "1.1"),
				),
			},
		},
	})
}

func TestResourceRuleDefinition(t *testing.T) {
	resource := CreateResourceRule()

	validateRuleResourceSchema(resource.Schema, t)

	if resource.Create == nil {
		t.Fatal("Create function expected")
	}
	if resource.Update == nil {
		t.Fatal("Update function expected")
	}
	if resource.Read == nil {
		t.Fatal("Read function expected")
	}
	if resource.Delete == nil {
		t.Fatal("Delete function expected")
	}
}

func validateRuleResourceSchema(schemaMap map[string]*schema.Schema, t *testing.T) {
	validateRequiredSchemaOfTypeString(RuleFieldName, schemaMap, t)
	validateRequiredSchemaOfTypeString(RuleFieldEntityType, schemaMap, t)
	validateRequiredSchemaOfTypeString(RuleFieldMetricName, schemaMap, t)
	validateOptionalSchemaOfTypeInt(RuleFieldRollup, schemaMap, t)
	validateRequiredSchemaOfTypeInt(RuleFieldWindow, schemaMap, t)
	validateRequiredSchemaOfTypeString(RuleFieldAggregation, schemaMap, t)
	validateRequiredSchemaOfTypeString(RuleFieldConditionOperator, schemaMap, t)
	validateRequiredSchemaOfTypeFloat(RuleFieldConditionValue, schemaMap, t)
}
