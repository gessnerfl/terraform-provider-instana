package instana_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"os"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const contentType = "Content-Type"
const trueAsString = "true"
const falseAsString = "false"

var providerConfig = `
provider "instana" {
	api_token 			= "test-token"
	endpoint 			= "localhost:%d"
    tls_skip_verify     = true
}
`

var testProviderFactory = map[string]func() (*schema.Provider, error){
	"instana": func() (*schema.Provider, error) { return Provider(), nil },
}

func appendProviderConfig(resourceConfig string, serverPort int64) string {
	return fmt.Sprintf(providerConfig, serverPort) + " \n\n" + resourceConfig
}

func createMockHttpServerForResource(resourcePath string, responseTemplate string, templateVars ...interface{}) testutils.TestHTTPServer {
	pathTemplate := resourcePath + "/{id}"
	httpServer := testutils.NewTestHTTPServer()
	responseHandler := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		path := resourcePath + "/" + vars["id"]
		callCount := getZeroBasedCallCount(httpServer, http.MethodPut, path)
		var json string
		if templateVars != nil {
			json = formatResponseTemplate(responseTemplate, vars["id"], callCount, templateVars...)
		} else {
			json = formatResponseTemplate(responseTemplate, vars["id"], callCount)
		}
		w.Header().Set(contentType, r.Header.Get(contentType))
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(json))
		if err != nil {
			log.Fatalf("failed to write response: %s", err)
		}
	}
	httpServer.AddRoute(http.MethodPut, pathTemplate, responseHandler)
	httpServer.AddRoute(http.MethodDelete, pathTemplate, responseHandler)
	httpServer.AddRoute(http.MethodGet, pathTemplate, responseHandler)
	return httpServer
}

type responseContentProvider interface {
	provide() ([]byte, error)
}

func newFileContentResponseProvider(filePath string) responseContentProvider {
	return &fileContentResponseProvider{filePath: filePath}
}

type fileContentResponseProvider struct {
	filePath string
}

func (p *fileContentResponseProvider) provide() ([]byte, error) {
	return os.ReadFile(p.filePath)
}

func newStringContentResponseProvider(content string) responseContentProvider {
	return &stringContentResponseProvider{content: content}
}

type stringContentResponseProvider struct {
	content string
}

func (p *stringContentResponseProvider) provide() ([]byte, error) {
	return []byte(p.content), nil
}

func createMockHttpServerForDataSource(resourcePath string, responseContent responseContentProvider) testutils.TestHTTPServer {
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodGet, resourcePath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(contentType, r.Header.Get(contentType))
		json, err := responseContent.provide()
		if err != nil {
			httpServer.WriteInternalServerError(w, err)
		} else {
			w.WriteHeader(http.StatusOK)
			httpServer.WriteJSONResponse(w, json)
		}
	})
	return httpServer
}

func formatResponseTemplate(template string, id string, iteration int, vars ...interface{}) string {
	allVars := make([]interface{}, len(vars)+2)
	allVars[0] = id
	allVars[1] = iteration
	for i, v := range vars {
		allVars[i+2] = v
	}
	return fmt.Sprintf(template, allVars...)
}

func getZeroBasedCallCount(httpServer testutils.TestHTTPServer, method string, path string) int {
	count := httpServer.GetCallCount(method, path)
	if count == 0 {
		return count
	}
	return count - 1
}

func testStepImport(resourceName string) resource.TestStep {
	return resource.TestStep{
		ResourceName:      resourceName,
		ImportState:       true,
		ImportStateVerify: true,
	}
}

func testStepImportWithCustomID(resourceName string, resourceID string) resource.TestStep {
	return resource.TestStep{
		ResourceName:      resourceName,
		ImportState:       true,
		ImportStateVerify: true,
		ImportStateId:     resourceID,
	}
}

const resourceName = "name"

func formatResourceName(iteration int) string {
	return fmt.Sprintf("name %d", iteration)
}

type testPair[A any, E any] struct {
	name     string
	input    A
	expected E
}

func setValueOnResourceData(t *testing.T, r *schema.ResourceData, key string, value interface{}) {
	err := r.Set(key, value)
	require.NoError(t, err)
}
