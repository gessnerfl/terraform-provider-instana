package instana_test

import (
	"fmt"
	"net/http"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const contentType = "Content-Type"

var testProviders = map[string]*schema.Provider{
	"instana": Provider(),
}

func createMockHttpServerForResource(resourcePath string, responseTemplate string, templateVars ...interface{}) testutils.TestHTTPServer {
	pathTemplate := resourcePath + "/{id}"
	testutils.DeactivateTLSServerCertificateVerification()
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
		w.Write([]byte(json))
	}
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

func toInterfaceSlice[T any](value []T) []interface{} {
	result := make([]interface{}, len(value))
	for i, v := range value {
		result[i] = v
	}
	return result
}

type testPair[A any, E any] struct {
	name     string
	input    A
	expected E
}
