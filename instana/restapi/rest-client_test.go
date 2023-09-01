package restapi_test

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/stretchr/testify/require"
)

const testPath = "/test"
const testID = "test-1234"
const testData = "testData"
const testPathWithID = testPath + "/" + testID

func TestShouldReturnDataForSuccessfulGetRequest(t *testing.T) {
	httpServer := setupAndStartHttpServerWithOKResponseCode(http.MethodGet, testPath)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	response, err := restClient.Get(testPath)

	verifySuccessResponseData(response, err, t)
}

func TestShouldReturnErrorMessageForGetRequestWhenStatusIsNotASuccessStatusAndNotEntityNotFound(t *testing.T) {
	statusCode := http.StatusBadRequest
	httpServer := setupAndStartHttpServer(http.MethodGet, testPath, statusCode)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	_, err := restClient.Get(testPath)

	verifyFailedCallWithStatusCodeIsResponse(err, statusCode, t)
}

func TestShouldReturnNotFoundErrorMessageForGetRequestWhenStatusIsNotEntityNotFound(t *testing.T) {
	httpServer := setupAndStartHttpServer(http.MethodGet, testPath, http.StatusNotFound)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	data, err := restClient.Get(testPath)

	verifyNotFoundResponse(data, err, t)
}

func TestShouldReturnDataForSuccessfulGetOneRequest(t *testing.T) {
	httpServer := setupAndStartHttpServerWithOKResponseCode(http.MethodGet, testPathWithID)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	response, err := restClient.GetOne(testID, testPath)

	verifySuccessResponseData(response, err, t)
}

func TestShouldReturnDataForSuccessfulGetOneRequestWhenResourcePathEndsWithASlash(t *testing.T) {
	httpServer := setupAndStartHttpServerWithOKResponseCode(http.MethodGet, testPathWithID)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	response, err := restClient.GetOne(testID, testPath+"/")

	verifySuccessResponseData(response, err, t)
}

func TestShouldReturnErrorMessageForGetOneRequestWhenStatusIsNotASuccessStatusAndNotEnityNotFound(t *testing.T) {
	statusCode := http.StatusBadRequest
	httpServer := setupAndStartHttpServer(http.MethodGet, testPathWithID, statusCode)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	_, err := restClient.GetOne(testID, testPath)

	verifyFailedCallWithStatusCodeIsResponse(err, statusCode, t)
}

func TestShouldReturnNotFoundErrorMessageForGetOneRequestWhenStatusIsNotEntityNotFound(t *testing.T) {
	httpServer := setupAndStartHttpServer(http.MethodGet, testPathWithID, http.StatusNotFound)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	data, err := restClient.GetOne(testID, testPath)

	verifyNotFoundResponse(data, err, t)
}

func TestShouldReturnDataForSuccessfulPostRequest(t *testing.T) {
	httpServer := setupAndStartHttpServerWithOKResponseCode(http.MethodPost, testPath)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	response, err := restClient.Post(testDataObject{id: testID}, testPath)

	verifySuccessResponseData(response, err, t)
}

func TestShouldReturnErrorMessageForPostRequestWhenStatusIsNotASuccessStatusAndNotEntityNotFound(t *testing.T) {
	statusCode := http.StatusBadRequest
	httpServer := setupAndStartHttpServer(http.MethodPost, testPath, statusCode)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	_, err := restClient.Post(testDataObject{id: testID}, testPath)

	verifyFailedCallWithStatusCodeIsResponse(err, statusCode, t)
}

func TestShouldReturnDataForSuccessfulPostWithIDRequest(t *testing.T) {
	httpServer := setupAndStartHttpServerWithOKResponseCode(http.MethodPost, testPathWithID)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	response, err := restClient.PostWithID(testDataObject{id: testID}, testPath)

	verifySuccessResponseData(response, err, t)
}

func TestShouldReturnErrorMessageForPostWithIDRequestWhenStatusIsNotASuccessStatusAndNotEntityNotFound(t *testing.T) {
	statusCode := http.StatusBadRequest
	httpServer := setupAndStartHttpServer(http.MethodPost, testPathWithID, statusCode)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	_, err := restClient.PostWithID(testDataObject{id: testID}, testPath)

	verifyFailedCallWithStatusCodeIsResponse(err, statusCode, t)
}

func TestShouldReturnDataForSuccessfulPutRequest(t *testing.T) {
	httpServer := setupAndStartHttpServerWithOKResponseCode(http.MethodPut, testPathWithID)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	response, err := restClient.Put(testDataObject{id: testID}, testPath)

	verifySuccessResponseData(response, err, t)
}

func TestShouldReturnErrorMessageForPutRequestWhenStatusIsNotASuccessStatusAndNotEntityNotFound(t *testing.T) {
	statusCode := http.StatusBadRequest
	httpServer := setupAndStartHttpServer(http.MethodPut, testPathWithID, statusCode)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	_, err := restClient.Put(testDataObject{id: testID}, testPath)

	verifyFailedCallWithStatusCodeIsResponse(err, statusCode, t)
}

func TestShouldReturnDataForSuccessfulPostByQueryRequestWhenNoQueryParametersAreProvided(t *testing.T) {
	queryParameters := map[string]string{}
	shouldReturnDataForSuccessfulPostByQueryRequest(t, queryParameters)
}

func TestShouldReturnDataForSuccessfulPostByQueryRequestWhenQueryParametersAreProvided(t *testing.T) {
	queryParameters := map[string]string{
		"a": "b",
		"c": "d",
	}
	shouldReturnDataForSuccessfulPostByQueryRequest(t, queryParameters)
}

func shouldReturnDataForSuccessfulPostByQueryRequest(t *testing.T, queryParameters map[string]string) {
	httpServer := setupAndStartHttpServerWithQueryParamerterCheck(http.MethodPost, testPath, queryParameters, http.StatusOK)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	response, err := restClient.PostByQuery(testPath, queryParameters)

	verifySuccessResponseData(response, err, t)
}

func TestShouldReturnErrorMessageForPostByQueryRequestWhenStatusIsNotASuccessStatusAndNotEntityNotFound(t *testing.T) {
	statusCode := http.StatusBadRequest
	queryParameters := map[string]string{
		"a": "b",
		"c": "d",
	}
	httpServer := setupAndStartHttpServerWithQueryParamerterCheck(http.MethodPost, testPath, queryParameters, statusCode)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	_, err := restClient.PostByQuery(testPath, queryParameters)

	verifyFailedCallWithStatusCodeIsResponse(err, statusCode, t)
}

func TestShouldReturnDataForSuccessfulPutByQueryRequestWhenNoQueryParametersAreProvided(t *testing.T) {
	queryParameters := map[string]string{}
	shouldReturnDataForSuccessfulPutByQueryRequest(t, queryParameters)
}

func TestShouldReturnDataForSuccessfulPutByQueryRequestWhenQueryParametersAreProvided(t *testing.T) {
	queryParameters := map[string]string{
		"a": "b",
		"c": "d",
	}
	shouldReturnDataForSuccessfulPutByQueryRequest(t, queryParameters)
}

func shouldReturnDataForSuccessfulPutByQueryRequest(t *testing.T, queryParameters map[string]string) {
	httpServer := setupAndStartHttpServerWithQueryParamerterCheck(http.MethodPut, testPathWithID, queryParameters, http.StatusOK)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	response, err := restClient.PutByQuery(testPath, testID, queryParameters)

	verifySuccessResponseData(response, err, t)
}

func TestShouldReturnEntityNotFoundErrorForPutByQueryRequestWhenStatusIsNotFound(t *testing.T) {
	statusCode := http.StatusNotFound
	queryParameters := map[string]string{
		"a": "b",
		"c": "d",
	}
	httpServer := setupAndStartHttpServerWithQueryParamerterCheck(http.MethodPut, testPathWithID, queryParameters, statusCode)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	_, err := restClient.PutByQuery(testPath, testID, queryParameters)

	verifyFailedCallWithStatusCodeIsResponse(err, statusCode, t)
}

func TestShouldReturnErrorMessageForPutByQueryRequestWhenStatusIsNotASuccessStatusAndNotEntityNotFound(t *testing.T) {
	statusCode := http.StatusBadRequest
	queryParameters := map[string]string{
		"a": "b",
		"c": "d",
	}
	httpServer := setupAndStartHttpServerWithQueryParamerterCheck(http.MethodPut, testPathWithID, queryParameters, statusCode)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	_, err := restClient.PutByQuery(testPath, testID, queryParameters)

	verifyFailedCallWithStatusCodeIsResponse(err, statusCode, t)
}

type testDataObject struct {
	id string
}

// GetIDForResourcePath implementation of InstanaDataObject
func (tdo testDataObject) GetIDForResourcePath() string {
	return tdo.id
}

func TestShouldReturnNothingForSuccessfulDeleteRequest(t *testing.T) {
	httpServer := setupAndStartHttpServerWithOKResponseCode(http.MethodDelete, testPathWithID)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	err := restClient.Delete(testID, testPath)

	require.Nil(t, err)
}

func TestShouldReturnErrorMessageForDeleteRequestWhenStatusIsNotASuccessStatusAndNotEntityNotFound(t *testing.T) {
	statusCode := http.StatusBadRequest
	httpServer := setupAndStartHttpServer(http.MethodDelete, testPathWithID, statusCode)
	defer httpServer.Close()

	restClient := createSut(httpServer)
	err := restClient.Delete(testID, testPath)

	verifyFailedCallWithStatusCodeIsResponse(err, statusCode, t)
}

func setupAndStartHttpServerWithOKResponseCode(httpMethod string, fullPath string) testutils.TestHTTPServer {
	return setupAndStartHttpServer(httpMethod, fullPath, 200)
}

func setupAndStartHttpServer(httpMethod string, fullPath string, statusCode int) testutils.TestHTTPServer {
	return doSetupAndStartHttpServer(httpMethod, fullPath, statusCode, func(r *http.Request) error { return nil })
}

func setupAndStartHttpServerWithQueryParamerterCheck(httpMethod string, fullPath string, queryParameters map[string]string, statusCode int) testutils.TestHTTPServer {
	return doSetupAndStartHttpServer(httpMethod, fullPath, statusCode, func(r *http.Request) error {
		for k, v := range queryParameters {
			val := r.URL.Query().Get(k)
			if val != v {
				return fmt.Errorf("Expected query parameter %s to be defined with value '%s'; current value is '%s'", k, v, val)
			}
		}
		return nil
	})
}

func doSetupAndStartHttpServer(httpMethod string, fullPath string, statusCode int, additionalChecks func(r *http.Request) error) testutils.TestHTTPServer {
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(httpMethod, fullPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		err := additionalChecks(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				fmt.Printf("failed to write error response; %s\n", err)
			}
			_, err = w.Write([]byte(testData))
			if err != nil {
				fmt.Printf("failed to write test data; %s\n", err)
			}
		} else {
			w.WriteHeader(statusCode)
			_, err = w.Write([]byte(testData))
			if err != nil {
				fmt.Printf("failed to write response; %s\n", err)
			}
		}
	})
	httpServer.Start()
	return httpServer
}

func createSut(httpServer testutils.TestHTTPServer) RestClient {
	return NewClient("api-token", fmt.Sprintf("localhost:%d", httpServer.GetPort()), true)
}

func verifyNotFoundResponse(data []byte, err error, t *testing.T) {
	require.Equal(t, ErrEntityNotFound, err)

	require.NotNil(t, data)
	require.GreaterOrEqual(t, 0, len(data))
}

func verifySuccessResponseData(response []byte, err error, t *testing.T) {
	require.Nil(t, err)
	responseString := string(response)
	require.Equal(t, testData, responseString)
}

func verifyFailedCallWithStatusCodeIsResponse(err error, statusCode int, t *testing.T) {
	require.NotNil(t, err)
	require.Contains(t, err.Error(), strconv.Itoa(statusCode))
}
