package instana_test

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

var testAlertingChannelEmailProviders = map[string]terraform.ResourceProvider{
	"instana": Provider(),
}

const resourceAlertingChannelEmailDefinitionTemplate = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:{{PORT}}"
  default_name_prefix = "prefix"
  default_name_suffix = "suffix"
}

resource "instana_alerting_channel_email" "example" {
  name = "name {{ITERATOR}}"
  emails = [ "EMAIL1", "EMAIL2" ]
}
`

const alertingChannelEmailServerResponseTemplate = `
{
	"id"     : "{{id}}",
	"name"   : "prefix name suffix",
	"kind"   : "EMAIL",
	"emails" : [ "EMAIL1", "EMAIL2" ]
}
`

const alertingChannelEmailApiPath = restapi.AlertingChannelsResourcePath + "/{id}"
const testAlertingChannelEmailDefinition = "instana_alerting_channel_email.example"
const alertingChannelEmailID = "alerting-channel-email-id"

func TestCRUDOfAlertingChannelEmailResourceWithMockServer(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, alertingChannelEmailApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, alertingChannelEmailApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, alertingChannelEmailApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		json := strings.ReplaceAll(alertingChannelEmailServerResponseTemplate, "{{id}}", vars["id"])
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceDefinitionWithoutName := strings.ReplaceAll(resourceAlertingChannelEmailDefinitionTemplate, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))
	resourceDefinitionWithoutName0 := strings.ReplaceAll(resourceDefinitionWithoutName, "{{ITERATOR}}", "0")
	resourceDefinitionWithoutName1 := strings.ReplaceAll(resourceDefinitionWithoutName, "{{ITERATOR}}", "1")

	resource.UnitTest(t, resource.TestCase{
		Providers: testAlertingChannelEmailProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: resourceDefinitionWithoutName0,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAlertingChannelEmailDefinition, "id"),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, AlertingChannelFieldName, "name 0"),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, AlertingChannelFieldFullName, "prefix name 0 suffix"),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, AlertingChannelEmailFieldEmails+".0", "EMAIL1"),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, AlertingChannelEmailFieldEmails+".1", "EMAIL2"),
				),
			},
			resource.TestStep{
				Config: resourceDefinitionWithoutName1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAlertingChannelEmailDefinition, "id"),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, AlertingChannelFieldName, "name 1"),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, AlertingChannelFieldFullName, "prefix name 1 suffix"),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, AlertingChannelEmailFieldEmails+".0", "EMAIL1"),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, AlertingChannelEmailFieldEmails+".1", "EMAIL2"),
				),
			},
		},
	})
}

func TestResourceAlertingChannelEmailDefinition(t *testing.T) {
	resource := CreateResourceAlertingChannelEmail()

	validateAlertingChannelEmailResourceSchema(resource.Schema, t)

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

func validateAlertingChannelEmailResourceSchema(schemaMap map[string]*schema.Schema, t *testing.T) {
	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelFieldName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingChannelFieldFullName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeListOfStrings(AlertingChannelEmailFieldEmails)
}

func TestShouldSuccessfullyReadAlertingChannelEmailFromInstanaAPIWhenBaseDataIsReturned(t *testing.T) {
	expectedModel := createTestAlertingChannelEmailModel()
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		resourceData := testHelper.CreateEmptyAlertingChannelEmailResourceData()
		resourceData.SetId(alertingChannelEmailID)
		mockAlertingChannelEmailApi := mocks.NewMockAlertingChannelResource(ctrl)

		mockInstanaAPI.EXPECT().AlertingChannels().Return(mockAlertingChannelEmailApi).Times(1)
		mockAlertingChannelEmailApi.EXPECT().GetOne(gomock.Eq(alertingChannelEmailID)).Return(expectedModel, nil).Times(1)

		err := ReadAlertingChannelEmail(resourceData, providerMeta)

		if err != nil {
			t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
		}
		verifyAlertingChannelEmailModelAppliedToResource(expectedModel, resourceData, t)
	})
}

func TestShouldFailToReadAlertingChannelEmailFromInstanaAPIWhenIDIsMissing(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		resourceData := testHelper.CreateEmptyAlertingChannelEmailResourceData()

		err := ReadAlertingChannelEmail(resourceData, providerMeta)

		if err == nil || !strings.HasPrefix(err.Error(), "ID of alerting channel email") {
			t.Fatal("Expected error to occur because of missing id")
		}
	})
}

func TestShouldFailToReadAlertingChannelEmailFromInstanaAPIAndDeleteResourceWhenRoleDoesNotExist(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		resourceData := testHelper.CreateEmptyAlertingChannelEmailResourceData()
		resourceData.SetId(alertingChannelEmailID)
		mockAlertingChannelEmailApi := mocks.NewMockAlertingChannelResource(ctrl)

		mockInstanaAPI.EXPECT().AlertingChannels().Return(mockAlertingChannelEmailApi).Times(1)
		mockAlertingChannelEmailApi.EXPECT().GetOne(gomock.Eq(alertingChannelEmailID)).Return(restapi.AlertingChannel{}, restapi.ErrEntityNotFound).Times(1)

		err := ReadAlertingChannelEmail(resourceData, providerMeta)

		if err != nil {
			t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
		}
		if len(resourceData.Id()) > 0 {
			t.Fatal("Expected ID to be cleaned to destroy resource")
		}
	})
}

func TestShouldFailToReadAlertingChannelEmailFromInstanaAPIAndReturnErrorWhenAPICallFails(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		resourceData := testHelper.CreateEmptyAlertingChannelEmailResourceData()
		resourceData.SetId(alertingChannelEmailID)
		expectedError := errors.New("test")
		mockAlertingChannelEmailApi := mocks.NewMockAlertingChannelResource(ctrl)

		mockInstanaAPI.EXPECT().AlertingChannels().Return(mockAlertingChannelEmailApi).Times(1)
		mockAlertingChannelEmailApi.EXPECT().GetOne(gomock.Eq(alertingChannelEmailID)).Return(restapi.AlertingChannel{}, expectedError).Times(1)

		err := ReadAlertingChannelEmail(resourceData, providerMeta)

		if err == nil || err != expectedError {
			t.Fatal("Expected error should be returned")
		}
		if len(resourceData.Id()) == 0 {
			t.Fatal("Expected ID should still be set")
		}
	})
}

func TestShouldCreateAlertingChannelEmailThroughInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		data := createTestAlertingChannelEmailData()
		resourceData := testHelper.CreateAlertingChannelEmailResourceData(data)
		expectedModel := createTestAlertingChannelEmailModel()
		mockAlertingChannelEmailApi := mocks.NewMockAlertingChannelResource(ctrl)

		mockInstanaAPI.EXPECT().AlertingChannels().Return(mockAlertingChannelEmailApi).Times(1)
		mockResourceNameFormatter.EXPECT().Format(data[AlertingChannelFieldName]).Return(data[AlertingChannelFieldName]).Times(1)
		mockAlertingChannelEmailApi.EXPECT().Upsert(gomock.AssignableToTypeOf(restapi.AlertingChannel{})).Return(expectedModel, nil).Times(1)

		err := CreateAlertingChannelEmail(resourceData, providerMeta)

		if err != nil {
			t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
		}
		verifyAlertingChannelEmailModelAppliedToResource(expectedModel, resourceData, t)
	})
}

func TestShouldReturnErrorWhenCreateAlertingChannelEmailFailsThroughInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		data := createTestAlertingChannelEmailData()
		resourceData := testHelper.CreateAlertingChannelEmailResourceData(data)
		expectedError := errors.New("test")
		mockAlertingChannelEmailApi := mocks.NewMockAlertingChannelResource(ctrl)

		mockInstanaAPI.EXPECT().AlertingChannels().Return(mockAlertingChannelEmailApi).Times(1)
		mockResourceNameFormatter.EXPECT().Format(data[AlertingChannelFieldName]).Return(data[AlertingChannelFieldName]).Times(1)
		mockAlertingChannelEmailApi.EXPECT().Upsert(gomock.AssignableToTypeOf(restapi.AlertingChannel{})).Return(restapi.AlertingChannel{}, expectedError).Times(1)

		err := CreateAlertingChannelEmail(resourceData, providerMeta)

		if err == nil || expectedError != err {
			t.Fatal("Expected definned error to be returned")
		}
	})
}

func TestShouldDeleteAlertingChannelEmailThroughInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		id := "test-id"
		data := createTestAlertingChannelEmailData()
		resourceData := testHelper.CreateAlertingChannelEmailResourceData(data)
		resourceData.SetId(id)
		mockAlertingChannelEmailApi := mocks.NewMockAlertingChannelResource(ctrl)

		mockInstanaAPI.EXPECT().AlertingChannels().Return(mockAlertingChannelEmailApi).Times(1)
		mockResourceNameFormatter.EXPECT().Format(data[AlertingChannelFieldName]).Return(data[AlertingChannelFieldName]).Times(1)
		mockAlertingChannelEmailApi.EXPECT().DeleteByID(gomock.Eq(id)).Return(nil).Times(1)

		err := DeleteAlertingChannelEmail(resourceData, providerMeta)

		if err != nil {
			t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
		}
		if len(resourceData.Id()) > 0 {
			t.Fatal("Expected ID to be cleaned to destroy resource")
		}
	})
}

func TestShouldReturnErrorWhenDeleteAlertingChannelEmailFailsThroughInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		id := "test-id"
		data := createTestAlertingChannelEmailData()
		resourceData := testHelper.CreateAlertingChannelEmailResourceData(data)
		resourceData.SetId(id)
		expectedError := errors.New("test")
		mockAlertingChannelEmailApi := mocks.NewMockAlertingChannelResource(ctrl)

		mockInstanaAPI.EXPECT().AlertingChannels().Return(mockAlertingChannelEmailApi).Times(1)
		mockResourceNameFormatter.EXPECT().Format(data[AlertingChannelFieldName]).Return(data[AlertingChannelFieldName]).Times(1)
		mockAlertingChannelEmailApi.EXPECT().DeleteByID(gomock.Eq(id)).Return(expectedError).Times(1)

		err := DeleteAlertingChannelEmail(resourceData, providerMeta)

		if err == nil || err != expectedError {
			t.Fatal("Expected error to be returned")
		}
		if len(resourceData.Id()) == 0 {
			t.Fatal("Expected ID not to be cleaned to avoid resource is destroy")
		}
	})
}

func verifyAlertingChannelEmailModelAppliedToResource(model restapi.AlertingChannel, resourceData *schema.ResourceData, t *testing.T) {
	if model.ID != resourceData.Id() {
		t.Fatal("Expected ID to be identical")
	}
	if model.Name != resourceData.Get(AlertingChannelFieldFullName).(string) {
		t.Fatal("Expected Full Name to match Name of API")
	}
	if !cmp.Equal(model.Emails, ReadStringArrayParameterFromResource(resourceData, AlertingChannelEmailFieldEmails)) {
		t.Fatal("Expected Emails to be identical")
	}
}

func createTestAlertingChannelEmailModel() restapi.AlertingChannel {
	return restapi.AlertingChannel{
		ID:     "id",
		Name:   "name",
		Emails: []string{"Email1", "Email2"},
	}
}

func createTestAlertingChannelEmailData() map[string]interface{} {
	emails := make([]interface{}, 2)
	emails[0] = "Email1"
	emails[1] = "Email2"

	data := make(map[string]interface{})
	data[AlertingChannelFieldName] = "name"
	data[AlertingChannelEmailFieldEmails] = emails
	return data
}
