package api

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
)

//NewClient creates a new instance of the Instana REST API client
func NewClient(apiToken string, host string) RestClient {
	return &RestClientImpl{
		apiToken: apiToken,
		host:     host,
	}
}

//RestClientImpl is a helper class to interact with Instana REST API
type RestClientImpl struct {
	apiToken string
	host     string
}

var emptyResponse = make([]byte, 0)

//GetOne request the resource with the given ID
func (client *RestClientImpl) GetOne(id string, resourcePath string) ([]byte, error) {
	return client.get(client.buildResourceURL(resourcePath, id))
}

//GetAll requests all objects of the given resource
func (client *RestClientImpl) GetAll(resourcePath string) ([]byte, error) {
	return client.get(client.buildURL(resourcePath))
}

func (client *RestClientImpl) get(url string) ([]byte, error) {
	log.Infof("Call GET %s", url)
	resp, err := client.newResty().Get(url)
	return client.verifyResponseAndReturnBodyAsBytes(resp, err)
}

//Put executes a HTTP PUT request to create or update the given resource
func (client *RestClientImpl) Put(data InstanaDataObject, resourcePath string) ([]byte, error) {
	url := client.buildResourceURL(resourcePath, data.GetID())
	log.Infof("Call PUT %s", url)
	resp, err := client.newResty().SetBody(data).Put(url)
	return client.verifyResponseAndReturnBodyAsBytes(resp, err)
}

func (client *RestClientImpl) verifyResponseAndReturnBodyAsBytes(resp *resty.Response, err error) ([]byte, error) {
	if err != nil {
		return emptyResponse, fmt.Errorf("failed to call Instana API; status code = %d; status message = %s, %s", resp.StatusCode(), resp.Status(), err)
	}
	statusCode := resp.StatusCode()
	if statusCode < 200 || statusCode >= 300 {
		return emptyResponse, fmt.Errorf("failed to call Instana API; status code = %d; status message = %s", statusCode, resp.Status())
	}
	return resp.Body(), nil
}

//Delete executes a HTTP DELETE request to delete the resource with the given ID
func (client *RestClientImpl) Delete(resourceID string, resourceBasePath string) error {
	url := client.buildResourceURL(resourceBasePath, resourceID)
	log.Infof("Call DELETE %s", url)
	resp, err := client.newResty().Delete(url)
	return client.processUpdateResponse(resp, err)
}

func (client *RestClientImpl) processUpdateResponse(resp *resty.Response, err error) error {
	if err != nil {
		return fmt.Errorf("failed to call Instana API; status code = %d; status message = %s, %s", resp.StatusCode(), resp.Status(), err)
	}
	statusCode := resp.StatusCode()
	if statusCode < 200 || statusCode >= 300 {
		return fmt.Errorf("failed to call Instana API; status code = %d; status message = %s", statusCode, resp.Status())
	}
	return nil
}

func (client *RestClientImpl) newResty() *resty.Request {
	return resty.R().SetHeader("Accept", "application/json").SetHeader("Authorization", fmt.Sprintf("apiToken %s", client.apiToken))
}

func (client *RestClientImpl) buildResourceURL(resourceBasePath string, id string) string {
	pattern := "%s/%s"
	if strings.HasSuffix(resourceBasePath, "/") {
		pattern = "%s%s"
	}
	resourcePath := fmt.Sprintf(pattern, resourceBasePath, id)
	return client.buildURL(resourcePath)
}

func (client *RestClientImpl) buildURL(resourcePath string) string {
	apiPath := "api"
	if !strings.HasPrefix(resourcePath, "/") {
		apiPath = apiPath + "/"
	}
	return fmt.Sprintf("https://%s/%s%s", client.host, apiPath, resourcePath)
}
