package instana_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/gessnerfl/terraform-provider-instana/utils"
)

const websiteMonitoringConfigTerraformTemplate = `
resource "instana_website_monitoring_config" "example_website_monitoring_config" {
	name = "name %d"
}
`

const (
	websiteMonitoringConfigApiPath    = restapi.WebsiteMonitoringConfigResourcePath + "/{id}"
	websiteMonitoringConfigDefinition = "instana_website_monitoring_config.example_website_monitoring_config"
	websiteMonitoringConfigFullName   = resourceFullName
)

func TestCRUDOfWebsiteMonitoringConfiguration(t *testing.T) {
	server := newWebsiteMonitoringConfigTestServer()
	defer server.Close()
	server.Start()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: appendProviderConfig(fmt.Sprintf(websiteMonitoringConfigTerraformTemplate, 0), server.GetPort()),
				Check:  resource.ComposeTestCheckFunc(createWebsiteMonitoringConfigTestCheckFunctions(0)...),
			},
			testStepImport(websiteMonitoringConfigDefinition),
			{
				Config: appendProviderConfig(fmt.Sprintf(websiteMonitoringConfigTerraformTemplate, 1), server.GetPort()),
				Check:  resource.ComposeTestCheckFunc(createWebsiteMonitoringConfigTestCheckFunctions(1)...),
			},
			testStepImport(websiteMonitoringConfigDefinition),
		},
	})
}

func createWebsiteMonitoringConfigTestCheckFunctions(iteration int) []resource.TestCheckFunc {
	testCheckFunctions := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(websiteMonitoringConfigDefinition, "id"),
		resource.TestCheckResourceAttr(websiteMonitoringConfigDefinition, WebsiteMonitoringConfigFieldName, fmt.Sprintf("name %d", iteration)),
		resource.TestCheckResourceAttr(websiteMonitoringConfigDefinition, WebsiteMonitoringConfigFieldAppName, fmt.Sprintf("name %d", iteration)),
	}
	return testCheckFunctions
}

func newWebsiteMonitoringConfigTestServer() *websiteMonitoringConfigTestServer {
	return &websiteMonitoringConfigTestServer{httpServer: testutils.NewTestHTTPServer()}
}

type websiteMonitoringConfigTestServer struct {
	httpServer  testutils.TestHTTPServer
	serverState *restapi.WebsiteMonitoringConfig
}

func (s *websiteMonitoringConfigTestServer) Start() {
	s.httpServer = testutils.NewTestHTTPServer()
	s.httpServer.AddRoute(http.MethodPost, restapi.WebsiteMonitoringConfigResourcePath, s.onPost)
	s.httpServer.AddRoute(http.MethodPut, websiteMonitoringConfigApiPath, s.onPut)
	s.httpServer.AddRoute(http.MethodDelete, websiteMonitoringConfigApiPath, testutils.EchoHandlerFunc)
	s.httpServer.AddRoute(http.MethodGet, websiteMonitoringConfigApiPath, s.onGet)
	s.httpServer.Start()
}

func (s *websiteMonitoringConfigTestServer) GetPort() int64 {
	if s.httpServer != nil {
		return s.httpServer.GetPort()
	}
	return -1
}

// GetCallCount returns the call counter for the given method and path
func (s *websiteMonitoringConfigTestServer) GetCallCount(method string, path string) int {
	if s.httpServer != nil {
		return s.httpServer.GetCallCount(method, path)
	}
	return 9
}

func (s *websiteMonitoringConfigTestServer) onPost(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name != "" {
		s.serverState = &restapi.WebsiteMonitoringConfig{
			ID:      utils.RandomString(10),
			Name:    name,
			AppName: name,
		}

		err := json.NewEncoder(w).Encode(s.serverState)
		if err != nil {
			fmt.Printf("failed to encode json; %s\n", err)
		}
		w.Header().Set(contentType, r.Header.Get(contentType))
		w.WriteHeader(http.StatusOK)
	} else {
		_, err := w.Write([]byte("Name is missing"))
		if err != nil {
			fmt.Printf("failed to write response; %s\n", err)
		}
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (s *websiteMonitoringConfigTestServer) onPut(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if s.serverState != nil && vars["id"] == s.serverState.ID {
		name := r.URL.Query().Get("name")

		if name != "" {
			s.serverState.Name = name
			s.serverState.AppName = name

			err := json.NewEncoder(w).Encode(s.serverState)
			if err != nil {
				fmt.Printf("failed to encode json; %s\n", err)
			}
			w.Header().Set(contentType, r.Header.Get(contentType))
			w.WriteHeader(http.StatusOK)
		} else {
			_, err := w.Write([]byte("Name is missing"))
			if err != nil {
				fmt.Printf("failed to write response; %s\n", err)
			}
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		_, err := w.Write([]byte("Entity with id %s not found"))
		if err != nil {
			fmt.Printf("failed to write response; %s\n", err)
		}
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *websiteMonitoringConfigTestServer) onGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if s.serverState != nil && vars["id"] == s.serverState.ID {
		err := json.NewEncoder(w).Encode(s.serverState)
		if err != nil {
			fmt.Printf("failed to encode json; %s\n", err)
		}
		w.Header().Set(contentType, r.Header.Get(contentType))
		w.WriteHeader(http.StatusOK)
	} else {
		_, err := w.Write([]byte("Entity with id %s not found"))
		if err != nil {
			fmt.Printf("failed to write response; %s\n", err)
		}
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *websiteMonitoringConfigTestServer) Close() {
	if s.httpServer != nil {
		s.httpServer.Close()
	}
}

func TestResourceWebsiteMonitoringConfigDefinition(t *testing.T) {
	resourceHandle := NewWebsiteMonitoringConfigResourceHandle()

	schemaMap := resourceHandle.MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(WebsiteMonitoringConfigFieldName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(WebsiteMonitoringConfigFieldAppName)
}

func TestShouldUpdateResourceStateForWebsiteMonitoringConfig(t *testing.T) {
	testHelper := NewTestHelper[*restapi.WebsiteMonitoringConfig](t)
	resourceHandle := NewWebsiteMonitoringConfigResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	fullname := resourceFullName
	appname := "appname"
	data := restapi.WebsiteMonitoringConfig{
		ID:      "id",
		Name:    fullname,
		AppName: appname,
	}

	err := resourceHandle.UpdateState(resourceData, &data, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.Equal(t, "id", resourceData.Id(), "id should be equal")
	require.Equal(t, "name", resourceData.Get(WebsiteMonitoringConfigFieldName))
	require.Equal(t, fullname, resourceData.Get(WebsiteMonitoringConfigFieldFullName))
	require.Equal(t, appname, resourceData.Get(WebsiteMonitoringConfigFieldAppName))
}

func TestShouldConvertStateOfWebsiteMonitoringConfigToDataModel(t *testing.T) {
	testHelper := NewTestHelper[*restapi.WebsiteMonitoringConfig](t)
	resourceHandle := NewWebsiteMonitoringConfigResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId("id")
	setValueOnResourceData(t, resourceData, WebsiteMonitoringConfigFieldName, "name")
	setValueOnResourceData(t, resourceData, WebsiteMonitoringConfigFieldFullName, websiteMonitoringConfigFullName)

	model, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.Nil(t, err)
	require.IsType(t, &restapi.WebsiteMonitoringConfig{}, model)
	require.Equal(t, "id", model.GetIDForResourcePath())
	require.Equal(t, websiteMonitoringConfigFullName, model.Name)
}

func TestWebsiteMonitoringConfigkShouldHaveSchemaVersionZero(t *testing.T) {
	require.Equal(t, 1, NewWebsiteMonitoringConfigResourceHandle().MetaData().SchemaVersion)
}

func TestWebsiteMonitoringConfigShouldHaveOneStateUpgrader(t *testing.T) {
	require.Equal(t, 1, len(NewWebsiteMonitoringConfigResourceHandle().StateUpgraders()))
}

func TestWebsiteMonitoringConfigShouldMigrateFullnameToNameWhenExecutingFirstStateUpgraderAndFullnameIsAvailable(t *testing.T) {
	input := map[string]interface{}{
		"full_name": "test",
	}
	result, err := NewWebsiteMonitoringConfigResourceHandle().StateUpgraders()[0].Upgrade(nil, input, nil)

	require.NoError(t, err)
	require.Len(t, result, 1)
	require.NotContains(t, result, WebsiteMonitoringConfigFieldFullName)
	require.Contains(t, result, WebsiteMonitoringConfigFieldName)
	require.Equal(t, "test", result[WebsiteMonitoringConfigFieldName])
}

func TestWebsiteMonitoringConfigShouldDoNothingWhenExecutingFirstStateUpgraderAndFullnameIsAvailable(t *testing.T) {
	input := map[string]interface{}{
		"name": "test",
	}
	result, err := NewWebsiteMonitoringConfigResourceHandle().StateUpgraders()[0].Upgrade(nil, input, nil)

	require.NoError(t, err)
	require.Equal(t, input, result)
}

func TestShouldReturnCorrectResourceNameForWebsiteMonitoringConfig(t *testing.T) {
	name := NewWebsiteMonitoringConfigResourceHandle().MetaData().ResourceName

	require.Equal(t, name, "instana_website_monitoring_config")
}
