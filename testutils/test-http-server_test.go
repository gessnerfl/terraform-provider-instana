package testutils_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	testutils "github.com/gessnerfl/terraform-provider-instana/testutils"
)

func TestShouldStartNewInstanceWithDynamicPortAndStopTheServerOnClose(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	path := "/test"
	server := testutils.NewTestHTTPServer()
	server.AddRoute(http.MethodPost, path, testutils.EchoHandlerFunc)
	server.Start()
	defer server.Close()

	url := fmt.Sprintf("https://localhost:%d%s", server.GetPort(), path)
	testString := "test string"

	resp, err := http.Post(url, "test/plain", strings.NewReader(testString))
	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expected http status code 200 but got %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Expected to get http response")
	}
	responseString := string(responseBytes)
	if testString != responseString {
		t.Fatalf("Expected to get '%s' but got '%s'", testString, responseString)
	}
}

func TestShouldCreateRandomPortNumber(t *testing.T) {
	result := testutils.RandomPort()

	if result < testutils.MinPortNumber || result > testutils.MaxPortNumber {
		t.Fatalf("Expected port number between %d and %d but got %d", testutils.MinPortNumber, testutils.MaxPortNumber, result)
	}
}
