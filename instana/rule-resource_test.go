package instana_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	testutils "github.com/gessnerfl/terraform-provider-instana/test-utils"
)

var testProviders = map[string]terraform.ResourceProvider{
	"instana": Provider(),
}

const resourceDefinition = `
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

func TestRender(t *testing.T) {
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
		Providers: testProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: resourceDefinition,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("instana_rule.example", "id"),
					resource.TestCheckResourceAttr("instana_rule.example", FieldName, "name"),
					resource.TestCheckResourceAttr("instana_rule.example", FieldEntityType, "entity_type"),
					resource.TestCheckResourceAttr("instana_rule.example", FieldMetricName, "metric_name"),
					resource.TestCheckResourceAttr("instana_rule.example", FieldRollup, "100"),
					resource.TestCheckResourceAttr("instana_rule.example", FieldWindow, "20000"),
					resource.TestCheckResourceAttr("instana_rule.example", FieldAggregation, "sum"),
					resource.TestCheckResourceAttr("instana_rule.example", FieldConditionOperator, ">"),
					resource.TestCheckResourceAttr("instana_rule.example", FieldConditionValue, "1.1"),
				),
			},
		},
	})
}
