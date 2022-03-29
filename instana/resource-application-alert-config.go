package instana

import (
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

//ResourceInstanaApplicationAlertConfig the name of the terraform-provider-instana resource to manage application alert config
const ResourceInstanaApplicationAlertConfig = "instana_application_alert_config"

const (
	//ApplicationAlertConfigFieldAlertChannelIDs constant value for field alerting_channel_ids of resource instana_application_alert_config
	ApplicationAlertConfigFieldAlertChannelIDs = "alerting_channel_ids"
	//ApplicationAlertConfigFieldApplications constant value for field applications of resource instana_application_alert_config
	ApplicationAlertConfigFieldApplications = "applications"
	//ApplicationAlertConfigFieldApplicationsApplicationID constant value for field applications.application_id of resource instana_application_alert_config
	ApplicationAlertConfigFieldApplicationsApplicationID = "application_id"
	//ApplicationAlertConfigFieldApplicationsInclusive constant value for field applications.inclusive of resource instana_application_alert_config
	ApplicationAlertConfigFieldApplicationsInclusive = "inclusive"
	//ApplicationAlertConfigFieldApplicationsServices constant value for field applications.services of resource instana_application_alert_config
	ApplicationAlertConfigFieldApplicationsServices = "services"
	//ApplicationAlertConfigFieldApplicationsServicesServiceID constant value for field applications.services.service_id of resource instana_application_alert_config
	ApplicationAlertConfigFieldApplicationsServicesServiceID = "service_id"
	//ApplicationAlertConfigFieldApplicationsServicesEndpoints constant value for field applications.services.endpoints of resource instana_application_alert_config
	ApplicationAlertConfigFieldApplicationsServicesEndpoints = "endpoints"
	//ApplicationAlertConfigFieldApplicationsServicesEndpointsEndpointID constant value for field applications.services.endpoints.endpoint_id of resource instana_application_alert_config
	ApplicationAlertConfigFieldApplicationsServicesEndpointsEndpointID = "endpoint_id"
	//ApplicationAlertConfigFieldBoundaryScope constant value for field boundary_scope of resource instana_application_alert_config
	ApplicationAlertConfigFieldBoundaryScope = "boundary_scope"
	//ApplicationAlertConfigFieldCustomPayloadFields constant value for field custom_payload_fields of resource instana_application_alert_config
	ApplicationAlertConfigFieldCustomPayloadFields = "custom_payload_fields"
	//ApplicationAlertConfigFieldCustomPayloadFieldsType constant value for field custom_payload_fields.type of resource instana_application_alert_config
	ApplicationAlertConfigFieldCustomPayloadFieldsType = "type"
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
	//ApplicationAlertConfigFieldThreshold constant value for field threshold of resource instana_application_alert_config
	ApplicationAlertConfigFieldThreshold = "threshold"
	//ApplicationAlertConfigFieldThresholdLastUpdated constant value for field threshold.*.last_updated of resource instana_application_alert_config
	ApplicationAlertConfigFieldThresholdLastUpdated = "last_updated"
	//ApplicationAlertConfigFieldThresholdOperator constant value for field threshold.*.operator of resource instana_application_alert_config
	ApplicationAlertConfigFieldThresholdOperator = "operator"
	//ApplicationAlertConfigFieldThresholdHistoricBaseline constant value for field threshold.historic_baseline of resource instana_application_alert_config
	ApplicationAlertConfigFieldThresholdHistoricBaseline = "historic_baseline"
	//ApplicationAlertConfigFieldThresholdHistoricBaselineBaseline constant value for field threshold.historic_baseline.baseline of resource instana_application_alert_config
	ApplicationAlertConfigFieldThresholdHistoricBaselineBaseline = "baseline"
	//ApplicationAlertConfigFieldThresholdHistoricBaselineDeviationFactor constant value for field threshold.historic_baseline.deviation_factor of resource instana_application_alert_config
	ApplicationAlertConfigFieldThresholdHistoricBaselineDeviationFactor = "deviation_factor"
	//ApplicationAlertConfigFieldThresholdHistoricBaselineSeasonality constant value for field threshold.historic_baseline.seasonality of resource instana_application_alert_config
	ApplicationAlertConfigFieldThresholdHistoricBaselineSeasonality = "seasonality"
	//ApplicationAlertConfigFieldThresholdStatic constant value for field threshold.static of resource instana_application_alert_config
	ApplicationAlertConfigFieldThresholdStatic = "static"
	//ApplicationAlertConfigFieldThresholdStaticValue constant value for field threshold.static.value of resource instana_application_alert_config
	ApplicationAlertConfigFieldThresholdStaticValue = "value"
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

	applicationAlertThresholdTypeKeys = []string{
		"threshold.0.historic_baseline",
		"threshold.0.static",
	}

	applicationAlertSchemaRequiredThresholdOperator = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		Description:  "The operator which will be applied to evaluate the threshold",
		ValidateFunc: validation.StringInSlice(restapi.SupportedThresholdOperators.ToStringSlice(), true),
	}

	applicationAlertSchemaOptionalThresholdLastUpdated = &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		ValidateFunc: validation.IntAtLeast(0),
		Description:  "The last updated value of the threshold",
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

//NewApplicationAlertConfigResourceHandle creates a new instance of the ResourceHandle for application alert configs
func NewApplicationAlertConfigResourceHandle() ResourceHandle {
	return &applicationAlertConfigResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaApplicationAlertConfig,
			Schema: map[string]*schema.Schema{
				ApplicationAlertConfigFieldAlertChannelIDs: {
					Type:     schema.TypeSet,
					MinItems: 0,
					MaxItems: 1024,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Optional:    true,
					Description: "List of IDs of alert channels defined in Instana.",
				},
				ApplicationAlertConfigFieldApplications: {
					Type:     schema.TypeList,
					Required: true,
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
								Description: "Defines whether this node and his child nodes are included (true) or excluded (false)",
							},
							ApplicationAlertConfigFieldApplicationsServices: {
								Type:     schema.TypeList,
								Required: true,
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
											Description: "Defines whether this node and his child nodes are included (true) or excluded (false)",
										},
										ApplicationAlertConfigFieldApplicationsServicesEndpoints: {
											Type:     schema.TypeList,
											Required: true,
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
														Description: "Defines whether this node and his child nodes are included (true) or excluded (false)",
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
				},
				ApplicationAlertConfigFieldBoundaryScope: {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(restapi.SupportedApplicationAlertConfigBoundaryScopes.ToStringSlice(), false),
					Description:  "The boundary scope of the application alert config",
				},
				ApplicationAlertConfigFieldCustomPayloadFields: {
					Type:     schema.TypeList,
					Optional: true,
					MinItems: 0,
					MaxItems: 20,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							ApplicationAlertConfigFieldCustomPayloadFieldsType: {
								Type:         schema.TypeString,
								Required:     true,
								ValidateFunc: validation.StringInSlice(restapi.SupportedCustomPayloadTypes.ToStringSlice(), false),
								Description:  "The type of the custom payload field",
							},
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
				},
				ApplicationAlertConfigFieldDescription: {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "The description text of the application alert config",
					ValidateFunc: validation.StringLenBetween(0, 65536),
				},
				ApplicationAlertConfigFieldEvaluationType: {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(restapi.SupportedApplicationAlertEvaluationTypes.ToStringSlice(), false),
					Description:  "The evaluation type of the application alert config",
				},
				ApplicationAlertConfigFieldGranularity: {
					Type:         schema.TypeInt,
					Optional:     true,
					Default:      restapi.Granularity600000,
					ValidateFunc: validation.IntInSlice(restapi.SupportedGranularities.ToIntSlice()),
					Description:  "The evaluation granularity used for detection of violations of the defined threshold. In other words, it defines the size of the tumbling window used",
				},
				ApplicationAlertConfigFieldIncludeInternal: {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Optional flag to indicate whether also internal calls are included in the scope or not. The default is false",
				},
				ApplicationAlertConfigFieldIncludeSynthetic: {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Optional flag to indicate whether also synthetic calls are included in the scope or not. The default is false",
				},
				ApplicationAlertConfigFieldName: {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "Name for the application alert configuration",
					ValidateFunc: validation.StringLenBetween(0, 256),
				},
				ApplicationAlertConfigFieldFullName: {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The full name field of the application alert config. The field is computed and contains the name which is sent to Instana. The computation depends on the configured default_name_prefix and default_name_suffix at provider level",
				},
				ApplicationAlertConfigFieldRule: {
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
				},
				ApplicationAlertConfigFieldSeverity: {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(restapi.SupportedSeverities.TerraformRepresentations(), false),
					Description:  "The severity of the alert when triggered",
				},
				ApplicationAlertConfigFieldTagFilter: {
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
				},
				ApplicationAlertConfigFieldThreshold: {
					Type:        schema.TypeList,
					MinItems:    1,
					MaxItems:    1,
					Required:    true,
					Description: "Indicates the type of threshold this alert rule is evaluated on.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							ApplicationAlertConfigFieldThresholdHistoricBaseline: {
								Type:        schema.TypeList,
								MinItems:    0,
								MaxItems:    1,
								Optional:    true,
								Description: "Threshold based on a historic baseline.",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										ApplicationAlertConfigFieldThresholdOperator:    applicationAlertSchemaRequiredThresholdOperator,
										ApplicationAlertConfigFieldThresholdLastUpdated: applicationAlertSchemaOptionalThresholdLastUpdated,
										ApplicationAlertConfigFieldThresholdHistoricBaselineBaseline: {
											Type:     schema.TypeSet,
											Optional: true,
											Elem: &schema.Schema{
												Type:     schema.TypeSet,
												Optional: false,
												Elem: &schema.Schema{
													Type: schema.TypeFloat,
												},
											},
											Description: "The baseline of the historic baseline threshold",
										},
										ApplicationAlertConfigFieldThresholdHistoricBaselineDeviationFactor: {
											Type:         schema.TypeFloat,
											Optional:     true,
											ValidateFunc: validation.FloatBetween(0.5, 16),
											Description:  "The baseline of the historic baseline threshold",
										},
										ApplicationAlertConfigFieldThresholdHistoricBaselineSeasonality: {
											Type:         schema.TypeString,
											Required:     true,
											ValidateFunc: validation.StringInSlice(restapi.SupportedThresholdSeasonalities.ToStringSlice(), false),
											Description:  "The seasonality of the historic baseline threshold",
										},
									},
								},
								ExactlyOneOf: applicationAlertThresholdTypeKeys,
							},
							ApplicationAlertConfigFieldThresholdStatic: {
								Type:        schema.TypeList,
								MinItems:    0,
								MaxItems:    1,
								Optional:    true,
								Description: "Static threshold definition",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										ApplicationAlertConfigFieldThresholdOperator:    applicationAlertSchemaRequiredThresholdOperator,
										ApplicationAlertConfigFieldThresholdLastUpdated: applicationAlertSchemaOptionalThresholdLastUpdated,
										ApplicationAlertConfigFieldThresholdStaticValue: {
											Type:        schema.TypeFloat,
											Optional:    true,
											Description: "The value of the static threshold",
										},
									},
								},
								ExactlyOneOf: applicationAlertThresholdTypeKeys,
							},
						},
					},
				},
				ApplicationAlertConfigFieldTimeThreshold: {
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
				},
				ApplicationAlertConfigFieldTriggering: {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Optional flag to indicate whether also an Incident is triggered or not. The default is false",
				},
			},
		},
	}
}

type applicationAlertConfigResource struct {
	metaData ResourceMetaData
}

func (r *applicationAlertConfigResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *applicationAlertConfigResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{}
}

func (r *applicationAlertConfigResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource {
	return api.ApplicationAlertConfigs()
}

func (r *applicationAlertConfigResource) SetComputedFields(d *schema.ResourceData) {
	//No computed fields defined
}

func (r *applicationAlertConfigResource) UpdateState(d *schema.ResourceData, obj restapi.InstanaDataObject, formatter utils.ResourceNameFormatter) error {
	config := obj.(*restapi.ApplicationAlertConfig)

	name := formatter.UndoFormat(config.Name)
	severity, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(config.Severity)
	if err != nil {
		return err
	}
	var normalizedTagFilterString *string
	if config.TagFilterExpression != nil {
		normalizedTagFilterString, err = r.mapTagFilterExpressionToSchema(config.TagFilterExpression.(restapi.TagFilterExpressionElement))
		if err != nil {
			return err
		}
	}

	d.Set(ApplicationAlertConfigFieldAlertChannelIDs, config.AlertChannelIDs)
	d.Set(ApplicationAlertConfigFieldApplications, r.mapApplicationsToSchema(config))
	d.Set(ApplicationAlertConfigFieldBoundaryScope, config.BoundaryScope)
	d.Set(ApplicationAlertConfigFieldCustomPayloadFields, r.mapCustomPayloadFieldsToSchema(config))
	d.Set(ApplicationAlertConfigFieldDescription, config.Description)
	d.Set(ApplicationAlertConfigFieldEvaluationType, config.EvaluationType)
	d.Set(ApplicationAlertConfigFieldGranularity, config.Granularity)
	d.Set(ApplicationAlertConfigFieldIncludeInternal, config.IncludeInternal)
	d.Set(ApplicationAlertConfigFieldIncludeSynthetic, config.IncludeSynthetic)
	d.Set(ApplicationAlertConfigFieldName, name)
	d.Set(ApplicationAlertConfigFieldFullName, config.Name)
	d.Set(ApplicationAlertConfigFieldRule, r.mapRuleToSchema(config))
	d.Set(ApplicationAlertConfigFieldSeverity, severity)
	d.Set(ApplicationAlertConfigFieldTagFilter, normalizedTagFilterString)
	d.Set(ApplicationAlertConfigFieldThreshold, r.mapThresholdToSchema(config))
	d.Set(ApplicationAlertConfigFieldTimeThreshold, r.mapTimeThresholdToSchema(config))
	d.Set(ApplicationAlertConfigFieldTriggering, config.Triggering)
	d.SetId(config.ID)
	return nil
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
		field[ApplicationAlertConfigFieldCustomPayloadFieldsType] = string(v.Type)
		field[ApplicationAlertConfigFieldCustomPayloadFieldsKey] = v.Key
		field[ApplicationAlertConfigFieldCustomPayloadFieldsValue] = v.Value
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

func (r *applicationAlertConfigResource) mapTagFilterExpressionToSchema(input restapi.TagFilterExpressionElement) (*string, error) {
	mapper := tagfilter.NewMapper()
	expr, err := mapper.FromAPIModel(input)
	if err != nil {
		return nil, err
	}
	renderedExpression := expr.Render()
	return &renderedExpression, nil
}

func (r *applicationAlertConfigResource) mapThresholdToSchema(config *restapi.ApplicationAlertConfig) []map[string]interface{} {
	thresholdConfig := make(map[string]interface{})
	thresholdConfig[ApplicationAlertConfigFieldThresholdOperator] = config.Threshold.Operator
	thresholdConfig[ApplicationAlertConfigFieldThresholdLastUpdated] = config.Threshold.LastUpdated

	if config.Threshold.Value != nil {
		thresholdConfig[ApplicationAlertConfigFieldThresholdStaticValue] = *config.Threshold.Value
	}
	if config.Threshold.Baseline != nil {
		thresholdConfig[ApplicationAlertConfigFieldThresholdHistoricBaselineBaseline] = *config.Threshold.Baseline
	}
	if config.Threshold.DeviationFactor != nil {
		thresholdConfig[ApplicationAlertConfigFieldThresholdHistoricBaselineDeviationFactor] = float64(*config.Threshold.DeviationFactor)
	}
	if config.Threshold.Seasonality != nil {
		thresholdConfig[ApplicationAlertConfigFieldThresholdHistoricBaselineSeasonality] = *config.Threshold.Seasonality
	}

	thresholdType := r.mapThresholdTypeToSchema(config.Threshold.Type)
	threshold := make(map[string]interface{})
	threshold[thresholdType] = []interface{}{thresholdConfig}
	result := make([]map[string]interface{}, 1)
	result[0] = threshold
	return result
}

func (r *applicationAlertConfigResource) mapThresholdTypeToSchema(input string) string {
	if input == "historicBaseline" {
		return ApplicationAlertConfigFieldThresholdHistoricBaseline
	} else if input == "staticThreshold" {
		return ApplicationAlertConfigFieldThresholdStatic
	}
	return input
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

func (r *applicationAlertConfigResource) MapStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	fullName := r.mapFullNameStringFromSchema(d, formatter)
	severity, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(d.Get(ApplicationAlertConfigFieldSeverity).(string))
	if err != nil {
		return nil, err
	}

	var tagFilter restapi.TagFilterExpressionElement
	tagFilterStr, ok := d.GetOk(ApplicationAlertConfigFieldTagFilter)
	if ok {
		tagFilter, err = r.mapTagFilterExpressionFromSchema(tagFilterStr.(string))
		if err != nil {
			return &restapi.ApplicationConfig{}, err
		}
	}

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
		Name:                  fullName,
		Rule:                  r.mapRuleFromSchema(d),
		Severity:              severity,
		TagFilterExpression:   tagFilter,
		Threshold:             r.mapThresholdFromSchema(d),
		TimeThreshold:         r.mapTimeThresholdFromSchema(d),
		Triggering:            d.Get(ApplicationAlertConfigFieldTriggering).(bool),
	}, nil
}

func (r *applicationAlertConfigResource) mapFullNameStringFromSchema(d *schema.ResourceData, formatter utils.ResourceNameFormatter) string {
	if d.HasChange(ApplicationAlertConfigFieldName) {
		return formatter.Format(d.Get(ApplicationAlertConfigFieldName).(string))
	}
	return d.Get(ApplicationAlertConfigFieldFullName).(string)
}

func (r *applicationAlertConfigResource) mapApplicationsFromSchema(d *schema.ResourceData) map[string]restapi.IncludedApplication {
	val := d.Get(ApplicationAlertConfigFieldApplications)
	result := make(map[string]restapi.IncludedApplication)
	if val != nil {
		for _, v := range val.([]interface{}) {
			app := r.mapApplicationFromSchema(v.(map[string]interface{}))
			result[app.ApplicationID] = app
		}
	}
	return result
}

func (r *applicationAlertConfigResource) mapApplicationFromSchema(appData map[string]interface{}) restapi.IncludedApplication {
	services := make(map[string]restapi.IncludedService)
	if appData[ApplicationAlertConfigFieldApplicationsServices] != nil {
		for _, v := range appData[ApplicationAlertConfigFieldApplicationsServices].([]interface{}) {
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
		for _, v := range appData[ApplicationAlertConfigFieldApplicationsServicesEndpoints].([]interface{}) {
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

func (r *applicationAlertConfigResource) mapCustomPayloadFieldsFromSchema(d *schema.ResourceData) []restapi.CustomPayloadField {
	val := d.Get(ApplicationAlertConfigFieldCustomPayloadFields)
	if val != nil {
		fields := val.([]interface{})
		result := make([]restapi.CustomPayloadField, len(fields))
		for i, v := range fields {
			field := v.(map[string]interface{})
			result[i] = restapi.CustomPayloadField{
				Type:  restapi.CustomPayloadType(field[ApplicationAlertConfigFieldCustomPayloadFieldsType].(string)),
				Key:   field[ApplicationAlertConfigFieldCustomPayloadFieldsKey].(string),
				Value: field[ApplicationAlertConfigFieldCustomPayloadFieldsValue].(string),
			}
		}
		return result
	}
	return []restapi.CustomPayloadField{}
}

func (r *applicationAlertConfigResource) mapRuleFromSchema(d *schema.ResourceData) restapi.ApplicationAlertRule {
	ruleSlice := d.Get(ApplicationAlertConfigFieldRule).([]interface{})
	rule := ruleSlice[0].(map[string]interface{})
	for alertType, v := range rule {
		configSlice := v.([]interface{})
		if len(configSlice) == 1 {
			config := configSlice[0].(map[string]interface{})
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
	}
	return restapi.ApplicationAlertRule{}
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

func (r *applicationAlertConfigResource) mapThresholdFromSchema(d *schema.ResourceData) restapi.Threshold {
	thresholdSlice := d.Get(ApplicationAlertConfigFieldThreshold).([]interface{})
	threshold := thresholdSlice[0].(map[string]interface{})
	for thresholdType, v := range threshold {
		configSlice := v.([]interface{})
		if len(configSlice) == 1 {
			config := configSlice[0].(map[string]interface{})
			var seasonalityPtr *restapi.ThresholdSeasonality
			if v, ok := config[ApplicationAlertConfigFieldThresholdHistoricBaselineSeasonality]; ok {
				seasonality := restapi.ThresholdSeasonality(v.(string))
				seasonalityPtr = &seasonality
			}
			var lastUpdatePtr *int64
			if v, ok := config[ApplicationAlertConfigFieldThresholdLastUpdated]; ok {
				lastUpdate := int64(v.(int))
				lastUpdatePtr = &lastUpdate
			}
			var valuePtr *float64
			if v, ok := config[ApplicationAlertConfigFieldThresholdStaticValue]; ok {
				value := v.(float64)
				valuePtr = &value
			}
			var deviationFactorPtr *float32
			if v, ok := config[ApplicationAlertConfigFieldThresholdHistoricBaselineDeviationFactor]; ok {
				deviationFactor := float32(v.(float64))
				deviationFactorPtr = &deviationFactor
			}
			var baselinePtr *[][]float64
			if v, ok := config[ApplicationAlertConfigFieldThresholdHistoricBaselineBaseline]; ok {
				baselineSet := v.(*schema.Set)
				if baselineSet.Len() > 0 {
					baseline := make([][]float64, baselineSet.Len())
					for i, val := range baselineSet.List() {
						baseline[i] = ConvertInterfaceSlice[float64](val.(*schema.Set).List())
					}
					baselinePtr = &baseline
				}
			}
			return restapi.Threshold{
				Type:            r.mapThresholdTypeFromSchema(thresholdType),
				Operator:        restapi.ThresholdOperator(config[ApplicationAlertConfigFieldThresholdOperator].(string)),
				LastUpdated:     lastUpdatePtr,
				Value:           valuePtr,
				DeviationFactor: deviationFactorPtr,
				Baseline:        baselinePtr,
				Seasonality:     seasonalityPtr,
			}
		}
	}
	return restapi.Threshold{}
}

func (r *applicationAlertConfigResource) mapThresholdTypeFromSchema(input string) string {
	if input == ApplicationAlertConfigFieldThresholdHistoricBaseline {
		return "historicBaseline"
	} else if input == ApplicationAlertConfigFieldThresholdStatic {
		return "staticThreshold"
	}
	return input
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
