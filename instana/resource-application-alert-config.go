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
	applicationAlertConfigFieldRuleMetricName  = "metric_name"
	applicationAlertConfigFieldRuleAggregation = "aggregation"
	applicationAlertConfigFieldRuleOperator    = "operator"

	applicationAlertConfigFieldThresholdOperator    = "operator"
	applicationAlertConfigFieldThresholdLastUpdated = "last_updated"

	applicationAlertConfigFieldTimeThresholdTimeWindow = "time_window"
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
				"alerting_channel_ids": {
					Type:     schema.TypeSet,
					MinItems: 0,
					MaxItems: 1024,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Optional:    true,
					Description: "List of IDs of alert channels defined in Instana.",
				},
				"applications": {
					Type:     schema.TypeMap,
					Required: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"applicationId": {
								Type:         schema.TypeString,
								Required:     true,
								Description:  "ID of the included application",
								ValidateFunc: validation.StringLenBetween(0, 64),
							},
							"inclusive": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     false,
								Description: "Defines whether this node and his child nodes are included (true) or excluded (false)",
							},
							"services": {
								Type:     schema.TypeMap,
								Required: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"serviceId": {
											Type:         schema.TypeString,
											Required:     true,
											Description:  "ID of the included service",
											ValidateFunc: validation.StringLenBetween(0, 64),
										},
										"inclusive": {
											Type:        schema.TypeBool,
											Optional:    true,
											Default:     false,
											Description: "Defines whether this node and his child nodes are included (true) or excluded (false)",
										},
										"endpoints": {
											Type:     schema.TypeMap,
											Required: true,
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													"endpointId": {
														Type:         schema.TypeString,
														Required:     true,
														Description:  "ID of the included endpoint",
														ValidateFunc: validation.StringLenBetween(0, 64),
													},
													"inclusive": {
														Type:        schema.TypeBool,
														Optional:    true,
														Default:     false,
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
				"boundary_scope": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(restapi.SupportedApplicationAlertConfigBoundaryScopes.ToStringSlice(), false),
					Description:  "The boundary scope of the application alert config",
				},
				"custom_payload_fields": {
					Type:     schema.TypeList,
					Optional: true,
					MinItems: 0,
					MaxItems: 20,
					Elem: schema.Resource{
						Schema: map[string]*schema.Schema{
							"key": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The key of the custom payload field",
							},
							"value": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The value of the custom payload field",
							},
						},
					},
					Description: "An optional list of custom payload fields (static key/value pairs added to the event)",
				},
				"description": {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "The description text of the application alert config",
					ValidateFunc: validation.StringLenBetween(0, 65536),
				},
				"evaluation_type": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(restapi.SupportedApplicationAlertEvaluationTypes.ToStringSlice(), false),
					Description:  "The evaluation type of the application alert config",
				},
				"granularity": {
					Type:         schema.TypeInt,
					Optional:     true,
					Default:      restapi.Granularity600000,
					ValidateFunc: validation.IntInSlice(restapi.SupportedGranularities.ToIntSlice()),
					Description:  "The evaluation granularity used for detection of violations of the defined threshold. In other words, it defines the size of the tumbling window used",
				},
				"include_internal": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Optional flag to indicate whether also internal calls are included in the scope or not. The default is false",
				},
				"include_synthetic": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Optional flag to indicate whether also synthetic calls are included in the scope or not. The default is false",
				},
				"name": {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "Name for the application alert configuration",
					ValidateFunc: validation.StringLenBetween(0, 256),
				},
				"rule": {
					Type:        schema.TypeList,
					MinItems:    1,
					MaxItems:    1,
					Required:    true,
					Description: "Indicates the type of rule this alert configuration is about.",
					Elem: schema.Resource{
						Schema: map[string]*schema.Schema{
							"error_rate": {
								Type:        schema.TypeList,
								MinItems:    0,
								MaxItems:    1,
								Optional:    true,
								Description: "Rule based on the error rate of the configured alert configuration target",
								Elem: schema.Resource{
									Schema: map[string]*schema.Schema{
										applicationAlertConfigFieldRuleMetricName:  applicationAlertSchemaRuleMetricName,
										applicationAlertConfigFieldRuleAggregation: applicationAlertSchemaOptionalRuleAggregation,
									},
								},
								ExactlyOneOf: applicationAlertRuleTypeKeys,
							},
							"logs": {
								Type:        schema.TypeList,
								MinItems:    0,
								MaxItems:    1,
								Optional:    true,
								Description: "Rule based on logs of the configured alert configuration target",
								Elem: schema.Resource{
									Schema: map[string]*schema.Schema{
										applicationAlertConfigFieldRuleMetricName:  applicationAlertSchemaRuleMetricName,
										applicationAlertConfigFieldRuleAggregation: applicationAlertSchemaOptionalRuleAggregation,
										"level": {
											Type:         schema.TypeString,
											Required:     true,
											Description:  "The log level for which this rule applies to",
											ValidateFunc: validation.StringInSlice(restapi.SupportedLogLevels.ToStringSlice(), true),
										},
										"message": {
											Type:        schema.TypeString,
											Optional:    true,
											Description: "The log message for which this rule applies to",
										},
										applicationAlertConfigFieldRuleOperator: applicationAlertSchemaRequiredRuleOperator,
									},
								},
								ExactlyOneOf: applicationAlertRuleTypeKeys,
							},
							"slowness": {
								Type:        schema.TypeList,
								MinItems:    0,
								MaxItems:    1,
								Optional:    true,
								Description: "Rule based on the slowness of the configured alert configuration target",
								Elem: schema.Resource{
									Schema: map[string]*schema.Schema{
										applicationAlertConfigFieldRuleMetricName:  applicationAlertSchemaRuleMetricName,
										applicationAlertConfigFieldRuleAggregation: applicationAlertSchemaRequiredRuleAggregation,
									},
								},
								ExactlyOneOf: applicationAlertRuleTypeKeys,
							},
							"status_code": {
								Type:        schema.TypeList,
								MinItems:    0,
								MaxItems:    1,
								Optional:    true,
								Description: "Rule based on the HTTP status code of the configured alert configuration target",
								Elem: schema.Resource{
									Schema: map[string]*schema.Schema{
										applicationAlertConfigFieldRuleMetricName:  applicationAlertSchemaRuleMetricName,
										applicationAlertConfigFieldRuleAggregation: applicationAlertSchemaOptionalRuleAggregation,
										"status_code_start": {
											Type:         schema.TypeInt,
											Optional:     true,
											Description:  "minimal HTTP status code applied for this rule",
											ValidateFunc: validation.IntAtLeast(1),
										},
										"status_code_end": {
											Type:         schema.TypeInt,
											Optional:     true,
											Description:  "maximum HTTP status code applied for this rule",
											ValidateFunc: validation.IntAtLeast(1),
										},
									},
								},
								ExactlyOneOf: applicationAlertRuleTypeKeys,
							},
							"throughput": {
								Type:        schema.TypeList,
								MinItems:    0,
								MaxItems:    1,
								Optional:    true,
								Description: "Rule based on the throughput of the configured alert configuration target",
								Elem: schema.Resource{
									Schema: map[string]*schema.Schema{
										applicationAlertConfigFieldRuleMetricName:  applicationAlertSchemaRuleMetricName,
										applicationAlertConfigFieldRuleAggregation: applicationAlertSchemaOptionalRuleAggregation,
									},
								},
								ExactlyOneOf: applicationAlertRuleTypeKeys,
							},
						},
					},
				},
				"severity": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(restapi.SupportedSeverities.TerraformRepresentations(), false),
					Description:  "The severity of the alert when triggered",
				},
				"tag_filter": {
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

				"threshold": {
					Type:        schema.TypeList,
					MinItems:    1,
					MaxItems:    1,
					Required:    true,
					Description: "Indicates the type of threshold this alert rule is evaluated on.",
					Elem: schema.Resource{
						Schema: map[string]*schema.Schema{
							"historic_baseline": {
								Type:        schema.TypeList,
								MinItems:    0,
								MaxItems:    1,
								Optional:    true,
								Description: "Threshold based on a historic baseline.",
								Elem: schema.Resource{
									Schema: map[string]*schema.Schema{
										applicationAlertConfigFieldThresholdOperator:    applicationAlertSchemaRequiredThresholdOperator,
										applicationAlertConfigFieldThresholdLastUpdated: applicationAlertSchemaOptionalThresholdLastUpdated,
										"baseline": {
											Type:     schema.TypeList,
											Optional: true,
											Elem: schema.Schema{
												Type:     schema.TypeList,
												Optional: true,
												Elem: schema.Schema{
													Type: schema.TypeFloat,
												},
											},
											Description: "The baseline of the historic baseline threshold",
										},
										"deviationFactor": {
											Type:         schema.TypeFloat,
											Optional:     true,
											ValidateFunc: validation.FloatBetween(0.5, 16),
											Description:  "The baseline of the historic baseline threshold",
										},
										"seasonality": {
											Type:         schema.TypeString,
											Required:     true,
											ValidateFunc: validation.StringInSlice(restapi.SupportedThresholdSeasonalities.ToStringSlice(), true),
											Description:  "The seasonality of the historic baseline threshold",
										},
									},
								},
								ExactlyOneOf: applicationAlertThresholdTypeKeys,
							},
							"static": {
								Type:        schema.TypeList,
								MinItems:    0,
								MaxItems:    1,
								Optional:    true,
								Description: "Static threshold definition",
								Elem: schema.Resource{
									Schema: map[string]*schema.Schema{
										applicationAlertConfigFieldThresholdOperator:    applicationAlertSchemaRequiredThresholdOperator,
										applicationAlertConfigFieldThresholdLastUpdated: applicationAlertSchemaOptionalThresholdLastUpdated,
										"value": {
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

				"time_threshold": {
					Type:        schema.TypeList,
					MinItems:    1,
					MaxItems:    1,
					Required:    true,
					Description: "Indicates the type of violation of the defined threshold.",
					Elem: schema.Resource{
						Schema: map[string]*schema.Schema{
							"request_impact": {
								Type:        schema.TypeList,
								MinItems:    0,
								MaxItems:    1,
								Optional:    true,
								Description: "Time threshold base on request impact",
								Elem: schema.Resource{
									Schema: map[string]*schema.Schema{
										applicationAlertConfigFieldTimeThresholdTimeWindow: applicationAlertSchemaRequiredTimeThresholdTimeWindow,
									},
								},
								ExactlyOneOf: applicationAlertTimeThresholdTypeKeys,
							},
							"violations_in_period": {
								Type:        schema.TypeList,
								MinItems:    0,
								MaxItems:    1,
								Optional:    true,
								Description: "Time threshold base on violations in period",
								Elem: schema.Resource{
									Schema: map[string]*schema.Schema{
										applicationAlertConfigFieldTimeThresholdTimeWindow: applicationAlertSchemaRequiredTimeThresholdTimeWindow,
										"violations": {
											Type:         schema.TypeInt,
											Optional:     true,
											ValidateFunc: validation.IntBetween(1, 12),
											Description:  "The violations appeared in the period",
										},
									},
								},
								ExactlyOneOf: applicationAlertTimeThresholdTypeKeys,
							},
							"violations_in_sequence": {
								Type:        schema.TypeList,
								MinItems:    0,
								MaxItems:    1,
								Optional:    true,
								Description: "Time threshold base on violations in sequence",
								Elem: schema.Resource{
									Schema: map[string]*schema.Schema{
										applicationAlertConfigFieldTimeThresholdTimeWindow: applicationAlertSchemaRequiredTimeThresholdTimeWindow,
									},
								},
								ExactlyOneOf: applicationAlertTimeThresholdTypeKeys,
							},
						},
					},
				},

				"triggering": {
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
	d.Set("alerting_channel_ids", config.AlertChannelIDs)
	return nil
}

func (r *applicationAlertConfigResource) MapStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	return nil, nil
}
