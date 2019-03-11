package instanaapi

import (
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"gopkg.in/resty.v1"
)

//New constructs a new instance of API
func New(apiToken string, tenant string) *API {
	return &API{
		apiToken: apiToken,
		tenant:   tenant,
	}
}

//API is the GO representation of the Instana API
type API struct {
	apiToken string
	tenant   string
}

type apiCallable func(*resty.Request) (*resty.Response, error)

func (api *API) newResty() *resty.Request {
	return resty.R().SetHeader("Accept", "application/json").SetHeader("Authorization", fmt.Sprintf("apiToken %s", api.apiToken))
}

func (api *API) buildURL(pathComponents ...string) string {
	path := strings.Join(pathComponents, "/")
	return fmt.Sprintf("https://%s-tis.instana.io/api/%s", api.tenant, path)
}

func (api *API) callAPI(callable apiCallable) (*resty.Response, error) {
	resp, err := callable(api.newResty())
	if err != nil {
		return resp, fmt.Errorf("failed to call Instana API; status code = %d; status message = %s", resp.StatusCode(), resp.Status())
	}

	statusCode := resp.StatusCode()
	if statusCode < 200 || statusCode >= 300 {
		return resp, fmt.Errorf("failed to call Instana API; status code = %d; status message = %s", statusCode, resp.Status())
	}

	log.Infof("Successfully called Instana REST API; path = %s, statusCode = %d, body = %s", resp.Request.URL, statusCode, resp.String())
	return resp, nil
}

//Rule is the representation of a custom rule in Instana
type Rule struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	EntityType        string  `json:"entityType"`
	MetricName        string  `json:"metricName"`
	Rollup            int     `json:"rollup"`
	Window            int     `json:"window"`
	Aggregation       string  `json:"aggregation"`
	ConditionOperator string  `json:"conditionOperator"`
	ConditionValue    float32 `json:"conditionValue"`
}

//DeleteRule deletes a custom rule
func (api *API) DeleteRule(rule Rule) error {
	return api.DeleteRuleByID(rule.ID)
}

//DeleteRuleByID deletes a custom rule by its ID
func (api *API) DeleteRuleByID(ruleID string) error {
	_, err := api.callAPI(func(req *resty.Request) (*resty.Response, error) {
		return req.Delete(api.buildURL("rules", ruleID))
	})
	return err
}

//UpsetRule creates or updates a custom rule
func (api *API) UpsetRule(rule Rule) error {
	_, err := api.callAPI(func(req *resty.Request) (*resty.Response, error) {
		return req.SetBody(rule).Put(api.buildURL("rules", rule.ID))
	})
	return err
}

//GetRules returns all configured custom rules from Instana API
func (api *API) GetRules() ([]Rule, error) {
	rules := make([]Rule, 0)

	resp, err := api.callAPI(func(req *resty.Request) (*resty.Response, error) { return req.Get(api.buildURL("rules")) })
	if err != nil {
		return rules, err
	}

	responseData := resp.Body()
	if err := json.Unmarshal(responseData, &rules); err != nil {
		return rules, fmt.Errorf("failed to parse json %s; %s", resp.String(), err)
	}

	return rules, nil
}
