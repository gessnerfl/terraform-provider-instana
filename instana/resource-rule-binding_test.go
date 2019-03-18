package instana_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	mocks "github.com/gessnerfl/terraform-provider-instana/mocks"
	testutils "github.com/gessnerfl/terraform-provider-instana/test-utils"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
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

func TestShouldSuccessfullyReadRuleBindingFromInstanaAPIWhenBaseDataIsReturned(t *testing.T) {
	expectedModel := createBaseTestRuleBindingModel()
	testShouldSuccessfullyReadRuleBindingFromInstanaAPI(expectedModel, t)
}

func TestShouldSuccessfullyReadRuleBindingFromInstanaAPIWhenFullDataIsReturned(t *testing.T) {
	expectedModel := createFullTestRuleBindingModel()
	testShouldSuccessfullyReadRuleBindingFromInstanaAPI(expectedModel, t)
}

func TestShouldSuccessfullyReadRuleBindingFromInstanaAPIWhenBaseDataWithQueryIsReturned(t *testing.T) {
	expectedModel := createTestRuleBindingModelWithQuery()
	testShouldSuccessfullyReadRuleBindingFromInstanaAPI(expectedModel, t)
}

func testShouldSuccessfullyReadRuleBindingFromInstanaAPI(expectedModel restapi.RuleBinding, t *testing.T) {
	resourceData := createEmptyRuleBindingResourceData(t)
	ruleBindingID := "rule-binding-id"
	resourceData.SetId(ruleBindingID)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRuleBindingApi := mocks.NewMockRuleBindingResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().RuleBindings().Return(mockRuleBindingApi).Times(1)
	mockRuleBindingApi.EXPECT().GetOne(gomock.Eq(ruleBindingID)).Return(expectedModel, nil).Times(1)

	err := ReadRuleBinding(resourceData, mockInstanaAPI)

	if err != nil {
		t.Fatalf("Expected no error to be returned, %s", err)
	}
	verifyModelAppliedToResource(expectedModel, resourceData, t)
}

func TestShouldFailToReadRuleBindingFromInstanaAPIWhenIDIsMissing(t *testing.T) {
	resourceData := createEmptyRuleBindingResourceData(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	err := ReadRuleBinding(resourceData, mockInstanaAPI)

	if err == nil || !strings.HasPrefix(err.Error(), "ID of rule binding") {
		t.Fatal("Expected error to occur because of missing id")
	}
}

func TestShouldFailToReadRuleBindingFromInstanaAPIAndDeleteResourceWhenBindingDoesNotExist(t *testing.T) {
	resourceData := createEmptyRuleBindingResourceData(t)
	ruleBindingID := "rule-binding-id"
	resourceData.SetId(ruleBindingID)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRuleBindingApi := mocks.NewMockRuleBindingResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().RuleBindings().Return(mockRuleBindingApi).Times(1)
	mockRuleBindingApi.EXPECT().GetOne(gomock.Eq(ruleBindingID)).Return(restapi.RuleBinding{}, restapi.ErrEntityNotFound).Times(1)

	err := ReadRuleBinding(resourceData, mockInstanaAPI)

	if err != nil {
		t.Fatalf("Expected no error to be returned, %s", err)
	}
	if len(resourceData.Id()) > 0 {
		t.Fatal("Expected ID to be cleaned to destroy resource")
	}
}

func TestShouldCreateRuleBindingThroughInstanaAPI(t *testing.T) {
	data := createFullTestRuleBindingData()
	resourceData := createRuleBindingResourceData(t, data)
	expectedModel := createFullTestRuleBindingModel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRuleBindingApi := mocks.NewMockRuleBindingResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().RuleBindings().Return(mockRuleBindingApi).Times(1)
	mockRuleBindingApi.EXPECT().Upsert(gomock.AssignableToTypeOf(restapi.RuleBinding{})).Return(expectedModel, nil).Times(1)

	err := CreateRuleBinding(resourceData, mockInstanaAPI)

	if err != nil {
		t.Fatalf("Expected no error to be returned, %s", err)
	}
	verifyModelAppliedToResource(expectedModel, resourceData, t)
}

func TestShouldReturnErrorWhenCreateRuleBindingFailsThroughInstanaAPI(t *testing.T) {
	data := createFullTestRuleBindingData()
	resourceData := createRuleBindingResourceData(t, data)
	expectedError := errors.New("test")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRuleBindingApi := mocks.NewMockRuleBindingResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().RuleBindings().Return(mockRuleBindingApi).Times(1)
	mockRuleBindingApi.EXPECT().Upsert(gomock.AssignableToTypeOf(restapi.RuleBinding{})).Return(restapi.RuleBinding{}, expectedError).Times(1)

	err := CreateRuleBinding(resourceData, mockInstanaAPI)

	if err == nil || expectedError != err {
		t.Fatal("Expected definned error to be returned")
	}
}

func TestShouldDeleteRuleBindingThroughInstanaAPI(t *testing.T) {
	id := "test-id"
	data := createFullTestRuleBindingData()
	resourceData := createRuleBindingResourceData(t, data)
	resourceData.SetId(id)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRuleBindingApi := mocks.NewMockRuleBindingResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().RuleBindings().Return(mockRuleBindingApi).Times(1)
	mockRuleBindingApi.EXPECT().DeleteByID(gomock.Eq(id)).Return(nil).Times(1)

	err := DeleteRuleBinding(resourceData, mockInstanaAPI)

	if err != nil {
		t.Fatalf("Expected no error to be returned, %s", err)
	}
	if len(resourceData.Id()) > 0 {
		t.Fatal("Expected ID to be cleaned to destroy resource")
	}
}

func TestShouldReturnErrorWhenDeleteRuleBindingFailsThroughInstanaAPI(t *testing.T) {
	id := "test-id"
	data := createFullTestRuleBindingData()
	resourceData := createRuleBindingResourceData(t, data)
	resourceData.SetId(id)
	expectedError := errors.New("test")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRuleBindingApi := mocks.NewMockRuleBindingResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().RuleBindings().Return(mockRuleBindingApi).Times(1)
	mockRuleBindingApi.EXPECT().DeleteByID(gomock.Eq(id)).Return(expectedError).Times(1)

	err := DeleteRuleBinding(resourceData, mockInstanaAPI)

	if err == nil || err != expectedError {
		t.Fatal("Expected error to be returned")
	}
	if len(resourceData.Id()) == 0 {
		t.Fatal("Expected ID not to be cleaned to avoid resource is destroy")
	}
}

func verifyModelAppliedToResource(model restapi.RuleBinding, resourceData *schema.ResourceData, t *testing.T) {
	if model.ID != resourceData.Id() {
		t.Fatal("Expected ID to be identical")
	}
	if model.Enabled != resourceData.Get(RuleBindingFieldEnabled).(bool) {
		t.Fatal("Expected Enabled to be identical")
	}
	if model.Triggering != resourceData.Get(RuleBindingFieldTriggering).(bool) {
		t.Fatal("Expected Triggering to be identical")
	}
	if model.Severity != resourceData.Get(RuleBindingFieldSeverity).(int) {
		t.Fatal("Expected Severity to be identical")
	}
	if model.Text != resourceData.Get(RuleBindingFieldText).(string) {
		t.Fatal("Expected Text to be identical")
	}
	if model.Description != resourceData.Get(RuleBindingFieldDescription).(string) {
		t.Fatal("Expected Description to be identical")
	}
	if model.ExpirationTime != resourceData.Get(RuleBindingFieldExpirationTime).(int) {
		t.Fatal("Expected ExpirationTime to be identical")
	}
	if model.Query != resourceData.Get(RuleBindingFieldQuery).(string) {
		t.Fatal("Expected Query to be identical")
	}
	if !cmp.Equal(model.RuleIds, ReadStringArrayParameterFromResource(resourceData, RuleBindingFieldRuleIds)) {
		t.Fatal("Expected RuleIds to be identical")
	}
}

func createFullTestRuleBindingModel() restapi.RuleBinding {
	data := createBaseTestRuleBindingModel()
	data.Enabled = true
	data.Triggering = true
	data.Query = "query"
	return data
}

func createTestRuleBindingModelWithQuery() restapi.RuleBinding {
	data := createBaseTestRuleBindingModel()
	data.Query = "query"
	return data
}

func createBaseTestRuleBindingModel() restapi.RuleBinding {
	return restapi.RuleBinding{
		ID:             "id",
		Severity:       5,
		Text:           "text",
		Description:    "description",
		ExpirationTime: 1234,
		RuleIds:        []string{"test-rule-id-1", "test-rule-id-2"},
	}
}

func createFullTestRuleBindingData() map[string]interface{} {
	data := make(map[string]interface{})
	data[RuleBindingFieldEnabled] = true
	data[RuleBindingFieldTriggering] = true
	data[RuleBindingFieldText] = "text"
	data[RuleBindingFieldDescription] = "description"
	data[RuleBindingFieldExpirationTime] = 1234
	data[RuleBindingFieldQuery] = "query"
	data[RuleBindingFieldRuleIds] = []string{"test-rule-id-1", "test-rule-id-2"}
	return data
}
