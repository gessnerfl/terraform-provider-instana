package restapi_test

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/stretchr/testify/assert"
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

func TestShouldReturnDataForSuccessfulGetOneRequestWhenResourcePathEndsWithASlash(t *testing.T) {
	httpServer := setupAndStartHttpServerWithOKResponseCode(http.MethodGet, testPathWithID)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	response, err := restClient.GetOne(testID, testPath+"/")

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

	assert.Nil(t, err)
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

func createSut(httpServer *testutils.TestHTTPServer) RestClient {
	return NewClient("api-token", fmt.Sprintf("localhost:%d", httpServer.GetPort()))
}

func verifyNotFoundResponse(data []byte, err error, t *testing.T) {
	assert.Equal(t, ErrEntityNotFound, err)

	assert.NotNil(t, data)
	assert.GreaterOrEqual(t, 0, len(data))
}

func verifySuccessfullGetOrPut(response []byte, err error, t *testing.T) {
	assert.Nil(t, err)
	responseString := string(response)
	assert.Equal(t, testData, responseString)
}

func verifyFailedCallWithStatusCodeIsResponse(err error, statusCode int, t *testing.T) {
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), strconv.Itoa(statusCode))
}
