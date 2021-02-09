package instana_test

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/gessnerfl/terraform-provider-instana/utils"

	"github.com/stretchr/testify/assert"
)

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

func TestCRUDOfAlertingChannelEmailResourceWithMockServer(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, alertingChannelEmailApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, alertingChannelEmailApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, alertingChannelEmailApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		json := strings.ReplaceAll(alertingChannelEmailServerResponseTemplate, "{{id}}", vars["id"])
		w.Header().Set(contentType, r.Header.Get(contentType))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceDefinitionWithoutName := strings.ReplaceAll(resourceAlertingChannelEmailDefinitionTemplate, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))
	resourceDefinitionWithoutName0 := strings.ReplaceAll(resourceDefinitionWithoutName, iteratorPlaceholder, "0")
	resourceDefinitionWithoutName1 := strings.ReplaceAll(resourceDefinitionWithoutName, iteratorPlaceholder, "1")

	hashFunctionEmails := schema.HashSchema(AlertingChannelEmailEmailsSchemaField.Elem.(*schema.Schema))
	emailAddress1 := "EMAIL1"
	emailAddress2 := "EMAIL2"
	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceDefinitionWithoutName0,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAlertingChannelEmailDefinition, "id"),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, AlertingChannelFieldName, "name 0"),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, AlertingChannelFieldFullName, "prefix name 0 suffix"),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, fmt.Sprintf("%s.%d", AlertingChannelEmailFieldEmails, hashFunctionEmails(emailAddress1)), emailAddress1),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, fmt.Sprintf("%s.%d", AlertingChannelEmailFieldEmails, hashFunctionEmails(emailAddress2)), emailAddress2),
				),
			},
			{
				Config: resourceDefinitionWithoutName1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAlertingChannelEmailDefinition, "id"),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, AlertingChannelFieldName, "name 1"),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, AlertingChannelFieldFullName, "prefix name 1 suffix"),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, fmt.Sprintf("%s.%d", AlertingChannelEmailFieldEmails, hashFunctionEmails(emailAddress1)), emailAddress1),
					resource.TestCheckResourceAttr(testAlertingChannelEmailDefinition, fmt.Sprintf("%s.%d", AlertingChannelEmailFieldEmails, hashFunctionEmails(emailAddress2)), emailAddress2),
				),
			},
		},
	})
}

func TestResourceAlertingChannelEmailDefinition(t *testing.T) {
	resource := NewAlertingChannelEmailResourceHandle()

	schemaMap := resource.Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(AlertingChannelFieldName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(AlertingChannelFieldFullName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeSetOfStrings(AlertingChannelEmailFieldEmails)
}

func TestShouldReturnCorrectResourceNameForAlertingChannelEmail(t *testing.T) {
	name := NewAlertingChannelEmailResourceHandle().ResourceName

	assert.Equal(t, "instana_alerting_channel_email", name, "Expected resource name to be instana_alerting_channel_email")
}

func TestAlertingChannelEmailResourceShouldHaveSchemaVersionOne(t *testing.T) {
	assert.Equal(t, 1, NewAlertingChannelEmailResourceHandle().SchemaVersion)
}

func TestAlertingChannelEmailShouldHaveOneStateUpgraderForVersionZero(t *testing.T) {
	resourceHandler := NewAlertingChannelEmailResourceHandle()

	assert.Equal(t, 1, len(resourceHandler.StateUpgraders))
	assert.Equal(t, 0, resourceHandler.StateUpgraders[0].Version)
}

func TestShouldReturnStateOfAlertingChannelEmailUnchangedWhenMigratingFromVersion0ToVersion1(t *testing.T) {
	emails := []interface{}{"email1", "email2"}
	name := "name"
	fullname := "fullname"
	rawData := make(map[string]interface{})
	rawData[AlertingChannelFieldName] = name
	rawData[AlertingChannelFieldFullName] = fullname
	rawData[AlertingChannelEmailFieldEmails] = emails
	meta := "dummy"

	result, err := NewAlertingChannelEmailResourceHandle().StateUpgraders[0].Upgrade(rawData, meta)

	assert.Nil(t, err)
	assert.Equal(t, rawData, result)
}

func TestShouldUpdateResourceStateForAlertingChannelEmail(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelEmailResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	data := restapi.AlertingChannel{
		ID:     "id",
		Name:   "name",
		Emails: []string{"email1", "email2"},
	}

	err := resourceHandle.UpdateState(resourceData, &data)

	assert.Nil(t, err)
	assert.Equal(t, "id", resourceData.Id(), "id should be equal")
	assert.Equal(t, "name", resourceData.Get(AlertingChannelFieldFullName), "name should be equal to full name")

	emails := resourceData.Get(AlertingChannelEmailFieldEmails).(*schema.Set)
	assert.Equal(t, 2, emails.Len())
	assert.Contains(t, emails.List(), "email1")
	assert.Contains(t, emails.List(), "email2")
}

func TestShouldConvertStateOfAlertingChannelEmailToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewAlertingChannelEmailResourceHandle()
	emails := []string{"email1", "email2"}
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId("id")
	resourceData.Set(AlertingChannelFieldName, "name")
	resourceData.Set(AlertingChannelFieldFullName, "prefix name suffix")
	resourceData.Set(AlertingChannelEmailFieldEmails, emails)

	model, err := resourceHandle.MapStateToDataObject(resourceData, utils.NewResourceNameFormatter("prefix ", " suffix"))

	assert.Nil(t, err)
	assert.IsType(t, &restapi.AlertingChannel{}, model, "Model should be an alerting channel")
	assert.Equal(t, "id", model.GetID())
	assert.Equal(t, "prefix name suffix", model.(*restapi.AlertingChannel).Name, "name should be equal to full name")
	assert.Len(t, model.(*restapi.AlertingChannel).Emails, 2)
	assert.Contains(t, model.(*restapi.AlertingChannel).Emails, "email1")
	assert.Contains(t, model.(*restapi.AlertingChannel).Emails, "email2")
}
