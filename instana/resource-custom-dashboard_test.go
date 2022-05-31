package instana_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"
	"net/http"
	"strings"
	"testing"
)

func TestCustomDashboardResource(t *testing.T) {
	terraformResourceInstanceName := ResourceInstanaCustomDashboard + ".example"
	inst := &customDashboardResourceTest{
		terraformResourceInstanceName: terraformResourceInstanceName,
		resourceHandle:                NewCustomDashboardResourceHandle(),
	}
	inst.run(t)
}

type customDashboardResourceTest struct {
	terraformResourceInstanceName string
	resourceHandle                ResourceHandle
}

func (test *customDashboardResourceTest) run(t *testing.T) {
	t.Run(fmt.Sprintf("CRUD integration test of %s", ResourceInstanaCustomDashboard), test.createIntegrationTest())
	t.Run(fmt.Sprintf("%s should have schema version zero", ResourceInstanaCustomDashboard), test.createTestResourceShouldHaveSchemaVersionZero())
	t.Run(fmt.Sprintf("%s should have no state upgrader", ResourceInstanaCustomDashboard), test.createTestResourceShouldHaveNoStateUpgrader())
	t.Run(fmt.Sprintf("%s should have correct resouce name", ResourceInstanaCustomDashboard), test.createTestResourceShouldHaveCorrectResourceName())
	t.Run(fmt.Sprintf("%s should successfully update state from model", ResourceInstanaCustomDashboard), test.createTestShouldSuccessfullyUpdateTerraformStateFromModel())
	t.Run(fmt.Sprintf("%s should successfully map state to model", ResourceInstanaCustomDashboard), test.createTestShouldSuccessfullyMapTerraformStateFromModel())
	t.Run(fmt.Sprintf("%s should successfully map state to model when no access rule is defined", ResourceInstanaCustomDashboard), test.createTestShouldSuccessfullyMapTerraformStateFromModelWhenNoAccessRuleIsDefined())
}

const customDashboardWidgetsJson = `[
    {
      "id": "6jK0w8KmdHtABCs3",
      "title": "Latency",
      "width": 4,
      "height": 13,
      "x": 4,
      "y": 26,
      "type": "chart",
      "config": {
        "y1": {
          "formatter": "millis.detailed",
          "renderer": "line",
          "metrics": [
            {
              "metric": "latency",
              "timeShift": 0,
              "tagFilters": [
                {
                  "stringValue": "my-app",
                  "name": "application.name",
                  "entity": "DESTINATION",
                  "operator": "EQUALS"
                },
                {
                  "name": "call.inbound_of_application",
                  "entity": "NOT_APPLICABLE",
                  "operator": "NOT_EMPTY"
                }
              ],
              "aggregation": "MEAN",
              "label": "Mean Latency",
              "source": "APPLICATION"
            },
            {
              "metric": "latency",
              "timeShift": 0,
              "tagFilters": [
                {
                  "stringValue": "my-app",
                  "name": "application.name",
                  "entity": "DESTINATION",
                  "operator": "EQUALS"
                },
                {
                  "name": "call.inbound_of_application",
                  "entity": "NOT_APPLICABLE",
                  "operator": "NOT_EMPTY"
                }
              ],
              "aggregation": "P99",
              "label": "99th latency",
              "source": "APPLICATION"
            }
          ]
        },
        "y2": {
          "formatter": "number.detailed",
          "renderer": "line",
          "metrics": []
        },
        "type": "TIME_SERIES"
      }
    }
  ]`

const customDashboardResponseJson = `
{
  "id": "%s",
  "title": "prefix name %d suffix",
  "accessRules": [
    {
      "accessType": "READ_WRITE",
      "relationType": "USER",
      "relatedId": "user-id-1"
    },
    {
      "accessType": "READ",
      "relationType": "GLOBAL",
      "relatedId": ""
    },
    {
      "accessType": "READ_WRITE",
      "relationType": "USER",
      "relatedId": "user-id-2"
    }
  ],
  "widgets": __WIDGETS__,
  "writable": false
}
`

const customDashboardResourceTemplate = `
resource "instana_custom_dashboard" "example" {
  title = "name %d"

  access_rule { 
	access_type = "READ_WRITE"
	relation_type = "USER"
	related_id = "user-id-1"
  }
  
  access_rule { 
	access_type = "READ"
	relation_type = "GLOBAL"
  }

  access_rule { 
	access_type = "READ_WRITE"
	relation_type = "USER"
	related_id = "user-id-2"
  }

  widgets = "__WIDGETS__"
}
`

func (test *customDashboardResourceTest) createIntegrationTest() func(t *testing.T) {
	return func(t *testing.T) {
		serverResponseTemplate := strings.ReplaceAll(customDashboardResponseJson, "__WIDGETS__", customDashboardWidgetsJson)

		id := RandomID()
		resourceInstanceRestAPIPath := restapi.CustomDashboardsResourcePath + "/{internal-id}"
		testutils.DeactivateTLSServerCertificateVerification()
		httpServer := testutils.NewTestHTTPServer()
		httpServer.AddRoute(http.MethodPost, restapi.CustomDashboardsResourcePath, func(w http.ResponseWriter, r *http.Request) {
			dashboard := &restapi.CustomDashboard{}
			err := json.NewDecoder(r.Body).Decode(dashboard)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				r.Write(bytes.NewBufferString("Failed to get request"))
			} else {
				dashboard.ID = id
				w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(dashboard)
			}
		})
		httpServer.AddRoute(http.MethodPut, resourceInstanceRestAPIPath, testutils.EchoHandlerFunc)
		httpServer.AddRoute(http.MethodDelete, resourceInstanceRestAPIPath, testutils.EchoHandlerFunc)
		httpServer.AddRoute(http.MethodGet, resourceInstanceRestAPIPath, func(w http.ResponseWriter, r *http.Request) {
			modCount := httpServer.GetCallCount(http.MethodPut, restapi.CustomDashboardsResourcePath+"/"+id)
			json := fmt.Sprintf(serverResponseTemplate, id, modCount)
			w.Header().Set(contentType, r.Header.Get(contentType))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(json))
		})
		httpServer.Start()
		defer httpServer.Close()

		resource.UnitTest(t, resource.TestCase{
			ProviderFactories: testProviderFactory,
			Steps: []resource.TestStep{
				test.createIntegrationTestStep(httpServer.GetPort(), 0, id),
				testStepImportWithCustomID(test.terraformResourceInstanceName, id),
				test.createIntegrationTestStep(httpServer.GetPort(), 1, id),
				testStepImportWithCustomID(test.terraformResourceInstanceName, id),
			},
		})
	}
}

func (test *customDashboardResourceTest) createIntegrationTestStep(httpPort int, iteration int, id string) resource.TestStep {
	widgetsDefinition := utils.RemoveNewLinesAndTabs(customDashboardWidgetsJson)
	resourceConfig := fmt.Sprintf(strings.ReplaceAll(customDashboardResourceTemplate, "__WIDGETS__", strings.ReplaceAll(widgetsDefinition, "\"", "\\\"")), iteration)
	normalizedWidgetsDefinition := NormalizeJSONString(widgetsDefinition)
	return resource.TestStep{
		Config: appendProviderConfig(resourceConfig, httpPort),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, "id", id),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, CustomDashboardFieldTitle, formatResourceName(iteration)),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, CustomDashboardFieldFullTitle, formatResourceFullName(iteration)),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, CustomDashboardFieldWidgets, normalizedWidgetsDefinition),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, CustomDashboardFieldAccessRule+".0."+CustomDashboardFieldAccessRuleAccessType, "READ_WRITE"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, CustomDashboardFieldAccessRule+".0."+CustomDashboardFieldAccessRuleRelatedID, "user-id-1"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, CustomDashboardFieldAccessRule+".0."+CustomDashboardFieldAccessRuleRelationType, "USER"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, CustomDashboardFieldAccessRule+".1."+CustomDashboardFieldAccessRuleAccessType, "READ"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, CustomDashboardFieldAccessRule+".1."+CustomDashboardFieldAccessRuleRelatedID, ""),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, CustomDashboardFieldAccessRule+".1."+CustomDashboardFieldAccessRuleRelationType, "GLOBAL"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, CustomDashboardFieldAccessRule+".2."+CustomDashboardFieldAccessRuleAccessType, "READ_WRITE"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, CustomDashboardFieldAccessRule+".2."+CustomDashboardFieldAccessRuleRelatedID, "user-id-2"),
			resource.TestCheckResourceAttr(test.terraformResourceInstanceName, CustomDashboardFieldAccessRule+".2."+CustomDashboardFieldAccessRuleRelationType, "USER"),
		),
	}
}

func (test *customDashboardResourceTest) createTestResourceShouldHaveSchemaVersionZero() func(t *testing.T) {
	return func(t *testing.T) {
		require.Equal(t, 0, test.resourceHandle.MetaData().SchemaVersion)
	}
}

func (test *customDashboardResourceTest) createTestResourceShouldHaveNoStateUpgrader() func(t *testing.T) {
	return func(t *testing.T) {
		require.Empty(t, test.resourceHandle.StateUpgraders())
	}
}

func (test *customDashboardResourceTest) createTestResourceShouldHaveCorrectResourceName() func(t *testing.T) {
	return func(t *testing.T) {
		require.Equal(t, test.resourceHandle.MetaData().ResourceName, "instana_custom_dashboard")
	}
}

func (test *customDashboardResourceTest) createTestShouldSuccessfullyUpdateTerraformStateFromModel() func(t *testing.T) {
	return func(t *testing.T) {
		userID := "user-id"
		dashboard := restapi.CustomDashboard{
			ID:      "dashboard-id",
			Title:   "prefix dashboard-title suffix",
			Widgets: "dashboard-widgets",
			AccessRules: []restapi.AccessRule{
				{AccessType: restapi.AccessTypeReadWrite, RelationType: restapi.RelationTypeUser, RelatedID: &userID},
				{AccessType: restapi.AccessTypeRead, RelationType: restapi.RelationTypeGlobal},
			},
		}

		testHelper := NewTestHelper(t)
		sut := test.resourceHandle
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

		err := sut.UpdateState(resourceData, &dashboard, testHelper.ResourceFormatter())

		require.NoError(t, err)
		require.Equal(t, "dashboard-id", resourceData.Id())
		require.Equal(t, "prefix dashboard-title suffix", resourceData.Get(CustomDashboardFieldFullTitle).(string))
		require.Equal(t, "dashboard-title", resourceData.Get(CustomDashboardFieldTitle).(string))
		require.Equal(t, "dashboard-widgets", resourceData.Get(CustomDashboardFieldWidgets).(string))
		require.Len(t, resourceData.Get(CustomDashboardFieldAccessRule).([]interface{}), 2)
		require.Equal(t, []interface{}{
			map[string]interface{}{
				CustomDashboardFieldAccessRuleAccessType:   "READ_WRITE",
				CustomDashboardFieldAccessRuleRelatedID:    "user-id",
				CustomDashboardFieldAccessRuleRelationType: "USER",
			},
			map[string]interface{}{
				CustomDashboardFieldAccessRuleAccessType:   "READ",
				CustomDashboardFieldAccessRuleRelatedID:    "",
				CustomDashboardFieldAccessRuleRelationType: "GLOBAL",
			},
		}, resourceData.Get(CustomDashboardFieldAccessRule).([]interface{}))
	}
}

func (test *customDashboardResourceTest) createTestShouldSuccessfullyMapTerraformStateFromModel() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper(t)
		sut := test.resourceHandle
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

		userID := "user-id"
		resourceData.SetId("dashboard-id")
		resourceData.Set(CustomDashboardFieldTitle, "dashboard-title")
		resourceData.Set(CustomDashboardFieldFullTitle, "prefix dashboard-title suffix")
		resourceData.Set(CustomDashboardFieldAccessRule, []interface{}{
			map[string]interface{}{
				CustomDashboardFieldAccessRuleAccessType:   "READ_WRITE",
				CustomDashboardFieldAccessRuleRelatedID:    userID,
				CustomDashboardFieldAccessRuleRelationType: "USER",
			},
			map[string]interface{}{
				CustomDashboardFieldAccessRuleAccessType:   "READ",
				CustomDashboardFieldAccessRuleRelationType: "GLOBAL",
			},
		})
		resourceData.Set(CustomDashboardFieldWidgets, "dashboard-widgets")

		result, err := sut.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

		require.NoError(t, err)
		require.Equal(t, &restapi.CustomDashboard{
			ID:      "dashboard-id",
			Title:   "prefix dashboard-title suffix",
			Widgets: json.RawMessage("dashboard-widgets"),
			AccessRules: []restapi.AccessRule{
				{AccessType: restapi.AccessTypeReadWrite, RelatedID: &userID, RelationType: restapi.RelationTypeUser},
				{AccessType: restapi.AccessTypeRead, RelationType: restapi.RelationTypeGlobal},
			},
		}, result)
	}

}

func (test *customDashboardResourceTest) createTestShouldSuccessfullyMapTerraformStateFromModelWhenNoAccessRuleIsDefined() func(t *testing.T) {
	return func(t *testing.T) {
		testHelper := NewTestHelper(t)
		sut := test.resourceHandle
		resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

		resourceData.SetId("dashboard-id")
		resourceData.Set(CustomDashboardFieldTitle, "dashboard-title")
		resourceData.Set(CustomDashboardFieldFullTitle, "prefix dashboard-title suffix")
		resourceData.Set(CustomDashboardFieldWidgets, "dashboard-widgets")

		result, err := sut.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

		require.NoError(t, err)
		require.Equal(t, &restapi.CustomDashboard{
			ID:          "dashboard-id",
			Title:       "prefix dashboard-title suffix",
			Widgets:     json.RawMessage("dashboard-widgets"),
			AccessRules: []restapi.AccessRule{},
		}, result)
	}

}
