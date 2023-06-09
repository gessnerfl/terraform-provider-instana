package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceInstanaSyntheticTest the name of the terraform-provider-instana resource to manage synthetic tests
const ResourceInstanaSyntheticMonitor = "instana_synthetic_monitor"

const (
	//SyntheticMonitorFieldLabel constant value for the schema field label
	SyntheticMonitorFieldLabel = "label"
	//SyntheticMonitorFieldDescription constant value for the computed schema field description
	SyntheticMonitorFieldDescription = "description"
	//SyntheticMonitorFieldActive constant value for the schema field active
	SyntheticMonitorFieldActive = "active"
	//SyntheticMonitorFieldApplicationID constant value for the schema field application_id
	SyntheticMonitorFieldApplicationID = "application_id"
	//SyntheticMonitorFieldConfiguration constant value for the schema field configuration
	SyntheticMonitorFieldConfiguration = "configuration"
	//SyntheticMonitorFieldCustomProperties constant value for the schema field custom_properties
	SyntheticMonitorFieldCustomProperties = "custom_properties"
	//SyntheticMonitorFieldLocations constant value for the schema field locations
	SyntheticMonitorFieldLocations = "locations"
	//SyntheticMonitorFieldPlaybackMode constant value for the schema field playback_mode
	SyntheticMonitorFieldPlaybackMode = "playback_mode"
	//SyntheticMonitorFieldTestFrequency constant value for the schema field test_frequency
	SyntheticMonitorFieldTestFrequency = "test_frequency"
	//SyntheticMonitorFieldConfigMarkSyntheticCall constant value for the schema field configuration.mark_synthetic_call
	SyntheticMonitorFieldConfigMarkSyntheticCall = "mark_synthetic_call"
	//SyntheticMonitorFieldConfigRetries constant value for the schema field configuration.retries
	SyntheticMonitorFieldConfigRetries = "retries"
	//SyntheticMonitorFieldConfigRetryInterval constant value for the schema field configuration.retry_interval
	SyntheticMonitorFieldConfigRetryInterval = "retry_interval"
	//SyntheticMonitorFieldConfigSyntheticType constant value for the schema field configuration.synthetic_type
	SyntheticMonitorFieldConfigSyntheticType = "synthetic_type"
	//SyntheticMonitorFieldConfigTimeout constant value for the schema field configuration.timeout
	SyntheticMonitorFieldConfigTimeout = "timeout"
	//SyntheticMonitorFieldConfigUrl constant value for the schema field configuration.url
	SyntheticMonitorFieldConfigUrl = "url"
	//SyntheticMonitorFieldConfigOperation constant value for the schema field configuration.operation
	SyntheticMonitorFieldConfigOperation = "operation"
	//SyntheticMonitorFieldConfigScript constant value for the schema field configuration.script
	SyntheticMonitorFieldConfigScript = "script"
)

// NewSyntheticTestResourceHandle creates the resource handle Synthetic Tests
func NewSyntheticTestResourceHandle() ResourceHandle {
	return &syntheticMonitorResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaSyntheticMonitor,
			Schema: map[string]*schema.Schema{
				SyntheticMonitorFieldLabel: {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "Friendly name of the Synthetic test",
					ValidateFunc: validation.StringLenBetween(0, 512),
				},
				SyntheticMonitorFieldDescription: {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "The description of the Synthetic test",
					ValidateFunc: validation.StringLenBetween(0, 512),
				},
				SyntheticMonitorFieldActive: {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     true,
					Description: "Indicates if the Synthetic test is started or not",
				},
				SyntheticMonitorFieldApplicationID: {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Unique identifier of the Application Perspective",
				},
				SyntheticMonitorFieldConfiguration: {
					Type:        schema.TypeList,
					MinItems:    1,
					MaxItems:    1,
					Required:    true,
					Description: "The configuration of the synthetic alert",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							SyntheticMonitorFieldConfigMarkSyntheticCall: {
								Type:        schema.TypeBool,
								Required:    true,
								Description: "Flag used to control if HTTP calls will be marked as synthetic calls",
							},
							SyntheticMonitorFieldConfigRetries: {
								Type:         schema.TypeInt,
								Optional:     true,
								Description:  "Indicates how many attempts will be allowed to get a successful connection",
								ValidateFunc: validation.IntBetween(0, 2),
							},
							SyntheticMonitorFieldConfigRetryInterval: {
								Type:         schema.TypeInt,
								Optional:     true,
								Description:  "The time interval between retries in seconds",
								ValidateFunc: validation.IntBetween(1, 10),
							},
							SyntheticMonitorFieldConfigSyntheticType: {
								Type:         schema.TypeString,
								Required:     true,
								Description:  "The type of the Synthetic test",
								ValidateFunc: validation.StringInSlice([]string{"HTTPAction", "HTTPScript", "BrowserScript", "WebpageAction", "WebpageScript", "DNSAction"}, true),
							},
							SyntheticMonitorFieldConfigTimeout: {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "The timeout to be used by the PoP playback engines running the test",
							},
							SyntheticMonitorFieldConfigUrl: {
								Type:         schema.TypeString,
								Optional:     true,
								Description:  "The URL is being tested",
								ValidateFunc: validation.IsURLWithHTTPorHTTPS,
							},
							SyntheticMonitorFieldConfigOperation: {
								Type:         schema.TypeString,
								Optional:     true,
								Default:      "GET",
								Description:  "The HTTP operation",
								ValidateFunc: validation.StringInSlice([]string{"GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT", "DELETE"}, true),
							},
							SyntheticMonitorFieldConfigScript: {
								Type:        schema.TypeString,
								Optional:    true,
								Description: " The Javascript content in plain text",
							},
						},
					},
				},
				SyntheticMonitorFieldCustomProperties: {
					Type:        schema.TypeMap,
					Optional:    true,
					Description: "Name/value pairs to provide additional information of the Synthetic test",
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				SyntheticMonitorFieldLocations: {
					Type:        schema.TypeSet,
					Required:    true,
					Description: "Array of the PoP location IDs",
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				SyntheticMonitorFieldPlaybackMode: {
					Type:         schema.TypeString,
					Optional:     true,
					Default:      "Simultaneous",
					Description:  "Defines how the Synthetic test should be executed across multiple PoPs",
					ValidateFunc: validation.StringInSlice([]string{"Simultaneous", "Staggered"}, true),
				},
				SyntheticMonitorFieldTestFrequency: {
					Type:         schema.TypeInt,
					Optional:     true,
					Default:      15,
					Description:  "How often the playback for a Synthetic test is scheduled",
					ValidateFunc: validation.IntBetween(1, 120),
				},
			},
			SchemaVersion: 0,
		},
	}
}

type syntheticMonitorResource struct {
	metaData ResourceMetaData
}

func (r *syntheticMonitorResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *syntheticMonitorResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{}
}

func (r *syntheticMonitorResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource {
	return api.SyntheticMonitorConfig()
}

func (r *syntheticMonitorResource) SetComputedFields(d *schema.ResourceData) {
	// d.Set(SyntheticTestFieldApplicationID, RandomID())
}

func (r *syntheticMonitorResource) UpdateState(d *schema.ResourceData, obj restapi.InstanaDataObject, formatter utils.ResourceNameFormatter) error {
	syntheticTest := obj.(*restapi.SyntheticMonitor)
	d.SetId(syntheticTest.ID)
	d.Set(SyntheticMonitorFieldConfiguration, r.mapConfigurationToSchema(syntheticTest))
	return nil
}

func (r *syntheticMonitorResource) MapStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	return &restapi.SyntheticMonitor{
		ID:               d.Id(),
		Label:            d.Get(SyntheticMonitorFieldLabel).(string),
		Description:      d.Get(SyntheticMonitorFieldDescription).(string),
		Active:           d.Get(SyntheticMonitorFieldActive).(bool),
		Configuration:    r.mapConfigurationFromSchema(d),
		CustomProperties: d.Get(SyntheticMonitorFieldCustomProperties).(map[string]interface{}),
		Locations:        ReadStringSetParameterFromResource(d, SyntheticMonitorFieldLocations),
		PlaybackMode:     d.Get(SyntheticMonitorFieldPlaybackMode).(string),
		TestFrequency:    int32(d.Get(SyntheticMonitorFieldTestFrequency).(int)),
	}, nil
}

func (r *syntheticMonitorResource) mapConfigurationToSchema(config *restapi.SyntheticMonitor) []map[string]interface{} {
	configuration := make(map[string]interface{})
	configuration[SyntheticMonitorFieldConfigMarkSyntheticCall] = config.Configuration.MarkSyntheticCall
	configuration[SyntheticMonitorFieldConfigSyntheticType] = config.Configuration.SyntheticType
	configuration[SyntheticMonitorFieldConfigTimeout] = config.Configuration.Timeout
	configuration[SyntheticMonitorFieldConfigRetries] = config.Configuration.Retries
	configuration[SyntheticMonitorFieldConfigRetryInterval] = config.Configuration.RetryInterval
	configuration[SyntheticMonitorFieldConfigUrl] = config.Configuration.URL
	configuration[SyntheticMonitorFieldConfigScript] = config.Configuration.Script
	configuration[SyntheticMonitorFieldConfigOperation] = config.Configuration.Operation
	result := make([]map[string]interface{}, 1)
	result[0] = configuration
	return result
}

func (r *syntheticMonitorResource) mapConfigurationFromSchema(d *schema.ResourceData) restapi.SyntheticTestConfig {
	syntheticTestConfigurationSlice := d.Get(SyntheticMonitorFieldConfiguration).([]interface{})
	syntheticTestConfig := syntheticTestConfigurationSlice[0].(map[string]interface{})

	return restapi.SyntheticTestConfig{
		MarkSyntheticCall: syntheticTestConfig[SyntheticMonitorFieldConfigMarkSyntheticCall].(bool),
		Retries:           int32(syntheticTestConfig[SyntheticMonitorFieldConfigRetries].(int)),
		RetryInterval:     int32(syntheticTestConfig[SyntheticMonitorFieldConfigRetryInterval].(int)),
		SyntheticType:     syntheticTestConfig[SyntheticMonitorFieldConfigSyntheticType].(string),
		Timeout:           syntheticTestConfig[SyntheticMonitorFieldConfigTimeout].(string),
		URL:               syntheticTestConfig[SyntheticMonitorFieldConfigUrl].(string),
		Operation:         syntheticTestConfig[SyntheticMonitorFieldConfigOperation].(string),
		Script:            syntheticTestConfig[SyntheticMonitorFieldConfigScript].(string),
	}
}
