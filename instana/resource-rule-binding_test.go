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

var testRuleBindingProviders = map[string]terraform.ResourceProvider{
	"instana": Provider(),
}

const resourceRuleBindingDefinition = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:8080"
}

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
`

func TestCRUDOfRuleBindingResourceWithMockServer(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, "/api/ruleBindings/{id}", testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, "/api/ruleBindings/{id}", testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, "/api/ruleBindings/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		json := strings.ReplaceAll(`
		{
			"id" : "{{id}}",
			"enabled" : true,
			"triggering" : true,
			"severity" : 5,
			"text" : "text",
			"description" : "description",
			"expirationTime" : 60000,
			"query" : "query",
			"ruleIds" : [ "rule-id-1", "rule-id-2" ]
		}
		`, "{{id}}", vars["id"])
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resource.UnitTest(t, resource.TestCase{
		Providers: testRuleBindingProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: resourceRuleBindingDefinition,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("instana_rule_binding.example", "id"),
					resource.TestCheckResourceAttr("instana_rule_binding.example", RuleBindingFieldEnabled, "true"),
					resource.TestCheckResourceAttr("instana_rule_binding.example", RuleBindingFieldTriggering, "true"),
					resource.TestCheckResourceAttr("instana_rule_binding.example", RuleBindingFieldSeverity, "5"),
					resource.TestCheckResourceAttr("instana_rule_binding.example", RuleBindingFieldText, "text"),
					resource.TestCheckResourceAttr("instana_rule_binding.example", RuleBindingFieldDescription, "description"),
					resource.TestCheckResourceAttr("instana_rule_binding.example", RuleBindingFieldExpirationTime, "60000"),
					resource.TestCheckResourceAttr("instana_rule_binding.example", RuleBindingFieldQuery, "query"),
					resource.TestCheckResourceAttr("instana_rule_binding.example", RuleBindingFieldRuleIds+".0", "rule-id-1"),
					resource.TestCheckResourceAttr("instana_rule_binding.example", RuleBindingFieldRuleIds+".1", "rule-id-2"),
				),
			},
		},
	})
}

func TestResourceRuleBindingDefinition(t *testing.T) {
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

func validateRuleBindingResourceSchema(schemaMap map[string]*schema.Schema, t *testing.T) {
	validateSchemaOfTypeBoolWithDefault(RuleBindingFieldEnabled, true, schemaMap, t)
	validateSchemaOfTypeBoolWithDefault(RuleBindingFieldTriggering, false, schemaMap, t)
	validateOptionalSchemaOfTypeInt(RuleBindingFieldSeverity, schemaMap, t)
	validateRequiredSchemaOfTypeString(RuleBindingFieldText, schemaMap, t)
	validateRequiredSchemaOfTypeString(RuleBindingFieldDescription, schemaMap, t)
	validateRequiredSchemaOfTypeInt(RuleBindingFieldExpirationTime, schemaMap, t)
	validateOptionalSchemaOfTypeString(RuleBindingFieldQuery, schemaMap, t)
	validateRequiredSchemaOfTypeListOfString(RuleBindingFieldRuleIds, schemaMap, t)
}
