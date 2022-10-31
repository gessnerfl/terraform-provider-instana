package instana_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"net/http"
	"testing"
)

func TestApplicationAlertConfig(t *testing.T) {
	commonTests := createApplicationAlertConfigTestFor("instana_application_alert_config", restapi.ApplicationAlertConfigsResourcePath, NewApplicationAlertConfigResourceHandle())
	commonTests.run(t)
}

const issue141Template = `
	resource "instana_application_alert_config" "issue141" {
  name              = "name %d"
  description       = "issue 141 description"
  boundary_scope    = "INBOUND"
  severity          = "warning"
  triggering        = false
  include_internal  = false
  include_synthetic = false
  alert_channel_ids = ["alert-channel-id-1"]
  granularity       = 1800000
  evaluation_type   = "PER_AP"

  application {
    application_id = "application-id-1"
    inclusive      = true
  }

  rule {
    throughput {
      metric_name = "calls"
      aggregation = "SUM"
    }
  }

  threshold {
    static {
      operator = "<="
      value    = 0
    }
  }

  time_threshold {
    violations_in_period {
      time_window = 1800000
      violations  = 1
    }
  }
}
`

const issue141JsonResponse = `
{
  "id": "%s",
  "name": "prefix name %d suffix",
  "description": "issue 141 description",
  "boundaryScope": "INBOUND",
  "applicationId": "application-id-1",
  "applications": {
    "application-id-1": {
      "applicationId": "application-id-1",
      "inclusive": true,
      "services": {}
    }
  },
  "severity": 5,
  "triggering": false,
  "tagFilterExpression": {
    "type": "EXPRESSION",
    "logicalOperator": "AND",
    "elements": []
  },
  "includeInternal": false,
  "includeSynthetic": false,
  "rule": {
    "alertType": "throughput",
    "metricName": "calls",
    "aggregation": "SUM"
  },
  "threshold": {
    "type": "staticThreshold",
    "operator": "<=",
    "value": 0.0,
    "lastUpdated": 0
  },
  "alertChannelIds": [
    "alert-channel-id-1"
  ],
  "granularity": 1800000,
  "timeThreshold": {
    "type": "violationsInPeriod",
    "timeWindow": 1800000,
    "violations": 1
  },
  "evaluationType": "PER_AP",
  "customPayloadFields": [],
  "created": 1665254479160,
  "readOnly": false,
  "enabled": true,
  "derivedFromGlobalAlert": false
}
`

func TestIssue141(t *testing.T) {
	id := RandomID()
	httpServer := testutils.NewTestHTTPServer()
	resourceRestAPIPath := restapi.ApplicationAlertConfigsResourcePath
	resourceInstanceRestAPIPath := resourceRestAPIPath + "/{internal-id}"
	httpServer.AddRoute(http.MethodPost, resourceRestAPIPath, func(w http.ResponseWriter, r *http.Request) {
		config := &restapi.ApplicationAlertConfig{}
		err := json.NewDecoder(r.Body).Decode(config)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			r.Write(bytes.NewBufferString("Failed to get request"))
		} else {
			config.ID = id
			w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(config)
		}
	})
	httpServer.AddRoute(http.MethodPost, resourceInstanceRestAPIPath, func(w http.ResponseWriter, r *http.Request) {
		testutils.EchoHandlerFunc(w, r)
	})
	httpServer.AddRoute(http.MethodDelete, resourceInstanceRestAPIPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, resourceInstanceRestAPIPath, func(w http.ResponseWriter, r *http.Request) {
		modCount := httpServer.GetCallCount(http.MethodPost, resourceRestAPIPath+"/"+id)
		json := fmt.Sprintf(issue141JsonResponse, id, modCount)
		w.Header().Set(contentType, r.Header.Get(contentType))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()
	terraformResourceInstanceName := "instana_application_alert_config.issue141"
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testProviderFactory,
		Steps: []resource.TestStep{
			createIssue141TestStep(terraformResourceInstanceName, httpServer.GetPort(), 0, id),
			testStepImportWithCustomID(terraformResourceInstanceName, id),
			createIssue141TestStep(terraformResourceInstanceName, httpServer.GetPort(), 1, id),
			testStepImportWithCustomID(terraformResourceInstanceName, id),
		},
	})
}

func createIssue141TestStep(terraformResourceInstanceName string, httpPort int, iteration int, id string) resource.TestStep {
	application1ApplicationId := fmt.Sprintf("%s.%d.%s", ApplicationAlertConfigFieldApplications, 0, ApplicationAlertConfigFieldApplicationsApplicationID)
	application1Inclusive := fmt.Sprintf("%s.%d.%s", ApplicationAlertConfigFieldApplications, 0, ApplicationAlertConfigFieldApplicationsInclusive)
	ruleThroughputMetricName := fmt.Sprintf("%s.%d.%s.%d.%s", ApplicationAlertConfigFieldRule, 0, ApplicationAlertConfigFieldRuleThroughput, 0, ApplicationAlertConfigFieldRuleMetricName)
	ruleThroughputAggregation := fmt.Sprintf("%s.%d.%s.%d.%s", ApplicationAlertConfigFieldRule, 0, ApplicationAlertConfigFieldRuleThroughput, 0, ApplicationAlertConfigFieldRuleAggregation)
	thresholdStaticOperator := fmt.Sprintf("%s.%d.%s.%d.%s", ResourceFieldThreshold, 0, ResourceFieldThresholdStatic, 0, ResourceFieldThresholdOperator)
	thresholdStaticValue := fmt.Sprintf("%s.%d.%s.%d.%s", ResourceFieldThreshold, 0, ResourceFieldThresholdStatic, 0, ResourceFieldThresholdStaticValue)
	timeThresholdViolationsInPeriodTimeWindow := fmt.Sprintf("%s.%d.%s.%d.%s", ApplicationAlertConfigFieldTimeThreshold, 0, ApplicationAlertConfigFieldTimeThresholdViolationsInPeriod, 0, ApplicationAlertConfigFieldTimeThresholdTimeWindow)
	timeThresholdViolationsInPeriodViolations := fmt.Sprintf("%s.%d.%s.%d.%s", ApplicationAlertConfigFieldTimeThreshold, 0, ApplicationAlertConfigFieldTimeThresholdViolationsInPeriod, 0, ApplicationAlertConfigFieldTimeThresholdViolationsInPeriodViolations)
	return resource.TestStep{
		Config: appendProviderConfig(fmt.Sprintf(issue141Template, iteration), httpPort),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(terraformResourceInstanceName, "id", id),
			resource.TestCheckResourceAttr(terraformResourceInstanceName, ApplicationAlertConfigFieldName, formatResourceName(iteration)),
			resource.TestCheckResourceAttr(terraformResourceInstanceName, ApplicationAlertConfigFieldFullName, formatResourceFullName(iteration)),
			resource.TestCheckResourceAttr(terraformResourceInstanceName, ApplicationAlertConfigFieldDescription, "issue 141 description"),
			resource.TestCheckResourceAttr(terraformResourceInstanceName, ApplicationAlertConfigFieldBoundaryScope, string(restapi.BoundaryScopeInbound)),
			resource.TestCheckResourceAttr(terraformResourceInstanceName, ApplicationAlertConfigFieldSeverity, restapi.SeverityWarning.GetTerraformRepresentation()),
			resource.TestCheckResourceAttr(terraformResourceInstanceName, ApplicationAlertConfigFieldTriggering, falseAsString),
			resource.TestCheckResourceAttr(terraformResourceInstanceName, ApplicationAlertConfigFieldIncludeInternal, falseAsString),
			resource.TestCheckResourceAttr(terraformResourceInstanceName, ApplicationAlertConfigFieldIncludeSynthetic, falseAsString),
			resource.TestCheckResourceAttr(terraformResourceInstanceName, ApplicationAlertConfigFieldAlertChannelIDs+".0", "alert-channel-id-1"),
			resource.TestCheckResourceAttr(terraformResourceInstanceName, ApplicationAlertConfigFieldGranularity, "1800000"),
			resource.TestCheckResourceAttr(terraformResourceInstanceName, ApplicationAlertConfigFieldEvaluationType, string(restapi.EvaluationTypePerApplication)),
			resource.TestCheckResourceAttr(terraformResourceInstanceName, application1ApplicationId, "application-id-1"),
			resource.TestCheckResourceAttr(terraformResourceInstanceName, application1Inclusive, trueAsString),
			resource.TestCheckResourceAttr(terraformResourceInstanceName, ruleThroughputMetricName, "calls"),
			resource.TestCheckResourceAttr(terraformResourceInstanceName, ruleThroughputAggregation, "SUM"),
			resource.TestCheckResourceAttr(terraformResourceInstanceName, thresholdStaticOperator, "<="),
			resource.TestCheckResourceAttr(terraformResourceInstanceName, thresholdStaticValue, "0"),
			resource.TestCheckResourceAttr(terraformResourceInstanceName, timeThresholdViolationsInPeriodTimeWindow, "1800000"),
			resource.TestCheckResourceAttr(terraformResourceInstanceName, timeThresholdViolationsInPeriodViolations, "1"),
		),
	}
}
