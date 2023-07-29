package instana_test

import (
	"fmt"
	"log"
	"net/http"

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
	default_name_prefix = "prefix"
	default_name_suffix = "suffix"
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
	httpServer.AddRoute(http.MethodGet, resourcePath, func(w http.ResponseWriter, r *http.Request) {
		//var json string
		//if templateVars != nil && len(templateVars) > 0 {
		//	json = fmt.Sprintf(responseTemplate, templateVars...)
		//} else {
		//	json = responseTemplate
		//}
		json := responseTemplate
		w.Header().Set(contentType, r.Header.Get(contentType))
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(json))
		if err != nil {
			log.Fatalf("failed to write response: %s", err)
		}
	})
	httpServer.AddRoute(http.MethodPut, pathTemplate, responseHandler)
	httpServer.AddRoute(http.MethodDelete, pathTemplate, responseHandler)
	httpServer.AddRoute(http.MethodGet, pathTemplate, responseHandler)
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
const resourceFullName = "prefix name suffix"

func formatResourceName(iteration int) string {
	return fmt.Sprintf("name %d", iteration)
}

func formatResourceFullName(iteration int) string {
	return fmt.Sprintf("prefix name %d suffix", iteration)
}

func copyMap(input map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range input {
		result[k] = v
	}

	return result
}

type testPair[A any, E any] struct {
	name     string
	input    A
	expected E
}
