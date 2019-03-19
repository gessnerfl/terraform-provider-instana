package services_test

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi/services"
	testutils "github.com/gessnerfl/terraform-provider-instana/test-utils"
)

const testPath = "/test"
const testID = "test-1234"
const testData = "testData"
const fullPathWithoutID = "/api" + testPath
const fullPathWithID = "/api" + testPath + "/" + testID

func TestShouldReturnDataForSuccessfulGetOneRequest(t *testing.T) {
	httpServer := setupAndStartHttpServerWithOKResponseCode(http.MethodGet, fullPathWithID)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	response, err := restClient.GetOne(testID, testPath)

	verifySuccessfullGetOrPut(response, err, t)
}

func TestShouldReturnErrorMessageForGetOneRequestWhenStatusIsNotASuccessStatusAndNotEnityNotFound(t *testing.T) {
	statusCode := http.StatusBadRequest
	httpServer := setupAndStartHttpServer(http.MethodGet, fullPathWithID, statusCode)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	_, err := restClient.GetOne(testID, testPath)

	verifyFailedCallWithStatusCodeIsResponse(err, statusCode, t)
}

func TestShouldReturnDataForSuccessfulGetAllRequest(t *testing.T) {
	httpServer := setupAndStartHttpServerWithOKResponseCode(http.MethodGet, fullPathWithoutID)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	response, err := restClient.GetAll(testPath)

	verifySuccessfullGetOrPut(response, err, t)
}

func TestShouldReturnErrorMessageForGetAllRequestWhenStatusIsNotASuccessStatusAndNotEnityNotFound(t *testing.T) {
	statusCode := http.StatusBadRequest
	httpServer := setupAndStartHttpServer(http.MethodGet, fullPathWithoutID, statusCode)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	_, err := restClient.GetAll(testPath)

	verifyFailedCallWithStatusCodeIsResponse(err, statusCode, t)
}

func TestShouldReturnDataForSuccessfulPutRequest(t *testing.T) {
	httpServer := setupAndStartHttpServerWithOKResponseCode(http.MethodPut, fullPathWithID)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	response, err := restClient.Put(testDataObject{id: testID}, testPath)

	verifySuccessfullGetOrPut(response, err, t)
}

func TestShouldReturnErrorMessageForPutRequestWhenStatusIsNotASuccessStatusAndNotEnityNotFound(t *testing.T) {
	statusCode := http.StatusBadRequest
	httpServer := setupAndStartHttpServer(http.MethodPut, fullPathWithID, statusCode)
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
	httpServer := setupAndStartHttpServerWithOKResponseCode(http.MethodDelete, fullPathWithID)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	err := restClient.Delete(testID, testPath)

	if err != nil {
		t.Fatalf("Expected no error to be returned but got %s", err)
	}
}

func TestShouldReturnErrorMessageForDeleteRequestWhenStatusIsNotASuccessStatusAndNotEnityNotFound(t *testing.T) {
	statusCode := http.StatusBadRequest
	httpServer := setupAndStartHttpServer(http.MethodDelete, fullPathWithID, statusCode)
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
