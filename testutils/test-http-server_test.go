package testutils_test

import (
	"crypto/tls"
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldStartNewInstanceWithDynamicPortAndStopTheServerOnClose(t *testing.T) {
	path := "/test"
	server := testutils.NewTestHTTPServer()
	server.AddRoute(http.MethodPost, path, testutils.EchoHandlerFunc)
	server.Start()
	defer server.Close()

	url := fmt.Sprintf("https://localhost:%d%s", server.GetPort(), path)
	testString := "test string"

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}} //nolint:gosec G402
	defer tr.CloseIdleConnections()
	client := &http.Client{Transport: tr}
	resp, err := client.Post(url, "test/plain", strings.NewReader(testString))

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
