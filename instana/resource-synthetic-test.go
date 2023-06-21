package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceInstanaSyntheticTest the name of the terraform-provider-instana resource to manage synthetic tests
const ResourceInstanaSyntheticTest = "instana_synthetic_test"

const (
	//SyntheticTestFieldLabel constant value for the schema field label
	SyntheticTestFieldLabel = "label"
	//SyntheticTestFieldDescription constant value for the computed schema field description
	SyntheticTestFieldDescription = "description"
	//SyntheticTestFieldActive constant value for the schema field active
	SyntheticTestFieldActive = "active"
	//SyntheticTestFieldConfiguration constant value for the schema field configuration
	SyntheticTestFieldConfiguration = "configuration"
	//SyntheticTestFieldCustomProperties constant value for the schema field custom_properties
	SyntheticTestFieldCustomProperties = "custom_properties"
	//SyntheticTestFieldLocations constant value for the schema field locations
	SyntheticTestFieldLocations = "locations"
	//SyntheticTestFieldPlaybackMode constant value for the schema field playback_mode
	SyntheticTestFieldPlaybackMode = "playback_mode"
	//SyntheticTestFieldTestFrequency constant value for the schema field test_frequency
	SyntheticTestFieldTestFrequency = "test_frequency"
	//SyntheticTestFieldConfigMarkSyntheticCall constant value for the schema field configuration.mark_synthetic_call
	SyntheticTestFieldConfigMarkSyntheticCall = "mark_synthetic_call"
	//SyntheticTestFieldConfigRetries constant value for the schema field configuration.retries
	SyntheticTestFieldConfigRetries = "retries"
	//SyntheticTestFieldConfigRetryInterval constant value for the schema field configuration.retry_interval
	SyntheticTestFieldConfigRetryInterval = "retry_interval"
	//SyntheticTestFieldConfigSyntheticType constant value for the schema field configuration.synthetic_type
	SyntheticTestFieldConfigSyntheticType = "synthetic_type"
	//SyntheticTestFieldConfigTimeout constant value for the schema field configuration.timeout
	SyntheticTestFieldConfigTimeout = "timeout"
	//SyntheticTestFieldConfigUrl constant value for the schema field configuration.url
	SyntheticTestFieldConfigUrl = "url"
	//SyntheticTestFieldConfigOperation constant value for the schema field configuration.operation
	SyntheticTestFieldConfigOperation = "operation"
	//SyntheticTestFieldConfigScript constant value for the schema field configuration.script
	SyntheticTestFieldConfigScript = "script"
)

// NewSyntheticTestResourceHandle creates the resource handle Synthetic Tests
func NewSyntheticTestResourceHandle() ResourceHandle {
	return &syntheticTestResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaSyntheticTest,
			Schema: map[string]*schema.Schema{
				SyntheticTestFieldLabel: {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "Friendly name of the Synthetic test",
					ValidateFunc: validation.StringLenBetween(0, 512),
				},
				SyntheticTestFieldDescription: {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "The description of the Synthetic test",
					ValidateFunc: validation.StringLenBetween(0, 512),
				},
				SyntheticTestFieldActive: {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     true,
					Description: "Indicates if the Synthetic test is started or not",
				},
				SyntheticTestFieldConfiguration: {
					Type:        schema.TypeList,
					MinItems:    1,
					MaxItems:    1,
					Required:    true,
					Description: "The configuration of the synthetic alert",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							SyntheticTestFieldConfigMarkSyntheticCall: {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     false,
								Description: "Flag used to control if HTTP calls will be marked as synthetic calls",
							},
							SyntheticTestFieldConfigRetries: {
								Type:         schema.TypeInt,
								Optional:     true,
								Default:      0,
								Description:  "Indicates how many attempts will be allowed to get a successful connection",
								ValidateFunc: validation.IntBetween(0, 2),
							},
							SyntheticTestFieldConfigRetryInterval: {
								Type:         schema.TypeInt,
								Optional:     true,
								Default:      1,
								Description:  "The time interval between retries in seconds",
								ValidateFunc: validation.IntBetween(1, 10),
							},
							SyntheticTestFieldConfigSyntheticType: {
								Type:         schema.TypeString,
								Required:     true,
								Description:  "The type of the Synthetic test",
								ValidateFunc: validation.StringInSlice([]string{"HTTPAction", "HTTPScript", "BrowserScript", "WebpageAction", "WebpageScript", "DNSAction"}, true),
							},
							SyntheticTestFieldConfigTimeout: {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "The timeout to be used by the PoP playback engines running the test",
							},
							SyntheticTestFieldConfigUrl: {
								Type:         schema.TypeString,
								Optional:     true,
								Description:  "The URL which is being tested",
								ValidateFunc: validation.IsURLWithHTTPorHTTPS,
							},
							SyntheticTestFieldConfigOperation: {
								Type:         schema.TypeString,
								Optional:     true,
								Default:      "GET",
								Description:  "The HTTP operation",
								ValidateFunc: validation.StringInSlice([]string{"GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT", "DELETE"}, true),
							},
							SyntheticTestFieldConfigScript: {
								Type:        schema.TypeString,
								Optional:    true,
								Description: " The Javascript content in plain text",
							},
						},
					},
				},
				SyntheticTestFieldCustomProperties: {
					Type:        schema.TypeMap,
					Optional:    true,
					Description: "Name/value pairs to provide additional information of the Synthetic test",
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				SyntheticTestFieldLocations: {
					Type:        schema.TypeSet,
					Required:    true,
					Description: "Array of the PoP location IDs",
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				SyntheticTestFieldPlaybackMode: {
					Type:         schema.TypeString,
					Optional:     true,
					Default:      "Simultaneous",
					Description:  "Defines how the Synthetic test should be executed across multiple PoPs",
					ValidateFunc: validation.StringInSlice([]string{"Simultaneous", "Staggered"}, true),
				},
				SyntheticTestFieldTestFrequency: {
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

type syntheticTestResource struct {
	metaData ResourceMetaData
}

func (r *syntheticTestResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *syntheticTestResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{}
}

func (r *syntheticTestResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource {
	return api.SyntheticTest()
}

func (r *syntheticTestResource) SetComputedFields(d *schema.ResourceData) {
	// No computed fields
}

func (r *syntheticTestResource) UpdateState(d *schema.ResourceData, obj restapi.InstanaDataObject, formatter utils.ResourceNameFormatter) error {
	syntheticTest := obj.(*restapi.SyntheticTest)
	d.SetId(syntheticTest.ID)
	d.Set(SyntheticTestFieldLabel, syntheticTest.Label)
	d.Set(SyntheticTestFieldActive, syntheticTest.Active)
	d.Set(SyntheticTestFieldCustomProperties, syntheticTest.CustomProperties)
	d.Set(SyntheticTestFieldLocations, syntheticTest.Locations)
	d.Set(SyntheticTestFieldPlaybackMode, syntheticTest.PlaybackMode)
	d.Set(SyntheticTestFieldTestFrequency, syntheticTest.TestFrequency)
	d.Set(SyntheticTestFieldConfiguration, r.mapConfigurationToSchema(syntheticTest))
	return nil
}

func (r *syntheticTestResource) MapStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	return &restapi.SyntheticTest{
		ID:               d.Id(),
		Label:            d.Get(SyntheticTestFieldLabel).(string),
		Description:      d.Get(SyntheticTestFieldDescription).(string),
		Active:           d.Get(SyntheticTestFieldActive).(bool),
		Configuration:    r.mapConfigurationFromSchema(d),
		CustomProperties: d.Get(SyntheticTestFieldCustomProperties).(map[string]interface{}),
		Locations:        ReadStringSetParameterFromResource(d, SyntheticTestFieldLocations),
		PlaybackMode:     d.Get(SyntheticTestFieldPlaybackMode).(string),
		TestFrequency:    int32(d.Get(SyntheticTestFieldTestFrequency).(int)),
	}, nil
}

func (r *syntheticTestResource) mapConfigurationToSchema(config *restapi.SyntheticTest) []map[string]interface{} {
	configuration := make(map[string]interface{})
	configuration[SyntheticTestFieldConfigMarkSyntheticCall] = config.Configuration.MarkSyntheticCall
	configuration[SyntheticTestFieldConfigSyntheticType] = config.Configuration.SyntheticType
	configuration[SyntheticTestFieldConfigTimeout] = config.Configuration.Timeout
	configuration[SyntheticTestFieldConfigRetries] = config.Configuration.Retries
	configuration[SyntheticTestFieldConfigRetryInterval] = config.Configuration.RetryInterval
	configuration[SyntheticTestFieldConfigUrl] = config.Configuration.URL
	configuration[SyntheticTestFieldConfigScript] = config.Configuration.Script
	configuration[SyntheticTestFieldConfigOperation] = config.Configuration.Operation
	result := make([]map[string]interface{}, 1)
	result[0] = configuration
	return result
}

func (r *syntheticTestResource) mapConfigurationFromSchema(d *schema.ResourceData) restapi.SyntheticTestConfig {
	syntheticTestConfigurationSlice := d.Get(SyntheticTestFieldConfiguration).([]interface{})
	syntheticTestConfig := syntheticTestConfigurationSlice[0].(map[string]interface{})

	return restapi.SyntheticTestConfig{
		MarkSyntheticCall: syntheticTestConfig[SyntheticTestFieldConfigMarkSyntheticCall].(bool),
		Retries:           int32(syntheticTestConfig[SyntheticTestFieldConfigRetries].(int)),
		RetryInterval:     int32(syntheticTestConfig[SyntheticTestFieldConfigRetryInterval].(int)),
		SyntheticType:     syntheticTestConfig[SyntheticTestFieldConfigSyntheticType].(string),
		Timeout:           syntheticTestConfig[SyntheticTestFieldConfigTimeout].(string),
		URL:               syntheticTestConfig[SyntheticTestFieldConfigUrl].(string),
		Operation:         syntheticTestConfig[SyntheticTestFieldConfigOperation].(string),
		Script:            syntheticTestConfig[SyntheticTestFieldConfigScript].(string),
	}
}
