package instana_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/gessnerfl/terraform-provider-instana/utils"
)

var testWebsiteMonitoringConfigProviders = map[string]terraform.ResourceProvider{
	"instana": Provider(),
}

const websiteMonitoringConfigTerraformTemplate = `
provider "instana" {
	api_token = "test-token"
	endpoint = "localhost:{{PORT}}"
	default_name_prefix = "prefix"
	default_name_suffix = "suffix"
}

resource "instana_website_monitoring_config" "example_website_monitoring_config" {
	name = "name {{ITERATOR}}"
}
`

const (
	websiteMonitoringConfigApiPath    = restapi.WebsiteMonitoringConfigResourcePath + "/{id}"
	websiteMonitoringConfigDefinition = "instana_website_monitoring_config.example_website_monitoring_config"
	websiteMonitoringConfigID         = "id"
	websiteMonitoringConfigName       = "name"
	websiteMonitoringConfigFullName   = "prefix name suffix"
)

func TestCRUDOfWebsiteMonitoringConfiguration(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()

	server := newWebsiteMonitoringConfigTestServer()
	defer server.Close()
	server.Start()

	resourceDefinitionWithoutName := strings.ReplaceAll(websiteMonitoringConfigTerraformTemplate, "{{PORT}}", strconv.Itoa(server.GetPort()))
	resourceDefinitionWithName0 := strings.ReplaceAll(resourceDefinitionWithoutName, iteratorPlaceholder, "0")
	resourceDefinitionWithName1 := strings.ReplaceAll(resourceDefinitionWithoutName, iteratorPlaceholder, "1")

	resource.UnitTest(t, resource.TestCase{
		Providers: testWebsiteMonitoringConfigProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceDefinitionWithName0,
				Check:  resource.ComposeTestCheckFunc(createWebsiteMonitoringConfigTestCheckFunctions(0)...),
			},
			{
				Config: resourceDefinitionWithName1,
				Check:  resource.ComposeTestCheckFunc(createWebsiteMonitoringConfigTestCheckFunctions(1)...),
			},
		},
	})
}

func createWebsiteMonitoringConfigTestCheckFunctions(iteration int) []resource.TestCheckFunc {
	testCheckFunctions := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(websiteMonitoringConfigDefinition, "id"),
		resource.TestCheckResourceAttr(websiteMonitoringConfigDefinition, WebsiteMonitoringConfigFieldName, fmt.Sprintf("name %d", iteration)),
		resource.TestCheckResourceAttr(websiteMonitoringConfigDefinition, WebsiteMonitoringConfigFieldFullName, fmt.Sprintf("prefix name %d suffix", iteration)),
		resource.TestCheckResourceAttr(websiteMonitoringConfigDefinition, WebsiteMonitoringConfigFieldAppName, fmt.Sprintf("prefix name %d suffix", iteration)),
	}
	return testCheckFunctions
}

func newWebsiteMonitoringConfigTestServer() *websiteMonitoringConfigTestServer {
	return &websiteMonitoringConfigTestServer{httpServer: testutils.NewTestHTTPServer()}
}

type websiteMonitoringConfigTestServer struct {
	httpServer  *testutils.TestHTTPServer
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

func (s *websiteMonitoringConfigTestServer) GetPort() int {
	if s.httpServer != nil {
		return s.httpServer.GetPort()
	}
	return -1
}

func (s *websiteMonitoringConfigTestServer) onPost(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name != "" {
		s.serverState = &restapi.WebsiteMonitoringConfig{
			ID:      utils.RandomString(10),
			Name:    name,
			AppName: name,
		}

		json.NewEncoder(w).Encode(s.serverState)
		w.Header().Set(contentType, r.Header.Get(contentType))
		w.WriteHeader(http.StatusOK)
	} else {
		w.Write([]byte("Name is missing"))
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

			json.NewEncoder(w).Encode(s.serverState)
			w.Header().Set(contentType, r.Header.Get(contentType))
			w.WriteHeader(http.StatusOK)
		} else {
			w.Write([]byte("Name is missing"))
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.Write([]byte("Entity with id %s not found"))
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *websiteMonitoringConfigTestServer) onGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if s.serverState != nil && vars["id"] == s.serverState.ID {
		json.NewEncoder(w).Encode(s.serverState)
		w.Header().Set(contentType, r.Header.Get(contentType))
		w.WriteHeader(http.StatusOK)
	} else {
		w.Write([]byte("Entity with id %s not found"))
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *websiteMonitoringConfigTestServer) Close() {
	if s.httpServer != nil {
		s.httpServer.Close()
	}
}
