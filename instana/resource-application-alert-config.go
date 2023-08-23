package instana

import (
	"context"
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceInstanaApplicationAlertConfig the name of the terraform-provider-instana resource to manage application alert configs
const ResourceInstanaApplicationAlertConfig = "instana_application_alert_config"

// ResourceInstanaGlobalApplicationAlertConfig the name of the terraform-provider-instana resource to manage global application alert configs
const ResourceInstanaGlobalApplicationAlertConfig = "instana_global_application_alert_config"

const (
	//ApplicationAlertConfigFieldAlertChannelIDs constant value for field alerting_channel_ids of resource instana_application_alert_config
	ApplicationAlertConfigFieldAlertChannelIDs = "alert_channel_ids"
	//ApplicationAlertConfigFieldApplications constant value for field applications of resource instana_application_alert_config
	ApplicationAlertConfigFieldApplications = "application"
	//ApplicationAlertConfigFieldApplicationsApplicationID constant value for field applications.application_id of resource instana_application_alert_config
	ApplicationAlertConfigFieldApplicationsApplicationID = "application_id"
	//ApplicationAlertConfigFieldApplicationsInclusive constant value for field applications.inclusive of resource instana_application_alert_config
	ApplicationAlertConfigFieldApplicationsInclusive            = "inclusive"
	applicationAlertConfigFieldApplicationsInclusiveDescription = "Defines whether this node and his child nodes are included (true) or excluded (false)"
	//ApplicationAlertConfigFieldApplicationsServices constant value for field applications.services of resource instana_application_alert_config
	ApplicationAlertConfigFieldApplicationsServices = "service"
	//ApplicationAlertConfigFieldApplicationsServicesServiceID constant value for field applications.services.service_id of resource instana_application_alert_config
	ApplicationAlertConfigFieldApplicationsServicesServiceID = "service_id"
	//ApplicationAlertConfigFieldApplicationsServicesEndpoints constant value for field applications.services.endpoints of resource instana_application_alert_config
	ApplicationAlertConfigFieldApplicationsServicesEndpoints = "endpoint"
	//ApplicationAlertConfigFieldApplicationsServicesEndpointsEndpointID constant value for field applications.services.endpoints.endpoint_id of resource instana_application_alert_config
	ApplicationAlertConfigFieldApplicationsServicesEndpointsEndpointID = "endpoint_id"
	//ApplicationAlertConfigFieldBoundaryScope constant value for field boundary_scope of resource instana_application_alert_config
	ApplicationAlertConfigFieldBoundaryScope = "boundary_scope"
	//ApplicationAlertConfigFieldCustomPayloadFields constant value for field custom_payload_fields of resource instana_application_alert_config
	ApplicationAlertConfigFieldCustomPayloadFields = "custom_payload_field"
	//ApplicationAlertConfigFieldCustomPayloadFieldsKey constant value for field custom_payload_fields.key of resource instana_application_alert_config
	ApplicationAlertConfigFieldCustomPayloadFieldsKey = "key"
	//ApplicationAlertConfigFieldCustomPayloadFieldsValue constant value for field custom_payload_fields.value of resource instana_application_alert_config
	ApplicationAlertConfigFieldCustomPayloadFieldsValue = "value"
	//ApplicationAlertConfigFieldDescription constant value for field description of resource instana_application_alert_config
	ApplicationAlertConfigFieldDescription = "description"
	//ApplicationAlertConfigFieldEvaluationType constant value for field evaluation_type of resource instana_application_alert_config
	ApplicationAlertConfigFieldEvaluationType = "evaluation_type"
	//ApplicationAlertConfigFieldGranularity constant value for field granularity of resource instana_application_alert_config
	ApplicationAlertConfigFieldGranularity = "granularity"
	//ApplicationAlertConfigFieldIncludeInternal constant value for field include_internal of resource instana_application_alert_config
	ApplicationAlertConfigFieldIncludeInternal = "include_internal"
	//ApplicationAlertConfigFieldIncludeSynthetic constant value for field include_synthetic of resource instana_application_alert_config
	ApplicationAlertConfigFieldIncludeSynthetic = "include_synthetic"
	//ApplicationAlertConfigFieldName constant value for field name of resource instana_application_alert_config
	ApplicationAlertConfigFieldName = "name"
	//ApplicationAlertConfigFieldFullName constant value for field full_name of resource instana_application_alert_config
	ApplicationAlertConfigFieldFullName = "full_name"
	//ApplicationAlertConfigFieldRule constant value for field rule of resource instana_application_alert_config
	ApplicationAlertConfigFieldRule = "rule"
	//ApplicationAlertConfigFieldRuleMetricName constant value for field rule.*.metric_name of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleMetricName = "metric_name"
	//ApplicationAlertConfigFieldRuleAggregation constant value for field rule.*.aggregation of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleAggregation = "aggregation"
	//ApplicationAlertConfigFieldRuleStableHash constant value for field rule.*.stable_hash of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleStableHash = "stable_hash"
	//ApplicationAlertConfigFieldRuleErrorRate constant value for field rule.error_rate of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleErrorRate = "error_rate"
	//ApplicationAlertConfigFieldRuleLogs constant value for field rule.logs of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleLogs = "logs"
	//ApplicationAlertConfigFieldRuleLogsLevel constant value for field rule.logs.level of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleLogsLevel = "level"
	//ApplicationAlertConfigFieldRuleLogsMessage constant value for field rule.logs.message of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleLogsMessage = "message"
	//ApplicationAlertConfigFieldRuleLogsOperator constant value for field rule.logs.operator of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleLogsOperator = "operator"
	//ApplicationAlertConfigFieldRuleSlowness constant value for field rule.slowness of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleSlowness = "slowness"
	//ApplicationAlertConfigFieldRuleStatusCode constant value for field rule.status_code of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleStatusCode = "status_code"
	//ApplicationAlertConfigFieldRuleStatusCodeStart constant value for field rule.status_code.status_code_start of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleStatusCodeStart = "status_code_start"
	//ApplicationAlertConfigFieldRuleStatusCodeEnd constant value for field rule.status_code.status_code_end of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleStatusCodeEnd = "status_code_end"
	//ApplicationAlertConfigFieldRuleThroughput constant value for field rule.throughput of resource instana_application_alert_config
	ApplicationAlertConfigFieldRuleThroughput = "throughput"
	//ApplicationAlertConfigFieldSeverity constant value for field severity of resource instana_application_alert_config
	ApplicationAlertConfigFieldSeverity = "severity"
	//ApplicationAlertConfigFieldTagFilter constant value for field tag_filter of resource instana_application_alert_config
	ApplicationAlertConfigFieldTagFilter = "tag_filter"
	//ApplicationAlertConfigFieldTimeThreshold constant value for field time_threshold of resource instana_application_alert_config
	ApplicationAlertConfigFieldTimeThreshold = "time_threshold"
	//ApplicationAlertConfigFieldTimeThresholdTimeWindow constant value for field time_threshold.time_window of resource instana_application_alert_config
	ApplicationAlertConfigFieldTimeThresholdTimeWindow = "time_window"
	//ApplicationAlertConfigFieldTimeThresholdRequestImpact constant value for field time_threshold.request_impact of resource instana_application_alert_config
	ApplicationAlertConfigFieldTimeThresholdRequestImpact = "request_impact"
	//ApplicationAlertConfigFieldTimeThresholdRequestImpactRequests constant value for field time_threshold.request_impact.requests of resource instana_application_alert_config
	ApplicationAlertConfigFieldTimeThresholdRequestImpactRequests = "requests"
	//ApplicationAlertConfigFieldTimeThresholdViolationsInPeriod constant value for field time_threshold.violations_in_period of resource instana_application_alert_config
	ApplicationAlertConfigFieldTimeThresholdViolationsInPeriod = "violations_in_period"
	//ApplicationAlertConfigFieldTimeThresholdViolationsInPeriodViolations constant value for field time_threshold.violations_in_period.violations of resource instana_application_alert_config
	ApplicationAlertConfigFieldTimeThresholdViolationsInPeriodViolations = "violations"
	//ApplicationAlertConfigFieldTimeThresholdViolationsInSequence constant value for field time_threshold.violations_in_sequence of resource instana_application_alert_config
	ApplicationAlertConfigFieldTimeThresholdViolationsInSequence = "violations_in_sequence"
	//ApplicationAlertConfigFieldTriggering constant value for field triggering of resource instana_application_alert_config
	ApplicationAlertConfigFieldTriggering = "triggering"
)

var (
	applicationAlertRuleTypeKeys = []string{
		"rule.0.error_rate",
		"rule.0.logs",
		"rule.0.slowness",
		"rule.0.status_code",
		"rule.0.throughput",
	}

	applicationAlertSchemaRuleMetricName = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The metric name of the application alert rule",
	}

	applicationAlertSchemaRequiredRuleAggregation = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice(restapi.SupportedAggregations.ToStringSlice(), true),
		Description:  "The aggregation function of the application alert rule",
	}

	applicationAlertSchemaRequiredRuleStableHash = &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "The stable hash used for the application alert rule",
	}

	applicationAlertSchemaOptionalRuleAggregation = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringInSlice(restapi.SupportedAggregations.ToStringSlice(), true),
		Description:  "The aggregation function of the application alert rule",
	}

	applicationAlertSchemaRequiredRuleOperator = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		Description:  "The operator which will be applied to evaluate this rule",
		ValidateFunc: validation.StringInSlice(restapi.SupportedExpressionOperators.ToStringSlice(), true),
	}

	applicationAlertTimeThresholdTypeKeys = []string{
		"time_threshold.0.request_impact",
		"time_threshold.0.violations_in_period",
		"time_threshold.0.violations_in_sequence",
	}

	applicationAlertSchemaRequiredTimeThresholdTimeWindow = &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "The time window if the time threshold",
	}
)

func applicationAlertConfigApplicationSchemaSetFunc(i interface{}) int {
	return schema.HashString(i.(map[string]interface{})[ApplicationAlertConfigFieldApplicationsApplicationID])
}
func applicationAlertConfigApplicationServiceSchemaSetFunc(i interface{}) int {
	return schema.HashString(i.(map[string]interface{})[ApplicationAlertConfigFieldApplicationsServicesServiceID])
}
func applicationAlertConfigApplicationServiceEndpointSchemaSetFunc(i interface{}) int {
	return schema.HashString(i.(map[string]interface{})[ApplicationAlertConfigFieldApplicationsServicesEndpointsEndpointID])
}

var (
	applicationAlertConfigSchemaAlertChannelIDs = &schema.Schema{
		Type:     schema.TypeSet,
		MinItems: 0,
		MaxItems: 1024,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional:    true,
		Description: "List of IDs of alert channels defined in Instana.",
	}
	applicationAlertConfigSchemaApplications = &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		Set:      applicationAlertConfigApplicationSchemaSetFunc,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				ApplicationAlertConfigFieldApplicationsApplicationID: {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "ID of the included application",
					ValidateFunc: validation.StringLenBetween(0, 64),
				},
				ApplicationAlertConfigFieldApplicationsInclusive: {
					Type:        schema.TypeBool,
					Required:    true,
					Description: applicationAlertConfigFieldApplicationsInclusiveDescription,
				},
				ApplicationAlertConfigFieldApplicationsServices: {
					Type:     schema.TypeSet,
					Optional: true,
					Set:      applicationAlertConfigApplicationServiceSchemaSetFunc,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							ApplicationAlertConfigFieldApplicationsServicesServiceID: {
								Type:         schema.TypeString,
								Required:     true,
								Description:  "ID of the included service",
								ValidateFunc: validation.StringLenBetween(0, 64),
							},
							ApplicationAlertConfigFieldApplicationsInclusive: {
								Type:        schema.TypeBool,
								Required:    true,
								Description: applicationAlertConfigFieldApplicationsInclusiveDescription,
							},
							ApplicationAlertConfigFieldApplicationsServicesEndpoints: {
								Type:     schema.TypeSet,
								Optional: true,
								Set:      applicationAlertConfigApplicationServiceEndpointSchemaSetFunc,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										ApplicationAlertConfigFieldApplicationsServicesEndpointsEndpointID: {
											Type:         schema.TypeString,
											Required:     true,
											Description:  "ID of the included endpoint",
											ValidateFunc: validation.StringLenBetween(0, 64),
										},
										ApplicationAlertConfigFieldApplicationsInclusive: {
											Type:        schema.TypeBool,
											Required:    true,
											Description: applicationAlertConfigFieldApplicationsInclusiveDescription,
										},
									},
								},
								Description: "Selection of endpoints in scope.",
							},
						},
					},
					Description: "Selection of services in scope.",
				},
			},
		},
		Description: "Selection of applications in scope.",
	}
	applicationAlertConfigSchemaBoundaryScope = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice(restapi.SupportedApplicationAlertConfigBoundaryScopes.ToStringSlice(), false),
		Description:  "The boundary scope of the application alert config",
	}
	applicationAlertConfigSchemaCustomPayloadFields = &schema.Schema{
		Type: schema.TypeSet,
		Set: func(i interface{}) int {
			return schema.HashString(i.(map[string]interface{})[ApplicationAlertConfigFieldCustomPayloadFieldsKey])
		},
		Optional: true,
		MinItems: 0,
		MaxItems: 20,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				/*
					ApplicationAlertConfigFieldCustomPayloadFieldsType: {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(restapi.SupportedCustomPayloadTypes.ToStringSlice(), false),
						Description:  "The type of the custom payload field",
					},
				*/
				ApplicationAlertConfigFieldCustomPayloadFieldsKey: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The key of the custom payload field",
				},
				ApplicationAlertConfigFieldCustomPayloadFieldsValue: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The value of the custom payload field",
				},
			},
		},
		Description: "An optional list of custom payload fields (static key/value pairs added to the event)",
	}
	applicationAlertConfigSchemaDescription = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		Description:  "The description text of the application alert config",
		ValidateFunc: validation.StringLenBetween(0, 65536),
	}
	applicationAlertConfigSchemaEvaluationType = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice(restapi.SupportedApplicationAlertEvaluationTypes.ToStringSlice(), false),
		Description:  "The evaluation type of the application alert config",
	}
	applicationAlertConfigSchemaGranularity = &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      restapi.Granularity600000,
		ValidateFunc: validation.IntInSlice(restapi.SupportedGranularities.ToIntSlice()),
		Description:  "The evaluation granularity used for detection of violations of the defined threshold. In other words, it defines the size of the tumbling window used",
	}
	applicationAlertConfigSchemaIncludeInternal = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Optional flag to indicate whether also internal calls are included in the scope or not. The default is false",
	}
	applicationAlertConfigSchemaIncludeSynthetic = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Optional flag to indicate whether also synthetic calls are included in the scope or not. The default is false",
	}
	applicationAlertConfigSchemaName = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		Description:  "Name for the application alert configuration",
		ValidateFunc: validation.StringLenBetween(0, 256),
	}
	applicationAlertConfigSchemaFullName = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The full name field of the application alert config. The field is computed and contains the name which is sent to Instana. The computation depends on the configured default_name_prefix and default_name_suffix at provider level",
	}
	applicationAlertConfigSchemaRule = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    1,
		MaxItems:    1,
		Required:    true,
		Description: "Indicates the type of rule this alert configuration is about.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				ApplicationAlertConfigFieldRuleErrorRate: {
					Type:        schema.TypeList,
					MinItems:    0,
					MaxItems:    1,
					Optional:    true,
					Description: "Rule based on the error rate of the configured alert configuration target",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							ApplicationAlertConfigFieldRuleMetricName:  applicationAlertSchemaRuleMetricName,
							ApplicationAlertConfigFieldRuleAggregation: applicationAlertSchemaOptionalRuleAggregation,
							ApplicationAlertConfigFieldRuleStableHash:  applicationAlertSchemaRequiredRuleStableHash,
						},
					},
					ExactlyOneOf: applicationAlertRuleTypeKeys,
				},
				ApplicationAlertConfigFieldRuleLogs: {
					Type:        schema.TypeList,
					MinItems:    0,
					MaxItems:    1,
					Optional:    true,
					Description: "Rule based on logs of the configured alert configuration target",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							ApplicationAlertConfigFieldRuleMetricName:  applicationAlertSchemaRuleMetricName,
							ApplicationAlertConfigFieldRuleAggregation: applicationAlertSchemaOptionalRuleAggregation,
							ApplicationAlertConfigFieldRuleStableHash:  applicationAlertSchemaRequiredRuleStableHash,
							ApplicationAlertConfigFieldRuleLogsLevel: {
								Type:         schema.TypeString,
								Required:     true,
								Description:  "The log level for which this rule applies to",
								ValidateFunc: validation.StringInSlice(restapi.SupportedLogLevels.ToStringSlice(), false),
							},
							ApplicationAlertConfigFieldRuleLogsMessage: {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "The log message for which this rule applies to",
							},
							ApplicationAlertConfigFieldRuleLogsOperator: applicationAlertSchemaRequiredRuleOperator,
						},
					},
					ExactlyOneOf: applicationAlertRuleTypeKeys,
				},
				ApplicationAlertConfigFieldRuleSlowness: {
					Type:        schema.TypeList,
					MinItems:    0,
					MaxItems:    1,
					Optional:    true,
					Description: "Rule based on the slowness of the configured alert configuration target",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							ApplicationAlertConfigFieldRuleMetricName:  applicationAlertSchemaRuleMetricName,
							ApplicationAlertConfigFieldRuleAggregation: applicationAlertSchemaRequiredRuleAggregation,
							ApplicationAlertConfigFieldRuleStableHash:  applicationAlertSchemaRequiredRuleStableHash,
						},
					},
					ExactlyOneOf: applicationAlertRuleTypeKeys,
				},
				ApplicationAlertConfigFieldRuleStatusCode: {
					Type:        schema.TypeList,
					MinItems:    0,
					MaxItems:    1,
					Optional:    true,
					Description: "Rule based on the HTTP status code of the configured alert configuration target",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							ApplicationAlertConfigFieldRuleMetricName:  applicationAlertSchemaRuleMetricName,
							ApplicationAlertConfigFieldRuleAggregation: applicationAlertSchemaOptionalRuleAggregation,
							ApplicationAlertConfigFieldRuleStableHash:  applicationAlertSchemaRequiredRuleStableHash,
							ApplicationAlertConfigFieldRuleStatusCodeStart: {
								Type:         schema.TypeInt,
								Optional:     true,
								Description:  "minimal HTTP status code applied for this rule",
								ValidateFunc: validation.IntAtLeast(1),
							},
							ApplicationAlertConfigFieldRuleStatusCodeEnd: {
								Type:         schema.TypeInt,
								Optional:     true,
								Description:  "maximum HTTP status code applied for this rule",
								ValidateFunc: validation.IntAtLeast(1),
							},
						},
					},
					ExactlyOneOf: applicationAlertRuleTypeKeys,
				},
				ApplicationAlertConfigFieldRuleThroughput: {
					Type:        schema.TypeList,
					MinItems:    0,
					MaxItems:    1,
					Optional:    true,
					Description: "Rule based on the throughput of the configured alert configuration target",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							ApplicationAlertConfigFieldRuleMetricName:  applicationAlertSchemaRuleMetricName,
							ApplicationAlertConfigFieldRuleAggregation: applicationAlertSchemaOptionalRuleAggregation,
							ApplicationAlertConfigFieldRuleStableHash:  applicationAlertSchemaRequiredRuleStableHash,
						},
					},
					ExactlyOneOf: applicationAlertRuleTypeKeys,
				},
			},
		},
	}
	applicationAlertConfigSchemaSeverity = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice(restapi.SupportedSeverities.TerraformRepresentations(), false),
		Description:  "The severity of the alert when triggered",
	}
	applicationAlertConfigSchemaTagFilter = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The tag filter of the application alert config",
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			normalized, err := tagfilter.Normalize(new)
			if err == nil {
				return normalized == old
			}
			return old == new
		},
		StateFunc: func(val interface{}) string {
			normalized, err := tagfilter.Normalize(val.(string))
			if err == nil {
				return normalized
			}
			return val.(string)
		},
		ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
			v := val.(string)
			if _, err := tagfilter.NewParser().Parse(v); err != nil {
				errs = append(errs, fmt.Errorf("%q is not a valid tag filter; %s", key, err))
			}

			return
		},
	}
	applicationAlertConfigSchemaTimeThreshold = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    1,
		MaxItems:    1,
		Required:    true,
		Description: "Indicates the type of violation of the defined threshold.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				ApplicationAlertConfigFieldTimeThresholdRequestImpact: {
					Type:        schema.TypeList,
					MinItems:    0,
					MaxItems:    1,
					Optional:    true,
					Description: "Time threshold base on request impact",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							ApplicationAlertConfigFieldTimeThresholdTimeWindow: applicationAlertSchemaRequiredTimeThresholdTimeWindow,
							ApplicationAlertConfigFieldTimeThresholdRequestImpactRequests: {
								Type:         schema.TypeInt,
								Optional:     true,
								ValidateFunc: validation.IntAtLeast(1),
								Description:  "The number of requests in the given window",
							},
						},
					},
					ExactlyOneOf: applicationAlertTimeThresholdTypeKeys,
				},
				ApplicationAlertConfigFieldTimeThresholdViolationsInPeriod: {
					Type:        schema.TypeList,
					MinItems:    0,
					MaxItems:    1,
					Optional:    true,
					Description: "Time threshold base on violations in period",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							ApplicationAlertConfigFieldTimeThresholdTimeWindow: applicationAlertSchemaRequiredTimeThresholdTimeWindow,
							ApplicationAlertConfigFieldTimeThresholdViolationsInPeriodViolations: {
								Type:         schema.TypeInt,
								Optional:     true,
								ValidateFunc: validation.IntBetween(1, 12),
								Description:  "The violations appeared in the period",
							},
						},
					},
					ExactlyOneOf: applicationAlertTimeThresholdTypeKeys,
				},
				ApplicationAlertConfigFieldTimeThresholdViolationsInSequence: {
					Type:        schema.TypeList,
					MinItems:    0,
					MaxItems:    1,
					Optional:    true,
					Description: "Time threshold base on violations in sequence",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							ApplicationAlertConfigFieldTimeThresholdTimeWindow: applicationAlertSchemaRequiredTimeThresholdTimeWindow,
						},
					},
					ExactlyOneOf: applicationAlertTimeThresholdTypeKeys,
				},
			},
		},
	}
	applicationAlertConfigSchemaTriggering = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Optional flag to indicate whether also an Incident is triggered or not. The default is false",
	}
)

var applicationAlertConfigResourceSchema = map[string]*schema.Schema{
	ApplicationAlertConfigFieldAlertChannelIDs:     applicationAlertConfigSchemaAlertChannelIDs,
	ApplicationAlertConfigFieldApplications:        applicationAlertConfigSchemaApplications,
	ApplicationAlertConfigFieldBoundaryScope:       applicationAlertConfigSchemaBoundaryScope,
	ApplicationAlertConfigFieldCustomPayloadFields: applicationAlertConfigSchemaCustomPayloadFields,
	ApplicationAlertConfigFieldDescription:         applicationAlertConfigSchemaDescription,
	ApplicationAlertConfigFieldEvaluationType:      applicationAlertConfigSchemaEvaluationType,
	ApplicationAlertConfigFieldGranularity:         applicationAlertConfigSchemaGranularity,
	ApplicationAlertConfigFieldIncludeInternal:     applicationAlertConfigSchemaIncludeInternal,
	ApplicationAlertConfigFieldIncludeSynthetic:    applicationAlertConfigSchemaIncludeSynthetic,
	ApplicationAlertConfigFieldName:                applicationAlertConfigSchemaName,
	ApplicationAlertConfigFieldRule:                applicationAlertConfigSchemaRule,
	ApplicationAlertConfigFieldSeverity:            applicationAlertConfigSchemaSeverity,
	ApplicationAlertConfigFieldTagFilter:           applicationAlertConfigSchemaTagFilter,
	ResourceFieldThreshold:                         thresholdSchema,
	ApplicationAlertConfigFieldTimeThreshold:       applicationAlertConfigSchemaTimeThreshold,
	ApplicationAlertConfigFieldTriggering:          applicationAlertConfigSchemaTriggering,
}

// NewApplicationAlertConfigResourceHandle creates a new instance of the ResourceHandle for application alert configs
func NewApplicationAlertConfigResourceHandle() ResourceHandle[*restapi.ApplicationAlertConfig] {
	return &applicationAlertConfigResource{
		metaData: ResourceMetaData{
			ResourceName:     ResourceInstanaApplicationAlertConfig,
			Schema:           applicationAlertConfigResourceSchema,
			SkipIDGeneration: true,
			SchemaVersion:    1,
		},
		resourceProvider: func(api restapi.InstanaAPI) restapi.RestResource[*restapi.ApplicationAlertConfig] {
			return api.ApplicationAlertConfigs()
		},
	}
}

// NewGlobalApplicationAlertConfigResourceHandle creates a new instance of the ResourceHandle for global application alert configs
func NewGlobalApplicationAlertConfigResourceHandle() ResourceHandle[*restapi.ApplicationAlertConfig] {
	return &applicationAlertConfigResource{
		metaData: ResourceMetaData{
			ResourceName:  ResourceInstanaGlobalApplicationAlertConfig,
			Schema:        applicationAlertConfigResourceSchema,
			SchemaVersion: 1,
		},
		resourceProvider: func(api restapi.InstanaAPI) restapi.RestResource[*restapi.ApplicationAlertConfig] {
			return api.GlobalApplicationAlertConfigs()
		},
	}
}

type applicationAlertConfigResource struct {
	metaData         ResourceMetaData
	resourceProvider func(api restapi.InstanaAPI) restapi.RestResource[*restapi.ApplicationAlertConfig]
}

func (r *applicationAlertConfigResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *applicationAlertConfigResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{
		{
			Type:    r.schemaV0().CoreConfigSchema().ImpliedType(),
			Upgrade: r.stateUpgradeV0,
			Version: 0,
		},
	}
}

func (r *applicationAlertConfigResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.ApplicationAlertConfig] {
	return r.resourceProvider(api)
}

func (r *applicationAlertConfigResource) SetComputedFields(_ *schema.ResourceData) error {
	return nil
}

func (r *applicationAlertConfigResource) UpdateState(d *schema.ResourceData, config *restapi.ApplicationAlertConfig) error {
	severity, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(config.Severity)
	if err != nil {
		return err
	}
	var normalizedTagFilterString *string
	if config.TagFilterExpression != nil {
		normalizedTagFilterString, err = tagfilter.MapTagFilterToNormalizedString(config.TagFilterExpression.(restapi.TagFilterExpressionElement))
		if err != nil {
			return err
		}
	}

	d.SetId(config.ID)
	return tfutils.UpdateState(d, map[string]interface{}{
		ApplicationAlertConfigFieldAlertChannelIDs:     config.AlertChannelIDs,
		ApplicationAlertConfigFieldApplications:        r.mapApplicationsToSchema(config),
		ApplicationAlertConfigFieldBoundaryScope:       config.BoundaryScope,
		ApplicationAlertConfigFieldCustomPayloadFields: r.mapCustomPayloadFieldsToSchema(config),
		ApplicationAlertConfigFieldDescription:         config.Description,
		ApplicationAlertConfigFieldEvaluationType:      config.EvaluationType,
		ApplicationAlertConfigFieldGranularity:         config.Granularity,
		ApplicationAlertConfigFieldIncludeInternal:     config.IncludeInternal,
		ApplicationAlertConfigFieldIncludeSynthetic:    config.IncludeSynthetic,
		ApplicationAlertConfigFieldName:                config.Name,
		ApplicationAlertConfigFieldRule:                r.mapRuleToSchema(config),
		ApplicationAlertConfigFieldSeverity:            severity,
		ApplicationAlertConfigFieldTagFilter:           normalizedTagFilterString,
		ResourceFieldThreshold:                         newThresholdMapper().toState(&config.Threshold),
		ApplicationAlertConfigFieldTimeThreshold:       r.mapTimeThresholdToSchema(config),
		ApplicationAlertConfigFieldTriggering:          config.Triggering,
	})
}

func (r *applicationAlertConfigResource) mapApplicationsToSchema(config *restapi.ApplicationAlertConfig) []interface{} {
	result := make([]interface{}, len(config.Applications))
	i := 0
	for _, v := range config.Applications {
		result[i] = r.mapApplicationToSchema(&v)
		i++
	}
	return result
}

func (r *applicationAlertConfigResource) mapApplicationToSchema(app *restapi.IncludedApplication) map[string]interface{} {
	result := make(map[string]interface{})
	result[ApplicationAlertConfigFieldApplicationsApplicationID] = app.ApplicationID
	result[ApplicationAlertConfigFieldApplicationsInclusive] = app.Inclusive

	services := make([]interface{}, len(app.Services))
	i := 0
	for _, v := range app.Services {
		services[i] = r.mapServiceToSchema(&v)
		i++
	}
	result[ApplicationAlertConfigFieldApplicationsServices] = services
	return result
}

func (r *applicationAlertConfigResource) mapServiceToSchema(service *restapi.IncludedService) map[string]interface{} {
	result := make(map[string]interface{})
	result[ApplicationAlertConfigFieldApplicationsServicesServiceID] = service.ServiceID
	result[ApplicationAlertConfigFieldApplicationsInclusive] = service.Inclusive

	endpoints := make([]interface{}, len(service.Endpoints))
	i := 0
	for _, v := range service.Endpoints {
		endpoints[i] = r.mapEndpointToSchema(&v)
		i++
	}
	result[ApplicationAlertConfigFieldApplicationsServicesEndpoints] = endpoints
	return result
}

func (r *applicationAlertConfigResource) mapEndpointToSchema(endpoint *restapi.IncludedEndpoint) map[string]interface{} {
	result := make(map[string]interface{})
	result[ApplicationAlertConfigFieldApplicationsServicesEndpointsEndpointID] = endpoint.EndpointID
	result[ApplicationAlertConfigFieldApplicationsInclusive] = endpoint.Inclusive
	return result
}

func (r *applicationAlertConfigResource) mapCustomPayloadFieldsToSchema(config *restapi.ApplicationAlertConfig) []map[string]string {
	result := make([]map[string]string, len(config.CustomerPayloadFields))
	for i, v := range config.CustomerPayloadFields {
		field := make(map[string]string)
		//field[ApplicationAlertConfigFieldCustomPayloadFieldsType] = string(v.Type)
		field[ApplicationAlertConfigFieldCustomPayloadFieldsKey] = v.Key
		if v.Type == restapi.DynamicCustomPayloadType {
			value := v.Value.(restapi.DynamicCustomPayloadFieldValue)
			key := ""
			if value.Key != nil {
				key = "=" + *value.Key
			}
			field[ApplicationAlertConfigFieldCustomPayloadFieldsValue] = value.TagName + key
		} else {
			field[ApplicationAlertConfigFieldCustomPayloadFieldsValue] = string(v.Value.(restapi.StaticStringCustomPayloadFieldValue))
		}
		result[i] = field
	}
	return result
}

func (r *applicationAlertConfigResource) mapRuleToSchema(config *restapi.ApplicationAlertConfig) []map[string]interface{} {
	ruleAttribute := make(map[string]interface{})
	ruleAttribute[ApplicationAlertConfigFieldRuleMetricName] = config.Rule.MetricName
	ruleAttribute[ApplicationAlertConfigFieldRuleAggregation] = config.Rule.Aggregation

	if config.Rule.StableHash != nil {
		ruleAttribute[ApplicationAlertConfigFieldRuleStableHash] = int(*config.Rule.StableHash)
	}

	if config.Rule.StatusCodeStart != nil {
		ruleAttribute[ApplicationAlertConfigFieldRuleStatusCodeStart] = int(*config.Rule.StatusCodeStart)
	}
	if config.Rule.StatusCodeEnd != nil {
		ruleAttribute[ApplicationAlertConfigFieldRuleStatusCodeEnd] = int(*config.Rule.StatusCodeEnd)
	}
	if config.Rule.Level != nil {
		ruleAttribute[ApplicationAlertConfigFieldRuleLogsLevel] = *config.Rule.Level
	}
	if config.Rule.Message != nil {
		ruleAttribute[ApplicationAlertConfigFieldRuleLogsMessage] = *config.Rule.Message
	}
	if config.Rule.Operator != nil {
		ruleAttribute[ApplicationAlertConfigFieldRuleLogsOperator] = *config.Rule.Operator
	}

	alertType := r.mapAlertTypeToSchema(config.Rule.AlertType)
	rule := make(map[string]interface{})
	rule[alertType] = []interface{}{ruleAttribute}
	result := make([]map[string]interface{}, 1)
	result[0] = rule
	return result
}

func (r *applicationAlertConfigResource) mapAlertTypeToSchema(alertType string) string {
	if alertType == "errorRate" {
		return ApplicationAlertConfigFieldRuleErrorRate
	} else if alertType == "statusCode" {
		return ApplicationAlertConfigFieldRuleStatusCode
	}
	return alertType
}

func (r *applicationAlertConfigResource) mapTimeThresholdToSchema(config *restapi.ApplicationAlertConfig) []map[string]interface{} {
	timeThresholdConfig := make(map[string]interface{})
	timeThresholdConfig[ApplicationAlertConfigFieldTimeThresholdTimeWindow] = config.TimeThreshold.TimeWindow

	if config.TimeThreshold.Violations != nil {
		timeThresholdConfig[ApplicationAlertConfigFieldTimeThresholdViolationsInPeriodViolations] = int(*config.TimeThreshold.Violations)
	}
	if config.TimeThreshold.Requests != nil {
		timeThresholdConfig[ApplicationAlertConfigFieldTimeThresholdRequestImpactRequests] = int(*config.TimeThreshold.Requests)
	}

	timeThresholdType := r.mapTimeThresholdTypeToSchema(config.TimeThreshold.Type)
	timeThreshold := make(map[string]interface{})
	timeThreshold[timeThresholdType] = []interface{}{timeThresholdConfig}
	result := make([]map[string]interface{}, 1)
	result[0] = timeThreshold
	return result
}

func (r *applicationAlertConfigResource) mapTimeThresholdTypeToSchema(input string) string {
	if input == "requestImpact" {
		return ApplicationAlertConfigFieldTimeThresholdRequestImpact
	} else if input == "violationsInPeriod" {
		return ApplicationAlertConfigFieldTimeThresholdViolationsInPeriod
	} else if input == "violationsInSequence" {
		return ApplicationAlertConfigFieldTimeThresholdViolationsInSequence
	}
	return input
}

func (r *applicationAlertConfigResource) MapStateToDataObject(d *schema.ResourceData) (*restapi.ApplicationAlertConfig, error) {
	severity, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(d.Get(ApplicationAlertConfigFieldSeverity).(string))
	if err != nil {
		return nil, err
	}

	var tagFilter restapi.TagFilterExpressionElement
	tagFilterStr, ok := d.GetOk(ApplicationAlertConfigFieldTagFilter)
	if ok {
		tagFilter, err = r.mapTagFilterExpressionFromSchema(tagFilterStr.(string))
		if err != nil {
			return &restapi.ApplicationAlertConfig{}, err
		}
	}

	threshold := newThresholdMapper().fromState(d)

	return &restapi.ApplicationAlertConfig{
		ID:                    d.Id(),
		AlertChannelIDs:       ReadStringSetParameterFromResource(d, ApplicationAlertConfigFieldAlertChannelIDs),
		Applications:          r.mapApplicationsFromSchema(d),
		BoundaryScope:         restapi.BoundaryScope(d.Get(ApplicationAlertConfigFieldBoundaryScope).(string)),
		CustomerPayloadFields: r.mapCustomPayloadFieldsFromSchema(d),
		Description:           d.Get(ApplicationAlertConfigFieldDescription).(string),
		EvaluationType:        restapi.ApplicationAlertEvaluationType(d.Get(ApplicationAlertConfigFieldEvaluationType).(string)),
		Granularity:           restapi.Granularity(d.Get(ApplicationAlertConfigFieldGranularity).(int)),
		IncludeInternal:       d.Get(ApplicationAlertConfigFieldIncludeInternal).(bool),
		IncludeSynthetic:      d.Get(ApplicationAlertConfigFieldIncludeSynthetic).(bool),
		Name:                  d.Get(ApplicationAlertConfigFieldName).(string),
		Rule:                  r.mapRuleFromSchema(d),
		Severity:              severity,
		TagFilterExpression:   tagFilter,
		Threshold:             *threshold,
		TimeThreshold:         r.mapTimeThresholdFromSchema(d),
		Triggering:            d.Get(ApplicationAlertConfigFieldTriggering).(bool),
	}, nil
}

func (r *applicationAlertConfigResource) mapApplicationsFromSchema(d *schema.ResourceData) map[string]restapi.IncludedApplication {
	val := d.Get(ApplicationAlertConfigFieldApplications)
	result := make(map[string]restapi.IncludedApplication)
	if val != nil {
		for _, v := range val.(*schema.Set).List() {
			app := r.mapApplicationFromSchema(v.(map[string]interface{}))
			result[app.ApplicationID] = app
		}
	}
	return result
}

func (r *applicationAlertConfigResource) mapApplicationFromSchema(appData map[string]interface{}) restapi.IncludedApplication {
	services := make(map[string]restapi.IncludedService)
	if appData[ApplicationAlertConfigFieldApplicationsServices] != nil {
		for _, v := range appData[ApplicationAlertConfigFieldApplicationsServices].(*schema.Set).List() {
			service := r.mapServiceFromSchema(v.(map[string]interface{}))
			services[service.ServiceID] = service
		}
	}
	return restapi.IncludedApplication{
		ApplicationID: appData[ApplicationAlertConfigFieldApplicationsApplicationID].(string),
		Inclusive:     appData[ApplicationAlertConfigFieldApplicationsInclusive].(bool),
		Services:      services,
	}
}

func (r *applicationAlertConfigResource) mapServiceFromSchema(appData map[string]interface{}) restapi.IncludedService {
	endpoints := make(map[string]restapi.IncludedEndpoint)
	if appData[ApplicationAlertConfigFieldApplicationsServicesEndpoints] != nil {
		for _, v := range appData[ApplicationAlertConfigFieldApplicationsServicesEndpoints].(*schema.Set).List() {
			endpoint := r.mapEndpointFromSchema(v.(map[string]interface{}))
			endpoints[endpoint.EndpointID] = endpoint
		}
	}
	return restapi.IncludedService{
		ServiceID: appData[ApplicationAlertConfigFieldApplicationsServicesServiceID].(string),
		Inclusive: appData[ApplicationAlertConfigFieldApplicationsInclusive].(bool),
		Endpoints: endpoints,
	}
}

func (r *applicationAlertConfigResource) mapEndpointFromSchema(appData map[string]interface{}) restapi.IncludedEndpoint {
	return restapi.IncludedEndpoint{
		EndpointID: appData[ApplicationAlertConfigFieldApplicationsServicesEndpointsEndpointID].(string),
		Inclusive:  appData[ApplicationAlertConfigFieldApplicationsInclusive].(bool),
	}
}

func (r *applicationAlertConfigResource) mapCustomPayloadFieldsFromSchema(d *schema.ResourceData) []restapi.CustomPayloadField[any] {
	val := d.Get(ApplicationAlertConfigFieldCustomPayloadFields)
	if val != nil {
		fields := val.(*schema.Set).List()
		result := make([]restapi.CustomPayloadField[any], len(fields))
		for i, v := range fields {
			field := v.(map[string]interface{})
			customPayloadFieldType := restapi.StaticCustomPayloadType
			key := field[ApplicationAlertConfigFieldCustomPayloadFieldsKey].(string)
			value := field[ApplicationAlertConfigFieldCustomPayloadFieldsValue].(string)

			/*
				if customPayloadFieldType == restapi.DynamicCustomPayloadType {
					parts := strings.Split(value, "=")
					dynamicValue := restapi.DynamicCustomPayloadFieldValue{TagName: parts[0]}
					if len(parts) == 2 {
						valueKey := parts[1]
						dynamicValue.Key = &valueKey
					}
					result[i] = restapi.CustomPayloadField[any]{
						Type:  customPayloadFieldType,
						Key:   key,
						Value: dynamicValue,
					}
				} else {

			*/
			result[i] = restapi.CustomPayloadField[any]{
				Type:  customPayloadFieldType,
				Key:   key,
				Value: restapi.StaticStringCustomPayloadFieldValue(value),
			}
			//}
		}
		return result
	}
	return []restapi.CustomPayloadField[any]{}
}

func (r *applicationAlertConfigResource) mapRuleFromSchema(d *schema.ResourceData) restapi.ApplicationAlertRule {
	ruleSlice := d.Get(ApplicationAlertConfigFieldRule).([]interface{})
	rule := ruleSlice[0].(map[string]interface{})
	for alertType, v := range rule {
		configSlice := v.([]interface{})
		if len(configSlice) == 1 {
			config := configSlice[0].(map[string]interface{})
			return r.mapRuleConfigFromSchema(config, alertType)
		}
	}
	return restapi.ApplicationAlertRule{}
}

func (r *applicationAlertConfigResource) mapRuleConfigFromSchema(config map[string]interface{}, alertType string) restapi.ApplicationAlertRule {
	var levelPtr *restapi.LogLevel
	if levelString, ok := config[ApplicationAlertConfigFieldRuleLogsLevel]; ok {
		level := restapi.LogLevel(levelString.(string))
		levelPtr = &level
	}
	var stableHashPtr *int32
	if v, ok := config[ApplicationAlertConfigFieldRuleStableHash]; ok {
		stableHash := int32(v.(int))
		stableHashPtr = &stableHash
	}
	var statusCodeStartPtr *int32
	if v, ok := config[ApplicationAlertConfigFieldRuleStatusCodeStart]; ok {
		statusCodeStart := int32(v.(int))
		statusCodeStartPtr = &statusCodeStart
	}
	var statusCodeEndPtr *int32
	if v, ok := config[ApplicationAlertConfigFieldRuleStatusCodeEnd]; ok {
		statusCodeEnd := int32(v.(int))
		statusCodeEndPtr = &statusCodeEnd
	}
	var messagePtr *string
	if v, ok := config[ApplicationAlertConfigFieldRuleLogsMessage]; ok {
		message := v.(string)
		messagePtr = &message
	}
	var operatorPtr *restapi.ExpressionOperator
	if v, ok := config[ApplicationAlertConfigFieldRuleLogsOperator]; ok {
		operator := restapi.ExpressionOperator(v.(string))
		operatorPtr = &operator
	}
	return restapi.ApplicationAlertRule{
		AlertType:       r.mapAlertTypeFromSchema(alertType),
		MetricName:      config[ApplicationAlertConfigFieldRuleMetricName].(string),
		Aggregation:     restapi.Aggregation(config[ApplicationAlertConfigFieldRuleAggregation].(string)),
		StableHash:      stableHashPtr,
		StatusCodeStart: statusCodeStartPtr,
		StatusCodeEnd:   statusCodeEndPtr,
		Level:           levelPtr,
		Message:         messagePtr,
		Operator:        operatorPtr,
	}
}

func (r *applicationAlertConfigResource) mapAlertTypeFromSchema(alertType string) string {
	if alertType == ApplicationAlertConfigFieldRuleErrorRate {
		return "errorRate"
	} else if alertType == ApplicationAlertConfigFieldRuleStatusCode {
		return "statusCode"
	}
	return alertType
}

func (r *applicationAlertConfigResource) mapTagFilterExpressionFromSchema(input string) (restapi.TagFilterExpressionElement, error) {
	parser := tagfilter.NewParser()
	expr, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}

	mapper := tagfilter.NewMapper()
	return mapper.ToAPIModel(expr), nil
}

func (r *applicationAlertConfigResource) mapTimeThresholdFromSchema(d *schema.ResourceData) restapi.TimeThreshold {
	timeThresholdSlice := d.Get(ApplicationAlertConfigFieldTimeThreshold).([]interface{})
	timeThreshold := timeThresholdSlice[0].(map[string]interface{})
	for timeThresholdType, v := range timeThreshold {
		configSlice := v.([]interface{})
		if len(configSlice) == 1 {
			config := configSlice[0].(map[string]interface{})
			var violationsPtr *int32
			if v, ok := config[ApplicationAlertConfigFieldTimeThresholdViolationsInPeriodViolations]; ok {
				violations := int32(v.(int))
				violationsPtr = &violations
			}
			var requestsPtr *int32
			if v, ok := config[ApplicationAlertConfigFieldTimeThresholdRequestImpactRequests]; ok {
				requests := int32(v.(int))
				requestsPtr = &requests
			}
			return restapi.TimeThreshold{
				Type:       r.mapTimeThresholdTypeFromSchema(timeThresholdType),
				TimeWindow: int64(config[ApplicationAlertConfigFieldTimeThresholdTimeWindow].(int)),
				Violations: violationsPtr,
				Requests:   requestsPtr,
			}
		}
	}
	return restapi.TimeThreshold{}
}

func (r *applicationAlertConfigResource) mapTimeThresholdTypeFromSchema(input string) string {
	if input == ApplicationAlertConfigFieldTimeThresholdRequestImpact {
		return "requestImpact"
	} else if input == ApplicationAlertConfigFieldTimeThresholdViolationsInPeriod {
		return "violationsInPeriod"
	} else if input == ApplicationAlertConfigFieldTimeThresholdViolationsInSequence {
		return "violationsInSequence"
	}
	return input
}

func (r *applicationAlertConfigResource) stateUpgradeV0(_ context.Context, state map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	if _, ok := state[ApplicationAlertConfigFieldFullName]; ok {
		state[ApplicationAlertConfigFieldName] = state[ApplicationAlertConfigFieldFullName]
		delete(state, ApplicationAlertConfigFieldFullName)
	}
	return state, nil
}

func (r *applicationAlertConfigResource) schemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			ApplicationAlertConfigFieldAlertChannelIDs:     applicationAlertConfigSchemaAlertChannelIDs,
			ApplicationAlertConfigFieldApplications:        applicationAlertConfigSchemaApplications,
			ApplicationAlertConfigFieldBoundaryScope:       applicationAlertConfigSchemaBoundaryScope,
			ApplicationAlertConfigFieldCustomPayloadFields: applicationAlertConfigSchemaCustomPayloadFields,
			ApplicationAlertConfigFieldDescription:         applicationAlertConfigSchemaDescription,
			ApplicationAlertConfigFieldEvaluationType:      applicationAlertConfigSchemaEvaluationType,
			ApplicationAlertConfigFieldGranularity:         applicationAlertConfigSchemaGranularity,
			ApplicationAlertConfigFieldIncludeInternal:     applicationAlertConfigSchemaIncludeInternal,
			ApplicationAlertConfigFieldIncludeSynthetic:    applicationAlertConfigSchemaIncludeSynthetic,
			ApplicationAlertConfigFieldName:                applicationAlertConfigSchemaName,
			ApplicationAlertConfigFieldFullName:            applicationAlertConfigSchemaFullName,
			ApplicationAlertConfigFieldRule:                applicationAlertConfigSchemaRule,
			ApplicationAlertConfigFieldSeverity:            applicationAlertConfigSchemaSeverity,
			ApplicationAlertConfigFieldTagFilter:           applicationAlertConfigSchemaTagFilter,
			ResourceFieldThreshold:                         thresholdSchema,
			ApplicationAlertConfigFieldTimeThreshold:       applicationAlertConfigSchemaTimeThreshold,
			ApplicationAlertConfigFieldTriggering:          applicationAlertConfigSchemaTriggering,
		},
	}
}
