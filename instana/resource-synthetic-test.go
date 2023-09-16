package instana

import (
	"errors"
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
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
	SyntheticTestFieldActive        = "active"
	SyntheticTestFieldApplicationID = "application_id"
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

	//SyntheticTestFieldConfigHttpScript constant value for the schema field configuration.http_script
	SyntheticTestFieldConfigHttpScript = "http_script"
	//SyntheticTestFieldConfigHttpAction constant value for the schema field configuration.http_action
	SyntheticTestFieldConfigHttpAction = "http_action"

	//SyntheticTestFieldConfigMarkSyntheticCall constant value for the schema field configuration.mark_synthetic_call
	SyntheticTestFieldConfigMarkSyntheticCall = "mark_synthetic_call"
	//SyntheticTestFieldConfigRetries constant value for the schema field configuration.retries
	SyntheticTestFieldConfigRetries = "retries"
	//SyntheticTestFieldConfigRetryInterval constant value for the schema field configuration.retry_interval
	SyntheticTestFieldConfigRetryInterval = "retry_interval"
	//SyntheticTestFieldConfigTimeout constant value for the schema field configuration.timeout
	SyntheticTestFieldConfigTimeout = "timeout"
	//SyntheticTestFieldConfigUrl constant value for the schema field configuration.url
	SyntheticTestFieldConfigUrl = "url"
	//SyntheticTestFieldConfigOperation constant value for the schema field configuration.operation
	SyntheticTestFieldConfigOperation = "operation"
	//SyntheticTestFieldConfigHeaders constant value for the schema field configuration.headers
	SyntheticTestFieldConfigHeaders = "headers"
	//SyntheticTestFieldConfigBody constant value for the schema field configuration.body
	SyntheticTestFieldConfigBody = "body"
	//SyntheticTestFieldConfigValidationString constant value for the schema field configuration.validation_string
	SyntheticTestFieldConfigValidationString = "validation_string"
	//SyntheticTestFieldConfigFollowRedirect constant value for the schema field configuration.follow_redirect
	SyntheticTestFieldConfigFollowRedirect = "follow_redirect"
	//SyntheticTestFieldConfigAllowInsecure constant value for the schema field configuration.allow_insecure
	SyntheticTestFieldConfigAllowInsecure = "allow_insecure"
	//SyntheticTestFieldConfigExpectStatus constant value for the schema field configuration.expect_status
	SyntheticTestFieldConfigExpectStatus = "expect_status"
	//SyntheticTestFieldConfigExpectMatch constant value for the schema field configuration.expect_match
	SyntheticTestFieldConfigExpectMatch = "expect_match"
	//SyntheticTestFieldConfigScript constant value for the schema field configuration.script
	SyntheticTestFieldConfigScript = "script"
)

var syntheticTestConfigurationOptions = []string{
	"configuration.0.http_script",
	"configuration.0.http_action",
}

const SyntheticCheckTypeHttpAction = "HTTPAction"
const SyntheticCheckTypeHttpScript = "HTTPScript"

var (
	syntheticTestSchemaConfigMarkSyntheticCall = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Flag used to control if HTTP calls will be marked as synthetic calls",
	}
	syntheticTestSchemaConfigRetries = &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      0,
		Description:  "Indicates how many attempts will be allowed to get a successful connection",
		ValidateFunc: validation.IntBetween(0, 2),
	}
	syntheticTestSchemaConfigRetryInterval = &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      1,
		Description:  "The time interval between retries in seconds",
		ValidateFunc: validation.IntBetween(1, 10),
	}
	syntheticTestSchemaConfigTimeout = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The timeout to be used by the PoP playback engines running the test",
	}
)

// NewSyntheticTestResourceHandle creates the resource handle Synthetic Tests
func NewSyntheticTestResourceHandle() ResourceHandle[*restapi.SyntheticTest] {
	return &syntheticTestResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaSyntheticTest,
			Schema: map[string]*schema.Schema{
				SyntheticTestFieldLabel: {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "Friendly name of the Synthetic test",
					ValidateFunc: validation.StringLenBetween(0, 128),
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
				SyntheticTestFieldApplicationID: {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Unique identifier of the Application Perspective.",
				},

				SyntheticTestFieldConfiguration: {
					Type:        schema.TypeList,
					MinItems:    1,
					MaxItems:    1,
					Required:    true,
					Description: "The configuration of the synthetic alert",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							SyntheticTestFieldConfigHttpAction: {
								Type:         schema.TypeList,
								MinItems:     0,
								MaxItems:     1,
								Optional:     true,
								Description:  "The configuration of the synthetic alert of type http script",
								ExactlyOneOf: syntheticTestConfigurationOptions,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										SyntheticTestFieldConfigMarkSyntheticCall: syntheticTestSchemaConfigMarkSyntheticCall,
										SyntheticTestFieldConfigRetries:           syntheticTestSchemaConfigRetries,
										SyntheticTestFieldConfigRetryInterval:     syntheticTestSchemaConfigRetryInterval,
										SyntheticTestFieldConfigTimeout:           syntheticTestSchemaConfigTimeout,
										SyntheticTestFieldConfigUrl: {
											Type:         schema.TypeString,
											Optional:     true,
											Description:  "The URL which is being tested",
											ValidateFunc: validation.IsURLWithHTTPorHTTPS,
										},
										SyntheticTestFieldConfigOperation: {
											Type:         schema.TypeString,
											Optional:     true,
											Description:  "The HTTP operation",
											ValidateFunc: validation.StringInSlice([]string{"GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT", "DELETE"}, true),
										},
										SyntheticTestFieldConfigHeaders: {
											Type:        schema.TypeMap,
											Optional:    true,
											Description: "An object with header/value pairs",
											Elem: &schema.Schema{
												Type: schema.TypeString,
											},
										},
										SyntheticTestFieldConfigBody: {
											Type:        schema.TypeString,
											Optional:    true,
											Description: " The body content to send with the operation",
										},
										SyntheticTestFieldConfigValidationString: {
											Type:        schema.TypeString,
											Optional:    true,
											Description: "An expression to be evaluated",
										},
										SyntheticTestFieldConfigFollowRedirect: {
											Type:        schema.TypeBool,
											Optional:    true,
											Default:     false,
											Description: "A boolean type, true by default; to allow redirect",
										},
										SyntheticTestFieldConfigAllowInsecure: {
											Type:        schema.TypeBool,
											Optional:    true,
											Default:     false,
											Description: "A boolean type, if set to true then allow insecure certificates",
										},
										SyntheticTestFieldConfigExpectStatus: {
											Type:        schema.TypeInt,
											Optional:    true,
											Description: "An integer type, by default, the Synthetic passes for any 2XX status code",
										},
										SyntheticTestFieldConfigExpectMatch: {
											Type:        schema.TypeString,
											Optional:    true,
											Description: "An optional regular expression string to be used to check the test response",
										},
									},
								},
							},
							SyntheticTestFieldConfigHttpScript: {
								Type:         schema.TypeList,
								MinItems:     0,
								MaxItems:     1,
								Optional:     true,
								Description:  "The configuration of the synthetic alert of type http action",
								ExactlyOneOf: syntheticTestConfigurationOptions,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										SyntheticTestFieldConfigMarkSyntheticCall: syntheticTestSchemaConfigMarkSyntheticCall,
										SyntheticTestFieldConfigRetries:           syntheticTestSchemaConfigRetries,
										SyntheticTestFieldConfigRetryInterval:     syntheticTestSchemaConfigRetryInterval,
										SyntheticTestFieldConfigTimeout:           syntheticTestSchemaConfigTimeout,
										SyntheticTestFieldConfigScript: {
											Type:        schema.TypeString,
											Required:    true,
											Description: "The Javascript content in plain text",
										},
									},
								},
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

func (r *syntheticTestResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.SyntheticTest] {
	return api.SyntheticTest()
}

func (r *syntheticTestResource) SetComputedFields(_ *schema.ResourceData) error {
	return nil
}

func (r *syntheticTestResource) UpdateState(d *schema.ResourceData, syntheticTest *restapi.SyntheticTest) error {
	config, err := r.mapConfigurationToSchema(syntheticTest)
	if err != nil {
		return err
	}
	d.SetId(syntheticTest.ID)
	return tfutils.UpdateState(d, map[string]interface{}{
		SyntheticTestFieldLabel:            syntheticTest.Label,
		SyntheticTestFieldActive:           syntheticTest.Active,
		SyntheticTestFieldDescription:      syntheticTest.Description,
		SyntheticTestFieldApplicationID:    syntheticTest.ApplicationID,
		SyntheticTestFieldCustomProperties: syntheticTest.CustomProperties,
		SyntheticTestFieldLocations:        syntheticTest.Locations,
		SyntheticTestFieldPlaybackMode:     syntheticTest.PlaybackMode,
		SyntheticTestFieldTestFrequency:    syntheticTest.TestFrequency,
		SyntheticTestFieldConfiguration:    []interface{}{config},
	})
}

func (r *syntheticTestResource) mapConfigurationToSchema(config *restapi.SyntheticTest) (map[string]interface{}, error) {
	configuration := make(map[string]interface{})
	configuration[SyntheticTestFieldConfigMarkSyntheticCall] = config.Configuration.MarkSyntheticCall
	configuration[SyntheticTestFieldConfigTimeout] = config.Configuration.Timeout
	configuration[SyntheticTestFieldConfigRetries] = config.Configuration.Retries
	configuration[SyntheticTestFieldConfigRetryInterval] = config.Configuration.RetryInterval

	switch config.Configuration.SyntheticType {
	case SyntheticCheckTypeHttpAction:
		configuration[SyntheticTestFieldConfigUrl] = config.Configuration.URL
		configuration[SyntheticTestFieldConfigOperation] = config.Configuration.Operation
		configuration[SyntheticTestFieldConfigHeaders] = config.Configuration.Headers
		configuration[SyntheticTestFieldConfigBody] = config.Configuration.Body
		configuration[SyntheticTestFieldConfigValidationString] = config.Configuration.ValidationString
		configuration[SyntheticTestFieldConfigFollowRedirect] = config.Configuration.FollowRedirect
		configuration[SyntheticTestFieldConfigAllowInsecure] = config.Configuration.AllowInsecure
		configuration[SyntheticTestFieldConfigExpectStatus] = config.Configuration.ExpectStatus
		configuration[SyntheticTestFieldConfigExpectMatch] = config.Configuration.ExpectMatch

		return map[string]interface{}{
			SyntheticTestFieldConfigHttpAction: []interface{}{
				configuration,
			},
		}, nil
	case SyntheticCheckTypeHttpScript:
		configuration[SyntheticTestFieldConfigScript] = config.Configuration.Script
		return map[string]interface{}{
			SyntheticTestFieldConfigHttpScript: []interface{}{
				configuration,
			},
		}, nil
	}

	return nil, fmt.Errorf("unsupported synthetic test of type %s received", config.Configuration.SyntheticType)
}

func (r *syntheticTestResource) MapStateToDataObject(d *schema.ResourceData) (*restapi.SyntheticTest, error) {
	appID, ok := d.GetOk(SyntheticTestFieldApplicationID)
	var applicationID *string
	if ok {
		tempAppID := appID.(string)
		applicationID = &tempAppID
	}
	configuration, err := r.mapConfigurationFromSchema(d)
	if err != nil {
		return nil, err
	}
	return &restapi.SyntheticTest{
		ID:               d.Id(),
		Label:            d.Get(SyntheticTestFieldLabel).(string),
		Description:      GetStringPointerFromResourceData(d, SyntheticTestFieldDescription),
		Active:           d.Get(SyntheticTestFieldActive).(bool),
		ApplicationID:    applicationID,
		Configuration:    configuration,
		CustomProperties: d.Get(SyntheticTestFieldCustomProperties).(map[string]interface{}),
		Locations:        ReadStringSetParameterFromResource(d, SyntheticTestFieldLocations),
		PlaybackMode:     d.Get(SyntheticTestFieldPlaybackMode).(string),
		TestFrequency:    GetInt32PointerFromResourceData(d, SyntheticTestFieldTestFrequency),
	}, nil
}

func (r *syntheticTestResource) mapConfigurationFromSchema(d *schema.ResourceData) (restapi.SyntheticTestConfig, error) {
	syntheticTestConfigurationSlice := d.Get(SyntheticTestFieldConfiguration).([]interface{})
	syntheticTestConfiguration := syntheticTestConfigurationSlice[0].(map[string]interface{})
	var syntheticTestType string
	var syntheticTestConfigData map[string]interface{}
	if val, ok := syntheticTestConfiguration[SyntheticTestFieldConfigHttpAction]; ok && len(val.([]interface{})) == 1 {
		syntheticTestType = SyntheticCheckTypeHttpAction
		syntheticTestConfigData = val.([]interface{})[0].(map[string]interface{})
	} else if val, ok := syntheticTestConfiguration[SyntheticTestFieldConfigHttpScript]; ok && len(val.([]interface{})) == 1 {
		syntheticTestType = SyntheticCheckTypeHttpScript
		syntheticTestConfigData = val.([]interface{})[0].(map[string]interface{})
	} else {
		return restapi.SyntheticTestConfig{}, errors.New("no supported synthetic test configuration provided")
	}

	expectedStatusAsInt := GetPointerFromMap[int](syntheticTestConfigData, SyntheticTestFieldConfigExpectStatus)
	var expectedStatus *int32
	if expectedStatusAsInt != nil {
		v := int32(*expectedStatusAsInt)
		expectedStatus = &v
	}
	headersRaw, ok := syntheticTestConfigData[SyntheticTestFieldConfigHeaders]
	var headers map[string]interface{}
	if ok {
		headers = headersRaw.(map[string]interface{})
	}
	return restapi.SyntheticTestConfig{
		MarkSyntheticCall: syntheticTestConfigData[SyntheticTestFieldConfigMarkSyntheticCall].(bool),
		Retries:           int32(syntheticTestConfigData[SyntheticTestFieldConfigRetries].(int)),
		RetryInterval:     int32(syntheticTestConfigData[SyntheticTestFieldConfigRetryInterval].(int)),
		SyntheticType:     syntheticTestType,
		Timeout:           GetPointerFromMap[string](syntheticTestConfigData, SyntheticTestFieldConfigTimeout),
		URL:               GetPointerFromMap[string](syntheticTestConfigData, SyntheticTestFieldConfigUrl),
		Operation:         GetPointerFromMap[string](syntheticTestConfigData, SyntheticTestFieldConfigOperation),
		Headers:           headers,
		Body:              GetPointerFromMap[string](syntheticTestConfigData, SyntheticTestFieldConfigBody),
		ValidationString:  GetPointerFromMap[string](syntheticTestConfigData, SyntheticTestFieldConfigValidationString),
		FollowRedirect:    GetPointerFromMap[bool](syntheticTestConfigData, SyntheticTestFieldConfigFollowRedirect),
		AllowInsecure:     GetPointerFromMap[bool](syntheticTestConfigData, SyntheticTestFieldConfigAllowInsecure),
		ExpectStatus:      expectedStatus,
		ExpectMatch:       GetPointerFromMap[string](syntheticTestConfigData, SyntheticTestFieldConfigExpectMatch),
		Script:            GetPointerFromMap[string](syntheticTestConfigData, SyntheticTestFieldConfigScript),
	}, nil
}
