package restapi

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	resty "gopkg.in/resty.v1"
)

//ErrEntityNotFound error message which is returned when the entity cannot be found at the server
var ErrEntityNotFound = errors.New("Failed to get resource from Instana API. 404 - Resource not found")

//RestClient interface to access REST resources of the Instana API
type RestClient interface {
	GetOne(id string, resourcePath string) ([]byte, error)
	Put(data InstanaDataObject, resourcePath string) ([]byte, error)
	Delete(resourceID string, resourceBasePath string) error
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

//NewClient creates a new instance of the Instana REST API client
func NewClient(apiToken string, host string) RestClient {
	restyClient := resty.New()

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

//GetOne request the resource with the given ID
func (client *restClientImpl) GetOne(id string, resourcePath string) ([]byte, error) {
	url := client.buildResourceURL(resourcePath, id)
	req := client.createRequest()
	return client.executeRequest(resty.MethodGet, url, req)
}

//Put executes a HTTP PUT request to create or update the given resource
func (client *restClientImpl) Put(data InstanaDataObject, resourcePath string) ([]byte, error) {
	url := client.buildResourceURL(resourcePath, data.GetID())
	req := client.createRequest().SetHeader("Content-Type", "application/json; charset=utf-8").SetBody(data)
	return client.executeRequestWithThrottling(resty.MethodPut, url, req)
}

//Delete executes a HTTP DELETE request to delete the resource with the given ID
func (client *restClientImpl) Delete(resourceID string, resourceBasePath string) error {
	url := client.buildResourceURL(resourceBasePath, resourceID)
	req := client.createRequest()
	_, err := client.executeRequestWithThrottling(resty.MethodDelete, url, req)
	return err
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
	throttle := time.Tick(client.throttleRate)
	for req := range client.throttledRequests {
		<-throttle
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
	log.Infof("Call %s %s", method, url)
	resp, err := req.Execute(method, url)
	if err != nil {
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
