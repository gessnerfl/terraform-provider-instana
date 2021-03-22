package instana_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

const testBuiltInEventSpecDataSource = "data.instana_builtin_event_spec.test"

const dataSourceBuiltinEventSpecificationDefinitionTemplate = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:{{PORT}}"
  default_name_prefix = "prefix"
  default_name_suffix = "suffix"
}

data "instana_builtin_event_spec" "test" {
  name = "System load too high"
  short_plugin_id = "host"
}
`

func TestDatasourceBuiltInEventsEndToEnd(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodGet, restapi.BuiltinEventSpecificationResourcePath, func(w http.ResponseWriter, r *http.Request) {
		wd, err := os.Getwd()
		if err != nil {
			httpServer.WriteInternalServerError(w, err)
		} else {
			json, err := ioutil.ReadFile(fmt.Sprintf("%s/data-source-builtin-event_test.http-response.json", wd))
			if err != nil {
				httpServer.WriteInternalServerError(w, err)
			} else {
				httpServer.WriteJSONResponse(w, json)
			}
		}

	})
	httpServer.Start()
	defer httpServer.Close()

	resourceDefinition := strings.ReplaceAll(dataSourceBuiltinEventSpecificationDefinitionTemplate, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))

	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceDefinition,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testBuiltInEventSpecDataSource, "id"),
					resource.TestCheckResourceAttr(testBuiltInEventSpecDataSource, BuiltinEventSpecificationFieldName, "System load too high"),
					resource.TestCheckResourceAttr(testBuiltInEventSpecDataSource, BuiltinEventSpecificationFieldDescription, "Checks whether the system load is too high, by comparing the load against 2 times the CPU cores of the machine."),
					resource.TestCheckResourceAttr(testBuiltInEventSpecDataSource, BuiltinEventSpecificationFieldShortPluginID, "host"),
					resource.TestCheckResourceAttr(testBuiltInEventSpecDataSource, BuiltinEventSpecificationFieldSeverity, restapi.SeverityWarning.GetTerraformRepresentation()),
					resource.TestCheckResourceAttr(testBuiltInEventSpecDataSource, BuiltinEventSpecificationFieldSeverityCode, strconv.Itoa(restapi.SeverityWarning.GetAPIRepresentation())),
					resource.TestCheckResourceAttr(testBuiltInEventSpecDataSource, BuiltinEventSpecificationFieldEnabled, "true"),
					resource.TestCheckResourceAttr(testBuiltInEventSpecDataSource, BuiltinEventSpecificationFieldTriggering, "false"),
				),
			},
		},
	})
}

func TestDataSourceBuiltinEventDefinition(t *testing.T) {
	sut := NewBuiltinEventDataSource().CreateResource()

	schemaAssert := testutils.NewTerraformSchemaAssert(sut.Schema, t)

	require.Equal(t, 0, sut.SchemaVersion)
	require.Equal(t, 7, len(sut.Schema))
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(BuiltinEventSpecificationFieldName)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(BuiltinEventSpecificationFieldShortPluginID)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(BuiltinEventSpecificationFieldDescription)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(BuiltinEventSpecificationFieldSeverity)
	schemaAssert.AssertSchemaIsComputedAndOfTypeInt(BuiltinEventSpecificationFieldSeverityCode)
	schemaAssert.AssertSchemaIsComputedAndOfTypeBool(BuiltinEventSpecificationFieldEnabled)
	schemaAssert.AssertSchemaIsComputedAndOfTypeBool(BuiltinEventSpecificationFieldTriggering)
}

func TestShouldSuccessFullyReadBuiltInEventWhenResponseContainsAnObjectWithTheExactNameAndShortPluginID(t *testing.T) {
	requestedName := "name-5"
	requestedPluginId := "plugin-id-5"

	resourceData, err := executeReadWithTenGeneratedObjectsAsResponse(requestedName, requestedPluginId, t)

	require.NoError(t, err)
	require.Equal(t, "id-5", resourceData.Id())
	require.Equal(t, requestedName, resourceData.Get(BuiltinEventSpecificationFieldName))
	require.Equal(t, "description-5", resourceData.Get(BuiltinEventSpecificationFieldDescription))
	require.Equal(t, requestedPluginId, resourceData.Get(BuiltinEventSpecificationFieldShortPluginID))
	require.Equal(t, restapi.SeverityWarning.GetTerraformRepresentation(), resourceData.Get(BuiltinEventSpecificationFieldSeverity))
	require.Equal(t, restapi.SeverityWarning.GetAPIRepresentation(), resourceData.Get(BuiltinEventSpecificationFieldSeverityCode))
	require.True(t, resourceData.Get(BuiltinEventSpecificationFieldEnabled).(bool))
	require.True(t, resourceData.Get(BuiltinEventSpecificationFieldTriggering).(bool))
}

func TestShouldFailtToReadBuiltInEventWhenNoObjectFromResponseMatchesTheNameAndShortPluginIDCombination(t *testing.T) {
	requestedName := "name-invalid"
	requestedPluginId := "plugin-id-invalid"

	_, err := executeReadWithTenGeneratedObjectsAsResponse(requestedName, requestedPluginId, t)

	require.Error(t, err)
	require.Contains(t, err.Error(), "no built in event found")
}

func executeReadWithTenGeneratedObjectsAsResponse(requestedName string, requestedPluginId string, t *testing.T) (*schema.ResourceData, error) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sut := NewBuiltinEventDataSource().CreateResource()

	response := createBuiltinEventSpecifications(10)
	builtInEventSpecificationAPI := mocks.NewMockReadOnlyRestResource(ctrl)
	builtInEventSpecificationAPI.EXPECT().GetAll().Times(1).Return(response, nil)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)
	mockInstanaAPI.EXPECT().BuiltinEventSpecifications().Times(1).Return(builtInEventSpecificationAPI)

	meta := &ProviderMeta{InstanaAPI: mockInstanaAPI}
	resourceData := schema.TestResourceDataRaw(t, sut.Schema, map[string]interface{}{BuiltinEventSpecificationFieldName: requestedName, BuiltinEventSpecificationFieldShortPluginID: requestedPluginId})

	err := sut.Read(resourceData, meta)
	if err != nil {
		return nil, err
	}
	return resourceData, nil
}

func TestShouldFailtToReadBuiltInEventWhenAPIRequestFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sut := NewBuiltinEventDataSource().CreateResource()

	expectedError := errors.New("test")
	requestedName := "name-1"
	requestedPluginId := "plugin-id-1"

	builtInEventSpecificationAPI := mocks.NewMockReadOnlyRestResource(ctrl)
	builtInEventSpecificationAPI.EXPECT().GetAll().Times(1).Return(nil, expectedError)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)
	mockInstanaAPI.EXPECT().BuiltinEventSpecifications().Times(1).Return(builtInEventSpecificationAPI)

	meta := &ProviderMeta{InstanaAPI: mockInstanaAPI}
	resourceData := schema.TestResourceDataRaw(t, sut.Schema, map[string]interface{}{BuiltinEventSpecificationFieldName: requestedName, BuiltinEventSpecificationFieldShortPluginID: requestedPluginId})

	err := sut.Read(resourceData, meta)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

func TestShouldFailtToReadBuiltInEventWhenSeverityCannotBeConvertedFromItsCodeRepresentation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sut := NewBuiltinEventDataSource().CreateResource()

	requestedName := "name-1"
	requestedPluginId := "plugin-id-1"

	builtinEvent := createBuiltinEventSpecification(1)
	builtinEvent.Severity = 100
	response := []restapi.InstanaDataObject{builtinEvent}
	builtInEventSpecificationAPI := mocks.NewMockReadOnlyRestResource(ctrl)
	builtInEventSpecificationAPI.EXPECT().GetAll().Times(1).Return(&response, nil)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)
	mockInstanaAPI.EXPECT().BuiltinEventSpecifications().Times(1).Return(builtInEventSpecificationAPI)

	meta := &ProviderMeta{InstanaAPI: mockInstanaAPI}
	resourceData := schema.TestResourceDataRaw(t, sut.Schema, map[string]interface{}{BuiltinEventSpecificationFieldName: requestedName, BuiltinEventSpecificationFieldShortPluginID: requestedPluginId})

	err := sut.Read(resourceData, meta)

	require.Error(t, err)
	require.Contains(t, err.Error(), "100 is not a valid severity")
}

func createBuiltinEventSpecifications(count int) *[]restapi.InstanaDataObject {
	result := make([]restapi.InstanaDataObject, count)
	for i := 0; i < count; i++ {
		result[i] = createBuiltinEventSpecification(i)
	}
	return &result
}

func createBuiltinEventSpecification(id int) restapi.BuiltinEventSpecification {
	description := fmt.Sprintf("description-%d", id)
	return restapi.BuiltinEventSpecification{
		ID:            fmt.Sprintf("id-%d", id),
		Name:          fmt.Sprintf("name-%d", id),
		Description:   &description,
		ShortPluginID: fmt.Sprintf("plugin-id-%d", id),
		Severity:      restapi.SeverityWarning.GetAPIRepresentation(),
		Triggering:    true,
		Enabled:       true,
	}
}
