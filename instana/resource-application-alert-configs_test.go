package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"testing"
)


const applicationConfigID = "application-config-id"

const sampleJson = "{\n  \"name\": \"Calls are slower than usual\",\n  \"description\": \"Calls are slower or equal to 16 ms based on latency (90th).\",\n  \"boundaryScope\": \"INBOUND\",\n  \"applicationId\": \"9BM1O6DFSjKRwbCA_IZCgQ\",\n  \"severity\": 5,\n  \"triggering\": false,\n  \"tagFilters\": [\n    {\n      \"type\": \"TAG_FILTER\",\n      \"name\": \"endpoint.name\",\n      \"stringValue\": \"com.vorwerk.nwot.device.firmwareupdate.core.FragmentChecksumScheduler\",\n      \"numberValue\": null,\n      \"booleanValue\": null,\n      \"operator\": \"EQUALS\",\n      \"entity\": \"NOT_APPLICABLE\"\n    },\n    {\n      \"type\": \"TAG_FILTER\",\n      \"name\": \"service.name\",\n      \"stringValue\": \"firmware-admin-demo-eu\",\n      \"numberValue\": null,\n      \"booleanValue\": null,\n      \"operator\": \"EQUALS\",\n      \"entity\": \"NOT_APPLICABLE\"\n    }\n  ],\n  \"rule\": {\n    \"alertType\": \"slowness\",\n    \"aggregation\": \"P90\",\n    \"metricName\": \"latency\"\n  },\n  \"alertChannelIds\": [],\n  \"granularity\": 600000,\n  \"timeThreshold\": {\n    \"type\": \"violationsInSequence\",\n    \"timeWindow\": 600000\n  },\n  \"id\": \"Q7noes5nSLeBeNmCNavoKg\",\n  \"created\": 1603789346130,\n  \"readOnly\": false,\n  \"enabled\": true,\n  \"threshold\": {\n    \"type\": \"staticThreshold\",\n    \"lastUpdated\": 0,\n    \"operator\": \">=\",\n    \"value\": 16\n  }\n}"
func Test_mapStateToDataObjectForApplicationAlertConfigs(t *testing.T) {

	testHelper := NewTestHelper(t)
	resourceHandle := NewApplicationAlertConfigsResourceHandle()
	unmarshaller := restapi.NewApplicationAlertConfigsUnmarshaller()
	instanaData, _ := unmarshaller.Unmarshal([]byte(sampleJson))
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceHandle.UpdateState(resourceData, instanaData)

	print(resourceData)
	//
	//
	//
	//resourceData.SetId(applicationConfigID)
	//resourceData.Set(ApplicationAlertConfigsFieldDescription, "Test Description")
	//respurceData.Set(ApplicationAlertConfigs)
	//resourceData.Set(ApplicationConfigFieldMatchSpecification, defaultMatchSpecification)
	//resourceData.Set(ApplicationConfigFieldScope, string(restapi.ApplicationConfigScopeIncludeNoDownstream))
	//resourceData.Set(ApplicationConfigFieldBoundaryScope, string(restapi.BoundaryScopeAll))
	//
	//result, err := resourceHandle.MapStateToDataObject(resourceData, utils.NewResourceNameFormatter("prefix ", " suffix"))
	//
	//
	//mapStateToDataObjectForApplicationAlertConfigs(tt.args.d, tt.args.formatter)
	//

	}

