package services_test

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi/services"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

const testPath = "/test"
const testID = "test-1234"
const testData = "testData"
const testPathWithID = testPath + "/" + testID

func TestShouldReturnDataForSuccessfulGetOneRequest(t *testing.T) {
	httpServer := setupAndStartHttpServerWithOKResponseCode(http.MethodGet, testPathWithID)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	response, err := restClient.GetOne(testID, testPath)

	verifySuccessfullGetOrPut(response, err, t)
}

func TestShouldReturnErrorMessageForGetOneRequestWhenStatusIsNotASuccessStatusAndNotEnityNotFound(t *testing.T) {
	statusCode := http.StatusBadRequest
	httpServer := setupAndStartHttpServer(http.MethodGet, testPathWithID, statusCode)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	_, err := restClient.GetOne(testID, testPath)

	verifyFailedCallWithStatusCodeIsResponse(err, statusCode, t)
}

func TestShouldReturnNotFoundErrorMessageForGetOneRequestWhenStatusIsNotEnityNotFound(t *testing.T) {
	httpServer := setupAndStartHttpServer(http.MethodGet, testPathWithID, http.StatusNotFound)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	data, err := restClient.GetOne(testID, testPath)

	verifyNotFoundResponse(data, err, t)
}

func TestShouldReturnDataForSuccessfulPutRequest(t *testing.T) {
	httpServer := setupAndStartHttpServerWithOKResponseCode(http.MethodPut, testPathWithID)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	response, err := restClient.Put(testDataObject{id: testID}, testPath)

	verifySuccessfullGetOrPut(response, err, t)
}

func TestShouldReturnErrorMessageForPutRequestWhenStatusIsNotASuccessStatusAndNotEnityNotFound(t *testing.T) {
	statusCode := http.StatusBadRequest
	httpServer := setupAndStartHttpServer(http.MethodPut, testPathWithID, statusCode)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	_, err := restClient.Put(testDataObject{id: testID}, testPath)

	verifyFailedCallWithStatusCodeIsResponse(err, statusCode, t)
}

type testDataObject struct {
	id string
}

//GetID implementation of InstanaDataObject
func (tdo testDataObject) GetID() string {
	return tdo.id
}

//Validate implementation of InstanaDataObject
func (tdo testDataObject) Validate() error {
	return nil
}

func TestShouldReturnNothingForSuccessfulDeleteRequest(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := setupAndStartHttpServerWithOKResponseCode(http.MethodDelete, testPathWithID)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	err := restClient.Delete(testID, testPath)

	if err != nil {
		t.Fatalf("Expected no error to be returned but got %s", err)
	}
}

func TestShouldReturnErrorMessageForDeleteRequestWhenStatusIsNotASuccessStatusAndNotEnityNotFound(t *testing.T) {
	statusCode := http.StatusBadRequest
	httpServer := setupAndStartHttpServer(http.MethodDelete, testPathWithID, statusCode)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	err := restClient.Delete(testID, testPath)

	verifyFailedCallWithStatusCodeIsResponse(err, statusCode, t)
}

func setupAndStartHttpServerWithOKResponseCode(httpMethod string, fullPath string) *testutils.TestHTTPServer {
	return setupAndStartHttpServer(httpMethod, fullPath, 200)
}

func setupAndStartHttpServer(httpMethod string, fullPath string, statusCode int) *testutils.TestHTTPServer {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(httpMethod, fullPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(statusCode)
		w.Write([]byte(testData))
	})
	httpServer.Start()
	return httpServer
}

func createSut(httpServer *testutils.TestHTTPServer) restapi.RestClient {
	return NewClient("api-token", fmt.Sprintf("localhost:%d", httpServer.GetPort()))
}

func verifyNotFoundResponse(data []byte, err error, t *testing.T) {
	if err != restapi.ErrEntityNotFound {
		t.Fatal("Expected error entity not found")
	}

	if data == nil || len(data) != 0 {
		t.Fatal("Expected empty data response")
	}
}

func verifySuccessfullGetOrPut(response []byte, err error, t *testing.T) {
	if err != nil {
		t.Fatalf("Expected no error to be returned but got %s", err)
	}
	responseString := string(response)
	if responseString != testData {
		t.Fatalf("Expected test data to be returned but got %s", responseString)
	}
}

func verifyFailedCallWithStatusCodeIsResponse(err error, statusCode int, t *testing.T) {
	if err == nil || !strings.Contains(err.Error(), strconv.Itoa(statusCode)) {
		t.Fatalf("Expected to receive error message with status Code %d", statusCode)
	}
}
