package testutils_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	testutils "github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/stretchr/testify/assert"
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

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	defer resp.Body.Close()
	responseBytes, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	responseString := string(responseBytes)

	assert.Equal(t, testString, responseString)
}

func TestShouldCreateRandomPortNumber(t *testing.T) {
	result := testutils.RandomPort()

	assert.LessOrEqual(t, result, testutils.MaxPortNumber)
	assert.GreaterOrEqual(t, result, testutils.MinPortNumber)
}
