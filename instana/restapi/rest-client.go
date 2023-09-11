package restapi

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	resty "gopkg.in/resty.v1"
)

// ErrEntityNotFound error message which is returned when the entity cannot be found at the server
var ErrEntityNotFound = errors.New("failed to get resource from Instana API. 404 - Resource not found")

const contentTypeHeader = "Content-Type"
const encodingApplicationJSON = "application/json; charset=utf-8"

// RestClient interface to access REST resources of the Instana API
type RestClient interface {
	Get(resourcePath string) ([]byte, error)
	GetOne(id string, resourcePath string) ([]byte, error)
	Post(data InstanaDataObject, resourcePath string) ([]byte, error)
	PostWithID(data InstanaDataObject, resourcePath string) ([]byte, error)
	Put(data InstanaDataObject, resourcePath string) ([]byte, error)
	Delete(resourceID string, resourceBasePath string) error
	PostByQuery(resourcePath string, queryParams map[string]string) ([]byte, error)
	PutByQuery(resourcePath string, is string, queryParams map[string]string) ([]byte, error)
}

type apiRequest struct {
	method          string
	url             string
	request         resty.Request
	responseChannel chan *apiResponse
	ctx             context.Context
}

type apiResponse struct {
	data []byte
	err  error
}

// NewClient creates a new instance of the Instana REST API client
func NewClient(apiToken string, host string, skipTlsVerification bool) RestClient {
	restyClient := resty.New()
	if skipTlsVerification {
		restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) //nolint:gosec
	}

	throttleRate := time.Second / 5 //5 write requests per second
	throttledRequests := make(chan *apiRequest, 1000)
	client := &restClientImpl{
		apiToken:          apiToken,
		host:              host,
		restyClient:       restyClient,
		throttledRequests: throttledRequests,
		throttleRate:      throttleRate,
	}

	go client.processThrottledRequests()
	return client
}

type restClientImpl struct {
	apiToken          string
	host              string
	restyClient       *resty.Client
	throttledRequests chan *apiRequest
	throttleRate      time.Duration
}

var emptyResponse = make([]byte, 0)

// Get request data via HTTP GET for the given resourcePath
func (client *restClientImpl) Get(resourcePath string) ([]byte, error) {
	url := client.buildURL(resourcePath)
	req := client.createRequest()
	return client.executeRequest(resty.MethodGet, url, req)
}

// GetOne request the resource with the given ID
func (client *restClientImpl) GetOne(id string, resourcePath string) ([]byte, error) {
	url := client.buildResourceURL(resourcePath, id)
	req := client.createRequest()
	return client.executeRequest(resty.MethodGet, url, req)
}

// Post executes a HTTP PUT request to create or update the given resource
func (client *restClientImpl) Post(data InstanaDataObject, resourcePath string) ([]byte, error) {
	url := client.buildURL(resourcePath)
	req := client.createRequest().SetHeader(contentTypeHeader, encodingApplicationJSON).SetBody(data)
	return client.executeRequestWithThrottling(resty.MethodPost, url, req)
}

// PostWithID executes a HTTP PUT request to create or update the given resource using the ID from the InstanaDataObject in the resource path
func (client *restClientImpl) PostWithID(data InstanaDataObject, resourcePath string) ([]byte, error) {
	url := client.buildResourceURL(resourcePath, data.GetIDForResourcePath())
	req := client.createRequest().SetHeader(contentTypeHeader, encodingApplicationJSON).SetBody(data)
	return client.executeRequestWithThrottling(resty.MethodPost, url, req)
}

// Put executes a HTTP PUT request to create or update the given resource
func (client *restClientImpl) Put(data InstanaDataObject, resourcePath string) ([]byte, error) {
	url := client.buildResourceURL(resourcePath, data.GetIDForResourcePath())
	req := client.createRequest().SetHeader(contentTypeHeader, encodingApplicationJSON).SetBody(data)
	return client.executeRequestWithThrottling(resty.MethodPut, url, req)
}

// Delete executes a HTTP DELETE request to delete the resource with the given ID
func (client *restClientImpl) Delete(resourceID string, resourceBasePath string) error {
	url := client.buildResourceURL(resourceBasePath, resourceID)
	req := client.createRequest()
	_, err := client.executeRequestWithThrottling(resty.MethodDelete, url, req)
	return err
}

// PostByQuery executes a HTTP POST request to create the resource by providing the data a query parameters
func (client *restClientImpl) PostByQuery(resourcePath string, queryParams map[string]string) ([]byte, error) {
	url := client.buildURL(resourcePath)
	req := client.createRequest()
	client.appendQueryParameters(req, queryParams)
	return client.executeRequest(resty.MethodPost, url, req)
}

// PutByQuery executes a HTTP PUT request to update the resource with the given ID by providing the data a query parameters
func (client *restClientImpl) PutByQuery(resourcePath string, id string, queryParams map[string]string) ([]byte, error) {
	url := client.buildResourceURL(resourcePath, id)
	req := client.createRequest()
	client.appendQueryParameters(req, queryParams)
	return client.executeRequest(resty.MethodPut, url, req)
}

func (client *restClientImpl) createRequest() *resty.Request {
	return client.restyClient.R().SetHeader("Accept", "application/json").SetHeader("Authorization", fmt.Sprintf("apiToken %s", client.apiToken))
}

func (client *restClientImpl) executeRequestWithThrottling(method string, url string, req *resty.Request) ([]byte, error) {
	responseChannel := make(chan *apiResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer close(responseChannel)
	defer cancel()

	client.throttledRequests <- &apiRequest{
		method:          method,
		url:             url,
		request:         *req,
		ctx:             ctx,
		responseChannel: responseChannel,
	}

	select {
	case r := <-responseChannel:
		return r.data, r.err
	case <-time.After(30 * time.Second):
		return nil, errors.New("API request timed out")
	}
}

func (client *restClientImpl) processThrottledRequests() {
	timer := time.NewTimer(client.throttleRate)
	for req := range client.throttledRequests {
		<-timer.C
		go client.handleThrottledAPIRequest(req)
	}
}

func (client *restClientImpl) handleThrottledAPIRequest(req *apiRequest) {
	data, err := client.executeRequest(req.method, req.url, &req.request)
	responseMessage := &apiResponse{
		data: data,
		err:  err,
	}
	select {
	case <-req.ctx.Done():
		return
	default:
		req.responseChannel <- responseMessage
	}
}

func (client *restClientImpl) executeRequest(method string, url string, req *resty.Request) ([]byte, error) {
	log.Printf("[DEBUG] Call %s %s\n", method, url)
	resp, err := req.Execute(method, url)
	if err != nil {
		if resp == nil {
			return emptyResponse, fmt.Errorf("failed to send HTTP %s request to Instana API; %s", method, err)
		}
		return emptyResponse, fmt.Errorf("failed to send HTTP %s request to Instana API; status code = %d; status message = %s; Headers %s, %s", method, resp.StatusCode(), resp.Status(), resp.Header(), err)
	}
	statusCode := resp.StatusCode()
	if statusCode == 404 {
		return emptyResponse, ErrEntityNotFound
	}
	if statusCode < 200 || statusCode >= 300 {
		return emptyResponse, fmt.Errorf("failed to send HTTP %s request to Instana API; status code = %d; status message = %s; Headers %s\nBody: %s", method, statusCode, resp.Status(), resp.Header(), resp.Body())
	}
	return resp.Body(), nil
}

func (client *restClientImpl) appendQueryParameters(req *resty.Request, queryParams map[string]string) {
	for k, v := range queryParams {
		req.QueryParam.Add(k, v)
	}
}

func (client *restClientImpl) buildResourceURL(resourceBasePath string, id string) string {
	pattern := "%s/%s"
	if strings.HasSuffix(resourceBasePath, "/") {
		pattern = "%s%s"
	}
	resourcePath := fmt.Sprintf(pattern, resourceBasePath, id)
	return client.buildURL(resourcePath)
}

func (client *restClientImpl) buildURL(resourcePath string) string {
	return fmt.Sprintf("https://%s%s", client.host, resourcePath)
}
