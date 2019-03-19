package services_test

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi/services"
	testutils "github.com/gessnerfl/terraform-provider-instana/test-utils"
)

const testPath = "/test"
const testID = "test-1234"
const testData = "testData"

func TestShouldReturnDataForSuccessfulGetOneRequest(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodGet, "/api"+testPath+"/"+testID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testData))
	})
	httpServer.Start()
	defer httpServer.Close()

	restClient := NewClient("api-token", fmt.Sprintf("localhost:%d", httpServer.GetPort()))
	response, err := restClient.GetOne(testID, testPath)

	if err != nil {
		t.Fatalf("Expected no error to be returned but got %s", err)
	}
	responseString := string(response)
	if responseString != testData {
		t.Fatalf("Expected test data to be returned but got %s", responseString)
	}
}

func TestShouldReturnErrorMessageForGetOneRequestWhenStatusIsNotASuccessStatusAndNotEnityNotFound(t *testing.T) {
	statusCode := http.StatusBadRequest
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodGet, "/api"+testPath+"/"+testID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(statusCode)
		w.Write([]byte(testData))
	})
	httpServer.Start()
	defer httpServer.Close()

	restClient := NewClient("api-token", fmt.Sprintf("localhost:%d", httpServer.GetPort()))
	_, err := restClient.GetOne(testID, testPath)

	if err == nil || !strings.Contains(err.Error(), strconv.Itoa(statusCode)) {
		t.Fatalf("Expected to receive error message with status Code %d", statusCode)
	}
}

func TestShouldReturnDataForSuccessfulGetAllRequest(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodGet, "/api"+testPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testData))
	})
	httpServer.Start()
	defer httpServer.Close()

	restClient := NewClient("api-token", fmt.Sprintf("localhost:%d", httpServer.GetPort()))
	response, err := restClient.GetAll(testPath)

	if err != nil {
		t.Fatalf("Expected no error to be returned but got %s", err)
	}
	responseString := string(response)
	if responseString != testData {
		t.Fatalf("Expected test data to be returned but got %s", responseString)
	}
}

func TestShouldReturnErrorMessageForGetAllRequestWhenStatusIsNotASuccessStatusAndNotEnityNotFound(t *testing.T) {
	statusCode := http.StatusBadRequest
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodGet, "/api"+testPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(statusCode)
		w.Write([]byte(testData))
	})
	httpServer.Start()
	defer httpServer.Close()

	restClient := NewClient("api-token", fmt.Sprintf("localhost:%d", httpServer.GetPort()))
	_, err := restClient.GetAll(testPath)

	if err == nil || !strings.Contains(err.Error(), strconv.Itoa(statusCode)) {
		t.Fatalf("Expected to receive error message with status Code %d", statusCode)
	}
}

func TestShouldReturnDataForSuccessfulPutRequest(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, "/api"+testPath+"/"+testID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testData))
	})
	httpServer.Start()
	defer httpServer.Close()

	restClient := NewClient("api-token", fmt.Sprintf("localhost:%d", httpServer.GetPort()))
	response, err := restClient.Put(testDataObject{id: testID}, testPath)

	if err != nil {
		t.Fatalf("Expected no error to be returned but got %s", err)
	}
	responseString := string(response)
	if responseString != testData {
		t.Fatalf("Expected test data to be returned but got %s", responseString)
	}
}

func TestShouldReturnErrorMessageForPutRequestWhenStatusIsNotASuccessStatusAndNotEnityNotFound(t *testing.T) {
	statusCode := http.StatusBadRequest
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, "/api"+testPath+"/"+testID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(statusCode)
		w.Write([]byte(testData))
	})
	httpServer.Start()
	defer httpServer.Close()

	restClient := NewClient("api-token", fmt.Sprintf("localhost:%d", httpServer.GetPort()))
	_, err := restClient.Put(testDataObject{id: testID}, testPath)

	if err == nil || !strings.Contains(err.Error(), strconv.Itoa(statusCode)) {
		t.Fatalf("Expected to receive error message with status Code %d", statusCode)
	}
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
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodDelete, "/api"+testPath+"/"+testID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
	})
	httpServer.Start()
	defer httpServer.Close()

	restClient := NewClient("api-token", fmt.Sprintf("localhost:%d", httpServer.GetPort()))
	err := restClient.Delete(testID, testPath)

	if err != nil {
		t.Fatalf("Expected no error to be returned but got %s", err)
	}
}

func TestShouldReturnErrorMessageForDeleteRequestWhenStatusIsNotASuccessStatusAndNotEnityNotFound(t *testing.T) {
	statusCode := http.StatusBadRequest
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodDelete, "/api"+testPath+"/"+testID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(statusCode)
	})
	httpServer.Start()
	defer httpServer.Close()

	restClient := NewClient("api-token", fmt.Sprintf("localhost:%d", httpServer.GetPort()))
	err := restClient.Delete(testID, testPath)

	if err == nil || !strings.Contains(err.Error(), strconv.Itoa(statusCode)) {
		t.Fatalf("Expected to receive error message with status Code %d", statusCode)
	}
}
